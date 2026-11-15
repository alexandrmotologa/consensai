package report

import (
	"fmt"
	"strings"
	"time"

	"github.com/alexandrmotologa/consensai/internal/ledger"
	"github.com/alexandrmotologa/consensai/internal/protocol"
)

// GeneratePrintableCertificateHTML renders a self-contained, high-typography institutional certificate.
func GeneratePrintableCertificateHTML(state *protocol.DeliberationState, l *ledger.Ledger) string {
	s := state.Synthesis
	headHash := "N/A"
	if l != nil && len(l.Blocks) > 0 {
		headHash = l.Blocks[len(l.Blocks)-1].Hash
	}

	merkleRoot := "N/A"
	if s != nil && s.MerkleRoot != "" {
		merkleRoot = s.MerkleRoot
	}

	var sb strings.Builder

	sb.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>consensai Audit Certificate - ` + state.ID + `</title>
  <style>
    body {
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
      margin: 40px;
      color: #1a202c;
      background: #ffffff;
      line-height: 1.5;
    }
    .header {
      border-bottom: 2px solid #0f172a;
      padding-bottom: 20px;
      margin-bottom: 30px;
      display: flex;
      justify-content: space-between;
      align-items: flex-end;
    }
    .title {
      font-size: 24px;
      font-weight: 800;
      letter-spacing: -0.5px;
      color: #0f172a;
      margin: 0;
    }
    .subtitle {
      font-size: 13px;
      color: #64748b;
      margin-top: 4px;
    }
    .seal {
      border: 2px solid #0f172a;
      padding: 6px 12px;
      font-family: monospace;
      font-size: 11px;
      font-weight: 700;
      letter-spacing: 1px;
      text-transform: uppercase;
    }
    .meta-grid {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 12px 24px;
      background: #f8fafc;
      border: 1px solid #e2e8f0;
      padding: 16px;
      border-radius: 6px;
      font-size: 13px;
      margin-bottom: 25px;
    }
    .meta-item { display: flex; justify-content: space-between; }
    .meta-label { color: #64748b; }
    .meta-value { font-family: monospace; font-weight: 600; color: #0f172a; }
    .decision-banner {
      padding: 16px;
      border-radius: 6px;
      margin-bottom: 25px;
      text-align: center;
    }
    .approve { background: #ecfdf5; border: 1px solid #10b981; color: #065f46; }
    .reject { background: #fff1f2; border: 1px solid #f43f5e; color: #9f1239; }
    .escalate { background: #fffbeb; border: 1px solid #f59e0b; color: #92400e; }
    .decision-title { font-size: 18px; font-weight: 800; margin: 0 0 6px 0; }
    .section-title {
      font-size: 15px;
      font-weight: 700;
      border-bottom: 1px solid #e2e8f0;
      padding-bottom: 6px;
      margin: 25px 0 12px 0;
      color: #0f172a;
    }
    table { width: 100%; border-collapse: collapse; font-size: 12px; margin-bottom: 20px; }
    th { background: #f1f5f9; text-align: left; padding: 8px; border: 1px solid #cbd5e1; font-weight: 600; }
    td { padding: 8px; border: 1px solid #cbd5e1; }
    .mono { font-family: monospace; font-size: 11px; }
    .badge { padding: 2px 6px; border-radius: 4px; font-weight: 600; font-size: 10px; }
    .badge-pass { background: #d1fae5; color: #065f46; }
    .badge-flag { background: #fee2e2; color: #991b1b; }
    .badge-act { background: #fef3c7; color: #92400e; }
    .print-btn {
      position: fixed;
      top: 20px;
      right: 20px;
      background: #0f172a;
      color: #ffffff;
      border: none;
      padding: 8px 16px;
      border-radius: 4px;
      cursor: pointer;
      font-weight: 600;
    }
    @media print {
      body { margin: 0; }
      .print-btn { display: none; }
    }
  </style>
</head>
<body>
  <button class="print-btn" onclick="window.print()">Print / Save as PDF</button>

  <div class="header">
    <div>
      <h1 class="title">consensai Compliance & Audit Certificate</h1>
      <div class="subtitle">Multi-Agent Epistemic Deliberation & Cryptographic Decision Trail</div>
    </div>
    <div class="seal">SHA-256 IMMUTABLE</div>
  </div>

  <div class="meta-grid">
    <div class="meta-item"><span class="meta-label">Session Identifier:</span><span class="meta-value">` + state.ID + `</span></div>
    <div class="meta-item"><span class="meta-label">Dossier Evaluated:</span><span class="meta-value">` + state.DossierTitle + `</span></div>
    <div class="meta-item"><span class="meta-label">Evaluation Timestamp:</span><span class="meta-value">` + state.CreatedAt.Format(time.RFC3339) + `</span></div>
    <div class="meta-item"><span class="meta-label">Epistemic Score (C_epistemic):</span><span class="meta-value">` + fmt.Sprintf("%.2f / 1.00", s.NormalizedScore) + `</span></div>
    <div class="meta-item"><span class="meta-label">Merkle Tree Root:</span><span class="meta-value">` + merkleRoot + `</span></div>
    <div class="meta-item"><span class="meta-label">Head Block Hash:</span><span class="meta-value">` + headHash + `</span></div>
  </div>`)

	// Decision Banner
	bannerClass := "escalate"
	if s.Decision == protocol.OutcomeApprove {
		bannerClass = "approve"
	} else if s.Decision == protocol.OutcomeReject {
		bannerClass = "reject"
	}

	sb.WriteString(fmt.Sprintf(`
  <div class="decision-banner %s">
    <div class="decision-title">FINAL GOVERNANCE DETERMINATION: %s</div>
    <div>%s</div>
  </div>`, bannerClass, s.Decision, s.ActionRationale))

	// Human Adjudication (if present)
	if state.Adjudication != nil {
		adj := state.Adjudication
		sb.WriteString(fmt.Sprintf(`
  <div class="section-title">Human-in-the-Loop (HITL) Adjudication (EU AI Act Art. 14)</div>
  <div class="meta-grid" style="background: #eff6ff; border-color: #bfdbfe;">
    <div class="meta-item"><span class="meta-label">Adjudicating Officer:</span><span class="meta-value">%s (%s)</span></div>
    <div class="meta-item"><span class="meta-label">Officer ID:</span><span class="meta-value">%s</span></div>
    <div class="meta-item"><span class="meta-label">Adjudication Decision:</span><span class="meta-value" style="color:#1d4ed8;">%s</span></div>
    <div class="meta-item"><span class="meta-label">Signature Hash:</span><span class="meta-value">%s</span></div>
  </div>
  <p><strong>Officer Justification:</strong> %s</p>`,
			adj.OfficerName, adj.OfficerRole, adj.OfficerID, adj.Decision, adj.SignatureHash, adj.Justification))

		if len(adj.Conditions) > 0 {
			sb.WriteString("<p><strong>Mandatory Operational Conditions:</strong></p><ul>")
			for _, cond := range adj.Conditions {
				sb.WriteString("<li>" + cond + "</li>")
			}
			sb.WriteString("</ul>")
		}
	}

	// Compliance Section
	if s.Compliance != nil && len(s.Compliance.Checks) > 0 {
		sb.WriteString(`
  <div class="section-title">Regulatory Compliance & Statutory Scorecard</div>
  <table>
    <thead>
      <tr>
        <th>Standard</th>
        <th>Description</th>
        <th>Status</th>
        <th>Auditor Regulatory Findings</th>
      </tr>
    </thead>
    <tbody>`)

		for _, chk := range s.Compliance.Checks {
			badgeClass := "badge-pass"
			if chk.Status == "FLAGGED" {
				badgeClass = "badge-flag"
			} else if chk.Status == "REQUIRES_HUMAN_ACTION" {
				badgeClass = "badge-act"
			}
			sb.WriteString(fmt.Sprintf(`
      <tr>
        <td class="mono">%s</td>
        <td>%s</td>
        <td><span class="badge %s">%s</span></td>
        <td>%s</td>
      </tr>`, chk.Standard, chk.Description, badgeClass, chk.Status, chk.Notes))
		}
		sb.WriteString(`</tbody></table>`)
	}

	// Cryptographic Block Ledger Chain
	if l != nil && len(l.Blocks) > 0 {
		sb.WriteString(`
  <div class="section-title">Cryptographic Hash Chain Sequence (SHA-256)</div>
  <table>
    <thead>
      <tr>
        <th>Index</th>
        <th>Type</th>
        <th>Previous Block Hash</th>
        <th>Committed Block Hash</th>
      </tr>
    </thead>
    <tbody>`)

		for _, b := range l.Blocks {
			sb.WriteString(fmt.Sprintf(`
      <tr>
        <td class="mono">#%d</td>
        <td class="mono">%s</td>
        <td class="mono">%s</td>
        <td class="mono"><strong>%s</strong></td>
      </tr>`, b.Index, b.Type, b.PrevHash, b.Hash))
		}
		sb.WriteString(`</tbody></table>`)
	}

	sb.WriteString(`
  <div style="margin-top: 40px; border-top: 1px solid #cbd5e1; padding-top: 12px; font-size: 11px; color: #64748b; text-align: center;">
    consensai Epistemic Verification Engine &bull; Non-Repudiable Cryptographic Ledger Certificate &bull; Produced: ` + time.Now().UTC().Format(time.RFC3339) + `
  </div>
</body>
</html>`)

	return sb.String()
}
