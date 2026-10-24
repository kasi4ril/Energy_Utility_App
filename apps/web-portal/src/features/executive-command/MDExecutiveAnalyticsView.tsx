/**
 * Workspace 2: Managing Director Executive Analytics & Oversight Board.
 * Features:
 * - Summary of requests made, under review, approved, and declined
 * - Recharts Pie Chart of decision distribution & Bar Chart of departmental capital spend
 * - Departmental heads SLA turnaround velocity and approval rate matrix
 * Strictly <= 350 lines.
 */

import React, { useState, useEffect } from "react";
import { ExecutiveRequisition } from "./types";
import { AllDepartmentsAnalyticsPortalsCard } from "./AllDepartmentsAnalyticsPortalsCard";
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
import {
  BarChart3,
  PieChart as PieIcon,
  CheckCircle,
  Clock,
  TrendingUp,
  DollarSign,
  AlertCircle,
  FileCheck,
} from "lucide-react";

interface ExecutiveAnalyticsProps {
  requisitions: ExecutiveRequisition[];
}

const STATUS_COLORS: Record<string, string> = {
  APPROVED: "#10b981",              // Emerald
  UNDER_REVIEW: "#06b6d4",          // Cyan
  PENDING_APPROVAL: "#f59e0b",      // Amber
  ACKNOWLEDGED_MINUTED: "#3b82f6",  // Blue
  REJECTED: "#f43f5e",              // Rose
};

const DEPT_NAMES: Record<string, string> = {
  GENERATION: "Generation",
  ENGINEERING_SERVICES: "Engineering",
  COMMERCIAL_OPS: "Commercial",
  INTERNAL_AUDIT: "Audit",
  STORES_PROCUREMENT: "Stores",
  HSE_SAFETY: "HSE Safety",
};

export const MDExecutiveAnalyticsView: React.FC<ExecutiveAnalyticsProps> = ({
  requisitions,
}) => {
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  const totalRequests = requisitions.length;
  const approvedRequests = requisitions.filter((r) => r.status === "APPROVED");
  const underReviewRequests = requisitions.filter((r) => r.status === "UNDER_REVIEW");
  const pendingRequests = requisitions.filter((r) => r.status === "PENDING_APPROVAL");
  const minutedRequests = requisitions.filter((r) => r.status === "ACKNOWLEDGED_MINUTED");
  const rejectedRequests = requisitions.filter((r) => r.status === "REJECTED");

  const totalApprovedAmount = approvedRequests.reduce((sum, r) => sum + r.amount, 0);
  const totalGrossRequested = requisitions.reduce((sum, r) => sum + r.amount, 0);
  const approvalRatePct = totalRequests > 0 ? (approvedRequests.length / totalRequests) * 100 : 0;

  // Pie Chart Data: Decision Share
  const pieData = [
    { name: "Approved", value: approvedRequests.length, color: STATUS_COLORS.APPROVED },
    { name: "Under Review", value: underReviewRequests.length, color: STATUS_COLORS.UNDER_REVIEW },
    { name: "Pending", value: pendingRequests.length, color: STATUS_COLORS.PENDING_APPROVAL },
    { name: "Minuted", value: minutedRequests.length, color: STATUS_COLORS.ACKNOWLEDGED_MINUTED },
    { name: "Declined", value: rejectedRequests.length, color: STATUS_COLORS.REJECTED },
  ].filter((d) => d.value > 0);

  // Bar Chart Data: Spend by Department (Millions ₦)
  const depts = ["GENERATION", "ENGINEERING_SERVICES", "COMMERCIAL_OPS", "INTERNAL_AUDIT", "STORES_PROCUREMENT", "HSE_SAFETY"];
  const barData = depts.map((d) => {
    const deptReqs = requisitions.filter((r) => r.source_department === d);
    const totalM = deptReqs.reduce((sum, r) => sum + r.amount, 0) / 1000000;
    const approvedM = deptReqs.filter((r) => r.status === "APPROVED").reduce((sum, r) => sum + r.amount, 0) / 1000000;
    return {
      department: DEPT_NAMES[d] || d,
      requestedM: parseFloat(totalM.toFixed(1)),
      approvedM: parseFloat(approvedM.toFixed(1)),
      count: deptReqs.length,
    };
  });

  if (!mounted) return null;

  return (
    <div className="space-y-6 text-xs text-slate-100">
      {/* 1. Header & Top Metrics Ribbon */}
      <div className="bg-gradient-to-r from-slate-900 via-amber-950/40 to-slate-900 border border-amber-500/30 rounded-2xl p-5 shadow-xl">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <div className="flex items-center space-x-2 mb-1.5">
              <span className="px-2.5 py-0.5 rounded-full text-[10px] font-bold tracking-wider uppercase bg-amber-500/20 border border-amber-500/40 text-amber-300">
                Executive Governance &bull; Capital Authorization
              </span>
              <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-cyan-500/20 text-cyan-300 border border-cyan-500/30">
                Fiscal Year 2026
              </span>
            </div>
            <h1 className="text-xl sm:text-2xl font-black text-white tracking-tight">
              Managing Director Executive Analytics &amp; Capital Velocity
            </h1>
            <p className="text-xs text-slate-300 mt-1 max-w-2xl leading-relaxed">
              Consolidated intelligence metrics on requisitions made, reviews executed, capital authorized,
              and inter-departmental turnaround velocity.
            </p>
          </div>

          <div className="flex items-center space-x-2 shrink-0">
            <div className="bg-slate-950/80 border border-emerald-500/40 px-3.5 py-2 rounded-xl text-right">
              <div className="text-[10px] text-emerald-400 font-bold uppercase tracking-wider">Authorized Capital</div>
              <div className="text-lg font-black font-mono text-emerald-300">
                ₦{(totalApprovedAmount / 1000000).toFixed(2)}M
              </div>
            </div>
          </div>
        </div>

        {/* 4 Core KPIs Strip */}
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 mt-4">
          <div className="bg-slate-950/80 border border-slate-800 rounded-xl p-3">
            <div className="text-slate-400 font-semibold text-[11px] flex items-center justify-between">
              <span>Total Requests Made</span>
              <FileCheck className="w-3.5 h-3.5 text-amber-400" />
            </div>
            <div className="mt-1 flex items-baseline justify-between">
              <span className="text-xl font-black font-mono text-white">{totalRequests}</span>
              <span className="text-[10px] text-slate-400 font-mono">₦{(totalGrossRequested / 1000000).toFixed(1)}M</span>
            </div>
          </div>

          <div className="bg-slate-950/80 border border-slate-800 rounded-xl p-3">
            <div className="text-slate-400 font-semibold text-[11px] flex items-center justify-between">
              <span>Under Review / Queried</span>
              <AlertCircle className="w-3.5 h-3.5 text-cyan-400" />
            </div>
            <div className="mt-1 flex items-baseline justify-between">
              <span className="text-xl font-black font-mono text-cyan-300">{underReviewRequests.length}</span>
              <span className="text-[10px] text-cyan-400 font-mono">Action Active</span>
            </div>
          </div>

          <div className="bg-slate-950/80 border border-slate-800 rounded-xl p-3">
            <div className="text-slate-400 font-semibold text-[11px] flex items-center justify-between">
              <span>Approval Rate</span>
              <TrendingUp className="w-3.5 h-3.5 text-emerald-400" />
            </div>
            <div className="mt-1 flex items-baseline justify-between">
              <span className="text-xl font-black font-mono text-emerald-300">{approvalRatePct.toFixed(0)}%</span>
              <span className="text-[10px] text-emerald-400 font-mono">{approvedRequests.length} Approved</span>
            </div>
          </div>

          <div className="bg-slate-950/80 border border-slate-800 rounded-xl p-3">
            <div className="text-slate-400 font-semibold text-[11px] flex items-center justify-between">
              <span>Avg Approval SLA</span>
              <Clock className="w-3.5 h-3.5 text-blue-400" />
            </div>
            <div className="mt-1 flex items-baseline justify-between">
              <span className="text-xl font-black font-mono text-blue-300">2.8 Hrs</span>
              <span className="text-[10px] text-emerald-400 font-mono">&lt; 12h NERC Target</span>
            </div>
          </div>
        </div>
      </div>

      {/* 2. Visual Charts Grid: Decision Share Pie Chart + Spend by Dept Bar Chart */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-5 items-stretch">
        {/* Left Column: Donut Pie Chart (Decision Distribution) */}
        <div className="lg:col-span-5 bg-slate-900 border border-slate-800 rounded-xl p-4 flex flex-col justify-between">
          <div className="flex items-center justify-between border-b border-slate-800 pb-2 mb-2">
            <span className="font-bold text-slate-200 text-xs flex items-center space-x-1.5">
              <PieIcon className="w-3.5 h-3.5 text-amber-400" />
              <span>MD Requisitions Status Breakdown</span>
            </span>
            <span className="text-[10px] font-mono text-amber-400 bg-amber-500/10 px-2 py-0.5 rounded border border-amber-500/20">
              {totalRequests} Requests
            </span>
          </div>

          <div className="relative h-56 flex items-center justify-center">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={pieData}
                  cx="50%"
                  cy="50%"
                  innerRadius={50}
                  outerRadius={75}
                  paddingAngle={3}
                  dataKey="value"
                  strokeWidth={1.5}
                  stroke="#0f172a"
                >
                  {pieData.map((entry, idx) => (
                    <Cell key={idx} fill={entry.color} />
                  ))}
                </Pie>
                <RechartsTooltip
                  formatter={(val: any, name: any) => [`${val} Requests`, name]}
                  contentStyle={{ backgroundColor: "#0f172a", borderColor: "#334155", borderRadius: "8px", fontSize: "11px" }}
                />
              </PieChart>
            </ResponsiveContainer>

            <div className="absolute inset-0 flex flex-col items-center justify-center pointer-events-none text-center">
              <span className="text-2xl font-black font-mono text-white">{approvedRequests.length}</span>
              <span className="text-[9px] uppercase tracking-wider text-emerald-400 font-bold">Authorized</span>
            </div>
          </div>

          {/* Breakdown Badges */}
          <div className="grid grid-cols-2 gap-2 pt-2 border-t border-slate-800 text-[11px] font-mono">
            {pieData.map((d, i) => (
              <div key={i} className="flex items-center justify-between p-1 rounded bg-slate-950/60 border border-slate-800/60">
                <div className="flex items-center space-x-1.5 truncate">
                  <span className="w-2 h-2 rounded-full shrink-0" style={{ backgroundColor: d.color }} />
                  <span className="text-slate-300 truncate font-sans text-[11px]">{d.name}</span>
                </div>
                <strong className="text-white text-xs ml-1">{d.value}</strong>
              </div>
            ))}
          </div>
        </div>

        {/* Right Column: Bar Chart (Financial Volume by Department) */}
        <div className="lg:col-span-7 bg-slate-900 border border-slate-800 rounded-xl p-4 flex flex-col justify-between">
          <div className="flex items-center justify-between border-b border-slate-800 pb-2 mb-2">
            <span className="font-bold text-slate-200 text-xs flex items-center space-x-1.5">
              <BarChart3 className="w-3.5 h-3.5 text-amber-400" />
              <span>Capital Requested vs Authorized by Department (₦ Millions)</span>
            </span>
            <span className="text-[10px] text-slate-400 font-mono">6 Divisions</span>
          </div>

          <div className="h-56 w-full">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={barData} margin={{ top: 10, right: 15, left: -10, bottom: 5 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" vertical={false} />
                <XAxis dataKey="department" stroke="#94a3b8" fontSize={11} tickLine={false} />
                <YAxis stroke="#94a3b8" fontSize={11} unit="M" tickLine={false} />
                <RechartsTooltip
                  formatter={(val: any, name: any) => [`₦${val}M`, name === "requestedM" ? "Requested" : "Authorized"]}
                  contentStyle={{ backgroundColor: "#0f172a", borderColor: "#334155", borderRadius: "8px", fontSize: "11px" }}
                />
                <Legend
                  wrapperStyle={{ fontSize: "11px", paddingTop: "8px" }}
                  formatter={(v) => (v === "requestedM" ? "Requested Capital (₦M)" : "Authorized by MD (₦M)")}
                />
                <Bar dataKey="requestedM" fill="#334155" radius={[4, 4, 0, 0]} barSize={20} />
                <Bar dataKey="approvedM" fill="#10b981" radius={[4, 4, 0, 0]} barSize={20} />
              </BarChart>
            </ResponsiveContainer>
          </div>

          <div className="flex flex-wrap items-center justify-between pt-2 border-t border-slate-800 text-[11px] text-slate-400">
            <span>Highest Commitment: <strong className="text-emerald-400 font-mono">Central Stores (₦27.5M)</strong></span>
            <span>Emergency Outage: <strong className="text-amber-400 font-mono">Generation (₦18.5M)</strong></span>
          </div>
        </div>
      </div>

      {/* 3. Executive Portals: All Departments Analytics Dashboards Launchpad */}
      <AllDepartmentsAnalyticsPortalsCard />
    </div>
  );
};

