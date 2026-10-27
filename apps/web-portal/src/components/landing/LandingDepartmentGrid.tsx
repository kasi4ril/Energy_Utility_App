/**
 * Operational Department Portals Launchpad Grid for NESCO Landing Page.
 * Showcases direct launchcards to all modernized departmental workspaces.
 */

import React from "react";
import Link from "next/link";

interface DepartmentCard {
  title: string;
  role: string;
  href: string;
  icon: string;
  badge: string;
  color: string;
  accent: string;
  description: string;
}

const DEPARTMENTS: DepartmentCard[] = [
  {
    title: "Executive Command Tier",
    role: "MD & CEO Command",
    href: "/portal/executive-command",
    icon: "👑",
    badge: "Tier 1",
    color: "from-amber-500/15 to-amber-950/30 border-amber-500/30",
    accent: "text-amber-400",
    description: "Requisitions desk with MD review, inter-departmental minuting, C-Suite capital analytics, and sovereign governance.",
  },
  {
    title: "Customer Care & KYC",
    role: "Customer Operations",
    href: "/portal/customer-care",
    icon: "📞",
    badge: "5 Workspaces",
    color: "from-blue-500/15 to-blue-950/30 border-blue-500/30",
    accent: "text-blue-400",
    description: "Period analytics, NIMC NIN national identity verification, smart landmark routing, and customer confirmation closure gate.",
  },
  {
    title: "Generation Master Station",
    role: "Hydrological Fleet",
    href: "/portal/generation-master",
    icon: "🌊",
    badge: "Plants 1–6",
    color: "from-cyan-500/15 to-cyan-950/30 border-cyan-500/30",
    accent: "text-cyan-400",
    description: "Real-time plant telemetry, active MW percentage share donut charts, hourly reading sheet, and cryptographic change-of-shift notes.",
  },
  {
    title: "Central Transmission Dispatch",
    role: "Grid Operations",
    href: "/portal/energy-management",
    icon: "🎛️",
    badge: "50.02 Hz",
    color: "from-violet-500/15 to-violet-950/30 border-violet-500/30",
    accent: "text-violet-400",
    description: "Dual-source bulk MW intake (Hydro + TCN), priority feeder load shedding (Bands A–E), and 5-step electrical safety PTW lock.",
  },
  {
    title: "Feeder Maintenance Unit",
    role: "33kV & 11kV Lines",
    href: "/portal/services/feeder",
    icon: "⚡",
    badge: "Outage Register",
    color: "from-emerald-500/15 to-emerald-950/30 border-emerald-500/30",
    accent: "text-emerald-400",
    description: "Reliability KPIs (SAIDI/SAIFI), calendar search outage register, permit exhaustion tracking, and GPS proof of work.",
  },
  {
    title: "Transformer Maintenance Unit",
    role: "Substation DTs",
    href: "/portal/services/transformer",
    icon: "🛢️",
    badge: "DGA Lab",
    color: "from-amber-600/15 to-amber-950/30 border-amber-600/30",
    accent: "text-amber-400",
    description: "DGA oil dielectric testing (BDV > 40kV), silica gel breathers, bushing gasket servicing, and IEEE C57 thermal standards.",
  },
  {
    title: "Metering Services Unit",
    role: "Revenue Metrology",
    href: "/portal/services/metering",
    icon: "⏱️",
    badge: "14.8K Fleet",
    color: "from-teal-500/15 to-teal-950/30 border-teal-500/30",
    accent: "text-teal-400",
    description: "STS Key Change Tokens (KCT1/KCT2), NEMSA-certified test bench recalibration, and bypass tamper investigations.",
  },
  {
    title: "Commercial & STS Vending",
    role: "Billing Operations",
    href: "/portal/billing-commercial",
    icon: "💳",
    badge: "STS-20",
    color: "from-teal-600/15 to-teal-950/30 border-teal-600/30",
    accent: "text-teal-400",
    description: "Zero-estimation metered billing, 20-digit STS tokens, NERC meter refund amortization, and revenue protection AI anomaly scans.",
  },
  {
    title: "Stores & Stockpile Desk",
    role: "Warehouse Logistics",
    href: "/portal/stores",
    icon: "📦",
    badge: "Dual-Path",
    color: "from-orange-500/15 to-orange-950/30 border-orange-500/30",
    accent: "text-orange-400",
    description: "Stockpile inventory ledgers, Dual-Path requisitions (In-Stock vs MD Approval), supplier waybills, and reverse logistics scrap.",
  },
  {
    title: "Internal Audit & Compliance",
    role: "Risk & Governance",
    href: "/portal/internal-audit",
    icon: "⚖️",
    badge: "SHA-256 Ledger",
    color: "from-rose-500/15 to-rose-950/30 border-rose-500/30",
    accent: "text-rose-400",
    description: "Cryptographic hash-chain inspection, 3-way staff claims audit (Roster + PTW + GPS), and salvage auction clearances.",
  },
];

export function LandingDepartmentGrid() {
  return (
    <section className="max-w-7xl w-full mx-auto px-6 py-10 select-none">
      <div className="flex flex-col sm:flex-row sm:items-end justify-between mb-8 gap-3">
        <div>
          <div className="text-[10px] uppercase font-bold text-cyan-400 font-mono tracking-widest">
            Sovereign Utility Workspaces
          </div>
          <h2 className="text-2xl font-black text-white tracking-tight mt-0.5">
            Operational Department Portals
          </h2>
          <p className="text-xs text-slate-400 mt-1 max-w-xl">
            Select an operational division to enter its 5-workspace layout standard with persistent left-hand navigation:
          </p>
        </div>

        <Link
          href="/portal/central-analytics"
          className="text-xs font-bold text-cyan-400 hover:text-cyan-300 transition flex items-center space-x-1"
        >
          <span>View Sovereign Central Analytics</span>
          <span>&rarr;</span>
        </Link>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-4">
        {DEPARTMENTS.map((dept) => (
          <Link
            key={dept.title}
            href={dept.href}
            className={`p-4 rounded-xl bg-gradient-to-b ${dept.color} border hover:border-slate-400 transition-all duration-200 group flex flex-col justify-between shadow-lg relative`}
          >
            <div>
              <div className="flex items-center justify-between mb-2.5">
                <span className="text-2xl group-hover:scale-110 transition-transform">{dept.icon}</span>
                <span className="text-[9px] font-mono uppercase tracking-wider px-2 py-0.5 rounded bg-slate-950/80 border border-slate-700 text-slate-300 font-bold">
                  {dept.badge}
                </span>
              </div>

              <div className="text-[10px] text-slate-400 font-mono uppercase tracking-wider">
                {dept.role}
              </div>

              <h3 className={`text-sm font-bold text-white group-hover:${dept.accent} transition-colors mt-0.5 mb-1.5`}>
                {dept.title}
              </h3>

              <p className="text-[11px] text-slate-300/80 leading-relaxed font-sans line-clamp-3">
                {dept.description}
              </p>
            </div>

            <div className="mt-4 pt-2.5 border-t border-slate-800/80 flex items-center justify-between text-[11px] font-bold text-slate-400 group-hover:text-white transition-colors">
              <span>Launch Portal</span>
              <span className="group-hover:translate-x-1 transition-transform">&rarr;</span>
            </div>
          </Link>
        ))}
      </div>
    </section>
  );
}

