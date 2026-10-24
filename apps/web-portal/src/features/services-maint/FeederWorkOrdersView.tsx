import React, { useState } from "react";
import { FeederWorkOrder } from "./types";
import { MOCK_FEEDER_WORK_ORDERS } from "./mockData";

export const FeederWorkOrdersView: React.FC = () => {
  const [orders, setOrders] = useState<FeederWorkOrder[]>(MOCK_FEEDER_WORK_ORDERS);
  const [statusFilter, setStatusFilter] = useState<string>("ALL");
  const [priorityFilter, setPriorityFilter] = useState<string>("ALL");
  const [selectedOrder, setSelectedOrder] = useState<FeederWorkOrder | null>(
    MOCK_FEEDER_WORK_ORDERS[0] || null
  );

  const filteredOrders = orders.filter((o) => {
    if (statusFilter !== "ALL" && o.status !== statusFilter) return false;
    if (priorityFilter !== "ALL" && o.priority !== priorityFilter) return false;
    return true;
  });

  const updateOrderStatus = (orderId: string, newStatus: FeederWorkOrder["status"]) => {
    const updated = orders.map((o) => (o.id === orderId ? { ...o, status: newStatus } : o));
    setOrders(updated);
    if (selectedOrder && selectedOrder.id === orderId) {
      setSelectedOrder({ ...selectedOrder, status: newStatus });
    }
  };

  const getPriorityBadge = (priority: FeederWorkOrder["priority"]) => {
    switch (priority) {
      case "CRITICAL":
        return <span className="px-2 py-0.5 rounded text-[10px] font-bold bg-rose-950 text-rose-300 border border-rose-800/50">CRITICAL</span>;
      case "HIGH":
        return <span className="px-2 py-0.5 rounded text-[10px] font-bold bg-amber-950 text-amber-300 border border-amber-800/50">HIGH</span>;
      case "MEDIUM":
        return <span className="px-2 py-0.5 rounded text-[10px] font-bold bg-sky-950 text-sky-300 border border-sky-800/50">MEDIUM</span>;
      default:
        return <span className="px-2 py-0.5 rounded text-[10px] font-bold bg-slate-800 text-slate-300">LOW</span>;
    }
  };

  const getStatusBadge = (status: FeederWorkOrder["status"]) => {
    switch (status) {
      case "ASSIGNED":
        return <span className="px-2 py-0.5 rounded text-[10px] font-bold bg-slate-800 text-slate-300">ASSIGNED</span>;
      case "IN_PROGRESS":
        return <span className="px-2 py-0.5 rounded text-[10px] font-bold bg-amber-900/80 text-amber-200 border border-amber-700/40">IN PROGRESS</span>;
      case "AWAITING_PROOF":
        return <span className="px-2 py-0.5 rounded text-[10px] font-bold bg-purple-900/80 text-purple-200 border border-purple-700/40">AWAITING PROOF</span>;
      case "ATTENDED_TO":
        return <span className="px-2 py-0.5 rounded text-[10px] font-bold bg-emerald-900/80 text-emerald-200 border border-emerald-700/40">ATTENDED TO</span>;
    }
  };

  return (
    <div className="space-y-6">
      {/* Top Banner */}
      <div className="bg-slate-900/80 border border-slate-800 rounded-xl p-5 shadow-lg flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <div className="inline-flex items-center gap-2 px-2.5 py-1 rounded-full text-xs font-semibold bg-amber-500/10 text-amber-300 border border-amber-500/30 mb-2">
            <span className="w-2 h-2 rounded-full bg-amber-400" />
            Field Engineering Dispatch • Feeder Unit
          </div>
          <h2 className="text-xl font-bold text-white tracking-tight">
            Active Feeder Work Orders & Crew Operations Queue
          </h2>
          <p className="text-sm text-slate-400 mt-1">
            Dispatch, PTW permits tracking, fault rectifications, and crew lifecycle supervision across transmission & distribution lines.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <div className="bg-slate-950/70 border border-slate-800 rounded-lg px-3 py-2 text-center">
            <div className="text-[10px] text-slate-400 uppercase font-semibold">In Progress</div>
            <div className="text-xl font-bold text-amber-400">
              {orders.filter((o) => o.status === "IN_PROGRESS").length}
            </div>
          </div>
          <div className="bg-slate-950/70 border border-slate-800 rounded-lg px-3 py-2 text-center">
            <div className="text-[10px] text-slate-400 uppercase font-semibold">Attended</div>
            <div className="text-xl font-bold text-emerald-400">
              {orders.filter((o) => o.status === "ATTENDED_TO").length}
            </div>
          </div>
        </div>
      </div>

      {/* Filter Bar with high-contrast dropdowns */}
      <div className="bg-slate-900/70 border border-slate-800 rounded-xl p-3.5 flex flex-wrap items-center justify-between gap-3">
        <div className="flex flex-wrap items-center gap-3 text-xs">
          <div>
            <label className="text-slate-400 font-semibold mr-2">Filter Status:</label>
            <select
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
              style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
              className="px-2.5 py-1.5 rounded-lg border border-slate-700 text-xs focus:outline-none"
            >
              <option value="ALL" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>All Statuses ({orders.length})</option>
              <option value="ASSIGNED" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Assigned</option>
              <option value="IN_PROGRESS" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>In Progress</option>
              <option value="AWAITING_PROOF" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Awaiting Proof</option>
              <option value="ATTENDED_TO" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Attended To</option>
            </select>
          </div>
          <div>
            <label className="text-slate-400 font-semibold mr-2">Priority:</label>
            <select
              value={priorityFilter}
              onChange={(e) => setPriorityFilter(e.target.value)}
              style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
              className="px-2.5 py-1.5 rounded-lg border border-slate-700 text-xs focus:outline-none"
            >
              <option value="ALL" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>All Priorities</option>
              <option value="CRITICAL" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Critical</option>
              <option value="HIGH" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>High</option>
              <option value="MEDIUM" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Medium</option>
              <option value="LOW" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>Low</option>
            </select>
          </div>
        </div>
        <div className="text-xs text-slate-400 font-mono">
          Showing {filteredOrders.length} of {orders.length} orders
        </div>
      </div>

      {/* Grid: Order List (5 cols) & Order Detail Card (7 cols) */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        {/* Orders List */}
        <div className="lg:col-span-5 bg-slate-900/70 border border-slate-800 rounded-xl p-4 shadow-sm space-y-2.5">
          <div className="text-xs font-semibold text-slate-400 uppercase tracking-wider pb-2 border-b border-slate-800">
            Work Order Queue
          </div>
          <div className="space-y-2 max-h-[600px] overflow-y-auto pr-1">
            {filteredOrders.map((o) => {
              const isSelected = selectedOrder?.id === o.id;
              return (
                <div
                  key={o.id}
                  onClick={() => setSelectedOrder(o)}
                  className={`p-3 rounded-lg border cursor-pointer transition-all ${
                    isSelected
                      ? "bg-amber-950/20 border-amber-500/60 shadow-md"
                      : "bg-slate-950/60 border-slate-800/80 hover:bg-slate-800/40"
                  }`}
                >
                  <div className="flex items-center justify-between">
                    <span className="font-mono font-bold text-xs text-white">{o.order_no}</span>
                    <div className="flex items-center gap-1.5">
                      {getPriorityBadge(o.priority)}
                      {getStatusBadge(o.status)}
                    </div>
                  </div>
                  <div className="text-xs font-semibold text-slate-200 mt-1 line-clamp-1">{o.title}</div>
                  <div className="text-[11px] text-amber-400 font-mono mt-0.5">Asset: {o.asset_id} • {o.corridor.replace("FEEDER_", "")}</div>
                  <div className="text-[11px] text-slate-400 flex items-center justify-between mt-1">
                    <span>{o.assigned_team_lead}</span>
                    {o.ptw_permit_no && (
                      <span className="text-emerald-400 font-mono text-[10px]">PTW: {o.ptw_permit_no}</span>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        </div>

        {/* Order Details & Operational Control */}
        <div className="lg:col-span-7 bg-slate-900/70 border border-slate-800 rounded-xl p-5 shadow-sm space-y-4">
          {selectedOrder ? (
            <>
              <div className="flex flex-col md:flex-row md:items-center justify-between gap-3 border-b border-slate-800 pb-4">
                <div>
                  <div className="flex items-center gap-2">
                    <span className="font-mono text-sm font-bold text-amber-400">{selectedOrder.order_no}</span>
                    {selectedOrder.ticket_ref && (
                      <span className="text-[11px] text-slate-400 font-mono">CC Ref: {selectedOrder.ticket_ref}</span>
                    )}
                  </div>
                  <h3 className="text-lg font-bold text-white mt-1">{selectedOrder.title}</h3>
                  <div className="text-xs text-slate-400 font-mono mt-0.5">Corridor: {selectedOrder.corridor}</div>
                </div>
                <div className="flex items-center gap-2">
                  {getPriorityBadge(selectedOrder.priority)}
                  {getStatusBadge(selectedOrder.status)}
                </div>
              </div>

              {/* Asset & Team Information */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-xs">
                <div className="bg-slate-950/70 border border-slate-800 rounded-lg p-3 space-y-2">
                  <div className="text-slate-400 font-semibold uppercase tracking-wider text-[10px]">Asset & Fault Type</div>
                  <div className="text-white font-bold text-sm">{selectedOrder.asset_id}</div>
                  <div className="text-amber-400 font-mono text-xs">{selectedOrder.fault_type}</div>
                  {selectedOrder.ptw_required && (
                    <div className="mt-2 pt-2 border-t border-slate-800/80 text-emerald-400 flex items-center gap-1.5 font-semibold">
                      <span>✓ Safety PTW Mandated</span>
                      {selectedOrder.ptw_permit_no && (
                        <span className="font-mono text-xs">({selectedOrder.ptw_permit_no})</span>
                      )}
                    </div>
                  )}
                </div>

                <div className="bg-slate-950/70 border border-slate-800 rounded-lg p-3 space-y-2">
                  <div className="text-slate-400 font-semibold uppercase tracking-wider text-[10px]">Crew Assignment & Time</div>
                  <div className="text-white font-bold text-sm">{selectedOrder.assigned_team_lead}</div>
                  <div className="text-slate-400 text-xs">Dispatched: {new Date(selectedOrder.created_at).toLocaleString()}</div>
                </div>
              </div>

              {/* Instructions Box */}
              <div className="bg-slate-950/70 border border-slate-800 rounded-lg p-3.5 space-y-1.5">
                <div className="text-xs font-semibold text-slate-300">Field Maintenance Instructions</div>
                <p className="text-xs text-slate-300 leading-relaxed font-mono bg-slate-900 p-2.5 rounded border border-slate-800">
                  {selectedOrder.instructions}
                </p>
              </div>

              {/* Status Transition Action Buttons */}
              <div className="bg-slate-950/90 border border-slate-800 rounded-lg p-4 space-y-2.5">
                <div className="text-xs font-semibold text-slate-300">Update Work Order Status</div>
                <div className="flex flex-wrap items-center gap-2.5 pt-1">
                  <button
                    onClick={() => updateOrderStatus(selectedOrder.id, "IN_PROGRESS")}
                    disabled={selectedOrder.status === "IN_PROGRESS"}
                    className="px-3 py-1.5 bg-amber-600 hover:bg-amber-700 disabled:opacity-40 text-white font-bold text-xs rounded transition-all"
                  >
                    Set IN PROGRESS
                  </button>
                  <button
                    onClick={() => updateOrderStatus(selectedOrder.id, "AWAITING_PROOF")}
                    disabled={selectedOrder.status === "AWAITING_PROOF"}
                    className="px-3 py-1.5 bg-purple-600 hover:bg-purple-700 disabled:opacity-40 text-white font-bold text-xs rounded transition-all"
                  >
                    Awaiting GPS / Photo Proof
                  </button>
                  <button
                    onClick={() => updateOrderStatus(selectedOrder.id, "ATTENDED_TO")}
                    disabled={selectedOrder.status === "ATTENDED_TO"}
                    className="px-3 py-1.5 bg-emerald-600 hover:bg-emerald-700 disabled:opacity-40 text-white font-bold text-xs rounded transition-all"
                  >
                    Mark ATTENDED TO (Close Out)
                  </button>
                </div>
              </div>
            </>
          ) : (
            <div className="text-center py-16 text-slate-500 text-sm">Select a work order from the left.</div>
          )}
        </div>
      </div>
    </div>
  );
};
export default FeederWorkOrdersView;

