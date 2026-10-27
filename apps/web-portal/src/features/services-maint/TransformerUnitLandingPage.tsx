/**
 * Transformer Maintenance Unit Landing Page.
 * Customer Care 5-Workspace Layout Standard:
 * Thin Top Header Bar + Persistent Left Navigation Sidebar + Modular Workspaces.
 */

import React, { useState, useEffect } from "react";
import Head from "next/head";
import { useRouter } from "next/router";
import {
  TransformerNav_LeftSidebarPane,
  TransformerWorkspace,
} from "./TransformerNav_LeftSidebarPane";
import { TransformerHeaderBar } from "./TransformerHeaderBar";
import { TransformerAnalyticsView } from "./TransformerAnalyticsView";
import { TransformerDGAView } from "./TransformerDGAView";
import { TransformerOverhaulsView } from "./TransformerOverhaulsView";
import { TransformerProofOfWorkView } from "./TransformerProofOfWorkView";
import { TransformerStandardsView } from "./TransformerStandardsView";
import { GridAssistantDock } from "../../components/shell/GridAssistantDock";
import { apiClient } from "../../lib/api-client";
import { parseTokenClaims } from "../../lib/auth";
import { AuthClaims } from "../../types/auth";

export default function TransformerUnitLandingPage() {
  const router = useRouter();
  const [claims, setClaims] = useState<AuthClaims | null>(null);
  const [activeWorkspace, setActiveWorkspace] = useState<TransformerWorkspace>("analytics");
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
      workspace === "dga-testing" ||
      workspace === "overhauls" ||
      workspace === "proof-of-work" ||
      workspace === "standards"
    ) {
      setActiveWorkspace(workspace as TransformerWorkspace);
    }
  }, [router.query]);

  const handleSelectWorkspace = (ws: TransformerWorkspace) => {
    setActiveWorkspace(ws);
    router.push(`/portal/services/transformer?workspace=${ws}`, undefined, { shallow: true });
  };

  return (
    <>
      <Head>
        <title>NESCO Services &bull; Transformer Maintenance Unit</title>
      </Head>

      <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans selection:bg-amber-500/30">
        {/* Thin Top Header Bar */}
        <TransformerHeaderBar
          claims={claims}
          onToggleChat={() => setIsChatOpen(!isChatOpen)}
          activeDtsCount={142}
          criticalOilAlertsCount={1}
        />

        {/* 2-Column Standard: Persistent Left Navigation Sidebar + Modular Workspace */}
        <div className="flex-1 flex overflow-hidden">
          {/* Persistent Left Sidebar */}
          <TransformerNav_LeftSidebarPane
            activeWorkspace={activeWorkspace}
            onSelectWorkspace={handleSelectWorkspace}
            activeDtsCount={142}
            criticalOilAlertsCount={1}
            activeOverhaulsCount={2}
          />

          {/* Main Operational Workspace */}
          <main className="flex-1 overflow-y-auto p-4 sm:p-6 max-w-7xl w-full mx-auto">
            {activeWorkspace === "analytics" && <TransformerAnalyticsView />}
            {activeWorkspace === "dga-testing" && <TransformerDGAView />}
            {activeWorkspace === "overhauls" && <TransformerOverhaulsView />}
            {activeWorkspace === "proof-of-work" && <TransformerProofOfWorkView />}
            {activeWorkspace === "standards" && <TransformerStandardsView />}
          </main>
        </div>

        {/* Omniscient AI Assistant Dock */}
        <GridAssistantDock
          isOpen={isChatOpen}
          onClose={() => setIsChatOpen(false)}
          staffRole={claims?.primary_role || "UNIT_HEAD_TRANSFORMER_MAINT"}
        />
      </div>
    </>
  );
}

