/**
 * Next.js Page: Services & Field Engineering Department Hub (/portal/services).
 * Unifies Feeder, Transformer, and Metering Engineering Units.
 */

import React, { useState, useEffect } from "react";
import Head from "next/head";
import Link from "next/link";
import { FeederHeaderBar } from "../../../features/services-maint/FeederHeaderBar";
import { GridAssistantDock } from "../../../components/shell/GridAssistantDock";
import { apiClient } from "../../../lib/api-client";
import { parseTokenClaims } from "../../../lib/auth";
import { AuthClaims } from "../../../types/auth";

export default function ServicesDepartmentHubPage() {
  const [claims, setClaims] = useState<AuthClaims | null>(null);
  const [isChatOpen, setIsChatOpen] = useState(false);

  useEffect(() => {
    const token = apiClient.getToken();
    if (token) setClaims(parseTokenClaims(token));
  }, []);

  const units = [
    {
      title: "Feeder Maintenance Unit",
      role: "High-Voltage Distribution Corridors",
      href: "/portal/services/feeder",
      icon: "⚡",
      color: "from-emerald-500/20 to-emerald-950/40 border-emerald-500/40",
      accent: "text-emerald-400",
      stat: "28 / 32 Feeders In-Service",
      statSub: "SAIDI: 6.4h &bull; SAIFI: 3.8 trips",
      description:
        "33kV & 11kV line patrols, Outage Request & Return Register with calendar search, ROW vegetation clearing, and anti-ghost GPS proof of work.",
    },
    {
      title: "Transformer Maintenance Unit",
      role: "Substation Distribution Transformers",
      href: "/portal/services/transformer",
      icon: "🛢️",
      color: "from-amber-500/20 to-amber-950/40 border-amber-500/40",
      accent: "text-amber-400",
      stat: "142 Active Substation DTs",
      statSub: "Oil Dielectric BDV > 40kV",
      description:
        "DGA oil dielectric testing, silica gel breather reactivation, bushing gasket seals, vacuum oil filtration, and IEEE C57 thermal standards.",
    },
    {
      title: "Metering Services Unit",
      role: "Revenue Metrology & Calibration",
      href: "/portal/services/metering",
      icon: "⏱️",
      badge: "STS-20",
      color: "from-cyan-500/20 to-cyan-950/40 border-cyan-500/40",
      accent: "text-cyan-400",
      stat: "14,850 Prepayment Meters",
      statSub: "94.8% TID Rollover Migrated",
      description:
        "STS Key Change Tokens (KCT1/KCT2), NEMSA-certified test bench recalibration, bypass tamper investigations, and holographic security seals.",
    },
  ];

  return (
    <>
      <Head>
        <title>NESCO Services &bull; Engineering Maintenance Hub</title>
      </Head>

      <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans selection:bg-emerald-500/30">
        {/* Thin Top Header Bar */}
        <FeederHeaderBar
          claims={claims}
          onToggleChat={() => setIsChatOpen(!isChatOpen)}
        />

        <main className="flex-1 p-6 max-w-7xl w-full mx-auto space-y-8">
          {/* Executive Overview Banner */}
          <div className="p-6 rounded-2xl bg-gradient-to-r from-emerald-950/60 via-slate-900/80 to-slate-950 border border-emerald-500/30 flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
            <div>
              <div className="inline-flex items-center space-x-2 px-3 py-1 rounded-full bg-emerald-950/80 border border-emerald-800/60 text-emerald-300 text-xs font-mono font-bold mb-2">
                <span>⚡</span>
                <span>Engineering Services &bull; Field Operations Command</span>
              </div>
              <h1 className="text-2xl font-black text-white tracking-tight">
                Engineering Services &amp; Maintenance Directorate
              </h1>
              <p className="text-xs text-slate-400 mt-1 max-w-2xl leading-relaxed">
                Supervising NESCO's distribution infrastructure: High-Voltage Feeder corridors, Substation Transformers, and Revenue Metering Population.
              </p>
            </div>

            <div className="flex items-center space-x-3">
              <Link
                href="/portal/central-analytics"
                className="px-4 py-2.5 rounded-xl bg-slate-900 hover:bg-slate-800 text-slate-200 border border-slate-700 font-bold text-xs transition"
              >
                Central Analytics &rarr;
              </Link>
            </div>
          </div>

          {/* 3 Core Engineering Units Launchpad */}
          <div>
            <div className="mb-4">
              <h2 className="text-sm font-bold uppercase tracking-wider text-slate-300 font-mono">
                Operational Engineering Units
              </h2>
              <p className="text-xs text-slate-400 mt-0.5">
                Select an engineering unit to enter its dedicated 5-workspace portal with persistent left navigation:
              </p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              {units.map((u) => (
                <Link
                  key={u.title}
                  href={u.href}
                  className={`p-6 rounded-2xl bg-gradient-to-b ${u.color} border hover:border-slate-400 transition-all duration-200 group flex flex-col justify-between space-y-6 shadow-xl`}
                >
                  <div className="space-y-4">
                    <div className="flex items-center justify-between">
                      <span className="text-3xl group-hover:scale-110 transition-transform">{u.icon}</span>
                      <span className="text-[10px] font-mono uppercase tracking-wider px-2 py-0.5 rounded bg-slate-950/80 border border-slate-700 text-slate-300">
                        {u.role}
                      </span>
                    </div>

                    <div>
                      <h3 className={`text-base font-bold text-white group-hover:${u.accent} transition-colors`}>
                        {u.title}
                      </h3>
                      <p className="text-xs text-slate-300 mt-1 leading-relaxed">
                        {u.description}
                      </p>
                    </div>

                    <div className="p-3 rounded-xl bg-slate-950/70 border border-slate-800/80">
                      <div className="text-xs font-bold text-white font-mono">{u.stat}</div>
                      <div className="text-[10px] text-slate-400 mt-0.5" dangerouslySetInnerHTML={{ __html: u.statSub }} />
                    </div>
                  </div>

                  <div className="pt-3 border-t border-slate-800/80 flex items-center justify-between text-xs font-bold text-slate-300 group-hover:text-white">
                    <span>Enter Unit Workspace</span>
                    <span className="group-hover:translate-x-1 transition-transform">&rarr;</span>
                  </div>
                </Link>
              ))}
            </div>
          </div>
        </main>

        {/* Omniscient AI Assistant Dock */}
        <GridAssistantDock
          isOpen={isChatOpen}
          onClose={() => setIsChatOpen(false)}
          staffRole={claims?.primary_role || "HEAD_ENGINEERING_SERVICES"}
        />
      </div>
    </>
  );
}

