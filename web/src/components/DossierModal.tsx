import React from 'react';

interface DossierModalProps {
  isOpen: boolean;
  onClose: () => void;
  dossierContent: string;
  dossierTitle: string;
}

export const DossierModal: React.FC<DossierModalProps> = ({
  isOpen,
  onClose,
  dossierContent,
  dossierTitle,
}) => {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/75 backdrop-blur-sm p-4">
      <div className="relative w-full max-w-3xl rounded-xl border border-border bg-surface shadow-2xl flex flex-col max-h-[85vh]">
        {/* Header */}
        <div className="flex items-center justify-between border-b border-border px-6 py-4">
          <div>
            <h3 className="text-base font-bold text-white font-mono">Dossier Ground Truth</h3>
            <p className="text-xs text-slate-400">{dossierTitle}</p>
          </div>
          <button
            onClick={onClose}
            className="rounded p-1 text-slate-400 hover:text-white transition-colors"
          >
            <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-y-auto p-6 font-mono text-xs text-slate-200">
          <pre className="whitespace-pre-wrap font-mono text-xs text-cyan-200/90 bg-background/60 p-4 rounded-lg border border-border">
            {dossierContent}
          </pre>
        </div>
      </div>
    </div>
  );
};
