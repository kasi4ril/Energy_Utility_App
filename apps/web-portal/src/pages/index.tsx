/**
 * NESCO Enterprise Utility Platform - Public Home & Landing Page (/).
 * Showcases live grid telemetry, architectural pillars, and direct departmental gateways.
 */

import React, { useState, useEffect } from "react";
import Head from "next/head";
import Link from "next/link";
import { apiClient } from "../lib/api-client";
import { parseTokenClaims } from "../lib/auth";
import { AuthClaims } from "../types/auth";

export default function HomePage() {
  const [claims, setClaims] = useState<AuthClaims | null>(null);

  useEffect(() => {
    const token = apiClient.getToken();
    if (token) {
      setClaims(parseTokenClaims(token));
    }
  }, []);

  const departments = [
    {
      title: "Executive Command Center",
      href: "/executive-command",
      role: "MD & C-Suite",
      icon: "👑",
      color: "from-amber-500/20 to-amber-950/40 border-amber-500/30",
      accent: "text-amber-400",
      description: "Live C-Suite pulse, role-segregated digital sign-off queue, and 4 CBM predictive maintenance prognostics.",
    },
    {
      title: "Customer Care & KYC",
      href: "/portal/customer-care",
      role: "Customer Operations",
      icon: "📞",
      color: "from-blue-500/20 to-blue-950/40 border-blue-500/30",
      accent: "text-blue-400",
      description: "Anti-ghost onboarding, premises GPS lock, direct SLA ticket routing, and 2-step customer confirmation gate.",
    },
    {
      title: "Services & Field Maintenance",
      href: "/portal/services/feeder",
      role: "Maintenance Engineering",
      icon: "⚡",
      color: "from-emerald-500/20 to-emerald-950/40 border-emerald-500/30",
      accent: "text-emerald-400",
      description: "33kV/11kV feeder patrols, transformer overhauls, metering disconnections, and field proof verification.",
    },
    {
      title: "Generation Master Station",
      href: "/portal/generation-master",
      role: "Generation Division",
      icon: "🌊",
      color: "from-cyan-500/20 to-cyan-950/40 border-cyan-500/30",
      accent: "text-cyan-400",
      description: "Remote call collation across Plants 1–6 (Kurra Falls), hydrological/thermal telemetry, and workshop loop.",
    },
    {
      title: "Central Energy Dispatch",
      href: "/portal/energy-management",
      role: "Grid Operations",
      icon: "🎛️",
      color: "from-violet-500/20 to-violet-950/40 border-violet-500/30",
      accent: "text-violet-400",
      description: "Dual-source bulk MW intake (Hydro + TCN), priority feeder rationing (Bands A–E), and 5-Step Safety PTW lock.",
    },
    {
      title: "Stores & Logistics Desk",
      href: "/portal/stores",
      role: "Warehouse Operations",
      icon: "📦",
      color: "from-orange-500/20 to-orange-950/40 border-orange-500/30",
      accent: "text-orange-400",
      description: "Dual-Path material requisitions (Path A in-stock vs Path B MD procurement), waybills, and scrap returns.",
    },
    {
      title: "Internal Audit Desk",
      href: "/portal/internal-audit",
      role: "Audit & Compliance",
      icon: "⚖️",
      color: "from-rose-500/20 to-rose-950/40 border-rose-500/30",
      accent: "text-rose-400",
      description: "Automated 3-way staff claims verification (Duty Roster + Safety PTW + GPS match) and inventory audit.",
    },
    {
      title: "Billing & STS Vending",
      href: "/portal/billing-commercial",
      role: "Commercial Operations",
      icon: "💳",
      color: "from-teal-500/20 to-teal-950/40 border-teal-500/30",
      accent: "text-teal-400",
      description: "Zero-estimation metered billing, 20-digit STS tokens, NERC refund amortization, and theft anomaly AI.",
    },
  ];

  return (
    <>
      <Head>
        <title>NESCO Enterprise Utility Platform &bull; Sovereign Grid OS</title>
        <meta name="description" content="Mission-Critical Utility Workflow & Grid Intelligence for NESCO Nigeria Ltd." />
      </Head>

      <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans selection:bg-emerald-500/30">
        {/* Top Announcement Bar */}
        <div className="bg-emerald-950/80 border-b border-emerald-500/30 py-2 px-6 flex items-center justify-between text-xs">
          <div className="flex items-center space-x-2">
            <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
            <span className="font-semibold text-emerald-300">GRID FREQUENCY:</span>
            <span className="font-mono text-white">50.02 Hz (Nominal Stability)</span>
            <span className="text-slate-600">|</span>
            <span className="text-slate-400 hidden sm:inline">BULK INTAKE:</span>
            <span className="font-mono text-emerald-300 hidden sm:inline">63.6 MW Total</span>
          </div>
          <div className="flex items-center space-x-3">
            {claims ? (
              <span className="text-slate-300">
                Logged in as: <strong className="text-white">{claims.full_name}</strong> ({claims.primary_role})
              </span>
            ) : (
              <Link href="/login" className="text-emerald-400 hover:text-emerald-300 font-medium transition">
                Staff Authentication &rarr;
              </Link>
            )}
          </div>
        </div>

        {/* Global Navigation */}
        <nav className="h-16 border-b border-slate-800/80 px-6 flex items-center justify-between max-w-7xl w-full mx-auto">
          <div className="flex items-center space-x-3">
            <div className="p-2 rounded-lg bg-emerald-500/10 border border-emerald-500/30 text-emerald-400 font-bold text-lg">
              ⚡
            </div>
            <div>
              <span className="text-base font-bold text-white tracking-wider">NESCO</span>
              <span className="text-xs text-slate-400 ml-2 font-mono">UTILITY PLATFORM</span>
            </div>
          </div>

          <div className="flex items-center space-x-4">
            <Link
              href="/field-pwa/tasks"
              className="text-xs text-slate-400 hover:text-white transition hidden md:flex items-center space-x-1"
            >
              <span>📱</span>
              <span>Field PWA</span>
            </Link>
            <Link
              href="/login"
              className="px-4 py-2 rounded-lg bg-emerald-600 hover:bg-emerald-500 text-white font-semibold text-xs transition shadow-lg shadow-emerald-900/30"
            >
              {claims ? "Access My Portal" : "Enterprise Sign In"}
            </Link>
          </div>
        </nav>

        {/* Hero Section */}
        <header className="py-16 px-6 text-center max-w-4xl mx-auto space-y-6">
          <div className="inline-flex items-center space-x-2 px-3 py-1 rounded-full bg-slate-800 border border-slate-700 text-slate-300 text-xs">
            <span className="w-1.5 h-1.5 rounded-full bg-emerald-400" />
            <span>Air-Gapped Sovereign Utility Architecture &bull; 100% On-Premise</span>
          </div>

          <h1 className="text-4xl sm:text-5xl font-extrabold text-white tracking-tight leading-tight">
            Next-Generation Grid Intelligence &amp; Operational Governance
          </h1>

          <p className="text-base sm:text-lg text-slate-400 max-w-2xl mx-auto leading-relaxed">
            Unifying power generation, transmission balancing, electrical maintenance, commercial billing, and C-Suite command under a cryptographic SHA-256 hash-chain ledger.
          </p>

          <div className="pt-4 flex flex-wrap items-center justify-center gap-4">
            <Link
              href="/executive-command"
              className="px-6 py-3 rounded-xl bg-emerald-600 hover:bg-emerald-500 text-white font-bold text-sm transition shadow-xl shadow-emerald-950 flex items-center space-x-2"
            >
              <span>👑</span>
              <span>Launch Executive Command</span>
            </Link>
            <Link
              href="/login"
              className="px-6 py-3 rounded-xl bg-slate-800 hover:bg-slate-700 border border-slate-700 text-slate-200 font-semibold text-sm transition flex items-center space-x-2"
            >
              <span>⚡</span>
              <span>All Departmental Portals</span>
            </Link>
            <Link
              href="/field-pwa/tasks"
              className="px-6 py-3 rounded-xl bg-slate-900 hover:bg-slate-800 border border-slate-800 text-slate-400 hover:text-white font-medium text-sm transition flex items-center space-x-2"
            >
              <span>📱</span>
              <span>Technician PWA</span>
            </Link>
          </div>
        </header>

        {/* Live Grid Metrics Pulse */}
        <section className="max-w-7xl w-full mx-auto px-6 py-8">
          <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-4">
            <div className="p-4 rounded-xl bg-slate-900/90 border border-slate-800">
              <div className="text-[11px] font-medium text-slate-400">Hydro Generation</div>
              <div className="text-xl font-bold text-emerald-400 mt-1">48.6 MW</div>
              <div className="text-[10px] text-slate-500 mt-0.5">Kurra Falls Plants 1–6</div>
            </div>
            <div className="p-4 rounded-xl bg-slate-900/90 border border-slate-800">
              <div className="text-[11px] font-medium text-slate-400">TCN Bulk Intake</div>
              <div className="text-xl font-bold text-cyan-400 mt-1">15.0 MW</div>
              <div className="text-[10px] text-slate-500 mt-0.5">Transmission Hook</div>
            </div>
            <div className="p-4 rounded-xl bg-slate-900/90 border border-slate-800">
              <div className="text-[11px] font-medium text-slate-400">Total Bulk Power</div>
              <div className="text-xl font-bold text-white mt-1">63.6 MW</div>
              <div className="text-[10px] text-emerald-400 mt-0.5">+2.1% Availability</div>
            </div>
            <div className="p-4 rounded-xl bg-slate-900/90 border border-slate-800">
              <div className="text-[11px] font-medium text-slate-400">Feeders In-Service</div>
              <div className="text-xl font-bold text-white mt-1">28 / 32</div>
              <div className="text-[10px] text-slate-500 mt-0.5">1 on PTW &bull; 1 Rationing</div>
            </div>
            <div className="p-4 rounded-xl bg-slate-900/90 border border-slate-800">
              <div className="text-[11px] font-medium text-slate-400">Shift Collections</div>
              <div className="text-xl font-bold text-emerald-400 mt-1">₦4.85M</div>
              <div className="text-[10px] text-slate-500 mt-0.5">91.2% Prepaid Share</div>
            </div>
            <div className="p-4 rounded-xl bg-slate-900/90 border border-slate-800">
              <div className="text-[11px] font-medium text-slate-400">Active Integrity Alerts</div>
              <div className="text-xl font-bold text-amber-400 mt-1">3 Alerts</div>
              <div className="text-[10px] text-slate-500 mt-0.5">AI Theft &amp; Churn Risk</div>
            </div>
          </div>
        </section>

        {/* Departmental Gateways Grid */}
        <section className="max-w-7xl w-full mx-auto px-6 py-12 flex-1">
          <div className="mb-8">
            <h2 className="text-xl font-bold text-white tracking-wide">Operational Departmental Portals</h2>
            <p className="text-xs text-slate-400 mt-1">
              Select an operational post to enter its 4-Zone Enterprise UI shell with real-time queues and intake forms:
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-5">
            {departments.map((dept) => (
              <Link
                key={dept.title}
                href={dept.href}
                className={`p-5 rounded-xl bg-gradient-to-b ${dept.color} border hover:border-slate-500 transition-all duration-200 group flex flex-col justify-between`}
              >
                <div>
                  <div className="flex items-center justify-between mb-3">
                    <span className="text-2xl">{dept.icon}</span>
                    <span className="text-[10px] font-mono uppercase tracking-wider px-2 py-0.5 rounded bg-slate-900/80 border border-slate-700 text-slate-300">
                      {dept.role}
                    </span>
                  </div>
                  <h3 className={`text-sm font-bold text-white group-hover:${dept.accent} transition-colors mb-2`}>
                    {dept.title}
                  </h3>
                  <p className="text-xs text-slate-300/80 leading-relaxed">
                    {dept.description}
                  </p>
                </div>

                <div className="mt-5 pt-3 border-t border-slate-800/60 flex items-center justify-between text-xs font-semibold text-slate-400 group-hover:text-white transition-colors">
                  <span>Launch Portal</span>
                  <span className="group-hover:translate-x-1 transition-transform">&rarr;</span>
                </div>
              </Link>
            ))}
          </div>
        </section>

        {/* Architecture Highlights */}
        <section className="bg-slate-900/60 border-t border-slate-800 py-12 px-6">
          <div className="max-w-7xl mx-auto grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 text-xs">
            <div className="space-y-2">
              <div className="text-emerald-400 font-bold text-sm flex items-center space-x-1.5">
                <span>🤖</span>
                <span>Air-Gapped AI Sovereign Sidecar</span>
              </div>
              <p className="text-slate-400 leading-relaxed">
                Local Ollama AI Guardian executing Maker-Checker fraud scoring, photo GPS validation, and @GridAssistant in-app conversational co-pilot without external cloud transmission.
              </p>
            </div>

            <div className="space-y-2">
              <div className="text-emerald-400 font-bold text-sm flex items-center space-x-1.5">
                <span>⛓️</span>
                <span>Chained Cryptographic WorkNotes</span>
              </div>
              <p className="text-slate-400 leading-relaxed">
                Immutable SHA-256 prev-hash ledger protecting all internal handovers, work orders, audit findings, and executive sign-off decisions with strict append-only constraints.
              </p>
            </div>

            <div className="space-y-2">
              <div className="text-emerald-400 font-bold text-sm flex items-center space-x-1.5">
                <span>🔒</span>
                <span>5-Step Safety PTW Protocol</span>
              </div>
              <p className="text-slate-400 leading-relaxed">
                Permit-to-Work safety lock: Request &rarr; Dual Approval &rarr; Physical Isolation &rarr; Safety Clearance &rarr; Substation Restoration, protecting line crews on high-voltage circuits.
              </p>
            </div>

            <div className="space-y-2">
              <div className="text-emerald-400 font-bold text-sm flex items-center space-x-1.5">
                <span>📍</span>
                <span>PostGIS Spatial Topology</span>
              </div>
              <p className="text-slate-400 leading-relaxed">
                Rigorous electrical grid hierarchy (Plant &rarr; 33kV Substation &rarr; 11kV Feeder &rarr; Transformer DT &rarr; Customer Meter) with spherical Haversine proof verification.
              </p>
            </div>
          </div>
        </section>

        {/* Global Footer */}
        <footer className="border-t border-slate-800/80 py-6 px-6 text-center text-xs text-slate-500">
          <div className="max-w-7xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-3">
            <div>
              &copy; {new Date().getFullYear()} <strong>NESCO Nigeria Ltd</strong> &bull; Established 1923 &bull; Pioneer Hydroelectric Power Utility in Nigeria
            </div>
            <div className="flex items-center space-x-4">
              <span>Platform v1.0.0</span>
              <span>&bull;</span>
              <span>Air-Gapped Sovereign Deployment</span>
            </div>
          </div>
        </footer>
      </div>
    </>
  );
}
