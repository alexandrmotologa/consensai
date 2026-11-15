package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alexandrmotologa/consensai/internal/protocol"
)

// SkepticAgent represents the adversarial risk counsel persona in the deliberation.
type SkepticAgent struct {
	client ProviderClient
}

func NewSkepticAgent(client ProviderClient) *SkepticAgent {
	return &SkepticAgent{client: client}
}

const skepticSystemPrompt = `You are The Skeptic, an adversarial risk auditor and fraud examiner.
Your mandate is to identify failure modes, covenant breaches, fraudulent indicators, hidden debt, and downside risks in the submitted dossier.
Operational rules:
1. Challenge unverified assumptions and identify stress-test vulnerabilities.
2. Ground each identified risk with citations from the raw dossier.
3. Formulate explicit, verifiable risk claims with citations.
4. Output valid JSON adhering to the specified schema.`

func (s *SkepticAgent) ExecuteRound1(ctx context.Context, dossierRaw string) (*protocol.AgentTurn, error) {
	userPrompt := fmt.Sprintf(`### Round 1: Blind Framing (Antithesis)
Evaluate the following dossier from an adversarial risk perspective.
Identify 2 to 5 failure scenarios, covenant vulnerabilities, or fraud indicators. Include specific citations.

Dossier:
%s

Output must be a valid JSON object matching:
{
  "summary": "Brief 1-2 sentence executive summary of risk assessment",
  "raw_content": "Detailed adversarial failure case analysis",
  "claims": [
    {
      "id": "claim_skp_1",
      "statement": "Clear factual risk assertion",
      "category": "risk_exposure|anomaly|covenant_failure|compliance",
      "citation": "Exact source field or text quote"
    }
  ]
}`, dossierRaw)

	resp, err := s.client.GenerateCompletion(ctx, CompletionRequest{
		SystemPrompt: skepticSystemPrompt,
		UserPrompt:   userPrompt,
		Temperature:  0.1,
		MaxTokens:    2048,
	})
	if err != nil {
		return nil, fmt.Errorf("skeptic round 1 error: %w", err)
	}

	return parseTurnResponse(resp, protocol.RoleSkeptic, 1)
}

func (s *SkepticAgent) ExecuteRound2(ctx context.Context, dossierRaw string, advocateClaims []protocol.Claim) (*protocol.AgentTurn, error) {
	advocateJSON, _ := json.MarshalIndent(advocateClaims, "", "  ")

	userPrompt := fmt.Sprintf(`### Round 2: Targeted Cross-Examination
The Advocate has submitted the following affirmative claims:
%s

Dossier:
%s

Your mandate:
1. Cross-examine the Advocate's claims against downside market stress and tail risks.
2. If the Advocate's evidence genuinely resolves an objection, record a formal concession.
3. Highlight remaining unmitigated risks.

Output must be a valid JSON object matching:
{
  "summary": "Summary of cross-examination and concessions",
  "raw_content": "Detailed risk rebuttal analysis",
  "rebuttals": [
    {
      "target_claim_id": "claim_adv_1",
      "counter_point": "Detailed critique of affirmative assertion",
      "citation": "Source citation",
      "concession": false
    }
  ],
  "concessions": [
    "Specific conceded point description"
  ]
}`, string(advocateJSON), dossierRaw)

	resp, err := s.client.GenerateCompletion(ctx, CompletionRequest{
		SystemPrompt: skepticSystemPrompt,
		UserPrompt:   userPrompt,
		Temperature:  0.1,
		MaxTokens:    2048,
	})
	if err != nil {
		return nil, fmt.Errorf("skeptic round 2 error: %w", err)
	}

	return parseTurnResponse(resp, protocol.RoleSkeptic, 2)
}
