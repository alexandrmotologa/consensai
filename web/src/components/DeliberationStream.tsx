import React from 'react';
import { AgentTurn, AgentRole } from '../types';

interface DeliberationStreamProps {
  turns: AgentTurn[];
  activeRole?: AgentRole;
  currentRound: number;
}

export const DeliberationStream: React.FC<DeliberationStreamProps> = ({
  turns,
  activeRole,
  currentRound,
}) => {
  const advocateTurns = turns.filter((t) => t.role === 'Advocate');
  const skepticTurns = turns.filter((t) => t.role === 'Skeptic');
  const auditorTurns = turns.filter((t) => t.role === 'Risk Auditor');

  return (
    <div className="rounded-xl border border-border bg-surface p-5 shadow-lg">
      <div className="flex items-center justify-between border-b border-border/50 pb-3 mb-4">
        <div>
          <h2 className="text-sm font-semibold tracking-wide text-white uppercase font-mono">
            Dialectical Debate Stream (3-Round Confrontation)
          </h2>
          <p className="text-xs text-slate-400">
            Adversarial cross-examination, concession logging, and factual grounding
          </p>
        </div>
        <div className="flex items-center space-x-2 text-xs font-mono">
          <span className="text-slate-400">CURRENT ROUND:</span>
          <span className="rounded bg-surface-elevated px-2 py-0.5 font-bold text-cyan-400 border border-border">
            {currentRound > 0 ? `ROUND ${currentRound}` : 'READY'}
          </span>
        </div>
      </div>

      {/* 3 Columns */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {/* Advocate Column */}
        <div className="flex flex-col rounded-lg bg-surface-elevated/40 border border-emerald-900/30 p-3">
          <div className="flex items-center justify-between border-b border-emerald-900/50 pb-2 mb-3">
            <div className="flex items-center space-x-2">
              <span className="h-2.5 w-2.5 rounded-full bg-emerald-400"></span>
              <h3 className="text-xs font-bold uppercase tracking-wider text-emerald-400 font-mono">
                The Advocate
              </h3>
            </div>
            <span className="text-[10px] text-slate-400">Affirmative Case</span>
          </div>

          <div className="space-y-3 flex-1 overflow-y-auto max-h-[550px] pr-1">
            {advocateTurns.length === 0 && (
              <div className="flex h-32 items-center justify-center text-xs text-slate-500 italic">
                {activeRole === 'Advocate' ? 'Formulating affirmative thesis...' : 'Awaiting turn...'}
              </div>
            )}
            {advocateTurns.map((turn, idx) => (
              <div
                key={idx}
                className="rounded-lg bg-surface border border-emerald-900/40 p-3 shadow-sm hover:border-emerald-700/60 transition-colors"
              >
                <div className="flex items-center justify-between mb-2">
                  <span className="rounded bg-emerald-950/80 px-1.5 py-0.5 text-[10px] font-mono font-semibold text-emerald-300 border border-emerald-800/60">
                    ROUND {turn.round}
                  </span>
                  <span className="text-[10px] text-slate-500 font-mono">
                    {new Date(turn.timestamp).toLocaleTimeString()}
                  </span>
                </div>

                <p className="text-xs text-slate-200 mb-2 leading-relaxed font-medium">
                  {turn.summary}
                </p>

                {turn.claims && turn.claims.length > 0 && (
                  <div className="mt-2 space-y-1.5">
                    <span className="text-[10px] font-mono text-emerald-500 font-semibold block">
                      DEFENDED CLAIMS:
                    </span>
                    {turn.claims.map((c) => (
                      <div
                        key={c.id}
                        className="rounded bg-surface-elevated/80 p-2 text-[11px] border border-border/50"
                      >
                        <div className="text-slate-200 font-medium mb-1">{c.statement}</div>
                        <div className="font-mono text-[10px] text-slate-400 truncate">
                          &gt; {c.citation}
                        </div>
                      </div>
                    ))}
                  </div>
                )}

                {turn.rebuttals && turn.rebuttals.length > 0 && (
                  <div className="mt-2 space-y-1.5">
                    <span className="text-[10px] font-mono text-cyan-400 font-semibold block">
                      REBUTTALS:
                    </span>
                    {turn.rebuttals.map((r, i) => (
                      <div
                        key={i}
                        className="rounded bg-surface-elevated/80 p-2 text-[11px] border border-border/50"
                      >
                        <div className="text-slate-300">{r.counter_point}</div>
                      </div>
                    ))}
                  </div>
                )}

                {turn.concessions && turn.concessions.length > 0 && (
                  <div className="mt-2 space-y-1">
                    <span className="text-[10px] font-mono text-amber-400 font-semibold block">
                      FORMAL CONCESSIONS:
                    </span>
                    {turn.concessions.map((con, i) => (
                      <div key={i} className="text-[11px] text-amber-200/90 italic pl-2 border-l border-amber-500/40">
                        {con}
                      </div>
                    ))}
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>

        {/* Skeptic Column */}
        <div className="flex flex-col rounded-lg bg-surface-elevated/40 border border-rose-900/30 p-3">
          <div className="flex items-center justify-between border-b border-rose-900/50 pb-2 mb-3">
            <div className="flex items-center space-x-2">
              <span className="h-2.5 w-2.5 rounded-full bg-rose-400"></span>
              <h3 className="text-xs font-bold uppercase tracking-wider text-rose-400 font-mono">
                The Skeptic
              </h3>
            </div>
            <span className="text-[10px] text-slate-400">Adversarial Risk</span>
          </div>

          <div className="space-y-3 flex-1 overflow-y-auto max-h-[550px] pr-1">
            {skepticTurns.length === 0 && (
              <div className="flex h-32 items-center justify-center text-xs text-slate-500 italic">
                {activeRole === 'Skeptic' ? 'Cross-examining failure modes...' : 'Awaiting turn...'}
              </div>
            )}
            {skepticTurns.map((turn, idx) => (
              <div
                key={idx}
                className="rounded-lg bg-surface border border-rose-900/40 p-3 shadow-sm hover:border-rose-700/60 transition-colors"
              >
                <div className="flex items-center justify-between mb-2">
                  <span className="rounded bg-rose-950/80 px-1.5 py-0.5 text-[10px] font-mono font-semibold text-rose-300 border border-rose-800/60">
                    ROUND {turn.round}
                  </span>
                  <span className="text-[10px] text-slate-500 font-mono">
                    {new Date(turn.timestamp).toLocaleTimeString()}
                  </span>
                </div>

                <p className="text-xs text-slate-200 mb-2 leading-relaxed font-medium">
                  {turn.summary}
                </p>

                {turn.claims && turn.claims.length > 0 && (
                  <div className="mt-2 space-y-1.5">
                    <span className="text-[10px] font-mono text-rose-400 font-semibold block">
                      IDENTIFIED RISKS & COVENANTS:
                    </span>
                    {turn.claims.map((c) => (
                      <div
                        key={c.id}
                        className="rounded bg-surface-elevated/80 p-2 text-[11px] border border-border/50"
                      >
                        <div className="text-slate-200 font-medium mb-1">{c.statement}</div>
                        <div className="font-mono text-[10px] text-slate-400 truncate">
                          &gt; {c.citation}
                        </div>
                      </div>
                    ))}
                  </div>
                )}

                {turn.rebuttals && turn.rebuttals.length > 0 && (
                  <div className="mt-2 space-y-1.5">
                    <span className="text-[10px] font-mono text-cyan-400 font-semibold block">
                      CROSS-EXAMINATION:
                    </span>
                    {turn.rebuttals.map((r, i) => (
                      <div
                        key={i}
                        className="rounded bg-surface-elevated/80 p-2 text-[11px] border border-border/50"
                      >
                        <div className="text-slate-300">{r.counter_point}</div>
                      </div>
                    ))}
                  </div>
                )}

                {turn.concessions && turn.concessions.length > 0 && (
                  <div className="mt-2 space-y-1">
                    <span className="text-[10px] font-mono text-amber-400 font-semibold block">
                      FORMAL CONCESSIONS:
                    </span>
                    {turn.concessions.map((con, i) => (
                      <div key={i} className="text-[11px] text-amber-200/90 italic pl-2 border-l border-amber-500/40">
                        {con}
                      </div>
                    ))}
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>

        {/* Risk Auditor Column */}
        <div className="flex flex-col rounded-lg bg-surface-elevated/40 border border-amber-900/30 p-3">
          <div className="flex items-center justify-between border-b border-amber-900/50 pb-2 mb-3">
            <div className="flex items-center space-x-2">
              <span className="h-2.5 w-2.5 rounded-full bg-amber-400"></span>
              <h3 className="text-xs font-bold uppercase tracking-wider text-amber-400 font-mono">
                The Risk Auditor
              </h3>
            </div>
            <span className="text-[10px] text-slate-400">Factual Grounding</span>
          </div>

          <div className="space-y-3 flex-1 overflow-y-auto max-h-[550px] pr-1">
            {auditorTurns.length === 0 && (
              <div className="flex h-32 items-center justify-center text-xs text-slate-500 italic">
                {activeRole === 'Risk Auditor' ? 'Auditing citations against dossier...' : 'Awaiting Round 3...'}
              </div>
            )}
            {auditorTurns.map((turn, idx) => (
              <div
                key={idx}
                className="rounded-lg bg-surface border border-amber-900/40 p-3 shadow-sm hover:border-amber-700/60 transition-colors"
              >
                <div className="flex items-center justify-between mb-2">
                  <span className="rounded bg-amber-950/80 px-1.5 py-0.5 text-[10px] font-mono font-semibold text-amber-300 border border-amber-800/60">
                    ROUND {turn.round} (AUDIT)
                  </span>
                  <span className="text-[10px] text-slate-500 font-mono">
                    {new Date(turn.timestamp).toLocaleTimeString()}
                  </span>
                </div>

                <p className="text-xs text-slate-200 mb-2 leading-relaxed font-medium">
                  {turn.summary}
                </p>

                <div className="text-[11px] text-slate-300 leading-relaxed bg-surface-elevated/70 p-2.5 rounded border border-border/60">
                  {turn.raw_content}
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};
