package protocol

import (
	"fmt"
	"math"

	"github.com/alexandrmotologa/consensai/internal/config"
)

// Synthesizer calculates the deterministic consensus output and Epistemic Convergence Score.
type Synthesizer struct {
	thresholds config.ThresholdConfig
}

func NewSynthesizer(thresholds config.ThresholdConfig) *Synthesizer {
	if thresholds.Approve == 0 {
		thresholds.Approve = 0.70
	}
	if thresholds.Reject == 0 {
		thresholds.Reject = 0.35
	}
	return &Synthesizer{thresholds: thresholds}
}

// ComputeEpistemicConsensus computes the mathematical convergence metric and decision.
func (s *Synthesizer) ComputeEpistemicConsensus(
	claims []Claim,
	findings []AuditFinding,
	turns []AgentTurn,
) *SynthesisResult {
	// Map findings by ClaimID for fast lookup
	findingMap := make(map[string]AuditFinding)
	for _, f := range findings {
		findingMap[f.ClaimID] = f
	}

	var advocateClaims []Claim
	var skepticClaims []Claim
	var verifiedAdvocateCount int
	var unmitigatedSkepticCount int
	var preservedDissents []string

	for _, c := range claims {
		finding, hasFinding := findingMap[c.ID]
		status := c.AuditedStatus
		if hasFinding {
			status = finding.Status
		}

		if c.AgentRole == RoleAdvocate {
			advocateClaims = append(advocateClaims, c)
			if status == StatusVerified {
				verifiedAdvocateCount++
			}
		} else if c.AgentRole == RoleSkeptic {
			skepticClaims = append(skepticClaims, c)
			if status == StatusUnmitigated || status == StatusPending || (!hasFinding && status != StatusConceded) {
				unmitigatedSkepticCount++
				dissent := fmt.Sprintf("[%s] %s (Citation: %s)", c.Category, c.Statement, c.Citation)
				if hasFinding && finding.Discrepancy != "" {
					dissent += fmt.Sprintf(" - Finding: %s", finding.Discrepancy)
				}
				preservedDissents = append(preservedDissents, dissent)
			}
		}
	}

	totalClaims := len(advocateClaims) + len(skepticClaims)
	var rawConvergence float64
	var normalized float64

	if totalClaims > 0 {
		// C_epistemic = (|A ∩ V| - |S_unmitigated|) / (|A| + |S|)
		rawConvergence = float64(verifiedAdvocateCount-unmitigatedSkepticCount) / float64(totalClaims)
		// Normalized to 0.0 .. 1.0 range: (raw + 1.0) / 2.0
		normalized = (rawConvergence + 1.0) / 2.0
		if normalized < 0.0 {
			normalized = 0.0
		} else if normalized > 1.0 {
			normalized = 1.0
		}
	} else {
		rawConvergence = 0.0
		normalized = 0.50
	}

	// Calculate Standard Error & 95% Confidence Interval
	var ci [2]float64
	if totalClaims > 0 {
		variance := (normalized * (1.0 - normalized)) / float64(totalClaims)
		se := math.Sqrt(math.Max(0.0, variance))
		z := 1.96
		ci[0] = math.Max(0.0, normalized-z*se)
		ci[1] = math.Min(1.0, normalized+z*se)
	} else {
		ci = [2]float64{0.0, 1.0}
	}

	// Evaluate against decision thresholds
	var decision DecisionOutcome
	var rationale string

	switch {
	case normalized >= s.thresholds.Approve:
		decision = OutcomeApprove
		rationale = fmt.Sprintf(
			"Normalized Epistemic Convergence Score (%.2f) exceeds approval threshold (%.2f). Affirmative claims demonstrate high documentary grounding (%d/%d verified) with manageable mitigated risks.",
			normalized, s.thresholds.Approve, verifiedAdvocateCount, len(advocateClaims),
		)

	case normalized <= s.thresholds.Reject:
		decision = OutcomeReject
		rationale = fmt.Sprintf(
			"Normalized Epistemic Convergence Score (%.2f) falls at or below rejection threshold (%.2f). %d unmitigated high-severity risks or failed authentication checks detected.",
			normalized, s.thresholds.Reject, unmitigatedSkepticCount,
		)

	default:
		decision = OutcomeEscalate
		rationale = fmt.Sprintf(
			"Normalized Epistemic Convergence Score (%.2f) falls in the escalation interval (%.2f - %.2f). Deliberation uncovered contested evidentiary points or marginal headroom requiring executive human oversight.",
			normalized, s.thresholds.Reject, s.thresholds.Approve,
		)
	}

	summary := fmt.Sprintf(
		"Protocol synthesized %d claims across %d rounds. Advocate verified claims: %d. Skeptic unmitigated risks: %d. Final Governance Decision: %s.",
		totalClaims, len(turns), verifiedAdvocateCount, unmitigatedSkepticCount, decision,
	)

	return &SynthesisResult{
		Decision:                  decision,
		EpistemicConvergenceScore: math.Round(rawConvergence*1000) / 1000,
		NormalizedScore:           math.Round(normalized*1000) / 1000,
		ConfidenceInterval: [2]float64{
			math.Round(ci[0]*1000) / 1000,
			math.Round(ci[1]*1000) / 1000,
		},
		TotalClaimsAudited:        totalClaims,
		VerifiedAffirmativeClaims: verifiedAdvocateCount,
		UnmitigatedAdversarialRisks: unmitigatedSkepticCount,
		PreservedDissents:         preservedDissents,
		ActionRationale:           rationale,
		ExecutiveSummary:          summary,
	}
}
