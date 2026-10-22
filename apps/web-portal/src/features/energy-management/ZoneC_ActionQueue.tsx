/**
 * Zone C: Main Operational Workspace (Left Pane - 65% Action Queue - Energy Management).
 * 33kV/11kV Feeder Reliability Matrix & 5-Step Safety PTW Permits Queue.
 */

import React, { useState } from "react";
import { Badge } from "../../components/ui/Badge";
import { Table, TableHead, TableBody, TableRow, TableHeaderCell, TableCell } from "../../components/ui/Table";
import { FeederStatusDTO, SafetyPermitDTO } from "./api";

interface ActionQueueProps {
  feeders: FeederStatusDTO[];
  permits: SafetyPermitDTO[];
  onSelectFeeder: (feeder: FeederStatusDTO) => void;
  onSelectPermit: (permit: SafetyPermitDTO) => void;
}

export const ZoneC_ActionQueue: React.FC<ActionQueueProps> = ({
  feeders,
  permits,
  onSelectFeeder,
  onSelectPermit,
}) => {
  const [tab, setTab] = useState<"FEEDERS" | "PTW">("FEEDERS");

  return (
    <div className="bg-slate-800/90 border border-slate-700/80 rounded-xl p-5 shadow-lg flex flex-col h-full">
      <div className="flex items-center justify-between pb-4 border-b border-slate-700/60">
        <div>
          <h3 className="text-sm font-bold text-white tracking-wide">
            {tab === "FEEDERS" ? "33kV & 11kV Feeder Grid Status" : "5-Step Safety Permit-to-Work Queue"}
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">Click any row to open immutable WorkNotes drawer</p>
        </div>

        <div className="flex gap-1.5">
          <button
            onClick={() => setTab("FEEDERS")}
            className={`px-3 py-1 rounded-lg text-xs font-semibold ${
              tab === "FEEDERS" ? "bg-emerald-600 text-white" : "bg-slate-700/60 text-slate-300"
            }`}
          >
            Feeders ({feeders.length})
          </button>
          <button
            onClick={() => setTab("PTW")}
            className={`px-3 py-1 rounded-lg text-xs font-semibold ${
              tab === "PTW" ? "bg-emerald-600 text-white" : "bg-slate-700/60 text-slate-300"
            }`}
          >
            Safety PTW ({permits.length})
          </button>
        </div>
      </div>

      <div className="mt-4 flex-1 overflow-y-auto">
        {tab === "FEEDERS" ? (
          <Table>
            <TableHead>
              <tr>
                <TableHeaderCell>Code</TableHeaderCell>
                <TableHeaderCell>Feeder Name</TableHeaderCell>
                <TableHeaderCell>Load MW</TableHeaderCell>
                <TableHeaderCell>Band</TableHeaderCell>
                <TableHeaderCell>Priority</TableHeaderCell>
                <TableHeaderCell>Status</TableHeaderCell>
              </tr>
            </TableHead>
            <TableBody>
              {feeders.map((f) => (
                <TableRow key={f.id} onClick={() => onSelectFeeder(f)}>
                  <TableCell className="font-mono text-emerald-400 font-semibold">{f.feeder_code}</TableCell>
                  <TableCell>{f.feeder_name}</TableCell>
                  <TableCell className="font-bold">{f.current_load_mw} MW</TableCell>
                  <TableCell>{f.tariff_band}</TableCell>
                  <TableCell>Tier {f.priority_tier}</TableCell>
                  <TableCell>
                    <Badge
                      variant={
                        f.state === "IN_SERVICE"
                          ? "success"
                          : f.state === "OUT_ON_PERMIT_PTW"
                          ? "warning"
                          : "danger"
                      }
                    >
                      {f.state.replace(/_/g, " ")}
                    </Badge>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        ) : (
          <Table>
            <TableHead>
              <tr>
                <TableHeaderCell>Permit #</TableHeaderCell>
                <TableHeaderCell>Feeder / Substation</TableHeaderCell>
                <TableHeaderCell>Isolation Point</TableHeaderCell>
                <TableHeaderCell>5-Step PTW Stage</TableHeaderCell>
              </tr>
            </TableHead>
            <TableBody>
              {permits.map((p) => (
                <TableRow key={p.id} onClick={() => onSelectPermit(p)}>
                  <TableCell className="font-mono text-amber-400 font-semibold">{p.permit_number}</TableCell>
                  <TableCell>
                    {p.feeder_id} &bull; {p.substation_code}
                  </TableCell>
                  <TableCell>{p.isolation_point}</TableCell>
                  <TableCell>
                    <Badge variant={p.step === "RESTORED" ? "success" : "warning"}>{p.step}</Badge>
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

