/**
 * Thin Top Header Bar for Transformer Maintenance Engineering Unit.
 * Features Sovereign Crest, real-time WAT clock, and fast links.
 */

import React, { useState, useEffect } from "react";
import Link from "next/link";
import { AuthClaims } from "../../types/auth";

interface TransformerHeaderBarProps {
  claims: AuthClaims | null;
  onToggleChat?: () => void;
  activeDtsCount?: number;
  criticalOilAlertsCount?: number;
}

export function TransformerHeaderBar({
  claims,
  onToggleChat,
  activeDtsCount = 142,
  criticalOilAlertsCount = 1,
}: TransformerHeaderBarProps) {
  const [watTime, setWatTime] = useState<string>("");

  useEffect(() => {
    const updateTime = () => {
      const now = new Date();
      const formatted = now.toLocaleTimeString("en-GB", {
        timeZone: "Africa/Lagos",
        hour12: false,
        hour: "2-digit",
        minute: "2-digit",
        second: "2-digit",
      });
      const dateFormatted = now.toLocaleDateString("en-GB", {
        timeZone: "Africa/Lagos",
        day: "2-digit",
        month: "short",
        year: "numeric",
      });
      setWatTime(`${dateFormatted} &bull; ${formatted} WAT`);
    };

    updateTime();
    const timer = setInterval(updateTime, 1000);
    return () => clearInterval(timer);
  }, []);

  return (
    <header className="h-14 bg-slate-950 border-b border-slate-800/80 px-4 flex items-center justify-between text-xs sticky top-0 z-40 select-none">
      {/* Far Left: Crest & Title */}
      <div className="flex items-center space-x-3">
        <Link
          href="/"
          className="flex items-center space-x-2 text-white hover:text-amber-400 transition group"
          title="Return to NESCO General Landing Page"
        >
          <div className="w-8 h-8 rounded-lg bg-amber-950/80 border border-amber-500/40 flex items-center justify-center text-amber-400 font-black shadow-inner shadow-amber-900/40 text-sm group-hover:scale-105 transition-transform">
            ⚡
          </div>
          <div>
            <div className="font-extrabold tracking-wider text-slate-100 flex items-center space-x-1.5">
              <span>NESCO</span>
              <span className="text-[10px] text-amber-400 font-mono font-normal tracking-normal px-1.5 py-0.2 rounded bg-amber-950/60 border border-amber-800/50">
                1923
              </span>
            </div>
            <div className="text-[10px] text-slate-400 font-medium">
              Services &bull; Transformer Unit
            </div>
          </div>
        </Link>

        <span className="text-slate-700 hidden sm:inline">|</span>

        {/* Real-time WAT Clock */}
        <div className="hidden md:flex items-center space-x-1.5 text-slate-400 font-mono text-[11px] bg-slate-900/60 px-2.5 py-1 rounded-md border border-slate-800">
          <span className="w-2 h-2 rounded-full bg-amber-400 animate-ping inline-block" />
          <span dangerouslySetInnerHTML={{ __html: watTime || "Loading WAT..." }} />
        </div>
      </div>

      {/* Middle: Transformer Fleet Ticker */}
      <div className="hidden lg:flex items-center space-x-3">
        <div className="flex items-center space-x-1.5 px-2.5 py-1 rounded bg-slate-900 border border-slate-800">
          <span className="text-[10px] text-slate-400 uppercase font-semibold">Active DT Fleet:</span>
          <span className="font-mono font-bold text-emerald-400">{activeDtsCount} DTs</span>
        </div>

        <div className="flex items-center space-x-1.5 px-2.5 py-1 rounded bg-slate-900 border border-slate-800">
          <span className="text-[10px] text-slate-400 uppercase font-semibold">DGA Dielectric Alerts:</span>
          <span className="font-mono font-bold text-amber-400">{criticalOilAlertsCount} BDV &lt; 30kV</span>
        </div>
      </div>

      {/* Far Right: Staff Badge, Shortcuts, Chat */}
      <div className="flex items-center space-x-3">
        {/* Navigation Return Links */}
        <div className="hidden sm:flex items-center space-x-2">
          <Link
            href="/"
            className="px-2.5 py-1 rounded bg-slate-900/80 hover:bg-slate-800 text-slate-300 hover:text-white border border-slate-800 text-[11px] font-medium transition"
          >
            &larr; Landing
          </Link>
          <Link
            href="/portal/central-analytics"
            className="px-2.5 py-1 rounded bg-amber-950/40 hover:bg-amber-900/60 text-amber-300 border border-amber-800/40 text-[11px] font-medium transition flex items-center space-x-1"
          >
            <span>📊</span>
            <span>Central Analytics</span>
          </Link>
        </div>

        {/* Staff Attribution Badge */}
        <div className="flex items-center space-x-2 pl-2 border-l border-slate-800">
          <div className="w-7 h-7 rounded-full bg-amber-900/50 border border-amber-500/40 flex items-center justify-center text-amber-300 font-bold text-xs">
            {claims?.full_name ? claims.full_name[0] : "T"}
          </div>
          <div className="hidden xl:block text-left">
            <div className="text-[11px] font-semibold text-slate-200 leading-tight">
              {claims?.full_name || "Engr. F. D. Musa"}
            </div>
            <div className="text-[10px] text-slate-400 font-mono leading-tight">
              {claims?.primary_role || "UNIT_HEAD_TRANSFORMER_MAINT"}
            </div>
          </div>
        </div>

        {/* Grid Assistant Chat Toggle */}
        {onToggleChat && (
          <button
            onClick={onToggleChat}
            className="p-1.5 rounded-lg bg-slate-900 hover:bg-slate-800 text-amber-400 border border-slate-800 hover:border-amber-500/40 transition flex items-center space-x-1"
            title="Ask @GridAssistant"
          >
            <span className="text-sm">🤖</span>
            <span className="hidden md:inline text-[10px] font-bold">Ask AI</span>
          </button>
        )}
      </div>
    </header>
  );
}

