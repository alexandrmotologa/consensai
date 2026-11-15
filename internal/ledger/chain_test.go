package ledger

import (
	"encoding/json"
	"testing"
)

func TestLedger_ValidChainVerification(t *testing.T) {
	l, err := NewLedger("sess_test_1", "dossier_test", `{"test": true}`)
	if err != nil {
		t.Fatalf("failed creating ledger: %v", err)
	}

	_, err = l.AddBlock(BlockTypeRound1, 1, map[string]string{"claim": "Advocate thesis"})
	if err != nil {
		t.Fatalf("failed adding block 1: %v", err)
	}

	_, err = l.AddBlock(BlockTypeRound2, 2, map[string]string{"rebuttal": "Skeptic cross"})
	if err != nil {
		t.Fatalf("failed adding block 2: %v", err)
	}

	valid, err := l.Verify()
	if !valid || err != nil {
		t.Fatalf("expected ledger to be valid, got valid=%v, err=%v", valid, err)
	}
}

func TestLedger_TamperDetectionOnPayload(t *testing.T) {
	l, err := NewLedger("sess_test_2", "dossier_test", `{"test": true}`)
	if err != nil {
		t.Fatalf("failed creating ledger: %v", err)
	}

	_, err = l.AddBlock(BlockTypeRound1, 1, map[string]string{"data": "original"})
	if err != nil {
		t.Fatalf("failed adding block 1: %v", err)
	}

	// Tamper with block 1 payload without updating hash
	l.Blocks[1].Payload = json.RawMessage(`{"data":"malicious_tamper"}`)

	valid, err := l.Verify()
	if valid || err == nil {
		t.Fatalf("expected tamper detection to catch modified payload, but verification succeeded")
	}
}

func TestLedger_TamperDetectionOnBrokenChain(t *testing.T) {
	l, err := NewLedger("sess_test_3", "dossier_test", `{"test": true}`)
	if err != nil {
		t.Fatalf("failed creating ledger: %v", err)
	}

	_, err = l.AddBlock(BlockTypeRound1, 1, map[string]string{"turn": "1"})
	if err != nil {
		t.Fatalf("failed adding block 1: %v", err)
	}

	_, err = l.AddBlock(BlockTypeRound2, 2, map[string]string{"turn": "2"})
	if err != nil {
		t.Fatalf("failed adding block 2: %v", err)
	}

	// Break prev_hash linkage
	l.Blocks[2].PrevHash = "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"

	valid, err := l.Verify()
	if valid || err == nil {
		t.Fatalf("expected verification failure due to broken chain link, got valid=%v", valid)
	}
}
