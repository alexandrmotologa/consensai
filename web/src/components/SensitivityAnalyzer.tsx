import React, { useState } from 'react';
import { SynthesisResult } from '../types';

interface SensitivityAnalyzerProps {
  sessionId: string;
  synthesis: SynthesisResult | null;
  onRecalculated: (newSynthesis: SynthesisResult) => void;
}

export const SensitivityAnalyzer: React.FC<SensitivityAnalyzerProps> = ({
  sessionId,
  synthesis,
  onRecalculated,
}) => {
  const [threshApprove, setThreshApprove] = useState<number>(0.7);
  const [threshReject, setThreshReject] = useState<number>(0.35);
  const [isRecalculating, setIsRecalculating] = useState<boolean>(false);

  const handleRecalculate = async (approve: number, reject: number) => {
    if (!sessionId) return;
    setIsRecalculating(true);
    try {
      const res = await fetch(`/api/deliberate/${sessionId}/recalculate`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          threshold_approve: approve,
          threshold_reject: reject,
        }),
      });

      if (res.ok) {
        const data = await res.json();
        onRecalculated(data);
      }
    } catch (err) {
      console.error('Recalculation error:', err);
    } finally {
      setIsRecalculating(false);
    }
  };

  const currentScore = synthesis ? synthesis.normalized_score : 0.5;

  return (
    <div className="rounded-xl border border-border bg-surface p-5 shadow-lg space-y-4">
      <div className="flex flex-col md:flex-row md:items-center md:justify-between border-b border-border/50 pb-3 gap-2">
        <div>
          <h2 className="text-sm font-semibold tracking-wide text-white uppercase font-mono">
            What-If Epistemic Sensitivity Analyzer
          </h2>
          <p className="text-xs text-slate-400">
            Dynamically stress-test risk thresholds without repeating LLM inference calls
          </p>
        </div>
        <div className="text-xs font-mono text-cyan-400">
          Normalized Score: <span className="font-bold text-white text-sm">{currentScore.toFixed(2)}</span>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 text-xs">
        {/* Rejection Threshold Slider */}
        <div className="space-y-2">
          <div className="flex justify-between font-mono">
            <span className="text-slate-400">Rejection Threshold (&theta;_reject):</span>
            <span className="text-rose-400 font-bold">&le; {threshReject.toFixed(2)}</span>
          </div>
          <input
            type="range"
            min="0.10"
            max="0.50"
            step="0.05"
            value={threshReject}
            onChange={(e) => {
              const val = parseFloat(e.target.value);
              setThreshReject(val);
              handleRecalculate(threshApprove, val);
            }}
            className="w-full accent-rose-500 cursor-pointer"
          />
          <div className="flex justify-between text-[10px] text-slate-500 font-mono">
            <span>0.10 (Permissive)</span>
            <span>0.35 (Default)</span>
            <span>0.50 (Strict)</span>
          </div>
        </div>

        {/* Approval Threshold Slider */}
        <div className="space-y-2">
          <div className="flex justify-between font-mono">
            <span className="text-slate-400">Approval Threshold (&theta;_approve):</span>
            <span className="text-emerald-400 font-bold">&ge; {threshApprove.toFixed(2)}</span>
          </div>
          <input
            type="range"
            min="0.55"
            max="0.90"
            step="0.05"
            value={threshApprove}
            onChange={(e) => {
              const val = parseFloat(e.target.value);
              setThreshApprove(val);
              handleRecalculate(val, threshReject);
            }}
            className="w-full accent-emerald-500 cursor-pointer"
          />
          <div className="flex justify-between text-[10px] text-slate-500 font-mono">
            <span>0.55 (Aggressive)</span>
            <span>0.70 (Default)</span>
            <span>0.90 (Conservative)</span>
          </div>
        </div>
      </div>

      {/* Sensitivity Delta Summary */}
      <div className="rounded-lg bg-surface-elevated/70 border border-border/80 p-3 text-xs font-mono space-y-1">
        <div className="text-slate-300">
          <span className="text-cyan-400 font-semibold">Risk Frontier Status: </span>
          {currentScore >= threshApprove ? (
            <span className="text-emerald-400 font-bold">
              Score ({currentScore.toFixed(2)}) clears approval barrier ({threshApprove.toFixed(2)}).
            </span>
          ) : currentScore <= threshReject ? (
            <span className="text-rose-400 font-bold">
              Score ({currentScore.toFixed(2)}) falls below rejection floor ({threshReject.toFixed(2)}).
            </span>
          ) : (
            <span className="text-amber-400 font-bold">
              Score ({currentScore.toFixed(2)}) requires human escalation [{threshReject.toFixed(2)} - {threshApprove.toFixed(2)}].
            </span>
          )}
        </div>
        <div className="text-[11px] text-slate-500">
          {currentScore < threshApprove && (
            <span>
              &bull; Needs a +{(threshApprove - currentScore).toFixed(2)} score delta or lowering &theta;_approve to {currentScore.toFixed(2)} for approval.
            </span>
          )}
        </div>
      </div>
    </div>
  );
};
