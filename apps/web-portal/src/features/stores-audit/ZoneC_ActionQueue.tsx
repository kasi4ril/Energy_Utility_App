/**
 * Zone C: Main Operational Workspace (Left Pane - 65% Action Queue - Stores & Audit).
 * Dual-Path Material Requisitions & Staff Allowance Claims Desk.
 */

import React from "react";
import { Badge } from "../../components/ui/Badge";
import { Table, TableHead, TableBody, TableRow, TableHeaderCell, TableCell } from "../../components/ui/Table";
import { RequisitionDTO } from "./api";

interface ActionQueueProps {
  requisitions: RequisitionDTO[];
  onSelectRequisition: (req: RequisitionDTO) => void;
}

export const ZoneC_ActionQueue: React.FC<ActionQueueProps> = ({ requisitions, onSelectRequisition }) => {
  return (
    <div className="bg-slate-800/90 border border-slate-700/80 rounded-xl p-5 shadow-lg flex flex-col h-full">
      <div className="pb-4 border-b border-slate-700/60 flex items-center justify-between">
        <div>
          <h3 className="text-sm font-bold text-white tracking-wide">
            Dual-Path Material Requisitions Queue
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Path A (In-Stock Audit Clearance) &bull; Path B (Out-of-Stock MD Procurement)
          </p>
        </div>
      </div>

      <div className="mt-4 flex-1 overflow-y-auto">
        {requisitions.length === 0 ? (
          <div className="text-center py-12 text-xs text-slate-500">
            No active material requisitions in queue.
          </div>
        ) : (
          <Table>
            <TableHead>
              <tr>
                <TableHeaderCell>Req #</TableHeaderCell>
                <TableHeaderCell>SKU</TableHeaderCell>
                <TableHeaderCell>Qty</TableHeaderCell>
                <TableHeaderCell>Type</TableHeaderCell>
                <TableHeaderCell>Audit Consent</TableHeaderCell>
                <TableHeaderCell>Status</TableHeaderCell>
              </tr>
            </TableHead>
            <TableBody>
              {requisitions.map((r) => (
                <TableRow key={r.id} onClick={() => onSelectRequisition(r)}>
                  <TableCell className="font-mono text-emerald-400 font-medium">{r.requisition_no}</TableCell>
                  <TableCell>{r.item_sku}</TableCell>
                  <TableCell className="font-bold">{r.quantity_requested}</TableCell>
                  <TableCell>
                    <Badge variant={r.is_out_of_stock_procurement ? "warning" : "info"}>
                      {r.is_out_of_stock_procurement ? "PATH B (MD PROCUREMENT)" : "PATH A (IN-STOCK)"}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    {r.internal_audit_consent ? (
                      <span className="text-emerald-400 font-bold text-[11px]">CLEARED</span>
                    ) : (
                      <span className="text-amber-400 font-bold text-[11px]">PENDING</span>
                    )}
                  </TableCell>
                  <TableCell>
                    <Badge variant={r.status === "DISPATCHED" ? "success" : "neutral"}>{r.status}</Badge>
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

