package protocol

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/alexandrmotologa/consensai/internal/config"
)

// EngineEvent represents an observable state transition during deliberation.
type EngineEvent struct {
	Type      string      `json:"type"`
	SessionID string      `json:"session_id"`
	Round     int         `json:"round"`
	Role      AgentRole   `json:"role,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// EventHandler receives real-time deliberation events.
type EventHandler func(event EngineEvent)

// AdvocateExecutor defines the interface required from the Advocate agent.
type AdvocateExecutor interface {
	ExecuteRound1(ctx context.Context, dossierRaw string) (*AgentTurn, error)
	ExecuteRound2(ctx context.Context, dossierRaw string, skepticClaims []Claim) (*AgentTurn, error)
}

// SkepticExecutor defines the interface required from the Skeptic agent.
type SkepticExecutor interface {
	ExecuteRound1(ctx context.Context, dossierRaw string) (*AgentTurn, error)
	ExecuteRound2(ctx context.Context, dossierRaw string, advocateClaims []Claim) (*AgentTurn, error)
}

// AuditorExecutor defines the interface required from the Auditor agent.
type AuditorExecutor interface {
	ExecuteRound3(ctx context.Context, dossierRaw string, allClaims []Claim, turns []AgentTurn) (*AgentTurn, []AuditFinding, error)
}

// Engine coordinates the multi-agent dialectical protocol.
type Engine struct {
	cfg         *config.ProtocolConfig
	advocate    AdvocateExecutor
	skeptic     SkepticExecutor
	auditor     AuditorExecutor
	synthesizer *Synthesizer
}

func NewEngine(
	cfg *config.ProtocolConfig,
	advocate AdvocateExecutor,
	skeptic SkepticExecutor,
	auditor AuditorExecutor,
) *Engine {
	return &Engine{
		cfg:         cfg,
		advocate:    advocate,
		skeptic:     skeptic,
		auditor:     auditor,
		synthesizer: NewSynthesizer(cfg.Thresholds),
	}
}

// RunDeliberation executes the full 3-round protocol over the provided dossier.
func (e *Engine) RunDeliberation(
	ctx context.Context,
	dossierID string,
	dossierTitle string,
	dossierRaw string,
	handler EventHandler,
) (*DeliberationState, error) {
	sessionID := generateSessionID()

	state := &DeliberationState{
		ID:           sessionID,
		DossierID:    dossierID,
		DossierTitle: dossierTitle,
		DossierRaw:   dossierRaw,
		CurrentRound: 0,
		CurrentPhase: "INITIALIZING",
		Turns:        make([]AgentTurn, 0, 6),
		AllClaims:    make([]Claim, 0, 12),
		CreatedAt:    time.Now().UTC(),
	}

	emit := func(eventType string, round int, role AgentRole, data interface{}) {
		if handler != nil {
			handler(EngineEvent{
				Type:      eventType,
				SessionID: sessionID,
				Round:     round,
				Role:      role,
				Timestamp: time.Now().UTC(),
				Data:      data,
			})
		}
	}

	emit("session_start", 0, "", map[string]interface{}{
		"session_id":    sessionID,
		"dossier_id":    dossierID,
		"dossier_title": dossierTitle,
	})

	// -------------------------------------------------------------------------
	// Round 1: Blind Framing (Thesis vs Antithesis)
	// -------------------------------------------------------------------------
	state.CurrentRound = 1
	state.CurrentPhase = "ROUND_1_BLIND_FRAMING"
	emit("phase_change", 1, "", "Round 1: Blind Framing (Advocate & Skeptic in isolation)")

	var advTurn1, skpTurn1 *AgentTurn
	var errAdv1, errSkp1 error
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		emit("turn_start", 1, RoleAdvocate, "Advocate evaluating affirmative case")
		advTurn1, errAdv1 = e.advocate.ExecuteRound1(ctx, dossierRaw)
	}()

	go func() {
		defer wg.Done()
		emit("turn_start", 1, RoleSkeptic, "Skeptic cross-examining risk factors")
		skpTurn1, errSkp1 = e.skeptic.ExecuteRound1(ctx, dossierRaw)
	}()

	wg.Wait()

	if errAdv1 != nil {
		return nil, fmt.Errorf("round 1 advocate failed: %w", errAdv1)
	}
	if errSkp1 != nil {
		return nil, fmt.Errorf("round 1 skeptic failed: %w", errSkp1)
	}

	state.Turns = append(state.Turns, *advTurn1, *skpTurn1)
	state.AllClaims = append(state.AllClaims, advTurn1.Claims...)
	state.AllClaims = append(state.AllClaims, skpTurn1.Claims...)

	emit("turn_complete", 1, RoleAdvocate, advTurn1)
	emit("turn_complete", 1, RoleSkeptic, skpTurn1)

	// -------------------------------------------------------------------------
	// Round 2: Targeted Cross-Rebuttal
	// -------------------------------------------------------------------------
	state.CurrentRound = 2
	state.CurrentPhase = "ROUND_2_CROSS_REBUTTAL"
	emit("phase_change", 2, "", "Round 2: Targeted Cross-Rebuttal and Concessions")

	var advTurn2, skpTurn2 *AgentTurn
	var errAdv2, errSkp2 error
	wg.Add(2)

	go func() {
		defer wg.Done()
		emit("turn_start", 2, RoleAdvocate, "Advocate formulating rebuttals against Skeptic objections")
		advTurn2, errAdv2 = e.advocate.ExecuteRound2(ctx, dossierRaw, skpTurn1.Claims)
	}()

	go func() {
		defer wg.Done()
		emit("turn_start", 2, RoleSkeptic, "Skeptic evaluating Advocate assertions")
		skpTurn2, errSkp2 = e.skeptic.ExecuteRound2(ctx, dossierRaw, advTurn1.Claims)
	}()

	wg.Wait()

	if errAdv2 != nil {
		return nil, fmt.Errorf("round 2 advocate failed: %w", errAdv2)
	}
	if errSkp2 != nil {
		return nil, fmt.Errorf("round 2 skeptic failed: %w", errSkp2)
	}

	state.Turns = append(state.Turns, *advTurn2, *skpTurn2)
	emit("turn_complete", 2, RoleAdvocate, advTurn2)
	emit("turn_complete", 2, RoleSkeptic, skpTurn2)

	// -------------------------------------------------------------------------
	// Round 3: Factual Grounding & Citation Audit
	// -------------------------------------------------------------------------
	state.CurrentRound = 3
	state.CurrentPhase = "ROUND_3_FACTUAL_AUDIT"
	emit("phase_change", 3, "", "Round 3: Grounded Audit and Citation Verification")

	emit("turn_start", 3, RoleAuditor, "Risk Auditor reconciling citations against raw dossier")
	audTurn3, findings, errAud3 := e.auditor.ExecuteRound3(ctx, dossierRaw, state.AllClaims, state.Turns)
	if errAud3 != nil {
		return nil, fmt.Errorf("round 3 auditor failed: %w", errAud3)
	}

	state.Turns = append(state.Turns, *audTurn3)
	state.AuditFindings = findings

	// Update claim statuses based on findings
	findingMap := make(map[string]AuditFinding)
	for _, f := range findings {
		findingMap[f.ClaimID] = f
	}
	for i := range state.AllClaims {
		if f, exists := findingMap[state.AllClaims[i].ID]; exists {
			state.AllClaims[i].AuditedStatus = f.Status
			state.AllClaims[i].AuditNotes = f.Discrepancy
		}
	}

	emit("turn_complete", 3, RoleAuditor, audTurn3)
	emit("audit_complete", 3, RoleAuditor, findings)

	// -------------------------------------------------------------------------
	// Synthesis: Deterministic Consensus & Score
	// -------------------------------------------------------------------------
	state.CurrentPhase = "SYNTHESIS"
	emit("phase_change", 4, RoleSynthesizer, "Deterministic Epistemic Consensus Synthesis")

	synthesis := e.synthesizer.ComputeEpistemicConsensus(state.AllClaims, state.AuditFindings, state.Turns)
	synthesis.Compliance = EvaluateCompliance(dossierID, dossierRaw, synthesis.Decision, state.AllClaims)
	state.Synthesis = synthesis
	completedAt := time.Now().UTC()
	state.CompletedAt = &completedAt
	state.CurrentPhase = "COMPLETED"

	emit("synthesis_complete", 4, RoleSynthesizer, synthesis)

	return state, nil
}

func generateSessionID() string {
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	return "sess_" + hex.EncodeToString(b)
}
