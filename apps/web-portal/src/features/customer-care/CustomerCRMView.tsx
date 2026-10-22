/**
 * Task View 4: Customer 360 CRM Directory Workspace.
 * Comprehensive consumer lookup, meter tracking, debt ledger standing, and profile inspection.
 */

import React, { useState } from "react";
import { CustomerRecord, KYCStatus, TariffBand } from "./types";
import { Badge } from "../../components/ui/Badge";
import { Table, TableHead, TableBody, TableRow, TableHeaderCell, TableCell } from "../../components/ui/Table";

interface CustomerCRMViewProps {
  customers: CustomerRecord[];
  onLogComplaintForCustomer?: (customer: CustomerRecord) => void;
  onNavigateToKYC?: () => void;
}

export const CustomerCRMView: React.FC<CustomerCRMViewProps> = ({
  customers,
  onLogComplaintForCustomer,
  onNavigateToKYC,
}) => {
  const [search, setSearch] = useState<string>("");
  const [statusFilter, setStatusFilter] = useState<string>("ALL");
  const [tariffFilter, setTariffFilter] = useState<string>("ALL");
  const [selectedCust, setSelectedCust] = useState<CustomerRecord | null>(null);

  const filteredCustomers = customers.filter((c) => {
    if (statusFilter !== "ALL" && c.kyc_status !== statusFilter) return false;
    if (tariffFilter !== "ALL" && c.tariff_band !== tariffFilter) return false;

    if (search.trim()) {
      const q = search.toLowerCase();
      const matchName = c.full_name.toLowerCase().includes(q);
      const matchCIN = c.cin.toLowerCase().includes(q);
      const matchAcc = c.account_no.toLowerCase().includes(q);
      const matchPhone = c.phone.toLowerCase().includes(q);
      const matchMeter = c.meter_serial.toLowerCase().includes(q);
      const matchDT = c.transformer_id.toLowerCase().includes(q);
      return matchName || matchCIN || matchAcc || matchPhone || matchMeter || matchDT;
    }
    return true;
  });

  return (
    <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-6 shadow-xl space-y-6">
      {/* Header and Controls */}
      <div className="flex flex-col md:flex-row md:items-center justify-between pb-4 border-b border-slate-800 gap-4">
        <div>
          <div className="flex items-center gap-2">
            <h2 className="text-base font-bold text-white tracking-wide">
              CRM &bull; Customer 360 Consumer Directory
            </h2>
            <span className="px-2 py-0.5 rounded text-[11px] font-semibold bg-sky-500/10 text-sky-400 border border-sky-500/20">
              {customers.length} Onboarded Consumers
            </span>
          </div>
          <p className="text-xs text-slate-400 mt-1">
            Search customer records, audit debt balances, verify PostGIS meter anchors, and initiate service actions.
          </p>
        </div>

        <div className="flex flex-wrap items-center gap-2.5">
          <input
            type="text"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search CIN, Account #, Name, Meter, DT..."
            className="px-3 py-1.5 rounded-lg bg-slate-950 border border-slate-700 text-xs text-white placeholder-slate-500 focus:outline-none focus:ring-1 focus:ring-emerald-500 w-64"
          />

          <select
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value)}
            className="px-3 py-1.5 rounded-lg bg-slate-950 border border-slate-700 text-xs text-slate-300"
          >
            <option value="ALL">All KYC Statuses</option>
            <option value="VERIFIED">Verified KYC</option>
            <option value="PENDING">Pending Verification</option>
            <option value="DEBTOR_FLAGGED">Debtor Flagged</option>
          </select>

          <select
            value={tariffFilter}
            onChange={(e) => setTariffFilter(e.target.value)}
            className="px-3 py-1.5 rounded-lg bg-slate-950 border border-slate-700 text-xs text-slate-300"
          >
            <option value="ALL">All Tariff Bands</option>
            <option value="BAND_A">Band A (20+ hrs)</option>
            <option value="BAND_B">Band B (16-20 hrs)</option>
            <option value="BAND_C">Band C (12-16 hrs)</option>
          </select>

          {onNavigateToKYC && (
            <button
              onClick={onNavigateToKYC}
              className="px-3 py-1.5 rounded-lg bg-sky-600 hover:bg-sky-500 text-white text-xs font-bold transition shadow"
            >
              + Onboard New Customer
            </button>
          )}
        </div>
      </div>

      {/* Main Content: Table + Customer Details Drawer */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <div className={selectedCust ? "lg:col-span-8 overflow-x-auto" : "lg:col-span-12 overflow-x-auto"}>
          {filteredCustomers.length === 0 ? (
            <div className="text-center py-16 text-slate-500 text-xs">
              No customer records match query.
            </div>
          ) : (
            <Table>
              <TableHead>
                <tr>
                  <TableHeaderCell>CIN &amp; Account</TableHeaderCell>
                  <TableHeaderCell>Consumer Name</TableHeaderCell>
                  <TableHeaderCell>Tariff &amp; Meter</TableHeaderCell>
                  <TableHeaderCell>DT Station</TableHeaderCell>
                  <TableHeaderCell>Ledger Standing</TableHeaderCell>
                  <TableHeaderCell>Status</TableHeaderCell>
                  <TableHeaderCell>Actions</TableHeaderCell>
                </tr>
              </TableHead>
              <TableBody>
                {filteredCustomers.map((c) => {
                  const isSelected = selectedCust?.id === c.id;
                  const isDebtor = c.ledger_balance < 0;
                  return (
                    <TableRow
                      key={c.id}
                      onClick={() => setSelectedCust(c)}
                      className={`cursor-pointer transition-colors ${
                        isSelected ? "bg-slate-800 ring-1 ring-sky-500/50" : "hover:bg-slate-800/60"
                      }`}
                    >
                      <TableCell>
                        <div className="font-mono text-xs font-bold text-sky-400">{c.cin}</div>
                        <div className="font-mono text-[11px] text-slate-400">{c.account_no}</div>
                      </TableCell>
                      <TableCell>
                        <div className="text-xs font-semibold text-white">{c.full_name}</div>
                        <div className="text-[11px] text-slate-400">{c.phone}</div>
                      </TableCell>
                      <TableCell>
                        <div className="text-xs font-bold text-amber-300">{c.tariff_band}</div>
                        <div className="font-mono text-[11px] text-slate-400">{c.meter_serial}</div>
                      </TableCell>
                      <TableCell>
                        <div className="font-mono text-[11px] text-slate-300">{c.transformer_id}</div>
                        <div className="text-[10px] text-slate-500">{c.feeder_line_code}</div>
                      </TableCell>
                      <TableCell>
                        {isDebtor ? (
                          <span className="px-2 py-0.5 rounded text-[11px] font-bold bg-rose-500/20 text-rose-300 border border-rose-500/30">
                            -₦{Math.abs(c.ledger_balance).toLocaleString()} DEBT
                          </span>
                        ) : (
                          <span className="text-xs font-semibold text-emerald-400">
                            ₦{c.ledger_balance.toLocaleString()}
                          </span>
                        )}
                      </TableCell>
                      <TableCell>
                        <Badge
                          variant={
                            c.kyc_status === "VERIFIED"
                              ? "success"
                              : c.kyc_status === "DEBTOR_FLAGGED"
                              ? "danger"
                              : "warning"
                          }
                        >
                          {c.kyc_status}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            setSelectedCust(c);
                          }}
                          className="px-2 py-1 rounded bg-slate-800 hover:bg-slate-700 text-slate-200 text-[11px] font-semibold border border-slate-700"
                        >
                          View 360 &rarr;
                        </button>
                      </TableCell>
                    </TableRow>
                  );
                })}
              </TableBody>
            </Table>
          )}
        </div>

        {/* Selected Customer 360 Card */}
        {selectedCust && (
          <div className="lg:col-span-4 bg-slate-950 border border-slate-800 rounded-xl p-5 space-y-4 animate-fadeIn">
            <div className="flex items-center justify-between border-b border-slate-800 pb-3">
              <h3 className="text-sm font-bold text-white flex items-center gap-2">
                <span>👤 Customer 360 Profile</span>
              </h3>
              <button
                onClick={() => setSelectedCust(null)}
                className="text-slate-400 hover:text-white text-xs font-bold"
              >
                ✕ Close
              </button>
            </div>

            <div className="space-y-3 text-xs">
              <div>
                <span className="text-[11px] text-slate-500 uppercase tracking-wider font-semibold">Customer Name</span>
                <p className="text-sm font-bold text-white">{selectedCust.full_name}</p>
                <p className="text-slate-400">{selectedCust.phone} &bull; {selectedCust.email || "No email"}</p>
              </div>

              <div className="grid grid-cols-2 gap-2 pt-1 border-t border-slate-900">
                <div>
                  <span className="text-[10px] text-slate-500 uppercase font-semibold">CIN Identifier</span>
                  <p className="font-mono font-bold text-sky-400">{selectedCust.cin}</p>
                </div>
                <div>
                  <span className="text-[10px] text-slate-500 uppercase font-semibold">Account Number</span>
                  <p className="font-mono text-slate-300">{selectedCust.account_no}</p>
                </div>
              </div>

              <div className="pt-1 border-t border-slate-900">
                <span className="text-[10px] text-slate-500 uppercase font-semibold">Service Premises</span>
                <p className="text-slate-300 mt-0.5">{selectedCust.service_address}</p>
                <p className="text-[11px] text-slate-500">Class: {selectedCust.building_class}</p>
              </div>

              <div className="grid grid-cols-2 gap-2 pt-1 border-t border-slate-900">
                <div>
                  <span className="text-[10px] text-slate-500 uppercase font-semibold">Connected DT</span>
                  <p className="font-mono text-slate-300 text-[11px]">{selectedCust.transformer_id}</p>
                </div>
                <div>
                  <span className="text-[10px] text-slate-500 uppercase font-semibold">Tariff Band</span>
                  <p className="font-bold text-amber-400">{selectedCust.tariff_band}</p>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-2 pt-1 border-t border-slate-900">
                <div>
                  <span className="text-[10px] text-slate-500 uppercase font-semibold">Meter Serial</span>
                  <p className="font-mono text-slate-300">{selectedCust.meter_serial}</p>
                  <p className="text-[10px] text-slate-500">{selectedCust.meter_type}</p>
                </div>
                <div>
                  <span className="text-[10px] text-slate-500 uppercase font-semibold">Ledger Balance</span>
                  <p className={`font-bold ${selectedCust.ledger_balance < 0 ? "text-rose-400" : "text-emerald-400"}`}>
                    ₦{selectedCust.ledger_balance.toLocaleString()}
                  </p>
                </div>
              </div>

              <div className="pt-1 border-t border-slate-900">
                <span className="text-[10px] text-slate-500 uppercase font-semibold">PostGIS GPS Coordinates</span>
                <p className="font-mono text-[11px] text-emerald-400">
                  Lat: {selectedCust.latitude.toFixed(4)}, Lon: {selectedCust.longitude.toFixed(4)} (Verified DT Anchor)
                </p>
              </div>

              <div className="pt-3 border-t border-slate-800 space-y-2">
                {onLogComplaintForCustomer && (
                  <button
                    onClick={() => onLogComplaintForCustomer(selectedCust)}
                    className="w-full py-2 bg-emerald-600 hover:bg-emerald-500 text-white rounded-lg font-bold text-xs shadow transition"
                  >
                    📝 Log Complaint For This Customer
                  </button>
                )}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};
