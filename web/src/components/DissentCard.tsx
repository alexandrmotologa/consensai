import React from 'react';

interface DissentCardProps {
  dissents: string[];
}

export const DissentCard: React.FC<DissentCardProps> = ({ dissents }) => {
  return (
    <div className="rounded-xl border border-rose-900/50 bg-gradient-to-br from-rose-950/20 to-surface p-5 shadow-lg">
      <div className="flex items-center space-x-2 border-b border-rose-900/40 pb-3 mb-3">
        <span className="flex h-6 w-6 items-center justify-center rounded-full bg-rose-950 border border-rose-600 text-rose-400 font-bold text-xs">
          !
        </span>
        <div>
          <h2 className="text-sm font-semibold tracking-wide text-rose-300 uppercase font-mono">
            Preserved Minority Dissents & Unmitigated Risks
          </h2>
          <p className="text-xs text-slate-400">
            Anti-Groupthink Guarantee: Adversarial objections preserved verbatim for executive audit
          </p>
        </div>
      </div>

      {dissents.length === 0 ? (
        <div className="py-4 text-xs text-slate-400 italic">
          Zero unmitigated adversarial risks remain. All counter-arguments were either formally conceded or refuted with documentary evidence.
        </div>
      ) : (
        <div className="space-y-2.5">
          {dissents.map((dissent, idx) => (
            <div
              key={idx}
              className="rounded-lg bg-surface border border-rose-900/40 p-3 text-xs text-rose-200/90 leading-relaxed shadow-sm"
            >
              <div className="flex items-start space-x-2">
                <span className="rounded bg-rose-950 px-1.5 py-0.5 text-[10px] font-mono font-bold text-rose-400 border border-rose-800">
                  #{idx + 1}
                </span>
                <span className="flex-1 font-mono text-[11px]">{dissent}</span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
