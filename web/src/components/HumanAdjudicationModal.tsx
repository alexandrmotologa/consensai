import React, { useState } from 'react';

interface HumanAdjudicationModalProps {
  isOpen: boolean;
  onClose: () => void;
  sessionId: string;
  onAdjudicated: (data: any) => void;
}

export const HumanAdjudicationModal: React.FC<HumanAdjudicationModalProps> = ({
  isOpen,
  onClose,
  sessionId,
  onAdjudicated,
}) => {
  const [officerId, setOfficerId] = useState<string>('CRO-8842');
  const [officerName, setOfficerName] = useState<string>('Elena Vance');
  const [officerRole, setOfficerRole] = useState<string>('Chief Risk Officer');
  const [decision, setDecision] = useState<string>('CONDITIONAL_APPROVE');
  const [justification, setJustification] = useState<string>(
    'Approved subject to establishing a mandatory $50,000 cash collateral reserve and quarterly covenant testing.'
  );
  const [conditionInput, setConditionInput] = useState<string>('');
  const [conditions, setConditions] = useState<string[]>([
    'Quarterly DSCR covenant testing floor at 1.20x',
    'Segregated $50,000 cash escrow reserve dedicated to supplier litigation',
  ]);
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  if (!isOpen) return null;

  const handleAddCondition = () => {
    if (conditionInput.trim()) {
      setConditions([...conditions, conditionInput.trim()]);
      setConditionInput('');
    }
  };

  const handleRemoveCondition = (index: number) => {
    setConditions(conditions.filter((_, i) => i !== index));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!sessionId) {
      setError('No active deliberation session to adjudicate.');
      return;
    }

    setIsSubmitting(true);
    setError(null);

    try {
      const res = await fetch(`/api/deliberate/${sessionId}/adjudicate`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          officer_id: officerId,
          officer_name: officerName,
          officer_role: officerRole,
          decision,
          justification,
          conditions,
        }),
      });

      if (!res.ok) {
        const errText = await res.text();
        throw new Error(errText || 'Adjudication failed');
      }

      const data = await res.json();
      onAdjudicated(data);
      onClose();
    } catch (err: any) {
      setError(err.message || 'Submission error');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm p-4">
      <div className="relative w-full max-w-2xl rounded-xl border border-border bg-surface shadow-2xl flex flex-col max-h-[90vh]">
        {/* Header */}
        <div className="flex items-center justify-between border-b border-border px-6 py-4">
          <div className="flex items-center space-x-2">
            <span className="flex h-6 w-6 items-center justify-center rounded-full bg-cyan-950 border border-cyan-500 text-cyan-300 font-bold text-xs">
              H
            </span>
            <div>
              <h3 className="text-base font-bold text-white font-mono">
                Human-in-the-Loop (HITL) Adjudication
              </h3>
              <p className="text-xs text-slate-400">EU AI Act Article 14 Governance Override</p>
            </div>
          </div>
          <button onClick={onClose} className="rounded p-1 text-slate-400 hover:text-white">
            <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>

        {/* Form Body */}
        <form onSubmit={handleSubmit} className="flex-1 overflow-y-auto p-6 space-y-4 text-xs">
          {error && (
            <div className="p-3 rounded-lg bg-rose-950/80 border border-rose-600 text-rose-200 font-mono">
              {error}
            </div>
          )}

          <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
            <div>
              <label className="block text-slate-400 font-medium mb-1">Officer ID</label>
              <input
                type="text"
                value={officerId}
                onChange={(e) => setOfficerId(e.target.value)}
                required
                className="w-full rounded bg-surface-elevated border border-border px-3 py-1.5 text-slate-200 font-mono focus:outline-none focus:ring-1 focus:ring-cyan-500"
              />
            </div>
            <div>
              <label className="block text-slate-400 font-medium mb-1">Officer Name</label>
              <input
                type="text"
                value={officerName}
                onChange={(e) => setOfficerName(e.target.value)}
                required
                className="w-full rounded bg-surface-elevated border border-border px-3 py-1.5 text-slate-200 focus:outline-none focus:ring-1 focus:ring-cyan-500"
              />
            </div>
            <div>
              <label className="block text-slate-400 font-medium mb-1">Title / Role</label>
              <input
                type="text"
                value={officerRole}
                onChange={(e) => setOfficerRole(e.target.value)}
                required
                className="w-full rounded bg-surface-elevated border border-border px-3 py-1.5 text-slate-200 focus:outline-none focus:ring-1 focus:ring-cyan-500"
              />
            </div>
          </div>

          <div>
            <label className="block text-slate-400 font-medium mb-1">Determination Outcome</label>
            <select
              value={decision}
              onChange={(e) => setDecision(e.target.value)}
              className="w-full rounded bg-surface-elevated border border-border px-3 py-2 text-slate-200 font-mono font-bold focus:outline-none focus:ring-1 focus:ring-cyan-500"
            >
              <option value="CONDITIONAL_APPROVE">CONDITIONAL APPROVAL (Require Mitigations)</option>
              <option value="OVERRIDE_APPROVE">OVERRIDE APPROVE (Executive Exemption)</option>
              <option value="OVERRIDE_REJECT">OVERRIDE REJECT (Block Transaction)</option>
              <option value="AFFIRM_REJECT">AFFIRM REJECTION (Final Dismissal)</option>
            </select>
          </div>

          <div>
            <label className="block text-slate-400 font-medium mb-1">
              Officer Justification & Legal Basis (Required)
            </label>
            <textarea
              rows={3}
              value={justification}
              onChange={(e) => setJustification(e.target.value)}
              required
              placeholder="State the formal underwriting or risk rationale for this determination..."
              className="w-full rounded bg-surface-elevated border border-border p-3 text-slate-200 focus:outline-none focus:ring-1 focus:ring-cyan-500"
            />
          </div>

          <div>
            <label className="block text-slate-400 font-medium mb-1">
              Mandatory Operational Conditions
            </label>
            <div className="flex space-x-2 mb-2">
              <input
                type="text"
                value={conditionInput}
                onChange={(e) => setConditionInput(e.target.value)}
                placeholder="Add compliance condition or covenant..."
                className="flex-1 rounded bg-surface-elevated border border-border px-3 py-1.5 text-slate-200 focus:outline-none focus:ring-1 focus:ring-cyan-500"
              />
              <button
                type="button"
                onClick={handleAddCondition}
                className="rounded bg-surface-elevated hover:bg-slate-700 px-3 py-1.5 text-slate-200 font-semibold"
              >
                Add
              </button>
            </div>
            <div className="space-y-1.5">
              {conditions.map((c, i) => (
                <div
                  key={i}
                  className="flex items-center justify-between rounded bg-surface-elevated/70 p-2 border border-border/60"
                >
                  <span className="text-slate-300 font-mono text-[11px]">&gt; {c}</span>
                  <button
                    type="button"
                    onClick={() => handleRemoveCondition(i)}
                    className="text-slate-500 hover:text-rose-400 text-xs px-1"
                  >
                    &times;
                  </button>
                </div>
              ))}
            </div>
          </div>

          <div className="border-t border-border pt-4 flex items-center justify-between">
            <span className="text-[11px] text-slate-500 font-mono">
              Committed as Block #5 with SHA-256 digital signature
            </span>
            <div className="flex space-x-2">
              <button
                type="button"
                onClick={onClose}
                className="rounded bg-surface-elevated hover:bg-slate-700 px-4 py-2 font-medium text-slate-300"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={isSubmitting}
                className="rounded bg-gradient-to-r from-emerald-600 to-teal-600 hover:from-emerald-500 hover:to-teal-500 px-5 py-2 font-bold text-white shadow-md disabled:opacity-50"
              >
                {isSubmitting ? 'Signing Block...' : 'Sign & Commit to Ledger'}
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
};
