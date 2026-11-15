package agent

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alexandrmotologa/consensai/internal/protocol"
)

// AdvocateAgent represents the affirmative counsel persona in the deliberation.
type AdvocateAgent struct {
	client ProviderClient
}

func NewAdvocateAgent(client ProviderClient) *AdvocateAgent {
	return &AdvocateAgent{client: client}
}

const advocateSystemPrompt = `You are The Advocate, an expert corporate counsel and high-stakes underwriting advocate.
Your mandate is to present the strongest possible affirmative case for approving the submitted dossier while maintaining strict factual honesty.
Operational rules:
1. Ground all claims in the documentary evidence provided in the dossier.
2. Formulate explicit, verifiable claims with exact citations (quotes or JSON field paths).
3. Do not invent or hallucinate facts, numbers, or regulatory waivers.
4. Output valid JSON adhering to the specified schema.`

func (a *AdvocateAgent) ExecuteRound1(ctx context.Context, dossierRaw string) (*protocol.AgentTurn, error) {
	userPrompt := fmt.Sprintf(`### Round 1: Blind Framing (Thesis)
Evaluate the following dossier and establish the affirmative case for approval.
Extract 2 to 5 factual, verifiable claims supporting approval. Include specific citations.

Dossier:
%s

Output must be a valid JSON object matching:
{
  "summary": "Brief 1-2 sentence executive summary of affirmative case",
  "raw_content": "Detailed analytical rationale defending approval",
  "claims": [
    {
      "id": "claim_adv_1",
      "statement": "Clear factual assertion",
      "category": "financial|collateral|compliance|operational",
      "citation": "Exact source field or text quote"
    }
  ]
}`, dossierRaw)

	resp, err := a.client.GenerateCompletion(ctx, CompletionRequest{
		SystemPrompt: advocateSystemPrompt,
		UserPrompt:   userPrompt,
		Temperature:  0.1,
		MaxTokens:    2048,
	})
	if err != nil {
		return nil, fmt.Errorf("advocate round 1 error: %w", err)
	}

	return parseTurnResponse(resp, protocol.RoleAdvocate, 1)
}

func (a *AdvocateAgent) ExecuteRound2(ctx context.Context, dossierRaw string, skepticClaims []protocol.Claim) (*protocol.AgentTurn, error) {
	skepticJSON, _ := json.MarshalIndent(skepticClaims, "", "  ")

	userPrompt := fmt.Sprintf(`### Round 2: Targeted Rebuttal
The Skeptic has raised the following adversarial risks and objections:
%s

Dossier:
%s

Your mandate:
1. Address each objection with documentary evidence from the dossier.
2. If an objection is insurmountable or mathematically valid, formally concede it.
3. Formulate structured rebuttals with citations.

Output must be a valid JSON object matching:
{
  "summary": "Summary of rebuttal defense and conceded items",
  "raw_content": "Detailed counter-argument analysis",
  "rebuttals": [
    {
      "target_claim_id": "claim_skp_1",
      "counter_point": "Detailed rebuttal point",
      "citation": "Source citation",
      "concession": false
    }
  ],
  "concessions": [
    "Specific conceded risk description"
  ]
}`, string(skepticJSON), dossierRaw)

	resp, err := a.client.GenerateCompletion(ctx, CompletionRequest{
		SystemPrompt: advocateSystemPrompt,
		UserPrompt:   userPrompt,
		Temperature:  0.1,
		MaxTokens:    2048,
	})
	if err != nil {
		return nil, fmt.Errorf("advocate round 2 error: %w", err)
	}

	return parseTurnResponse(resp, protocol.RoleAdvocate, 2)
}
