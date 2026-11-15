import React, { useState } from 'react';
import { Ledger, Block } from '../types';

interface LedgerVerifierProps {
  ledger: Ledger | null;
}

export const LedgerVerifier: React.FC<LedgerVerifierProps> = ({ ledger }) => {
  const [expandedBlock, setExpandedBlock] = useState<number | null>(null);
  const [verificationResult, setVerificationResult] = useState<{
    valid: boolean;
    block_count: number;
    head_hash: string;
    error?: string;
  } | null>(null);
  const [isVerifying, setIsVerifying] = useState<boolean>(false);
  const [isTampered, setIsTampered] = useState<boolean>(false);

  const toggleBlock = (index: number) => {
    setExpandedBlock(expandedBlock === index ? null : index);
  };

  const handleVerify = async (tamper: boolean = false) => {
    if (!ledger) return;
    setIsVerifying(true);
    setIsTampered(tamper);

    try {
      // Create clone of ledger for verification
      const ledgerToVerify = JSON.parse(JSON.stringify(ledger));
      if (tamper && ledgerToVerify.blocks.length > 1) {
        // Deliberately alter Block #1 payload to simulate an attack
        ledgerToVerify.blocks[1].payload = {
          ...ledgerToVerify.blocks[1].payload,
          tampered_attribute: "unauthorized_modification_exploit"
        };
      }

      const res = await fetch('/api/ledger/verify', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(ledgerToVerify),
      });

      const data = await res.json();
      setVerificationResult(data);
    } catch (err: any) {
      setVerificationResult({
        valid: false,
        block_count: ledger.blocks.length,
        head_hash: 'ERROR',
        error: err.message || 'Verification failed',
      });
    } finally {
      setIsVerifying(false);
    }
  };

  return (
    <div className="rounded-xl border border-border bg-surface p-5 shadow-lg">
      <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4 border-b border-border/50 pb-4 mb-4">
        <div>
          <h2 className="text-sm font-semibold tracking-wide text-white uppercase font-mono">
            Cryptographic SHA-256 Audit Trail
          </h2>
          <p className="text-xs text-slate-400">
            Immutable append-only hash chain linking every prompt, citation, and voting round
          </p>
        </div>

        {/* Action buttons */}
        <div className="flex items-center space-x-2">
          <button
            onClick={() => handleVerify(false)}
            disabled={!ledger || isVerifying}
            className="flex items-center space-x-1.5 rounded-lg bg-emerald-950/80 hover:bg-emerald-900/80 border border-emerald-700/80 px-3 py-1.5 text-xs font-semibold text-emerald-300 transition-colors disabled:opacity-50"
          >
            <span>Verify Hash Chain</span>
          </button>

          <button
            onClick={() => handleVerify(true)}
            disabled={!ledger || isVerifying}
            className="flex items-center space-x-1.5 rounded-lg bg-rose-950/80 hover:bg-rose-900/80 border border-rose-700/80 px-3 py-1.5 text-xs font-semibold text-rose-300 transition-colors disabled:opacity-50"
          >
            <span>Simulate Tamper</span>
          </button>
        </div>
      </div>

      {/* Verification Alert Banner */}
      {verificationResult && (
        <div
          className={`mb-4 p-3 rounded-lg border text-xs font-mono flex items-center justify-between ${
            verificationResult.valid
              ? 'bg-emerald-950/60 border-emerald-600 text-emerald-300'
              : 'bg-rose-950/70 border-rose-500 text-rose-200'
          }`}
        >
          <div className="flex items-center space-x-2">
            <span className="font-bold text-sm">
              {verificationResult.valid ? '[PASS]' : '[TAMPER DETECTED]'}
            </span>
            <span>
              {verificationResult.valid
                ? `Cryptographic integrity confirmed across ${verificationResult.block_count} blocks. Zero hash variance.`
                : `Security alert: ${verificationResult.error || 'Hash mismatch detected in block sequence.'}`}
            </span>
          </div>
          {isTampered && (
            <span className="text-[10px] text-amber-300 uppercase px-2 py-0.5 rounded bg-amber-950 border border-amber-700">
              Simulated Attack Test
            </span>
          )}
        </div>
      )}

      {/* Blocks sequence list */}
      {!ledger || ledger.blocks.length === 0 ? (
        <div className="py-8 text-center text-xs text-slate-500 italic">
          No ledger blocks committed yet. Start a deliberation to generate blocks.
        </div>
      ) : (
        <div className="space-y-3">
          {ledger.blocks.map((block: Block, idx: number) => {
            const isExpanded = expandedBlock === block.index;
            return (
              <div
                key={block.index}
                className="rounded-lg bg-surface-elevated/70 border border-border/80 overflow-hidden"
              >
                {/* Block header line */}
                <div
                  onClick={() => toggleBlock(block.index)}
                  className="flex items-center justify-between px-4 py-2.5 cursor-pointer hover:bg-surface-elevated transition-colors text-xs font-mono"
                >
                  <div className="flex items-center space-x-3">
                    <span className="rounded bg-cyan-950 border border-cyan-800 px-2 py-0.5 text-cyan-300 font-bold">
                      BLOCK #{block.index}
                    </span>
                    <span className="text-slate-300 font-medium">{block.type}</span>
                    <span className="text-slate-500 text-[11px]">Round {block.round}</span>
                  </div>

                  <div className="flex items-center space-x-4">
                    <div className="text-[11px] text-slate-400 hidden sm:block">
                      HASH: <span className="text-cyan-400">{block.hash.substring(0, 16)}...</span>
                    </div>
                    <span className="text-slate-500">{isExpanded ? '[-]' : '[+]'}</span>
                  </div>
                </div>

                {/* Expanded Details */}
                {isExpanded && (
                  <div className="border-t border-border px-4 py-3 bg-background/80 text-xs font-mono space-y-2">
                    <div>
                      <span className="text-slate-500">PREVIOUS HASH: </span>
                      <span className="text-slate-300 break-all">{block.prev_hash}</span>
                    </div>
                    <div>
                      <span className="text-slate-500">BLOCK HASH: </span>
                      <span className="text-cyan-300 break-all">{block.hash}</span>
                    </div>
                    <div>
                      <span className="text-slate-500">TIMESTAMP: </span>
                      <span className="text-slate-300">
                        {new Date(block.timestamp_ms).toISOString()}
                      </span>
                    </div>
                    <div>
                      <span className="text-slate-500 block mb-1">PAYLOAD DATA:</span>
                      <pre className="rounded bg-surface p-2 text-[11px] text-slate-300 overflow-x-auto max-h-48 border border-border">
                        {JSON.stringify(block.payload, null, 2)}
                      </pre>
                    </div>
                  </div>
                )}
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
};
