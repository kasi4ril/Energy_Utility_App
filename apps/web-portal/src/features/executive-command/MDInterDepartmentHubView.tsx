/**
 * Workspace 3: Inter-Departmental Information & Minuting Hub.
 * Manages cross-departmental communications, joint technical memos, and operational
 * advisories where the Managing Director needs to be in the picture.
 * Strictly <= 350 lines.
 */

import React, { useState } from "react";
import { ExecutiveInterDeptMemo, DepartmentSource } from "./types";
import {
  GitFork,
  FileText,
  AlertTriangle,
  Clock,
  Send,
  CheckCircle,
  Building2,
  Calendar,
} from "lucide-react";

interface InterDeptHubProps {
  memos: ExecutiveInterDeptMemo[];
  onMinuteMemo: (memoId: string, remarks: string) => void;
}

const DEPT_LABELS: Record<DepartmentSource, string> = {
  GENERATION: "Generation (Hydro Fleet)",
  ENGINEERING_SERVICES: "Engineering & Distribution",
  COMMERCIAL_OPS: "Commercial Operations",
  INTERNAL_AUDIT: "Internal Audit",
  STORES_PROCUREMENT: "Central Stores",
  HSE_SAFETY: "HSE Safety",
};

export const MDInterDepartmentHubView: React.FC<InterDeptHubProps> = ({
  memos,
  onMinuteMemo,
}) => {
  const [selectedMemoId, setSelectedMemoId] = useState<string>(memos[0]?.id || "");
  const [minuteText, setMinuteText] = useState<string>("");

  const activeMemo = memos.find((m) => m.id === selectedMemoId) || memos[0];

  const handleSaveMinute = () => {
    if (!activeMemo || !minuteText.trim()) return;
    onMinuteMemo(activeMemo.id, minuteText);
    setMinuteText("");
  };

  return (
    <div className="space-y-6 text-xs text-slate-100">
      {/* Header Banner */}
      <div className="bg-gradient-to-r from-slate-900 via-cyan-950/40 to-slate-900 border border-cyan-500/30 rounded-2xl p-5 shadow-xl">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <div className="flex items-center space-x-2 mb-1.5">
              <span className="px-2.5 py-0.5 rounded-full text-[10px] font-bold tracking-wider uppercase bg-cyan-500/20 border border-cyan-500/40 text-cyan-300">
                Cross-Unit Directives &bull; Executive Minuting
              </span>
              <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-emerald-500/20 text-emerald-300 border border-emerald-500/30">
                MD In The Picture
              </span>
            </div>
            <h1 className="text-xl sm:text-2xl font-black text-white tracking-tight">
              Inter-Departmental Information &amp; Minuting Hub
            </h1>
            <p className="text-xs text-slate-300 mt-1 max-w-2xl leading-relaxed">
              Operational advisories, joint engineering-generation protocols, and regulatory communications
              requiring executive cognizance, inter-departmental harmonization, and minuted directives.
            </p>
          </div>

          <div className="bg-slate-950/80 border border-cyan-500/40 px-3.5 py-2 rounded-xl text-right shrink-0">
            <div className="text-[10px] text-cyan-400 font-bold uppercase tracking-wider">Active Memorandums</div>
            <div className="text-lg font-black font-mono text-white">
              {memos.length} Memos
            </div>
          </div>
        </div>
      </div>

      {/* Main Grid: Memos List (Left) + Detail & Minuting Panel (Right) */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-5 items-start">
        {/* Left Column: Memo Selector List */}
        <div className="lg:col-span-5 space-y-2.5">
          <div className="text-[11px] font-bold text-slate-400 uppercase tracking-wider px-1">
            Incoming Advisories &amp; Joint Memos ({memos.length})
          </div>

          {memos.map((memo) => {
            const isSelected = memo.id === activeMemo?.id;
            return (
              <div
                key={memo.id}
                onClick={() => setSelectedMemoId(memo.id)}
                className={`p-3.5 rounded-xl border transition-all cursor-pointer space-y-2 ${
                  isSelected
                    ? "bg-slate-800/90 border-cyan-500 shadow-md ring-1 ring-cyan-500/40"
                    : "bg-slate-900/80 border-slate-800 hover:border-slate-700"
                }`}
              >
                <div className="flex items-center justify-between">
                  <span className="font-mono font-bold text-[11px] text-cyan-400">
                    {memo.memo_no}
                  </span>
                  <span
                    className={`text-[9px] font-bold px-1.5 py-0.2 rounded uppercase font-mono ${
                      memo.priority === "CRITICAL"
                        ? "bg-rose-500/20 text-rose-300 border border-rose-500/40"
                        : memo.priority === "HIGH"
                        ? "bg-amber-500/20 text-amber-300"
                        : "bg-slate-700 text-slate-300"
                    }`}
                  >
                    {memo.priority}
                  </span>
                </div>

                <h3 className="font-bold text-white text-xs line-clamp-2 leading-snug">
                  {memo.title}
                </h3>

                <div className="flex items-center justify-between text-[10px] pt-1 border-t border-slate-800/80 text-slate-400">
                  <span className="truncate max-w-[200px]">
                    {memo.originating_departments.map((d) => DEPT_LABELS[d] || d).join(" & ")}
                  </span>
                  <span
                    className={`font-semibold ${
                      memo.status === "NOTED_AND_MINUTED" ? "text-emerald-400" : "text-amber-400"
                    }`}
                  >
                    {memo.status === "NOTED_AND_MINUTED" ? "✓ Minuted" : "Awaiting MD"}
                  </span>
                </div>
              </div>
            );
          })}
        </div>

        {/* Right Column: Active Memo Detail & Minute Box */}
        <div className="lg:col-span-7">
          {activeMemo ? (
            <div className="bg-slate-900 border border-slate-800 rounded-xl p-5 shadow-xl space-y-4">
              {/* Header */}
              <div className="pb-3 border-b border-slate-800">
                <div className="flex items-center justify-between">
                  <span className="font-mono text-xs font-bold text-cyan-400 bg-cyan-500/10 px-2 py-0.5 rounded border border-cyan-500/20">
                    {activeMemo.memo_no}
                  </span>
                  <span className="text-[10px] text-slate-400 flex items-center space-x-1 font-mono">
                    <Calendar className="w-3 h-3 text-cyan-400" />
                    <span>{new Date(activeMemo.date).toLocaleDateString("en-GB")} WAT</span>
                  </span>
                </div>
                <h2 className="text-base font-bold text-white mt-2 leading-snug">
                  {activeMemo.title}
                </h2>
                <div className="text-[11px] text-slate-400 mt-1 flex items-center space-x-2">
                  <Building2 className="w-3.5 h-3.5 text-cyan-400 shrink-0" />
                  <span>
                    Origin:{" "}
                    <strong className="text-slate-200">
                      {activeMemo.originating_departments.map((d) => DEPT_LABELS[d] || d).join(" & ")}
                    </strong>
                  </span>
                </div>
              </div>

              {/* Executive Summary */}
              <div className="bg-slate-950/70 border border-slate-800 rounded-lg p-3.5 space-y-2">
                <div>
                  <span className="font-bold text-slate-300 uppercase text-[10px]">Executive Summary:</span>
                  <p className="text-slate-200 mt-0.5 leading-relaxed">{activeMemo.executive_summary}</p>
                </div>
                <div>
                  <span className="font-bold text-cyan-400 uppercase text-[10px]">Technical &amp; Operational Impact:</span>
                  <p className="text-slate-300 mt-0.5 leading-relaxed">{activeMemo.technical_impact}</p>
                </div>
              </div>

              {/* Existing Remarks if already minuted */}
              {activeMemo.md_minute_remarks && (
                <div className="bg-emerald-950/30 border border-emerald-500/40 rounded-lg p-3.5 space-y-1.5">
                  <div className="flex items-center justify-between text-xs font-bold text-emerald-300">
                    <span className="flex items-center space-x-1.5">
                      <CheckCircle className="w-4 h-4 text-emerald-400" />
                      <span>MD Executive Minute Remarks (Enacted)</span>
                    </span>
                    <span className="text-[10px] font-mono text-emerald-400">Engr. Dr. E. K. Dung</span>
                  </div>
                  <p className="text-slate-200 text-xs italic">"{activeMemo.md_minute_remarks}"</p>
                  {activeMemo.assigned_action_heads && (
                    <div className="text-[10px] text-slate-400 pt-1 border-t border-emerald-900/60">
                      Assigned Action Heads:{" "}
                      <strong className="text-emerald-300">{activeMemo.assigned_action_heads.join(", ")}</strong>
                    </div>
                  )}
                </div>
              )}

              {/* MD Minuting Directives Input */}
              <div className="bg-slate-950/80 border border-slate-800 rounded-lg p-3.5 space-y-2.5">
                <div className="flex items-center justify-between text-xs font-bold text-cyan-300">
                  <span className="flex items-center space-x-1.5">
                    <Clock className="w-3.5 h-3.5 text-cyan-400" />
                    <span>Minute Executive Directive to Concerned Departments</span>
                  </span>
                  <span className="text-[10px] text-slate-400">Notifies Heads Instantly</span>
                </div>

                <textarea
                  value={minuteText}
                  onChange={(e) => setMinuteText(e.target.value)}
                  placeholder="Enter executive directive minuted to concerned departmental heads (e.g., 'Noted. GM Generation and Head Engineering Services to convene readiness test by Friday morning')..."
                  rows={3}
                  className="w-full bg-slate-900 border border-slate-700 rounded-lg p-2.5 text-xs text-white placeholder-slate-500 focus:outline-none focus:border-cyan-500"
                />

                <div className="flex justify-end">
                  <button
                    onClick={handleSaveMinute}
                    disabled={!minuteText.trim()}
                    className="px-4 py-2 rounded-lg bg-gradient-to-r from-cyan-600 to-cyan-500 hover:from-cyan-500 hover:to-cyan-400 text-white font-bold text-xs flex items-center space-x-1.5 shadow transition disabled:opacity-50"
                  >
                    <Send className="w-3.5 h-3.5" />
                    <span>Stamp Directive &amp; Minute to Departments</span>
                  </button>
                </div>
              </div>
            </div>
          ) : (
            <div className="bg-slate-900 border border-slate-800 rounded-xl p-12 text-center text-slate-500">
              Select a memo from the list to view and minute directives.
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

