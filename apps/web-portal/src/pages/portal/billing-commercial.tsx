/**
 * Next.js Page: Commercial Operations, Billing & STS Prepayment (/portal/billing-commercial).
 * Customer Care 5-Workspace Layout Standard:
 * Thin Top Header Bar + Persistent Left Navigation Sidebar + Modular Workspaces.
 */

import React, { useState, useEffect } from "react";
import Head from "next/head";
import { useRouter } from "next/router";
import { BillingWorkspace } from "../../features/billing-commercial/types";
import { BillingHeaderBar } from "../../features/billing-commercial/BillingHeaderBar";
import { BillingNav_LeftSidebarPane } from "../../features/billing-commercial/BillingNav_LeftSidebarPane";
import { CommercialAnalyticsView } from "../../features/billing-commercial/CommercialAnalyticsView";
import { CommercialVendingTerminalView } from "../../features/billing-commercial/CommercialVendingTerminalView";
import { CommercialRefundsLedgerView } from "../../features/billing-commercial/CommercialRefundsLedgerView";
import { CommercialDisputeDeskView } from "../../features/billing-commercial/CommercialDisputeDeskView";
import { CommercialTheftAlertsView } from "../../features/billing-commercial/CommercialTheftAlertsView";
import { GridAssistantDock } from "../../components/shell/GridAssistantDock";
import { apiClient } from "../../lib/api-client";
import { parseTokenClaims } from "../../lib/auth";
import { AuthClaims } from "../../types/auth";

export default function BillingCommercialPage() {
  const router = useRouter();
  const [claims, setClaims] = useState<AuthClaims | null>(null);
  const [activeWorkspace, setActiveWorkspace] = useState<BillingWorkspace>("analytics");
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
      workspace === "vending" ||
      workspace === "refunds" ||
      workspace === "disputes" ||
      workspace === "theft-ai"
    ) {
      setActiveWorkspace(workspace as BillingWorkspace);
    }
  }, [router.query]);

  const handleSelectWorkspace = (ws: BillingWorkspace) => {
    setActiveWorkspace(ws);
    router.push(`/portal/billing-commercial?workspace=${ws}`, undefined, { shallow: true });
  };

  return (
    <>
      <Head>
        <title>NESCO Commercial Operations &bull; Billing &amp; STS Vending</title>
      </Head>

      <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans selection:bg-teal-500/30">
        {/* Thin Top Header Bar */}
        <BillingHeaderBar
          claims={claims}
          onToggleChat={() => setIsChatOpen(!isChatOpen)}
          todayCollectionsNgn={4850000}
          collectionEfficiencyPct={94.2}
        />

        {/* 2-Column Standard: Persistent Left Navigation Sidebar + Modular Workspace */}
        <div className="flex-1 flex overflow-hidden">
          {/* Persistent Left Sidebar */}
          <BillingNav_LeftSidebarPane
            activeWorkspace={activeWorkspace}
            onSelectWorkspace={handleSelectWorkspace}
            pendingDisputesCount={1}
            criticalTheftAlertsCount={2}
            activeRefundsCount={2}
          />

          {/* Main Operational Workspace */}
          <main className="flex-1 overflow-y-auto p-4 sm:p-6 max-w-7xl w-full mx-auto">
            {activeWorkspace === "analytics" && <CommercialAnalyticsView />}
            {activeWorkspace === "vending" && <CommercialVendingTerminalView />}
            {activeWorkspace === "refunds" && <CommercialRefundsLedgerView />}
            {activeWorkspace === "disputes" && <CommercialDisputeDeskView />}
            {activeWorkspace === "theft-ai" && <CommercialTheftAlertsView />}
          </main>
        </div>

        {/* Omniscient AI Assistant Dock */}
        <GridAssistantDock
          isOpen={isChatOpen}
          onClose={() => setIsChatOpen(false)}
          staffRole={claims?.primary_role || "HEAD_COMMERCIAL_OPERATIONS"}
        />
      </div>
    </>
  );
}
