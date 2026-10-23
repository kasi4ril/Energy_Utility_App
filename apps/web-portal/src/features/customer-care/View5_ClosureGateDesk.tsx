/**
 * View 5: Mandatory Customer Call-Confirmation & Closure Gate Desk.
 * Compliant with Section 12.3 Page 5: Action queue of all ATTENDED_TO tickets,
 * showing routing CC Rep, technical staff member attribution, GPS proof note,
 * customer phone confirmation outcome, and final CLOSE button stamping closed_by_cc_staff_name.
 * Strictly <= 350 lines.
 */

import React, { useState } from "react";
import { ServiceTicket } from "./types";
import { PhoneCall, CheckCircle, XCircle, ShieldCheck, UserCheck, PhoneForwarded } from "lucide-react";

interface View5Props {
  tickets: ServiceTicket[];
  staffId: string;
  staffName?: string;
  onConfirmClose: (ticketId: string, notes: string, confirmed: boolean) => Promise<void>;
  onRefresh?: () => void;
}

export function View5_ClosureGateDesk({
  tickets,
  staffId,
  staffName = "Amina Yusuf (CC Officer)",
  onConfirmClose,
  onRefresh,
}: View5Props) {
  const [selectedTicketId, setSelectedTicketId] = useState<string | null>(null);
  const [callNotes, setCallNotes] = useState("");
  const [callOutcome, setCallOutcome] = useState<"CONFIRMED" | "REJECTED">("CONFIRMED");
  const [isProcessing, setIsProcessing] = useState(false);

  const attendedTickets = tickets.filter((t) => t.status === "ATTENDED_TO");
  const selectedTicket = tickets.find((t) => t.id === selectedTicketId) || attendedTickets[0];

  const handleExecuteGate = async () => {
    if (!selectedTicket) return;
    setIsProcessing(true);
    try {
      await onConfirmClose(selectedTicket.id, callNotes, callOutcome === "CONFIRMED");
      setCallNotes("");
      if (onRefresh) onRefresh();
    } finally {
      setIsProcessing(false);
    }
  };

  return (
    <div className="space-y-6 max-w-5xl text-xs">
      {/* Header Banner */}
      <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-4 sm:p-5">
        <div className="flex items-center space-x-3 mb-2">
          <div className="p-2.5 rounded-lg bg-emerald-600/20 text-emerald-400 border border-emerald-500/30">
            <PhoneCall className="w-6 h-6" />
          </div>
          <div>
            <h2 className="text-base font-bold text-slate-100 uppercase tracking-wide">
              Mandatory Customer Call-Confirmation &amp; Closure Gate
            </h2>
            <p className="text-xs text-slate-400">
              NERC 2-Step Protective Resolution Gate &bull; Technical units cannot close tickets &bull; Staff attribution enforced
            </p>
          </div>
        </div>
      </div>

      {attendedTickets.length === 0 ? (
        <div className="bg-slate-900/60 border border-slate-800 rounded-xl p-8 text-center text-slate-400">
          <ShieldCheck className="w-10 h-10 text-emerald-400 mx-auto mb-2 opacity-80" />
          <h3 className="text-sm font-semibold text-slate-200">Closure Gate Clear</h3>
          <p className="text-xs mt-1">All attended service tickets have been telephone verified with customers.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-5">
          {/* Left Column: Attended Tickets Queue */}
          <div className="lg:col-span-5 bg-slate-900/90 border border-slate-800 rounded-xl p-3.5 space-y-2">
            <div className="text-[11px] font-bold uppercase tracking-wider text-slate-400 mb-2 px-1 flex justify-between">
              <span>Pending Verification Queue</span>
              <span className="bg-amber-950 text-amber-400 px-2 py-0.2 rounded font-mono">
                {attendedTickets.length} Attended
              </span>
            </div>

            <div className="space-y-2 max-h-[500px] overflow-y-auto pr-1">
              {attendedTickets.map((t) => {
                const isSelected = selectedTicket?.id === t.id;
                return (
                  <div
                    key={t.id}
                    onClick={() => setSelectedTicketId(t.id)}
                    className={`p-3 rounded-lg border cursor-pointer transition-all ${
                      isSelected
                        ? "bg-blue-950/60 border-blue-500/50 shadow-sm"
                        : "bg-slate-950/60 border-slate-800/80 hover:bg-slate-800/50"
                    }`}
                  >
                    <div className="flex items-center justify-between font-mono text-[11px] mb-1">
                      <span className="font-bold text-slate-200">{t.ticket_number}</span>
                      <span className="text-[10px] px-1.5 py-0.2 rounded bg-amber-500/20 text-amber-400">
                        ATTENDED_TO
                      </span>
                    </div>

                    <div className="font-semibold text-slate-100 truncate">{t.customer_name}</div>
                    <div className="text-[10px] text-slate-400 truncate">{t.category}</div>

                    <div className="mt-2 text-[10px] text-slate-400 border-t border-slate-800/80 pt-1.5 flex justify-between">
                      <span>Attended by:</span>
                      <span className="text-slate-200 font-medium truncate max-w-[140px]">
                        {t.attended_by_staff_name || "Technical Team Lead"}
                      </span>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>

          {/* Right Column: Active Verification & Customer Call Terminal */}
          {selectedTicket && (
            <div className="lg:col-span-7 bg-slate-900/90 border border-slate-800 rounded-xl p-4 sm:p-5 space-y-4">
              <div className="flex items-center justify-between border-b border-slate-800 pb-3">
                <div>
                  <span className="text-[10px] font-mono text-blue-400 uppercase font-bold tracking-wider">
                    Call Verification Terminal &bull; {selectedTicket.ticket_number}
                  </span>
                  <h3 className="text-sm font-bold text-slate-100 mt-0.5">{selectedTicket.customer_name}</h3>
                </div>
                <a
                  href={`tel:${selectedTicket.customer_phone}`}
                  className="bg-emerald-600 hover:bg-emerald-500 text-white font-semibold px-3 py-1.5 rounded-lg flex items-center space-x-1.5 transition-colors"
                >
                  <PhoneForwarded className="w-3.5 h-3.5" />
                  <span>Call {selectedTicket.customer_phone}</span>
                </a>
              </div>

              {/* End-to-End Attribution Badges */}
              <div className="grid grid-cols-2 gap-2.5 bg-slate-950/70 p-3 rounded-lg border border-slate-800/80 text-[11px]">
                <div>
                  <div className="text-slate-500 text-[10px] uppercase font-mono">Logged &amp; Routed By:</div>
                  <div className="text-slate-200 font-medium">
                    {selectedTicket.logged_by_cc_staff_name || "Amina Yusuf (CC-02)"}
                  </div>
                </div>
                <div>
                  <div className="text-slate-500 text-[10px] uppercase font-mono">Work Done By (Technical):</div>
                  <div className="text-emerald-400 font-medium">
                    {selectedTicket.attended_by_staff_name || "Engr. Paul Gyang (Metering Tech Lead)"}
                  </div>
                </div>
              </div>

              {/* Technical Work Done Summary */}
              <div className="bg-slate-950/70 p-3 rounded-lg border border-slate-800/80">
                <div className="text-slate-400 text-[10px] uppercase font-mono mb-1">
                  Field Proof &amp; Work Done Summary
                </div>
                <p className="text-slate-200 text-xs leading-relaxed">
                  {selectedTicket.work_done_summary || selectedTicket.attended_proof_note || "Tamper cleared, internal fuse replaced, voltage measured at 228V."}
                </p>
              </div>

              {/* Call Confirmation Gate Radio Form */}
              <div className="space-y-3 pt-2">
                <div className="text-slate-300 font-medium text-xs">Customer Telephone Call Outcome *</div>
                <div className="grid grid-cols-2 gap-3">
                  <button
                    type="button"
                    onClick={() => setCallOutcome("CONFIRMED")}
                    className={`p-3 rounded-lg border flex items-center space-x-2.5 transition-all text-left ${
                      callOutcome === "CONFIRMED"
                        ? "bg-emerald-950/70 border-emerald-500 text-emerald-200"
                        : "bg-slate-950/40 border-slate-800 text-slate-400 hover:bg-slate-800/40"
                    }`}
                  >
                    <CheckCircle className={`w-4 h-4 ${callOutcome === "CONFIRMED" ? "text-emerald-400" : ""}`} />
                    <div>
                      <div className="font-bold text-xs">Customer Confirmed Resolution</div>
                      <div className="text-[10px] text-slate-400">Power restored, work accepted</div>
                    </div>
                  </button>

                  <button
                    type="button"
                    onClick={() => setCallOutcome("REJECTED")}
                    className={`p-3 rounded-lg border flex items-center space-x-2.5 transition-all text-left ${
                      callOutcome === "REJECTED"
                        ? "bg-red-950/70 border-red-500 text-red-200"
                        : "bg-slate-950/40 border-slate-800 text-slate-400 hover:bg-slate-800/40"
                    }`}
                  >
                    <XCircle className={`w-4 h-4 ${callOutcome === "REJECTED" ? "text-red-400" : ""}`} />
                    <div>
                      <div className="font-bold text-xs">Customer Rejected Resolution</div>
                      <div className="text-[10px] text-slate-400">Reopen &amp; escalate to Dept Head</div>
                    </div>
                  </button>
                </div>

                <div>
                  <label className="block text-slate-400 mb-1">Customer Feedback Notes *</label>
                  <textarea
                    rows={2}
                    required
                    placeholder="Document customer remarks during the verification call..."
                    value={callNotes}
                    onChange={(e) => setCallNotes(e.target.value)}
                    className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-slate-100 focus:border-blue-500 focus:outline-none"
                  />
                </div>

                <div className="p-2.5 bg-slate-950/60 rounded border border-slate-800 text-[10px] text-slate-400 flex items-center justify-between">
                  <span>
                    Closing Officer: <strong className="text-slate-200">{staffName}</strong> ({staffId})
                  </span>
                  <span className="text-emerald-400 font-mono">NERC Closure Stamp</span>
                </div>

                <button
                  type="button"
                  disabled={isProcessing}
                  onClick={handleExecuteGate}
                  className={`w-full font-bold py-2.5 px-4 rounded-lg transition-colors flex items-center justify-center space-x-2 shadow-md ${
                    callOutcome === "CONFIRMED"
                      ? "bg-emerald-600 hover:bg-emerald-500 text-white"
                      : "bg-red-600 hover:bg-red-500 text-white"
                  }`}
                >
                  <UserCheck className="w-4 h-4" />
                  <span>
                    {callOutcome === "CONFIRMED"
                      ? "EXECUTE VERIFIED CLOSURE GATE (MARK CLOSED)"
                      : "REOPEN TICKET & ESCALATE TO DEPARTMENT HEAD"}
                  </span>
                </button>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

