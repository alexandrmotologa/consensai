import React from 'react';
import { SynthesisResult } from '../types';

interface ConvergenceGaugeProps {
  synthesis: SynthesisResult | null;
  isRunning: boolean;
  currentPhase: string;
}

export const ConvergenceGauge: React.FC<ConvergenceGaugeProps> = ({
  synthesis,
  isRunning,
  currentPhase,
}) => {
  const score = synthesis ? synthesis.normalized_score : 0.5;
  const rawScore = synthesis ? synthesis.epistemic_convergence_score : 0.0;
  const ci = synthesis ? synthesis.confidence_interval : [0.0, 1.0];
  const decision = synthesis ? synthesis.decision : null;

  // Gauge calculation (180 degree semi-circle: from -90 to +90 degrees)
  // angle = -90 + (score * 180)
  const angle = -90 + score * 180;
  const ciStartAngle = -90 + ci[0] * 180;
  const ciEndAngle = -90 + ci[1] * 180;

  const getDecisionBadge = () => {
    if (!decision) {
      if (isRunning) {
        return (
          <div className="flex items-center space-x-2 rounded-full bg-cyan-950/80 border border-cyan-800 px-4 py-1 text-xs font-semibold text-cyan-300">
            <span className="h-2 w-2 rounded-full bg-cyan-400 animate-ping"></span>
            <span>DELIBERATING ({currentPhase})</span>
          </div>
        );
      }
      return (
        <div className="rounded-full bg-slate-800/80 border border-slate-700 px-4 py-1 text-xs font-semibold text-slate-400">
          AWAITING DELIBERATION
        </div>
      );
    }

    switch (decision) {
      case 'APPROVE':
        return (
          <div className="flex items-center space-x-2 rounded-full bg-emerald-950/90 border border-emerald-500/80 px-4 py-1 text-xs font-bold text-emerald-300 shadow-lg shadow-emerald-500/20">
            <span className="h-2 w-2 rounded-full bg-emerald-400"></span>
            <span>APPROVED</span>
          </div>
        );
      case 'REJECT':
        return (
          <div className="flex items-center space-x-2 rounded-full bg-rose-950/90 border border-rose-500/80 px-4 py-1 text-xs font-bold text-rose-300 shadow-lg shadow-rose-500/20">
            <span className="h-2 w-2 rounded-full bg-rose-400"></span>
            <span>REJECTED</span>
          </div>
        );
      default:
        return (
          <div className="flex items-center space-x-2 rounded-full bg-amber-950/90 border border-amber-500/80 px-4 py-1 text-xs font-bold text-amber-300 shadow-lg shadow-amber-500/20">
            <span className="h-2 w-2 rounded-full bg-amber-400 animate-pulse"></span>
            <span>ESCALATE TO HUMAN OVERSIGHT</span>
          </div>
        );
    }
  };

  return (
    <div className="rounded-xl border border-border bg-surface p-5 shadow-lg flex flex-col items-center">
      <div className="w-full flex items-center justify-between border-b border-border/50 pb-3 mb-4">
        <div>
          <h2 className="text-sm font-semibold tracking-wide text-white uppercase font-mono">
            Epistemic Convergence Metric (C_epistemic)
          </h2>
          <p className="text-xs text-slate-400">
            Calibrated deterministic agreement across dialectical agent turns
          </p>
        </div>
        {getDecisionBadge()}
      </div>

      {/* SVG Semi-Circle Gauge */}
      <div className="relative w-64 h-36 flex items-center justify-center">
        <svg viewBox="0 0 200 110" className="w-full h-full overflow-visible">
          {/* Background Arc */}
          <path
            d="M 20 100 A 80 80 0 0 1 180 100"
            fill="none"
            stroke="#1f2937"
            strokeWidth="14"
            strokeLinecap="round"
          />

          {/* Reject Zone (0.00 to 0.35) -> 0 to 63 degrees */}
          <path
            d="M 20 100 A 80 80 0 0 1 54 44"
            fill="none"
            stroke="#f43f5e"
            strokeWidth="14"
            strokeOpacity="0.4"
          />

          {/* Escalate Zone (0.35 to 0.70) -> 63 to 126 degrees */}
          <path
            d="M 54 44 A 80 80 0 0 1 146 44"
            fill="none"
            stroke="#f59e0b"
            strokeWidth="14"
            strokeOpacity="0.4"
          />

          {/* Approve Zone (0.70 to 1.00) -> 126 to 180 degrees */}
          <path
            d="M 146 44 A 80 80 0 0 1 180 100"
            fill="none"
            stroke="#10b981"
            strokeWidth="14"
            strokeOpacity="0.4"
          />

          {/* Active Score Arc */}
          <path
            d={`M 20 100 A 80 80 0 0 1 ${100 + 80 * Math.sin((angle * Math.PI) / 180)} ${
              100 - 80 * Math.cos((angle * Math.PI) / 180)
            }`}
            fill="none"
            stroke={score >= 0.7 ? '#10b981' : score <= 0.35 ? '#f43f5e' : '#f59e0b'}
            strokeWidth="14"
            strokeLinecap="round"
          />

          {/* Threshold marker 0.35 (Reject Boundary) */}
          <line
            x1="52"
            y1="42"
            x2="57"
            y2="47"
            stroke="#ffffff"
            strokeWidth="2"
            strokeDasharray="2 2"
          />

          {/* Threshold marker 0.70 (Approve Boundary) */}
          <line
            x1="148"
            y1="42"
            x2="143"
            y2="47"
            stroke="#ffffff"
            strokeWidth="2"
            strokeDasharray="2 2"
          />

          {/* Center Pivot */}
          <circle cx="100" cy="100" r="7" fill="#38bdf8" />
          <circle cx="100" cy="100" r="3" fill="#ffffff" />

          {/* Needle */}
          <line
            x1="100"
            y1="100"
            x2={100 + 68 * Math.sin((angle * Math.PI) / 180)}
            y2={100 - 68 * Math.cos((angle * Math.PI) / 180)}
            stroke="#ffffff"
            strokeWidth="2.5"
            strokeLinecap="round"
          />
        </svg>

        {/* Big Numeric Score Display */}
        <div className="absolute bottom-0 flex flex-col items-center">
          <span className="text-3xl font-extrabold font-mono text-white tracking-tight">
            {score.toFixed(2)}
          </span>
          <span className="text-[10px] font-mono text-slate-400">
            Raw: {rawScore > 0 ? `+${rawScore.toFixed(3)}` : rawScore.toFixed(3)}
          </span>
        </div>
      </div>

      {/* Threshold and Confidence Interval Indicators */}
      <div className="mt-4 grid grid-cols-3 gap-3 w-full max-w-lg text-center font-mono text-xs">
        <div className="rounded-lg bg-rose-950/30 border border-rose-900/40 p-2">
          <div className="text-slate-400 text-[10px]">REJECT ZONE</div>
          <div className="text-rose-400 font-semibold">&le; 0.35</div>
        </div>
        <div className="rounded-lg bg-amber-950/30 border border-amber-900/40 p-2">
          <div className="text-slate-400 text-[10px]">ESCALATE (HUMAN)</div>
          <div className="text-amber-400 font-semibold">0.35 - 0.70</div>
        </div>
        <div className="rounded-lg bg-emerald-950/30 border border-emerald-900/40 p-2">
          <div className="text-slate-400 text-[10px]">APPROVE ZONE</div>
          <div className="text-emerald-400 font-semibold">&ge; 0.70</div>
        </div>
      </div>

      {/* Detailed metrics footer */}
      {synthesis && (
        <div className="mt-4 w-full flex flex-wrap items-center justify-between border-t border-border/50 pt-3 text-xs text-slate-300">
          <div>
            <span className="text-slate-400">95% Confidence Interval: </span>
            <span className="font-mono text-cyan-300">
              [{ci[0].toFixed(2)}, {ci[1].toFixed(2)}]
            </span>
          </div>
          <div>
            <span className="text-slate-400">Affirmative Grounding: </span>
            <span className="font-mono text-emerald-400">
              {synthesis.verified_affirmative_claims} / {synthesis.total_claims_audited} verified
            </span>
          </div>
          <div>
            <span className="text-slate-400">Unmitigated Risks: </span>
            <span className={`font-mono ${synthesis.unmitigated_adversarial_risks > 0 ? 'text-rose-400 font-bold' : 'text-slate-400'}`}>
              {synthesis.unmitigated_adversarial_risks}
            </span>
          </div>
        </div>
      )}
    </div>
  );
};
