/**
 * Architectural & Engineering Pillars for NESCO Sovereign Utility Platform.
 * Highlights local air-gapped AI, cryptographic hash-chaining, PostGIS, and 5-step PTW.
 */

import React from "react";

export function LandingArchitecturePillars() {
  const pillars = [
    {
      icon: "🤖",
      title: "Air-Gapped AI Sovereign Sidecar",
      badge: "Local Dual-Brain",
      description:
        "100% on-premise Ollama AI Guardian executing Maker-Checker fraud scoring, photo GPS validation, and @GridAssistant in-app conversational co-pilot without external cloud transmission.",
    },
    {
      icon: "⛓️",
      title: "Chained Cryptographic WorkNotes",
      badge: "SHA-256 Ledger",
      description:
        "Immutable previous-hash chained ledger protecting shift handovers, fault observations, and executive sign-off decisions with strict database trigger-level append-only constraints.",
    },
    {
      icon: "🔒",
      title: "5-Step Safety PTW Protocol",
      badge: "NERC Grid Code",
      description:
        "Permit-to-Work safety lock: Requisition &rarr; HSE Approval &rarr; Physical Substation Isolation &rarr; Work In Progress &rarr; Earth Removal Clearance &rarr; Line Reclosing.",
    },
    {
      icon: "📍",
      title: "PostGIS Spatial Grid Topology",
      badge: "Haversine Distance",
      description:
        "Rigorous 5-tier electrical grid hierarchy (Hydro Plant &rarr; 33kV Substation &rarr; 11kV Feeder &rarr; Transformer DT &rarr; Customer Meter) with spherical GPS distance verification.",
    },
  ];

  return (
    <section className="bg-slate-900/60 border-t border-slate-800/80 py-12 px-6 select-none">
      <div className="max-w-7xl mx-auto space-y-8">
        <div>
          <div className="text-[10px] uppercase font-bold text-cyan-400 font-mono tracking-widest">
            Architectural Guarantees
          </div>
          <h2 className="text-xl font-bold text-white tracking-tight mt-0.5">
            Mission-Critical Utility Governance Architecture
          </h2>
          <p className="text-xs text-slate-400 mt-1 max-w-xl">
            Engineered specifically for the operational demands of Nigeria's pioneer hydroelectric power utility.
          </p>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-5">
          {pillars.map((p) => (
            <div
              key={p.title}
              className="p-5 rounded-xl bg-slate-950/70 border border-slate-800/90 space-y-3 hover:border-slate-700 transition"
            >
              <div className="flex items-center justify-between">
                <span className="text-2xl">{p.icon}</span>
                <span className="text-[9px] font-mono uppercase px-2 py-0.5 rounded bg-slate-900 border border-slate-700 text-slate-400">
                  {p.badge}
                </span>
              </div>
              <h3 className="text-sm font-bold text-white">{p.title}</h3>
              <p className="text-xs text-slate-400 leading-relaxed font-sans" dangerouslySetInnerHTML={{ __html: p.description }} />
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

