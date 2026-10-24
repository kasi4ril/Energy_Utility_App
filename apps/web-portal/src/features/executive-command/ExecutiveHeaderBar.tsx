/**
 * Thin Top Header Bar for Executive Command Tier (Managing Director & CEO).
 * Features official sovereign utility crest, real-time WAT clock, staff attribution badge,
 * NERC Governance stamp, and rapid navigation links. Strictly <= 350 lines.
 */

import React, { useState, useEffect } from "react";
import Link from "next/link";
import { AuthClaims } from "../../types/auth";
import { Clock, ShieldCheck, Bell, Sparkles, BarChart3, ArrowLeft, Crown } from "lucide-react";

interface ExecutiveHeaderBarProps {
  claims: AuthClaims | null;
  pendingApprovalsCount?: number;
  onToggleChat?: () => void;
}

export const ExecutiveHeaderBar: React.FC<ExecutiveHeaderBarProps> = ({
  claims,
  pendingApprovalsCount = 5,
  onToggleChat,
}) => {
  const [timeStr, setTimeStr] = useState<string>("");
  const [dateStr, setDateStr] = useState<string>("");

  useEffect(() => {
    const updateTime = () => {
      const now = new Date();
      setTimeStr(
        now.toLocaleTimeString("en-GB", {
          timeZone: "Africa/Lagos",
          hour: "2-digit",
          minute: "2-digit",
          second: "2-digit",
        }) + " WAT"
      );
      setDateStr(
        now.toLocaleDateString("en-GB", {
          timeZone: "Africa/Lagos",
          weekday: "short",
          day: "numeric",
          month: "short",
          year: "numeric",
        })
      );
    };

    updateTime();
    const interval = setInterval(updateTime, 1000);
    return () => clearInterval(interval);
  }, []);

  const executiveName = claims?.full_name || "Engr. Dr. E. K. Dung, FNSE, NPOM";
  const executiveId = claims?.staff_id || "STAFF-MD-01";
  const roleName = "MANAGING DIRECTOR & CEO";

  return (
    <header className="w-full bg-slate-900/95 backdrop-blur-md border-b border-amber-500/30 px-4 py-2 flex flex-wrap items-center justify-between gap-3 text-xs shadow-md sticky top-0 z-40 shrink-0">
      {/* Left: NESCO Crest & MD Executive Identity */}
      <div className="flex items-center space-x-3.5">
        <Link
          href="/"
          className="flex items-center space-x-3 px-3.5 py-1.5 bg-gradient-to-r from-amber-950 via-slate-900 to-slate-950 border border-amber-500/50 rounded-xl shadow-lg hover:border-amber-400 transition-all shrink-0 group"
          title="NESCO 1923 • Return to Main Landing Page"
        >
          <div className="w-10 h-10 rounded-lg bg-amber-500/20 border border-amber-500/50 flex items-center justify-center text-xl text-amber-400 group-hover:scale-105 transition-transform shadow-[0_0_16px_rgba(245,158,11,0.35)] shrink-0">
            👑
          </div>
          <div className="leading-tight">
            <div className="flex items-center space-x-1.5">
              <span className="font-black tracking-widest text-white text-lg sm:text-xl font-mono">
                NESCO
              </span>
              <span className="text-[10px] bg-amber-500/20 text-amber-300 px-1.5 py-0.2 rounded font-extrabold font-mono border border-amber-500/40">
                1923
              </span>
            </div>
            <div className="text-[10px] text-amber-400 font-bold tracking-tight uppercase font-mono">
              Executive Command
            </div>
          </div>
        </Link>

        <div className="h-10 w-px bg-slate-800 hidden sm:block" />

        <div>
          <h1 className="font-bold tracking-wider text-slate-100 text-sm sm:text-base uppercase flex items-center space-x-2">
            <span>Office of the Managing Director &amp; CEO</span>
            <span className="text-[10px] bg-amber-500/20 text-amber-300 px-2 py-0.5 rounded-full border border-amber-500/40 font-mono font-bold">
              Tier-1 Command
            </span>
          </h1>
          <p className="text-[10px] text-slate-400 font-mono tracking-tight">
            Executive Approvals, Inter-Departmental Minuting &amp; Corporate Governance
          </p>
        </div>
      </div>

      {/* Middle: Live Ticking WAT Clock */}
      <div className="hidden lg:flex items-center space-x-2 bg-slate-950/70 border border-slate-800 px-3 py-1.5 rounded-md text-slate-300 font-mono">
        <Clock className="w-3.5 h-3.5 text-amber-400 animate-pulse" />
        <span>{dateStr}</span>
        <span className="text-slate-600">|</span>
        <span className="text-amber-400 font-semibold">{timeStr}</span>
      </div>

      {/* Right: Quick Links, Notification Bell & MD Attribution */}
      <div className="flex items-center space-x-2.5">
        {/* Quick link: General Landing Page */}
        <Link
          href="/"
          className="px-2.5 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 border border-slate-700 flex items-center space-x-1.5 text-xs font-semibold transition"
          title="Return to Main General Landing Page"
        >
          <ArrowLeft className="w-3.5 h-3.5 text-amber-400" />
          <span className="hidden sm:inline">Landing Page</span>
        </Link>

        {/* Quick link: Central Analytics */}
        <Link
          href="/portal/central-analytics"
          className="px-2.5 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 border border-slate-700 flex items-center space-x-1.5 text-xs font-semibold transition"
          title="Company-Wide Central Analytics Board"
        >
          <BarChart3 className="w-3.5 h-3.5 text-cyan-400" />
          <span className="hidden sm:inline">Central Analytics</span>
        </Link>

        {/* Pending Approvals Bell Alert */}
        <div
          className={`flex items-center space-x-1.5 px-2.5 py-1.5 rounded-lg border text-xs font-mono font-bold ${
            pendingApprovalsCount > 0
              ? "bg-amber-500/20 text-amber-300 border-amber-500/40 animate-pulse"
              : "bg-slate-800 text-slate-400 border-slate-700"
          }`}
          title={`${pendingApprovalsCount} Requisitions Awaiting MD Sign-Off`}
        >
          <Bell className="w-3.5 h-3.5 text-amber-400" />
          <span>{pendingApprovalsCount} Due</span>
        </div>

        {/* AI Co-Pilot Drawer Toggle */}
        {onToggleChat && (
          <button
            onClick={onToggleChat}
            className="flex items-center space-x-1.5 bg-gradient-to-r from-amber-600 to-amber-700 hover:from-amber-500 hover:to-amber-600 text-white font-bold px-3 py-1.5 rounded-lg shadow-md transition-all active:scale-95"
            title="Open GridAssistant AI Executive Co-Pilot"
          >
            <Sparkles className="w-3.5 h-3.5 text-amber-200 animate-spin" />
            <span className="hidden md:inline">Executive AI</span>
          </button>
        )}

        {/* MD Attribution Profile Pill */}
        <div className="flex items-center space-x-2 bg-slate-950/80 border border-amber-500/40 px-3 py-1 rounded-lg">
          <div className="w-7 h-7 rounded-full bg-amber-500/20 text-amber-400 border border-amber-500/40 flex items-center justify-center font-bold text-xs">
            MD
          </div>
          <div className="text-left leading-tight hidden xl:block">
            <div className="font-bold text-white text-[11px] truncate max-w-[170px]">
              {executiveName}
            </div>
            <div className="text-[9px] text-amber-400 font-mono tracking-tight font-bold">
              {roleName} &bull; {executiveId}
            </div>
          </div>
        </div>
      </div>
    </header>
  );
};

