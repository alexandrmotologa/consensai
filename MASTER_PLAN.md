# consensai: Architecture and Engineering Master Plan
> Multi-Agent Adversarial Debate & Calibrated Consensus Decision Protocol for High-Stakes AI Automation
> Target: Single-Binary Go Orchestrator (< 25MB RAM) with Embedded React 19 Deliberation Cockpit & Cryptographic Audit Ledger

---

## 1. Problem Statement and Motivation

### 1.1 The Vulnerability of Single-Agent Decisions in High-Stakes Workflows
Enterprises are increasingly deploying LLMs to automate mission-critical, high-liability decisions: credit underwriting, commercial loan approvals, fraud adjudication, insurance claim resolutions, and regulatory compliance audits. However, single-agent architectures exhibit critical vulnerabilities:
1. **Sycophancy and Confirmation Bias**: LLMs exhibit a tendency to validate user assumptions or double down on initial premises rather than rigorously seeking falsifying evidence.
2. **Ungrounded Confidence**: Single-model outputs produce high-confidence prose even when hallucinating factual figures, statutory clauses, or risk assessments.
3. **Fragility Under Edge Cases**: Complex decisions often hinge on subtle contractual conditions or regulatory exemptions that a single pass easily overlooks.
4. **Lack of Auditability**: Standard LLM completions provide only a final answer, lacking a transparent record of counter-arguments considered and dismissed during deliberation.

### 1.2 The Failure of Naive Self-Reflection
Common techniques like Chain-of-Thought (CoT) or prompting an LLM to "critique yourself" fail in high-stakes environments:
- **Shared Latent Biases**: A model reviewing its own output shares the exact same pre-training parametric blind spots.
- **Superficial Self-Correction**: When prompted to double-check, models often flip valid answers to incorrect ones or make purely stylistic edits without probing factual risks.

### 1.3 The consensai Solution
`consensai` replaces monolithic LLM decisions with a structured, multi-agent dialectical debate protocol inspired by formal judicial arbitration. Implemented as a lightweight, zero-dependency Go orchestrator:
- **Specialized Adversarial Roles**: Instantiates three distinct agent personas with opposing objective functions:
  - **The Advocate (Thesis)**: Marshals supporting evidence and argues the affirmative case.
  - **The Skeptic (Antithesis)**: Rigorously cross-examines assumptions, identifies systemic edge cases, and calculates downside exposure.
  - **The Risk Auditor (Verification)**: Deterministically validates factual claims against supplied source documents and audits quantitative assertions.
- **Three-Round Dialectical Protocol**: Enforces structured argumentation: Initial Claims $\rightarrow$ Direct Cross-Rebuttal $\rightarrow$ Grounded Audit & Synthesis.
- **Calibrated Epistemic Uncertainty**: Generates calibrated confidence intervals derived from agent convergence rather than raw token log-probabilities.
- **Preserved Minority Dissents**: Prevents artificial consensus by explicitly surfacing dissenting views and unmitigated risks in executive summaries.
- **Cryptographic Audit Ledger**: Chained SHA-256 blocks record every prompt, response, citation, and voting round to provide non-repudiable audit trails.

---

## 2. High-Level System Architecture

```
       +-------------------------------------------------------------+
       |             High-Stakes Decision Input Dossier              |
       |  (Loan Application, Legal Claim, Underwriting Data, Docs)   |
       +-------------------------------------------------------------+
                                      |
                                      v
       +-------------------------------------------------------------+
       |                     consensai Orchestrator                  |
       |                   (High-Performance Go CLI)                |
       |                                                             |
       |  +-------------------------------------------------------+  |
       |  |                    Round 1: Thesis                    |  |
       |  |   [The Advocate]             [The Skeptic]            |  |
       |  |   Affirmative Case           Adversarial Failure Case |  |
       |  +-------------------------------------------------------+  |
       |                              |                              |
       |                              v                              |
       |  +-------------------------------------------------------+  |
       |  |               Round 2: Cross-Rebuttal                 |  |
       |  |   The Advocate addresses Skeptic's edge cases         |  |
       |  |   The Skeptic concedes verified points                |  |
       |  +-------------------------------------------------------+  |
       |                              |                              |
       |                              v                              |
       |  +-------------------------------------------------------+  |
       |  |            Round 3: Grounded Audit & Falsification    |  |
       |  |   [The Risk Auditor]                                  |  |
       |  |   Verifies citations against raw source dossier       |  |
       |  |   Audits math, logic, and regulatory constraints      |  |
       |  +-------------------------------------------------------+  |
       |                              |                              |
       |                              v                              |
       |  +-------------------------------------------------------+  |
       |  |             Deterministic Consensus Synthesizer       |  |
       |  |   - Weighted Epistemic Convergence Score              |  |
       |  |   - Final Actionable Decision (Approve/Reject/Escalate) |  |
       |  |   - Preserved Dissenting Opinions & Unmitigated Risks |  |
       |  +-------------------------------------------------------+  |
       |                              |                              |
       |                              v                              |
       |  +-------------------------------------------------------+  |
       |  |              SHA-256 Chained Audit Ledger             |  |
       |  +-------------------------------------------------------+  |
       |                                                             |
       +------------------------------|------------------------------+
                                      |
                                      v
       +-------------------------------------------------------------+
       |                 Output & Governance Surfaces                |
       |  1. Embedded React 19 Cockpit (Local HTTP :6060)            |
       |  2. Executive Decision Memo (Markdown & PDF)                |
       |  3. Non-Repudiable Cryptographic Ledger (audit_ledger.json) |
       +-------------------------------------------------------------+
```

---

## 3. Deliberation Protocol & Epistemic Formulation

```
+-----------------------------------+-----------------------------------------------------------------+
| Protocol Phase                    | Operational Mandate                                             |
+-----------------------------------+-----------------------------------------------------------------+
| Round 1: Blind Framing            | Advocate and Skeptic receive dossier in isolation. Advocate     |
|                                   | maximizes approval utility; Skeptic seeks grounds for rejection.|
| Round 2: Targeted Rebuttal        | Agents exchange arguments. Advocate must provide verifiable    |
|                                   | counter-evidence for every flaw; Skeptic must concede or escalate.|
| Round 3: Factual Grounding Audit  | Auditor verifies each cited figure against input documents.     |
|                                   | Any unsupported claim is struck from the record with a penalty. |
| Synthesis: Consensus Calibration  | Synthesizer computes convergence ratio and produces decision:  |
|                                   | APPROVE, REJECT, or ESCALATE TO HUMAN OVERSIGHT.                |
+-----------------------------------+-----------------------------------------------------------------+
```

### 3.1 Epistemic Convergence Metric ($C_{\text{epistemic}}$)
Let $\mathcal{A}$ represent the set of claims defended by the Advocate, $\mathcal{S}$ the set of risks raised by the Skeptic, and $\mathcal{V}$ the set of claims verified by the Auditor:

$$C_{\text{epistemic}} = \frac{|\mathcal{A} \cap \mathcal{V}| - |\mathcal{S}_{\text{unmitigated}}|}{|\mathcal{A}| + |\mathcal{S}|}$$

When $C_{\text{epistemic}} \ge \theta_{\text{approve}}$, the decision is marked approved; when $C_{\text{epistemic}} \le \theta_{\text{reject}}$, it is rejected; intermediate scores trigger automatic escalation to human review with highlighted friction points.

---

## 4. Repository Structure

```
consensai/
├── cmd/
│   └── consensai/
│       └── main.go                 # CLI entry point, subcommands: debate, audit, verify, serve
├── internal/
│   ├── config/
│   │   ├── config.go               # Configuration parser for model endpoints and thresholds
│   │   └── schema.go               # Debate parameters, role definitions, and convergence bounds
│   ├── agent/
│   │   ├── advocate.go             # Affirmative case persona prompt and reasoning loop
│   │   ├── skeptic.go              # Adversarial risk persona prompt and cross-examiner
│   │   ├── auditor.go              # Evidence citation verifier and factual grounding checker
│   │   └── client.go               # LLM client interface (OpenAI, Anthropic, Gemini, Ollama)
│   ├── protocol/
│   │   ├── engine.go               # 3-round debate orchestrator and turn management
│   │   ├── state.go                # Debate transcript state machine and step progression
│   │   └── synthesis.go            # Deterministic decision consensus and convergence scorer
│   ├── ledger/
│   │   ├── block.go                # Cryptographic block structure (SHA-256 hash chaining)
│   │   ├── chain.go                # Append-only audit trail with tamper-detection verification
│   │   └── export.go               # JSON and markdown export of non-repudiable audit logs
│   ├── report/
│   │   ├── memo.go                 # Executive decision memo generator
│   │   └── terminal.go             # Rich terminal output displaying real-time debate rounds
│   └── server/
│       ├── api.go                  # REST API for interactive review and triggers
│       ├── embed.go                # Go embed filesystem for React 19 cockpit
│       └── sse.go                  # Server-Sent Events stream for live debate turns
├── web/
│   ├── src/
│   │   ├── App.tsx                 # React 19 root application
│   │   ├── components/
│   │   │   ├── DeliberationStream.tsx # Multi-column live debate feed (Advocate / Skeptic / Auditor)
│   │   │   ├── ConvergenceGauge.tsx   # Visual epistemic confidence dial with threshold markers
│   │   │   ├── ClaimMatrix.tsx        # Interactive table of claims, rebuttals, and audit statuses
│   │   │   ├── DissentCard.tsx        # Highlight box for preserved minority viewpoints
│   │   │   └── LedgerVerifier.tsx     # Cryptographic integrity checker with hash verification
│   │   └── pages/
│   │       └── Dashboard.tsx          # Main debate monitoring console
│   ├── package.json
│   ├── tsconfig.json
│   └── vite.config.ts
├── examples/
│   ├── consensai.yaml              # Sample configuration file
│   └── dossiers/                   # Sample enterprise decision cases
│       ├── commercial_loan_500k.json
│       ├── fraud_investigation_wire.json
│       └── hipaa_compliance_review.md
├── Makefile
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```

---

## 5. Implementation Roadmap & Status

### Phase 1: Go Engine Foundation and Configuration [COMPLETED]
- Initialized Go module (`github.com/alexandrmotologa/consensai`) with minimal dependencies.
- Built `internal/config` supporting YAML parameters, environment variable overrides, and decision thresholds.
- Implemented LLM provider client interface in `internal/agent/client.go` supporting OpenAI, Anthropic, Gemini, Ollama, and built-in offline mock provider.

### Phase 2: Multi-Agent Personas and Protocol State Machine [COMPLETED]
- Authored system and user instructions for `advocate.go`, `skeptic.go`, and `auditor.go`.
- Built the 3-round protocol state machine in `internal/protocol/engine.go` managing prompt context boundaries, turn sequences, and event emission.

### Phase 3: Factual Grounding Audit and Deterministic Synthesizer [COMPLETED]
- Implemented `auditor.go` matching cited claims against original source text using lexical matching and citation parsing.
- Built `internal/protocol/synthesis.go` calculating the Epistemic Convergence Score $C_{\text{epistemic}}$ with 95% confidence intervals and isolating unmitigated minority dissents.

### Phase 4: Cryptographic SHA-256 Audit Ledger [COMPLETED]
- Built append-only cryptographic ledger in `internal/ledger/chain.go` where every round transcript is hashed with the previous block hash using canonical byte preservation.
- Implemented tamper-detection verification subcommand (`consensai verify --ledger audit.json`).

### Phase 5: Executive Decision Memo & Terminal Interface [COMPLETED]
- Implemented CLI terminal formatting displaying real-time agent turns with distinct colors, citations, and summaries.
- Generated structured Markdown and JSON decision memos containing actionable recommendations, convergence scores, claim matrices, and preserved dissents.

### Phase 6: Embedded React 19 Deliberation Cockpit [COMPLETED]
- Built high-density web UI in `web/` featuring a multi-agent 3-column debate viewer, epistemic convergence dial gauge, claim verification matrix, preserved dissent card, and interactive ledger verifier.
- Embedded assets directly into the Go binary via `//go:embed` and exposed on local port `:6060`.

### Phase 7: Verification, Test Suite, and Documentation [COMPLETED]
- Wrote end-to-end integration tests running debates on sample loan and fraud dossiers (`go test -v ./...`).
- Validated memory usage via benchmark test: 0.25 MB heap allocation, staying well below the 25MB RAM ceiling.
- Authored humanized README and documentation in `docs/` (`ARCHITECTURE.md`, `DECISION_PROTOCOL.md`, `API.md`).

---

## 6. Architecture Improvements & Enhancements Added

1. **Deterministic Offline Heuristic Reasoning Engine**:
   - Enables full offline testing and local demonstrations across credit, fraud, and HIPAA compliance without requiring paid external API keys.
   - Evaluates real financial and operational metrics (DSCR, EBITDA, LTV, IP geolocations, DKIM headers, 90-day vs 6-year retention).

2. **Byte-Level Canonical JSON Serialization for Blockchain Hash Reproducibility**:
   - Implemented raw byte message serialization (`json.RawMessage`) and Unix millisecond epoch timestamps for block payloads, eliminating float/int and struct-versus-map hashing divergences across disk persistence cycles.

3. **Interactive Cryptographic Tamper Simulator**:
   - Integrated live tamper attack simulation in the web cockpit and CLI verification suite, enabling security teams to observe automated hash mismatch alerts when records are altered.

4. **Calibrated Epistemic Uncertainty & Confidence Interval Formulation**:
   - Integrated standard error derivation $SE = \sqrt{C(1-C)/N}$ and 95% confidence interval calculations $[CI_{95}]$ into both the Go synthesizer and the SVG circular gauge in the React 19 cockpit.
