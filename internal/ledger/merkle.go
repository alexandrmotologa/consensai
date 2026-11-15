package ledger

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// MerkleNode represents a single node in the cryptographic Merkle tree.
type MerkleNode struct {
	Left  *MerkleNode `json:"left,omitempty"`
	Right *MerkleNode `json:"right,omitempty"`
	Hash  string      `json:"hash"`
}

// MerkleTree computes a verifiable cryptographic tree over arbitrary leaf hashes.
type MerkleTree struct {
	Root   *MerkleNode `json:"root"`
	Leaves []string    `json:"leaves"`
}

// NewMerkleTree builds a balanced Merkle tree from an array of leaf string data or hashes.
func NewMerkleTree(data []string) (*MerkleTree, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("cannot build Merkle tree from empty dataset")
	}

	var leaves []string
	var nodes []*MerkleNode

	for _, item := range data {
		h := sha256.Sum256([]byte(item))
		hashStr := hex.EncodeToString(h[:])
		leaves = append(leaves, hashStr)
		nodes = append(nodes, &MerkleNode{Hash: hashStr})
	}

	// Duplicate last element if odd count to ensure binary balance
	if len(nodes)%2 != 0 {
		nodes = append(nodes, nodes[len(nodes)-1])
	}

	root := buildTreeLevel(nodes)
	return &MerkleTree{
		Root:   root,
		Leaves: leaves,
	}, nil
}

func buildTreeLevel(nodes []*MerkleNode) *MerkleNode {
	if len(nodes) == 1 {
		return nodes[0]
	}

	var nextLevel []*MerkleNode
	for i := 0; i < len(nodes); i += 2 {
		left := nodes[i]
		right := left
		if i+1 < len(nodes) {
			right = nodes[i+1]
		}

		combined := left.Hash + right.Hash
		h := sha256.Sum256([]byte(combined))
		parentHash := hex.EncodeToString(h[:])

		nextLevel = append(nextLevel, &MerkleNode{
			Left:  left,
			Right: right,
			Hash:  parentHash,
		})
	}

	return buildTreeLevel(nextLevel)
}

// RootHash returns the hex-encoded Merkle root hash.
func (m *MerkleTree) RootHash() string {
	if m.Root == nil {
		return ""
	}
	return m.Root.Hash
}

// GenerateProof produces an inclusion proof for a given leaf index.
func (m *MerkleTree) GenerateProof(leafIndex int) ([]string, error) {
	if leafIndex < 0 || leafIndex >= len(m.Leaves) {
		return nil, fmt.Errorf("leaf index %d out of bounds (0-%d)", leafIndex, len(m.Leaves)-1)
	}

	var proof []string
	currentLevel := make([]string, len(m.Leaves))
	copy(currentLevel, m.Leaves)

	idx := leafIndex
	for len(currentLevel) > 1 {
		if len(currentLevel)%2 != 0 {
			currentLevel = append(currentLevel, currentLevel[len(currentLevel)-1])
		}

		pairIdx := idx ^ 1
		proof = append(proof, currentLevel[pairIdx])

		var nextLevel []string
		for i := 0; i < len(currentLevel); i += 2 {
			combined := currentLevel[i] + currentLevel[i+1]
			h := sha256.Sum256([]byte(combined))
			nextLevel = append(nextLevel, hex.EncodeToString(h[:]))
		}

		idx /= 2
		currentLevel = nextLevel
	}

	return proof, nil
}

// VerifyProof verifies if a leaf hash belongs to the Merkle tree root using the proof path.
func VerifyProof(leafHash string, proof []string, rootHash string, index int) bool {
	current := leafHash
	idx := index

	for _, sibling := range proof {
		var combined string
		if idx%2 == 0 {
			combined = current + sibling
		} else {
			combined = sibling + current
		}
		h := sha256.Sum256([]byte(combined))
		current = hex.EncodeToString(h[:])
		idx /= 2
	}

	return current == rootHash
}
