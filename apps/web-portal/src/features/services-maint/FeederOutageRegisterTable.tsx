import React, { useState } from "react";
import { OutageStatus, FeederClosureOutcome, FeederOutageRequest } from "./types";
import { FeederRecordReturnModal } from "./FeederRecordReturnModal";

interface OutageRegisterTableProps {
  requests: FeederOutageRequest[];
  onRecordReturn?: (reqId: string, returnTime: string, actionText: string, outcome: FeederClosureOutcome) => void;
}

const MONTHS = [
  { val: "01", name: "Jan - January" },
  { val: "02", name: "Feb - February" },
  { val: "03", name: "Mar - March" },
  { val: "04", name: "Apr - April" },
  { val: "05", name: "May - May" },
  { val: "06", name: "Jun - June" },
  { val: "07", name: "Jul - July" },
  { val: "08", name: "Aug - August" },
  { val: "09", name: "Sep - September" },
  { val: "10", name: "Oct - October" },
  { val: "11", name: "Nov - November" },
  { val: "12", name: "Dec - December" },
];

export const FeederOutageRegisterTable: React.FC<OutageRegisterTableProps> = ({ requests, onRecordReturn }) => {
  const [selectedYear, setSelectedYear] = useState<string>("2026");
  const [selectedMonth, setSelectedMonth] = useState<string>("ALL");
  const [selectedDay, setSelectedDay] = useState<string>("ALL");
  const [statusFilter, setStatusFilter] = useState<"ALL" | "APPROVED" | "PENDING" | "DECLINED">("ALL");
  const [activeReturnReq, setActiveReturnReq] = useState<FeederOutageRequest | null>(null);

  const isApproved = (s: OutageStatus) => s === "APPROVED_PTW_ISSUED" || s === "ISOLATED_IN_PROGRESS" || s === "COMPLETED_LINE_ENERGIZED";
  const isDeclined = (s: OutageStatus) => s === "REJECTED_BY_DISPATCH";
  const isPending = (s: OutageStatus) => s === "PENDING_DISPATCH";

  const approvedCount = requests.filter((r) => isApproved(r.status)).length;
  const pendingCount = requests.filter((r) => isPending(r.status)).length;
  const declinedCount = requests.filter((r) => isDeclined(r.status)).length;

  const maxDays = selectedYear !== "ALL" && selectedMonth !== "ALL"
    ? new Date(parseInt(selectedYear, 10), parseInt(selectedMonth, 10), 0).getDate()
    : 31;
  const daysList = Array.from({ length: maxDays }, (_, i) => String(i + 1).padStart(2, "0"));

  const handleDatePick = (dateStr: string) => {
    if (!dateStr) { setSelectedYear("ALL"); setSelectedMonth("ALL"); setSelectedDay("ALL"); return; }
    const [y, m, d] = dateStr.split("-");
    setSelectedYear(y); setSelectedMonth(m); setSelectedDay(d);
  };

  const filteredRequests = requests.filter((r) => {
    const [rYear, rMonth, rDay] = r.requested_date.split("-");
    if (selectedYear !== "ALL" && rYear !== selectedYear) return false;
    if (selectedMonth !== "ALL" && rMonth !== selectedMonth) return false;
    if (selectedDay !== "ALL" && rDay !== selectedDay) return false;
    if (statusFilter === "APPROVED" && !isApproved(r.status)) return false;
    if (statusFilter === "PENDING" && !isPending(r.status)) return false;
    if (statusFilter === "DECLINED" && !isDeclined(r.status)) return false;
    return true;
  });

  const getStatusBadge = (status: OutageStatus) => {
    switch (status) {
      case "APPROVED_PTW_ISSUED":
        return <span className="px-2 py-0.5 rounded text-[10px] font-bold bg-emerald-900/90 text-emerald-200 border border-emerald-500/50">APPROVED (PTW)</span>;
      case "COMPLETED_LINE_ENERGIZED":
        return <span className="px-2 py-0.5 rounded text-[10px] font-bold bg-blue-900/90 text-blue-200 border border-blue-500/50">APPROVED (ENERGIZED)</span>;
      case "ISOLATED_IN_PROGRESS":
        return <span className="px-2 py-0.5 rounded text-[10px] font-bold bg-sky-900/90 text-sky-200 border border-sky-500/50">APPROVED (ISOLATED)</span>;
      case "PENDING_DISPATCH":
        return <span className="px-2 py-0.5 rounded text-[10px] font-bold bg-amber-900/90 text-amber-200 border border-amber-500/50">PENDING REQUEST</span>;
      default:
        return <span className="px-2 py-0.5 rounded text-[10px] font-bold bg-rose-900/90 text-rose-200 border border-rose-500/50">DECLINED</span>;
    }
  };

  const getClosureBadge = (outcome?: FeederClosureOutcome) => {
    switch (outcome) {
      case "CLOSED_ENERGIZED_IMMEDIATELY":
        return <span className="inline-block px-1.5 py-0.5 rounded text-[9px] font-bold bg-emerald-950 text-emerald-300 border border-emerald-600/50 mt-0.5">✓ CLOSED IMMEDIATELY</span>;
      case "LEFT_ON_LOAD_SHEDDING":
        return <span className="inline-block px-1.5 py-0.5 rounded text-[9px] font-bold bg-amber-950 text-amber-300 border border-amber-600/50 mt-0.5">⚠ LEFT ON LOAD SHED</span>;
      case "TRIPPED_ON_CLOSURE":
        return <span className="inline-block px-1.5 py-0.5 rounded text-[9px] font-bold bg-rose-950 text-rose-300 border border-rose-600/50 mt-0.5">⚡ TRIPPED ON CLOSURE</span>;
      case "PENDING_SURRENDER":
        return <span className="inline-block px-1.5 py-0.5 rounded text-[9px] font-bold bg-slate-800 text-slate-300 border border-slate-700/50 mt-0.5">⏳ PENDING SURRENDER</span>;
      default:
        return null;
    }
  };

  const hasCalendarFilter = selectedYear !== "ALL" || selectedMonth !== "ALL" || selectedDay !== "ALL";

  return (
    <div className="bg-slate-900/70 border border-slate-800 rounded-xl p-5 shadow-sm space-y-3.5">
      <div className="flex items-center justify-between border-b border-slate-800 pb-2.5">
        <div>
          <h3 className="text-base font-semibold text-white flex items-center gap-2">
            <svg className="w-4 h-4 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
            Outage Request and Return Register
          </h3>
          <p className="text-xs text-slate-400">Permit granting, permit return upon exhaustion, and post-closure action log.</p>
        </div>
        <span className="text-xs font-semibold px-2 py-0.5 rounded bg-slate-800 text-slate-300 font-mono">
          {filteredRequests.length} of {requests.length} records
        </span>
      </div>

      {/* Calendar Date Query Toolbar */}
      <div className="bg-slate-950/80 border border-slate-800 rounded-lg p-3 space-y-2.5">
        <div className="text-[11px] font-bold uppercase tracking-wider text-amber-400/90 flex items-center justify-between">
          <span>Calendar Date Query (Year • Month • Day)</span>
          <span className="text-slate-400 lowercase font-normal">All past &amp; scheduled outages</span>
        </div>

        <div className="grid grid-cols-2 sm:grid-cols-12 gap-2 text-xs">
          <div className="col-span-1 sm:col-span-3">
            <label className="block text-[10px] text-slate-400 font-semibold mb-0.5">Year</label>
            <select
              value={selectedYear}
              onChange={(e) => setSelectedYear(e.target.value)}
              style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
              className="w-full px-2 py-1.5 rounded border border-slate-700 text-xs font-semibold focus:outline-none"
            >
              <option value="ALL" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>All Years</option>
              <option value="2026" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>2026</option>
              <option value="2025" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>2025</option>
              <option value="2024" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>2024</option>
            </select>
          </div>

          <div className="col-span-1 sm:col-span-3">
            <label className="block text-[10px] text-slate-400 font-semibold mb-0.5">Month</label>
            <select
              value={selectedMonth}
              onChange={(e) => setSelectedMonth(e.target.value)}
              style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
              className="w-full px-2 py-1.5 rounded border border-slate-700 text-xs font-semibold focus:outline-none"
            >
              <option value="ALL" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>All Months</option>
              {MONTHS.map((m) => (
                <option key={m.val} value={m.val} style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>{m.name}</option>
              ))}
            </select>
          </div>

          <div className="col-span-1 sm:col-span-3">
            <label className="block text-[10px] text-slate-400 font-semibold mb-0.5">Day</label>
            <select
              value={selectedDay}
              onChange={(e) => setSelectedDay(e.target.value)}
              style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
              className="w-full px-2 py-1.5 rounded border border-slate-700 text-xs font-semibold focus:outline-none"
            >
              <option value="ALL" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>All Days</option>
              {daysList.map((d) => (
                <option key={d} value={d} style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Day {parseInt(d, 10)} ({d})</option>
              ))}
            </select>
          </div>

          <div className="col-span-1 sm:col-span-3">
            <label className="block text-[10px] text-slate-400 font-semibold mb-0.5">Approval Status</label>
            <select
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value as any)}
              style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
              className="w-full px-2 py-1.5 rounded border border-slate-700 text-xs font-semibold focus:outline-none"
            >
              <option value="ALL" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>All Statuses</option>
              <option value="APPROVED" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Approved ({approvedCount})</option>
              <option value="PENDING" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Pending ({pendingCount})</option>
              <option value="DECLINED" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Declined ({declinedCount})</option>
            </select>
          </div>
        </div>

        {/* Shortcuts Bar */}
        <div className="flex flex-wrap items-center justify-between gap-2 pt-1 border-t border-slate-800/80 text-[11px]">
          <div className="flex items-center gap-1.5">
            <span className="text-slate-400">Shortcuts:</span>
            <button
              type="button"
              onClick={() => { setSelectedYear("2026"); setSelectedMonth("10"); setSelectedDay("24"); }}
              className="px-2 py-0.5 bg-slate-800 hover:bg-slate-700 text-amber-300 rounded font-mono"
            >
              Today (24 Oct)
            </button>
            <button
              type="button"
              onClick={() => { setSelectedYear("2026"); setSelectedMonth("10"); setSelectedDay("ALL"); }}
              className="px-2 py-0.5 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded font-mono"
            >
              Oct 2026
            </button>
            <button
              type="button"
              onClick={() => { setSelectedYear("2026"); setSelectedMonth("09"); setSelectedDay("ALL"); }}
              className="px-2 py-0.5 bg-slate-800 hover:bg-slate-700 text-slate-300 rounded font-mono"
            >
              Sep 2026
            </button>
          </div>

          <div className="flex items-center gap-2">
            <input
              type="date"
              title="Pick exact date"
              onChange={(e) => handleDatePick(e.target.value)}
              style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
              className="px-1.5 py-0.5 rounded border border-slate-700 text-[11px] focus:outline-none"
            />
            {hasCalendarFilter && (
              <button
                type="button"
                onClick={() => { setSelectedYear("ALL"); setSelectedMonth("ALL"); setSelectedDay("ALL"); setStatusFilter("ALL"); }}
                className="text-amber-400 hover:text-amber-300 font-semibold underline"
              >
                Clear Calendar
              </button>
            )}
          </div>
        </div>
      </div>

      {/* Table */}
      <div className="overflow-x-auto">
        <table className="w-full text-left text-xs border-collapse">
          <thead>
            <tr className="border-b border-slate-800 text-slate-400 uppercase tracking-wider font-semibold">
              <th className="py-2 px-2.5">Req # / Corridor</th>
              <th className="py-2 px-2.5">Category &amp; Scope</th>
              <th className="py-2 px-2.5">Grant &amp; Return Times (WAT)</th>
              <th className="py-2 px-2.5">Action Taken on Return</th>
              <th className="py-2 px-2.5">Status &amp; Closure</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-800/60 text-slate-300">
            {filteredRequests.length > 0 ? (
              filteredRequests.map((r) => {
                const isReturned = !!r.time_returned_wat || !!r.surrendering_time_wat;
                const canRecordReturn = !isReturned && (r.status === "APPROVED_PTW_ISSUED" || r.status === "ISOLATED_IN_PROGRESS");
                return (
                  <tr key={r.id} className="hover:bg-slate-800/40 transition-colors">
                    <td className="py-2.5 px-2.5">
                      <div className="font-semibold text-white">{r.request_no}</div>
                      <div className="text-[11px] text-amber-400 font-mono">{r.corridor.replace("FEEDER_", "")}</div>
                      {r.ptw_issued_no && <div className="text-[10px] text-emerald-400 font-mono">PTW: {r.ptw_issued_no}</div>}
                    </td>
                    <td className="py-2.5 px-2.5">
                      <div className="text-slate-200">{r.outage_type.replace(/_/g, " ")}</div>
                      <span className={`inline-block text-[10px] px-1.5 py-0.2 rounded font-bold mt-0.5 ${r.urgency === "EMERGENCY" ? "bg-rose-950 text-rose-300" : r.urgency === "URGENT" ? "bg-amber-950 text-amber-300" : "bg-slate-800 text-slate-300"}`}>
                        {r.urgency}
                      </span>
                      <div className="text-[11px] text-slate-400 mt-0.5">{r.estimated_load_mw} MW • {r.affected_dts_count} DTs</div>
                    </td>
                    <td className="py-2.5 px-2.5 whitespace-nowrap">
                      <div className="font-mono text-slate-200">{r.requested_date} ({r.duration_hours}h)</div>
                      <div className="text-[11px] text-emerald-400 font-mono mt-0.5">
                        Granted: {r.opening_time_wat || "Not yet opened"}
                      </div>
                      <div className="text-[11px] text-amber-300 font-mono">
                        Returned: {r.time_returned_wat || r.surrendering_time_wat || (r.status === "REJECTED_BY_DISPATCH" ? "N/A (Declined)" : "Permit In Progress")}
                      </div>
                    </td>
                    <td className="py-2.5 px-2.5 max-w-[200px]">
                      {r.action_taken_on_return ? (
                        <div className="text-[11px] text-slate-200 bg-slate-950/70 p-1.5 rounded border border-slate-800/80 leading-snug">
                          {r.action_taken_on_return}
                        </div>
                      ) : (
                        <span className="text-[11px] text-slate-500 italic">
                          {r.status === "REJECTED_BY_DISPATCH" ? "Declined by dispatch" : "Pending return & exhaustion"}
                        </span>
                      )}
                    </td>
                    <td className="py-2.5 px-2.5">
                      <div>{getStatusBadge(r.status)}</div>
                      <div>{getClosureBadge(r.closure_outcome)}</div>
                      {canRecordReturn && (
                        <button
                          type="button"
                          onClick={() => setActiveReturnReq(r)}
                          className="mt-1.5 px-2 py-0.5 bg-amber-500 hover:bg-amber-600 text-slate-950 font-bold rounded text-[10px] transition-all block"
                        >
                          Record Return
                        </button>
                      )}
                    </td>
                  </tr>
                );
              })
            ) : (
              <tr>
                <td colSpan={5} className="py-8 text-center text-slate-500">
                  No records found in return register for selected date and filter.
                  <button
                    onClick={() => { setSelectedYear("ALL"); setSelectedMonth("ALL"); setSelectedDay("ALL"); setStatusFilter("ALL"); }}
                    className="block mx-auto mt-2 text-xs text-amber-400 hover:underline"
                  >
                    Reset calendar to view all historical and scheduled requests
                  </button>
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {activeReturnReq && onRecordReturn && (
        <FeederRecordReturnModal
          request={activeReturnReq}
          onClose={() => setActiveReturnReq(null)}
          onSubmit={(reqId, retTime, actText, outc) => {
            onRecordReturn(reqId, retTime, actText, outc);
            setActiveReturnReq(null);
          }}
        />
      )}
    </div>
  );
};
export default FeederOutageRegisterTable;
