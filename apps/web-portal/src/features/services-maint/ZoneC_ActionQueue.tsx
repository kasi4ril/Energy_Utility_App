/**
 * Zone C: Main Operational Workspace (Left Pane - 65% Action Queue - Services & Maintenance).
 * Live queue of dispatched work orders, PTW isolations, and proof-of-work verifications.
 */

import React, { useState } from "react";
import { Badge } from "../../components/ui/Badge";
import { Table, TableHead, TableBody, TableRow, TableHeaderCell, TableCell } from "../../components/ui/Table";
import { WorkOrderDTO } from "./api";

interface ActionQueueProps {
  orders: WorkOrderDTO[];
  onSelectOrder: (order: WorkOrderDTO) => void;
  onVerifyProof?: (orderId: string) => void;
}

export const ZoneC_ActionQueue: React.FC<ActionQueueProps> = ({
  orders,
  onSelectOrder,
  onVerifyProof,
}) => {
  const [filter, setFilter] = useState("ALL");

  const filtered = orders.filter((o) => {
    if (filter === "PENDING_PROOF") return o.status === "PROOF_SUBMITTED";
    if (filter === "IN_PROGRESS") return o.status === "ASSIGNED" || o.status === "IN_PROGRESS";
    return true;
  });

  return (
    <div className="bg-slate-800/90 border border-slate-700/80 rounded-xl p-5 shadow-lg flex flex-col h-full">
      <div className="flex items-center justify-between pb-4 border-b border-slate-700/60">
        <div>
          <h3 className="text-sm font-bold text-white tracking-wide">Work Orders Queue</h3>
          <p className="text-xs text-slate-400 mt-0.5">Click to view immutable WorkNotes hash-chain</p>
        </div>
        <div className="flex gap-1.5">
          {["ALL", "IN_PROGRESS", "PENDING_PROOF"].map((tab) => (
            <button
              key={tab}
              onClick={() => setFilter(tab)}
              className={`px-3 py-1 rounded-lg text-xs font-semibold ${
                filter === tab ? "bg-emerald-600 text-white" : "bg-slate-700/60 text-slate-300"
              }`}
            >
              {tab.replace("_", " ")}
            </button>
          ))}
        </div>
      </div>

      <div className="mt-4 flex-1 overflow-y-auto">
        {filtered.length === 0 ? (
          <div className="text-center py-12 text-xs text-slate-500">
            No work orders in current queue. All units clear.
          </div>
        ) : (
          <Table>
            <TableHead>
              <tr>
                <TableHeaderCell>Order #</TableHeaderCell>
                <TableHeaderCell>Asset / Feeder</TableHeaderCell>
                <TableHeaderCell>Team Lead</TableHeaderCell>
                <TableHeaderCell>PTW</TableHeaderCell>
                <TableHeaderCell>Priority</TableHeaderCell>
                <TableHeaderCell>Status</TableHeaderCell>
                <TableHeaderCell>Action</TableHeaderCell>
              </tr>
            </TableHead>
            <TableBody>
              {filtered.map((o) => (
                <TableRow key={o.id} onClick={() => onSelectOrder(o)}>
                  <TableCell className="font-mono text-emerald-400 font-medium">{o.order_no}</TableCell>
                  <TableCell>{o.asset_id}</TableCell>
                  <TableCell>{o.team_lead_id || "UNASSIGNED"}</TableCell>
                  <TableCell>
                    {o.ptw_required ? (
                      <span className="text-[10px] text-amber-400 font-bold">REQUIRED</span>
                    ) : (
                      <span className="text-[10px] text-slate-500">NO</span>
                    )}
                  </TableCell>
                  <TableCell>
                    <Badge variant={o.priority === "CRITICAL" ? "danger" : "warning"}>{o.priority}</Badge>
                  </TableCell>
                  <TableCell>
                    <Badge variant={o.status === "ATTENDED_TO" ? "success" : "info"}>{o.status}</Badge>
                  </TableCell>
                  <TableCell>
                    {o.status === "PROOF_SUBMITTED" && onVerifyProof && (
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          onVerifyProof(o.id);
                        }}
                        className="px-2 py-1 bg-emerald-600 hover:bg-emerald-500 text-white rounded text-[10px] font-bold"
                      >
                        Verify GPS
                      </button>
                    )}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </div>
    </div>
  );
};

