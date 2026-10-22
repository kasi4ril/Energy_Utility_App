/**
 * Next.js Page: Generation Department & Master Station Command Hub (/portal/generation-master).
 * Enterprise Multi-Task Portal with Mission Header, Task Navigation, and Modular Workspaces.
 */

import React, { useState, useEffect } from "react";
import Head from "next/head";
import { TopCommandBar } from "../../components/shell/TopCommandBar";
import { GridAssistantDock } from "../../components/shell/GridAssistantDock";
import { WorkNotesDrawer } from "../../components/shell/WorkNotesDrawer";
import { GenerationMasterHeader } from "../../features/generation-master/GenerationMasterHeader";
import {
  GenerationMasterNavigation,
} from "../../features/generation-master/GenerationMasterNavigation";
import { PlantFleetLiveBoard } from "../../features/generation-master/PlantFleetLiveBoard";
import { ShiftCollationDesk } from "../../features/generation-master/ShiftCollationDesk";
import { PlantFaultDesk } from "../../features/generation-master/PlantFaultDesk";
import { WorkshopFabricationView } from "../../features/generation-master/WorkshopFabricationView";
import { HydrologyAnalyticsView } from "../../features/generation-master/HydrologyAnalyticsView";
import {
  fetchPlantFleet,
  fetchShiftRollcalls,
  fetchPlantFaults,
  fetchGenerationAnalytics,
} from "../../features/generation-master/api";
import {
  GenerationTask,
  HydroPlantSummary,
  ShiftRollcallRecord,
  PlantFaultJob,
  GenerationAnalytics,
} from "../../features/generation-master/types";
import { parseTokenClaims } from "../../lib/auth";
import { AuthClaims } from "../../types/auth";
import { apiClient } from "../../lib/api-client";

export default function GenerationMasterPage() {
  const [claims, setClaims] = useState<AuthClaims | null>(null);
  const [activeTask, setActiveTask] = useState<GenerationTask>("fleet-board");
  const [plants, setPlants] = useState<HydroPlantSummary[]>([]);
  const [rollcalls, setRollcalls] = useState<ShiftRollcallRecord[]>([]);
  const [faults, setFaults] = useState<PlantFaultJob[]>([]);
  const [analytics, setAnalytics] = useState<GenerationAnalytics | null>(null);
  const [selectedPlant, setSelectedPlant] = useState<HydroPlantSummary | null>(null);
  const [targetCollationPlantId, setTargetCollationPlantId] = useState<string>("PLANT_01_KURRA_FALLS");
  const [isChatOpen, setIsChatOpen] = useState(false);

  const loadData = async () => {
    const [pFleet, rList, fList, aData] = await Promise.all([
      fetchPlantFleet(),
      fetchShiftRollcalls(),
      fetchPlantFaults(),
      fetchGenerationAnalytics(),
    ]);
    setPlants(pFleet);
    setRollcalls(rList);
    setFaults(fList);
    setAnalytics(aData);
  };

  useEffect(() => {
    const token = apiClient.getToken();
    if (token) setClaims(parseTokenClaims(token));
    loadData();
  }, []);

  const totalMW = plants.reduce((sum, p) => sum + p.current_mw, 0);
  const activeUnits = plants.reduce((sum, p) => sum + p.active_units, 0);
  const totalUnits = plants.reduce((sum, p) => sum + p.total_units, 0);
  const activePlantsCount = plants.filter((p) => p.state === "IN_SERVICE").length;
  const openFaultsCount = faults.filter((f) => f.status !== "RESOLVED_REINSTALLED").length;
  const activeWorkshopJobsCount = faults.filter(
    (f) => f.requires_workshop && f.workshop_status !== "COMPLETED_DELIVERED"
  ).length;

  return (
    <>
      <Head>
        <title>NESCO Generation &bull; Master Station Command Hub (Plants 1–6)</title>
      </Head>

      <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans">
        <TopCommandBar
          claims={claims}
          departmentTitle="Generation & Master Station Command Hub"
          onToggleChat={() => setIsChatOpen(!isChatOpen)}
        />

        <main className="flex-1 p-4 sm:p-6 space-y-6 max-w-7xl w-full mx-auto">
          {/* Executive Overview & Regulatory Mission */}
          <GenerationMasterHeader
            totalMW={totalMW}
            inServiceCount={activeUnits}
            totalUnitsCount={totalUnits}
            onRefresh={loadData}
          />

          {/* Operational Task Navigation Pane */}
          <GenerationMasterNavigation
            activeTask={activeTask}
            onSelectTask={(task) => setActiveTask(task)}
            activePlantsCount={activePlantsCount}
            openFaultsCount={openFaultsCount}
            activeWorkshopJobsCount={activeWorkshopJobsCount}
          />

          {/* Dedicated Task Workspace Views */}
          <div className="transition-all duration-200">
            {activeTask === "fleet-board" && (
              <PlantFleetLiveBoard
                plants={plants}
                onSelectPlant={(p) => setSelectedPlant(p)}
                onNavigateToCollation={(pId) => {
                  setTargetCollationPlantId(pId);
                  setActiveTask("call-collation");
                }}
                onNavigateToFaultReport={() => {
                  setActiveTask("plant-faults");
                }}
              />
            )}

            {activeTask === "call-collation" && (
              <ShiftCollationDesk
                plants={plants}
                rollcalls={rollcalls}
                initialPlantId={targetCollationPlantId}
                onSuccess={loadData}
              />
            )}

            {activeTask === "plant-faults" && (
              <PlantFaultDesk
                plants={plants}
                faults={faults}
                onSuccess={loadData}
              />
            )}

            {activeTask === "workshop" && (
              <WorkshopFabricationView
                faults={faults}
                onSuccess={loadData}
              />
            )}

            {activeTask === "analytics" && analytics && (
              <HydrologyAnalyticsView
                analytics={analytics}
                onRefresh={loadData}
              />
            )}
          </div>
        </main>

        <GridAssistantDock
          isOpen={isChatOpen}
          onClose={() => setIsChatOpen(false)}
          staffRole={claims?.primary_role || "MASTER_STATION_OPERATOR"}
        />

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
