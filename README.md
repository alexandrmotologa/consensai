# consensai

Multi-agent adversarial debate and calibrated consensus decision protocol for high-stakes AI automation.

consensai coordinates three specialized agents (The Advocate, The Skeptic, and The Risk Auditor) across three structured rounds of confrontation, citation audit, and factual grounding. A deterministic consensus synthesizer calculates an Epistemic Convergence Score, outputs an actionable decision (Approve, Reject, or Escalate to Human Oversight), preserves dissenting arguments, and generates a SHA-256 cryptographic audit ledger.

The system is delivered as a lightweight Go binary with an embedded React 19 web cockpit and interactive command line tools, with memory footprint below 25MB RAM.

---

## Key Features

- **Adversarial Multi-Agent Roles**: Splits reasoning across an affirmative Advocate, a cross-examining Skeptic, and an evidence-grounded Risk Auditor.
- **Three-Round Dialectical Protocol**: Enforces isolated framing, targeted cross-rebuttals, and source citation audits.
- **Deterministic Consensus Synthesizer**: Epistemic Convergence Score calculation with confidence intervals and calibrated threshold gating.
- **Preserved Minority Dissents**: Retains unmitigated risks and counter-evidence in executive summaries even during approval outcomes.
- **Cryptographic Audit Ledger**: Chained SHA-256 blocks record raw dossiers, prompt transcripts, citations, and voting records for tamper-evident verification.
- **Embedded Web Cockpit**: React 19 single-page application embedded directly into the Go executable, running on local port 6060.
- **Provider Agnostic**: Built-in support for OpenAI, Anthropic Claude, Google Gemini, Ollama, and a fully offline deterministic mock engine.

---

## Architecture Overview

```
               High-Stakes Decision Input Dossier
           (Credit Application, Fraud Investigation, BAA)
                               |
                               v
                     consensai Orchestrator
             (Single-Binary Go Engine, < 25MB RAM)
                               |
   +---------------------------+---------------------------+
   | Round 1: Thesis & Antithesis                          |
   | Advocate (Affirmative Case) | Skeptic (Risk Assessment)|
   +-------------------------------------------------------+
                               |
   +---------------------------+---------------------------+
   | Round 2: Targeted Cross-Rebuttal                      |
   | Evidence Exchanges and Documented Concessions         |
   +-------------------------------------------------------+
                               |
   +---------------------------+---------------------------+
   | Round 3: Factual Grounding & Citation Audit          |
   | Risk Auditor Verifies Claims Against Source Text     |
   +-------------------------------------------------------+
                               |
   +---------------------------+---------------------------+
   | Deterministic Consensus Synthesizer                   |
   | Calculates Epistemic Convergence Score (C_epistemic)  |
   | Decision: APPROVE / REJECT / ESCALATE                 |
   +-------------------------------------------------------+
                               |
   +---------------------------+---------------------------+
   | SHA-256 Chained Cryptographic Ledger                  |
   | Tamper-Proof Audit Trail with Merkle Proofs          |
   +-------------------------------------------------------+
                               |
            +------------------+------------------+
            |                                     |
            v                                     v
   Embedded React 19 Cockpit             Executive Decision Memo
       (http://localhost:6060)               (Markdown, JSON, HTML)
```

---

## Installation

### Prerequisites

- Go 1.22 or later
- Node.js 20 or later (for building the web dashboard)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/alexandrmotologa/consensai.git
cd consensai

# Build the web frontend
cd web
npm install
npm run build
cd ..

# Build the single Go binary
go build -o bin/consensai.exe ./cmd/consensai
```

---

## Usage

### Run a Deliberation from the Command Line

Run an adversarial debate over a business dossier:

```bash
consensai debate --dossier examples/dossiers/commercial_loan_500k.json
```

Specify a custom LLM provider and export formats:

```bash
consensai debate \
  --dossier examples/dossiers/wire_fraud_investigation.json \
  --provider openai \
  --model gpt-4o \
  --output report.md \
  --ledger audit.json
```

### Launch the Embedded Web Cockpit

Start the HTTP server on port 6060:

```bash
consensai serve --port 6060
```

Open `http://localhost:6060` in any modern web browser to access:
- Live multi-column debate stream.
- Epistemic convergence dial with adjustable threshold markers.
- Interactive claim and citation verification matrix.
- Preserved minority dissent cards.
- Cryptographic ledger explorer with one-click verification.

### Verify Ledger Integrity

Verify the cryptographic integrity of an audit ledger:

```bash
consensai verify --ledger audit.json
```

If any block, transcript, or citation in the chain was altered, the verification tool reports the exact block index and hash mismatch.

### Generate Executive Decision Memo

Extract a formal decision memo from a completed ledger:

```bash
consensai memo --ledger audit.json --format markdown --out memo.md
```

---

## Epistemic Convergence Metric

The orchestrator computes the Epistemic Convergence Score ($C_{\text{epistemic}}$):

$$C_{\text{epistemic}} = \frac{|\mathcal{A} \cap \mathcal{V}| - |\mathcal{S}_{\text{unmitigated}}|}{|\mathcal{A}| + |\mathcal{S}|}$$

Where:
- $\mathcal{A}$: Set of claims defended by the Advocate.
- $\mathcal{S}$: Set of risks and failure points raised by the Skeptic.
- $\mathcal{V}$: Set of claims verified against the source dossier by the Risk Auditor.
- $\mathcal{S}_{\text{unmitigated}}$: Validated risks left unmitigated by the Advocate.

### Decision Boundaries

- **APPROVE**: $C_{\text{epistemic}} \ge 0.70$
- **REJECT**: $C_{\text{epistemic}} \le 0.35$
- **ESCALATE TO HUMAN OVERSIGHT**: $0.35 < C_{\text{epistemic}} < 0.70$

---

## Configuration

Configuration is managed via `consensai.yaml` or command line flags:

```yaml
version: "1.0"
server:
  port: 6060
  cors_allowed_origins: ["*"]

provider:
  default: "mock" # options: mock, openai, anthropic, gemini, ollama
  api_key_env: "OPENAI_API_KEY"
  model: "gpt-4o"
  temperature: 0.2

protocol:
  max_rounds: 3
  thresholds:
    approve: 0.70
    reject: 0.35
  strict_citations: true
  min_citation_coverage: 0.80

ledger:
  storage_path: "./ledgers"
  hash_algorithm: "sha256"
```

---

## Directory Layout

```
consensai/
├── cmd/
│   └── consensai/           # CLI entry point (main.go)
├── internal/
│   ├── config/              # Configuration schemas and parser
│   ├── agent/               # Agent roles, prompt templates, and LLM clients
│   ├── protocol/            # Dialectical engine, state machine, and synthesizer
│   ├── ledger/              # SHA-256 block structure, chain validation, exports
│   ├── report/              # Memo generator and terminal ANSI printer
│   └── server/              # HTTP API, SSE streaming, embedded assets
├── web/                     # React 19 dashboard (Vite, TypeScript, Tailwind)
├── examples/
│   ├── consensai.yaml       # Example configuration
│   └── dossiers/            # Sample evaluation cases (credit, fraud, compliance)
└── docs/                    # Technical specifications and math models
```

---

## License

Apache License 2.0. See [LICENSE](LICENSE) for details.
