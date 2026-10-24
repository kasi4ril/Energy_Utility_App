/**
 * Active Plants Output Charts Card:
 * 1. Pie Chart displaying each active plant's output as a percentage of total generating power:
 *    - At the moment (Real-Time MW Output Share)
 *    - For the day / Last 24 Hours (Daily MWh Delivered Energy Share)
 * 2. Bar Chart representing the output of all active plants (MW Generated vs Installed Capacity).
 */

import React, { useState, useEffect } from "react";
import { HydroPlantSummary, GenerationAnalytics } from "./types";
import {
  PieChart,
  Pie,
  Cell,
  Tooltip as RechartsTooltip,
  ResponsiveContainer,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Legend,
} from "recharts";
import { Zap, PieChart as PieIcon, BarChart3, Clock, Calendar } from "lucide-react";

interface ChartsCardProps {
  plants: HydroPlantSummary[];
  analytics?: GenerationAnalytics | null;
}

const PLANT_COLORS: Record<string, string> = {
  PLANT_01_KURRA_FALLS: "#06b6d4", // Cyan
  PLANT_02_KURRA_FALLS: "#3b82f6", // Blue
  PLANT_03_KWALL: "#10b981",       // Emerald
  PLANT_04_KWALL: "#f59e0b",       // Amber
  PLANT_05_ANOMA: "#a855f7",       // Purple
  PLANT_06_JERE: "#f43f5e",        // Rose
};

const SHORT_NAMES: Record<string, string> = {
  PLANT_01_KURRA_FALLS: "Kurra 1",
  PLANT_02_KURRA_FALLS: "Kurra 2",
  PLANT_03_KWALL: "Kwall 1",
  PLANT_04_KWALL: "Kwall 2",
  PLANT_05_ANOMA: "Anoma",
  PLANT_06_JERE: "Jere",
};

export const ActivePlantsOutputChartsCard: React.FC<ChartsCardProps> = ({
  plants,
  analytics,
}) => {
  const [mounted, setMounted] = useState(false);
  const [timeframe, setTimeframe] = useState<"MOMENT" | "DAY">("MOMENT");

  useEffect(() => {
    setMounted(true);
  }, []);

  // Filter only active generating plants
  const activePlants = plants.filter(
    (p) => p.state === "IN_SERVICE" && p.current_mw > 0
  );

  const totalCurrentMW = activePlants.reduce((sum, p) => sum + p.current_mw, 0);
  const totalDailyMWh = analytics?.daily_mwh_delivered || (totalCurrentMW * 23.25);

  // Prepare Pie Chart data
  const pieData = activePlants.map((plant) => {
    const momentPct = totalCurrentMW > 0 ? (plant.current_mw / totalCurrentMW) * 100 : 0;
    
    // Estimate 24-hour daily delivered MWh based on plant output share
    const plantShare = analytics?.plants_generation_share?.find((s) => s.plant_id === plant.id);
    const dayPct = plantShare ? plantShare.pct : momentPct;
    const dailyMWh = (dayPct / 100) * totalDailyMWh;

    return {
      id: plant.id,
      name: plant.name,
      shortName: SHORT_NAMES[plant.id] || plant.name,
      value: timeframe === "MOMENT" ? parseFloat(plant.current_mw.toFixed(1)) : parseFloat(dailyMWh.toFixed(1)),
      percentage: timeframe === "MOMENT" ? momentPct : dayPct,
      unit: timeframe === "MOMENT" ? "MW" : "MWh",
      color: PLANT_COLORS[plant.id] || "#06b6d4",
    };
  });

  // Prepare Bar Chart data indicating all active plants
  const barData = activePlants.map((plant) => ({
    name: SHORT_NAMES[plant.id] || plant.name,
    fullName: plant.name,
    outputMW: plant.current_mw,
    capacityMW: plant.installed_capacity_mw,
    utilizationPct: plant.installed_capacity_mw > 0
      ? ((plant.current_mw / plant.installed_capacity_mw) * 100).toFixed(1)
      : "0",
    turbine: plant.turbine_type,
  }));

  if (!mounted) return null;

  return (
    <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-5 shadow-xl space-y-5">
      {/* Card Header with Timeframe Switcher */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between pb-3 border-b border-slate-800 gap-3">
        <div>
          <div className="flex items-center space-x-2">
            <Zap className="w-4 h-4 text-cyan-400" />
            <h2 className="text-sm font-bold text-white tracking-wide">
              Active Plants Output &bull; Power Distribution &amp; Fleet Comparison
            </h2>
          </div>
          <p className="text-[11px] text-slate-400 mt-0.5">
            Real-time percentage share of active generation at the moment and for the day (last 24 hours), with plant output bar comparison.
          </p>
        </div>

        {/* Timeframe Toggle Buttons */}
        <div className="flex items-center space-x-1.5 bg-slate-950 p-1 rounded-lg border border-slate-800 shrink-0">
          <button
            onClick={() => setTimeframe("MOMENT")}
            className={`px-3 py-1.5 rounded-md text-xs font-semibold flex items-center space-x-1.5 transition ${
              timeframe === "MOMENT"
                ? "bg-cyan-600 text-white shadow"
                : "text-slate-400 hover:text-slate-200"
            }`}
          >
            <Clock className="w-3.5 h-3.5" />
            <span>At the Moment (MW %)</span>
          </button>

          <button
            onClick={() => setTimeframe("DAY")}
            className={`px-3 py-1.5 rounded-md text-xs font-semibold flex items-center space-x-1.5 transition ${
              timeframe === "DAY"
                ? "bg-cyan-600 text-white shadow"
                : "text-slate-400 hover:text-slate-200"
            }`}
          >
            <Calendar className="w-3.5 h-3.5" />
            <span>Last 24 Hours (MWh %)</span>
          </button>
        </div>
      </div>

      {/* Visual Charts Grid: 1 Pie Chart + 1 Bar Chart */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-5 items-stretch">
        {/* Left Column: Pie Chart (Percentage of Total Generating Power) */}
        <div className="lg:col-span-5 bg-slate-950/80 border border-slate-800 rounded-xl p-4 flex flex-col justify-between">
          <div className="flex items-center justify-between border-b border-slate-800/80 pb-2 mb-2">
            <span className="font-bold text-slate-200 text-xs flex items-center space-x-1.5">
              <PieIcon className="w-3.5 h-3.5 text-cyan-400" />
              <span>
                {timeframe === "MOMENT"
                  ? "Current Power Share (% of Total MW)"
                  : "24h Delivered Energy Share (% of MWh)"}
              </span>
            </span>
            <span className="text-[10px] font-mono text-cyan-400 bg-cyan-500/10 px-2 py-0.5 rounded border border-cyan-500/20">
              {timeframe === "MOMENT" ? "Real-Time MW" : "Last 24h Total"}
            </span>
          </div>

          {/* Donut Pie Chart Container */}
          <div className="relative h-56 flex items-center justify-center">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={pieData}
                  cx="50%"
                  cy="50%"
                  innerRadius={55}
                  outerRadius={78}
                  paddingAngle={3}
                  dataKey="value"
                  strokeWidth={1.5}
                  stroke="#090d16"
                >
                  {pieData.map((entry) => (
                    <Cell key={entry.id} fill={entry.color} />
                  ))}
                </Pie>
                <RechartsTooltip
                  formatter={(val: any, _name: any, item: any) => [
                    `${val} ${item.payload.unit} (${item.payload.percentage.toFixed(1)}%)`,
                    item.payload.name,
                  ]}
                  contentStyle={{
                    backgroundColor: "#090d16",
                    borderColor: "#334155",
                    borderRadius: "8px",
                    fontSize: "11px",
                    color: "#f8fafc",
                  }}
                />
              </PieChart>
            </ResponsiveContainer>

            {/* Central Metric in Donut */}
            <div className="absolute inset-0 flex flex-col items-center justify-center pointer-events-none text-center">
              <span className="text-xl font-black font-mono text-white">
                {timeframe === "MOMENT"
                  ? `${totalCurrentMW.toFixed(1)} MW`
                  : `${totalDailyMWh.toFixed(0)} MWh`}
              </span>
              <span className="text-[9px] uppercase tracking-wider text-cyan-400 font-bold mt-0.5">
                {timeframe === "MOMENT" ? "Active Generation" : "24h Energy Delivered"}
              </span>
            </div>
          </div>

          {/* Pie Chart Legend Breakdown */}
          <div className="grid grid-cols-2 gap-2 pt-3 border-t border-slate-800/80 text-[11px] font-mono">
            {pieData.map((item) => (
              <div key={item.id} className="flex items-center justify-between p-1.5 rounded bg-slate-900/60 border border-slate-800/60">
                <div className="flex items-center space-x-1.5 truncate">
                  <span
                    className="w-2.5 h-2.5 rounded-full shrink-0"
                    style={{ backgroundColor: item.color }}
                  />
                  <span className="text-slate-300 truncate font-sans text-[11px]">{item.shortName}</span>
                </div>
                <div className="text-right shrink-0 ml-1">
                  <strong className="text-white text-xs">{item.percentage.toFixed(1)}%</strong>
                  <div className="text-[9px] text-slate-400 font-normal">{item.value} {item.unit}</div>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Right Column: Bar Chart Indicating All Active Plants Output */}
        <div className="lg:col-span-7 bg-slate-950/80 border border-slate-800 rounded-xl p-4 flex flex-col justify-between">
          <div className="flex items-center justify-between border-b border-slate-800/80 pb-2 mb-2">
            <span className="font-bold text-slate-200 text-xs flex items-center space-x-1.5">
              <BarChart3 className="w-3.5 h-3.5 text-cyan-400" />
              <span>Active Plants Output vs. Nameplate Installed Capacity</span>
            </span>
            <span className="text-[10px] text-slate-400 font-mono">
              {activePlants.length} Stations Online
            </span>
          </div>

          {/* Bar Chart Container */}
          <div className="h-56 w-full">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart
                data={barData}
                margin={{ top: 10, right: 15, left: -10, bottom: 5 }}
              >
                <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" vertical={false} />
                <XAxis
                  dataKey="name"
                  stroke="#94a3b8"
                  fontSize={11}
                  tickLine={false}
                />
                <YAxis
                  stroke="#94a3b8"
                  fontSize={11}
                  unit=" MW"
                  tickLine={false}
                />
                <RechartsTooltip
                  content={({ active, payload, label }) => {
                    if (active && payload && payload.length) {
                      const data = payload[0].payload;
                      return (
                        <div className="bg-slate-900 border border-slate-700 p-2.5 rounded-lg shadow-xl text-xs space-y-1">
                          <div className="font-bold text-white border-b border-slate-800 pb-1">
                            {data.fullName}
                          </div>
                          <div className="text-cyan-300 font-mono">
                            Power Output: <strong>{data.outputMW.toFixed(1)} MW</strong>
                          </div>
                          <div className="text-slate-400 font-mono">
                            Installed Limit: <strong>{data.capacityMW.toFixed(1)} MW</strong>
                          </div>
                          <div className="text-emerald-400 font-mono">
                            Utilization: <strong>{data.utilizationPct}%</strong>
                          </div>
                          <div className="text-slate-400 text-[10px]">
                            Turbine: {data.turbine}
                          </div>
                        </div>
                      );
                    }
                    return null;
                  }}
                />
                <Legend
                  wrapperStyle={{ fontSize: "11px", paddingTop: "8px" }}
                  formatter={(value) => (
                    <span className="text-slate-300 font-medium">
                      {value === "outputMW" ? "Current Active Output (MW)" : "Installed Capacity (MW)"}
                    </span>
                  )}
                />
                <Bar
                  dataKey="outputMW"
                  name="outputMW"
                  fill="#06b6d4"
                  radius={[4, 4, 0, 0]}
                  barSize={24}
                />
                <Bar
                  dataKey="capacityMW"
                  name="capacityMW"
                  fill="#334155"
                  radius={[4, 4, 0, 0]}
                  barSize={24}
                />
              </BarChart>
            </ResponsiveContainer>
          </div>

          {/* Quick Summary Strip */}
          <div className="flex flex-wrap items-center justify-between pt-3 border-t border-slate-800/80 text-[11px] text-slate-400">
            <div>
              Total Active Output: <strong className="text-cyan-300 font-mono">{totalCurrentMW.toFixed(1)} MW</strong>
            </div>
            <div>
              Avg Station Output: <strong className="text-slate-200 font-mono">{(totalCurrentMW / activePlants.length).toFixed(1)} MW</strong>
            </div>
            <div>
              Lead Generator: <strong className="text-emerald-400 font-mono">Kurra 1 ({activePlants[0]?.current_mw} MW)</strong>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

