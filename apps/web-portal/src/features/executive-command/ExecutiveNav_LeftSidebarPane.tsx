/**
 * Persistent Left-Hand Navigation Sidebar Pane for Executive Command Tier.
 * Positioned on the left side facing the screen, providing one-click access
 * across the MD's 5 core workspaces. Strictly <= 350 lines.
 */

import React from "react";
import { ExecutiveWorkspace } from "./types";
import {
  CheckCircle2,
  BarChart3,
  GitFork,
  Landmark,
  ShieldCheck,
  ChevronRight,
  ShieldAlert,
  ArrowLeft,
  Zap,
  Users,
  Activity,
  Wrench,
  CreditCard,
  Package,
  ExternalLink,
} from "lucide-react";
import Link from "next/link";

interface LeftSidebarProps {
  activeWorkspace: ExecutiveWorkspace;
  onSelectWorkspace: (ws: ExecutiveWorkspace) => void;
  pendingApprovalsCount: number;
  underReviewCount: number;
  interDeptActionCount?: number;
  emergencyAlertCount?: number;
}

interface NavItem {
  id: ExecutiveWorkspace;
  label: string;
  sublabel: string;
  icon: React.ElementType;
  badge?: number;
  badgeColor?: string;
  path: string;
}

export const ExecutiveNav_LeftSidebarPane: React.FC<LeftSidebarProps> = ({
  activeWorkspace,
  onSelectWorkspace,
  pendingApprovalsCount,
  underReviewCount,
  interDeptActionCount = 1,
  emergencyAlertCount = 2,
}) => {
  const navItems: NavItem[] = [
    {
      id: "approvals",
      label: "Approvals & Requisitions",
      sublabel: "Dept Heads Requests & Sign-Off",
      icon: CheckCircle2,
      badge: pendingApprovalsCount,
      badgeColor: "bg-amber-600 text-white animate-pulse",
      path: "/portal/executive-command/approvals",
    },
    {
      id: "analytics",
      label: "Executive Analytics Board",
      sublabel: "Requests, Spend & Decision Speed",
      icon: BarChart3,
      path: "/portal/executive-command/analytics",
    },
    {
      id: "inter-department",
      label: "Inter-Dept Information Hub",
      sublabel: "Cross-Unit Memos & MD Minuting",
      icon: GitFork,
      badge: interDeptActionCount,
      badgeColor: "bg-cyan-600 text-white",
      path: "/portal/executive-command/inter-department",
    },
    {
      id: "governance",
      label: "MD Professional Governance",
      sublabel: "NERC, Board & Concessions",
      icon: Landmark,
      path: "/portal/executive-command/governance",
    },
    {
      id: "ledger",
      label: "Decisions Audit Ledger",
      sublabel: "SHA-256 Chained Approval Stamping",
      icon: ShieldCheck,
      path: "/portal/executive-command/ledger",
    },
  ];

  return (
    <aside className="w-full md:w-64 bg-slate-900 border-r border-slate-800 flex flex-col shrink-0 h-full md:sticky md:top-14 md:h-[calc(100vh-3.5rem)] overflow-y-auto">
      {/* Sidebar Header */}
      <div className="p-3.5 border-b border-slate-800 bg-slate-950/60">
        <div className="text-[11px] font-bold tracking-wider text-amber-400 uppercase">
          MD Executive Command
        </div>
        <div className="text-xs text-slate-200 font-semibold mt-0.5">
          5-Tier Executive Governance
        </div>
      </div>

      {/* Navigation Links */}
      <nav className="p-2 space-y-1.5 flex-1">
        {navItems.map((item) => {
          const Icon = item.icon;
          const isActive = activeWorkspace === item.id;

          return (
            <button
              key={item.id}
              onClick={() => onSelectWorkspace(item.id)}
              className={`w-full flex items-center justify-between p-2.5 rounded-lg text-left transition-all group ${
                isActive
                  ? "bg-amber-500/15 border border-amber-500/40 text-amber-300 shadow-sm"
                  : "hover:bg-slate-800/80 text-slate-300 border border-transparent"
              }`}
            >
              <div className="flex items-start space-x-2.5 min-w-0">
                <div
                  className={`mt-0.5 p-1.5 rounded-md ${
                    isActive
                      ? "bg-amber-600 text-white"
                      : "bg-slate-800 text-slate-400 group-hover:text-slate-200"
                  }`}
                >
                  <Icon className="w-4 h-4" />
                </div>
                <div className="truncate">
                  <div
                    className={`text-xs font-semibold truncate ${
                      isActive ? "text-amber-300" : "text-slate-200"
                    }`}
                  >
                    {item.label}
                  </div>
                  <div className="text-[10px] text-slate-400 truncate">
                    {item.sublabel}
                  </div>
                </div>
              </div>

              <div className="flex items-center space-x-1 shrink-0 ml-2">
                {item.badge !== undefined && item.badge > 0 && (
                  <span
                    className={`text-[10px] font-bold px-1.5 py-0.5 rounded-full ${
                      item.badgeColor || "bg-slate-700 text-slate-200"
                    }`}
                  >
                    {item.badge}
                  </span>
                )}
                <ChevronRight
                  className={`w-3.5 h-3.5 transition-transform ${
                    isActive ? "text-amber-400 translate-x-0.5" : "text-slate-600"
                  }`}
                />
              </div>
            </button>
          );
        })}
      </nav>

      {/* Executive Portals: All Departments Landing Pages & Analytics Dashboards */}
      <div className="p-2 border-t border-slate-800/80 space-y-1">
        <div className="px-2 py-1 text-[10px] font-bold uppercase tracking-wider text-amber-400 flex items-center justify-between">
          <span>Executive Portals</span>
          <span className="text-[9px] font-mono text-slate-400 bg-slate-800 px-1.5 py-0.2 rounded">8 Dashboards</span>
        </div>

        {/* 1. General Landing Page & Central Analytics */}
        <div className="flex items-center justify-between p-1.5 rounded-md text-xs text-slate-300 hover:bg-slate-800/60 transition group">
          <Link href="/" className="flex items-center space-x-2 truncate hover:text-white" title="Main Utility Landing Page">
            <ArrowLeft className="w-3.5 h-3.5 text-amber-400 shrink-0" />
            <span className="truncate">General Landing</span>
          </Link>
          <Link
            href="/portal/central-analytics"
            className="text-[9px] font-bold font-mono px-1.5 py-0.5 rounded bg-amber-500/10 text-amber-400 border border-amber-500/30 hover:bg-amber-500/20 shrink-0"
            title="Company-Wide Central Analytics Dashboard"
          >
            Analytics
          </Link>
        </div>

        {/* 2. Generation Department (Hydro Fleet) */}
        <div className="flex items-center justify-between p-1.5 rounded-md text-xs text-slate-300 hover:bg-slate-800/60 transition group">
          <Link href="/portal/generation-master" className="flex items-center space-x-2 truncate hover:text-white" title="Generation Department Landing Page">
            <Zap className="w-3.5 h-3.5 text-cyan-400 shrink-0" />
            <span className="truncate">Generation (Hydro)</span>
          </Link>
          <Link
            href="/portal/generation-master/analytics"
            className="text-[9px] font-bold font-mono px-1.5 py-0.5 rounded bg-cyan-500/10 text-cyan-400 border border-cyan-500/30 hover:bg-cyan-500/20 shrink-0"
            title="Baseload Hydro Fleet Telemetry & MW Analytics Dashboard"
          >
            Analytics
          </Link>
        </div>

        {/* 3. Customer Care & Service Delivery */}
        <div className="flex items-center justify-between p-1.5 rounded-md text-xs text-slate-300 hover:bg-slate-800/60 transition group">
          <Link href="/portal/customer-care" className="flex items-center space-x-2 truncate hover:text-white" title="Customer Care Department Landing Page">
            <Users className="w-3.5 h-3.5 text-blue-400 shrink-0" />
            <span className="truncate">Customer Care</span>
          </Link>
          <Link
            href="/portal/customer-care/analytics"
            className="text-[9px] font-bold font-mono px-1.5 py-0.5 rounded bg-blue-500/10 text-blue-400 border border-blue-500/30 hover:bg-blue-500/20 shrink-0"
            title="Customer Care KPI Velocity & Ollama Analytics Dashboard"
          >
            Analytics
          </Link>
        </div>

        {/* 4. Energy Management & Central Dispatch (33kV) */}
        <div className="flex items-center justify-between p-1.5 rounded-md text-xs text-slate-300 hover:bg-slate-800/60 transition group">
          <Link href="/portal/energy-management" className="flex items-center space-x-2 truncate hover:text-white" title="Energy Dispatch & Grid Dashboard">
            <Activity className="w-3.5 h-3.5 text-emerald-400 shrink-0" />
            <span className="truncate">Energy Dispatch (33kV)</span>
          </Link>
          <Link
            href="/portal/energy-management"
            className="text-[9px] font-bold font-mono px-1.5 py-0.5 rounded bg-emerald-500/10 text-emerald-400 border border-emerald-500/30 hover:bg-emerald-500/20 shrink-0"
            title="33kV Grid Feeder Balancing & PTW Dashboard"
          >
            Analytics
          </Link>
        </div>

        {/* 5. Engineering & Maintenance Services */}
        <div className="flex items-center justify-between p-1.5 rounded-md text-xs text-slate-300 hover:bg-slate-800/60 transition group">
          <Link href="/portal/services" className="flex items-center space-x-2 truncate hover:text-white" title="Engineering & Distribution Maintenance Portal">
            <Wrench className="w-3.5 h-3.5 text-orange-400 shrink-0" />
            <span className="truncate">Engineering Services</span>
          </Link>
          <Link
            href="/portal/services"
            className="text-[9px] font-bold font-mono px-1.5 py-0.5 rounded bg-orange-500/10 text-orange-400 border border-orange-500/30 hover:bg-orange-500/20 shrink-0"
            title="Feeders, Transformers & Metering Work Orders Dashboard"
          >
            Analytics
          </Link>
        </div>

        {/* 6. Commercial Operations & Billing */}
        <div className="flex items-center justify-between p-1.5 rounded-md text-xs text-slate-300 hover:bg-slate-800/60 transition group">
          <Link href="/portal/billing-commercial" className="flex items-center space-x-2 truncate hover:text-white" title="Commercial Operations Landing Page">
            <CreditCard className="w-3.5 h-3.5 text-purple-400 shrink-0" />
            <span className="truncate">Commercial &amp; Billing</span>
          </Link>
          <Link
            href="/portal/billing-commercial"
            className="text-[9px] font-bold font-mono px-1.5 py-0.5 rounded bg-purple-500/10 text-purple-400 border border-purple-500/30 hover:bg-purple-500/20 shrink-0"
            title="STS Token Vending, Meter Fleet & Revenue Analytics Dashboard"
          >
            Analytics
          </Link>
        </div>

        {/* 7. Internal Audit & Compliance */}
        <div className="flex items-center justify-between p-1.5 rounded-md text-xs text-slate-300 hover:bg-slate-800/60 transition group">
          <Link href="/portal/internal-audit" className="flex items-center space-x-2 truncate hover:text-white" title="Internal Audit & Governance Console">
            <ShieldCheck className="w-3.5 h-3.5 text-rose-400 shrink-0" />
            <span className="truncate">Internal Audit</span>
          </Link>
          <Link
            href="/portal/internal-audit"
            className="text-[9px] font-bold font-mono px-1.5 py-0.5 rounded bg-rose-500/10 text-rose-400 border border-rose-500/30 hover:bg-rose-500/20 shrink-0"
            title="Fraud Detection, Scrap Clearance & Compliance Dashboard"
          >
            Analytics
          </Link>
        </div>

        {/* 8. Central Stores & Procurement */}
        <div className="flex items-center justify-between p-1.5 rounded-md text-xs text-slate-300 hover:bg-slate-800/60 transition group">
          <Link href="/portal/stores" className="flex items-center space-x-2 truncate hover:text-white" title="Central Stores & Inventory Landing Page">
            <Package className="w-3.5 h-3.5 text-teal-400 shrink-0" />
            <span className="truncate">Stores &amp; Inventory</span>
          </Link>
          <Link
            href="/portal/stores"
            className="text-[9px] font-bold font-mono px-1.5 py-0.5 rounded bg-teal-500/10 text-teal-400 border border-teal-500/30 hover:bg-teal-500/20 shrink-0"
            title="Critical Spares Stockpile Radar & Procurement Dashboard"
          >
            Analytics
          </Link>
        </div>
      </div>

      {/* Emergency Active Radar Pill at bottom of sidebar */}
      {emergencyAlertCount > 0 && (
        <div className="p-3 m-2 rounded-lg bg-rose-950/40 border border-rose-500/40 text-rose-300 text-xs">
          <div className="flex items-center space-x-2 font-bold text-rose-400 mb-1">
            <ShieldAlert className="w-4 h-4 animate-bounce text-rose-400" />
            <span>Critical Grid Outage Radar</span>
          </div>
          <p className="text-[10px] text-rose-200/80 leading-relaxed">
            Kurra Falls 2 tripped on runner cavitation pitting. Rayfield 33kV high-risk tree clearance active under emergency SLA.
          </p>
        </div>
      )}

      {/* Regulatory Badge */}
      <div className="p-3 border-t border-slate-800 text-[10px] text-slate-400 text-center font-mono bg-slate-950/40">
        NERC Sovereign Utility Governance &bull; v2026.4
      </div>
    </aside>
  );
};

