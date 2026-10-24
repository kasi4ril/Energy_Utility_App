/**
 * Next.js Page: Generation Department & Master Station Command Hub (/portal/generation-master).
 * Features persistent left-hand navigation pane with direct endpoints for:
 * 1. Real-Time Generation Analytics (water levels, MW vs capacity, health, outages)
 * 2. Hourly Reading Sheet Form (direct database ingestion)
 * 3. Change of Shift Desk (physical muster & immutable handover notes)
 * 4. Supervisor Requisition Section (direct notes & spares requests to immediate supervisor)
 */

import React, { useState, useEffect } from "react";
import Head from "next/head";
import Link from "next/link";
import { useRouter } from "next/router";
import { GenerationNav_LeftSidebarPane } from "../../features/generation-master/GenerationNav_LeftSidebarPane";
import { GenerationAnalyticsView } from "../../features/generation-master/GenerationAnalyticsView";
import { ReadingSheetForm } from "../../features/generation-master/ReadingSheetForm";
import { ChangeOfShiftForm } from "../../features/generation-master/ChangeOfShiftForm";
import { SupervisorRequisitionView } from "../../features/generation-master/SupervisorRequisitionView";
import { PlantFleetLiveBoard } from "../../features/generation-master/PlantFleetLiveBoard";
import { WorkshopFabricationView } from "../../features/generation-master/WorkshopFabricationView";
import { GridAssistantDock } from "../../components/shell/GridAssistantDock";
import { WorkNotesDrawer } from "../../components/shell/WorkNotesDrawer";
import {
  fetchPlantFleet,
  fetchShiftRollcalls,
  fetchPlantFaults,
  fetchGenerationAnalytics,
  fetchSupervisorRequisitions,
} from "../../features/generation-master/api";
import {
  GenerationTask,
  HydroPlantSummary,
  ShiftRollcallRecord,
  PlantFaultJob,
  GenerationAnalytics,
  SupervisorRequisition,
} from "../../features/generation-master/types";
import { parseTokenClaims } from "../../lib/auth";
import { AuthClaims } from "../../types/auth";
import { apiClient } from "../../lib/api-client";
import { Clock, Radio, ShieldCheck, Sparkles, Zap, ArrowLeft, BarChart3 } from "lucide-react";

interface GenerationMasterPageProps {
  defaultView?: GenerationTask;
}

export default function GenerationMasterPage({ defaultView = "analytics" }: GenerationMasterPageProps) {
  const router = useRouter();
  const [claims, setClaims] = useState<AuthClaims | null>(null);
  const [activeTask, setActiveTask] = useState<GenerationTask>(defaultView);
  const [plants, setPlants] = useState<HydroPlantSummary[]>([]);
  const [rollcalls, setRollcalls] = useState<ShiftRollcallRecord[]>([]);
  const [faults, setFaults] = useState<PlantFaultJob[]>([]);
  const [requisitions, setRequisitions] = useState<SupervisorRequisition[]>([]);
  const [analytics, setAnalytics] = useState<GenerationAnalytics | null>(null);
  const [selectedPlant, setSelectedPlant] = useState<HydroPlantSummary | null>(null);
  const [isChatOpen, setIsChatOpen] = useState(false);

  // Time state for WAT clock
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

  const loadData = async () => {
    const [pFleet, rList, fList, aData, reqList] = await Promise.all([
      fetchPlantFleet(),
      fetchShiftRollcalls(),
      fetchPlantFaults(),
      fetchGenerationAnalytics(),
      fetchSupervisorRequisitions(),
    ]);
    setPlants(pFleet);
    setRollcalls(rList);
    setFaults(fList);
    setAnalytics(aData);
    setRequisitions(reqList);
  };

  useEffect(() => {
    const token = apiClient.getToken();
    if (token) setClaims(parseTokenClaims(token));
    loadData();

    // Check query param for initial view
    if (router.query.view && typeof router.query.view === "string") {
      setActiveTask(router.query.view as GenerationTask);
    }
  }, [router.query.view]);

  const totalMW = plants.reduce((sum, p) => sum + p.current_mw, 0);
  const activePlantsCount = plants.filter((p) => p.state === "IN_SERVICE").length;
  const openFaultsCount = faults.filter((f) => f.status !== "RESOLVED_REINSTALLED").length;
  const pendingRequisitionsCount = requisitions.filter(
    (r) => r.status === "PENDING_SUPERVISOR_REVIEW" || r.status === "SUPERVISOR_ACKNOWLEDGED"
  ).length;

  const staffId = claims?.staff_id || "STAFF-GEN-OP-01";
  const staffName = claims?.full_name || "Ibrahim Pam (Master Operator)";

  return (
    <>
      <Head>
        <title>NESCO Generation &bull; Master Station Command Hub (Plants 1–6)</title>
      </Head>

      <div className="h-screen bg-slate-950 text-slate-100 flex flex-col font-sans overflow-hidden">
        {/* Top Header Bar */}
        <header className="w-full bg-slate-900/95 backdrop-blur-md border-b border-slate-800 px-4 py-2 flex flex-wrap items-center justify-between gap-3 text-xs shadow-md sticky top-0 z-40 shrink-0">
          {/* Left: NESCO Crest & Department Title */}
          <div className="flex items-center space-x-3.5">
            <Link
              href="/"
              className="flex items-center space-x-3 px-3 py-1 bg-gradient-to-r from-cyan-950 via-slate-900 to-slate-950 border border-cyan-500/50 rounded-xl shadow-lg hover:border-cyan-400 transition-all shrink-0 group"
              title="Return to Main Landing Page"
            >
              <div className="w-9 h-9 rounded-lg bg-cyan-500/20 border border-cyan-500/50 flex items-center justify-center text-xl text-cyan-400 group-hover:scale-105 transition-transform shadow-[0_0_12px_rgba(6,182,212,0.35)] shrink-0">
                ⚡
              </div>
              <div className="leading-tight">
                <div className="flex items-center space-x-1.5">
                  <span className="font-black tracking-widest text-white text-base sm:text-lg font-mono">
                    NESCO
                  </span>
                  <span className="text-[9px] bg-cyan-500/20 text-cyan-300 px-1 py-0.2 rounded font-extrabold font-mono border border-cyan-500/40">
                    1923
                  </span>
                </div>
                <div className="text-[9px] text-cyan-400 font-bold tracking-tight uppercase font-mono">
                  Baseload Hydro
                </div>
              </div>
            </Link>

            <div className="h-8 w-px bg-slate-800 hidden sm:block" />

            <div>
              <div className="flex items-center space-x-1.5">
                <h1 className="font-bold tracking-wider text-slate-100 text-xs sm:text-sm uppercase">
                  Generation Department &amp; Master Station Hub
                </h1>
                <span className="px-1.5 py-0.2 rounded text-[9px] font-bold bg-cyan-500/10 text-cyan-400 border border-cyan-500/20">
                  Unit S05
                </span>
              </div>
              <p className="text-[10px] text-slate-400 font-mono tracking-tight hidden sm:block">
                Cascade Hydroelectric Power Stations 1–6 &bull; NERC Grid Code Operations
              </p>
            </div>
          </div>

          {/* Middle: Live Ticking WAT Clock & Frequency */}
          <div className="hidden xl:flex items-center space-x-3 bg-slate-950/70 border border-slate-800 px-3 py-1.5 rounded-md text-slate-300 font-mono text-[11px]">
            <div className="flex items-center space-x-1.5">
              <Clock className="w-3.5 h-3.5 text-cyan-400 animate-pulse" />
              <span>{dateStr}</span>
              <span className="text-slate-600">|</span>
              <span className="text-emerald-400 font-bold">{timeStr}</span>
            </div>
            <span className="text-slate-600">&bull;</span>
            <div className="flex items-center space-x-1">
              <span className="text-slate-400">Total MW:</span>
              <span className="text-cyan-300 font-bold">{totalMW.toFixed(1)} MW</span>
            </div>
            <span className="text-slate-600">&bull;</span>
            <div className="flex items-center space-x-1">
              <span className="text-slate-400">Freq:</span>
              <span className="text-emerald-400 font-bold">50.02 Hz</span>
            </div>
          </div>

          {/* Right: Navigation to General Landing & General Analytics + Operator Badge & Co-Pilot */}
          <div className="flex items-center space-x-2.5">
            {/* Back to General Landing Page Button */}
            <Link
              href="/"
              className="flex items-center space-x-1.5 bg-slate-850 hover:bg-slate-750 text-slate-200 hover:text-white border border-slate-700 px-3 py-1.5 rounded-lg text-xs font-semibold transition shadow-sm"
              title="Return to General Landing Page"
            >
              <ArrowLeft className="w-3.5 h-3.5 text-cyan-400" />
              <span className="hidden sm:inline">General Landing Page</span>
              <span className="sm:hidden">Back</span>
            </Link>

            {/* General Analytics Page Button */}
            <Link
              href="/portal/central-analytics"
              className="flex items-center space-x-1.5 bg-emerald-950/80 hover:bg-emerald-900/90 text-emerald-300 border border-emerald-500/50 px-3 py-1.5 rounded-lg font-semibold text-xs transition shadow-sm"
              title="View General Company Analytics & Daily Storyboard"
            >
              <BarChart3 className="w-3.5 h-3.5 text-emerald-400" />
              <span className="hidden sm:inline">General Analytics Page</span>
              <span className="sm:hidden">Analytics</span>
            </Link>

            <div className="h-6 w-px bg-slate-800 hidden md:block" />

            <div className="text-right hidden lg:block">
              <div className="text-[11px] font-bold text-slate-200">{staffName}</div>
              <div className="text-[10px] text-cyan-400 font-mono">{staffId} &bull; Master Operator</div>
            </div>

            <button
              onClick={() => setIsChatOpen(!isChatOpen)}
              className="px-3 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 border border-slate-700 text-cyan-300 text-xs font-semibold flex items-center space-x-1.5 transition"
            >
              <Sparkles className="w-3.5 h-3.5 text-cyan-400" />
              <span>@GridAssistant</span>
            </button>
          </div>
        </header>

        {/* Main 2-Column Operational Layout */}
        <div className="flex-1 flex flex-col md:flex-row w-full overflow-hidden">
          {/* Left-Hand Navigation Sidebar Pane (Endpoints to Core Activities) */}
          <GenerationNav_LeftSidebarPane
            activeTask={activeTask}
            onSelectTask={(task) => {
              setActiveTask(task);
              router.push(
                { pathname: "/portal/generation-master", query: { view: task } },
                undefined,
                { shallow: true }
              );
            }}
            activePlantsCount={activePlantsCount}
            openFaultsCount={openFaultsCount}
            pendingRequisitionsCount={pendingRequisitionsCount}
          />

          {/* Active Operating View Workspace - Independent Scrollable Area */}
          <main className="flex-1 p-4 sm:p-6 overflow-y-auto min-w-0">
            {/* Activity 1: Real-Time Analytics & Database Summary */}
            {activeTask === "analytics" && (
              <GenerationAnalyticsView
                plants={plants}
                analytics={analytics}
                onRefresh={loadData}
                onNavigateToReadingSheet={() => {
                  setActiveTask("reading-sheet");
                  router.push({ pathname: "/portal/generation-master", query: { view: "reading-sheet" } }, undefined, { shallow: true });
                }}
              />
            )}

            {/* Activity 2: Hourly Reading Sheet Form */}
            {activeTask === "reading-sheet" && (
              <ReadingSheetForm
                plants={plants}
                staffId={staffId}
                staffName={staffName}
                onSuccess={loadData}
              />
            )}

            {/* Activity 3: Change of Shift Form (Immutable Notes) */}
            {activeTask === "change-of-shift" && (
              <ChangeOfShiftForm
                plants={plants}
                staffId={staffId}
                staffName={staffName}
                onSuccess={loadData}
              />
            )}

            {/* Activity 4: Supervisor Requisition Section */}
            {activeTask === "requisition" && (
              <SupervisorRequisitionView
                plants={plants}
                staffId={staffId}
                staffName={staffName}
                onSuccess={loadData}
              />
            )}

            {/* Auxiliary View: 6-Plant Cascade Fleet Board */}
            {activeTask === "fleet-board" && (
              <PlantFleetLiveBoard
                plants={plants}
                onSelectPlant={(p) => setSelectedPlant(p)}
                onNavigateToCollation={() => setActiveTask("reading-sheet")}
                onNavigateToFaultReport={() => setActiveTask("requisition")}
              />
            )}

            {/* Auxiliary View: Central Machine Workshop */}
            {activeTask === "workshop" && (
              <WorkshopFabricationView
                faults={faults}
                onSuccess={loadData}
              />
            )}
          </main>
        </div>

        {/* AI Co-Pilot Dock */}
        <GridAssistantDock
          isOpen={isChatOpen}
          onClose={() => setIsChatOpen(false)}
          staffRole={claims?.primary_role || "MASTER_STATION_OPERATOR"}
        />

        {/* WorkNotes Drawer */}
        {selectedPlant && claims && (
          <WorkNotesDrawer
            isOpen={true}
            onClose={() => setSelectedPlant(null)}
            entityType="plant_log"
            entityId={selectedPlant.id}
            authorStaffId={claims.staff_id}
            authorRole={claims.primary_role}
          />
        )}
      </div>
    </>
  );
}
