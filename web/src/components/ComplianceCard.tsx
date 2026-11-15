import React from 'react';
import { ComplianceReport } from '../types';

interface ComplianceCardProps {
  report?: ComplianceReport;
}

export const ComplianceCard: React.FC<ComplianceCardProps> = ({ report }) => {
  if (!report || report.checks.length === 0) {
    return (
      <div className="rounded-xl border border-border bg-surface p-5 shadow-lg">
        <h2 className="text-sm font-semibold tracking-wide text-white uppercase font-mono border-b border-border/50 pb-2 mb-3">
          Statutory Regulatory Conformance
        </h2>
        <div className="py-4 text-xs text-slate-500 italic text-center">
          Regulatory rules will be evaluated upon deliberation synthesis.
        </div>
      </div>
    );
  }

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'COMPLIANT':
        return (
          <span className="px-2 py-0.5 rounded text-[10px] font-mono font-bold bg-emerald-950/80 text-emerald-300 border border-emerald-800">
            COMPLIANT
          </span>
        );
      case 'FLAGGED':
        return (
          <span className="px-2 py-0.5 rounded text-[10px] font-mono font-bold bg-rose-950/80 text-rose-300 border border-rose-800">
            FLAGGED
          </span>
        );
      case 'REQUIRES_HUMAN_ACTION':
        return (
          <span className="px-2 py-0.5 rounded text-[10px] font-mono font-bold bg-amber-950/80 text-amber-300 border border-amber-800">
            REQUIRES ACTION
          </span>
        );
      default:
        return (
          <span className="px-2 py-0.5 rounded text-[10px] font-mono font-medium bg-slate-800 text-slate-300 border border-slate-700">
            N/A
          </span>
        );
    }
  };

  return (
    <div className="rounded-xl border border-border bg-surface p-5 shadow-lg">
      <div className="flex items-center justify-between border-b border-border/50 pb-3 mb-4">
        <div>
          <h2 className="text-sm font-semibold tracking-wide text-white uppercase font-mono">
            Statutory Regulatory Conformance Scorecard
          </h2>
          <p className="text-xs text-slate-400">
            Automated verification against EU AI Act Art. 14, Basel III, FCRA, and HIPAA mandates
          </p>
        </div>
        <div className="flex items-center space-x-2">
          <span className="text-xs text-slate-400 font-mono">OVERALL:</span>
          {getStatusBadge(report.overall_status)}
        </div>
      </div>

      <div className="space-y-3">
        {report.checks.map((chk, idx) => (
          <div
            key={idx}
            className="rounded-lg bg-surface-elevated/70 border border-border/80 p-3 text-xs space-y-1.5"
          >
            <div className="flex items-center justify-between">
              <span className="font-mono font-bold text-cyan-400">{chk.standard}</span>
              {getStatusBadge(chk.status)}
            </div>
            <p className="text-slate-300 font-medium text-[11px]">{chk.description}</p>
            <div className="bg-background/60 p-2 rounded border border-border/50 text-[11px] font-mono text-slate-400">
              &gt; {chk.notes}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
