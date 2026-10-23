/**
 * View 3: Smart Complaint Form & Departmental Routing Desk (/portal/customer-care/complain).
 * Strictly <= 350 lines.
 */
import React, { useState } from "react";
import { createServiceTicket } from "./api";
import { CustomerRecord, TicketChannel, TicketCategory, TicketPriority } from "./types";
import { AlertTriangle, Send, Flame, ShieldAlert, CheckCircle, Search } from "lucide-react";

interface View3Props {
  customers: CustomerRecord[];
  staffId: string;
  staffName?: string;
  onTicketLogged?: () => void;
}

export function View3_SmartComplaintDesk({
  customers,
  staffId,
  staffName = "Amina Yusuf (CC Officer)",
  onTicketLogged,
}: View3Props) {
  const [channel, setChannel] = useState<TicketChannel>("CALL_CENTER");
  const [accountQuery, setAccountQuery] = useState("");
  const [matchedCustomer, setMatchedCustomer] = useState<CustomerRecord | null>(null);
  const [customerName, setCustomerName] = useState("");
  const [customerPhone, setCustomerPhone] = useState("");
  const [serviceAddress, setServiceAddress] = useState("");
  const [landmark, setLandmark] = useState("");
  const [category, setCategory] = useState<TicketCategory>("LINE_SNAP");
  const [subCategory, setSubCategory] = useState("Low-hanging conductor across road");
  const [isEmergency, setIsEmergency] = useState(false);
  const [assignedUnitName, setAssignedUnitName] = useState("Feeder Maintenance Unit Head (Services)");
  const [description, setDescription] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [successInfo, setSuccessInfo] = useState<string | null>(null);

  const populateFromCustomer = (cust: CustomerRecord) => {
    setMatchedCustomer(cust);
    setCustomerName(cust.full_name);
    setCustomerPhone(cust.phone);
    setServiceAddress(cust.service_address);
    setLandmark(cust.landmark || "");
  };

  const handleAccountInputChange = (val: string) => {
    setAccountQuery(val);
    const q = val.trim().toLowerCase();
    if (q.length >= 4) {
      const found = customers.find(
        (c) => c.account_no.toLowerCase().includes(q) || c.meter_serial.toLowerCase().includes(q) || c.phone.includes(q)
      );
      if (found) populateFromCustomer(found);
    }
  };

  const handleSearchAccount = () => {
    const q = accountQuery.trim().toLowerCase();
    const found = customers.find(
      (c) => c.account_no.toLowerCase().includes(q) || c.meter_serial.toLowerCase().includes(q) || c.phone.includes(q)
    );
    if (found) {
      populateFromCustomer(found);
    } else {
      alert("No customer account matched. You may proceed as guest/walk-in.");
    }
  };

  const handleCategoryChange = (newCat: TicketCategory) => {
    setCategory(newCat);
    const map: Record<string, { u: string; s: string; em?: boolean }> = {
      FALLEN_POLE: { u: "Feeder Maintenance Unit Head (Services)", s: "Broken electric pole" },
      LINE_SNAP: { u: "Feeder Maintenance Unit Head (Services)", s: "Conductor snap / line fault" },
      TRANSFORMER_FAULT: { u: "Transformer Maintenance Unit Head (Services)", s: "Transformer fault / blowout" },
      OIL_LEAK: { u: "Transformer Maintenance Unit Head (Services)", s: "Dielectric oil leakage" },
      FUSE_BLOWN: { u: "Transformer Maintenance Unit Head (Services)", s: "Drop-out fuse replacement" },
      FAULTY_METER: { u: "Meter Maintenance & Installation Unit Head (Services)", s: "STS tamper error 30" },
      NEW_METER_SCHEDULE: { u: "Meter Maintenance & Installation Unit Head (Services)", s: "MAP survey appointment" },
      WRONG_BILLING: { u: "Billing Unit Head (Commercial)", s: "Estimated billing dispute" },
      UNCREDITED_PAYMENT: { u: "STS Vending & Billing IT Support Desk (Commercial)", s: "Missing vending token" },
      SAFETY_HAZARD_EMERGENCY: { u: "Feeder Maintenance Unit Head (Services)", s: "Live wire contact hazard", em: true },
    };
    const c = map[newCat];
    if (c) {
      setAssignedUnitName(c.u);
      setSubCategory(c.s);
      if (c.em) setIsEmergency(true);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setSuccessInfo(null);
    try {
      const payload = {
        customer_id: matchedCustomer?.id, customer_name: customerName,
        customer_phone: customerPhone, service_address: serviceAddress, landmark,
        channel, category, sub_category: subCategory, is_safety_emergency: isEmergency,
        priority: (isEmergency ? "CRITICAL" : "MEDIUM") as TicketPriority,
        assigned_unit_id: "UNIT-DEPT-DIRECT", assigned_unit_name: assignedUnitName,
        description, logged_by_cc_staff_id: staffId, logged_by_cc_staff_name: staffName,
        routed_by_cc_staff_id: staffId, routed_by_cc_staff_name: staffName,
      };
      const res = await createServiceTicket(payload);
      setSuccessInfo(`Ticket ${res.ticket_number} routed successfully to ${assignedUnitName}! SLA: ${isEmergency ? "<2 Hours (EMERGENCY)" : "4 Hours"}`);
      if (onTicketLogged) onTicketLogged();
      setDescription("");
      setIsEmergency(false);
    } catch (err: any) {
      alert(err.message || "Failed to create ticket");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-4 sm:p-6 max-w-4xl">
      <div className="flex items-center space-x-3 mb-6 pb-4 border-b border-slate-800">
        <div className="p-2.5 rounded-lg bg-blue-600/20 text-blue-400 border border-blue-500/30">
          <AlertTriangle className="w-6 h-6" />
        </div>
        <div>
          <h2 className="text-base font-bold text-slate-100 uppercase tracking-wide">
            Smart Complaint Intake &amp; Departmental Routing Desk
          </h2>
          <p className="text-xs text-slate-400">
            Omnichannel triage &bull; Auto-suggested routing &bull; NERC &lt;2-Hour Safety Emergency Gate
          </p>
        </div>
      </div>

      {/* Safety Hazard Trigger Banner */}
      <div
        onClick={() => {
          setIsEmergency(!isEmergency);
          if (!isEmergency) {
            setCategory("SAFETY_HAZARD_EMERGENCY");
            setAssignedUnitName("Feeder Maintenance Unit Head (Services)");
          }
        }}
        className={`p-3.5 mb-5 rounded-xl border cursor-pointer transition-all flex items-center justify-between ${
          isEmergency ? "bg-red-950/80 border-red-500 text-red-200 animate-pulse" : "bg-slate-950/50 border-slate-800 hover:border-red-500/40 text-slate-300"
        }`}
      >
        <div className="flex items-center space-x-3">
          <Flame className={`w-5 h-5 ${isEmergency ? "text-red-400 animate-bounce" : "text-amber-500"}`} />
          <div>
            <div className="text-xs font-bold uppercase tracking-wider flex items-center space-x-2">
              <span>Life-Threatening Safety Hazard Triage</span>
              <span className="text-[10px] bg-red-600 text-white px-2 py-0.2 rounded-full font-bold">&lt; 2-Hour SLA</span>
            </div>
            <p className="text-[11px] text-slate-400">Live Wire Down, Transformer Fire, or Public Contact Risk.</p>
          </div>
        </div>
        <div className="text-xs font-mono font-bold text-red-400">
          {isEmergency ? "[ACTIVATED - 2HR CLOCK]" : "[CLICK TO ACTIVATE]"}
        </div>
      </div>

      {successInfo && (
        <div className="p-3.5 mb-5 rounded-lg bg-emerald-950/60 border border-emerald-500/40 text-emerald-300 text-xs flex items-center space-x-2.5">
          <CheckCircle className="w-4 h-4 shrink-0 text-emerald-400" />
          <span>{successInfo}</span>
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-4 text-xs">
        {/* Row 1: Omnichannel Source & Meter Verification */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label className="block text-slate-300 font-medium mb-1">Omnichannel Source *</label>
            <select
              value={channel}
              onChange={(e) => setChannel(e.target.value as TicketChannel)}
              className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-slate-100 focus:border-blue-500 focus:outline-none"
            >
              <option value="CALL_CENTER">Call Center (Toll-Free Hotline)</option>
              <option value="WALK_IN">Walk-In Customer Care Desk</option>
              <option value="ONLINE_PORTAL">Online Consumer Web Portal</option>
              <option value="WHATSAPP">WhatsApp Utility Bot (+234 800-NESCO)</option>
            </select>
          </div>

          <div className="md:col-span-2">
            <label className="block text-slate-300 font-medium mb-1">
              Account No. / Meter Verification {matchedCustomer && <span className="text-[10px] text-emerald-400 font-mono">(Verified &bull; Auto-Populated)</span>}
            </label>
            <div className="flex space-x-2">
              <input
                type="text"
                placeholder="Type Account No, Meter Serial or Phone to auto-populate..."
                value={accountQuery}
                onChange={(e) => handleAccountInputChange(e.target.value)}
                className="flex-1 bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-slate-100 font-mono focus:border-blue-500 focus:outline-none"
              />
              <button
                type="button"
                onClick={handleSearchAccount}
                className="bg-slate-800 hover:bg-slate-700 text-slate-200 px-3 py-2 rounded-lg transition-colors flex items-center space-x-1"
              >
                <Search className="w-3.5 h-3.5" />
                <span>Verify</span>
              </button>
            </div>
          </div>
        </div>

        {/* Row 2: Customer Details & Landmark */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-3">
          <div>
            <label className="block text-slate-300 font-medium mb-1">Complainant Name *</label>
            <input
              type="text"
              required
              value={customerName}
              onChange={(e) => setCustomerName(e.target.value)}
              className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-slate-100 focus:border-blue-500 focus:outline-none"
            />
          </div>
          <div>
            <label className="block text-slate-300 font-medium mb-1">Contact Phone *</label>
            <input
              type="text"
              required
              value={customerPhone}
              onChange={(e) => setCustomerPhone(e.target.value)}
              className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-slate-100 focus:border-blue-500 focus:outline-none"
            />
          </div>
          <div>
            <label className="block text-slate-300 font-medium mb-1">Service Address *</label>
            <input
              type="text"
              required
              value={serviceAddress}
              onChange={(e) => setServiceAddress(e.target.value)}
              className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-slate-100 focus:border-blue-500 focus:outline-none"
            />
          </div>
          <div>
            <label className="block text-slate-300 font-medium mb-1">Major Landmark</label>
            <input
              type="text"
              placeholder="e.g. Near Central Mosque"
              value={landmark}
              onChange={(e) => setLandmark(e.target.value)}
              className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-slate-100 focus:border-blue-500 focus:outline-none"
            />
          </div>
        </div>

        {/* Row 3: Category & Departmental Routing Selector */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label className="block text-slate-300 font-medium mb-1">Complaint Category *</label>
            <select
              value={category}
              onChange={(e) => handleCategoryChange(e.target.value as TicketCategory)}
              className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-slate-100 focus:border-blue-500 focus:outline-none"
            >
              <option value="LINE_SNAP">Line Snap / Fallen Conductor (Feeder)</option>
              <option value="FALLEN_POLE">Fallen Electric Pole (Feeder)</option>
              <option value="TRANSFORMER_FAULT">Transformer Smoking / Fault (Substation)</option>
              <option value="FUSE_BLOWN">Drop-Out Fuse Blown (Substation)</option>
              <option value="OIL_LEAK">Transformer Oil Leak (Substation)</option>
              <option value="FAULTY_METER">Prepaid STS Fault / Tamper (Metering)</option>
              <option value="NEW_METER_SCHEDULE">New Meter Schedule Appointment (Metering)</option>
              <option value="WRONG_BILLING">Disputed Billing / Overcharge (Billing)</option>
              <option value="UNCREDITED_PAYMENT">Uncredited Vending Token (Vending IT)</option>
              <option value="SAFETY_HAZARD_EMERGENCY">SAFETY HAZARD EMERGENCY (HSE &lt;2H)</option>
            </select>
          </div>

          <div>
            <label className="block text-slate-300 font-medium mb-1">Sub-Category *</label>
            <input
              type="text"
              required
              value={subCategory}
              onChange={(e) => setSubCategory(e.target.value)}
              className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-slate-100 focus:border-blue-500 focus:outline-none"
            />
          </div>

          <div>
            <label className="block text-slate-300 font-medium mb-1">Target Technical Unit Head *</label>
            <select
              value={assignedUnitName}
              onChange={(e) => setAssignedUnitName(e.target.value)}
              className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-blue-300 font-semibold focus:border-blue-500 focus:outline-none"
            >
              <option value="Feeder Maintenance Unit Head (Services)">Feeder Maintenance Unit Head (Services)</option>
              <option value="Transformer Maintenance Unit Head (Services)">Transformer Maintenance Unit Head (Services)</option>
              <option value="Meter Maintenance & Installation Unit Head (Services)">Meter Maintenance Unit Head (Services)</option>
              <option value="Billing Unit Head (Commercial)">Billing Unit Head (Commercial)</option>
              <option value="STS Vending & Billing IT Support Desk (Commercial)">STS Vending &amp; IT Support (Commercial)</option>
            </select>
          </div>
        </div>

        {/* Description Field */}
        <div>
          <label className="block text-slate-300 font-medium mb-1">Detailed Technical Description &bull; Symptoms *</label>
          <textarea
            rows={3}
            required
            placeholder="Document customer narration, exact pole number, burnt jumper smell, sparks, or meter screen error code..."
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            className="w-full bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-slate-100 focus:border-blue-500 focus:outline-none"
          />
        </div>

        {/* Staff Attribution Footer Stamp */}
        <div className="p-3 rounded-lg bg-slate-950/70 border border-slate-800 flex items-center justify-between text-slate-400">
          <div className="flex items-center space-x-2">
            <ShieldAlert className="w-4 h-4 text-emerald-400" />
            <span>
              Attribution: Logged &amp; Routed by <strong>{staffName}</strong> ({staffId})
            </span>
          </div>
          <span className="text-[10px] text-slate-500 font-mono">Immutable NERC Audit Trail</span>
        </div>

        <button
          type="submit"
          disabled={isSubmitting}
          className={`w-full font-bold py-2.5 px-4 rounded-lg transition-colors flex items-center justify-center space-x-2 shadow-md ${
            isEmergency
              ? "bg-red-600 hover:bg-red-500 text-white"
              : "bg-blue-600 hover:bg-blue-500 text-white"
          }`}
        >
          {isSubmitting ? (
            <span>Routing to Unit Head...</span>
          ) : (
            <>
              <Send className="w-4 h-4" />
              <span>
                {isEmergency ? "DISPATCH PRIORITY EMERGENCY TICKET (<2H SLA)" : "Log Ticket & Route to Designated Unit Head"}
              </span>
            </>
          )}
        </button>
      </form>
    </div>
  );
}
