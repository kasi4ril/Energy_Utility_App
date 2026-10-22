/**
 * Task View 2: Dedicated Complaint & Fault Intake Form.
 * Allows officers to enter customer complaints, lookup accounts, and route directly to engineering.
 */

import React, { useState } from "react";
import { CustomerRecord, NewComplaintPayload, TicketCategory, TicketChannel, TicketPriority } from "./types";
import { createServiceTicket } from "./api";

interface LogComplaintFormProps {
  customers: CustomerRecord[];
  onTicketCreated?: () => void;
  onNavigateToQueue?: () => void;
}

export const LogComplaintForm: React.FC<LogComplaintFormProps> = ({
  customers,
  onTicketCreated,
  onNavigateToQueue,
}) => {
  const [selectedCustId, setSelectedCustId] = useState<string>("");
  const [customerName, setCustomerName] = useState<string>("");
  const [customerPhone, setCustomerPhone] = useState<string>("");
  const [serviceAddress, setServiceAddress] = useState<string>("");
  const [feederCode, setFeederCode] = useState<string>("FDR-33KV-VOM");
  const [dtCode, setDtCode] = useState<string>("TX-DT-500KVA-MARKET");
  const [channel, setChannel] = useState<TicketChannel>("CALL_CENTER");
  const [category, setCategory] = useState<TicketCategory>("LINE_SNAP");
  const [priority, setPriority] = useState<TicketPriority>("HIGH");
  const [assignedUnit, setAssignedUnit] = useState<string>("FEEDER_MAINT");
  const [description, setDescription] = useState<string>("");
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [generatedTicketNo, setGeneratedTicketNo] = useState<string | null>(null);

  const handleSelectCustomer = (id: string) => {
    setSelectedCustId(id);
    const found = customers.find((c) => c.id === id);
    if (found) {
      setCustomerName(found.full_name);
      setCustomerPhone(found.phone);
      setServiceAddress(found.service_address);
      setFeederCode(found.feeder_line_code);
      setDtCode(found.transformer_id);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setGeneratedTicketNo(null);

    const payload: NewComplaintPayload = {
      customer_id: selectedCustId || undefined,
      customer_name: customerName,
      customer_phone: customerPhone,
      service_address: serviceAddress,
      feeder_line_code: feederCode,
      dt_station_code: dtCode,
      channel,
      category,
      priority,
      assigned_unit_id: assignedUnit,
      description,
    };

    try {
      const ticket = await createServiceTicket(payload);
      setGeneratedTicketNo(ticket.ticket_number);
      setDescription("");
      if (onTicketCreated) onTicketCreated();
    } catch (err: any) {
      alert(err.message || "Failed to log complaint");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-6 shadow-xl space-y-6">
      <div className="border-b border-slate-800 pb-4 flex flex-col sm:flex-row sm:items-center justify-between gap-2">
        <div>
          <h2 className="text-base font-bold text-white tracking-wide flex items-center gap-2">
            <span>📝 Omnichannel Complaint &amp; Fault Intake Workspace</span>
          </h2>
          <p className="text-xs text-slate-400 mt-0.5">
            Log grievances, technical outages, and metering faults with automatic engineering unit dispatch.
          </p>
        </div>
        {onNavigateToQueue && (
          <button
            onClick={onNavigateToQueue}
            className="text-xs text-emerald-400 hover:text-emerald-300 font-semibold flex items-center gap-1"
          >
            &larr; Return to Triage Queue
          </button>
        )}
      </div>

      {generatedTicketNo && (
        <div className="p-4 rounded-xl bg-emerald-500/10 border border-emerald-500/30 text-emerald-300 flex items-center justify-between text-xs animate-fadeIn">
          <div>
            <span className="font-bold">Complaint Logged Successfully!</span> Dispatch Ticket{" "}
            <span className="font-mono font-extrabold text-white bg-emerald-900/60 px-2 py-0.5 rounded ml-1">
              {generatedTicketNo}
            </span>{" "}
            has been routed to <span className="font-semibold text-white">{assignedUnit}</span>. SLA countdown initialized.
          </div>
          {onNavigateToQueue && (
            <button
              onClick={onNavigateToQueue}
              className="px-3 py-1 bg-emerald-600 hover:bg-emerald-500 text-white rounded font-bold text-xs shadow"
            >
              View In Queue &rarr;
            </button>
          )}
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-5 text-xs">
        {/* Customer Quick Lookup */}
        <div className="bg-slate-950/70 border border-slate-800 rounded-xl p-4 space-y-3">
          <div className="flex items-center justify-between">
            <label className="text-xs font-bold text-slate-200">
              1. Customer Lookup &amp; Relationship Link (Optional)
            </label>
            <span className="text-[11px] text-slate-400">Select existing account or leave blank for new caller</span>
          </div>

          <select
            value={selectedCustId}
            onChange={(e) => handleSelectCustomer(e.target.value)}
            className="w-full px-3 py-2 rounded-lg bg-slate-900 border border-slate-700 text-white focus:outline-none focus:ring-1 focus:ring-emerald-500"
          >
            <option value="">-- Manual Walk-In / Unregistered Consumer --</option>
            {customers.map((c) => (
              <option key={c.id} value={c.id}>
                {c.cin} &bull; {c.full_name} ({c.phone}) - {c.tariff_band} - Bal: ₦{c.ledger_balance.toLocaleString()}
              </option>
            ))}
          </select>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-3 pt-1">
            <div>
              <label className="block text-[11px] font-semibold text-slate-400 mb-1">Customer Full Name *</label>
              <input
                type="text"
                required
                value={customerName}
                onChange={(e) => setCustomerName(e.target.value)}
                placeholder="e.g. Pam Gyang"
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white placeholder-slate-500 focus:outline-none focus:ring-1 focus:ring-emerald-500"
              />
            </div>
            <div>
              <label className="block text-[11px] font-semibold text-slate-400 mb-1">Phone Number *</label>
              <input
                type="text"
                required
                value={customerPhone}
                onChange={(e) => setCustomerPhone(e.target.value)}
                placeholder="08031234567"
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white placeholder-slate-500 focus:outline-none focus:ring-1 focus:ring-emerald-500"
              />
            </div>
            <div>
              <label className="block text-[11px] font-semibold text-slate-400 mb-1">Premises Service Address *</label>
              <input
                type="text"
                required
                value={serviceAddress}
                onChange={(e) => setServiceAddress(e.target.value)}
                placeholder="Plot 10 Jos South, Plateau"
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white placeholder-slate-500 focus:outline-none focus:ring-1 focus:ring-emerald-500"
              />
            </div>
          </div>
        </div>

        {/* Classification & Routing Matrix */}
        <div className="bg-slate-950/70 border border-slate-800 rounded-xl p-4 space-y-3">
          <label className="text-xs font-bold text-slate-200 block">
            2. Channel, Fault Classification &amp; Engineering Target
          </label>

          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-3">
            <div>
              <label className="block text-[11px] font-semibold text-slate-400 mb-1">Intake Channel</label>
              <select
                value={channel}
                onChange={(e) => setChannel(e.target.value as TicketChannel)}
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white focus:outline-none focus:ring-1 focus:ring-emerald-500"
              >
                <option value="CALL_CENTER">Call Center (Inbound)</option>
                <option value="WALK_IN">Walk-In Customer Care Desk</option>
                <option value="ONLINE_PORTAL">Online Web Portal / Self-Service</option>
                <option value="FIELD_AGENT">Field Agent / Inspection</option>
              </select>
            </div>

            <div>
              <label className="block text-[11px] font-semibold text-slate-400 mb-1">Fault Category</label>
              <select
                value={category}
                onChange={(e) => setCategory(e.target.value as TicketCategory)}
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white focus:outline-none focus:ring-1 focus:ring-emerald-500"
              >
                <option value="LINE_SNAP">Line Snap / Fallen Conductor</option>
                <option value="FALLEN_POLE">Fallen Concrete/Wood Pole</option>
                <option value="TRANSFORMER_FAULT">Distribution Transformer Fault</option>
                <option value="FUSE_BLOWN">11kV/33kV Fuse Blown</option>
                <option value="OIL_LEAK">Transformer Oil Leak</option>
                <option value="FAULTY_METER">Faulty / Tampered Meter</option>
                <option value="NEW_METER_SCHEDULE">New Meter Schedule</option>
                <option value="WRONG_BILLING">Disputed Estimated Bill</option>
                <option value="UNCREDITED_PAYMENT">Uncredited Vending Token</option>
                <option value="GENERAL_INQUIRY">General Customer Inquiry</option>
              </select>
            </div>

            <div>
              <label className="block text-[11px] font-semibold text-slate-400 mb-1">NERC Priority</label>
              <select
                value={priority}
                onChange={(e) => setPriority(e.target.value as TicketPriority)}
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white focus:outline-none focus:ring-1 focus:ring-emerald-500"
              >
                <option value="CRITICAL">Critical (4h SLA - Hazard/Blackout)</option>
                <option value="HIGH">High (8h SLA)</option>
                <option value="MEDIUM">Medium (24h SLA)</option>
                <option value="LOW">Low (48h SLA)</option>
              </select>
            </div>

            <div>
              <label className="block text-[11px] font-semibold text-slate-400 mb-1">Assigned Maintenance Unit</label>
              <select
                value={assignedUnit}
                onChange={(e) => setAssignedUnit(e.target.value)}
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white focus:outline-none focus:ring-1 focus:ring-emerald-500"
              >
                <option value="FEEDER_MAINT">Feeder &amp; Overhead Lines</option>
                <option value="UNDERGROUND_CABLES">Underground Cable Unit</option>
                <option value="METERING_TEST">Metering &amp; Protection Test</option>
                <option value="SUBSTATION_MAINT">Substation Transformers</option>
                <option value="COMMERCIAL_OPS">Commercial &amp; Billing Desk</option>
              </select>
            </div>
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-2">
            <div>
              <label className="block text-[11px] font-semibold text-slate-400 mb-1">Feeder Line Code</label>
              <input
                type="text"
                value={feederCode}
                onChange={(e) => setFeederCode(e.target.value)}
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white font-mono"
              />
            </div>
            <div>
              <label className="block text-[11px] font-semibold text-slate-400 mb-1">Distribution Transformer (DT)</label>
              <input
                type="text"
                value={dtCode}
                onChange={(e) => setDtCode(e.target.value)}
                className="w-full px-3 py-2 rounded bg-slate-900 border border-slate-700 text-white font-mono"
              />
            </div>
          </div>
        </div>

        {/* Detailed Narrative */}
        <div>
          <label className="block text-[11px] font-semibold text-slate-300 mb-1">
            3. Detailed Grievance / Outage Incident Narrative *
          </label>
          <textarea
            required
            rows={4}
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Document precise fault location, visible physical damage (sparks, broken cross-arms), affected streets, and consumer statements..."
            className="w-full px-3 py-2.5 rounded-lg bg-slate-950 border border-slate-700 text-white placeholder-slate-500 focus:outline-none focus:ring-1 focus:ring-emerald-500"
          />
        </div>

        <div className="flex justify-end gap-3 pt-2">
          <button
            type="submit"
            disabled={submitting}
            className="px-6 py-2.5 bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 text-white font-bold text-xs rounded-lg transition shadow-lg flex items-center gap-2"
          >
            {submitting ? "Dispatching..." : "⚡ Log Complaint & Dispatch Field Ticket"}
          </button>
        </div>
      </form>
    </div>
  );
};

