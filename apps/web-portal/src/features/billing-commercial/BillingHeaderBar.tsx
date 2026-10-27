/**
 * Thin Top Header Bar for Commercial Operations & Billing.
 * Features Sovereign Crest, real-time WAT clock, and fast links.
 */

import React, { useState, useEffect } from "react";
import Link from "next/link";
import { AuthClaims } from "../../types/auth";

interface BillingHeaderBarProps {
  claims: AuthClaims | null;
  onToggleChat?: () => void;
  todayCollectionsNgn?: number;
  collectionEfficiencyPct?: number;
}

export function BillingHeaderBar({
  claims,
  onToggleChat,
  todayCollectionsNgn = 4850000,
  collectionEfficiencyPct = 94.2,
}: BillingHeaderBarProps) {
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
          className="flex items-center space-x-2 text-white hover:text-teal-400 transition group"
          title="Return to NESCO General Landing Page"
        >
          <div className="w-8 h-8 rounded-lg bg-teal-950/80 border border-teal-500/40 flex items-center justify-center text-teal-400 font-black shadow-inner shadow-teal-900/40 text-sm group-hover:scale-105 transition-transform">
            ⚡
          </div>
          <div>
            <div className="font-extrabold tracking-wider text-slate-100 flex items-center space-x-1.5">
              <span>NESCO</span>
              <span className="text-[10px] text-teal-400 font-mono font-normal tracking-normal px-1.5 py-0.2 rounded bg-teal-950/60 border border-teal-800/50">
                1923
              </span>
            </div>
            <div className="text-[10px] text-slate-400 font-medium">
              Commercial Operations &bull; Billing
            </div>
          </div>
        </Link>

        <span className="text-slate-700 hidden sm:inline">|</span>

        {/* Real-time WAT Clock */}
        <div className="hidden md:flex items-center space-x-1.5 text-slate-400 font-mono text-[11px] bg-slate-900/60 px-2.5 py-1 rounded-md border border-slate-800">
          <span className="w-2 h-2 rounded-full bg-teal-400 animate-ping inline-block" />
          <span dangerouslySetInnerHTML={{ __html: watTime || "Loading WAT..." }} />
        </div>
      </div>

      {/* Middle: Revenue Metrics */}
      <div className="hidden lg:flex items-center space-x-3">
        <div className="flex items-center space-x-1.5 px-2.5 py-1 rounded bg-slate-900 border border-slate-800">
          <span className="text-[10px] text-slate-400 uppercase font-semibold">Today's Collections:</span>
          <span className="font-mono font-bold text-emerald-400">
            &#8358;{(todayCollectionsNgn / 1000000).toFixed(2)}M
          </span>
        </div>

        <div className="flex items-center space-x-1.5 px-2.5 py-1 rounded bg-slate-900 border border-slate-800">
          <span className="text-[10px] text-slate-400 uppercase font-semibold">Collection Efficiency:</span>
          <span className="font-mono font-bold text-teal-300">{collectionEfficiencyPct}%</span>
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
            className="px-2.5 py-1 rounded bg-teal-950/40 hover:bg-teal-900/60 text-teal-300 border border-teal-800/40 text-[11px] font-medium transition flex items-center space-x-1"
          >
            <span>📊</span>
            <span>Central Analytics</span>
          </Link>
        </div>

        {/* Staff Attribution Badge */}
        <div className="flex items-center space-x-2 pl-2 border-l border-slate-800">
          <div className="w-7 h-7 rounded-full bg-teal-900/50 border border-teal-500/40 flex items-center justify-center text-teal-300 font-bold text-xs">
            {claims?.full_name ? claims.full_name[0] : "K"}
          </div>
          <div className="hidden xl:block text-left">
            <div className="text-[11px] font-semibold text-slate-200 leading-tight">
              {claims?.full_name || "Mr. K. O. Alabi"}
            </div>
            <div className="text-[10px] text-slate-400 font-mono leading-tight">
              {claims?.primary_role || "HEAD_COMMERCIAL_OPERATIONS"}
            </div>
          </div>
        </div>

        {/* Grid Assistant Chat Toggle */}
        {onToggleChat && (
          <button
            onClick={onToggleChat}
            className="p-1.5 rounded-lg bg-slate-900 hover:bg-slate-800 text-teal-400 border border-slate-800 hover:border-teal-500/40 transition flex items-center space-x-1"
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

