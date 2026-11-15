# HTTP and Server-Sent Events API

The consensai orchestrator provides an HTTP server on port 6060 with REST endpoints and real-time Server-Sent Events (SSE).

## Endpoints

### 1. Server Status
- **Method**: `GET`
- **Path**: `/api/status`
- **Response**:
```json
{
  "status": "online",
  "version": "1.0.0",
  "uptime_seconds": 128,
  "allocated_memory_bytes": 14285900,
  "active_sessions": 1
}
```

### 2. List Dossiers
- **Method**: `GET`
- **Path**: `/api/dossiers`
- **Response**:
```json
[
  {
    "id": "commercial_loan_500k",
    "title": "Commercial Real Estate Loan ($500,000)",
    "domain": "credit_underwriting",
    "risk_level": "medium",
    "file_path": "examples/dossiers/commercial_loan_500k.json"
  }
]
```

### 3. Get Dossier Details
- **Method**: `GET`
- **Path**: `/api/dossiers/:id`
- **Response**: Returns raw JSON or Markdown content of the requested dossier.

### 4. Start Deliberation
- **Method**: `POST`
- **Path**: `/api/deliberate`
- **Request Body**:
```json
{
  "dossier_id": "commercial_loan_500k",
  "custom_dossier": null,
  "provider": "mock",
  "model": "heuristic-v1",
  "threshold_approve": 0.70,
  "threshold_reject": 0.35
}
```
- **Response**:
```json
{
  "session_id": "sess_89ef412a9c",
  "status": "started",
  "stream_url": "/api/deliberate/sess_89ef412a9c/stream"
}
```

### 5. Deliberation Event Stream (SSE)
- **Method**: `GET`
- **Path**: `/api/deliberate/:id/stream`
- **Protocol**: Server-Sent Events (`text/event-stream`)
- **Event Types**:
  - `status`: Deliberation status changes (`framing`, `rebuttal`, `audit`, `synthesizing`, `completed`).
  - `turn`: New agent contribution from Advocate, Skeptic, or Auditor.
  - `convergence`: Real-time score update and confidence interval.
  - `claim_update`: Individual claim validation status change.
  - `block_committed`: SHA-256 block added to the ledger.
  - `completed`: Final decision payload with memo summary.

### 6. Retrieve Ledger
- **Method**: `GET`
- **Path**: `/api/deliberate/:id/ledger`
- **Response**: Complete cryptographic ledger object with all blocks and hash chains.

### 7. Verify Ledger
- **Method**: `POST`
- **Path**: `/api/ledger/verify`
- **Request Body**: Raw JSON of an audit ledger.
- **Response**:
```json
{
  "valid": true,
  "block_count": 5,
  "head_hash": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
  "error": null
}
```

### 8. Download Decision Memo
- **Method**: `GET`
- **Path**: `/api/deliberate/:id/memo?format=markdown`
- **Response**: Formatted Markdown or HTML decision memo.

### 9. Human-in-the-Loop Adjudication
- **Method**: `POST`
- **Path**: `/api/deliberate/:id/adjudicate`
- **Request Body**:
```json
{
  "officer_id": "CRO-8842",
  "officer_name": "Elena Vance",
  "officer_role": "Chief Risk Officer",
  "decision": "CONDITIONAL_APPROVE",
  "justification": "Approved subject to $50,000 cash escrow cushion.",
  "conditions": ["Quarterly DSCR covenant test"]
}
```
- **Response**:
```json
{
  "status": "adjudicated",
  "block_index": 5,
  "head_hash": "a1b2c3d4...",
  "signature_hash": "e5f6..."
}
```

### 10. Sensitivity Recalculation (What-If Analysis)
- **Method**: `POST`
- **Path**: `/api/deliberate/:id/recalculate`
- **Request Body**:
```json
{
  "threshold_approve": 0.65,
  "threshold_reject": 0.30
}
```
- **Response**: Recomputed `SynthesisResult` with updated decision boundaries without re-invoking LLMs.

### 11. Printable Compliance Certificate
- **Method**: `GET`
- **Path**: `/api/deliberate/:id/certificate`
- **Response**: Self-contained printable HTML document with official seal, Merkle root hash, and statutory scorecard.
