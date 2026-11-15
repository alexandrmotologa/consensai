import React, { useState, useEffect, useRef } from 'react';
import { Header } from './components/Header';
import { ConvergenceGauge } from './components/ConvergenceGauge';
import { DeliberationStream } from './components/DeliberationStream';
import { ClaimMatrix } from './components/ClaimMatrix';
import { DissentCard } from './components/DissentCard';
import { LedgerVerifier } from './components/LedgerVerifier';
import { MemoModal } from './components/MemoModal';
import { DossierModal } from './components/DossierModal';
import { ComplianceCard } from './components/ComplianceCard';
import { HumanAdjudicationModal } from './components/HumanAdjudicationModal';
import { SensitivityAnalyzer } from './components/SensitivityAnalyzer';
import { DossierBuilder } from './components/DossierBuilder';
import { ModelArenaModal } from './components/ModelArenaModal';
import {
  DossierSummary,
  SystemStatus,
  AgentTurn,
  Claim,
  SynthesisResult,
  Ledger,
  AgentRole,
  EngineEvent,
  ModelArenaConfig,
  HumanAdjudication,
} from './types';

export const App: React.FC = () => {
  const [dossiers, setDossiers] = useState<DossierSummary[]>([]);
  const [selectedDossier, setSelectedDossier] = useState<string>('commercial_loan_500k');
  const [systemStatus, setSystemStatus] = useState<SystemStatus | null>(null);
  const [isRunning, setIsRunning] = useState<boolean>(false);
  const [activeSessionId, setActiveSessionId] = useState<string | null>(null);

  // Deliberation state
  const [currentRound, setCurrentRound] = useState<number>(0);
  const [currentPhase, setCurrentPhase] = useState<string>('READY');
  const [activeRole, setActiveRole] = useState<AgentRole | undefined>(undefined);
  const [turns, setTurns] = useState<AgentTurn[]>([]);
  const [claims, setClaims] = useState<Claim[]>([]);
  const [synthesis, setSynthesis] = useState<SynthesisResult | null>(null);
  const [ledger, setLedger] = useState<Ledger | null>(null);
  const [adjudication, setAdjudication] = useState<HumanAdjudication | null>(null);

  // Modals & Panels
  const [isMemoModalOpen, setIsMemoModalOpen] = useState<boolean>(false);
  const [isDossierModalOpen, setIsDossierModalOpen] = useState<boolean>(false);
  const [isArenaModalOpen, setIsArenaModalOpen] = useState<boolean>(false);
  const [isAdjudicationModalOpen, setIsAdjudicationModalOpen] = useState<boolean>(false);
  const [memoMarkdown, setMemoMarkdown] = useState<string>('');
  const [dossierContent, setDossierContent] = useState<string>('');

  // Model Arena Config
  const [arenaConfig, setArenaConfig] = useState<ModelArenaConfig>({
    advocate_provider: 'mock',
    advocate_model: 'heuristic-v1',
    skeptic_provider: 'mock',
    skeptic_model: 'heuristic-v1',
    auditor_provider: 'mock',
    auditor_model: 'heuristic-v1',
  });

  // Active view tab
  const [activeTab, setActiveTab] = useState<
    'deliberation' | 'matrix' | 'compliance' | 'builder' | 'ledger'
  >('deliberation');

  const sseRef = useRef<EventSource | null>(null);

  useEffect(() => {
    fetchDossiers();
    fetchSystemStatus();
    const interval = setInterval(fetchSystemStatus, 5000);
    return () => clearInterval(interval);
  }, []);

  useEffect(() => {
    if (selectedDossier) {
      fetchDossierContent(selectedDossier);
    }
  }, [selectedDossier]);

  const fetchDossiers = async () => {
    try {
      const res = await fetch('/api/dossiers');
      if (res.ok) {
        const data = await res.json();
        setDossiers(data);
        if (data.length > 0 && !selectedDossier) {
          setSelectedDossier(data[0].id);
        }
      }
    } catch (err) {
      console.error('Failed to fetch dossiers:', err);
    }
  };

  const fetchSystemStatus = async () => {
    try {
      const res = await fetch('/api/status');
      if (res.ok) {
        const data = await res.json();
        setSystemStatus(data);
      }
    } catch (err) {
      // Backend starting or offline
    }
  };

  const fetchDossierContent = async (id: string) => {
    try {
      const res = await fetch(`/api/dossiers/${id}`);
      if (res.ok) {
        const data = await res.text();
        setDossierContent(data);
      }
    } catch (err) {
      console.error('Failed to fetch dossier content:', err);
    }
  };

  const startDeliberation = async (customPayload?: any) => {
    if (isRunning) return;

    setIsRunning(true);
    setTurns([]);
    setClaims([]);
    setSynthesis(null);
    setLedger(null);
    setAdjudication(null);
    setCurrentRound(0);
    setCurrentPhase('INITIALIZING');
    setActiveRole(undefined);
    setMemoMarkdown('');

    try {
      const payload: any = {
        advocate_provider: arenaConfig.advocate_provider,
        advocate_model: arenaConfig.advocate_model,
        skeptic_provider: arenaConfig.skeptic_provider,
        skeptic_model: arenaConfig.skeptic_model,
        auditor_provider: arenaConfig.auditor_provider,
        auditor_model: arenaConfig.auditor_model,
      };

      if (customPayload) {
        payload.custom_dossier = customPayload;
      } else {
        payload.dossier_id = selectedDossier;
      }

      const res = await fetch('/api/deliberate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });

      if (!res.ok) {
        throw new Error('Failed to start deliberation');
      }

      connectSSE();
      setActiveTab('deliberation');
    } catch (err) {
      console.error('Error starting deliberation:', err);
      setIsRunning(false);
    }
  };

  const connectSSE = () => {
    if (sseRef.current) {
      sseRef.current.close();
    }

    const sse = new EventSource(`/api/deliberate/active/stream`);
    sseRef.current = sse;

    sse.addEventListener('session_start', (e: MessageEvent) => {
      const data = JSON.parse(e.data);
      if (data.data?.session_id) {
        setActiveSessionId(data.data.session_id);
      }
    });

    sse.addEventListener('phase_change', (e: MessageEvent) => {
      const data: EngineEvent = JSON.parse(e.data);
      setCurrentRound(data.round);
      if (typeof data.data === 'string') {
        setCurrentPhase(data.data);
      }
    });

    sse.addEventListener('turn_start', (e: MessageEvent) => {
      const data: EngineEvent = JSON.parse(e.data);
      setActiveRole(data.role);
    });

    sse.addEventListener('turn_complete', (e: MessageEvent) => {
      const data: EngineEvent = JSON.parse(e.data);
      const turn = data.data as AgentTurn;
      setTurns((prev) => [...prev, turn]);
      if (turn.claims) {
        setClaims((prev) => [...prev, ...turn.claims!]);
      }
    });

    sse.addEventListener('audit_complete', (e: MessageEvent) => {
      const data: EngineEvent = JSON.parse(e.data);
      const findings = data.data as any[];
      const findingMap = new Map();
      findings.forEach((f) => findingMap.set(f.claim_id, f));

      setClaims((prev) =>
        prev.map((c) => {
          if (findingMap.has(c.id)) {
            const f = findingMap.get(c.id);
            return {
              ...c,
              audited_status: f.status,
              audit_notes: f.discrepancy,
            };
          }
          return c;
        })
      );
    });

    sse.addEventListener('synthesis_complete', (e: MessageEvent) => {
      const data: EngineEvent = JSON.parse(e.data);
      const synth = data.data as SynthesisResult;
      setSynthesis(synth);
    });

    sse.addEventListener('ledger_committed', async (e: MessageEvent) => {
      const data: EngineEvent = JSON.parse(e.data);
      const sessId = data.session_id;
      if (sessId) {
        setActiveSessionId(sessId);
        try {
          const lRes = await fetch(`/api/deliberate/${sessId}/ledger`);
          if (lRes.ok) {
            const lData = await lRes.json();
            setLedger(lData);
          }
          const mRes = await fetch(`/api/deliberate/${sessId}/memo`);
          if (mRes.ok) {
            const mData = await mRes.text();
            setMemoMarkdown(mData);
          }
        } catch (err) {
          console.error('Error fetching ledger or memo:', err);
        }
      }
    });

    sse.addEventListener('adjudication_committed', async (e: MessageEvent) => {
      const data: EngineEvent = JSON.parse(e.data);
      if (data.data?.adjudication) {
        setAdjudication(data.data.adjudication);
      }
      if (activeSessionId) {
        const lRes = await fetch(`/api/deliberate/${activeSessionId}/ledger`);
        if (lRes.ok) {
          setLedger(await lRes.json());
        }
        const mRes = await fetch(`/api/deliberate/${activeSessionId}/memo`);
        if (mRes.ok) {
          setMemoMarkdown(await mRes.text());
        }
      }
    });

    sse.addEventListener('session_completed', async (e: MessageEvent) => {
      setIsRunning(false);
      setCurrentPhase('COMPLETED');
      setActiveRole(undefined);
      try {
        const data = JSON.parse(e.data);
        const sid = data.session_id || data.data?.id;
        if (sid) {
          setActiveSessionId(sid);
          const lRes = await fetch(`/api/deliberate/${sid}/ledger`);
          if (lRes.ok) setLedger(await lRes.json());
          const mRes = await fetch(`/api/deliberate/${sid}/memo`);
          if (mRes.ok) setMemoMarkdown(await mRes.text());
        }
      } catch (err) {
        console.error('Error handling session_completed:', err);
      }
      if (sseRef.current) {
        sseRef.current.close();
      }
    });

    sse.onerror = () => {
      setIsRunning(false);
      if (sseRef.current) {
        sseRef.current.close();
      }
    };
  };

  const handleExportCertificate = () => {
    if (activeSessionId) {
      window.open(`/api/deliberate/${activeSessionId}/certificate`, '_blank');
    }
  };

  const handleAdjudicated = async (data: any) => {
    if (data.adjudication) {
      setAdjudication(data.adjudication);
    }
    if (activeSessionId) {
      const lRes = await fetch(`/api/deliberate/${activeSessionId}/ledger`);
      if (lRes.ok) {
        setLedger(await lRes.json());
      }
      const mRes = await fetch(`/api/deliberate/${activeSessionId}/memo`);
      if (mRes.ok) {
        setMemoMarkdown(await mRes.text());
      }
    }
  };

  const activeDossierTitle =
    dossiers.find((d) => d.id === selectedDossier)?.title || selectedDossier;

  return (
    <div className="min-h-screen bg-background text-slate-100 flex flex-col font-sans">
      {/* Top Header */}
      <Header
        dossiers={dossiers}
        selectedDossier={selectedDossier}
        onSelectDossier={setSelectedDossier}
        onRunDeliberation={() => startDeliberation()}
        isRunning={isRunning}
        systemStatus={systemStatus}
        onOpenDossierModal={() => setIsDossierModalOpen(true)}
        onOpenMemoModal={() => setIsMemoModalOpen(true)}
        onOpenArenaModal={() => setIsArenaModalOpen(true)}
        onOpenAdjudicationModal={() => setIsAdjudicationModalOpen(true)}
        onExportCertificate={handleExportCertificate}
        hasCompletedMemo={Boolean(memoMarkdown)}
        canAdjudicate={Boolean(activeSessionId && synthesis)}
      />

      {/* Main Workspace Layout */}
      <main className="flex-1 p-6 space-y-6 max-w-7xl mx-auto w-full">
        {/* Top Metric Cards Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2">
            <ConvergenceGauge
              synthesis={synthesis}
              isRunning={isRunning}
              currentPhase={currentPhase}
            />
          </div>
          <div className="space-y-4">
            <DissentCard dissents={synthesis?.preserved_dissents || []} />
            {adjudication && (
              <div className="rounded-xl border border-indigo-500/60 bg-gradient-to-br from-indigo-950/40 to-surface p-4 shadow-lg text-xs space-y-2">
                <div className="flex items-center justify-between border-b border-indigo-500/40 pb-2">
                  <span className="font-mono font-bold text-indigo-300 uppercase">
                    Human Officer Adjudicated
                  </span>
                  <span className="rounded bg-indigo-950 px-2 py-0.5 text-[10px] font-bold text-indigo-200 border border-indigo-700">
                    {adjudication.decision}
                  </span>
                </div>
                <div className="text-slate-300 font-medium">
                  {adjudication.officer_name} ({adjudication.officer_role})
                </div>
                <p className="text-slate-400 text-[11px] italic">
                  "{adjudication.justification}"
                </p>
                <div className="text-[10px] font-mono text-indigo-400 truncate">
                  SIG: {adjudication.signature_hash}
                </div>
              </div>
            )}
          </div>
        </div>

        {/* View Switcher Tabs */}
        <div className="flex flex-wrap items-center gap-2 border-b border-border pb-2 text-xs font-mono">
          <button
            onClick={() => setActiveTab('deliberation')}
            className={`px-3 py-1.5 rounded-lg transition-colors font-medium ${
              activeTab === 'deliberation'
                ? 'bg-cyan-950 text-cyan-300 border border-cyan-800'
                : 'text-slate-400 hover:text-slate-200'
            }`}
          >
            Deliberation Feed (3 Rounds)
          </button>

          <button
            onClick={() => setActiveTab('matrix')}
            className={`px-3 py-1.5 rounded-lg transition-colors font-medium ${
              activeTab === 'matrix'
                ? 'bg-cyan-950 text-cyan-300 border border-cyan-800'
                : 'text-slate-400 hover:text-slate-200'
            }`}
          >
            Claim Grounding Matrix ({claims.length})
          </button>

          <button
            onClick={() => setActiveTab('compliance')}
            className={`px-3 py-1.5 rounded-lg transition-colors font-medium ${
              activeTab === 'compliance'
                ? 'bg-cyan-950 text-cyan-300 border border-cyan-800'
                : 'text-slate-400 hover:text-slate-200'
            }`}
          >
            Regulatory Compliance ({synthesis?.compliance?.checks.length || 0})
          </button>

          <button
            onClick={() => setActiveTab('builder')}
            className={`px-3 py-1.5 rounded-lg transition-colors font-medium ${
              activeTab === 'builder'
                ? 'bg-cyan-950 text-cyan-300 border border-cyan-800'
                : 'text-slate-400 hover:text-slate-200'
            }`}
          >
            Dossier Builder / Ingest
          </button>

          <button
            onClick={() => setActiveTab('ledger')}
            className={`px-3 py-1.5 rounded-lg transition-colors font-medium ${
              activeTab === 'ledger'
                ? 'bg-cyan-950 text-cyan-300 border border-cyan-800'
                : 'text-slate-400 hover:text-slate-200'
            }`}
          >
            Cryptographic SHA-256 Ledger
          </button>
        </div>

        {/* Active Tab Content */}
        {activeTab === 'deliberation' && (
          <div className="space-y-6">
            <DeliberationStream
              turns={turns}
              activeRole={activeRole}
              currentRound={currentRound}
            />

            {activeSessionId && synthesis && (
              <SensitivityAnalyzer
                sessionId={activeSessionId}
                synthesis={synthesis}
                onRecalculated={setSynthesis}
              />
            )}
          </div>
        )}

        {activeTab === 'matrix' && <ClaimMatrix claims={claims} />}

        {activeTab === 'compliance' && <ComplianceCard report={synthesis?.compliance} />}

        {activeTab === 'builder' && (
          <DossierBuilder
            onRunCustom={(custom) => startDeliberation(custom)}
            isRunning={isRunning}
          />
        )}

        {activeTab === 'ledger' && (
          <LedgerVerifier
            ledger={ledger}
            merkleRoot={synthesis?.merkle_root}
          />
        )}
      </main>

      {/* Footer */}
      <footer className="border-t border-border bg-surface px-6 py-3 text-center text-xs text-slate-500 font-mono">
        consensai | High-Stakes Decision Protocol | Multi-Agent Epistemic Consensus &bull; Resident Memory &lt; 25MB RAM
      </footer>

      {/* Modals */}
      <MemoModal
        isOpen={isMemoModalOpen}
        onClose={() => setIsMemoModalOpen(false)}
        memoMarkdown={memoMarkdown}
      />

      <DossierModal
        isOpen={isDossierModalOpen}
        onClose={() => setIsDossierModalOpen(false)}
        dossierContent={dossierContent}
        dossierTitle={activeDossierTitle}
      />

      <ModelArenaModal
        isOpen={isArenaModalOpen}
        onClose={() => setIsArenaModalOpen(false)}
        config={arenaConfig}
        onSave={setArenaConfig}
      />

      <HumanAdjudicationModal
        isOpen={isAdjudicationModalOpen}
        onClose={() => setIsAdjudicationModalOpen(false)}
        sessionId={activeSessionId || ''}
        onAdjudicated={handleAdjudicated}
      />
    </div>
  );
};
