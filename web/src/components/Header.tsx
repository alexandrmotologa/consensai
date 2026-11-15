import React from 'react';
import { DossierSummary, SystemStatus } from '../types';

interface HeaderProps {
  dossiers: DossierSummary[];
  selectedDossier: string;
  onSelectDossier: (id: string) => void;
  onRunDeliberation: () => void;
  isRunning: boolean;
  systemStatus: SystemStatus | null;
  onOpenDossierModal: () => void;
  onOpenMemoModal: () => void;
  hasCompletedMemo: boolean;
}

export const Header: React.FC<HeaderProps> = ({
  dossiers,
  selectedDossier,
  onSelectDossier,
  onRunDeliberation,
  isRunning,
  systemStatus,
  onOpenDossierModal,
  onOpenMemoModal,
  hasCompletedMemo,
}) => {
  return (
    <header className="border-b border-border bg-surface px-6 py-4">
      <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        {/* Brand and Tagline */}
        <div className="flex items-center space-x-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-gradient-to-br from-cyan-500 to-blue-600 text-white font-bold text-lg shadow-md shadow-cyan-500/20">
            C
          </div>
          <div>
            <div className="flex items-center space-x-2">
              <h1 className="text-xl font-bold tracking-tight text-white font-mono">consensai</h1>
              <span className="rounded bg-cyan-950/80 px-2 py-0.5 text-[11px] font-semibold text-cyan-400 border border-cyan-800/60">
                v1.0.0
              </span>
            </div>
            <p className="text-xs text-slate-400">
              Adversarial Multi-Agent Epistemic Deliberation & SHA-256 Ledger
            </p>
          </div>
        </div>

        {/* Runtime Performance Stats */}
        {systemStatus && (
          <div className="hidden lg:flex items-center space-x-4 px-4 py-1.5 rounded-lg bg-surface-elevated/70 border border-border/60 text-xs font-mono">
            <div className="flex items-center space-x-1.5">
              <span className="h-2 w-2 rounded-full bg-emerald-400 animate-pulse"></span>
              <span className="text-slate-400">STATUS:</span>
              <span className="text-emerald-400 font-medium">ONLINE</span>
            </div>
            <span className="text-slate-600">|</span>
            <div>
              <span className="text-slate-400">RAM: </span>
              <span className={`font-medium ${systemStatus.allocated_memory_mb > 25 ? 'text-amber-400' : 'text-cyan-400'}`}>
                {systemStatus.allocated_memory_mb.toFixed(2)} MB
              </span>
              <span className="text-slate-500"> (&lt; 25MB target)</span>
            </div>
            <span className="text-slate-600">|</span>
            <div>
              <span className="text-slate-400">UPTIME: </span>
              <span className="text-slate-200">{systemStatus.uptime_seconds}s</span>
            </div>
          </div>
        )}

        {/* Controls */}
        <div className="flex flex-wrap items-center gap-3">
          {/* Dossier Picker */}
          <div className="flex items-center space-x-2">
            <label className="text-xs font-medium text-slate-400">Dossier:</label>
            <select
              value={selectedDossier}
              onChange={(e) => onSelectDossier(e.target.value)}
              disabled={isRunning}
              className="rounded-lg bg-surface-elevated border border-border px-3 py-1.5 text-xs text-slate-200 focus:outline-none focus:ring-1 focus:ring-cyan-500 disabled:opacity-50"
            >
              {dossiers.map((d) => (
                <option key={d.id} value={d.id}>
                  {d.title}
                </option>
              ))}
            </select>
          </div>

          <button
            onClick={onOpenDossierModal}
            className="rounded-lg bg-surface-elevated hover:bg-slate-700/80 border border-border px-3 py-1.5 text-xs font-medium text-slate-300 transition-colors"
          >
            Inspect Dossier
          </button>

          {hasCompletedMemo && (
            <button
              onClick={onOpenMemoModal}
              className="rounded-lg bg-cyan-950/60 hover:bg-cyan-900/60 border border-cyan-700/60 px-3 py-1.5 text-xs font-medium text-cyan-300 transition-colors flex items-center space-x-1"
            >
              <span>Decision Memo</span>
            </button>
          )}

          {/* Run Action */}
          <button
            onClick={onRunDeliberation}
            disabled={isRunning}
            className="flex items-center space-x-2 rounded-lg bg-gradient-to-r from-cyan-600 to-blue-600 hover:from-cyan-500 hover:to-blue-500 px-4 py-1.5 text-xs font-semibold text-white shadow-md shadow-cyan-600/30 transition-all disabled:opacity-50"
          >
            {isRunning ? (
              <>
                <svg className="animate-spin h-3.5 w-3.5 text-white" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"></path>
                </svg>
                <span>Deliberating...</span>
              </>
            ) : (
              <>
                <svg className="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"></path>
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
                <span>Run Deliberation</span>
              </>
            )}
          </button>
        </div>
      </div>
    </header>
  );
};
