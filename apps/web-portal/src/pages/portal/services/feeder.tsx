/**
 * Next.js Page: Services & Maintenance Feeder Unit (/portal/services/feeder).
 * Built with Customer Care Structure: Thin Top Header Bar + Persistent Left-Hand Navigation Sidebar Pane + 5-Point Workspaces.
 * Strictly <= 350 lines.
 */

import React, { useState, useEffect } from "react";
import Head from "next/head";
import { useRouter } from "next/router";
import { FeederHeaderBar } from "../../../features/services-maint/FeederHeaderBar";
import { FeederNav_LeftSidebarPane } from "../../../features/services-maint/FeederNav_LeftSidebarPane";
import { FeederAnalyticsView } from "../../../features/services-maint/FeederAnalyticsView";
import { FeederOutageRequestView } from "../../../features/services-maint/FeederOutageRequestView";
import { FeederProofOfWorkView } from "../../../features/services-maint/FeederProofOfWorkView";
import { FeederStandardsView } from "../../../features/services-maint/FeederStandardsView";
import { FeederWorkOrdersView } from "../../../features/services-maint/FeederWorkOrdersView";
import { GridAssistantDock } from "../../../components/shell/GridAssistantDock";
import { FeederWorkspace } from "../../../features/services-maint/types";
import { MOCK_FEEDER_WORK_ORDERS, MOCK_FEEDER_PROOFS, MOCK_FEEDER_OUTAGE_REQUESTS } from "../../../features/services-maint/mockData";
import { parseTokenClaims } from "../../../lib/auth";
import { AuthClaims } from "../../../types/auth";
import { apiClient } from "../../../lib/api-client";

interface FeederPageProps {
  initialWorkspace?: FeederWorkspace;
}

export default function FeederUnitLandingPage({ initialWorkspace = "analytics" }: FeederPageProps) {
  const router = useRouter();
  const [claims, setClaims] = useState<AuthClaims | null>(null);
  const [activeWorkspace, setActiveWorkspace] = useState<FeederWorkspace>(initialWorkspace);
  const [isChatOpen, setIsChatOpen] = useState(false);

  useEffect(() => {
    const token = apiClient.getToken();
    if (token) setClaims(parseTokenClaims(token));

    if (router.query.view && typeof router.query.view === "string") {
      const qView = router.query.view as FeederWorkspace;
      if (["analytics", "outage-request", "proof-of-work", "standards", "work-orders"].includes(qView)) {
        setActiveWorkspace(qView);
      }
    }
  }, [router.query.view]);

  const handleSelectWorkspace = (ws: FeederWorkspace) => {
    setActiveWorkspace(ws);
    router.push({ pathname: "/portal/services/feeder", query: { view: ws } }, undefined, {
      shallow: true,
    });
  };

  const activeOrdersCount = MOCK_FEEDER_WORK_ORDERS.filter((o) => o.status === "IN_PROGRESS" || o.status === "ASSIGNED").length;
  const pendingProofCount = MOCK_FEEDER_PROOFS.filter((p) => p.verification_status === "PENDING_UNIT_HEAD").length;
  const activeOutageCount = MOCK_FEEDER_OUTAGE_REQUESTS.filter((r) => r.status === "APPROVED_PTW_ISSUED" || r.status === "ISOLATED_IN_PROGRESS").length;

  return (
    <>
      <Head>
        <title>NESCO Feeder Unit &bull; Services & Maintenance</title>
      </Head>

      <div className="h-screen bg-slate-950 text-slate-100 flex flex-col font-sans overflow-hidden">
        {/* 1. Thin Top Header Bar - Static Stand-Still */}
        <FeederHeaderBar
          claims={claims}
          activeTripsCount={1}
          onToggleChat={() => setIsChatOpen(!isChatOpen)}
        />

        {/* 2. Main 2-Column Customer-Care Standard Layout */}
        <div className="flex-1 flex flex-col md:flex-row w-full overflow-hidden">
          {/* Left-Hand Navigation Sidebar Pane - Persistent Static Stand-Still */}
          <FeederNav_LeftSidebarPane
            activeWorkspace={activeWorkspace}
            onSelectWorkspace={handleSelectWorkspace}
            activeOrdersCount={activeOrdersCount}
            pendingProofCount={pendingProofCount}
            activeOutageCount={activeOutageCount}
          />

          {/* Active Operating View Workspace - Independent Scrollable Area */}
          <main className="flex-1 p-4 sm:p-6 overflow-y-auto min-w-0 bg-slate-950/95">
            {activeWorkspace === "analytics" && <FeederAnalyticsView />}
            {activeWorkspace === "outage-request" && <FeederOutageRequestView />}
            {activeWorkspace === "proof-of-work" && <FeederProofOfWorkView />}
            {activeWorkspace === "standards" && <FeederStandardsView />}
            {activeWorkspace === "work-orders" && <FeederWorkOrdersView />}
          </main>
        </div>

        {/* AI Guardian / Assistant Dock */}
        <GridAssistantDock
          isOpen={isChatOpen}
          onClose={() => setIsChatOpen(false)}
          staffRole={claims?.primary_role || "UNIT_HEAD_FEEDER_MAINT"}
        />
      </div>
    </>
  );
}
