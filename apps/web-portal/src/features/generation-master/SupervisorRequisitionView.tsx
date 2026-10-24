/**
 * Activity 4: Supervisor Requisition & Notes Desk.
 * Allows Master Station Operators and Plant Engineers to write official notes
 * and material/maintenance requisitions directly to immediate supervisors:
 * - Electrical Unit Head, Mechanical Unit Head, General Manager, Workshop Lead
 * - Detailed technical justification, symptoms, risk assessment
 * - Itemized materials/spares requisition builder
 * - Interactive tracking of supervisor acknowledgments and GM authorizations
 */

import React, { useState, useEffect } from "react";
import {
  HydroPlantSummary,
  PlantID,
  RequisitionSupervisor,
  RequisitionPriority,
  RequisitionType,
  RequisitionItem,
  SupervisorRequisition,
} from "./types";
import { fetchSupervisorRequisitions, submitSupervisorRequisition } from "./api";
import {
  FileText,
  Send,
  Plus,
  Trash2,
  CheckCircle,
  Clock,
  AlertTriangle,
  UserCheck,
  ShieldCheck,
  Package,
  Wrench,
  Search,
} from "lucide-react";

interface SupervisorRequisitionProps {
  plants: HydroPlantSummary[];
  staffId?: string;
  staffName?: string;
  onSuccess?: () => void;
}

export const SupervisorRequisitionView: React.FC<SupervisorRequisitionProps> = ({
  plants,
  staffId = "STAFF-GEN-OP-01",
  staffName = "Ibrahim Pam (Master Operator)",
  onSuccess,
}) => {
  const [requisitions, setRequisitions] = useState<SupervisorRequisition[]>([]);
  const [loading, setLoading] = useState(true);
  const [filterStatus, setFilterStatus] = useState<string>("ALL");

  // Form State
  const [plantId, setPlantId] = useState<PlantID>("PLANT_01_KURRA_FALLS");
  const [supervisorRole, setSupervisorRole] = useState<RequisitionSupervisor>("MECHANICAL_UNIT_HEAD");
  const [requisitionType, setRequisitionType] = useState<RequisitionType>("MATERIAL_SPARES");
  const [priority, setPriority] = useState<RequisitionPriority>("URGENT");
  const [title, setTitle] = useState("");
  const [notesToSupervisor, setNotesToSupervisor] = useState("");

  // Itemized Bill of Materials
  const [items, setItems] = useState<RequisitionItem[]>([
    { item_name: "", spec_or_part_no: "", quantity: 1, unit: "Pieces" },
  ]);

  const [submitting, setSubmitting] = useState(false);
  const [toastMessage, setToastMessage] = useState<string | null>(null);

  const getSupervisorDisplayName = (role: RequisitionSupervisor): string => {
    switch (role) {
      case "MECHANICAL_UNIT_HEAD":
        return "Engr. S. Lar (Mechanical Unit Head)";
      case "ELECTRICAL_UNIT_HEAD":
        return "Engr. D. Gyang (Electrical Unit Head)";
      case "GM_GENERATION":
        return "Engr. M. B. Jauro (General Manager - Generation)";
      case "WORKSHOP_LEAD":
        return "Mr. Bitrus Dalyop (Central Workshop Lead)";
      case "STORES_OFFICER":
        return "Store Officer (Central Stores Depot)";
      default:
        return "Immediate Supervisor";
    }
  };

  const loadData = async () => {
    setLoading(true);
    try {
      const records = await fetchSupervisorRequisitions();
      setRequisitions(records);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  const handleAddItem = () => {
    setItems([...items, { item_name: "", spec_or_part_no: "", quantity: 1, unit: "Pieces" }]);
  };

  const handleRemoveItem = (index: number) => {
    setItems(items.filter((_, i) => i !== index));
  };

  const handleItemChange = (index: number, field: keyof RequisitionItem, value: any) => {
    const updated = [...items];
    updated[index] = { ...updated[index], [field]: value };
    setItems(updated);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim() || !notesToSupervisor.trim()) {
      alert("Please enter a subject title and notes to your immediate supervisor.");
      return;
    }

    setSubmitting(true);
    setToastMessage(null);

    try {
      const validItems = items.filter((it) => it.item_name.trim().length > 0);
      const newReq = await submitSupervisorRequisition({
        plant_id: plantId,
        supervisor_role: supervisorRole,
        supervisor_name: getSupervisorDisplayName(supervisorRole),
        requisition_type: requisitionType,
        priority: priority,
        title: title.trim(),
        notes_to_supervisor: notesToSupervisor.trim(),
        items_requested: validItems.length > 0 ? validItems : undefined,
        author_staff_id: staffId,
        author_staff_name: staffName,
      });

      setRequisitions([newReq, ...requisitions]);
      setTitle("");
      setNotesToSupervisor("");
      setItems([{ item_name: "", spec_or_part_no: "", quantity: 1, unit: "Pieces" }]);
      setToastMessage(`Requisition ${newReq.requisition_no} successfully routed to ${newReq.supervisor_name}!`);
      if (onSuccess) onSuccess();
    } catch (err: any) {
      alert(err.message || "Failed to submit requisition");
    } finally {
      setSubmitting(false);
    }
  };

  const filteredRequisitions = filterStatus === "ALL"
    ? requisitions
    : requisitions.filter((r) => r.status === filterStatus);

  return (
    <div className="space-y-6 text-xs text-slate-100 max-w-5xl">
      {/* Header Banner */}
      <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-5 shadow-xl">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
          <div className="flex items-center space-x-3">
            <div className="p-2.5 rounded-lg bg-purple-500/20 text-purple-400 border border-purple-500/30">
              <FileText className="w-6 h-6" />
            </div>
            <div>
              <div className="flex items-center space-x-2">
                <span className="text-[10px] font-bold uppercase tracking-wider text-purple-400 bg-purple-500/10 px-2 py-0.5 rounded border border-purple-500/20">
                  Form S05-REQ &bull; Supervisor Communications &amp; Material Gate
                </span>
                <span className="text-[10px] font-mono text-slate-400">Hierarchy Routing Active</span>
              </div>
              <h1 className="text-lg sm:text-xl font-black text-white mt-1">
                Supervisor Requisition &amp; Technical Notes Section
              </h1>
              <p className="text-slate-400 text-xs mt-0.5">
                Official notes and spare parts requisitions submitted by the Master Station Operator to Unit Heads and the General Manager.
              </p>
            </div>
          </div>

          <div className="bg-slate-950/80 border border-slate-800 px-3.5 py-2 rounded-lg text-right shrink-0">
            <div className="text-[10px] text-slate-400">Submitting Staff:</div>
            <div className="font-bold text-slate-200 text-xs">{staffName}</div>
            <div className="text-[10px] font-mono text-purple-400">{staffId}</div>
          </div>
        </div>
      </div>

      {/* Success Notification */}
      {toastMessage && (
        <div className="p-4 rounded-xl bg-emerald-950/40 border border-emerald-500/60 shadow-lg flex items-start justify-between">
          <div className="flex items-center space-x-3">
            <CheckCircle className="w-5 h-5 text-emerald-400 shrink-0" />
            <div className="text-xs text-emerald-200 font-semibold">{toastMessage}</div>
          </div>
          <button
            onClick={() => setToastMessage(null)}
            className="text-xs text-emerald-400 hover:text-white font-bold ml-4"
          >
            Dismiss
          </button>
        </div>
      )}

      {/* Requisition Intake Form */}
      <form onSubmit={handleSubmit} className="bg-slate-900/90 border border-slate-800 rounded-xl p-5 shadow-lg space-y-4">
        <h2 className="text-xs font-bold text-purple-400 uppercase tracking-wider flex items-center space-x-2 pb-2 border-b border-slate-800">
          <UserCheck className="w-4 h-4" />
          <span>Write Official Note / Requisition to Immediate Supervisor</span>
        </h2>

        {/* Supervisor & Category Selector */}
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div>
            <label className="block text-[11px] font-semibold text-slate-300 mb-1">
              Target Immediate Supervisor *
            </label>
            <select
              value={supervisorRole}
              onChange={(e) => setSupervisorRole(e.target.value as RequisitionSupervisor)}
              className="w-full bg-slate-950 border border-slate-700 rounded-lg p-2.5 text-xs text-white focus:outline-none focus:border-purple-500"
            >
              <option value="MECHANICAL_UNIT_HEAD">Engr. S. Lar &bull; Mechanical Unit Head</option>
              <option value="ELECTRICAL_UNIT_HEAD">Engr. D. Gyang &bull; Electrical Unit Head</option>
              <option value="GM_GENERATION">Engr. M. B. Jauro &bull; General Manager (Gen)</option>
              <option value="WORKSHOP_LEAD">Mr. Bitrus Dalyop &bull; Central Workshop Lead</option>
              <option value="STORES_OFFICER">Store Officer &bull; Central Stores Depot</option>
            </select>
          </div>

          <div>
            <label className="block text-[11px] font-semibold text-slate-300 mb-1">
              Requisition / Note Type *
            </label>
            <select
              value={requisitionType}
              onChange={(e) => setRequisitionType(e.target.value as RequisitionType)}
              className="w-full bg-slate-950 border border-slate-700 rounded-lg p-2.5 text-xs text-white focus:outline-none focus:border-purple-500"
            >
              <option value="MATERIAL_SPARES">Material &amp; Spare Parts Requisition</option>
              <option value="CRITICAL_ALERT">Critical Machine Alert / Urgent Note</option>
              <option value="OUTAGE_REQUEST">Planned Outage &amp; Overhaul Request</option>
              <option value="WORKSHOP_ORDER">Workshop Machining / Coil Winding Order</option>
            </select>
          </div>

          <div>
            <label className="block text-[11px] font-semibold text-slate-300 mb-1">
              Priority Urgency *
            </label>
            <select
              value={priority}
              onChange={(e) => setPriority(e.target.value as RequisitionPriority)}
              className={`w-full border rounded-lg p-2.5 text-xs font-bold focus:outline-none ${
                priority === "EMERGENCY_OUTAGE"
                  ? "bg-slate-950 border-rose-500/80 text-rose-300"
                  : priority === "URGENT"
                  ? "bg-slate-950 border-amber-500/80 text-amber-300"
                  : "bg-slate-950 border-slate-700 text-slate-200"
              }`}
            >
              <option value="ROUTINE">Routine Maintenance</option>
              <option value="URGENT">Urgent (Vibration / Temp Warning)</option>
              <option value="EMERGENCY_OUTAGE">EMERGENCY OUTAGE (Generator Tripped)</option>
            </select>
          </div>
        </div>

        {/* Affected Plant & Subject Title */}
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 pt-1">
          <div>
            <label className="block text-[11px] font-semibold text-slate-300 mb-1">
              Affected Generating Station *
            </label>
            <select
              value={plantId}
              onChange={(e) => setPlantId(e.target.value as PlantID)}
              className="w-full bg-slate-950 border border-slate-700 rounded-lg p-2.5 text-xs text-white focus:outline-none focus:border-purple-500"
            >
              {plants.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.name}
                </option>
              ))}
            </select>
          </div>

          <div className="sm:col-span-2">
            <label className="block text-[11px] font-semibold text-slate-300 mb-1">
              Requisition Title / Subject *
            </label>
            <input
              type="text"
              required
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="e.g. Request for Emergency Turbine Runner Bushings & Synthetic Oil (Unit 2)"
              className="w-full bg-slate-950 border border-slate-700 rounded-lg p-2.5 text-xs text-white placeholder-slate-500 focus:outline-none focus:border-purple-500"
            />
          </div>
        </div>

        {/* Notes to Immediate Supervisor */}
        <div>
          <label className="block text-[11px] font-semibold text-slate-300 mb-1">
            Notes &amp; Technical Justification to Immediate Supervisor *
          </label>
          <p className="text-[10px] text-slate-400 mb-1.5">
            Describe the operational observation, vibration/thermal readings, risk to grid baseload, and the required supervisor action.
          </p>
          <textarea
            rows={4}
            required
            value={notesToSupervisor}
            onChange={(e) => setNotesToSupervisor(e.target.value)}
            placeholder="e.g. Unit 2 guide bearing vibration has climbed to 4.7 mm/s during afternoon peak. Stator winding temperature steady at 72°C. Requesting immediate approval for phosphor bronze bushings and workshop shaft alignment before forced outage occurs."
            className="w-full bg-slate-950 border border-slate-700 rounded-lg p-3 text-xs text-white placeholder-slate-500 focus:outline-none focus:border-purple-500 font-sans leading-relaxed"
          />
        </div>

        {/* Bill of Materials / Spares Requisition Builder */}
        <div className="p-3.5 bg-slate-950/70 border border-slate-800 rounded-lg space-y-3">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-2">
              <Package className="w-4 h-4 text-purple-400" />
              <span className="text-xs font-bold text-slate-200">
                Itemized Bill of Materials / Spares Requested (Optional)
              </span>
            </div>
            <button
              type="button"
              onClick={handleAddItem}
              className="px-2.5 py-1 rounded bg-slate-800 hover:bg-slate-700 border border-slate-700 text-purple-300 text-[10px] font-semibold flex items-center space-x-1"
            >
              <Plus className="w-3 h-3" />
              <span>Add Item</span>
            </button>
          </div>

          <div className="space-y-2">
            {items.map((item, idx) => (
              <div key={idx} className="grid grid-cols-12 gap-2 items-center">
                <div className="col-span-5">
                  <input
                    type="text"
                    placeholder="Item Description (e.g. Phosphor Bronze Bushing)"
                    value={item.item_name}
                    onChange={(e) => handleItemChange(idx, "item_name", e.target.value)}
                    className="w-full bg-slate-900 border border-slate-700 rounded p-2 text-xs text-white placeholder-slate-600"
                  />
                </div>
                <div className="col-span-4">
                  <input
                    type="text"
                    placeholder="Spec / Part No. (e.g. BS-1400-PB1)"
                    value={item.spec_or_part_no}
                    onChange={(e) => handleItemChange(idx, "spec_or_part_no", e.target.value)}
                    className="w-full bg-slate-900 border border-slate-700 rounded p-2 text-xs text-white placeholder-slate-600 font-mono"
                  />
                </div>
                <div className="col-span-1">
                  <input
                    type="number"
                    min="1"
                    value={item.quantity}
                    onChange={(e) => handleItemChange(idx, "quantity", parseInt(e.target.value) || 1)}
                    className="w-full bg-slate-900 border border-slate-700 rounded p-2 text-xs text-white text-center font-mono"
                  />
                </div>
                <div className="col-span-1">
                  <select
                    value={item.unit}
                    onChange={(e) => handleItemChange(idx, "unit", e.target.value)}
                    className="w-full bg-slate-900 border border-slate-700 rounded p-2 text-[10px] text-white"
                  >
                    <option value="Pieces">Pcs</option>
                    <option value="Kg">Kg</option>
                    <option value="Litres">Ltrs</option>
                    <option value="Sets">Sets</option>
                    <option value="Drums">Drums</option>
                  </select>
                </div>
                <div className="col-span-1 text-center">
                  {items.length > 1 && (
                    <button
                      type="button"
                      onClick={() => handleRemoveItem(idx)}
                      className="p-1 text-slate-500 hover:text-rose-400"
                    >
                      <Trash2 className="w-3.5 h-3.5" />
                    </button>
                  )}
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Submit */}
        <div className="flex items-center justify-between pt-2">
          <span className="text-[10px] text-slate-400 font-mono">
            Direct routing to: <strong className="text-purple-300">{getSupervisorDisplayName(supervisorRole)}</strong>
          </span>

          <button
            type="submit"
            disabled={submitting || !title.trim()}
            className="px-6 py-2.5 rounded-xl bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-500 hover:to-indigo-500 text-white font-bold text-xs flex items-center space-x-2 shadow-lg transition-all disabled:opacity-50"
          >
            <Send className="w-4 h-4" />
            <span>Submit Requisition to Supervisor</span>
          </button>
        </div>
      </form>

      {/* Submitted Requisitions & Status Tracker */}
      <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-5 shadow-lg space-y-4">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 pb-3 border-b border-slate-800">
          <div>
            <h3 className="text-sm font-bold text-white flex items-center space-x-2">
              <Clock className="w-4 h-4 text-purple-400" />
              <span>Supervisor Requisitions &amp; Notes Queue</span>
            </h3>
            <p className="text-[11px] text-slate-400 mt-0.5">
              Live tracking of submitted notes, supervisor acknowledgments, and General Manager authorizations.
            </p>
          </div>

          <div className="flex items-center space-x-2">
            <span className="text-[10px] text-slate-400">Filter:</span>
            <select
              value={filterStatus}
              onChange={(e) => setFilterStatus(e.target.value)}
              className="bg-slate-950 border border-slate-700 rounded-lg px-2.5 py-1 text-slate-200 text-xs focus:outline-none"
            >
              <option value="ALL">All Requisitions ({requisitions.length})</option>
              <option value="PENDING_SUPERVISOR_REVIEW">Pending Review</option>
              <option value="SUPERVISOR_ACKNOWLEDGED">Supervisor Acknowledged</option>
              <option value="GM_APPROVED">GM Approved</option>
            </select>
          </div>
        </div>

        <div className="space-y-3.5">
          {filteredRequisitions.map((req) => (
            <div
              key={req.id}
              className="p-4 rounded-xl bg-slate-950/80 border border-slate-800 space-y-3 transition"
            >
              {/* Header Bar */}
              <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
                <div className="flex items-center space-x-2">
                  <span className="font-mono font-bold text-purple-400 text-xs">
                    {req.requisition_no}
                  </span>
                  <span className="px-2 py-0.2 rounded text-[10px] font-bold uppercase bg-slate-800 text-slate-300">
                    {req.plant_id.replace(/_/g, " ")}
                  </span>
                  <span
                    className={`px-2 py-0.2 rounded text-[9px] font-extrabold uppercase ${
                      req.priority === "EMERGENCY_OUTAGE"
                        ? "bg-rose-500/20 text-rose-300 border border-rose-500/30"
                        : req.priority === "URGENT"
                        ? "bg-amber-500/20 text-amber-300 border border-amber-500/30"
                        : "bg-slate-700 text-slate-300"
                    }`}
                  >
                    {req.priority.replace(/_/g, " ")}
                  </span>
                </div>

                <span
                  className={`px-2.5 py-0.5 rounded-full text-[10px] font-bold font-mono ${
                    req.status === "GM_APPROVED"
                      ? "bg-emerald-500/20 text-emerald-300 border border-emerald-500/30"
                      : req.status === "SUPERVISOR_ACKNOWLEDGED"
                      ? "bg-blue-500/20 text-blue-300 border border-blue-500/30"
                      : "bg-amber-500/20 text-amber-300 border border-amber-500/30"
                  }`}
                >
                  {req.status.replace(/_/g, " ")}
                </span>
              </div>

              {/* Title & Supervisor Info */}
              <div>
                <h4 className="font-bold text-white text-xs sm:text-sm">{req.title}</h4>
                <div className="text-[10px] text-slate-400 mt-0.5">
                  Routed to: <strong className="text-slate-200">{req.supervisor_name}</strong> &bull; Submitted by: <span className="text-slate-300">{req.author_staff_name}</span>
                </div>
              </div>

              {/* Notes to Supervisor Content */}
              <div className="p-3 bg-slate-900/70 border border-slate-800/80 rounded-lg text-xs text-slate-200 leading-relaxed font-sans">
                <span className="text-[10px] text-slate-400 font-bold block mb-1 uppercase tracking-wider">
                  Notes to Immediate Supervisor:
                </span>
                {req.notes_to_supervisor}
              </div>

              {/* Requested Materials Table if present */}
              {req.items_requested && req.items_requested.length > 0 && (
                <div className="bg-slate-900/50 border border-slate-800/60 rounded-lg p-2.5 space-y-1.5">
                  <span className="text-[10px] font-bold text-slate-400 uppercase tracking-wider block">
                    Requested Materials / Spares ({req.items_requested.length}):
                  </span>
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 text-[11px]">
                    {req.items_requested.map((it, i) => (
                      <div key={i} className="flex items-center justify-between p-1.5 bg-slate-950/70 rounded border border-slate-800">
                        <span className="text-slate-200 font-medium">{it.item_name} ({it.spec_or_part_no})</span>
                        <span className="font-mono text-cyan-300 font-bold">{it.quantity} {it.unit}</span>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* Supervisor Feedback / Action Notes */}
              {req.supervisor_response && (
                <div className="p-3 bg-emerald-950/20 border border-emerald-500/30 rounded-lg text-xs text-emerald-200">
                  <span className="text-[10px] text-emerald-400 font-bold block mb-0.5 uppercase tracking-wider flex items-center space-x-1">
                    <CheckCircle className="w-3 h-3" />
                    <span>Supervisor Acknowledgment / GM Directive:</span>
                  </span>
                  {req.supervisor_response}
                </div>
              )}

              {/* Footer Timestamp */}
              <div className="flex items-center justify-between text-[10px] text-slate-500 font-mono pt-1 border-t border-slate-800/80">
                <span>Submitted at: {new Date(req.created_at).toLocaleString("en-GB", { timeZone: "Africa/Lagos" })} WAT</span>
                <span className="text-purple-400 font-semibold">{req.requisition_type.replace(/_/g, " ")}</span>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};
