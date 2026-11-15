package protocol

import (
	"time"
)

// AgentRole defines the role of an agent in the dialectic debate.
type AgentRole string

const (
	RoleAdvocate   AgentRole = "Advocate"
	RoleSkeptic    AgentRole = "Skeptic"
	RoleAuditor    AgentRole = "Risk Auditor"
	RoleSynthesizer AgentRole = "Synthesizer"
)

// ClaimStatus defines the verification status of a factual assertion.
type ClaimStatus string

const (
	StatusPending     ClaimStatus = "PENDING"
	StatusVerified    ClaimStatus = "VERIFIED"
	StatusUnverified  ClaimStatus = "UNVERIFIED"
	StatusHallucinated ClaimStatus = "HALLUCINATED"
	StatusConceded    ClaimStatus = "CONCEDED"
	StatusUnmitigated ClaimStatus = "UNMITIGATED"
)

// Claim represents a structured factual assertion made during a debate turn.
type Claim struct {
	ID            string      `json:"id"`
	AgentRole     AgentRole   `json:"agent_role"`
	Statement     string      `json:"statement"`
	Category      string      `json:"category"` // e.g. financial, risk, compliance, operational
	Citation      string      `json:"citation"` // exact citation or quote from source dossier
	AuditedStatus ClaimStatus `json:"audited_status"`
	AuditNotes    string      `json:"audit_notes,omitempty"`
}

// Rebuttal represents a counter-argument to a specific claim.
type Rebuttal struct {
	TargetClaimID string `json:"target_claim_id"`
	CounterPoint  string `json:"counter_point"`
	Citation      string `json:"citation,omitempty"`
	Concession    bool   `json:"concession"`
}

// AgentTurn represents the complete output from an agent during one round.
type AgentTurn struct {
	Round       int        `json:"round"`
	Role        AgentRole  `json:"role"`
	Timestamp   time.Time  `json:"timestamp"`
	Summary     string     `json:"summary"`
	RawContent  string     `json:"raw_content"`
	Claims      []Claim    `json:"claims,omitempty"`
	Rebuttals   []Rebuttal `json:"rebuttals,omitempty"`
	Concessions []string   `json:"concessions,omitempty"`
}

// AuditFinding stores the result of verifying an assertion against the dossier.
type AuditFinding struct {
	ClaimID         string      `json:"claim_id"`
	AgentRole       AgentRole   `json:"agent_role"`
	Statement       string      `json:"statement"`
	CitationGiven   string      `json:"citation_given"`
	FoundInDossier  bool        `json:"found_in_dossier"`
	Discrepancy     string      `json:"discrepancy,omitempty"`
	Status          ClaimStatus `json:"status"`
	ConfidenceScore float64     `json:"confidence_score"`
}

// DecisionOutcome defines the three final governance actions.
type DecisionOutcome string

const (
	OutcomeApprove  DecisionOutcome = "APPROVE"
	OutcomeReject   DecisionOutcome = "REJECT"
	OutcomeEscalate DecisionOutcome = "ESCALATE TO HUMAN OVERSIGHT"
)

// ComplianceCheck represents an automated statutory or policy audit verification.
type ComplianceCheck struct {
	Standard    string `json:"standard"` // e.g. EU_AI_ACT_ART_14, BASEL_III_DSCR, FCRA_ADVERSE_ACTION, HIPAA_45_CFR_164
	Description string `json:"description"`
	Status      string `json:"status"` // COMPLIANT, FLAGGED, REQUIRES_HUMAN_ACTION, NOT_APPLICABLE
	Notes       string `json:"notes"`
}

// ComplianceReport summarizes regulatory conformance.
type ComplianceReport struct {
	OverallStatus string            `json:"overall_status"`
	Checks        []ComplianceCheck `json:"checks"`
}

// HumanAdjudication records an executive human officer's signed override or conditional approval.
type HumanAdjudication struct {
	OfficerID     string    `json:"officer_id"`
	OfficerName   string    `json:"officer_name"`
	OfficerRole   string    `json:"officer_role"`
	Decision      string    `json:"decision"` // CONDITIONAL_APPROVE, OVERRIDE_APPROVE, OVERRIDE_REJECT, AFFIRM_REJECT
	Justification string    `json:"justification"`
	Conditions    []string  `json:"conditions,omitempty"`
	SignatureHash string    `json:"signature_hash"`
	Timestamp     time.Time `json:"timestamp"`
}

// SynthesisResult stores the deterministic consensus output.
type SynthesisResult struct {
	Decision                  DecisionOutcome   `json:"decision"`
	EpistemicConvergenceScore float64           `json:"epistemic_convergence_score"` // raw (-1.0 to 1.0)
	NormalizedScore           float64           `json:"normalized_score"`            // scaled (0.0 to 1.0)
	ConfidenceInterval        [2]float64        `json:"confidence_interval"`
	TotalClaimsAudited        int               `json:"total_claims_audited"`
	VerifiedAffirmativeClaims int               `json:"verified_affirmative_claims"`
	UnmitigatedAdversarialRisks int             `json:"unmitigated_adversarial_risks"`
	PreservedDissents         []string          `json:"preserved_dissents"`
	ActionRationale           string            `json:"action_rationale"`
	ExecutiveSummary          string            `json:"executive_summary"`
	Compliance                *ComplianceReport `json:"compliance,omitempty"`
	MerkleRoot                string            `json:"merkle_root,omitempty"`
}

// DeliberationState tracks the progression of an entire decision workflow.
type DeliberationState struct {
	ID            string             `json:"id"`
	DossierID     string             `json:"dossier_id"`
	DossierTitle  string             `json:"dossier_title"`
	DossierRaw    string             `json:"dossier_raw"`
	CurrentRound  int                `json:"current_round"`
	CurrentPhase  string             `json:"current_phase"`
	Turns         []AgentTurn        `json:"turns"`
	AllClaims     []Claim            `json:"all_claims"`
	AuditFindings []AuditFinding     `json:"audit_findings"`
	Synthesis     *SynthesisResult   `json:"synthesis,omitempty"`
	Adjudication  *HumanAdjudication `json:"adjudication,omitempty"`
	CreatedAt     time.Time          `json:"created_at"`
	CompletedAt   *time.Time         `json:"completed_at,omitempty"`
}
