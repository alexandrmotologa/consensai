package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/alexandrmotologa/consensai/internal/protocol"
)

// AuditorAgent represents the judicial evidence clerk persona in the deliberation.
type AuditorAgent struct {
	client ProviderClient
}

func NewAuditorAgent(client ProviderClient) *AuditorAgent {
	return &AuditorAgent{client: client}
}

const auditorSystemPrompt = `You are The Risk Auditor, an impartial judicial referee and forensic evidence verifier.
Your mandate is to verify all assertions made by The Advocate and The Skeptic against the raw input dossier.
Operational rules:
1. Examine each claim and its cited evidence against the raw dossier.
2. If an assertion is fully supported by the dossier, mark it VERIFIED.
3. If an assertion is contradicted or absent from the dossier, mark it UNVERIFIED or HALLUCINATED.
4. If an assertion was formally conceded by opposing counsel, mark it CONCEDED.
5. If an identified adversarial risk remains unmitigated, mark it UNMITIGATED.
6. Assign a calibrated confidence score (0.0 to 1.0) to each audited claim.
7. Output valid JSON adhering to the specified schema.`

func (a *AuditorAgent) ExecuteRound3(
	ctx context.Context,
	dossierRaw string,
	allClaims []protocol.Claim,
	turns []protocol.AgentTurn,
) (*protocol.AgentTurn, []protocol.AuditFinding, error) {
	claimsJSON, _ := json.MarshalIndent(allClaims, "", "  ")

	userPrompt := fmt.Sprintf(`### Round 3: Factual Grounding & Citation Audit
Evaluate the following complete set of claims and rebuttals against the source dossier.

All Claims from Rounds 1 & 2:
%s

Raw Source Dossier:
%s

Your mandate:
1. Cross-reference every cited field and statement against the raw text.
2. Determine if the claim is VERIFIED, UNVERIFIED, HALLUCINATED, CONCEDED, or UNMITIGATED.
3. Report any factual discrepancy or mathematical variance.

Output must be a valid JSON object matching:
{
  "summary": "Executive audit summary of factual grounding and remaining unmitigated risks",
  "raw_content": "Detailed audit findings and citation reconciliation",
  "audit_findings": [
    {
      "claim_id": "claim_adv_1",
      "agent_role": "Advocate",
      "statement": "Original claim statement",
      "citation_given": "Citation text or path",
      "found_in_dossier": true,
      "discrepancy": "None or specific discrepancy",
      "status": "VERIFIED|UNVERIFIED|HALLUCINATED|CONCEDED|UNMITIGATED",
      "confidence_score": 0.95
    }
  ]
}`, string(claimsJSON), dossierRaw)

	resp, err := a.client.GenerateCompletion(ctx, CompletionRequest{
		SystemPrompt: auditorSystemPrompt,
		UserPrompt:   userPrompt,
		Temperature:  0.0,
		MaxTokens:    3000,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("auditor round 3 error: %w", err)
	}

	turn, findings, err := parseAuditorResponse(resp, 3)
	if err != nil {
		// Fallback to deterministic verification if response parsing fails
		findings = FallbackGroundingAudit(dossierRaw, allClaims)
		turn = &protocol.AgentTurn{
			Round:      3,
			Role:       protocol.RoleAuditor,
			Timestamp:  time.Now().UTC(),
			Summary:    "Automated deterministic citation audit completed.",
			RawContent: "Audit verified claims using lexical matching against raw dossier text.",
		}
	}

	return turn, findings, nil
}

// FallbackGroundingAudit performs direct substring and lexical verification against the dossier.
func FallbackGroundingAudit(dossierRaw string, claims []protocol.Claim) []protocol.AuditFinding {
	dossierLower := strings.ToLower(dossierRaw)
	var findings []AuditFindingWrapper

	for _, c := range claims {
		citationClean := strings.TrimSpace(strings.ToLower(c.Citation))
		found := false
		if citationClean != "" && len(citationClean) > 3 {
			// Check if citation or parts of it appear in dossier
			parts := strings.FieldsFunc(citationClean, func(r rune) bool {
				return r == '=' || r == ';' || r == ':' || r == ','
			})
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if len(part) > 3 && strings.Contains(dossierLower, part) {
					found = true
					break
				}
			}
		}

		status := protocol.StatusUnverified
		confidence := 0.60
		discrepancy := "Citation could not be verified by exact lexical match"

		if found {
			if c.AgentRole == protocol.RoleAdvocate {
				status = protocol.StatusVerified
				confidence = 0.95
				discrepancy = "Verified against raw dossier text"
			} else {
				status = protocol.StatusUnmitigated
				confidence = 0.90
				discrepancy = "Adversarial risk substantiated by source text"
			}
		}

		findings = append(findings, AuditFindingWrapper{
			ClaimID:         c.ID,
			AgentRole:       c.AgentRole,
			Statement:       c.Statement,
			CitationGiven:   c.Citation,
			FoundInDossier:  found,
			Discrepancy:     discrepancy,
			Status:          status,
			ConfidenceScore: confidence,
		})
	}

	res := make([]protocol.AuditFinding, len(findings))
	for i, f := range findings {
		res[i] = protocol.AuditFinding{
			ClaimID:         f.ClaimID,
			AgentRole:       f.AgentRole,
			Statement:       f.Statement,
			CitationGiven:   f.CitationGiven,
			FoundInDossier:  f.FoundInDossier,
			Discrepancy:     f.Discrepancy,
			Status:          f.Status,
			ConfidenceScore: f.ConfidenceScore,
		}
	}
	return res
}

type AuditFindingWrapper struct {
	ClaimID         string               `json:"claim_id"`
	AgentRole       protocol.AgentRole   `json:"agent_role"`
	Statement       string               `json:"statement"`
	CitationGiven   string               `json:"citation_given"`
	FoundInDossier  bool                 `json:"found_in_dossier"`
	Discrepancy     string               `json:"discrepancy"`
	Status          protocol.ClaimStatus `json:"status"`
	ConfidenceScore float64              `json:"confidence_score"`
}

func parseTurnResponse(raw string, role protocol.AgentRole, round int) (*protocol.AgentTurn, error) {
	cleaned := cleanJSON(raw)

	var payload struct {
		Summary     string              `json:"summary"`
		RawContent  string              `json:"raw_content"`
		Claims      []protocol.Claim    `json:"claims"`
		Rebuttals   []protocol.Rebuttal `json:"rebuttals"`
		Concessions []string            `json:"concessions"`
	}

	if err := json.Unmarshal([]byte(cleaned), &payload); err != nil {
		// If raw unmarshal fails, create structured turn from text
		return &protocol.AgentTurn{
			Round:      round,
			Role:       role,
			Timestamp:  time.Now().UTC(),
			Summary:    fmt.Sprintf("%s completed turn for round %d", role, round),
			RawContent: raw,
		}, nil
	}

	for i := range payload.Claims {
		payload.Claims[i].AgentRole = role
		payload.Claims[i].AuditedStatus = protocol.StatusPending
	}

	return &protocol.AgentTurn{
		Round:       round,
		Role:        role,
		Timestamp:   time.Now().UTC(),
		Summary:     payload.Summary,
		RawContent:  payload.RawContent,
		Claims:      payload.Claims,
		Rebuttals:   payload.Rebuttals,
		Concessions: payload.Concessions,
	}, nil
}

func parseAuditorResponse(raw string, round int) (*protocol.AgentTurn, []protocol.AuditFinding, error) {
	cleaned := cleanJSON(raw)

	var payload struct {
		Summary       string                  `json:"summary"`
		RawContent    string                  `json:"raw_content"`
		AuditFindings []protocol.AuditFinding `json:"audit_findings"`
	}

	if err := json.Unmarshal([]byte(cleaned), &payload); err != nil {
		return nil, nil, err
	}

	turn := &protocol.AgentTurn{
		Round:      round,
		Role:       protocol.RoleAuditor,
		Timestamp:  time.Now().UTC(),
		Summary:    payload.Summary,
		RawContent: payload.RawContent,
	}

	return turn, payload.AuditFindings, nil
}

func cleanJSON(s string) string {
	s = strings.TrimSpace(s)
	if idx := strings.Index(s, "```json"); idx != -1 {
		s = s[idx+7:]
		if end := strings.Index(s, "```"); end != -1 {
			s = s[:end]
		}
	} else if idx := strings.Index(s, "```"); idx != -1 {
		s = s[idx+3:]
		if end := strings.Index(s, "```"); end != -1 {
			s = s[:end]
		}
	}
	return strings.TrimSpace(s)
}
