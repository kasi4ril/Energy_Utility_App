/**
 * Zone C: Main Operational Workspace (Left Pane - 65% Action Queue - Generation Master).
 * Live telemetry logs and plant workshop reconditioning jobs.
 */

import React from "react";
import { Badge } from "../../components/ui/Badge";
import { Table, TableHead, TableBody, TableRow, TableHeaderCell, TableCell } from "../../components/ui/Table";
import { PlantTelemetryDTO } from "./api";

interface ActionQueueProps {
  logs: PlantTelemetryDTO[];
  onSelectLog: (log: PlantTelemetryDTO) => void;
}

export const ZoneC_ActionQueue: React.FC<ActionQueueProps> = ({ logs, onSelectLog }) => {
  return (
    <div className="bg-slate-800/90 border border-slate-700/80 rounded-xl p-5 shadow-lg flex flex-col h-full">
      <div className="pb-4 border-b border-slate-700/60 flex items-center justify-between">
        <div>
          <h3 className="text-sm font-bold text-white tracking-wide">6-Plant Station Collation Logs</h3>
          <p className="text-xs text-slate-400 mt-0.5">Real-time MW, Hydro Head & Vibration Telemetry</p>
        </div>
        <span className="text-xs font-semibold px-2.5 py-1 bg-slate-700/60 text-slate-300 rounded-lg">
          Live Shift Telemetry
        </span>
      </div>

      <div className="mt-4 flex-1 overflow-y-auto">
        {logs.length === 0 ? (
          <div className="text-center py-12 text-xs text-slate-500">
            Awaiting hourly station log collation from Plants 1–6.
          </div>
        ) : (
          <Table>
            <TableHead>
              <tr>
                <TableHeaderCell>Plant ID</TableHeaderCell>
                <TableHeaderCell>MW Generated</TableHeaderCell>
                <TableHeaderCell>Water Level (m)</TableHeaderCell>
                <TableHeaderCell>Vibration (mm/s)</TableHeaderCell>
                <TableHeaderCell>Stator Temp (°C)</TableHeaderCell>
                <TableHeaderCell>Status</TableHeaderCell>
              </tr>
            </TableHead>
            <TableBody>
              {logs.map((l) => (
                <TableRow key={l.id} onClick={() => onSelectLog(l)}>
                  <TableCell className="font-mono text-emerald-400 font-semibold">{l.plant_id}</TableCell>
                  <TableCell className="font-bold">{l.mw_generated} MW</TableCell>
                  <TableCell>{(l as any).water_level_meters ?? l.water_level_m} m</TableCell>
                  <TableCell>
                    <span className={l.bearing_vibration_mm_s >= 4.5 ? "text-rose-400 font-bold" : ""}>
                      {l.bearing_vibration_mm_s} mm/s
                    </span>
                  </TableCell>
                  <TableCell>{(l as any).stator_temp_celsius ?? l.stator_winding_temp_c}°C</TableCell>
                  <TableCell>
                    <Badge variant={l.plant_state === "IN_SERVICE" || (l as any).in_service ? "success" : "danger"}>
                      {l.plant_state || ((l as any).in_service ? "IN SERVICE" : "OUT")}
                    </Badge>
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

