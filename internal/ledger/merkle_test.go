package ledger

import (
	"testing"
)

func TestMerkleTree_BuildAndVerifyProof(t *testing.T) {
	data := []string{
		"claim_1:Advocate:Statement A:VERIFIED",
		"claim_2:Advocate:Statement B:VERIFIED",
		"claim_3:Skeptic:Risk C:UNMITIGATED",
		"claim_4:Skeptic:Risk D:CONCEDED",
	}

	tree, err := NewMerkleTree(data)
	if err != nil {
		t.Fatalf("failed building Merkle tree: %v", err)
	}

	root := tree.RootHash()
	if root == "" {
		t.Fatalf("expected non-empty root hash")
	}

	// Verify proof for leaf 2
	proof, err := tree.GenerateProof(2)
	if err != nil {
		t.Fatalf("failed generating proof: %v", err)
	}

	valid := VerifyProof(tree.Leaves[2], proof, root, 2)
	if !valid {
		t.Fatalf("expected valid proof verification for leaf 2")
	}

	// Verify failure on fake root
	fakeValid := VerifyProof(tree.Leaves[2], proof, "0000000000000000000000000000000000000000000000000000000000000000", 2)
	if fakeValid {
		t.Fatalf("expected proof verification to fail against fake root")
	}
}
