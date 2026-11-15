package ledger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ExportJSON serializes the ledger to formatted JSON bytes.
func (l *Ledger) ExportJSON() ([]byte, error) {
	return json.MarshalIndent(l, "", "  ")
}

// SaveToFile writes the ledger to the specified path, creating directories if needed.
func (l *Ledger) SaveToFile(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating directory %s: %w", dir, err)
	}

	data, err := l.ExportJSON()
	if err != nil {
		return fmt.Errorf("serializing ledger: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing ledger file: %w", err)
	}

	return nil
}

// LoadFromFile reads and parses a ledger file from disk.
func LoadFromFile(path string) (*Ledger, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading ledger file %s: %w", path, err)
	}

	var l Ledger
	if err := json.Unmarshal(data, &l); err != nil {
		return nil, fmt.Errorf("parsing ledger JSON: %w", err)
	}

	return &l, nil
}

// GenerateAuditCertificate creates a formal cryptographic proof certificate.
func (l *Ledger) GenerateAuditCertificate() string {
	valid, err := l.Verify()
	statusText := "VERIFIED IMMUTABLE"
	if !valid || err != nil {
		statusText = fmt.Sprintf("TAMPER DETECTED: %v", err)
	}

	headHash := "NONE"
	if len(l.Blocks) > 0 {
		headHash = l.Blocks[len(l.Blocks)-1].Hash
	}

	var sb strings.Builder
	sb.WriteString("================================================================================\n")
	sb.WriteString("                 CONSENSAI CRYPTOGRAPHIC AUDIT CERTIFICATE                      \n")
	sb.WriteString("================================================================================\n\n")
	sb.WriteString(fmt.Sprintf("Session ID:        %s\n", l.SessionID))
	sb.WriteString(fmt.Sprintf("Dossier ID:        %s\n", l.DossierID))
	sb.WriteString(fmt.Sprintf("Dossier SHA-256:   %s\n", l.DossierHash))
	sb.WriteString(fmt.Sprintf("Ledger Created At: %s\n", l.CreatedAt.Format(time.RFC3339)))
	sb.WriteString(fmt.Sprintf("Block Count:       %d\n", len(l.Blocks)))
	sb.WriteString(fmt.Sprintf("Head Block Hash:   %s\n", headHash))
	sb.WriteString(fmt.Sprintf("Integrity Status:  %s\n\n", statusText))
	sb.WriteString("--------------------------------------------------------------------------------\n")
	sb.WriteString("BLOCK SEQUENCE & HASH CHAIN:\n")
	sb.WriteString("--------------------------------------------------------------------------------\n")

	for _, b := range l.Blocks {
		sb.WriteString(fmt.Sprintf("Block #%d | Type: %-22s | Round: %d\n", b.Index, b.Type, b.Round))
		sb.WriteString(fmt.Sprintf("  Timestamp: %s\n", b.FormatTimestamp()))
		sb.WriteString(fmt.Sprintf("  Prev Hash: %s\n", b.PrevHash))
		sb.WriteString(fmt.Sprintf("  Hash:      %s\n", b.Hash))
		sb.WriteString("\n")
	}

	sb.WriteString("================================================================================\n")
	return sb.String()
}
