/**
 * Executive Portals & All Departments Analytics Dashboards Launchpad Card.
 * Provides the Managing Director one-click access to all departmental landing pages
 * and specialized analytics dashboards across NESCO Sovereign Utility.
 * Strictly <= 350 lines.
 */

import React from "react";
import Link from "next/link";
import {
  BarChart3,
  Zap,
  Users,
  Activity,
  Wrench,
  CreditCard,
  ShieldCheck,
  Package,
  ExternalLink,
  ArrowUpRight,
} from "lucide-react";

interface DeptPortal {
  id: string;
  name: string;
  tagline: string;
  landingUrl: string;
  analyticsUrl: string;
  icon: React.ElementType;
  badge: string;
  badgeColor: string;
  accentColor: string;
}

const DEPARTMENT_PORTALS: DeptPortal[] = [
  {
    id: "central",
    name: "Central Company-Wide Analytics",
    tagline: "Executive Data Storytelling & Cross-Department Pulse",
    landingUrl: "/",
    analyticsUrl: "/portal/central-analytics",
    icon: BarChart3,
    badge: "Company-Wide",
    badgeColor: "bg-amber-500/20 text-amber-300 border-amber-500/30",
    accentColor: "border-amber-500/40 text-amber-400",
  },
  {
    id: "generation",
    name: "Generation Department (Hydro Fleet)",
    tagline: "Kurra, Kwall, Anoma, Jere MW output & machine health",
    landingUrl: "/portal/generation-master",
    analyticsUrl: "/portal/generation-master/analytics",
    icon: Zap,
    badge: "Baseload Hydro",
    badgeColor: "bg-cyan-500/20 text-cyan-300 border-cyan-500/30",
    accentColor: "border-cyan-500/40 text-cyan-400",
  },
  {
    id: "customercare",
    name: "Customer Care & Service Delivery",
    tagline: "KYC onboarding, smart triage, token vending & closures",
    landingUrl: "/portal/customer-care",
    analyticsUrl: "/portal/customer-care/analytics",
    icon: Users,
    badge: "CRM & Service",
    badgeColor: "bg-blue-500/20 text-blue-300 border-blue-500/30",
    accentColor: "border-blue-500/40 text-blue-400",
  },
  {
    id: "dispatch",
    name: "Energy Management & Dispatch",
    tagline: "33kV feeder balancing, load rationing & TCN grid hook",
    landingUrl: "/portal/energy-management",
    analyticsUrl: "/portal/energy-management",
    icon: Activity,
    badge: "Grid Dispatch",
    badgeColor: "bg-emerald-500/20 text-emerald-300 border-emerald-500/30",
    accentColor: "border-emerald-500/40 text-emerald-400",
  },
  {
    id: "services",
    name: "Engineering & Maintenance Services (Feeder Unit)",
    tagline: "33kV/11kV corridors, PTW outages, GPS proofs & line patrols",
    landingUrl: "/portal/services/feeder",
    analyticsUrl: "/portal/services/feeder/analytics",
    icon: Wrench,
    badge: "Feeder Unit",
    badgeColor: "bg-orange-500/20 text-orange-300 border-orange-500/30",
    accentColor: "border-orange-500/40 text-orange-400",
  },
  {
    id: "commercial",
    name: "Commercial Operations & Billing",
    tagline: "STS token vending, tariff classes & revenue collections",
    landingUrl: "/portal/billing-commercial",
    analyticsUrl: "/portal/billing-commercial",
    icon: CreditCard,
    badge: "Revenue Hub",
    badgeColor: "bg-purple-500/20 text-purple-300 border-purple-500/30",
    accentColor: "border-purple-500/40 text-purple-400",
  },
  {
    id: "audit",
    name: "Internal Audit & Compliance",
    tagline: "Discrepancy verification, inventory checks & claims",
    landingUrl: "/portal/internal-audit",
    analyticsUrl: "/portal/internal-audit",
    icon: ShieldCheck,
    badge: "Audit Gate",
    badgeColor: "bg-rose-500/20 text-rose-300 border-rose-500/30",
    accentColor: "border-rose-500/40 text-rose-400",
  },
  {
    id: "stores",
    name: "Central Stores & Procurement",
    tagline: "Stockpile radar, transformer spares & conductor stock",
    landingUrl: "/portal/stores",
    analyticsUrl: "/portal/stores",
    icon: Package,
    badge: "Stores & Stock",
    badgeColor: "bg-teal-500/20 text-teal-300 border-teal-500/30",
    accentColor: "border-teal-500/40 text-teal-400",
  },
];

export const AllDepartmentsAnalyticsPortalsCard: React.FC = () => {
  return (
    <div className="bg-slate-900 border border-slate-800 rounded-xl p-5 shadow-xl space-y-4">
      {/* Card Header */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between pb-3 border-b border-slate-800 gap-2">
        <div>
          <div className="flex items-center space-x-2">
            <BarChart3 className="w-4 h-4 text-amber-400" />
            <h2 className="text-sm font-bold text-white tracking-wide">
              Executive Portals &bull; All Departmental Analytics Dashboards
            </h2>
          </div>
          <p className="text-[11px] text-slate-400 mt-0.5">
            Direct executive access to live operational landing pages and specialized analytics dashboards across all 8 enterprise divisions.
          </p>
        </div>
        <span className="text-[10px] font-mono text-amber-400 bg-amber-500/10 px-2 py-0.5 rounded border border-amber-500/20 shrink-0">
          8 Live Portals Connected
        </span>
      </div>

      {/* Grid of Departmental Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3.5">
        {DEPARTMENT_PORTALS.map((dept) => {
          const Icon = dept.icon;
          return (
            <div
              key={dept.id}
              className="bg-slate-950/80 border border-slate-800 hover:border-slate-700 rounded-xl p-3.5 flex flex-col justify-between space-y-3 transition group"
            >
              <div>
                <div className="flex items-center justify-between mb-2">
                  <div className={`p-2 rounded-lg bg-slate-900 border ${dept.accentColor}`}>
                    <Icon className="w-4 h-4" />
                  </div>
                  <span className={`text-[9px] font-bold px-1.5 py-0.2 rounded border font-mono ${dept.badgeColor}`}>
                    {dept.badge}
                  </span>
                </div>
                <h3 className="font-bold text-white text-xs leading-snug group-hover:text-amber-300 transition-colors">
                  {dept.name}
                </h3>
                <p className="text-[10px] text-slate-400 mt-1 line-clamp-2 leading-relaxed">
                  {dept.tagline}
                </p>
              </div>

              {/* Action Buttons: Analytics Dashboard & Landing Page */}
              <div className="pt-2 border-t border-slate-800/80 space-y-1.5 font-mono text-[10px]">
                <Link
                  href={dept.analyticsUrl}
                  className="w-full py-1.5 px-2.5 rounded bg-slate-900 hover:bg-slate-800 text-cyan-300 hover:text-cyan-200 border border-cyan-500/30 flex items-center justify-between transition font-semibold"
                >
                  <span className="flex items-center space-x-1.5">
                    <BarChart3 className="w-3 h-3 text-cyan-400" />
                    <span>Analytics Dashboard</span>
                  </span>
                  <ArrowUpRight className="w-3 h-3 text-cyan-400" />
                </Link>

                <Link
                  href={dept.landingUrl}
                  className="w-full py-1 px-2.5 rounded bg-slate-950 hover:bg-slate-900 text-slate-400 hover:text-slate-200 border border-slate-800 flex items-center justify-between transition"
                >
                  <span className="flex items-center space-x-1.5">
                    <ExternalLink className="w-3 h-3 text-slate-500" />
                    <span>Department Landing</span>
                  </span>
                  <span className="text-slate-500">&rarr;</span>
                </Link>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};

