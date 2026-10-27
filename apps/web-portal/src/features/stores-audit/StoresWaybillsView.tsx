/**
 * Goods Received & Supplier Waybills Terminal for Stores & Inventory.
 * Manages material intake, quality control inspection, and warehouse binning.
 */

import React, { useState } from "react";
import { MOCK_WAYBILLS } from "./storesMockData";
import { SupplierWaybillItem } from "./storesTypes";

export function StoresWaybillsView() {
  const [waybills, setWaybills] = useState<SupplierWaybillItem[]>(MOCK_WAYBILLS);
  const [supplierName, setSupplierName] = useState("");
  const [itemDescription, setItemDescription] = useState("");
  const [quantity, setQuantity] = useState("");
  const [poRef, setPoRef] = useState("");
  const [binAssigned, setBinAssigned] = useState("Bay 1");

  const handleAddWaybill = (e: React.FormEvent) => {
    e.preventDefault();
    const q = parseInt(quantity, 10);
    if (!supplierName || !itemDescription || isNaN(q) || q <= 0) {
      alert("Please provide valid consignment details.");
      return;
    }

    const newWb: SupplierWaybillItem = {
      id: `WB-${Date.now()}`,
      waybill_number: `WB-SUP-2026-${Math.floor(450 + Math.random() * 50)}`,
      supplier_name: supplierName,
      item_description: itemDescription,
      quantity_delivered: q,
      po_reference: poRef || "PO-NESCO-DIRECT",
      inspection_status: "PASSED_QC",
      store_bin_assigned: binAssigned,
      received_by_staff_id: "STAFF-STR-01",
      delivered_date_wat: "Just Now",
    };

    setWaybills([newWb, ...waybills]);
    setSupplierName("");
    setItemDescription("");
    setQuantity("");
    setPoRef("");
    alert("Goods received, QC verified, and bin location assigned successfully!");
  };

  return (
    <div className="space-y-6">
      {/* Intake Form */}
      <div className="p-5 rounded-xl bg-slate-900/80 border border-slate-800">
        <h3 className="text-sm font-bold text-white flex items-center space-x-2 mb-3">
          <span>🚚</span>
          <span>Goods Received Note (GRN) &amp; Supplier Delivery Ingestion</span>
        </h3>

        <form onSubmit={handleAddWaybill} className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 text-xs">
          <div>
            <label className="text-slate-400 font-semibold block mb-1">Supplier Name</label>
            <input
              type="text"
              value={supplierName}
              onChange={(e) => setSupplierName(e.target.value)}
              placeholder="e.g. Kabelmetal Nigeria PLC"
              className="w-full px-3 py-2 rounded-lg bg-[#0f172a] text-[#f8fafc] border border-slate-700 font-bold focus:border-orange-400 focus:outline-none"
              required
            />
          </div>

          <div>
            <label className="text-slate-400 font-semibold block mb-1">Purchase Order (PO) Reference</label>
            <input
              type="text"
              value={poRef}
              onChange={(e) => setPoRef(e.target.value)}
              placeholder="e.g. PO-NESCO-2026-109"
              className="w-full px-3 py-2 rounded-lg bg-[#0f172a] text-[#f8fafc] border border-slate-700 font-mono focus:border-orange-400 focus:outline-none"
            />
          </div>

          <div>
            <label className="text-slate-400 font-semibold block mb-1">Delivered Quantity</label>
            <input
              type="number"
              value={quantity}
              onChange={(e) => setQuantity(e.target.value)}
              placeholder="e.g. 500"
              className="w-full px-3 py-2 rounded-lg bg-[#0f172a] text-[#f8fafc] border border-slate-700 font-mono font-bold text-orange-400 focus:border-orange-400 focus:outline-none"
              required
            />
          </div>

          <div>
            <label className="text-slate-400 font-semibold block mb-1">Assigned Warehouse Bin</label>
            <select
              value={binAssigned}
              onChange={(e) => setBinAssigned(e.target.value)}
              className="w-full px-3 py-2 rounded-lg bg-[#0f172a] text-[#f8fafc] border border-slate-700 font-semibold focus:border-orange-400 focus:outline-none"
            >
              <option value="Bay 1 &bull; Bin A-04" className="bg-[#0f172a] text-[#f8fafc]">Bay 1 &bull; Insulator Racks</option>
              <option value="Bay 2 &bull; Rack C-01" className="bg-[#0f172a] text-[#f8fafc]">Bay 2 &bull; Cable & Conductor Drums</option>
              <option value="Bay 3 &bull; Drawer S-12" className="bg-[#0f172a] text-[#f8fafc]">Bay 3 &bull; Switchgear & Fuses</option>
              <option value="Chemical Shed &bull; Bay 4" className="bg-[#0f172a] text-[#f8fafc]">Chemical Shed &bull; Oil Drums</option>
            </select>
          </div>

          <div className="sm:col-span-2 lg:col-span-3">
            <label className="text-slate-400 font-semibold block mb-1">Material Description &amp; Specifications</label>
            <input
              type="text"
              value={itemDescription}
              onChange={(e) => setItemDescription(e.target.value)}
              placeholder="e.g. 150mm² ACSR 'Wolf' Aluminum Conductor (Brand New Drums, QC Inspected)"
              className="w-full px-3 py-2 rounded-lg bg-[#0f172a] text-[#f8fafc] border border-slate-700 focus:border-orange-400 focus:outline-none"
              required
            />
          </div>

          <div className="flex items-end">
            <button
              type="submit"
              className="w-full py-2 px-4 rounded-lg bg-orange-600 hover:bg-orange-500 text-white font-bold transition shadow-lg shadow-orange-950 flex items-center justify-center space-x-1"
            >
              <span>📥</span>
              <span>Accept &amp; Bin Stock</span>
            </button>
          </div>
        </form>
      </div>

      {/* Waybills Table */}
      <div className="rounded-xl border border-slate-800 overflow-hidden bg-slate-900/60">
        <div className="p-3.5 border-b border-slate-800 flex items-center justify-between">
          <span className="text-xs font-bold text-white uppercase tracking-wider font-mono">
            Supplier Waybills &amp; Consignment Intake Register
          </span>
          <span className="text-xs text-slate-400 font-mono">QC Inspected</span>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left text-xs">
            <thead className="bg-slate-900 text-slate-400 border-b border-slate-800 uppercase font-mono text-[10px]">
              <tr>
                <th className="py-3 px-4">Waybill No</th>
                <th className="py-3 px-3">Supplier &bull; PO</th>
                <th className="py-3 px-3">Material Description</th>
                <th className="py-3 px-3">Quantity</th>
                <th className="py-3 px-3">Store Bin</th>
                <th className="py-3 px-4 text-right">QC Status</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-800/80 text-slate-200">
              {waybills.map((w) => (
                <tr key={w.id} className="hover:bg-slate-800/40 transition">
                  <td className="py-3 px-4 font-mono font-bold text-orange-400">
                    {w.waybill_number}
                    <span className="block text-[10px] text-slate-500 font-normal">{w.delivered_date_wat}</span>
                  </td>
                  <td className="py-3 px-3">
                    <div className="font-bold text-white">{w.supplier_name}</div>
                    <div className="text-[10px] text-slate-500 font-mono">{w.po_reference}</div>
                  </td>
                  <td className="py-3 px-3 text-slate-300 max-w-xs">{w.item_description}</td>
                  <td className="py-3 px-3 font-mono font-bold text-white">
                    {w.quantity_delivered.toLocaleString()} Units
                  </td>
                  <td className="py-3 px-3 font-mono text-teal-300 text-[11px]" dangerouslySetInnerHTML={{ __html: w.store_bin_assigned }} />
                  <td className="py-3 px-4 text-right">
                    <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-emerald-950 text-emerald-400 border border-emerald-800">
                      PASSED QC
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

