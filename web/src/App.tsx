import React, { useState, useEffect, useRef } from 'react';
import { Header } from './components/Header';
import { ConvergenceGauge } from './components/ConvergenceGauge';
import { DeliberationStream } from './components/DeliberationStream';
import { ClaimMatrix } from './components/ClaimMatrix';
import { DissentCard } from './components/DissentCard';
import { LedgerVerifier } from './components/LedgerVerifier';
import { MemoModal } from './components/MemoModal';
import { DossierModal } from './components/DossierModal';
import {
  DossierSummary,
  SystemStatus,
  AgentTurn,
  Claim,
  SynthesisResult,
  Ledger,
  AgentRole,
  EngineEvent,
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

  // Modals
  const [isMemoModalOpen, setIsMemoModalOpen] = useState<boolean>(false);
  const [isDossierModalOpen, setIsDossierModalOpen] = useState<boolean>(false);
  const [memoMarkdown, setMemoMarkdown] = useState<string>('');
  const [dossierContent, setDossierContent] = useState<string>('');

  // Tabs for main content view
  const [activeTab, setActiveTab] = useState<'deliberation' | 'matrix' | 'ledger'>('deliberation');

  const sseRef = useRef<EventSource | null>(null);

  // Fetch initial dossiers and system status
  useEffect(() => {
    fetchDossiers();
    fetchSystemStatus();
    const interval = setInterval(fetchSystemStatus, 5000);
    return () => clearInterval(interval);
  }, []);

  // Update dossier content when selected dossier changes
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
      // Backend may be starting or offline
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

  const startDeliberation = async () => {
    if (isRunning) return;

    // Reset current state
    setIsRunning(true);
    setTurns([]);
    setClaims([]);
    setSynthesis(null);
    setLedger(null);
    setCurrentRound(0);
    setCurrentPhase('INITIALIZING');
    setActiveRole(undefined);
    setMemoMarkdown('');

    try {
      const res = await fetch('/api/deliberate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          dossier_id: selectedDossier,
          provider: 'mock',
        }),
      });

      if (!res.ok) {
        throw new Error('Failed to start deliberation');
      }

      // Connect to SSE stream
      // We listen to the global broadcast stream or poll session
      connectSSE();
    } catch (err) {
      console.error('Error starting deliberation:', err);
      setIsRunning(false);
    }
  };

  const connectSSE = () => {
    if (sseRef.current) {
      sseRef.current.close();
    }

    // Connect to SSE stream
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
      // Update claim statuses with audit findings
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
        // Fetch full committed ledger
        try {
          const lRes = await fetch(`/api/deliberate/${sessId}/ledger`);
          if (lRes.ok) {
            const lData = await lRes.json();
            setLedger(lData);
          }
          // Fetch memo
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

    sse.addEventListener('session_completed', () => {
      setIsRunning(false);
      setCurrentPhase('COMPLETED');
      setActiveRole(undefined);
      if (sseRef.current) {
        sseRef.current.close();
      }
    });

    sse.onerror = () => {
      // Stream closed or error
      setIsRunning(false);
      if (sseRef.current) {
        sseRef.current.close();
      }
    };
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
        onRunDeliberation={startDeliberation}
        isRunning={isRunning}
        systemStatus={systemStatus}
        onOpenDossierModal={() => setIsDossierModalOpen(true)}
        onOpenMemoModal={() => setIsMemoModalOpen(true)}
        hasCompletedMemo={Boolean(memoMarkdown)}
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
          <div>
            <DissentCard dissents={synthesis?.preserved_dissents || []} />
          </div>
        </div>

        {/* View Switcher Tabs */}
        <div className="flex items-center space-x-2 border-b border-border pb-2 text-xs font-mono">
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
          <DeliberationStream
            turns={turns}
            activeRole={activeRole}
            currentRound={currentRound}
          />
        )}

        {activeTab === 'matrix' && <ClaimMatrix claims={claims} />}

        {activeTab === 'ledger' && <LedgerVerifier ledger={ledger} />}
      </main>

      {/* Footer */}
      <footer className="border-t border-border bg-surface px-6 py-3 text-center text-xs text-slate-500 font-mono">
        consensai | High-Stakes Decision Protocol | Single Binary Go Orchestrator (&lt; 25MB RAM)
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
    </div>
  );
};
