import React, { useState } from 'react';
import { Claim, ClaimStatus } from '../types';

interface ClaimMatrixProps {
  claims: Claim[];
}

export const ClaimMatrix: React.FC<ClaimMatrixProps> = ({ claims }) => {
  const [filterStatus, setFilterStatus] = useState<string>('ALL');
  const [searchQuery, setSearchQuery] = useState<string>('');

  const filteredClaims = claims.filter((c) => {
    const matchesStatus = filterStatus === 'ALL' || c.audited_status === filterStatus;
    const matchesSearch =
      c.statement.toLowerCase().includes(searchQuery.toLowerCase()) ||
      c.citation.toLowerCase().includes(searchQuery.toLowerCase()) ||
      c.id.toLowerCase().includes(searchQuery.toLowerCase());
    return matchesStatus && matchesSearch;
  });

  const getStatusBadge = (status: ClaimStatus) => {
    switch (status) {
      case 'VERIFIED':
        return (
          <span className="inline-flex items-center px-2 py-0.5 rounded text-[11px] font-mono font-bold bg-emerald-950/90 text-emerald-300 border border-emerald-800">
            VERIFIED
          </span>
        );
      case 'CONCEDED':
        return (
          <span className="inline-flex items-center px-2 py-0.5 rounded text-[11px] font-mono font-bold bg-amber-950/90 text-amber-300 border border-amber-800">
            CONCEDED
          </span>
        );
      case 'UNMITIGATED':
        return (
          <span className="inline-flex items-center px-2 py-0.5 rounded text-[11px] font-mono font-bold bg-rose-950/90 text-rose-300 border border-rose-800">
            UNMITIGATED RISK
          </span>
        );
      case 'HALLUCINATED':
        return (
          <span className="inline-flex items-center px-2 py-0.5 rounded text-[11px] font-mono font-bold bg-red-950/90 text-red-300 border border-red-800">
            HALLUCINATED
          </span>
        );
      case 'UNVERIFIED':
        return (
          <span className="inline-flex items-center px-2 py-0.5 rounded text-[11px] font-mono font-medium bg-slate-800 text-slate-300 border border-slate-700">
            UNVERIFIED
          </span>
        );
      default:
        return (
          <span className="inline-flex items-center px-2 py-0.5 rounded text-[11px] font-mono font-medium bg-slate-900 text-slate-400 border border-slate-800">
            PENDING AUDIT
          </span>
        );
    }
  };

  return (
    <div className="rounded-xl border border-border bg-surface p-5 shadow-lg">
      <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4 border-b border-border/50 pb-4 mb-4">
        <div>
          <h2 className="text-sm font-semibold tracking-wide text-white uppercase font-mono">
            Audited Claim Verification Matrix
          </h2>
          <p className="text-xs text-slate-400">
            Forensic reconciliation of agent assertions against raw dossier citations
          </p>
        </div>

        {/* Filters */}
        <div className="flex flex-wrap items-center gap-2">
          <input
            type="text"
            placeholder="Search claims or citations..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="rounded-lg bg-surface-elevated border border-border px-3 py-1 text-xs text-slate-200 placeholder-slate-500 focus:outline-none focus:ring-1 focus:ring-cyan-500 w-48"
          />

          <select
            value={filterStatus}
            onChange={(e) => setFilterStatus(e.target.value)}
            className="rounded-lg bg-surface-elevated border border-border px-3 py-1 text-xs text-slate-200 focus:outline-none focus:ring-1 focus:ring-cyan-500 font-mono"
          >
            <option value="ALL">ALL STATUSES ({claims.length})</option>
            <option value="VERIFIED">VERIFIED</option>
            <option value="UNMITIGATED">UNMITIGATED</option>
            <option value="CONCEDED">CONCEDED</option>
            <option value="UNVERIFIED">UNVERIFIED</option>
            <option value="PENDING">PENDING</option>
          </select>
        </div>
      </div>

      {/* Table */}
      <div className="overflow-x-auto">
        <table className="w-full text-left text-xs">
          <thead className="bg-surface-elevated text-slate-400 font-mono text-[11px] uppercase tracking-wider border-b border-border">
            <tr>
              <th className="px-3 py-2.5">ID</th>
              <th className="px-3 py-2.5">Role</th>
              <th className="px-3 py-2.5">Statement</th>
              <th className="px-3 py-2.5">Documentary Citation</th>
              <th className="px-3 py-2.5">Audit Status</th>
              <th className="px-3 py-2.5">Forensic Finding</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/60">
            {filteredClaims.length === 0 ? (
              <tr>
                <td colSpan={6} className="px-4 py-8 text-center text-slate-500 italic">
                  No claims match the selected criteria or awaiting debate start.
                </td>
              </tr>
            ) : (
              filteredClaims.map((c) => (
                <tr key={c.id} className="hover:bg-surface-elevated/40 transition-colors">
                  <td className="px-3 py-2.5 font-mono font-medium text-cyan-400 whitespace-nowrap">
                    {c.id}
                  </td>
                  <td className="px-3 py-2.5 whitespace-nowrap">
                    <span
                      className={`inline-block px-1.5 py-0.5 rounded text-[10px] font-mono font-bold ${
                        c.agent_role === 'Advocate'
                          ? 'bg-emerald-950/80 text-emerald-400 border border-emerald-800/60'
                          : 'bg-rose-950/80 text-rose-400 border border-rose-800/60'
                      }`}
                    >
                      {c.agent_role}
                    </span>
                  </td>
                  <td className="px-3 py-2.5 font-medium text-slate-200 max-w-xs md:max-w-md">
                    {c.statement}
                  </td>
                  <td className="px-3 py-2.5 font-mono text-slate-400 text-[11px] max-w-xs truncate">
                    {c.citation}
                  </td>
                  <td className="px-3 py-2.5 whitespace-nowrap">{getStatusBadge(c.audited_status)}</td>
                  <td className="px-3 py-2.5 text-slate-300 text-[11px] max-w-xs">
                    {c.audit_notes || 'Consistent with raw input text'}
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};
