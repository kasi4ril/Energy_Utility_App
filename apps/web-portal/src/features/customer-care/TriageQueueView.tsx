/**
 * Task View 1: Active Triage Queue & SLA Countdown.
 * Features live queue, SLA status tags, and 2-step closure verification modal.
 */

import React, { useState } from "react";
import { ServiceTicket } from "./types";
import { Badge } from "../../components/ui/Badge";
import { Table, TableHead, TableBody, TableRow, TableHeaderCell, TableCell } from "../../components/ui/Table";

interface TriageQueueViewProps {
  tickets: ServiceTicket[];
  onSelectTicket: (ticket: ServiceTicket) => void;
  onConfirmClose: (ticketId: string, notes: string) => Promise<void>;
  onNavigateToLogComplaint?: () => void;
}

export const TriageQueueView: React.FC<TriageQueueViewProps> = ({
  tickets,
  onSelectTicket,
  onConfirmClose,
  onNavigateToLogComplaint,
}) => {
  const [filter, setFilter] = useState<string>("ALL");
  const [search, setSearch] = useState<string>("");
  const [confirmingTicket, setConfirmingTicket] = useState<ServiceTicket | null>(null);
  const [feedbackNotes, setFeedbackNotes] = useState<string>("");
  const [submittingClose, setSubmittingClose] = useState<boolean>(false);

  const filteredTickets = tickets.filter((t) => {
    if (filter === "ATTENDED_TO" && t.status !== "ATTENDED_TO") return false;
    if (filter === "OPEN" && t.status !== "ROUTED" && t.status !== "LOGGED") return false;
    if (filter === "CLOSED" && t.status !== "CLOSED") return false;

    if (search.trim()) {
      const q = search.toLowerCase();
      const matchNo = t.ticket_number.toLowerCase().includes(q);
      const matchCust = (t.customer_name || "").toLowerCase().includes(q);
      const matchCat = t.category.toLowerCase().includes(q);
      const matchDesc = t.description.toLowerCase().includes(q);
      return matchNo || matchCust || matchCat || matchDesc;
    }
    return true;
  });

  const handleExecuteClose = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!confirmingTicket) return;
    setSubmittingClose(true);
    try {
      await onConfirmClose(confirmingTicket.id, feedbackNotes || "Customer verified complete restoration.");
      setConfirmingTicket(null);
      setFeedbackNotes("");
    } catch {
      // handled in parent
    } finally {
      setSubmittingClose(false);
    }
  };

  return (
    <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-6 shadow-xl space-y-5">
      {/* Header Bar */}
      <div className="flex flex-col md:flex-row md:items-center justify-between pb-4 border-b border-slate-800 gap-4">
        <div>
          <div className="flex items-center gap-2">
            <h2 className="text-base font-bold text-white tracking-wide">
              Active Triage Queue &bull; NERC SLA Countdown
            </h2>
            <span className="px-2 py-0.5 rounded text-[11px] font-semibold bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
              Live Queue
            </span>
          </div>
          <p className="text-xs text-slate-400 mt-1">
            Real-time feed of logged faults. Tickets in ATTENDED_TO require two-step verbal customer sign-off to close.
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-2.5">
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search tickets, customers, faults..."
            className="px-3 py-1.5 rounded-lg bg-slate-950 border border-slate-700 text-xs text-white placeholder-slate-500 focus:outline-none focus:ring-1 focus:ring-emerald-500"
          />

          <div className="flex rounded-lg bg-slate-950 p-0.5 border border-slate-800">
            {["ALL", "OPEN", "ATTENDED_TO", "CLOSED"].map((tab) => (
              <button
                key={tab}
                onClick={() => setFilter(tab)}
                className={`px-3 py-1 rounded text-xs font-semibold transition-colors ${
                  filter === tab
                    ? "bg-emerald-600 text-white"
                    : "text-slate-400 hover:text-slate-200"
                }`}
              >
                {tab.replace("_", " ")}
              </button>
            ))}
          </div>

          {onNavigateToLogComplaint && (
            <button
              onClick={onNavigateToLogComplaint}
              className="px-3 py-1.5 rounded-lg bg-emerald-600 hover:bg-emerald-500 text-white text-xs font-bold transition shadow"
            >
              + Log New Complaint
            </button>
          )}
        </div>
      </div>

      {/* Ticket Table */}
      <div className="overflow-x-auto min-h-[300px]">
        {filteredTickets.length === 0 ? (
          <div className="text-center py-16 text-slate-500 text-xs">
            No complaints matching query. Operational queue is currently clear.
          </div>
        ) : (
          <Table>
            <TableHead>
              <tr>
                <TableHeaderCell>Ticket #</TableHeaderCell>
                <TableHeaderCell>Customer &amp; Phone</TableHeaderCell>
                <TableHeaderCell>Category</TableHeaderCell>
                <TableHeaderCell>Assigned Unit</TableHeaderCell>
                <TableHeaderCell>Priority</TableHeaderCell>
                <TableHeaderCell>Status</TableHeaderCell>
                <TableHeaderCell>Action Gate</TableHeaderCell>
              </tr>
            </TableHead>
            <TableBody>
              {filteredTickets.map((t) => (
                <TableRow
                  key={t.id}
                  onClick={() => onSelectTicket(t)}
                  className="cursor-pointer hover:bg-slate-800/70 transition-colors"
                >
                  <TableCell className="font-mono font-bold text-emerald-400 text-xs">
                    {t.ticket_number}
                  </TableCell>
                  <TableCell>
                    <div className="text-xs font-semibold text-white">{t.customer_name || "Walk-In Consumer"}</div>
                    <div className="text-[11px] text-slate-400">{t.customer_phone || "No phone"}</div>
                  </TableCell>
                  <TableCell>
                    <span className="text-xs text-slate-300 font-medium">
                      {t.category.replace(/_/g, " ")}
                    </span>
                  </TableCell>
                  <TableCell>
                    <span className="text-xs text-slate-400 font-mono">
                      {t.unit_type || "TRIAGE_DESK"}
                    </span>
                  </TableCell>
                  <TableCell>
                    <Badge
                      variant={
                        t.priority === "CRITICAL"
                          ? "danger"
                          : t.priority === "HIGH"
                          ? "warning"
                          : "neutral"
                      }
                    >
                      {t.priority}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <Badge
                      variant={
                        t.status === "CLOSED"
                          ? "success"
                          : t.status === "ATTENDED_TO"
                          ? "warning"
                          : "info"
                      }
                    >
                      {t.status}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    {t.status === "ATTENDED_TO" ? (
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          setConfirmingTicket(t);
                        }}
                        className="px-2.5 py-1 rounded bg-amber-600 hover:bg-amber-500 text-white text-[11px] font-bold shadow transition animate-pulse"
                      >
                        📞 Call &amp; Confirm Closure
                      </button>
                    ) : t.status === "CLOSED" ? (
                      <span className="text-[11px] text-emerald-400 font-medium">
                        ✓ Verified Closed
                      </span>
                    ) : (
                      <span className="text-[11px] text-slate-500 font-mono">
                        Field Unit En Route
                      </span>
                    )}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </div>

      {/* 2-Step Confirmation Modal */}
      {confirmingTicket && (
        <div className="fixed inset-0 z-50 bg-black/80 flex items-center justify-center p-4">
          <div className="bg-slate-900 border border-slate-700 rounded-xl max-w-lg w-full p-6 space-y-4 shadow-2xl">
            <div className="flex items-center justify-between border-b border-slate-800 pb-3">
              <h3 className="text-sm font-bold text-white flex items-center gap-2">
                <span>🛡️ NERC 2-Step Ticket Closure Gate</span>
              </h3>
              <button
                onClick={() => setConfirmingTicket(null)}
                className="text-slate-400 hover:text-white text-xs font-bold"
              >
                ✕ Close
              </button>
            </div>

            <div className="bg-amber-500/10 border border-amber-500/30 rounded-lg p-3 text-xs text-amber-200 space-y-1">
              <p className="font-semibold">Mandatory Regulatory Protocol:</p>
              <p className="text-[11px] text-amber-300">
                You must place a phone call to the customer at{" "}
                <span className="font-bold underline">{confirmingTicket.customer_phone}</span> to confirm
                that power was restored or service issue was resolved to their satisfaction.
              </p>
            </div>

            <div className="space-y-2 text-xs">
              <div className="grid grid-cols-2 gap-2 text-slate-300">
                <div><span className="text-slate-500">Ticket:</span> {confirmingTicket.ticket_number}</div>
                <div><span className="text-slate-500">Customer:</span> {confirmingTicket.customer_name}</div>
                <div className="col-span-2">
                  <span className="text-slate-500">Field Unit Proof Note:</span>{" "}
                  <span className="text-emerald-400 font-mono">{confirmingTicket.attended_proof_note || "Repairs completed on site"}</span>
                </div>
              </div>

              <form onSubmit={handleExecuteClose} className="space-y-3 pt-2">
                <label className="block text-[11px] font-semibold text-slate-300">
                  Customer Verbal Feedback &amp; Closure Notes:
                </label>
                <textarea
                  required
                  rows={3}
                  value={feedbackNotes}
                  onChange={(e) => setFeedbackNotes(e.target.value)}
                  placeholder="e.g. Spoke with customer. They confirmed electricity was restored at 14:20 and meter display is functioning correctly."
                  className="w-full px-3 py-2 rounded bg-slate-950 border border-slate-700 text-xs text-white placeholder-slate-500 focus:outline-none focus:ring-1 focus:ring-emerald-500"
                />

                <div className="flex justify-end gap-2 pt-2">
                  <button
                    type="button"
                    onClick={() => setConfirmingTicket(null)}
                    className="px-3 py-1.5 rounded bg-slate-800 text-slate-300 text-xs font-semibold hover:bg-slate-700"
                  >
                    Cancel
                  </button>
                  <button
                    type="submit"
                    disabled={submittingClose}
                    className="px-4 py-1.5 rounded bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 text-white text-xs font-bold transition shadow"
                  >
                    {submittingClose ? "Verifying..." : "Confirm & Permanently Close Ticket"}
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

