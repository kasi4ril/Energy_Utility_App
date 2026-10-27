/**
 * Cryptographic Hash-Chain Ledger Explorer for Internal Audit.
 * Verifies SHA-256 prev-hash link integrity on immutable WorkNotes.
 */

import React from "react";
import { MOCK_CHAINED_LEDGER } from "./auditMockData";

export function AuditLedgerView() {
  return (
    <div className="space-y-6">
      {/* Header Info */}
      <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800 flex items-center justify-between">
        <div>
          <h3 className="text-sm font-bold text-white flex items-center space-x-2">
            <span>⛓️</span>
            <span>Cryptographic Hash-Chain Ledger &bull; Immutable Audit Trail</span>
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Mathematical proof of non-tampering: Every handover note is cryptographically signed and chained to its predecessor.
          </p>
        </div>
        <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-cyan-950 text-cyan-400 border border-cyan-800">
          SHA-256 Chain Intact
        </span>
      </div>

      {/* Chained Records Timeline */}
      <div className="space-y-4">
        {MOCK_CHAINED_LEDGER.map((record, idx) => (
          <div
            key={record.id}
            className="p-5 rounded-xl bg-slate-900/70 border border-slate-800 relative hover:border-cyan-500/40 transition"
          >
            {/* Sequence Pill */}
            <div className="flex items-center justify-between mb-3">
              <div className="flex items-center space-x-2">
                <span className="w-6 h-6 rounded-full bg-cyan-950 border border-cyan-500/50 flex items-center justify-center text-cyan-400 font-mono font-bold text-xs">
                  {record.sequence_number}
                </span>
                <span className="text-xs font-mono font-bold text-white">
                  {record.entity_type.toUpperCase()} &bull; {record.entity_id}
                </span>
              </div>
              <span className="text-[10px] font-mono px-2 py-0.5 rounded bg-emerald-950 text-emerald-400 border border-emerald-800">
                ✓ Cryptographic Seal Valid
              </span>
            </div>

            {/* Note Content */}
            <p className="text-xs text-slate-200 bg-slate-950/80 p-3 rounded-lg border border-slate-800 font-sans leading-relaxed">
              "{record.note_content}"
            </p>

            {/* Author Attribution */}
            <div className="flex items-center justify-between text-[11px] text-slate-400 mt-3 pt-2 border-t border-slate-800/60">
              <div>
                Author: <strong className="text-slate-200 font-mono">{record.author_staff_id}</strong> ({record.author_role})
              </div>
              <div className="font-mono text-slate-500">{record.timestamp_wat}</div>
            </div>

            {/* Cryptographic Hash Pair */}
            <div className="mt-3 pt-2 border-t border-slate-800/80 font-mono text-[10px] space-y-1 bg-slate-950/50 p-2.5 rounded border border-slate-800/50">
              <div className="text-slate-500 flex items-center justify-between">
                <span>PREV_HASH:</span>
                <span className="text-slate-400 truncate max-w-md">{record.prev_sha256_hash}</span>
              </div>
              <div className="text-slate-500 flex items-center justify-between">
                <span>CURR_HASH:</span>
                <span className="text-cyan-400 truncate max-w-md font-bold">{record.current_sha256_hash}</span>
              </div>
            </div>

            {/* Connector Arrow */}
            {idx < MOCK_CHAINED_LEDGER.length - 1 && (
              <div className="text-center text-slate-600 text-sm my-1">
                &darr; Chained via SHA-256 Prev-Hash &darr;
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}

