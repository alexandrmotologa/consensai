package ledger

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// BlockType defines the semantic type of a cryptographic block.
type BlockType string

const (
	BlockTypeGenesis   BlockType = "GENESIS"
	BlockTypeRound1    BlockType = "ROUND_1_FRAMING"
	BlockTypeRound2    BlockType = "ROUND_2_REBUTTAL"
	BlockTypeRound3    BlockType = "ROUND_3_AUDIT"
	BlockTypeSynthesis         BlockType = "SYNTHESIS_DECISION"
	BlockTypeHumanAdjudication BlockType = "HUMAN_ADJUDICATION"
)

// Block represents an immutable cryptographic record in the audit chain.
type Block struct {
	Index       int             `json:"index"`
	TimestampMs int64           `json:"timestamp_ms"`
	Round       int             `json:"round"`
	Type        BlockType       `json:"type"`
	PrevHash    string          `json:"prev_hash"`
	Hash        string          `json:"hash"`
	Payload     json.RawMessage `json:"payload"`
}

// CanonicalRepresentation returns the deterministic byte slice used for hashing.
func (b *Block) CanonicalRepresentation() ([]byte, error) {
	type canonicalBlockData struct {
		Index       int             `json:"index"`
		TimestampMs int64           `json:"timestamp_ms"`
		Round       int             `json:"round"`
		Type        BlockType       `json:"type"`
		PrevHash    string          `json:"prev_hash"`
		Payload     json.RawMessage `json:"payload"`
	}

	bd := canonicalBlockData{
		Index:       b.Index,
		TimestampMs: b.TimestampMs,
		Round:       b.Round,
		Type:        b.Type,
		PrevHash:    b.PrevHash,
		Payload:     b.Payload,
	}

	return json.Marshal(bd)
}

// CalculateHash deterministically computes the SHA-256 hash of the block.
func (b *Block) CalculateHash() (string, error) {
	canonicalBytes, err := b.CanonicalRepresentation()
	if err != nil {
		return "", fmt.Errorf("canonicalizing block %d: %w", b.Index, err)
	}

	h := sha256.Sum256(canonicalBytes)
	return hex.EncodeToString(h[:]), nil
}

// FormatTimestamp returns a human-readable RFC3339 representation of the block timestamp.
func (b *Block) FormatTimestamp() string {
	return time.UnixMilli(b.TimestampMs).UTC().Format(time.RFC3339)
}
