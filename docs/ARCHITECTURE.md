# System Architecture

consensai provides a multi-agent dialectical debate protocol and deterministic consensus engine for high-stakes enterprise decision workflows.

## Design Philosophy

Monolithic LLM completions suffer from sycophancy, confirmation bias, ungrounded confidence, and silent omissions. In regulated domains such as commercial credit underwriting, wire fraud adjudication, and healthcare compliance, unvalidated model outputs introduce systemic liability.

consensai addresses these failure modes by decoupling deliberation into adversarial roles with opposing objective functions:
1. The Advocate maximizes the affirmative justification for the dossier.
2. The Skeptic seeks contractual breaches, debt burdens, fraud indicators, and downside tail risks.
3. The Risk Auditor extracts verifiable assertions and reconciles them against the raw input documents.

A deterministic synthesizer calculates the final consensus score without stochastic variance.

```
                          Decision Dossier
                                 |
                                 v
               +-----------------------------------+
               |       consensai Orchestrator      |
               |       (Go Engine, < 25MB RAM)     |
               +-----------------------------------+
                                 |
                +----------------+----------------+
                |                                 |
                v                                 v
        [Round 1: Thesis]               [Round 1: Antithesis]
          The Advocate                      The Skeptic
                \                                 /
                 +---------------+---------------+
                                 |
                                 v
                     [Round 2: Cross-Rebuttal]
                     Concessions & Defense
                                 |
                                 v
                 [Round 3: Grounded Citation Audit]
                          The Risk Auditor
                                 |
                                 v
                [Deterministic Consensus Synthesizer]
                 - Epistemic Convergence Score
                 - Preserved Minority Dissents
                 - Decision: APPROVE / REJECT / ESCALATE
                                 |
                                 v
                [SHA-256 Chained Cryptographic Ledger]
                                 |
               +-----------------+-----------------+
               |                                   |
               v                                   v
      Embedded React 19 UI              Executive Decision Memo
     (Server-Sent Events)               (Markdown, JSON, HTML)
```

## Memory and Performance Envelope

The Go orchestrator is engineered to operate under strict resource constraints:
- Resident Set Size (RSS) memory profile under 25 megabytes during active deliberations.
- Zero external runtime dependencies; all web assets are embedded into the compiled binary via Go `embed`.
- Event streaming through standard HTTP Server-Sent Events (SSE) with minimal buffer allocation.
- Local execution option via deterministic heuristic mock engine that requires no external network calls or API keys.

## Core Components

### 1. Agent Runtime (`internal/agent`)
Defines the `Agent` and `ProviderClient` interfaces. It supports OpenAI, Anthropic Claude, Google Gemini, Ollama, and an offline mock provider. Each persona receives tailored system prompts with explicit operational boundaries and output schema requirements.

### 2. Protocol Engine (`internal/protocol`)
A deterministic finite state machine that enforces the three-round dialectic sequence:
- Round 1: Blind Framing. The Advocate and Skeptic evaluate the dossier in isolation to prevent anchoring bias.
- Round 2: Targeted Rebuttal. The agents exchange structured arguments. The Advocate must answer every identified flaw with verified citations; the Skeptic must record formal concessions when assertions are validated.
- Round 3: Factual Grounding Audit. The Risk Auditor extracts quantitative and factual claims and performs lexical and numerical matching against the raw dossier text.

### 3. Epistemic Synthesizer (`internal/protocol/synthesis.go`)
Calculates the Epistemic Convergence Score ($C_{\text{epistemic}}$), evaluates threshold boundaries ($\theta_{\text{approve}}$, $\theta_{\text{reject}}$), and isolates unresolved minority dissents.

### 4. Cryptographic Ledger (`internal/ledger`)
Maintains an append-only block chain. Every round transcript, citation audit result, and synthesis block includes the SHA-256 hash of the preceding block. The chain can be serialized to JSON and verified independently by third-party auditors.

### 5. Web Cockpit and Server (`internal/server`, `web/`)
Provides a local HTTP dashboard on port 6060. The frontend is built with React 19, TypeScript, and Tailwind CSS. It streams debate events in real time using Server-Sent Events, displays the convergence gauge, and provides interactive ledger tamper testing.
