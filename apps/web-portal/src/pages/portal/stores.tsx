/**
 * Next.js Page: Store Stock Ledger & Dual-Path Requisitions (/portal/stores).
 * Customer Care 5-Workspace Layout Standard:
 * Thin Top Header Bar + Persistent Left Navigation Sidebar + Modular Workspaces.
 */

import React, { useState, useEffect } from "react";
import Head from "next/head";
import { useRouter } from "next/router";
import { StoresWorkspace } from "../../features/stores-audit/storesTypes";
import { StoresHeaderBar } from "../../features/stores-audit/StoresHeaderBar";
import { StoresNav_LeftSidebarPane } from "../../features/stores-audit/StoresNav_LeftSidebarPane";
import { StoresAnalyticsView } from "../../features/stores-audit/StoresAnalyticsView";
import { StoresRequisitionsView } from "../../features/stores-audit/StoresRequisitionsView";
import { StoresWaybillsView } from "../../features/stores-audit/StoresWaybillsView";
import { StoresScrapReturnsView } from "../../features/stores-audit/StoresScrapReturnsView";
import { StoresClaimsCheckView } from "../../features/stores-audit/StoresClaimsCheckView";
import { GridAssistantDock } from "../../components/shell/GridAssistantDock";
import { apiClient } from "../../lib/api-client";
import { parseTokenClaims } from "../../lib/auth";
import { AuthClaims } from "../../types/auth";

export default function StoresPortalPage() {
  const router = useRouter();
  const [claims, setClaims] = useState<AuthClaims | null>(null);
  const [activeWorkspace, setActiveWorkspace] = useState<StoresWorkspace>("analytics");
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
      workspace === "requisitions" ||
      workspace === "waybills" ||
      workspace === "scrap-returns" ||
      workspace === "claims"
    ) {
      setActiveWorkspace(workspace as StoresWorkspace);
    }
  }, [router.query]);

  const handleSelectWorkspace = (ws: StoresWorkspace) => {
    setActiveWorkspace(ws);
    router.push(`/portal/stores?workspace=${ws}`, undefined, { shallow: true });
  };

  return (
    <>
      <Head>
        <title>NESCO Stores &bull; Stock Ledger &amp; Material Logistics</title>
      </Head>

      <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans selection:bg-orange-500/30">
        {/* Thin Top Header Bar */}
        <StoresHeaderBar
          claims={claims}
          onToggleChat={() => setIsChatOpen(!isChatOpen)}
          inventoryValuationNgn={48500000}
          lowStockItemsCount={2}
        />

        {/* 2-Column Standard: Persistent Left Navigation Sidebar + Modular Workspace */}
        <div className="flex-1 flex overflow-hidden">
          {/* Persistent Left Sidebar */}
          <StoresNav_LeftSidebarPane
            activeWorkspace={activeWorkspace}
            onSelectWorkspace={handleSelectWorkspace}
            pendingRequisitionsCount={1}
            lowStockItemsCount={2}
            pendingClaimsCount={1}
          />

          {/* Main Operational Workspace */}
          <main className="flex-1 overflow-y-auto p-4 sm:p-6 max-w-7xl w-full mx-auto">
            {activeWorkspace === "analytics" && <StoresAnalyticsView />}
            {activeWorkspace === "requisitions" && <StoresRequisitionsView />}
            {activeWorkspace === "waybills" && <StoresWaybillsView />}
            {activeWorkspace === "scrap-returns" && <StoresScrapReturnsView />}
            {activeWorkspace === "claims" && <StoresClaimsCheckView />}
          </main>
        </div>

        {/* Omniscient AI Assistant Dock */}
        <GridAssistantDock
          isOpen={isChatOpen}
          onClose={() => setIsChatOpen(false)}
          staffRole={claims?.primary_role || "HEAD_CENTRAL_STORES"}
        />
      </div>
    </>
  );
}
