package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alexandrmotologa/consensai/internal/agent"
	"github.com/alexandrmotologa/consensai/internal/config"
	"github.com/alexandrmotologa/consensai/internal/ledger"
	"github.com/alexandrmotologa/consensai/internal/protocol"
	"github.com/alexandrmotologa/consensai/internal/report"
	"github.com/alexandrmotologa/consensai/internal/server"
)

const Version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	cmd := os.Args[1]
	switch cmd {
	case "debate":
		runDebateCmd(os.Args[2:])
	case "serve":
		runServeCmd(os.Args[2:])
	case "verify":
		runVerifyCmd(os.Args[2:])
	case "memo":
		runMemoCmd(os.Args[2:])
	case "version", "--version", "-v":
		fmt.Printf("consensai v%s (Go orchestrator, Epistemic Consensus Engine)\n", Version)
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Fprintf(os.Stderr, "Unknown subcommand %q. Run 'consensai help' for usage.\n", cmd)
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Printf(`consensai v%s - Multi-Agent Adversarial Debate & Epistemic Consensus Protocol

Usage:
  consensai <command> [flags]

Available Commands:
  debate     Run a 3-round dialectical debate over a high-stakes decision dossier
  serve      Launch the local HTTP server and embedded React 19 web cockpit (port 6060)
  verify     Audit cryptographic SHA-256 integrity of an audit ledger file
  memo       Generate an executive decision memo from a completed ledger
  version    Show version information

Flags for 'debate':
  --dossier   Path to decision dossier JSON (default: examples/dossiers/commercial_loan_500k.json)
  --config    Path to configuration YAML (default: examples/consensai.yaml)
  --provider  LLM provider: mock, openai, anthropic, gemini, ollama (default: mock)
  --model     Model name (e.g. gpt-4o, claude-3-5-sonnet, gemini-1.5-pro)
  --output    Path to export Markdown decision memo
  --ledger    Path to export JSON cryptographic ledger
  --quiet     Suppress real-time terminal stream output

Flags for 'serve':
  --port      HTTP server port (default: 6060)
  --config    Path to configuration YAML
  --dossiers  Directory containing decision dossiers (default: examples/dossiers)

Flags for 'verify':
  --ledger    Path to JSON ledger file to audit (required)

Flags for 'memo':
  --ledger    Path to JSON ledger file (required)
  --out       Path to save output memo (default: stdout)
`, Version)
}

func runDebateCmd(args []string) {
	fs := flag.NewFlagSet("debate", flag.ExitOnError)
	dossierPath := fs.String("dossier", "examples/dossiers/commercial_loan_500k.json", "Path to decision dossier JSON")
	configPath := fs.String("config", "examples/consensai.yaml", "Path to configuration YAML")
	providerName := fs.String("provider", "", "LLM provider (mock, openai, anthropic, gemini, ollama)")
	modelName := fs.String("model", "", "Model name")
	outputPath := fs.String("output", "", "Path to save Markdown decision memo")
	ledgerPath := fs.String("ledger", "", "Path to save JSON audit ledger")
	quiet := fs.Bool("quiet", false, "Suppress interactive terminal output")
	_ = fs.Parse(args)

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	if *providerName != "" {
		cfg.Provider.Default = *providerName
	}
	if *modelName != "" {
		cfg.Provider.Model = *modelName
	}

	dossierData, err := os.ReadFile(*dossierPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading dossier file %s: %v\n", *dossierPath, err)
		os.Exit(1)
	}

	var parsed struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}
	_ = json.Unmarshal(dossierData, &parsed)
	dossierID := parsed.ID
	if dossierID == "" {
		dossierID = strings.TrimSuffix(filepath.Base(*dossierPath), filepath.Ext(*dossierPath))
	}
	dossierTitle := parsed.Title
	if dossierTitle == "" {
		dossierTitle = dossierID
	}

	client, err := agent.NewProviderClient(&cfg.Provider)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating provider client: %v\n", err)
		os.Exit(1)
	}

	renderer := report.NewTerminalRenderer()
	if !*quiet {
		renderer.PrintBanner()
		fmt.Printf("Evaluating Dossier: %s (%s)\n", dossierTitle, *dossierPath)
		fmt.Printf("Inference Provider: %s\n\n", client.Name())
	}

	advocate := agent.NewAdvocateAgent(client)
	skeptic := agent.NewSkepticAgent(client)
	auditor := agent.NewAuditorAgent(client)
	engine := protocol.NewEngine(&cfg.Protocol, advocate, skeptic, auditor)

	state, err := engine.RunDeliberation(
		context.Background(),
		dossierID,
		dossierTitle,
		string(dossierData),
		func(evt protocol.EngineEvent) {
			if *quiet {
				return
			}
			switch evt.Type {
			case "phase_change":
				if text, ok := evt.Data.(string); ok {
					renderer.PrintPhase(text)
				}
			case "turn_complete":
				if turn, ok := evt.Data.(*protocol.AgentTurn); ok {
					renderer.PrintTurn(turn)
				}
			case "audit_complete":
				if findings, ok := evt.Data.([]protocol.AuditFinding); ok {
					renderer.PrintFindings(findings)
				}
			}
		},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Deliberation failed: %v\n", err)
		os.Exit(1)
	}

	// Create and commit cryptographic ledger
	chain, err := ledger.CreateFromState(state)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating audit ledger: %v\n", err)
		os.Exit(1)
	}

	headHash := chain.Blocks[len(chain.Blocks)-1].Hash
	if !*quiet {
		renderer.PrintSynthesis(state.Synthesis, headHash)
	}

	if *ledgerPath != "" {
		if err := chain.SaveToFile(*ledgerPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving ledger: %v\n", err)
		} else {
			fmt.Printf("Cryptographic audit ledger saved to %s\n", *ledgerPath)
		}
	}

	if *outputPath != "" {
		memo := report.NewDecisionMemo(state)
		md := memo.GenerateMarkdown(headHash)
		if err := os.WriteFile(*outputPath, []byte(md), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving decision memo: %v\n", err)
		} else {
			fmt.Printf("Executive decision memo saved to %s\n", *outputPath)
		}
	}
}

func runServeCmd(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	port := fs.Int("port", 6060, "HTTP server port")
	configPath := fs.String("config", "examples/consensai.yaml", "Path to configuration YAML")
	dossiersDir := fs.String("dossiers", "examples/dossiers", "Directory containing decision dossiers")
	_ = fs.Parse(args)

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning loading config, using defaults: %v\n", err)
		cfg = config.DefaultConfig()
	}

	if *port != 0 {
		cfg.Server.Port = *port
	}
	if *dossiersDir != "" {
		cfg.Server.DossiersDir = *dossiersDir
	}

	srv := server.NewServer(cfg)
	addr := fmt.Sprintf(":%d", cfg.Server.Port)

	fmt.Println("========================================================================")
	fmt.Printf(" consensai Server listening on http://localhost:%d\n", cfg.Server.Port)
	fmt.Printf(" Web Cockpit: http://localhost:%d\n", cfg.Server.Port)
	fmt.Printf(" Dossiers directory: %s\n", cfg.Server.DossiersDir)
	fmt.Printf(" Inference provider: %s\n", cfg.Provider.Default)
	fmt.Println("========================================================================")

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      srv.Routes(),
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeoutSeconds) * time.Second,
	}

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

func runVerifyCmd(args []string) {
	fs := flag.NewFlagSet("verify", flag.ExitOnError)
	ledgerPath := fs.String("ledger", "", "Path to JSON audit ledger file (required)")
	_ = fs.Parse(args)

	if *ledgerPath == "" {
		fmt.Fprintln(os.Stderr, "Error: --ledger flag is required")
		os.Exit(1)
	}

	l, err := ledger.LoadFromFile(*ledgerPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading ledger: %v\n", err)
		os.Exit(1)
	}

	cert := l.GenerateAuditCertificate()
	fmt.Println(cert)

	valid, err := l.Verify()
	if !valid || err != nil {
		fmt.Fprintf(os.Stderr, "FAILURE: Cryptographic ledger verification failed: %v\n", err)
		os.Exit(2)
	}

	fmt.Println("SUCCESS: All cryptographic blocks verified. Zero hash discrepancies or tampering detected.")
}

func runMemoCmd(args []string) {
	fs := flag.NewFlagSet("memo", flag.ExitOnError)
	ledgerPath := fs.String("ledger", "", "Path to JSON audit ledger file (required)")
	outPath := fs.String("out", "", "Path to save decision memo (default: stdout)")
	_ = fs.Parse(args)

	if *ledgerPath == "" {
		fmt.Fprintln(os.Stderr, "Error: --ledger flag is required")
		os.Exit(1)
	}

	l, err := ledger.LoadFromFile(*ledgerPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading ledger: %v\n", err)
		os.Exit(1)
	}

	// Extract state, claims, and findings from blocks
	var synthesis *protocol.SynthesisResult
	var allClaims []protocol.Claim
	findingMap := make(map[string]protocol.AuditFinding)

	for _, b := range l.Blocks {
		switch b.Type {
		case ledger.BlockTypeRound1:
			var p struct {
				Turns []protocol.AgentTurn `json:"turns"`
			}
			if err := json.Unmarshal(b.Payload, &p); err == nil {
				for _, t := range p.Turns {
					allClaims = append(allClaims, t.Claims...)
				}
			}
		case ledger.BlockTypeRound3:
			var p struct {
				AuditFindings []protocol.AuditFinding `json:"audit_findings"`
			}
			if err := json.Unmarshal(b.Payload, &p); err == nil {
				for _, f := range p.AuditFindings {
					findingMap[f.ClaimID] = f
				}
			}
		case ledger.BlockTypeSynthesis:
			var p struct {
				Synthesis *protocol.SynthesisResult `json:"synthesis"`
			}
			if err := json.Unmarshal(b.Payload, &p); err == nil {
				synthesis = p.Synthesis
			}
		}
	}

	for i := range allClaims {
		if f, exists := findingMap[allClaims[i].ID]; exists {
			allClaims[i].AuditedStatus = f.Status
			allClaims[i].AuditNotes = f.Discrepancy
		}
	}

	if synthesis == nil {
		fmt.Fprintln(os.Stderr, "No synthesis block found in ledger")
		os.Exit(1)
	}

	state := &protocol.DeliberationState{
		ID:            l.SessionID,
		DossierID:     l.DossierID,
		DossierTitle:  "Deliberation for " + l.DossierID,
		Synthesis:     synthesis,
		AllClaims:     allClaims,
		CreatedAt:     l.CreatedAt,
	}

	headHash := l.Blocks[len(l.Blocks)-1].Hash
	memo := report.NewDecisionMemo(state)
	md := memo.GenerateMarkdown(headHash)

	if *outPath != "" {
		if err := os.WriteFile(*outPath, []byte(md), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving memo: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Executive decision memo exported to %s\n", *outPath)
	} else {
		fmt.Println(md)
	}
}
