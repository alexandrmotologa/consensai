package agent

import (
	"context"
	"strings"
	"testing"
)

func TestMockProvider_DossierReasoning(t *testing.T) {
	mock := NewMockProvider()
	ctx := context.Background()

	// Test Advocate Round 1
	advResp, err := mock.GenerateCompletion(ctx, CompletionRequest{
		SystemPrompt: advocateSystemPrompt,
		UserPrompt:   "### Round 1: Blind Framing (Thesis)\nEvaluate commercial_loan_500k.json",
	})
	if err != nil {
		t.Fatalf("advocate round 1 failed: %v", err)
	}
	if !strings.Contains(advResp, "claim_adv_1") {
		t.Fatalf("expected claim_adv_1 in advocate response, got: %s", advResp)
	}

	// Test Skeptic Round 1
	skpResp, err := mock.GenerateCompletion(ctx, CompletionRequest{
		SystemPrompt: skepticSystemPrompt,
		UserPrompt:   "### Round 1: Blind Framing (Antithesis)\nEvaluate commercial_loan_500k.json",
	})
	if err != nil {
		t.Fatalf("skeptic round 1 failed: %v", err)
	}
	if !strings.Contains(skpResp, "claim_skp_1") {
		t.Fatalf("expected claim_skp_1 in skeptic response, got: %s", skpResp)
	}

	// Test Auditor Round 3
	audResp, err := mock.GenerateCompletion(ctx, CompletionRequest{
		SystemPrompt: auditorSystemPrompt,
		UserPrompt:   "### Round 3: Factual Grounding\nEvaluate commercial_loan_500k.json",
	})
	if err != nil {
		t.Fatalf("auditor round 3 failed: %v", err)
	}
	if !strings.Contains(audResp, "audit_findings") {
		t.Fatalf("expected audit_findings in auditor response, got: %s", audResp)
	}
}
