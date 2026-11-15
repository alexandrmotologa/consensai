import React, { useState } from 'react';

interface DossierBuilderProps {
  onRunCustom: (customDossier: any) => void;
  isRunning: boolean;
}

const TEMPLATES: Record<string, any> = {
  jumbo_mortgage_950k: {
    id: "jumbo_mortgage_950k",
    title: "Residential Jumbo Mortgage Application ($950,000)",
    domain: "credit_underwriting",
    risk_tier: "tier_2_medium",
    borrower: {
      name: "Arthur Pendelton",
      annual_w2_income_usd: 340000,
      liquid_reserves_usd: 210000,
      credit_score: 765,
      debt_to_income_ratio_pct: 38.4
    },
    property: {
      address: "742 Evergreen Terrace, Austin, TX",
      purchase_price_usd: 1250000,
      appraised_value_usd: 1280000,
      loan_to_value_ratio_pct: 74.2
    },
    covenants: [
      "Escrow impounds for property taxes and hazard insurance required",
      "Minimum 12 months post-closing principal and interest reserves"
    ],
    adverse_factors: [
      {
        factor: "Recent Employment Transition",
        description: "Borrower shifted from corporate VP to 1099 consulting 14 months ago.",
        severity: "medium",
        mitigation_notes: "Borrower maintains 3-year signed retainer agreements totaling $300k annually."
      }
    ]
  },
  aml_pep_screening: {
    id: "aml_pep_screening",
    title: "Cross-Border Energy Asset Trade ($4,200,000)",
    domain: "compliance_aml",
    risk_tier: "tier_1_high",
    transaction: {
      amount_usd: 4200000,
      remitter: "Caspian Caspian Energy Holdings",
      beneficiary: "Bosphorus Petroleum Logistics Ltd",
      settlement_bank: "Standard Chartered Singapore"
    },
    pep_findings: [
      {
        entity: "Bosphorus Petroleum Logistics Ltd",
        flag: "Beneficial owner is direct relative of foreign minister of energy.",
        severity: "critical"
      }
    ],
    regulatory_mandates: [
      "Enhanced Due Diligence (EDD) mandatory under FinCEN 31 CFR § 1010.620",
      "Senior management sign-off required prior to funds clearance"
    ]
  }
};

export const DossierBuilder: React.FC<DossierBuilderProps> = ({ onRunCustom, isRunning }) => {
  const [jsonText, setJsonText] = useState<string>(
    JSON.stringify(TEMPLATES.jumbo_mortgage_950k, null, 2)
  );
  const [parseError, setParseError] = useState<string | null>(null);

  const handleTemplateSelect = (key: string) => {
    if (TEMPLATES[key]) {
      setJsonText(JSON.stringify(TEMPLATES[key], null, 2));
      setParseError(null);
    }
  };

  const handleFileUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = (event) => {
      try {
        const content = event.target?.result as string;
        // Validate if json
        const parsed = JSON.parse(content);
        setJsonText(JSON.stringify(parsed, null, 2));
        setParseError(null);
      } catch (err: any) {
        setParseError('Uploaded file is not valid JSON: ' + err.message);
      }
    };
    reader.readAsText(file);
  };

  const handleRun = () => {
    try {
      const parsed = JSON.parse(jsonText);
      setParseError(null);
      onRunCustom(parsed);
    } catch (err: any) {
      setParseError('Syntax Error in Dossier JSON: ' + err.message);
    }
  };

  return (
    <div className="rounded-xl border border-border bg-surface p-5 shadow-lg space-y-4">
      <div className="flex flex-col md:flex-row md:items-center md:justify-between border-b border-border/50 pb-3 gap-2">
        <div>
          <h2 className="text-sm font-semibold tracking-wide text-white uppercase font-mono">
            Interactive Dossier Builder & Ingestion Playground
          </h2>
          <p className="text-xs text-slate-400">
            Author, load enterprise templates, or upload custom evaluation dossiers (.json)
          </p>
        </div>

        {/* Templates */}
        <div className="flex flex-wrap items-center gap-2">
          <span className="text-xs text-slate-400 font-mono">Templates:</span>
          <button
            onClick={() => handleTemplateSelect('jumbo_mortgage_950k')}
            className="rounded bg-surface-elevated hover:bg-slate-700 px-2.5 py-1 text-[11px] font-mono text-cyan-300 border border-border"
          >
            Jumbo Mortgage ($950k)
          </button>
          <button
            onClick={() => handleTemplateSelect('aml_pep_screening')}
            className="rounded bg-surface-elevated hover:bg-slate-700 px-2.5 py-1 text-[11px] font-mono text-rose-300 border border-border"
          >
            AML PEP Trade ($4.2M)
          </button>
        </div>
      </div>

      {parseError && (
        <div className="p-3 rounded-lg bg-rose-950/80 border border-rose-600 text-rose-200 text-xs font-mono">
          {parseError}
        </div>
      )}

      {/* Upload button & Editor */}
      <div className="flex items-center justify-between text-xs font-mono">
        <label className="flex items-center space-x-2 cursor-pointer rounded bg-surface-elevated hover:bg-slate-700 px-3 py-1.5 border border-border text-slate-300">
          <svg className="h-4 w-4 text-cyan-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"></path>
          </svg>
          <span>Upload Custom JSON Dossier</span>
          <input type="file" accept=".json,.txt" onChange={handleFileUpload} className="hidden" />
        </label>

        <span className="text-slate-500 text-[11px]">JSON Schema Verified</span>
      </div>

      <textarea
        rows={12}
        value={jsonText}
        onChange={(e) => setJsonText(e.target.value)}
        className="w-full rounded-lg bg-background border border-border p-4 font-mono text-xs text-cyan-200/90 leading-relaxed focus:outline-none focus:ring-1 focus:ring-cyan-500"
      />

      <div className="flex justify-end">
        <button
          onClick={handleRun}
          disabled={isRunning}
          className="rounded-lg bg-gradient-to-r from-cyan-600 to-blue-600 hover:from-cyan-500 hover:to-blue-500 px-6 py-2 text-xs font-bold text-white shadow-md disabled:opacity-50 transition-all flex items-center space-x-2"
        >
          <span>Execute Deliberation on Ingested Dossier</span>
        </button>
      </div>
    </div>
  );
};
