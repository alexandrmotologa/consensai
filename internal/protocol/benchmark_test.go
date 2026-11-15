package protocol

import (
	"context"
	"fmt"
	"runtime"
	"testing"

	"github.com/alexandrmotologa/consensai/internal/config"
)

// MockAdvocate for benchmark
type benchAdvocate struct{}

func (b *benchAdvocate) ExecuteRound1(ctx context.Context, dossierRaw string) (*AgentTurn, error) {
	return &AgentTurn{
		Round: 1, Role: RoleAdvocate, Summary: "Affirmative",
		Claims: []Claim{{ID: "c1", AgentRole: RoleAdvocate, Statement: "Statement 1", Citation: "path.1"}},
	}, nil
}
func (b *benchAdvocate) ExecuteRound2(ctx context.Context, dossierRaw string, skepticClaims []Claim) (*AgentTurn, error) {
	return &AgentTurn{Round: 2, Role: RoleAdvocate, Summary: "Rebuttal"}, nil
}

// MockSkeptic for benchmark
type benchSkeptic struct{}

func (b *benchSkeptic) ExecuteRound1(ctx context.Context, dossierRaw string) (*AgentTurn, error) {
	return &AgentTurn{
		Round: 1, Role: RoleSkeptic, Summary: "Antithesis",
		Claims: []Claim{{ID: "s1", AgentRole: RoleSkeptic, Statement: "Risk 1", Citation: "risk.path"}},
	}, nil
}
func (b *benchSkeptic) ExecuteRound2(ctx context.Context, dossierRaw string, advocateClaims []Claim) (*AgentTurn, error) {
	return &AgentTurn{Round: 2, Role: RoleSkeptic, Summary: "Cross"}, nil
}

// MockAuditor for benchmark
type benchAuditor struct{}

func (b *benchAuditor) ExecuteRound3(ctx context.Context, dossierRaw string, allClaims []Claim, turns []AgentTurn) (*AgentTurn, []AuditFinding, error) {
	return &AgentTurn{Round: 3, Role: RoleAuditor, Summary: "Audit complete"},
		[]AuditFinding{{ClaimID: "c1", Status: StatusVerified}, {ClaimID: "s1", Status: StatusUnmitigated}},
		nil
}

func TestMemoryFootprintUnder25MB(t *testing.T) {
	runtime.GC()
	var mBefore runtime.MemStats
	runtime.ReadMemStats(&mBefore)

	cfg := &config.ProtocolConfig{
		MaxRounds: 3,
		Thresholds: config.ThresholdConfig{
			Approve: 0.70,
			Reject:  0.35,
		},
	}

	engine := NewEngine(cfg, &benchAdvocate{}, &benchSkeptic{}, &benchAuditor{})

	// Run 100 consecutive deliberations
	for i := 0; i < 100; i++ {
		dossier := fmt.Sprintf(`{"id":"dossier_%d","data":"sample financial payload"}`, i)
		_, err := engine.RunDeliberation(context.Background(), fmt.Sprintf("dossier_%d", i), "Bench", dossier, nil)
		if err != nil {
			t.Fatalf("deliberation failed: %v", err)
		}
	}

	runtime.GC()
	var mAfter runtime.MemStats
	runtime.ReadMemStats(&mAfter)

	allocMB := float64(mAfter.Alloc) / (1024 * 1024)
	sysMB := float64(mAfter.Sys) / (1024 * 1024)

	t.Logf("Allocated Memory: %.2f MB, System Memory: %.2f MB", allocMB, sysMB)

	// Invariant: Memory footprint must be strictly under 25MB RAM
	if allocMB > 25.0 {
		t.Fatalf("memory allocation exceeded 25MB: %.2f MB", allocMB)
	}
}
