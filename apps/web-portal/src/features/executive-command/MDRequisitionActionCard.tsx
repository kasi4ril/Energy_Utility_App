/**
 * Requisition Action & Review Card for Managing Director.
 * Provides:
 * 1. Room for MD Review & Technical Queries.
 * 2. Room for Inter-Departmental Minuting & Acknowledgment.
 * 3. Authoritative Digital Approval & Authorize Button.
 * Strictly <= 350 lines.
 */

import React, { useState } from "react";
import { ExecutiveRequisition, DepartmentSource } from "./types";
import {
  CheckCircle,
  FileText,
  AlertTriangle,
  GitFork,
  Send,
  XCircle,
  Clock,
  ShieldCheck,
} from "lucide-react";

interface ActionCardProps {
  requisition: ExecutiveRequisition;
  onApprove: (id: string, notes: string) => void;
  onReview: (id: string, notes: string) => void;
  onMinute: (id: string, departments: DepartmentSource[], notes: string) => void;
  onReject: (id: string, notes: string) => void;
}

const DEPARTMENTS: { id: DepartmentSource; label: string }[] = [
  { id: "GENERATION", label: "GM Generation" },
  { id: "ENGINEERING_SERVICES", label: "Head Engineering" },
  { id: "COMMERCIAL_OPS", label: "Commercial Ops" },
  { id: "INTERNAL_AUDIT", label: "Internal Audit" },
  { id: "STORES_PROCUREMENT", label: "Stores & Stock" },
  { id: "HSE_SAFETY", label: "HSE Safety" },
];

export const MDRequisitionActionCard: React.FC<ActionCardProps> = ({
  requisition,
  onApprove,
  onReview,
  onMinute,
  onReject,
}) => {
  const [reviewNotes, setReviewNotes] = useState(requisition.md_review_notes || "");
  const [minuteNotes, setMinuteNotes] = useState(requisition.minuted_notes || "");
  const [selectedDepts, setSelectedDepts] = useState<DepartmentSource[]>(
    requisition.minuted_to_departments || []
  );

  const toggleDept = (dept: DepartmentSource) => {
    setSelectedDepts((prev) =>
      prev.includes(dept) ? prev.filter((d) => d !== dept) : [...prev, dept]
    );
  };

  const isApproved = requisition.status === "APPROVED";
  const isRejected = requisition.status === "REJECTED";
  const isMinuted = requisition.status === "ACKNOWLEDGED_MINUTED";

  return (
    <div className="bg-slate-900 border border-slate-800 rounded-xl p-5 shadow-xl space-y-4">
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-start justify-between pb-3 border-b border-slate-800 gap-3">
        <div>
          <div className="flex items-center space-x-2">
            <span className="font-mono text-xs font-bold text-amber-400 bg-amber-500/10 px-2 py-0.5 rounded border border-amber-500/20">
              {requisition.request_no}
            </span>
            <span
              className={`text-[10px] font-bold px-2 py-0.5 rounded uppercase font-mono ${
                requisition.urgency === "EMERGENCY_OUTAGE"
                  ? "bg-rose-500/20 text-rose-300 border border-rose-500/40 animate-pulse"
                  : requisition.urgency === "URGENT"
                  ? "bg-amber-500/20 text-amber-300"
                  : "bg-slate-700 text-slate-300"
              }`}
            >
              {requisition.urgency.replace(/_/g, " ")}
            </span>
            <span
              className={`text-[10px] font-bold px-2 py-0.5 rounded uppercase ${
                isApproved
                  ? "bg-emerald-500/20 text-emerald-300 border border-emerald-500/30"
                  : isRejected
                  ? "bg-rose-500/20 text-rose-300"
                  : isMinuted
                  ? "bg-cyan-500/20 text-cyan-300"
                  : "bg-amber-500/20 text-amber-300"
              }`}
            >
              {requisition.status.replace(/_/g, " ")}
            </span>
          </div>
          <h2 className="text-base font-bold text-white mt-1.5 leading-snug">
            {requisition.title}
          </h2>
          <div className="text-[11px] text-slate-400 mt-0.5">
            Submitted by: <strong className="text-slate-200">{requisition.submitting_head_name}</strong> ({requisition.submitting_head_role}) &bull;{" "}
            {new Date(requisition.submission_date).toLocaleString("en-GB")} WAT
          </div>
        </div>

        <div className="text-left sm:text-right shrink-0">
          <div className="text-[10px] text-slate-400 uppercase font-bold">Total Financial Commitment</div>
          <div className="text-xl font-black font-mono text-emerald-400">
            ₦{requisition.amount.toLocaleString()}
          </div>
          <div className="text-[10px] text-slate-400 font-mono mt-0.5">
            Category: <span className="text-slate-300">{requisition.category.replace(/_/g, " ")}</span>
          </div>
        </div>
      </div>

      {/* Technical Justification & Summary */}
      <div className="bg-slate-950/70 border border-slate-800 rounded-lg p-3.5 space-y-2 text-xs">
        <div>
          <span className="font-bold text-slate-300 uppercase text-[10px]">Executive Summary:</span>
          <p className="text-slate-200 mt-0.5">{requisition.summary}</p>
        </div>
        <div>
          <span className="font-bold text-slate-400 uppercase text-[10px]">Technical Justification &amp; Grid Impact:</span>
          <p className="text-slate-300 mt-0.5 leading-relaxed">{requisition.technical_justification}</p>
        </div>
      </div>

      {/* Itemized Line Items Table */}
      {requisition.line_items.length > 0 && (
        <div className="border border-slate-800 rounded-lg overflow-hidden">
          <div className="bg-slate-950 px-3 py-2 text-[11px] font-bold text-slate-400 uppercase flex items-center justify-between">
            <span className="flex items-center space-x-1.5">
              <FileText className="w-3.5 h-3.5 text-amber-400" />
              <span>Itemized Bill of Materials &amp; Spares ({requisition.line_items.length})</span>
            </span>
            <span className="font-mono text-slate-300">Total: ₦{requisition.amount.toLocaleString()}</span>
          </div>
          <table className="w-full text-left text-xs font-mono">
            <thead className="bg-slate-900/80 text-slate-400 uppercase text-[10px] border-b border-slate-800">
              <tr>
                <th className="p-2 font-sans">Item Description</th>
                <th className="p-2">Qty</th>
                <th className="p-2">Unit Price</th>
                <th className="p-2 text-right">Total Price</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/60 text-slate-300">
              {requisition.line_items.map((item) => (
                <tr key={item.id} className="hover:bg-slate-800/30">
                  <td className="p-2 font-sans font-medium text-slate-200">{item.item}</td>
                  <td className="p-2">{item.quantity} {item.unit}</td>
                  <td className="p-2">₦{item.unit_price.toLocaleString()}</td>
                  <td className="p-2 text-right font-bold text-emerald-400">₦{item.total_price.toLocaleString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* 1. Room for MD Review & Directives */}
      <div className="bg-slate-950/80 border border-slate-800 rounded-lg p-3.5 space-y-2">
        <div className="flex items-center justify-between text-xs">
          <span className="font-bold text-amber-300 flex items-center space-x-1.5">
            <Clock className="w-3.5 h-3.5 text-amber-400" />
            <span>Room for MD Review &amp; Executive Directives</span>
          </span>
          <span className="text-[10px] text-slate-400">Recorded in Immutable Audit Log</span>
        </div>
        <textarea
          value={reviewNotes}
          onChange={(e) => setReviewNotes(e.target.value)}
          placeholder="Enter MD technical observations, operational queries, conditions of approval, or guidance to the department head..."
          rows={2}
          disabled={isApproved || isRejected}
          className="w-full bg-slate-900 border border-slate-700 rounded-lg p-2.5 text-xs text-white placeholder-slate-500 focus:outline-none focus:border-amber-500 disabled:opacity-60"
        />
        {!isApproved && !isRejected && (
          <div className="flex justify-end">
            <button
              onClick={() => onReview(requisition.id, reviewNotes)}
              disabled={!reviewNotes.trim()}
              className="px-3 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-amber-300 border border-amber-500/30 text-xs font-semibold flex items-center space-x-1.5 transition disabled:opacity-50"
            >
              <Send className="w-3 h-3 text-amber-400" />
              <span>Save Review &amp; Query Department Head</span>
            </button>
          </div>
        )}
      </div>

      {/* 2. Room for Inter-Departmental Acknowledgment & Minuting */}
      <div className="bg-slate-950/80 border border-slate-800 rounded-lg p-3.5 space-y-2.5">
        <div className="flex items-center justify-between text-xs">
          <span className="font-bold text-cyan-300 flex items-center space-x-1.5">
            <GitFork className="w-3.5 h-3.5 text-cyan-400" />
            <span>Inter-Departmental Minuting &amp; Relationship Acknowledgment</span>
          </span>
          <span className="text-[10px] text-slate-400">Brings other Heads into the picture</span>
        </div>

        <p className="text-[11px] text-slate-400">
          Select related departments that need to be informed, coordinated with, or assigned complementary tasks:
        </p>

        <div className="flex flex-wrap gap-1.5">
          {DEPARTMENTS.map((dept) => {
            const isSelected = selectedDepts.includes(dept.id);
            return (
              <button
                key={dept.id}
                type="button"
                onClick={() => toggleDept(dept.id)}
                disabled={isApproved || isRejected}
                className={`px-2.5 py-1 rounded-md text-[11px] font-semibold transition border ${
                  isSelected
                    ? "bg-cyan-600 text-white border-cyan-500 shadow"
                    : "bg-slate-900 text-slate-400 border-slate-700 hover:text-slate-200"
                }`}
              >
                {isSelected ? `✓ ${dept.label}` : `+ ${dept.label}`}
              </button>
            );
          })}
        </div>

        <input
          type="text"
          value={minuteNotes}
          onChange={(e) => setMinuteNotes(e.target.value)}
          placeholder="Directive minuted across departments (e.g. 'Audit to verify quotes; Central Dispatch to isolate feeder')..."
          disabled={isApproved || isRejected}
          className="w-full bg-slate-900 border border-slate-700 rounded-lg px-2.5 py-1.5 text-xs text-white placeholder-slate-500 focus:outline-none focus:border-cyan-500 disabled:opacity-60"
        />

        {!isApproved && !isRejected && (
          <div className="flex justify-end">
            <button
              onClick={() => onMinute(requisition.id, selectedDepts, minuteNotes)}
              disabled={selectedDepts.length === 0}
              className="px-3 py-1.5 rounded-lg bg-cyan-700 hover:bg-cyan-600 text-white text-xs font-semibold flex items-center space-x-1.5 transition disabled:opacity-50"
            >
              <GitFork className="w-3 h-3 text-cyan-200" />
              <span>Acknowledge &amp; Minute to Selected Departments</span>
            </button>
          </div>
        )}
      </div>

      {/* 3. Authoritative Approval & Rejection Decision Bar */}
      <div className="pt-2 flex flex-col sm:flex-row sm:items-center justify-between gap-3 border-t border-slate-800">
        <div className="flex items-center space-x-2 text-xs text-slate-400 font-mono">
          <ShieldCheck className="w-4 h-4 text-emerald-400" />
          <span>MD Authority Limit: <strong>₦50,000,000</strong> per Requisition</span>
        </div>

        <div className="flex items-center space-x-2">
          {!isApproved && !isRejected ? (
            <>
              <button
                onClick={() => onReject(requisition.id, reviewNotes || "Declined by MD")}
                className="px-3.5 py-2 rounded-lg bg-rose-950/80 hover:bg-rose-900 text-rose-300 border border-rose-700/60 font-semibold text-xs flex items-center space-x-1.5 transition"
              >
                <XCircle className="w-3.5 h-3.5 text-rose-400" />
                <span>Decline / Return</span>
              </button>

              <button
                onClick={() => onApprove(requisition.id, reviewNotes || "Executive approval granted.")}
                className="px-5 py-2 rounded-lg bg-gradient-to-r from-emerald-600 to-emerald-500 hover:from-emerald-500 hover:to-emerald-400 text-white font-bold text-xs flex items-center space-x-2 shadow-lg shadow-emerald-950/50 transition active:scale-95"
              >
                <CheckCircle className="w-4 h-4 text-white" />
                <span>Approve &amp; Authorize Expenditure</span>
              </button>
            </>
          ) : isApproved ? (
            <div className="flex items-center space-x-2 px-3.5 py-1.5 rounded-lg bg-emerald-500/20 text-emerald-300 border border-emerald-500/40 text-xs font-mono font-bold">
              <CheckCircle className="w-4 h-4 text-emerald-400" />
              <span>Authorized by {requisition.decision_by || "MD/CEO"}</span>
            </div>
          ) : (
            <div className="flex items-center space-x-2 px-3.5 py-1.5 rounded-lg bg-rose-500/20 text-rose-300 border border-rose-500/40 text-xs font-mono font-bold">
              <XCircle className="w-4 h-4 text-rose-400" />
              <span>Declined / Returned for Revision</span>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

