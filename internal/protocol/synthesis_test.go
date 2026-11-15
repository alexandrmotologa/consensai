package protocol

import (
	"testing"

	"github.com/alexandrmotologa/consensai/internal/config"
)

func TestSynthesizer_ApprovalOutcome(t *testing.T) {
	thresholds := config.ThresholdConfig{
		Approve: 0.70,
		Reject:  0.35,
	}
	synth := NewSynthesizer(thresholds)

	claims := []Claim{
		{ID: "c1", AgentRole: RoleAdvocate, AuditedStatus: StatusVerified},
		{ID: "c2", AgentRole: RoleAdvocate, AuditedStatus: StatusVerified},
		{ID: "c3", AgentRole: RoleAdvocate, AuditedStatus: StatusVerified},
		{ID: "c4", AgentRole: RoleSkeptic, AuditedStatus: StatusConceded},
	}

	findings := []AuditFinding{
		{ClaimID: "c1", Status: StatusVerified},
		{ClaimID: "c2", Status: StatusVerified},
		{ClaimID: "c3", Status: StatusVerified},
		{ClaimID: "c4", Status: StatusConceded},
	}

	res := synth.ComputeEpistemicConsensus(claims, findings, nil)

	if res.Decision != OutcomeApprove {
		t.Fatalf("expected decision %s, got %s (score: %.2f)", OutcomeApprove, res.Decision, res.NormalizedScore)
	}

	if res.NormalizedScore < 0.70 {
		t.Fatalf("expected normalized score >= 0.70, got %.2f", res.NormalizedScore)
	}

	if len(res.PreservedDissents) != 0 {
		t.Fatalf("expected 0 preserved dissents, got %d", len(res.PreservedDissents))
	}
}

func TestSynthesizer_RejectOutcome(t *testing.T) {
	thresholds := config.ThresholdConfig{
		Approve: 0.70,
		Reject:  0.35,
	}
	synth := NewSynthesizer(thresholds)

	claims := []Claim{
		{ID: "c1", AgentRole: RoleAdvocate, AuditedStatus: StatusUnverified},
		{ID: "c2", AgentRole: RoleSkeptic, AuditedStatus: StatusUnmitigated},
		{ID: "c3", AgentRole: RoleSkeptic, AuditedStatus: StatusUnmitigated},
	}

	findings := []AuditFinding{
		{ClaimID: "c1", Status: StatusUnverified},
		{ClaimID: "c2", Status: StatusUnmitigated},
		{ClaimID: "c3", Status: StatusUnmitigated},
	}

	res := synth.ComputeEpistemicConsensus(claims, findings, nil)

	if res.Decision != OutcomeReject {
		t.Fatalf("expected decision %s, got %s (score: %.2f)", OutcomeReject, res.Decision, res.NormalizedScore)
	}

	if res.NormalizedScore > 0.35 {
		t.Fatalf("expected normalized score <= 0.35, got %.2f", res.NormalizedScore)
	}

	if len(res.PreservedDissents) != 2 {
		t.Fatalf("expected 2 preserved dissents, got %d", len(res.PreservedDissents))
	}
}

func TestSynthesizer_EscalateOutcome(t *testing.T) {
	thresholds := config.ThresholdConfig{
		Approve: 0.70,
		Reject:  0.35,
	}
	synth := NewSynthesizer(thresholds)

	claims := []Claim{
		{ID: "c1", AgentRole: RoleAdvocate, AuditedStatus: StatusVerified},
		{ID: "c2", AgentRole: RoleAdvocate, AuditedStatus: StatusVerified},
		{ID: "c3", AgentRole: RoleSkeptic, AuditedStatus: StatusUnmitigated},
	}

	findings := []AuditFinding{
		{ClaimID: "c1", Status: StatusVerified},
		{ClaimID: "c2", Status: StatusVerified},
		{ClaimID: "c3", Status: StatusUnmitigated},
	}

	res := synth.ComputeEpistemicConsensus(claims, findings, nil)

	if res.Decision != OutcomeEscalate {
		t.Fatalf("expected decision %s, got %s (score: %.2f)", OutcomeEscalate, res.Decision, res.NormalizedScore)
	}
}
