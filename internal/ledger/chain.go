package ledger

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alexandrmotologa/consensai/internal/protocol"
)

const GenesisPrevHash = "0000000000000000000000000000000000000000000000000000000000000000"

// Ledger maintains an append-only, SHA-256 hash-chained cryptographic audit trail.
type Ledger struct {
	SessionID   string    `json:"session_id"`
	DossierID   string    `json:"dossier_id"`
	DossierHash string    `json:"dossier_hash"`
	CreatedAt   time.Time `json:"created_at"`
	Blocks      []*Block  `json:"blocks"`
}

// NewLedger initializes a fresh audit ledger with a genesis block.
func NewLedger(sessionID, dossierID, dossierRaw string) (*Ledger, error) {
	dossierHashBytes := sha256.Sum256([]byte(dossierRaw))
	dossierHash := hex.EncodeToString(dossierHashBytes[:])

	l := &Ledger{
		SessionID:   sessionID,
		DossierID:   dossierID,
		DossierHash: dossierHash,
		CreatedAt:   time.Now().UTC(),
		Blocks:      make([]*Block, 0, 8),
	}

	genesisPayload := map[string]interface{}{
		"session_id":   sessionID,
		"dossier_id":   dossierID,
		"dossier_hash": dossierHash,
		"description":  "Genesis block initializing consensai epistemic debate ledger.",
	}

	_, err := l.AddBlock(BlockTypeGenesis, 0, genesisPayload)
	if err != nil {
		return nil, fmt.Errorf("creating genesis block: %w", err)
	}

	return l, nil
}

// AddBlock appends a new block to the chain, computing its cryptographic link to the head block.
func (l *Ledger) AddBlock(blockType BlockType, round int, payload interface{}) (*Block, error) {
	prevHash := GenesisPrevHash
	index := 0

	if len(l.Blocks) > 0 {
		head := l.Blocks[len(l.Blocks)-1]
		prevHash = head.Hash
		index = head.Index + 1
	}

	var rawPayload json.RawMessage
	switch p := payload.(type) {
	case json.RawMessage:
		rawPayload = p
	case []byte:
		rawPayload = p
	default:
		bytes, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("marshaling payload for block %d: %w", index, err)
		}
		rawPayload = bytes
	}

	block := &Block{
		Index:       index,
		TimestampMs: time.Now().UTC().UnixMilli(),
		Round:       round,
		Type:        blockType,
		PrevHash:    prevHash,
		Payload:     rawPayload,
	}

	hash, err := block.CalculateHash()
	if err != nil {
		return nil, fmt.Errorf("hashing block: %w", err)
	}
	block.Hash = hash

	l.Blocks = append(l.Blocks, block)
	return block, nil
}

// Verify audits the integrity of the entire ledger chain.
// It detects any modifications, omitted records, or broken hash linkages.
func (l *Ledger) Verify() (bool, error) {
	if len(l.Blocks) == 0 {
		return false, fmt.Errorf("ledger contains zero blocks")
	}

	for i, b := range l.Blocks {
		if b.Index != i {
			return false, fmt.Errorf("block index mismatch at position %d: expected %d, got %d", i, i, b.Index)
		}

		if i == 0 {
			if b.PrevHash != GenesisPrevHash {
				return false, fmt.Errorf("genesis block has invalid prev_hash: %s", b.PrevHash)
			}
		} else {
			prev := l.Blocks[i-1]
			if b.PrevHash != prev.Hash {
				return false, fmt.Errorf("broken chain at block %d: prev_hash %s does not match block %d hash %s", i, b.PrevHash, i-1, prev.Hash)
			}
		}

		calculatedHash, err := b.CalculateHash()
		if err != nil {
			return false, fmt.Errorf("failed recalculating hash for block %d: %w", i, err)
		}

		if calculatedHash != b.Hash {
			return false, fmt.Errorf("cryptographic tamper detected at block %d: recorded hash %s != calculated hash %s", i, b.Hash, calculatedHash)
		}
	}

	return true, nil
}

// AddHumanAdjudication appends a signed human adjudication block to the ledger chain.
func (l *Ledger) AddHumanAdjudication(adj *protocol.HumanAdjudication) (*Block, error) {
	return l.AddBlock(BlockTypeHumanAdjudication, 5, map[string]interface{}{
		"phase":        "HUMAN_ADJUDICATION",
		"adjudication": adj,
	})
}

// CreateFromState constructs a complete verified ledger from a finalized deliberation state.
func CreateFromState(state *protocol.DeliberationState) (*Ledger, error) {
	l, err := NewLedger(state.ID, state.DossierID, state.DossierRaw)
	if err != nil {
		return nil, err
	}

	// Compute Merkle root over all audited claims
	if len(state.AllClaims) > 0 {
		var claimHashes []string
		for _, c := range state.AllClaims {
			claimHashes = append(claimHashes, fmt.Sprintf("%s:%s:%s:%s", c.ID, c.AgentRole, c.Statement, c.AuditedStatus))
		}
		if tree, err := NewMerkleTree(claimHashes); err == nil && state.Synthesis != nil {
			state.Synthesis.MerkleRoot = tree.RootHash()
		}
	}

	// Filter turns by round
	var r1Turns, r2Turns, r3Turns []protocol.AgentTurn
	for _, t := range state.Turns {
		switch t.Round {
		case 1:
			r1Turns = append(r1Turns, t)
		case 2:
			r2Turns = append(r2Turns, t)
		case 3:
			r3Turns = append(r3Turns, t)
		}
	}

	// Block 1: Round 1 Blind Framing
	if len(r1Turns) > 0 {
		if _, err := l.AddBlock(BlockTypeRound1, 1, map[string]interface{}{
			"phase": "ROUND_1_BLIND_FRAMING",
			"turns": r1Turns,
		}); err != nil {
			return nil, err
		}
	}

	// Block 2: Round 2 Cross-Rebuttal
	if len(r2Turns) > 0 {
		if _, err := l.AddBlock(BlockTypeRound2, 2, map[string]interface{}{
			"phase": "ROUND_2_CROSS_REBUTTAL",
			"turns": r2Turns,
		}); err != nil {
			return nil, err
		}
	}

	// Block 3: Round 3 Factual Grounding Audit
	if len(r3Turns) > 0 || len(state.AuditFindings) > 0 {
		if _, err := l.AddBlock(BlockTypeRound3, 3, map[string]interface{}{
			"phase":          "ROUND_3_FACTUAL_AUDIT",
			"turns":          r3Turns,
			"audit_findings": state.AuditFindings,
		}); err != nil {
			return nil, err
		}
	}

	// Block 4: Synthesis Decision
	if state.Synthesis != nil {
		if _, err := l.AddBlock(BlockTypeSynthesis, 4, map[string]interface{}{
			"phase":     "SYNTHESIS",
			"synthesis": state.Synthesis,
		}); err != nil {
			return nil, err
		}
	}

	// Block 5: Human Adjudication (if already recorded)
	if state.Adjudication != nil {
		if _, err := l.AddHumanAdjudication(state.Adjudication); err != nil {
			return nil, err
		}
	}

	return l, nil
}
