/**
 * Out of Service Plants Alert Banner:
 * Highlights tripped or shut-down hydro generating stations, root causes,
 * lost capacity, and GM workshop review status.
 */

import React from "react";
import { HydroPlantSummary } from "./types";
import { ShieldAlert } from "lucide-react";

interface OutOfServiceAlertBannerProps {
  outOfServicePlants: HydroPlantSummary[];
}

export const OutOfServiceAlertBanner: React.FC<OutOfServiceAlertBannerProps> = ({
  outOfServicePlants,
}) => {
  if (outOfServicePlants.length === 0) return null;

  return (
    <div className="bg-rose-950/40 border border-rose-500/60 rounded-xl p-4 shadow-lg space-y-3">
      <div className="flex items-center justify-between">
        <div className="flex items-center space-x-2.5">
          <div className="p-1.5 rounded-lg bg-rose-500/20 text-rose-400 border border-rose-500/40 animate-pulse">
            <ShieldAlert className="w-5 h-5" />
          </div>
          <div>
            <h3 className="font-bold text-rose-200 text-sm tracking-wide">
              OUT OF SERVICE GENERATING PLANTS ({outOfServicePlants.length})
            </h3>
            <p className="text-[11px] text-rose-300/80">
              Real-time outage triage: Immediate engineering attention and GM material authorizations active.
            </p>
          </div>
        </div>
        <span className="px-2.5 py-1 rounded-full text-[10px] font-extrabold uppercase bg-rose-500/30 text-rose-200 border border-rose-500/50">
          Active Outage
        </span>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-3 pt-1">
        {outOfServicePlants.map((plant) => (
          <div
            key={plant.id}
            className="bg-slate-900/90 border border-rose-800/80 rounded-lg p-3.5 space-y-2"
          >
            <div className="flex items-center justify-between">
              <span className="font-bold text-white text-xs">{plant.name}</span>
              <span className="px-2 py-0.5 rounded text-[10px] font-mono font-bold bg-rose-500/20 text-rose-300 border border-rose-500/40">
                {plant.outage_reason || "TRIPPED_ON_FAULT"}
              </span>
            </div>
            <div className="text-[11px] text-slate-300">
              <strong className="text-slate-400">Root Cause:</strong> Turbine runner hub cavitation pitting and shaft sleeve groove. Dispatched to Central Machine Workshop.
            </div>
            <div className="flex items-center justify-between text-[10px] text-slate-400 pt-1 border-t border-slate-800">
              <span>Lost Capacity: <strong className="text-rose-300 font-mono">{plant.installed_capacity_mw} MW</strong></span>
              <span>Turbine: <strong className="text-slate-200">{plant.turbine_type}</strong></span>
              <span className="text-amber-400 font-semibold">Under GM Review</span>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

