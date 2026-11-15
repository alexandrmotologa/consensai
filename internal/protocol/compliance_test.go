package protocol

import (
	"testing"
)

func TestEvaluateCompliance_CreditUnderwriting(t *testing.T) {
	dossier := `{"id": "commercial_loan", "annual_revenue": 2840000, "dscr": 1.48}`
	rep := EvaluateCompliance("commercial_loan_500k", dossier, OutcomeEscalate, nil)

	if rep == nil {
		t.Fatalf("expected compliance report, got nil")
	}

	foundEU := false
	foundBasel := false

	for _, c := range rep.Checks {
		if c.Standard == "EU_AI_ACT_ARTICLE_14" {
			foundEU = true
			if c.Status != "REQUIRES_HUMAN_ACTION" {
				t.Fatalf("expected EU AI Act check to require human action on escalation, got %s", c.Status)
			}
		}
		if c.Standard == "BASEL_III_CAPITAL_ADEQUACY" {
			foundBasel = true
		}
	}

	if !foundEU || !foundBasel {
		t.Fatalf("expected both EU AI Act and Basel III checks to be present")
	}
}
