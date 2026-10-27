/**
 * Next.js Page: Internal Audit Gate & 3-Way Claims Desk (/portal/internal-audit).
 * Customer Care 5-Workspace Layout Standard:
 * Thin Top Header Bar + Persistent Left Navigation Sidebar + Modular Workspaces.
 */

import React, { useState, useEffect } from "react";
import Head from "next/head";
import { useRouter } from "next/router";
import { AuditWorkspace } from "../../features/stores-audit/auditTypes";
import { AuditHeaderBar } from "../../features/stores-audit/AuditHeaderBar";
import { AuditNav_LeftSidebarPane } from "../../features/stores-audit/AuditNav_LeftSidebarPane";
import { AuditAnalyticsView } from "../../features/stores-audit/AuditAnalyticsView";
import { AuditLedgerView } from "../../features/stores-audit/AuditLedgerView";
import { AuditClaimsGateView } from "../../features/stores-audit/AuditClaimsGateView";
import { AuditScrapClearanceView } from "../../features/stores-audit/AuditScrapClearanceView";
import { AuditComplianceView } from "../../features/stores-audit/AuditComplianceView";
import { GridAssistantDock } from "../../components/shell/GridAssistantDock";
import { apiClient } from "../../lib/api-client";
import { parseTokenClaims } from "../../lib/auth";
import { AuthClaims } from "../../types/auth";

export default function InternalAuditPage() {
  const router = useRouter();
  const [claims, setClaims] = useState<AuthClaims | null>(null);
  const [activeWorkspace, setActiveWorkspace] = useState<AuditWorkspace>("analytics");
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
      workspace === "ledger" ||
      workspace === "claims-audit" ||
      workspace === "scrap-audit" ||
      workspace === "compliance"
    ) {
      setActiveWorkspace(workspace as AuditWorkspace);
    }
  }, [router.query]);

  const handleSelectWorkspace = (ws: AuditWorkspace) => {
    setActiveWorkspace(ws);
    router.push(`/portal/internal-audit?workspace=${ws}`, undefined, { shallow: true });
  };

  return (
    <>
      <Head>
        <title>NESCO Internal Audit &bull; Compliance &amp; Risk Governance</title>
      </Head>

      <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans selection:bg-rose-500/30">
        {/* Thin Top Header Bar */}
        <AuditHeaderBar
          claims={claims}
          onToggleChat={() => setIsChatOpen(!isChatOpen)}
          complianceScorePct={96.8}
          openObservationsCount={1}
        />

        {/* 2-Column Standard: Persistent Left Navigation Sidebar + Modular Workspace */}
        <div className="flex-1 flex overflow-hidden">
          {/* Persistent Left Sidebar */}
          <AuditNav_LeftSidebarPane
            activeWorkspace={activeWorkspace}
            onSelectWorkspace={handleSelectWorkspace}
            openFindingsCount={1}
            pendingClaimsCount={1}
            unauthorizedScrapCount={0}
          />

          {/* Main Operational Workspace */}
          <main className="flex-1 overflow-y-auto p-4 sm:p-6 max-w-7xl w-full mx-auto">
            {activeWorkspace === "analytics" && <AuditAnalyticsView />}
            {activeWorkspace === "ledger" && <AuditLedgerView />}
            {activeWorkspace === "claims-audit" && <AuditClaimsGateView />}
            {activeWorkspace === "scrap-audit" && <AuditScrapClearanceView />}
            {activeWorkspace === "compliance" && <AuditComplianceView />}
          </main>
        </div>

        {/* Omniscient AI Assistant Dock */}
        <GridAssistantDock
          isOpen={isChatOpen}
          onClose={() => setIsChatOpen(false)}
          staffRole={claims?.primary_role || "HEAD_INTERNAL_AUDIT"}
        />
      </div>
    </>
  );
}
