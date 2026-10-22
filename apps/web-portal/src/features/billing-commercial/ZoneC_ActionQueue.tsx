/**
 * Zone C: Main Operational Workspace (Left Pane - 65% Action Queue - Billing & Commercial).
 * Live Vending Audit Logs & Revenue Protection AI Anomaly Radar.
 */

import React, { useState } from "react";
import { Badge } from "../../components/ui/Badge";
import { Table, TableHead, TableBody, TableRow, TableHeaderCell, TableCell } from "../../components/ui/Table";
import { VendingTransactionDTO, TheftAlertDTO } from "./api";

interface ActionQueueProps {
  transactions: VendingTransactionDTO[];
  alerts: TheftAlertDTO[];
  onSelectTx: (tx: VendingTransactionDTO) => void;
  onSelectAlert: (alert: TheftAlertDTO) => void;
}

export const ZoneC_ActionQueue: React.FC<ActionQueueProps> = ({
  transactions,
  alerts,
  onSelectTx,
  onSelectAlert,
}) => {
  const [tab, setTab] = useState<"TX" | "THEFT">("TX");

  return (
    <div className="bg-slate-800/90 border border-slate-700/80 rounded-xl p-5 shadow-lg flex flex-col h-full">
      <div className="flex items-center justify-between pb-4 border-b border-slate-700/60">
        <div>
          <h3 className="text-sm font-bold text-white tracking-wide">
            {tab === "TX" ? "Metered Vending Ledger (20-Digit STS Tokens)" : "Revenue Protection AI Anomaly Radar"}
          </h3>
          <p className="text-xs text-slate-400 mt-0.5">Click any row to open immutable WorkNotes drawer</p>
        </div>

        <div className="flex gap-1.5">
          <button
            onClick={() => setTab("TX")}
            className={`px-3 py-1 rounded-lg text-xs font-semibold ${
              tab === "TX" ? "bg-emerald-600 text-white" : "bg-slate-700/60 text-slate-300"
            }`}
          >
            Vending ({transactions.length})
          </button>
          <button
            onClick={() => setTab("THEFT")}
            className={`px-3 py-1 rounded-lg text-xs font-semibold ${
              tab === "THEFT" ? "bg-emerald-600 text-white" : "bg-slate-700/60 text-slate-300"
            }`}
          >
            Theft Alarms ({alerts.length})
          </button>
        </div>
      </div>

      <div className="mt-4 flex-1 overflow-y-auto">
        {tab === "TX" ? (
          <Table>
            <TableHead>
              <tr>
                <TableHeaderCell>Ref #</TableHeaderCell>
                <TableHeaderCell>Meter #</TableHeaderCell>
                <TableHeaderCell>Amount</TableHeaderCell>
                <TableHeaderCell>kWh Vended</TableHeaderCell>
                <TableHeaderCell>20-Digit STS Token</TableHeaderCell>
                <TableHeaderCell>Platform</TableHeaderCell>
              </tr>
            </TableHead>
            <TableBody>
              {transactions.map((t) => (
                <TableRow key={t.id} onClick={() => onSelectTx(t)}>
                  <TableCell className="font-mono text-emerald-400 font-medium">{t.transaction_ref}</TableCell>
                  <TableCell>{t.meter_number}</TableCell>
                  <TableCell className="font-bold">₦{t.amount_paid.toLocaleString()}</TableCell>
                  <TableCell>{t.kwh_vended} kWh</TableCell>
                  <TableCell className="font-mono text-amber-300 tracking-wider">
                    {t.sts_token_20_digit}
                  </TableCell>
                  <TableCell>
                    <Badge variant="info">{t.payment_platform_source}</Badge>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        ) : (
          <Table>
            <TableHead>
              <tr>
                <TableHeaderCell>Target ID</TableHeaderCell>
                <TableHeaderCell>Anomaly Type</TableHeaderCell>
                <TableHeaderCell>Confidence Score</TableHeaderCell>
                <TableHeaderCell>Evidence Summary</TableHeaderCell>
                <TableHeaderCell>Status</TableHeaderCell>
              </tr>
            </TableHead>
            <TableBody>
              {alerts.map((a) => (
                <TableRow key={a.id} onClick={() => onSelectAlert(a)}>
                  <TableCell className="font-mono text-rose-400 font-medium">{a.customer_or_dt_id}</TableCell>
                  <TableCell>{a.alert_type.replace(/_/g, " ")}</TableCell>
                  <TableCell className="font-bold text-rose-300">{a.confidence_score}%</TableCell>
                  <TableCell className="truncate max-w-xs">{a.evidence_summary}</TableCell>
                  <TableCell>
                    <Badge variant={a.status === "RESOLVED" ? "success" : "danger"}>{a.status}</Badge>
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

