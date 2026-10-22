/**
 * Zone C: Main Operational Workspace (Left Pane - 65% Action Queue).
 * Live-updating queue of incoming complaints and attended-to verification tickets.
 */

import React, { useState } from "react";
import { Badge } from "../../components/ui/Badge";
import { Table, TableHead, TableBody, TableRow, TableHeaderCell, TableCell } from "../../components/ui/Table";
import { ServiceTicketDTO } from "./api";

interface ActionQueueProps {
  tickets: ServiceTicketDTO[];
  onSelectTicket: (ticket: ServiceTicketDTO) => void;
  onConfirmClose?: (ticketId: string) => void;
}

export const ZoneC_ActionQueue: React.FC<ActionQueueProps> = ({
  tickets,
  onSelectTicket,
  onConfirmClose,
}) => {
  const [filter, setFilter] = useState<string>("ALL");

  const filteredTickets = tickets.filter((t) => {
    if (filter === "ATTENDED_TO") return t.status === "ATTENDED_TO";
    if (filter === "OPEN") return t.status === "LOGGED" || t.status === "ROUTED";
    return true;
  });

  return (
    <div className="bg-slate-800/90 border border-slate-700/80 rounded-xl p-5 shadow-lg flex flex-col h-full">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between pb-4 border-b border-slate-700/60 gap-3">
        <div>
          <h3 className="text-sm font-bold text-white tracking-wide">
            Operational Action Queue &bull; Customer Care
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Select any ticket row to inspect immutable WorkNotes hash-chain
          </p>
        </div>

        <div className="flex gap-1.5">
          {["ALL", "OPEN", "ATTENDED_TO"].map((tab) => (
            <button
              key={tab}
              onClick={() => setFilter(tab)}
              className={`px-3 py-1 rounded-lg text-xs font-semibold transition-colors ${
                filter === tab
                  ? "bg-emerald-600 text-white"
                  : "bg-slate-700/60 text-slate-300 hover:bg-slate-700"
              }`}
            >
              {tab.replace("_", " ")}
            </button>
          ))}
        </div>
      </div>

      <div className="mt-4 flex-1 overflow-y-auto">
        {filteredTickets.length === 0 ? (
          <div className="text-center py-12 text-xs text-slate-500">
            No active tickets matching filter. All customer queues clear.
          </div>
        ) : (
          <Table>
            <TableHead>
              <tr>
                <TableHeaderCell>Ticket #</TableHeaderCell>
                <TableHeaderCell>Category</TableHeaderCell>
                <TableHeaderCell>Assigned Unit</TableHeaderCell>
                <TableHeaderCell>Priority</TableHeaderCell>
                <TableHeaderCell>Status</TableHeaderCell>
                <TableHeaderCell>Action</TableHeaderCell>
              </tr>
            </TableHead>
            <TableBody>
              {filteredTickets.map((t) => (
                <TableRow key={t.id} onClick={() => onSelectTicket(t)}>
                  <TableCell className="font-mono font-medium text-emerald-400">
                    {t.ticket_number}
                  </TableCell>
                  <TableCell>{t.category}</TableCell>
                  <TableCell>{t.unit_type || "TRIAGE"}</TableCell>
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
                    {t.status === "ATTENDED_TO" && onConfirmClose && (
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          onConfirmClose(t.id);
                        }}
                        className="px-2 py-1 rounded bg-amber-600 hover:bg-amber-500 text-white text-[10px] font-bold"
                      >
                        Call & Confirm
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

