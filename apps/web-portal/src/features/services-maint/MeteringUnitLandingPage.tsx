/**
 * Metering Services Unit Landing Page.
 * Customer Care 5-Workspace Layout Standard:
 * Thin Top Header Bar + Persistent Left Navigation Sidebar + Modular Workspaces.
 */

import React, { useState, useEffect } from "react";
import Head from "next/head";
import { useRouter } from "next/router";
import {
  MeteringNav_LeftSidebarPane,
  MeteringWorkspace,
} from "./MeteringNav_LeftSidebarPane";
import { MeteringHeaderBar } from "./MeteringHeaderBar";
import { MeteringAnalyticsView } from "./MeteringAnalyticsView";
import { MeteringKCTTerminalView } from "./MeteringKCTTerminalView";
import { MeteringTestBenchView } from "./MeteringTestBenchView";
import { MeteringTamperInvestigationView } from "./MeteringTamperInvestigationView";
import { MeteringStandardsView } from "./MeteringStandardsView";
import { GridAssistantDock } from "../../components/shell/GridAssistantDock";
import { apiClient } from "../../lib/api-client";
import { parseTokenClaims } from "../../lib/auth";
import { AuthClaims } from "../../types/auth";

export default function MeteringUnitLandingPage() {
  const router = useRouter();
  const [claims, setClaims] = useState<AuthClaims | null>(null);
  const [activeWorkspace, setActiveWorkspace] = useState<MeteringWorkspace>("analytics");
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
      workspace === "kct-terminal" ||
      workspace === "testing-bench" ||
      workspace === "tamper-desk" ||
      workspace === "standards"
    ) {
      setActiveWorkspace(workspace as MeteringWorkspace);
    }
  }, [router.query]);

  const handleSelectWorkspace = (ws: MeteringWorkspace) => {
    setActiveWorkspace(ws);
    router.push(`/portal/services/metering?workspace=${ws}`, undefined, { shallow: true });
  };

  return (
    <>
      <Head>
        <title>NESCO Services &bull; Metering Services Unit</title>
      </Head>

      <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans selection:bg-cyan-500/30">
        {/* Thin Top Header Bar */}
        <MeteringHeaderBar
          claims={claims}
          onToggleChat={() => setIsChatOpen(!isChatOpen)}
          totalMetersFleet={14850}
          activeTamperAlertsCount={2}
        />

        {/* 2-Column Standard: Persistent Left Navigation Sidebar + Modular Workspace */}
        <div className="flex-1 flex overflow-hidden">
          {/* Persistent Left Sidebar */}
          <MeteringNav_LeftSidebarPane
            activeWorkspace={activeWorkspace}
            onSelectWorkspace={handleSelectWorkspace}
            totalMetersFleet={14850}
            pendingKctsCount={3}
            activeTampersCount={2}
          />

          {/* Main Operational Workspace */}
          <main className="flex-1 overflow-y-auto p-4 sm:p-6 max-w-7xl w-full mx-auto">
            {activeWorkspace === "analytics" && <MeteringAnalyticsView />}
            {activeWorkspace === "kct-terminal" && <MeteringKCTTerminalView />}
            {activeWorkspace === "testing-bench" && <MeteringTestBenchView />}
            {activeWorkspace === "tamper-desk" && <MeteringTamperInvestigationView />}
            {activeWorkspace === "standards" && <MeteringStandardsView />}
          </main>
        </div>

        {/* Omniscient AI Assistant Dock */}
        <GridAssistantDock
          isOpen={isChatOpen}
          onClose={() => setIsChatOpen(false)}
          staffRole={claims?.primary_role || "UNIT_HEAD_METERING_SERVICES"}
        />
      </div>
    </>
  );
}

