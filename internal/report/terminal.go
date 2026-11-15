package report

import (
	"fmt"
	"strings"

	"github.com/alexandrmotologa/consensai/internal/protocol"
)

// ANSI color codes for terminal styling.
const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorDim    = "\033[2m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

// TerminalRenderer prints structured deliberation events to standard output.
type TerminalRenderer struct{}

func NewTerminalRenderer() *TerminalRenderer {
	return &TerminalRenderer{}
}

func (r *TerminalRenderer) PrintBanner() {
	fmt.Printf("%s%s========================================================================%s\n", colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s   consensai : Adversarial Multi-Agent Epistemic Deliberation Engine   %s\n", colorBold, colorWhite, colorReset)
	fmt.Printf("%s%s========================================================================%s\n\n", colorBold, colorCyan, colorReset)
}

func (r *TerminalRenderer) PrintPhase(phase string) {
	fmt.Printf("\n%s%s>>> [PHASE] %s%s\n", colorBold, colorYellow, phase, colorReset)
}

func (r *TerminalRenderer) PrintTurn(turn *protocol.AgentTurn) {
	var roleColor string
	switch turn.Role {
	case protocol.RoleAdvocate:
		roleColor = colorGreen
	case protocol.RoleSkeptic:
		roleColor = colorRed
	case protocol.RoleAuditor:
		roleColor = colorCyan
	default:
		roleColor = colorWhite
	}

	fmt.Printf("\n%s%s=== Round %d : %s ===%s\n", colorBold, roleColor, turn.Round, turn.Role, colorReset)
	fmt.Printf("%sSummary:%s %s\n", colorBold, colorReset, turn.Summary)

	if len(turn.Claims) > 0 {
		fmt.Printf("\n  %sClaims Submitted:%s\n", colorDim, colorReset)
		for _, c := range turn.Claims {
			fmt.Printf("   * [%s%s%s] %s\n     %sCitation: %s%s\n",
				roleColor, c.ID, colorReset, c.Statement, colorDim, c.Citation, colorReset)
		}
	}

	if len(turn.Rebuttals) > 0 {
		fmt.Printf("\n  %sRebuttals & Counter-Evidence:%s\n", colorDim, colorReset)
		for _, reb := range turn.Rebuttals {
			concessionTag := ""
			if reb.Concession {
				concessionTag = fmt.Sprintf(" %s(CONCEDED)%s", colorYellow, colorReset)
			}
			fmt.Printf("   * Against [%s]: %s%s\n", reb.TargetClaimID, reb.CounterPoint, concessionTag)
		}
	}

	if len(turn.Concessions) > 0 {
		fmt.Printf("\n  %sFormal Concessions Recorded:%s\n", colorYellow, colorReset)
		for _, con := range turn.Concessions {
			fmt.Printf("   - %s\n", con)
		}
	}
}

func (r *TerminalRenderer) PrintFindings(findings []protocol.AuditFinding) {
	fmt.Printf("\n%s%s=== Round 3 : Risk Auditor Grounding Findings ===%s\n", colorBold, colorCyan, colorReset)
	for _, f := range findings {
		var statusBadge string
		switch f.Status {
		case protocol.StatusVerified:
			statusBadge = fmt.Sprintf("%s[VERIFIED]%s", colorGreen, colorReset)
		case protocol.StatusConceded:
			statusBadge = fmt.Sprintf("%s[CONCEDED]%s", colorYellow, colorReset)
		case protocol.StatusUnmitigated:
			statusBadge = fmt.Sprintf("%s[UNMITIGATED RISK]%s", colorRed, colorReset)
		case protocol.StatusHallucinated:
			statusBadge = fmt.Sprintf("%s[HALLUCINATED]%s", colorRed, colorReset)
		default:
			statusBadge = fmt.Sprintf("%s[%s]%s", colorDim, f.Status, colorReset)
		}

		fmt.Printf("  * %s %s: %s\n    %sFinding: %s (Confidence: %.0f%%)%s\n",
			statusBadge, f.ClaimID, f.Statement, colorDim, f.Discrepancy, f.ConfidenceScore*100, colorReset)
	}
}

func (r *TerminalRenderer) PrintSynthesis(s *protocol.SynthesisResult, headHash string) {
	var decColor string
	switch s.Decision {
	case protocol.OutcomeApprove:
		decColor = colorGreen
	case protocol.OutcomeReject:
		decColor = colorRed
	default:
		decColor = colorYellow
	}

	fmt.Printf("\n%s========================================================================%s\n", colorBold, colorReset)
	fmt.Printf("%s%s                  FINAL GOVERNANCE SYNTHESIS                     %s\n", colorBold, colorWhite, colorReset)
	fmt.Printf("%s========================================================================%s\n\n", colorBold, colorReset)

	fmt.Printf("  %sDECISION OUTCOME:%s %s%s%s%s\n\n", colorBold, colorReset, colorBold, decColor, s.Decision, colorReset)
	fmt.Printf("  %sEpistemic Convergence Score:%s  %.2f / 1.00\n", colorBold, colorReset, s.NormalizedScore)
	fmt.Printf("  %sConfidence Interval (95%%):%s    [%.2f, %.2f]\n", colorBold, colorReset, s.ConfidenceInterval[0], s.ConfidenceInterval[1])
	fmt.Printf("  %sVerified Affirmative Claims:%s  %d / %d\n", colorBold, colorReset, s.VerifiedAffirmativeClaims, s.TotalClaimsAudited)
	fmt.Printf("  %sUnmitigated Adversarial Risks:%s %d\n", colorBold, colorReset, s.UnmitigatedAdversarialRisks)
	fmt.Printf("  %sLedger Head SHA-256:%s           %s\n\n", colorBold, colorReset, headHash)

	fmt.Printf("  %sRationale:%s\n  %s\n\n", colorBold, colorReset, s.ActionRationale)

	if len(s.PreservedDissents) > 0 {
		fmt.Printf("  %sPreserved Minority Dissents:%s\n", colorRed, colorReset)
		for i, d := range s.PreservedDissents {
			fmt.Printf("    %d. %s\n", i+1, d)
		}
		fmt.Println()
	}

	fmt.Printf("%s------------------------------------------------------------------------%s\n", strings.Repeat("-", 10), colorReset)
}
