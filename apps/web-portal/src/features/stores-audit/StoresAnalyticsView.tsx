/**
 * Stockpile & Material Analytics Dashboard for Stores & Inventory.
 * Visualizes warehouse inventory valuation, reorder thresholds, and bin locations.
 */

import React from "react";
import {
  PieChart,
  Pie,
  Cell,
  ResponsiveContainer,
  Tooltip as RechartsTooltip,
} from "recharts";
import { MOCK_STOCKPILE } from "./storesMockData";

const CATEGORY_VALUATION = [
  { name: "Conductors & Cables", value: 18900000, color: "#f97316" },
  { name: "Insulators & Hardware", value: 2220000, color: "#06b6d4" },
  { name: "Transformer Diala Oil", value: 6720000, color: "#eab308" },
  { name: "Switchgear & Fuses", value: 1440000, color: "#10b981" },
  { name: "Safety PPE & Tools", value: 1260000, color: "#8b5cf6" },
];

export function StoresAnalyticsView() {
  const totalValuation = CATEGORY_VALUATION.reduce((acc, c) => acc + c.value, 0);

  return (
    <div className="space-y-6">
      {/* Top Inventory KPI Ribbon */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
          <div className="text-[11px] font-semibold text-slate-400">Total Stockpile Value</div>
          <div className="text-2xl font-black text-amber-400 mt-1">
            &#8358;{(totalValuation / 1000000).toFixed(2)}M
          </div>
          <div className="text-[10px] text-slate-500 mt-0.5">5 Bounded Asset Classes</div>
        </div>

        <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
          <div className="text-[11px] font-semibold text-slate-400">Low Stock SKUs</div>
          <div className="text-2xl font-black text-rose-400 mt-1">2 Items</div>
          <div className="text-[10px] text-rose-400 mt-0.5">Below Minimum Buffer Level</div>
        </div>

        <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
          <div className="text-[11px] font-semibold text-slate-400">Stock Turnover Velocity</div>
          <div className="text-2xl font-black text-emerald-400 mt-1">18.4 Days</div>
          <div className="text-[10px] text-emerald-400 mt-0.5">Rapid Field Consumption</div>
        </div>

        <div className="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
          <div className="text-[11px] font-semibold text-slate-400">Salvage Yard Value</div>
          <div className="text-2xl font-black text-teal-400 mt-1">&#8358;2.98M</div>
          <div className="text-[10px] text-slate-400 mt-0.5">Pending Auction Audit</div>
        </div>
      </div>

      {/* Category Breakdown & Stockpile Table */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
        {/* Category Valuation Pie Chart */}
        <div className="lg:col-span-5 p-5 rounded-xl bg-slate-900/70 border border-slate-800 flex flex-col justify-between">
          <div>
            <h3 className="text-sm font-bold text-white">Stockpile Value by Asset Category</h3>
            <p className="text-xs text-slate-400 mt-0.5">
              Capital distribution across conductors, oils, and switchgear.
            </p>
          </div>

          <div className="h-60 my-2">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={CATEGORY_VALUATION}
                  cx="50%"
                  cy="50%"
                  innerRadius={55}
                  outerRadius={85}
                  paddingAngle={3}
                  dataKey="value"
                >
                  {CATEGORY_VALUATION.map((entry, idx) => (
                    <Cell key={`cell-${idx}`} fill={entry.color} />
                  ))}
                </Pie>
                <RechartsTooltip
                  content={({ active, payload }) => {
                    if (active && payload && payload.length) {
                      const data = payload[0].payload;
                      const pct = ((data.value / totalValuation) * 100).toFixed(1);
                      return (
                        <div className="p-2 rounded bg-slate-950 border border-slate-700 text-xs shadow-xl">
                          <div className="font-bold text-white">{data.name}</div>
                          <div className="text-orange-400 font-mono">
                            &#8358;{(data.value / 1000000).toFixed(2)}M ({pct}%)
                          </div>
                        </div>
                      );
                    }
                    return null;
                  }}
                />
              </PieChart>
            </ResponsiveContainer>
          </div>

          <div className="space-y-1.5 text-[11px] border-t border-slate-800/80 pt-3">
            {CATEGORY_VALUATION.map((c) => (
              <div key={c.name} className="flex items-center justify-between">
                <span className="flex items-center space-x-2 truncate">
                  <span className="w-2 h-2 rounded-full flex-shrink-0" style={{ backgroundColor: c.color }} />
                  <span className="text-slate-300 truncate">{c.name}</span>
                </span>
                <span className="font-mono font-bold text-white ml-2">
                  &#8358;{(c.value / 1000000).toFixed(2)}M
                </span>
              </div>
            ))}
          </div>
        </div>

        {/* Stock Items Inventory Table */}
        <div className="lg:col-span-7 rounded-xl border border-slate-800 overflow-hidden bg-slate-900/60 flex flex-col justify-between">
          <div className="p-3.5 border-b border-slate-800 flex items-center justify-between">
            <span className="text-xs font-bold text-white uppercase tracking-wider font-mono">
              Key Warehouse Inventory SKUs
            </span>
            <span className="text-xs text-slate-400 font-mono">Physical Ledger</span>
          </div>

          <div className="overflow-x-auto flex-1">
            <table className="w-full text-left text-xs">
              <thead className="bg-slate-900 text-slate-400 border-b border-slate-800 uppercase font-mono text-[10px]">
                <tr>
                  <th className="py-2.5 px-3">Item &bull; SKU</th>
                  <th className="py-2.5 px-3">Bin Location</th>
                  <th className="py-2.5 px-3">Stock / Min</th>
                  <th className="py-2.5 px-3 text-right">Status</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-800/80 text-slate-200">
                {MOCK_STOCKPILE.map((item) => (
                  <tr key={item.id} className="hover:bg-slate-800/40 transition">
                    <td className="py-2.5 px-3">
                      <div className="font-bold text-white truncate max-w-xs">{item.item_name}</div>
                      <div className="text-[10px] text-slate-500 font-mono">{item.sku}</div>
                    </td>
                    <td className="py-2.5 px-3 font-mono text-slate-300 text-[11px]">
                      {item.bin_location}
                    </td>
                    <td className="py-2.5 px-3 font-mono">
                      <span className="font-bold text-white">
                        {item.current_quantity.toLocaleString()} {item.unit_of_measure}
                      </span>
                      <span className="block text-[10px] text-slate-500">
                        Min: {item.reorder_threshold.toLocaleString()}
                      </span>
                    </td>
                    <td className="py-2.5 px-3 text-right">
                      {item.reorder_status === "HEALTHY" && (
                        <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-emerald-950 text-emerald-400 border border-emerald-800">
                          HEALTHY
                        </span>
                      )}
                      {item.reorder_status === "REORDER_TRIGGERED" && (
                        <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-rose-950 text-rose-400 border border-rose-800">
                          REORDER DUE
                        </span>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <div className="p-3 border-t border-slate-800/80 text-[11px] text-slate-400 flex items-center justify-between">
            <span>Automated replenishment triggers active.</span>
            <span className="font-mono text-orange-400 font-semibold">Stores Audit Verified</span>
          </div>
        </div>
      </div>
    </div>
  );
}

