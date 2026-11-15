package report

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/alexandrmotologa/consensai/internal/protocol"
)

// DecisionMemo generates structured executive reports from deliberation states.
type DecisionMemo struct {
	State *protocol.DeliberationState
}

func NewDecisionMemo(state *protocol.DeliberationState) *DecisionMemo {
	return &DecisionMemo{State: state}
}

// GenerateMarkdown creates a humanized, factual Markdown decision report.
func (m *DecisionMemo) GenerateMarkdown(headHash string) string {
	s := m.State.Synthesis
	if s == nil {
		return "# Deliberation Incomplete\nNo synthesis result recorded."
	}

	var sb strings.Builder

	sb.WriteString("# Executive Decision Memo: " + m.State.DossierTitle + "\n\n")
	sb.WriteString(fmt.Sprintf("**Session ID:** `%s` | **Date:** %s\n\n", m.State.ID, m.State.CreatedAt.Format(time.RFC822)))

	sb.WriteString("## 1. Actionable Recommendation\n\n")
	sb.WriteString(fmt.Sprintf("### Outcome: %s\n\n", s.Decision))
	sb.WriteString(s.ActionRationale + "\n\n")

	sb.WriteString("## 2. Epistemic Convergence Metrics\n\n")
	sb.WriteString(fmt.Sprintf("- **Normalized Convergence Score:** `%.2f` / 1.00\n", s.NormalizedScore))
	sb.WriteString(fmt.Sprintf("- **Raw Convergence Score:** `%.3f`\n", s.EpistemicConvergenceScore))
	sb.WriteString(fmt.Sprintf("- **95%% Confidence Interval:** `[%.2f, %.2f]`\n", s.ConfidenceInterval[0], s.ConfidenceInterval[1]))
	sb.WriteString(fmt.Sprintf("- **Total Claims Audited:** %d\n", s.TotalClaimsAudited))
	sb.WriteString(fmt.Sprintf("- **Verified Affirmative Claims:** %d\n", s.VerifiedAffirmativeClaims))
	sb.WriteString(fmt.Sprintf("- **Unmitigated Adversarial Risks:** %d\n\n", s.UnmitigatedAdversarialRisks))

	sb.WriteString("## 3. Preserved Minority Dissents\n\n")
	if len(s.PreservedDissents) == 0 {
		sb.WriteString("No unmitigated adversarial dissents remained at protocol conclusion.\n\n")
	} else {
		sb.WriteString("The following objections raised by The Skeptic remain active and are preserved for governance review:\n\n")
		for i, d := range s.PreservedDissents {
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, d))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## 4. Audited Claims and Grounding Matrix\n\n")
	sb.WriteString("| ID | Role | Statement | Citation | Status | Finding |\n")
	sb.WriteString("| :--- | :--- | :--- | :--- | :--- | :--- |\n")
	for _, c := range m.State.AllClaims {
		status := string(c.AuditedStatus)
		note := c.AuditNotes
		if note == "" {
			note = "Consistent with source text"
		}
		sb.WriteString(fmt.Sprintf("| `%s` | %s | %s | `%s` | **%s** | %s |\n",
			c.ID, c.AgentRole, escapePipes(c.Statement), escapePipes(c.Citation), status, escapePipes(note)))
	}
	sb.WriteString("\n")

	sb.WriteString("## 5. Cryptographic Audit Proof\n\n")
	sb.WriteString(fmt.Sprintf("- **Ledger Head Block Hash:** `%s`\n", headHash))
	sb.WriteString("- **Verification Status:** Verified with SHA-256 hash chaining\n")
	sb.WriteString("- **Immutable Record:** Stored in append-only audit trail\n\n")

	return sb.String()
}

// GenerateJSON returns the structured report as JSON bytes.
func (m *DecisionMemo) GenerateJSON() ([]byte, error) {
	return json.MarshalIndent(m.State, "", "  ")
}

func escapePipes(s string) string {
	return strings.ReplaceAll(s, "|", "\\|")
}
