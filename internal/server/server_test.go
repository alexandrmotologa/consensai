package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alexandrmotologa/consensai/internal/config"
)

func TestServer_Integration(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Server.DossiersDir = "../../examples/dossiers"

	srv := NewServer(cfg)
	handler := srv.Routes()

	// 1. Test GET /api/status
	reqStatus := httptest.NewRequest("GET", "/api/status", nil)
	recStatus := httptest.NewRecorder()
	handler.ServeHTTP(recStatus, reqStatus)

	if recStatus.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", recStatus.Code, recStatus.Body.String())
	}

	var statusMap map[string]interface{}
	if err := json.Unmarshal(recStatus.Body.Bytes(), &statusMap); err != nil {
		t.Fatalf("unmarshaling status: %v", err)
	}
	if statusMap["status"] != "online" {
		t.Fatalf("expected status 'online', got %v", statusMap["status"])
	}

	// 2. Test GET /api/dossiers
	reqDossiers := httptest.NewRequest("GET", "/api/dossiers", nil)
	recDossiers := httptest.NewRecorder()
	handler.ServeHTTP(recDossiers, reqDossiers)

	if recDossiers.Code != http.StatusOK {
		t.Fatalf("expected dossiers 200, got %d", recDossiers.Code)
	}

	var dossiers []map[string]interface{}
	if err := json.Unmarshal(recDossiers.Body.Bytes(), &dossiers); err != nil {
		t.Fatalf("unmarshaling dossiers: %v", err)
	}
	if len(dossiers) == 0 {
		t.Fatalf("expected at least 1 dossier, got 0")
	}

	// 3. Test POST /api/deliberate with sample dossier
	delibReqPayload := map[string]interface{}{
		"dossier_id": "commercial_loan_500k",
		"provider":   "mock",
	}
	delibReqBytes, _ := json.Marshal(delibReqPayload)
	reqDelib := httptest.NewRequest("POST", "/api/deliberate", bytes.NewReader(delibReqBytes))
	recDelib := httptest.NewRecorder()
	handler.ServeHTTP(recDelib, reqDelib)

	if recDelib.Code != http.StatusOK {
		t.Fatalf("expected deliberate 200, got %d: %s", recDelib.Code, recDelib.Body.String())
	}

	// Wait briefly for background goroutine deliberation to commit
	time.Sleep(300 * time.Millisecond)

	// 4. Test GET / (embedded React UI)
	reqIndex := httptest.NewRequest("GET", "/", nil)
	recIndex := httptest.NewRecorder()
	handler.ServeHTTP(recIndex, reqIndex)

	if recIndex.Code != http.StatusOK {
		t.Fatalf("expected index 200, got %d", recIndex.Code)
	}
	if !strings.Contains(recIndex.Body.String(), "consensai") {
		t.Fatalf("expected index.html to contain 'consensai', got: %s", recIndex.Body.String())
	}
}
