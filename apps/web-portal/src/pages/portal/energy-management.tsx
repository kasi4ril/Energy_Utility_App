/**
 * Next.js Page: Energy Management & Central Transmission Dispatch (/portal/energy-management).
 * Customer Care 5-Workspace Layout Standard:
 * Thin Top Header Bar + Persistent Left Navigation Sidebar + Modular Workspaces.
 */

import React, { useState, useEffect } from "react";
import Head from "next/head";
import { useRouter } from "next/router";
import { EnergyWorkspace } from "../../features/energy-management/types";
import { EnergyHeaderBar } from "../../features/energy-management/EnergyHeaderBar";
import { EnergyNav_LeftSidebarPane } from "../../features/energy-management/EnergyNav_LeftSidebarPane";
import { EnergyAnalyticsView } from "../../features/energy-management/EnergyAnalyticsView";
import { EnergyLoadSheddingView } from "../../features/energy-management/EnergyLoadSheddingView";
import { EnergySafetyPTWView } from "../../features/energy-management/EnergySafetyPTWView";
import { EnergyDispatchLogView } from "../../features/energy-management/EnergyDispatchLogView";
import { EnergyOutageClearanceView } from "../../features/energy-management/EnergyOutageClearanceView";
import { GridAssistantDock } from "../../components/shell/GridAssistantDock";
import { apiClient } from "../../lib/api-client";
import { parseTokenClaims } from "../../lib/auth";
import { AuthClaims } from "../../types/auth";

export default function EnergyManagementPage() {
  const router = useRouter();
  const [claims, setClaims] = useState<AuthClaims | null>(null);
  const [activeWorkspace, setActiveWorkspace] = useState<EnergyWorkspace>("analytics");
  const [isChatOpen, setIsChatOpen] = useState(false);

  useEffect(() => {
    const token = apiClient.getToken();
    if (token) setClaims(parseTokenClaims(token));
  }, []);

  // Sync workspace with query
  useEffect(() => {
    const { workspace } = router.query;
    if (
      workspace === "analytics" ||
      workspace === "load-shedding" ||
      workspace === "safety-ptw" ||
      workspace === "dispatch-log" ||
      workspace === "outage-clearance"
    ) {
      setActiveWorkspace(workspace as EnergyWorkspace);
    }
  }, [router.query]);

  const handleSelectWorkspace = (ws: EnergyWorkspace) => {
    setActiveWorkspace(ws);
    router.push(`/portal/energy-management?workspace=${ws}`, undefined, { shallow: true });
  };

  return (
    <>
      <Head>
        <title>NESCO Central Dispatch &bull; Energy Management</title>
      </Head>

      <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans selection:bg-cyan-500/30">
        {/* Thin Top Header Bar */}
        <EnergyHeaderBar
          claims={claims}
          onToggleChat={() => setIsChatOpen(!isChatOpen)}
          frequencyHz={50.02}
          bulkIntakeMw={63.6}
        />

        {/* 2-Column Standard: Persistent Left Navigation Sidebar + Modular Workspace */}
        <div className="flex-1 flex overflow-hidden">
          {/* Persistent Left Sidebar */}
          <EnergyNav_LeftSidebarPane
            activeWorkspace={activeWorkspace}
            onSelectWorkspace={handleSelectWorkspace}
            activeFeedersCount={6}
            sheddedFeedersCount={1}
            activePTWCount={1}
            pendingClearancesCount={1}
          />

          {/* Main Operational Workspace */}
          <main className="flex-1 overflow-y-auto p-4 sm:p-6 max-w-7xl w-full mx-auto">
            {activeWorkspace === "analytics" && <EnergyAnalyticsView />}
            {activeWorkspace === "load-shedding" && <EnergyLoadSheddingView />}
            {activeWorkspace === "safety-ptw" && <EnergySafetyPTWView />}
            {activeWorkspace === "dispatch-log" && <EnergyDispatchLogView />}
            {activeWorkspace === "outage-clearance" && <EnergyOutageClearanceView />}
          </main>
        </div>

        {/* Omniscient AI Assistant Dock */}
        <GridAssistantDock
          isOpen={isChatOpen}
          onClose={() => setIsChatOpen(false)}
          staffRole={claims?.primary_role || "CHIEF_DISPATCH_OFFICER"}
        />
      </div>
    </>
  );
}
