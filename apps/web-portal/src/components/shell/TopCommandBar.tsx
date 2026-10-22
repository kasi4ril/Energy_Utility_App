/**
 * Zone A: Top Command Bar (Global).
 * Displays Logged-in Staff Name, Role Badge, Active Org Unit Scope,
 * Shift Handover Status Indicator, and SLA Notification Bell.
 */

import React from "react";
import { AuthClaims } from "../../types/auth";
import { logoutStaff } from "../../features/auth/api";

interface TopCommandBarProps {
  claims: AuthClaims | null;
  departmentTitle: string;
  onToggleChat: () => void;
  unreadAlertsCount?: number;
}

export const TopCommandBar: React.FC<TopCommandBarProps> = ({
  claims,
  departmentTitle,
  onToggleChat,
  unreadAlertsCount = 0,
}) => {
  return (
    <header className="h-16 bg-slate-900 border-b border-slate-800 px-6 flex items-center justify-between sticky top-0 z-40">
      <div className="flex items-center space-x-4">
        <div className="flex items-center space-x-2 text-emerald-400 font-bold tracking-wider text-base">
          <span className="p-1.5 rounded-lg bg-emerald-500/10 border border-emerald-500/30">
            ⚡ NESCO
          </span>
          <span className="text-slate-300">/</span>
          <span className="text-white text-sm font-semibold">{departmentTitle}</span>
        </div>

        <div className="hidden md:flex items-center space-x-2 pl-4 border-l border-slate-800 text-xs">
          <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
          <span className="text-slate-400">Shift Handover:</span>
          <span className="text-emerald-300 font-medium">ACTIVE (08:00 - 16:00)</span>
        </div>
      </div>

      <div className="flex items-center space-x-4">
        {claims && (
          <div className="flex items-center space-x-3 text-right">
            <div>
              <div className="text-xs font-medium text-white">{claims.full_name}</div>
              <div className="text-[10px] text-slate-400 font-mono">
                {claims.staff_number} &bull; {claims.org_unit_id}
              </div>
            </div>
            <span className="px-2.5 py-1 rounded-full text-[10px] font-bold tracking-wider uppercase bg-emerald-900/60 text-emerald-300 border border-emerald-700/50">
              {claims.primary_role}
            </span>
          </div>
        )}

        <button
          onClick={onToggleChat}
          className="relative p-2 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 hover:text-white transition-colors"
          title="Toggle @GridAssistant Chat Dock"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
          </svg>
          {unreadAlertsCount > 0 && (
            <span className="absolute -top-1 -right-1 w-4 h-4 rounded-full bg-rose-500 text-white text-[10px] font-bold flex items-center justify-center">
              {unreadAlertsCount}
            </span>
          )}
        </button>

        <button
          onClick={() => {
            logoutStaff();
            window.location.href = "/login";
          }}
          className="p-2 rounded-lg bg-slate-800 hover:bg-red-900/40 text-slate-400 hover:text-red-300 transition-colors text-xs"
          title="Sign Out"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
          </svg>
        </button>
      </div>
    </header>
  );
};

