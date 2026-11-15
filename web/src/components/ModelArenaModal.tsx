import React, { useState } from 'react';
import { ModelArenaConfig } from '../types';

interface ModelArenaModalProps {
  isOpen: boolean;
  onClose: () => void;
  config: ModelArenaConfig;
  onSave: (config: ModelArenaConfig) => void;
}

const PROVIDERS = [
  { id: 'mock', name: 'Deterministic Mock Engine (Offline / Instant)' },
  { id: 'openai', name: 'OpenAI API (GPT-4o / GPT-4o-mini)' },
  { id: 'anthropic', name: 'Anthropic API (Claude 3.5 Sonnet)' },
  { id: 'gemini', name: 'Google Gemini (Gemini 1.5 Pro)' },
  { id: 'ollama', name: 'Ollama Local (Llama 3.1 / Mistral)' },
];

export const ModelArenaModal: React.FC<ModelArenaModalProps> = ({
  isOpen,
  onClose,
  config,
  onSave,
}) => {
  const [advProvider, setAdvProvider] = useState<string>(config.advocate_provider || 'mock');
  const [advModel, setAdvModel] = useState<string>(config.advocate_model || 'heuristic-v1');

  const [skpProvider, setSkpProvider] = useState<string>(config.skeptic_provider || 'mock');
  const [skpModel, setSkpModel] = useState<string>(config.skeptic_model || 'heuristic-v1');

  const [audProvider, setAudProvider] = useState<string>(config.auditor_provider || 'mock');
  const [audModel, setAudModel] = useState<string>(config.auditor_model || 'heuristic-v1');

  if (!isOpen) return null;

  const handleApplyPreset = (preset: 'mock' | 'cross' | 'openai') => {
    if (preset === 'mock') {
      setAdvProvider('mock');
      setAdvModel('heuristic-v1');
      setSkpProvider('mock');
      setSkpModel('heuristic-v1');
      setAudProvider('mock');
      setAudModel('heuristic-v1');
    } else if (preset === 'cross') {
      setAdvProvider('anthropic');
      setAdvModel('claude-3-5-sonnet-20241022');
      setSkpProvider('openai');
      setSkpModel('gpt-4o');
      setAudProvider('mock');
      setAudModel('heuristic-v1');
    } else if (preset === 'openai') {
      setAdvProvider('openai');
      setAdvModel('gpt-4o');
      setSkpProvider('openai');
      setSkpModel('gpt-4o');
      setAudProvider('openai');
      setAudModel('gpt-4o');
    }
  };

  const handleSave = () => {
    onSave({
      advocate_provider: advProvider,
      advocate_model: advModel,
      skeptic_provider: skpProvider,
      skeptic_model: skpModel,
      auditor_provider: audProvider,
      auditor_model: audModel,
    });
    onClose();
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm p-4">
      <div className="relative w-full max-w-2xl rounded-xl border border-border bg-surface shadow-2xl flex flex-col max-h-[90vh]">
        {/* Header */}
        <div className="flex items-center justify-between border-b border-border px-6 py-4">
          <div className="flex items-center space-x-2">
            <span className="flex h-6 w-6 items-center justify-center rounded-full bg-cyan-950 border border-cyan-500 text-cyan-300 font-bold text-xs">
              A
            </span>
            <div>
              <h3 className="text-base font-bold text-white font-mono">
                Multi-Model Adversarial Arena
              </h3>
              <p className="text-xs text-slate-400">Eliminate vendor bias by cross-pairing models</p>
            </div>
          </div>
          <button onClick={onClose} className="rounded p-1 text-slate-400 hover:text-white">
            <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>

        {/* Body */}
        <div className="flex-1 overflow-y-auto p-6 space-y-5 text-xs">
          {/* Presets */}
          <div className="flex items-center space-x-2 border-b border-border/60 pb-3">
            <span className="text-slate-400 font-mono">Presets:</span>
            <button
              onClick={() => handleApplyPreset('mock')}
              className="rounded bg-surface-elevated hover:bg-slate-700 px-2.5 py-1 text-[11px] font-mono text-cyan-300 border border-border"
            >
              Offline Deterministic
            </button>
            <button
              onClick={() => handleApplyPreset('cross')}
              className="rounded bg-surface-elevated hover:bg-slate-700 px-2.5 py-1 text-[11px] font-mono text-emerald-300 border border-border"
            >
              Cross-Vendor (Claude vs GPT-4o)
            </button>
            <button
              onClick={() => handleApplyPreset('openai')}
              className="rounded bg-surface-elevated hover:bg-slate-700 px-2.5 py-1 text-[11px] font-mono text-slate-300 border border-border"
            >
              All OpenAI
            </button>
          </div>

          {/* Advocate Role */}
          <div className="rounded-lg bg-surface-elevated/70 border border-emerald-900/40 p-4 space-y-3">
            <div className="flex items-center space-x-2">
              <span className="h-2 w-2 rounded-full bg-emerald-400"></span>
              <span className="font-mono font-bold text-emerald-400 uppercase">The Advocate</span>
              <span className="text-slate-400">(Affirmative Defense)</span>
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div>
                <label className="block text-slate-400 mb-1">Provider</label>
                <select
                  value={advProvider}
                  onChange={(e) => setAdvProvider(e.target.value)}
                  className="w-full rounded bg-surface border border-border px-3 py-1.5 text-slate-200 font-mono"
                >
                  {PROVIDERS.map((p) => (
                    <option key={p.id} value={p.id}>
                      {p.name}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-slate-400 mb-1">Model Name</label>
                <input
                  type="text"
                  value={advModel}
                  onChange={(e) => setAdvModel(e.target.value)}
                  className="w-full rounded bg-surface border border-border px-3 py-1.5 text-slate-200 font-mono"
                />
              </div>
            </div>
          </div>

          {/* Skeptic Role */}
          <div className="rounded-lg bg-surface-elevated/70 border border-rose-900/40 p-4 space-y-3">
            <div className="flex items-center space-x-2">
              <span className="h-2 w-2 rounded-full bg-rose-400"></span>
              <span className="font-mono font-bold text-rose-400 uppercase">The Skeptic</span>
              <span className="text-slate-400">(Adversarial Cross-Examiner)</span>
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div>
                <label className="block text-slate-400 mb-1">Provider</label>
                <select
                  value={skpProvider}
                  onChange={(e) => setSkpProvider(e.target.value)}
                  className="w-full rounded bg-surface border border-border px-3 py-1.5 text-slate-200 font-mono"
                >
                  {PROVIDERS.map((p) => (
                    <option key={p.id} value={p.id}>
                      {p.name}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-slate-400 mb-1">Model Name</label>
                <input
                  type="text"
                  value={skpModel}
                  onChange={(e) => setSkpModel(e.target.value)}
                  className="w-full rounded bg-surface border border-border px-3 py-1.5 text-slate-200 font-mono"
                />
              </div>
            </div>
          </div>

          {/* Risk Auditor Role */}
          <div className="rounded-lg bg-surface-elevated/70 border border-amber-900/40 p-4 space-y-3">
            <div className="flex items-center space-x-2">
              <span className="h-2 w-2 rounded-full bg-amber-400"></span>
              <span className="font-mono font-bold text-amber-400 uppercase">The Risk Auditor</span>
              <span className="text-slate-400">(Forensic Grounding Verifier)</span>
            </div>
            <div className="grid grid-cols-2 gap-3">
              <div>
                <label className="block text-slate-400 mb-1">Provider</label>
                <select
                  value={audProvider}
                  onChange={(e) => setAudProvider(e.target.value)}
                  className="w-full rounded bg-surface border border-border px-3 py-1.5 text-slate-200 font-mono"
                >
                  {PROVIDERS.map((p) => (
                    <option key={p.id} value={p.id}>
                      {p.name}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-slate-400 mb-1">Model Name</label>
                <input
                  type="text"
                  value={audModel}
                  onChange={(e) => setAudModel(e.target.value)}
                  className="w-full rounded bg-surface border border-border px-3 py-1.5 text-slate-200 font-mono"
                />
              </div>
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="border-t border-border px-6 py-4 flex justify-end space-x-2">
          <button
            onClick={onClose}
            className="rounded bg-surface-elevated hover:bg-slate-700 px-4 py-2 font-medium text-slate-300 text-xs"
          >
            Cancel
          </button>
          <button
            onClick={handleSave}
            className="rounded bg-gradient-to-r from-cyan-600 to-blue-600 hover:from-cyan-500 hover:to-blue-500 px-5 py-2 font-bold text-white text-xs shadow-md"
          >
            Apply Arena Configuration
          </button>
        </div>
      </div>
    </div>
  );
};
