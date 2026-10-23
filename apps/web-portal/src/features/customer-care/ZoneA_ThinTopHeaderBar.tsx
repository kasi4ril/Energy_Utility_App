/**
 * Thin Top Header Bar for Customer Care & Service Delivery Department (Zone A).
 * Displays official department title, real-time ticking WAT clock, staff attribution badge,
 * and live SLA breach alarm bell. Strictly <= 350 lines.
 */

import React, { useState, useEffect } from "react";
import Link from "next/link";
import { AuthClaims } from "../../types/auth";
import { Clock, ShieldCheck, Bell, Sparkles, BarChart3, ArrowLeft } from "lucide-react";

interface ZoneAProps {
  claims: AuthClaims | null;
  breachCount?: number;
  onToggleChat?: () => void;
}

export function ZoneA_ThinTopHeaderBar({
  claims,
  breachCount = 1,
  onToggleChat,
}: ZoneAProps) {
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

  const staffName = claims?.full_name || "Amina Yusuf";
  const staffId = claims?.staff_id || "STAFF-CC-02";
  const roleName = claims?.primary_role?.replace(/_/g, " ") || "CUSTOMER CARE OFFICER";

  return (
    <header className="w-full bg-slate-900/95 backdrop-blur-md border-b border-slate-800 px-4 py-2 flex flex-wrap items-center justify-between gap-3 text-xs shadow-md sticky top-0 z-40 shrink-0">
      {/* Left: Big Prominent NESCO Logo Crest & Department Identity */}
      <div className="flex items-center space-x-3.5">
        <Link
          href="/"
          className="flex items-center space-x-3 px-3.5 py-1.5 bg-gradient-to-r from-emerald-950 via-slate-900 to-slate-950 border border-emerald-500/50 rounded-xl shadow-lg hover:border-emerald-400 transition-all shrink-0 group"
          title="NESCO 1923 • Return to Main Landing Page"
        >
          <div className="w-11 h-11 rounded-lg bg-emerald-500/20 border border-emerald-500/50 flex items-center justify-center text-2xl text-emerald-400 group-hover:scale-105 transition-transform shadow-[0_0_16px_rgba(16,185,129,0.35)] shrink-0">
            ⚡
          </div>
          <div className="leading-tight">
            <div className="flex items-center space-x-1.5">
              <span className="font-black tracking-widest text-white text-lg sm:text-xl font-mono">
                NESCO
              </span>
              <span className="text-[10px] bg-emerald-500/20 text-emerald-300 px-1.5 py-0.2 rounded font-extrabold font-mono border border-emerald-500/40">
                1923
              </span>
            </div>
            <div className="text-[10px] text-emerald-400 font-bold tracking-tight uppercase font-mono">
              Sovereign Utility
            </div>
          </div>
        </Link>
        <div className="h-10 w-px bg-slate-800 hidden sm:block" />
        <div>
          <h1 className="font-bold tracking-wider text-slate-100 text-sm sm:text-base uppercase">
            Customer Care &amp; Service Delivery Department
          </h1>
          <p className="text-[10px] text-slate-400 font-mono tracking-tight">
            NESCO Enterprise Sovereign Grid Utility &bull; NERC Compliance Center
          </p>
        </div>
      </div>

      {/* Middle: Live Ticking Africa/Lagos Clock */}
      <div className="hidden lg:flex items-center space-x-2 bg-slate-950/70 border border-slate-800 px-3 py-1.5 rounded-md text-slate-300 font-mono">
        <Clock className="w-3.5 h-3.5 text-blue-400 animate-pulse" />
        <span>{dateStr}</span>
        <span className="text-slate-600">|</span>
        <span className="text-emerald-400 font-semibold">{timeStr}</span>
      </div>

      {/* Right: Navigation Buttons (Shifted Right), Staff Attribution Badge & Alerts */}
      <div className="flex items-center space-x-2.5">
        {/* Back Button (Renamed from Landing Page, Shifted Right) */}
        <Link
          href="/"
          className="flex items-center space-x-1.5 bg-slate-800 hover:bg-slate-700 text-slate-200 hover:text-white border border-slate-700 px-3 py-1.5 rounded-lg text-xs font-semibold transition-colors shadow-sm"
          title="Back to Main Landing Page"
        >
          <ArrowLeft className="w-3.5 h-3.5 text-blue-400" />
          <span>Back</span>
        </Link>

        {/* Central Daily Analytics Button (Shifted Right) */}
        <Link
          href="/portal/central-analytics"
          className="flex items-center space-x-1.5 bg-emerald-950/80 hover:bg-emerald-900/90 text-emerald-300 border border-emerald-500/50 px-3 py-1.5 rounded-lg font-semibold text-xs transition-colors shadow-sm"
          title="Overall Company Data Analytics Story for the Day"
        >
          <BarChart3 className="w-3.5 h-3.5 text-emerald-400" />
          <span className="hidden sm:inline">Central Analytics</span>
        </Link>

        {/* SLA Breach Bell */}
        <div className="relative group">
          <button
            className="p-1.5 rounded-md bg-slate-800 hover:bg-slate-700 text-slate-300 relative transition-colors"
            title="SLA Breach & Escalation Radar"
          >
            <Bell className="w-4 h-4 text-amber-400" />
            {breachCount > 0 && (
              <span className="absolute -top-1 -right-1 px-1.5 py-0.2 bg-red-600 text-white text-[9px] font-bold rounded-full animate-bounce">
                {breachCount}
              </span>
            )}
          </button>
        </div>

        {/* AI GridAssistant Co-Pilot Button */}
        {onToggleChat && (
          <button
            onClick={onToggleChat}
            className="flex items-center space-x-1.5 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-500 hover:to-indigo-500 text-white px-2.5 py-1.5 rounded-md font-medium shadow-sm transition-all text-[11px]"
          >
            <Sparkles className="w-3.5 h-3.5" />
            <span>@GridAssistant</span>
          </button>
        )}

        {/* Staff Attribution Card */}
        <div className="flex items-center space-x-2.5 bg-slate-950/80 border border-slate-800 px-2.5 py-1 rounded-md">
          <div className="w-6 h-6 rounded-full bg-blue-500/20 text-blue-400 flex items-center justify-center font-bold text-[10px] border border-blue-500/40">
            {staffName.charAt(0)}
          </div>
          <div className="text-left leading-tight">
            <div className="font-semibold text-slate-200 text-[11px] flex items-center space-x-1">
              <span>{staffName}</span>
              <ShieldCheck className="w-3 h-3 text-emerald-400 inline" />
            </div>
            <div className="text-[10px] text-slate-400 font-mono">
              {staffId} &bull; <span className="text-blue-400">{roleName}</span> &bull; 08:00–16:00
            </div>
          </div>
        </div>
      </div>
    </header>
  );
}

