package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alexandrmotologa/consensai/internal/config"
)

func TestServer_Integration(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Server.DossiersDir = "../../examples/dossiers"

	srv := NewServer(cfg)
	handler := srv.Routes()

	// 1. Test GET /api/status
	reqStatus := httptest.NewRequest("GET", "/api/status", nil)
	recStatus := httptest.NewRecorder()
	handler.ServeHTTP(recStatus, reqStatus)

	if recStatus.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", recStatus.Code, recStatus.Body.String())
	}

	var statusMap map[string]interface{}
	if err := json.Unmarshal(recStatus.Body.Bytes(), &statusMap); err != nil {
		t.Fatalf("unmarshaling status: %v", err)
	}
	if statusMap["status"] != "online" {
		t.Fatalf("expected status 'online', got %v", statusMap["status"])
	}

	// 2. Test GET /api/dossiers
	reqDossiers := httptest.NewRequest("GET", "/api/dossiers", nil)
	recDossiers := httptest.NewRecorder()
	handler.ServeHTTP(recDossiers, reqDossiers)

	if recDossiers.Code != http.StatusOK {
		t.Fatalf("expected dossiers 200, got %d", recDossiers.Code)
	}

	// 3. Test POST /api/deliberate with sample dossier
	delibReqPayload := map[string]interface{}{
		"dossier_id": "commercial_loan_500k",
		"provider":   "mock",
	}
	delibReqBytes, _ := json.Marshal(delibReqPayload)
	reqDelib := httptest.NewRequest("POST", "/api/deliberate", bytes.NewReader(delibReqBytes))
	recDelib := httptest.NewRecorder()
	handler.ServeHTTP(recDelib, reqDelib)

	if recDelib.Code != http.StatusOK {
		t.Fatalf("expected deliberate 200, got %d: %s", recDelib.Code, recDelib.Body.String())
	}

	// Wait for background deliberation to finish and commit
	time.Sleep(350 * time.Millisecond)

	srv.mu.RLock()
	var sessionId string
	for id := range srv.sessions {
		sessionId = id
		break
	}
	srv.mu.RUnlock()

	if sessionId == "" {
		t.Fatalf("expected session to be created")
	}

	// 4. Test POST /api/deliberate/{id}/recalculate (Sensitivity Analysis)
	recalcPayload := map[string]interface{}{
		"threshold_approve": 0.60,
		"threshold_reject":  0.30,
	}
	recalcBytes, _ := json.Marshal(recalcPayload)
	reqRecalc := httptest.NewRequest("POST", "/api/deliberate/"+sessionId+"/recalculate", bytes.NewReader(recalcBytes))
	recRecalc := httptest.NewRecorder()
	handler.ServeHTTP(recRecalc, reqRecalc)

	if recRecalc.Code != http.StatusOK {
		t.Fatalf("expected recalculate 200, got %d: %s", recRecalc.Code, recRecalc.Body.String())
	}

	var recalcRes map[string]interface{}
	if err := json.Unmarshal(recRecalc.Body.Bytes(), &recalcRes); err != nil {
		t.Fatalf("unmarshaling recalculate response: %v", err)
	}
	if recalcRes["decision"] == nil {
		t.Fatalf("expected decision in recalculate result")
	}

	// 5. Test POST /api/deliberate/{id}/adjudicate (Human-in-the-Loop Override)
	adjPayload := map[string]interface{}{
		"officer_id":     "OFFICER-442",
		"officer_name":   "Mark Stone",
		"officer_role":   "Senior Credit Officer",
		"decision":       "CONDITIONAL_APPROVE",
		"justification":  "Approved with $50k escrow cushion.",
		"conditions":     []string{"1.20x quarterly DSCR test"},
	}
	adjBytes, _ := json.Marshal(adjPayload)
	reqAdj := httptest.NewRequest("POST", "/api/deliberate/"+sessionId+"/adjudicate", bytes.NewReader(adjBytes))
	recAdj := httptest.NewRecorder()
	handler.ServeHTTP(recAdj, reqAdj)

	if recAdj.Code != http.StatusOK {
		t.Fatalf("expected adjudicate 200, got %d: %s", recAdj.Code, recAdj.Body.String())
	}

	// 6. Test GET /api/deliberate/{id}/certificate (Printable Certificate)
	reqCert := httptest.NewRequest("GET", "/api/deliberate/"+sessionId+"/certificate", nil)
	recCert := httptest.NewRecorder()
	handler.ServeHTTP(recCert, reqCert)

	if recCert.Code != http.StatusOK {
		t.Fatalf("expected certificate 200, got %d", recCert.Code)
	}
	if !strings.Contains(recCert.Body.String(), "consensai Compliance & Audit Certificate") {
		t.Fatalf("expected certificate HTML to contain certificate title")
	}

	// 7. Test GET / (embedded React UI)
	reqIndex := httptest.NewRequest("GET", "/", nil)
	recIndex := httptest.NewRecorder()
	handler.ServeHTTP(recIndex, reqIndex)

	if recIndex.Code != http.StatusOK {
		t.Fatalf("expected index 200, got %d", recIndex.Code)
	}
	if !strings.Contains(recIndex.Body.String(), "consensai") {
		t.Fatalf("expected index.html to contain 'consensai', got: %s", recIndex.Body.String())
	}
}
