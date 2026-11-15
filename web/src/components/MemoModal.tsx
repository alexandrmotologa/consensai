import React, { useState } from 'react';

interface MemoModalProps {
  isOpen: boolean;
  onClose: () => void;
  memoMarkdown: string;
}

export const MemoModal: React.FC<MemoModalProps> = ({ isOpen, onClose, memoMarkdown }) => {
  const [copied, setCopied] = useState<boolean>(false);

  if (!isOpen) return null;

  const handleCopy = () => {
    navigator.clipboard.writeText(memoMarkdown);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/75 backdrop-blur-sm p-4">
      <div className="relative w-full max-w-3xl rounded-xl border border-border bg-surface shadow-2xl flex flex-col max-h-[85vh]">
        {/* Header */}
        <div className="flex items-center justify-between border-b border-border px-6 py-4">
          <div className="flex items-center space-x-2">
            <h3 className="text-base font-bold text-white font-mono">Executive Decision Memo</h3>
            <span className="rounded bg-cyan-950 px-2 py-0.5 text-xs text-cyan-400 font-mono">
              MARKDOWN
            </span>
          </div>
          <div className="flex items-center space-x-2">
            <button
              onClick={handleCopy}
              className="rounded bg-surface-elevated hover:bg-slate-700 px-3 py-1 text-xs text-slate-200 font-mono transition-colors"
            >
              {copied ? 'Copied!' : 'Copy to Clipboard'}
            </button>
            <button
              onClick={onClose}
              className="rounded p-1 text-slate-400 hover:text-white transition-colors"
            >
              <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12"></path>
              </svg>
            </button>
          </div>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-y-auto p-6 font-mono text-xs text-slate-200 leading-relaxed space-y-4">
          <pre className="whitespace-pre-wrap font-mono text-xs text-slate-200 bg-background/60 p-4 rounded-lg border border-border">
            {memoMarkdown}
          </pre>
        </div>
      </div>
    </div>
  );
};
