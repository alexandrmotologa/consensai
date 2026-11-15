package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/alexandrmotologa/consensai/internal/agent"
	"github.com/alexandrmotologa/consensai/internal/config"
	"github.com/alexandrmotologa/consensai/internal/ledger"
	"github.com/alexandrmotologa/consensai/internal/protocol"
	"github.com/alexandrmotologa/consensai/internal/report"
)

// SessionRecord holds the state and cryptographic ledger for a deliberation.
type SessionRecord struct {
	ID        string                      `json:"id"`
	State     *protocol.DeliberationState `json:"state"`
	Ledger    *ledger.Ledger              `json:"ledger"`
	CreatedAt time.Time                   `json:"created_at"`
}

// Server provides the HTTP REST API and Server-Sent Events infrastructure.
type Server struct {
	cfg          *config.Config
	startTime    time.Time
	mu           sync.RWMutex
	sessions     map[string]*SessionRecord
	broadcasters map[string]*Broadcaster
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg:          cfg,
		startTime:    time.Now().UTC(),
		sessions:     make(map[string]*SessionRecord),
		broadcasters: make(map[string]*Broadcaster),
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	// REST API Endpoints
	mux.HandleFunc("GET /api/status", s.handleStatus)
	mux.HandleFunc("GET /api/dossiers", s.handleListDossiers)
	mux.HandleFunc("GET /api/dossiers/{id}", s.handleGetDossier)
	mux.HandleFunc("POST /api/deliberate", s.handleDeliberate)
	mux.HandleFunc("GET /api/deliberate/{id}/stream", s.handleStream)
	mux.HandleFunc("GET /api/deliberate/{id}/ledger", s.handleGetLedger)
	mux.HandleFunc("POST /api/ledger/verify", s.handleVerifyLedger)
	mux.HandleFunc("GET /api/deliberate/{id}/memo", s.handleGetMemo)

	// Static Web Cockpit
	mux.Handle("/", EmbeddedWebHandler())

	return s.corsMiddleware(mux)
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	s.mu.RLock()
	activeSessions := len(s.sessions)
	s.mu.RUnlock()

	status := map[string]interface{}{
		"status":                 "online",
		"version":                "1.0.0",
		"uptime_seconds":         int(time.Since(s.startTime).Seconds()),
		"allocated_memory_bytes": m.Alloc,
		"allocated_memory_mb":    float64(m.Alloc) / (1024 * 1024),
		"sys_memory_mb":          float64(m.Sys) / (1024 * 1024),
		"goroutines":             runtime.NumGoroutine(),
		"active_sessions":        activeSessions,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}

func (s *Server) handleListDossiers(w http.ResponseWriter, r *http.Request) {
	dossiersDir := s.cfg.Server.DossiersDir
	if dossiersDir == "" {
		dossiersDir = "./examples/dossiers"
	}

	type dossierSummary struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		Domain    string `json:"domain"`
		RiskTier  string `json:"risk_tier"`
		FilePath  string `json:"file_path"`
	}

	var results []dossierSummary

	entries, err := os.ReadDir(dossiersDir)
	if err == nil {
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
				continue
			}

			fullPath := filepath.Join(dossiersDir, e.Name())
			data, err := os.ReadFile(fullPath)
			if err != nil {
				continue
			}

			var parsed struct {
				ID       string `json:"id"`
				Title    string `json:"title"`
				Domain   string `json:"domain"`
				RiskTier string `json:"risk_tier"`
			}
			if err := json.Unmarshal(data, &parsed); err == nil {
				if parsed.ID == "" {
					parsed.ID = strings.TrimSuffix(e.Name(), ".json")
				}
				results = append(results, dossierSummary{
					ID:       parsed.ID,
					Title:    parsed.Title,
					Domain:   parsed.Domain,
					RiskTier: parsed.RiskTier,
					FilePath: fullPath,
				})
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(results)
}

func (s *Server) handleGetDossier(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	dossiersDir := s.cfg.Server.DossiersDir
	if dossiersDir == "" {
		dossiersDir = "./examples/dossiers"
	}

	targetPath := filepath.Join(dossiersDir, id+".json")
	data, err := os.ReadFile(targetPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Dossier %s not found: %v", id, err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

type DeliberateRequest struct {
	DossierID        string                 `json:"dossier_id"`
	CustomDossier    map[string]interface{} `json:"custom_dossier,omitempty"`
	Provider         string                 `json:"provider,omitempty"`
	Model            string                 `json:"model,omitempty"`
	ThresholdApprove float64                `json:"threshold_approve,omitempty"`
	ThresholdReject  float64                `json:"threshold_reject,omitempty"`
}

func (s *Server) handleDeliberate(w http.ResponseWriter, r *http.Request) {
	var req DeliberateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	var dossierRaw string
	var dossierTitle string
	var dossierID string

	if len(req.CustomDossier) > 0 {
		rawBytes, err := json.MarshalIndent(req.CustomDossier, "", "  ")
		if err != nil {
			http.Error(w, "Formatting custom dossier: "+err.Error(), http.StatusBadRequest)
			return
		}
		dossierRaw = string(rawBytes)
		dossierID = "custom_dossier"
		if id, ok := req.CustomDossier["id"].(string); ok && id != "" {
			dossierID = id
		}
		dossierTitle = "Custom Dossier Review"
		if title, ok := req.CustomDossier["title"].(string); ok && title != "" {
			dossierTitle = title
		}
	} else if req.DossierID != "" {
		dossiersDir := s.cfg.Server.DossiersDir
		if dossiersDir == "" {
			dossiersDir = "./examples/dossiers"
		}
		targetPath := filepath.Join(dossiersDir, req.DossierID+".json")
		data, err := os.ReadFile(targetPath)
		if err != nil {
			http.Error(w, fmt.Sprintf("Dossier %s not found: %v", req.DossierID, err), http.StatusNotFound)
			return
		}
		dossierRaw = string(data)
		dossierID = req.DossierID

		var parsed struct {
			Title string `json:"title"`
		}
		_ = json.Unmarshal(data, &parsed)
		dossierTitle = parsed.Title
		if dossierTitle == "" {
			dossierTitle = req.DossierID
		}
	} else {
		http.Error(w, "Either dossier_id or custom_dossier is required", http.StatusBadRequest)
		return
	}

	// Prepare provider
	providerCfg := s.cfg.Provider
	if req.Provider != "" {
		providerCfg.Default = req.Provider
	}
	if req.Model != "" {
		providerCfg.Model = req.Model
	}

	client, err := agent.NewProviderClient(&providerCfg)
	if err != nil {
		http.Error(w, "Initializing provider client: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare protocol parameters
	protoCfg := s.cfg.Protocol
	if req.ThresholdApprove > 0 {
		protoCfg.Thresholds.Approve = req.ThresholdApprove
	}
	if req.ThresholdReject > 0 {
		protoCfg.Thresholds.Reject = req.ThresholdReject
	}

	advocate := agent.NewAdvocateAgent(client)
	skeptic := agent.NewSkepticAgent(client)
	auditor := agent.NewAuditorAgent(client)
	engine := protocol.NewEngine(&protoCfg, advocate, skeptic, auditor)

	// Create broadcaster and register under active
	broadcaster := NewBroadcaster()
	s.mu.Lock()
	s.broadcasters["active"] = broadcaster
	s.mu.Unlock()

	// Launch async deliberation
	go func() {
		state, err := engine.RunDeliberation(
			context.Background(),
			dossierID,
			dossierTitle,
			dossierRaw,
			func(evt protocol.EngineEvent) {
				if evt.SessionID != "" {
					s.mu.Lock()
					s.broadcasters[evt.SessionID] = broadcaster
					s.mu.Unlock()
				}
				broadcaster.Broadcast(evt)
			},
		)

		if err != nil {
			broadcaster.Broadcast(protocol.EngineEvent{
				Type:      "error",
				Timestamp: time.Now().UTC(),
				Data:      err.Error(),
			})
			return
		}

		// Commit cryptographic ledger
		ledgerObj, err := ledger.CreateFromState(state)
		if err == nil {
			s.mu.Lock()
			s.sessions[state.ID] = &SessionRecord{
				ID:        state.ID,
				State:     state,
				Ledger:    ledgerObj,
				CreatedAt: time.Now().UTC(),
			}
			s.mu.Unlock()

			// Broadcast committed block
			broadcaster.Broadcast(protocol.EngineEvent{
				Type:      "ledger_committed",
				SessionID: state.ID,
				Timestamp: time.Now().UTC(),
				Data: map[string]interface{}{
					"session_id":  state.ID,
					"block_count": len(ledgerObj.Blocks),
					"head_hash":   ledgerObj.Blocks[len(ledgerObj.Blocks)-1].Hash,
				},
			})
		}
	}()

	// Register broadcaster
	// Wait a tiny moment for sessionID generation or broadcast start
	// In practice we can return immediately with session start hook
	respPayload := map[string]interface{}{
		"status":     "started",
		"message":    "Deliberation started across 3 rounds",
		"dossier_id": dossierID,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(respPayload)
}

func (s *Server) handleStream(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	s.mu.RLock()
	broadcaster, ok := s.broadcasters[id]
	s.mu.RUnlock()

	if !ok {
		s.mu.Lock()
		broadcaster = NewBroadcaster()
		s.broadcasters[id] = broadcaster
		s.mu.Unlock()
	}

	broadcaster.ServeHTTP(w, r)
}

func (s *Server) handleGetLedger(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	s.mu.RLock()
	rec, ok := s.sessions[id]
	s.mu.RUnlock()

	if !ok || rec.Ledger == nil {
		http.Error(w, fmt.Sprintf("Ledger for session %s not found", id), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(rec.Ledger)
}

func (s *Server) handleVerifyLedger(w http.ResponseWriter, r *http.Request) {
	var l ledger.Ledger
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		http.Error(w, "Invalid ledger JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	valid, err := l.Verify()
	headHash := ""
	if len(l.Blocks) > 0 {
		headHash = l.Blocks[len(l.Blocks)-1].Hash
	}

	res := map[string]interface{}{
		"valid":       valid,
		"block_count": len(l.Blocks),
		"head_hash":   headHash,
		"verified_at": time.Now().UTC().Format(time.RFC3339),
	}
	if err != nil {
		res["error"] = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func (s *Server) handleGetMemo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	s.mu.RLock()
	rec, ok := s.sessions[id]
	s.mu.RUnlock()

	if !ok || rec.State == nil {
		http.Error(w, fmt.Sprintf("Session %s not found", id), http.StatusNotFound)
		return
	}

	headHash := "UNCOMMITTED"
	if rec.Ledger != nil && len(rec.Ledger.Blocks) > 0 {
		headHash = rec.Ledger.Blocks[len(rec.Ledger.Blocks)-1].Hash
	}

	memo := report.NewDecisionMemo(rec.State)
	md := memo.GenerateMarkdown(headHash)

	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	_, _ = w.Write([]byte(md))
}
