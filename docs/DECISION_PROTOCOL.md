# Deliberation Protocol and Epistemic Formulation

consensai implements a structured dialectical arbitration protocol designed to convert adversarial agent deliberations into a deterministic, auditable decision.

## Protocol Lifecycle

The deliberation protocol progresses through four distinct phases.

| Phase | Persona Involved | Input Context | Output Artifact |
| :--- | :--- | :--- | :--- |
| Round 1: Blind Framing | The Advocate, The Skeptic | Raw input dossier | Thesis claims, adversarial counter-claims |
| Round 2: Cross-Rebuttal | The Advocate, The Skeptic | Opponent Round 1 claims | Rebuttals, counter-evidence, formal concessions |
| Round 3: Factual Audit | The Risk Auditor | Raw dossier, all claims | Verified claims, unverified claims, citations |
| Synthesis | Deterministic Synthesizer | Audited claim set | Epistemic Convergence Score, final decision |

### Phase 1: Blind Framing
The Advocate and Skeptic process the raw dossier simultaneously without observing each other's reasoning. This guarantees that initial threat modeling and affirmative arguments are formed independently, preventing conversational anchoring.

### Phase 2: Targeted Rebuttal
The Advocate receives the full list of objections from the Skeptic. For each objection, the Advocate must provide documentary counter-evidence or concede the point. Simultaneously, the Skeptic reviews the Advocate's affirmative claims and highlights unaddressed risks.

### Phase 3: Factual Grounding Audit
The Risk Auditor operates as an impartial judicial clerk. The Auditor parses all factual statements, numerical ratios (such as Debt Service Coverage Ratio, loan-to-value, IP geolocation distances), and regulatory citations made by both parties. Each assertion is matched against the raw source dossier.
- Validated claims receive a status of `VERIFIED`.
- Claims lacking textual ground truth receive a status of `UNVERIFIED` or `HALLUCINATED`.
- Struck claims carry a penalty against the claiming agent.

## Mathematical Formulation

### Epistemic Convergence Metric

Let $\mathcal{A}$ represent the set of claims defended by the Advocate, $\mathcal{S}$ the set of risks raised by the Skeptic, $\mathcal{V}$ the set of claims verified by the Auditor, and $\mathcal{S}_{\text{unmitigated}}$ the subset of Skeptic risks that lack verified mitigation:

$$C_{\text{epistemic}} = \frac{|\mathcal{A} \cap \mathcal{V}| - |\mathcal{S}_{\text{unmitigated}}|}{|\mathcal{A}| + |\mathcal{S}|}$$

When normalized, $C_{\text{epistemic}}$ ranges from $0.0$ to $1.0$:

$$C_{\text{normalized}} = \max\left(0.0, \min\left(1.0, \frac{C_{\text{epistemic}} + 1.0}{2.0}\right)\right)$$

### Decision Thresholds

The synthesizer maps $C_{\text{normalized}}$ to one of three final governance actions:

1. **APPROVE**: $C_{\text{normalized}} \ge \theta_{\text{approve}}$ (default: $0.70$).
   The affirmative case is strongly supported by verified documentary evidence, and identified risks have documented mitigations.
2. **REJECT**: $C_{\text{normalized}} \le \theta_{\text{reject}}$ (default: $0.35$).
   The dossier contains unmitigated risks, failed credit covenants, or unverified factual claims.
3. **ESCALATE TO HUMAN OVERSIGHT**: $\theta_{\text{reject}} < C_{\text{normalized}} < \theta_{\text{approve}}$.
   The debate revealed significant friction, unresolved objections, or borderline factual coverage that requires human executive judgment.

### Calibrated Confidence Interval

To quantify uncertainty across deliberations, the synthesizer calculates a standard error estimate based on total audited assertions $N = |\mathcal{A}| + |\mathcal{S}|$:

$$SE = \sqrt{\frac{C(1 - C)}{N}}$$

The 95% confidence interval ($z = 1.96$) is:

$$CI_{95} = \left[ \max(0, C - 1.96 \cdot SE), \min(1, C + 1.96 \cdot SE) \right]$$

## Preservation of Minority Dissents

Standard majority voting or naive LLM summarization tends to discard outlier warnings. consensai requires that every unmitigated risk identified by the Skeptic and verified by the Auditor is preserved verbatim in the final decision memo and ledger, even when the final outcome is an approval.
