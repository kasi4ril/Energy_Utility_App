/**
 * Zone C: Main Operational Workspace (Left Pane - 65% Action Queue - Executive Command).
 * One-Click Digital Sign-Off Queue with Cryptographic Role Segregation of Duties.
 */

import React, { useState } from "react";
import { Badge } from "../../components/ui/Badge";
import { Table, TableHead, TableBody, TableRow, TableHeaderCell, TableCell } from "../../components/ui/Table";
import { ActionItemDTO } from "./api";

interface ActionQueueProps {
  actionItems: ActionItemDTO[];
  userRole: string;
  onSelectAction: (item: ActionItemDTO) => void;
  onDecideAction: (id: string, approve: boolean) => void;
}

export const ZoneC_ActionQueue: React.FC<ActionQueueProps> = ({
  actionItems = [],
  userRole,
  onSelectAction,
  onDecideAction,
}) => {
  const [filterRole, setFilterRole] = useState(userRole || "ALL");

  const filtered = (actionItems || []).filter((i) => {
    if (filterRole === "ALL") return true;
    return i.target_role === filterRole;
  });

  return (
    <div className="bg-slate-800/90 border border-slate-700/80 rounded-xl p-5 shadow-lg flex flex-col h-full">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between pb-4 border-b border-slate-700/60 gap-3">
        <div>
          <h3 className="text-sm font-bold text-white tracking-wide">
            Executive Digital Sign-Off Queue
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">
            Role Segregated: MD &bull; Commercial Ops &bull; Audit &bull; Services &bull; Generation
          </p>
        </div>

        <select
          value={filterRole}
          onChange={(e) => setFilterRole(e.target.value)}
          className="px-3 py-1.5 rounded-lg bg-slate-900 border border-slate-700 text-xs text-white focus:outline-none focus:ring-1 focus:ring-emerald-500"
        >
          <option value="ALL">All Executive Roles</option>
          <option value="MD">Managing Director (MD)</option>
          <option value="HEAD_COMMERCIAL_OPS">Head Commercial Ops</option>
          <option value="HEAD_INTERNAL_AUDIT">Head Internal Audit</option>
          <option value="HEAD_SERVICES">Head Services / Safety</option>
          <option value="GM_GENERATION">GM Generation</option>
        </select>
      </div>

      <div className="mt-4 flex-1 overflow-y-auto">
        {filtered.length === 0 ? (
          <div className="text-center py-12 text-xs text-slate-500">
            No pending action items for this executive role.
          </div>
        ) : (
          <Table>
            <TableHead>
              <tr>
                <TableHeaderCell>Action Type</TableHeaderCell>
                <TableHeaderCell>Target Role</TableHeaderCell>
                <TableHeaderCell>Department</TableHeaderCell>
                <TableHeaderCell>Amount</TableHeaderCell>
                <TableHeaderCell>Summary</TableHeaderCell>
                <TableHeaderCell>Digital Decision</TableHeaderCell>
              </tr>
            </TableHead>
            <TableBody>
              {filtered.map((item) => {
                const canSign = userRole === item.target_role;
                return (
                  <TableRow key={item.id} onClick={() => onSelectAction(item)}>
                    <TableCell className="font-mono text-emerald-400 font-semibold">
                      {item.action_type.replace(/_/g, " ")}
                    </TableCell>
                    <TableCell>
                      <Badge variant={canSign ? "success" : "neutral"}>{item.target_role}</Badge>
                    </TableCell>
                    <TableCell>{item.source_department}</TableCell>
                    <TableCell className="font-bold">
                      {item.amount ? `₦${item.amount.toLocaleString()}` : "—"}
                    </TableCell>
                    <TableCell className="truncate max-w-xs">{item.summary}</TableCell>
                    <TableCell>
                      {canSign ? (
                        <div className="flex gap-1.5" onClick={(e) => e.stopPropagation()}>
                          <button
                            onClick={() => onDecideAction(item.id, true)}
                            className="px-2.5 py-1 bg-emerald-600 hover:bg-emerald-500 text-white rounded text-[11px] font-bold shadow"
                          >
                            Approve
                          </button>
                          <button
                            onClick={() => onDecideAction(item.id, false)}
                            className="px-2.5 py-1 bg-rose-700 hover:bg-rose-600 text-white rounded text-[11px] font-bold"
                          >
                            Reject
                          </button>
                        </div>
                      ) : (
                        <span className="text-[10px] text-slate-500 italic">Role Restricted</span>
                      )}
                    </TableCell>
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        )}
      </div>
    </div>
  );
};

