export type AgentRole = 'Advocate' | 'Skeptic' | 'Risk Auditor' | 'Synthesizer';

export type ClaimStatus =
  | 'PENDING'
  | 'VERIFIED'
  | 'UNVERIFIED'
  | 'HALLUCINATED'
  | 'CONCEDED'
  | 'UNMITIGATED';

export type DecisionOutcome =
  | 'APPROVE'
  | 'REJECT'
  | 'ESCALATE TO HUMAN OVERSIGHT';

export interface Claim {
  id: string;
  agent_role: AgentRole;
  statement: string;
  category: string;
  citation: string;
  audited_status: ClaimStatus;
  audit_notes?: string;
}

export interface Rebuttal {
  target_claim_id: string;
  counter_point: string;
  citation?: string;
  concession: boolean;
}

export interface AgentTurn {
  round: number;
  role: AgentRole;
  timestamp: string;
  summary: string;
  raw_content: string;
  claims?: Claim[];
  rebuttals?: Rebuttal[];
  concessions?: string[];
}

export interface AuditFinding {
  claim_id: string;
  agent_role: AgentRole;
  statement: string;
  citation_given: string;
  found_in_dossier: boolean;
  discrepancy?: string;
  status: ClaimStatus;
  confidence_score: number;
}

export interface SynthesisResult {
  decision: DecisionOutcome;
  epistemic_convergence_score: number;
  normalized_score: number;
  confidence_interval: [number, number];
  total_claims_audited: number;
  verified_affirmative_claims: number;
  unmitigated_adversarial_risks: number;
  preserved_dissents: string[];
  action_rationale: string;
  executive_summary: string;
}

export interface Block {
  index: number;
  timestamp_ms: number;
  round: number;
  type: string;
  prev_hash: string;
  hash: string;
  payload: any;
}

export interface Ledger {
  session_id: string;
  dossier_id: string;
  dossier_hash: string;
  created_at: string;
  blocks: Block[];
}

export interface DossierSummary {
  id: string;
  title: string;
  domain: string;
  risk_tier: string;
  file_path: string;
}

export interface SystemStatus {
  status: string;
  version: string;
  uptime_seconds: number;
  allocated_memory_mb: number;
  sys_memory_mb: number;
  goroutines: number;
  active_sessions: number;
}

export interface EngineEvent {
  type: string;
  session_id: string;
  round: number;
  role?: AgentRole;
  timestamp: string;
  data: any;
}
