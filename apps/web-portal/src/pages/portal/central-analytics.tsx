/**
 * Central Daily Analytics Dashboard (/portal/central-analytics).
 * Unified Sovereign Data Storytelling for the Day across All NESCO Departments.
 * Generation, Dispatch, Maintenance, Customer Care, Commercial, and Stores.
 * Strictly <= 350 lines.
 */

import React, { useState, useEffect } from "react";
import Head from "next/head";
import Link from "next/link";
import { TopCommandBar } from "../../components/shell/TopCommandBar";
import { GridAssistantDock } from "../../components/shell/GridAssistantDock";
import { parseTokenClaims } from "../../lib/auth";
import { AuthClaims } from "../../types/auth";
import { apiClient } from "../../lib/api-client";
import {
  Sparkles,
  TrendingUp,
  Zap,
  Activity,
  Users,
  ShieldCheck,
  Package,
  DollarSign,
  ArrowRight,
  CheckCircle2,
  Clock,
  RefreshCw,
} from "lucide-react";

export default function CentralAnalyticsPage() {
  const [claims, setClaims] = useState<AuthClaims | null>(null);
  const [isChatOpen, setIsChatOpen] = useState(false);
  const [selectedDept, setSelectedDept] = useState<string>("ALL");
  const [lastUpdated, setLastUpdated] = useState<string>("");

  useEffect(() => {
    const token = apiClient.getToken();
    if (token) setClaims(parseTokenClaims(token));
    setLastUpdated(new Date().toLocaleTimeString("en-GB", { hour: "2-digit", minute: "2-digit" }));
  }, []);

  const departmentMetrics = [
    {
      id: "generation",
      name: "Generation Master Station",
      stat: "28.4 MW",
      sub: "Kurra & Jekko Hydro (99.8% Avail)",
      color: "emerald",
      icon: Zap,
      href: "/portal/generation-master",
    },
    {
      id: "dispatch",
      name: "Central Energy Dispatch",
      stat: "31.2 MW Peak",
      sub: "14 Active Feeders &bull; 0 Grid Collapse",
      color: "sky",
      icon: Activity,
      href: "/portal/energy-management",
    },
    {
      id: "services",
      name: "Field Engineering & Maintenance",
      stat: "2.1h MTTR",
      sub: "18 Work Orders Attended &bull; 100% PTW",
      color: "blue",
      icon: Clock,
      href: "/portal/services/feeder",
    },
    {
      id: "care",
      name: "Customer Care & Relations",
      stat: "92.8% SLA",
      sub: "84 Lodged &bull; 78 Resolved &bull; 14 KYC",
      color: "amber",
      icon: Users,
      href: "/portal/customer-care",
    },
    {
      id: "commercial",
      name: "Billing & Revenue Ops",
      stat: "₦42.8M Daily",
      sub: "3,420 STS Tokens &bull; 91.4% Vended",
      color: "purple",
      icon: DollarSign,
      href: "/portal/billing-commercial",
    },
    {
      id: "stores",
      name: "Warehouse & Material Stockpile",
      stat: "98.2% Ready",
      sub: "Conductors, 500kVA TX & Fuses Healthy",
      color: "emerald",
      icon: Package,
      href: "/portal/stores",
    },
  ];

  return (
    <>
      <Head>
        <title>NESCO &bull; Central Daily Analytics Storyboard</title>
      </Head>

      <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans">
        <TopCommandBar
          claims={claims}
          departmentTitle="Central Daily Analytics & Sovereign Storyboard"
          onToggleChat={() => setIsChatOpen(!isChatOpen)}
        />

        <main className="flex-1 p-4 sm:p-6 space-y-6 max-w-7xl w-full mx-auto">
          {/* Top Banner: Central Company Narrative */}
          <div className="p-5 rounded-2xl bg-gradient-to-r from-emerald-950/70 via-slate-900 to-blue-950/80 border border-emerald-500/30 shadow-xl relative overflow-hidden">
            <div className="flex items-start space-x-4">
              <div className="p-3 rounded-xl bg-emerald-500/20 text-emerald-400 border border-emerald-500/40 shrink-0">
                <Sparkles className="w-6 h-6 animate-spin" />
              </div>
              <div className="space-y-2 max-w-5xl">
                <div className="flex flex-wrap items-center gap-2">
                  <span className="text-xs font-bold uppercase tracking-wider text-emerald-400">
                    NESCO Sovereign AI &bull; Overall Company Data Analytics Story for the Day
                  </span>
                  <span className="text-[10px] bg-emerald-900/60 text-emerald-300 px-2.5 py-0.5 rounded-full border border-emerald-700/50 font-mono">
                    Updated: {lastUpdated} WAT
                  </span>
                  <span className="text-[10px] bg-blue-900/60 text-blue-300 px-2.5 py-0.5 rounded-full border border-blue-700/50 font-mono">
                    All 6 Operational Spheres Linked
                  </span>
                </div>
                <p className="text-sm text-slate-200 leading-relaxed font-sans">
                  Today, NESCO&apos;s sovereign hydro cascades at Kurra Falls and Jekko dispatched 28.4 MW continuous baseline, sustaining 100% uninterrupted availability on institutional Band A feeders (Vom &amp; Rayfield). Field maintenance resolved 18 emergency line and transformer drop-out issues with a record MTTR of 2.1 hours. Customer Care registered 14 new NIMC-verified KYC consumers and cleared 78 customer tickets with 92.8% SLA compliance. Commercial operations processed ₦42.8M in prepaid vending and postpaid recoveries with zero system downtime. Stores verified critical reserves of 100mm² ACSR conductors and transformer oil drums.
                </p>
              </div>
            </div>
          </div>

          {/* Departmental Daily Story Cards Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {departmentMetrics.map((dept) => {
              const Icon = dept.icon;
              return (
                <div
                  key={dept.id}
                  className="bg-slate-900/90 border border-slate-800 hover:border-emerald-500/50 rounded-xl p-4 transition-all duration-200 shadow-md group"
                >
                  <div className="flex items-center justify-between mb-3">
                    <div className="flex items-center space-x-2.5">
                      <div className="p-2 rounded-lg bg-slate-800 text-emerald-400 group-hover:bg-emerald-600/20 transition-colors">
                        <Icon className="w-4 h-4" />
                      </div>
                      <h4 className="text-xs font-bold text-slate-200 uppercase tracking-wide">
                        {dept.name}
                      </h4>
                    </div>
                    <Link
                      href={dept.href}
                      className="text-slate-400 hover:text-emerald-400 transition-colors"
                      title="Navigate to Department"
                    >
                      <ArrowRight className="w-4 h-4" />
                    </Link>
                  </div>

                  <div className="text-2xl font-extrabold text-slate-100 font-mono mb-1">
                    {dept.stat}
                  </div>
                  <div className="text-xs text-slate-400 flex items-center justify-between">
                    <span dangerouslySetInnerHTML={{ __html: dept.sub }} />
                  </div>
                </div>
              );
            })}
          </div>

          {/* Cross-Company Operational Pulse Breakdown */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-5">
            {/* Left Card: Customer Care & Distribution Performance */}
            <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-5 space-y-4">
              <div className="flex items-center justify-between border-b border-slate-800 pb-3">
                <h3 className="text-xs font-bold uppercase tracking-wider text-slate-200 flex items-center space-x-2">
                  <CheckCircle2 className="w-4 h-4 text-emerald-400" />
                  <span>Customer Complaints &amp; Distribution Velocity</span>
                </h3>
                <span className="text-[11px] text-emerald-400 font-mono font-bold">92.8% SLA Resolved</span>
              </div>

              <div className="space-y-3 text-xs">
                <div>
                  <div className="flex justify-between text-[11px] mb-1">
                    <span className="text-slate-300">Complaints Resolved vs Lodged (Today)</span>
                    <span className="font-mono text-emerald-400 font-bold">78 / 84 (92.8%)</span>
                  </div>
                  <div className="w-full bg-slate-800 h-2 rounded-full overflow-hidden">
                    <div className="bg-emerald-500 h-full rounded-full" style={{ width: "92.8%" }} />
                  </div>
                </div>

                <div>
                  <div className="flex justify-between text-[11px] mb-1">
                    <span className="text-slate-300">Feeder &amp; Transformer MTTR Speed</span>
                    <span className="font-mono text-sky-400 font-bold">2.1h vs 4.0h benchmark</span>
                  </div>
                  <div className="w-full bg-slate-800 h-2 rounded-full overflow-hidden">
                    <div className="bg-sky-500 h-full rounded-full" style={{ width: "85%" }} />
                  </div>
                </div>

                <div>
                  <div className="flex justify-between text-[11px] mb-1">
                    <span className="text-slate-300">New Customer Registrations (KYC GPS Verified)</span>
                    <span className="font-mono text-amber-400 font-bold">14 Onboarded</span>
                  </div>
                  <div className="w-full bg-slate-800 h-2 rounded-full overflow-hidden">
                    <div className="bg-amber-500 h-full rounded-full" style={{ width: "70%" }} />
                  </div>
                </div>
              </div>
            </div>

            {/* Right Card: Hydro Generation & Revenue Collection */}
            <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-5 space-y-4">
              <div className="flex items-center justify-between border-b border-slate-800 pb-3">
                <h3 className="text-xs font-bold uppercase tracking-wider text-slate-200 flex items-center space-x-2">
                  <DollarSign className="w-4 h-4 text-blue-400" />
                  <span>Hydro Generation &amp; Vending Revenue</span>
                </h3>
                <span className="text-[11px] text-blue-400 font-mono font-bold">₦42.8M Total Inflow</span>
              </div>

              <div className="space-y-3 text-xs">
                <div>
                  <div className="flex justify-between text-[11px] mb-1">
                    <span className="text-slate-300">Prepaid STS Token Vending Fulfillment</span>
                    <span className="font-mono text-emerald-400 font-bold">3,420 Tokens (99.9%)</span>
                  </div>
                  <div className="w-full bg-slate-800 h-2 rounded-full overflow-hidden">
                    <div className="bg-emerald-500 h-full rounded-full" style={{ width: "99.9%" }} />
                  </div>
                </div>

                <div>
                  <div className="flex justify-between text-[11px] mb-1">
                    <span className="text-slate-300">Hydro MW Generation Baseload</span>
                    <span className="font-mono text-blue-400 font-bold">28.4 MW / 32.0 MW Cap</span>
                  </div>
                  <div className="w-full bg-slate-800 h-2 rounded-full overflow-hidden">
                    <div className="bg-blue-500 h-full rounded-full" style={{ width: "88.75%" }} />
                  </div>
                </div>

                <div>
                  <div className="flex justify-between text-[11px] mb-1">
                    <span className="text-slate-300">Stores Essential Material Reserves</span>
                    <span className="font-mono text-purple-400 font-bold">98.2% In-Stock</span>
                  </div>
                  <div className="w-full bg-slate-800 h-2 rounded-full overflow-hidden">
                    <div className="bg-purple-500 h-full rounded-full" style={{ width: "98.2%" }} />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </main>

        <GridAssistantDock
          isOpen={isChatOpen}
          onClose={() => setIsChatOpen(false)}
          staffRole={claims?.primary_role || "EXECUTIVE_VIEWER"}
        />
      </div>
    </>
  );
}

