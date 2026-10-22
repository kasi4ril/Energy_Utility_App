/**
 * Next.js Page: Internal Audit Gate & 3-Way Claims Desk (/portal/internal-audit).
 * 4-Zone Enterprise Shell for Internal Auditors and Compliance Officers.
 */

import React, { useState, useEffect } from "react";
import Head from "next/head";
import { TopCommandBar } from "../../components/shell/TopCommandBar";
import { GridAssistantDock } from "../../components/shell/GridAssistantDock";
import { WorkNotesDrawer } from "../../components/shell/WorkNotesDrawer";
import { ZoneB_DashboardRibbon } from "../../features/stores-audit/ZoneB_DashboardRibbon";
import { ZoneC_ActionQueue } from "../../features/stores-audit/ZoneC_ActionQueue";
import { ZoneC_IntakeForm } from "../../features/stores-audit/ZoneC_IntakeForm";
import { fetchRequisitions, RequisitionDTO } from "../../features/stores-audit/api";
import { parseTokenClaims } from "../../lib/auth";
import { AuthClaims } from "../../types/auth";
import { apiClient } from "../../lib/api-client";

export default function InternalAuditPage() {
  const [claims, setClaims] = useState<AuthClaims | null>(null);
  const [requisitions, setRequisitions] = useState<RequisitionDTO[]>([]);
  const [isChatOpen, setIsChatOpen] = useState(false);
  const [selectedReq, setSelectedReq] = useState<RequisitionDTO | null>(null);

  useEffect(() => {
    const token = apiClient.getToken();
    if (token) setClaims(parseTokenClaims(token));

    const loadReqs = async () => {
      const data = await fetchRequisitions();
      setRequisitions(data);
    };
    loadReqs();
  }, []);

  return (
    <>
      <Head>
        <title>NESCO Internal Audit &bull; 3-Way Claims Verification</title>
      </Head>

      <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col">
        <TopCommandBar
          claims={claims}
          departmentTitle="Internal Audit & Compliance Desk"
          onToggleChat={() => setIsChatOpen(!isChatOpen)}
        />

        <main className="flex-1 p-6 space-y-6 max-w-7xl w-full mx-auto">
          <ZoneB_DashboardRibbon />

          <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 min-h-[500px]">
            <div className="lg:col-span-8">
              <ZoneC_ActionQueue
                requisitions={requisitions}
                onSelectRequisition={(r) => setSelectedReq(r)}
              />
            </div>
            <div className="lg:col-span-4">
              <ZoneC_IntakeForm
                staffId={claims?.staff_id || "STAFF-AUDIT-01"}
                onSuccess={async () => {
                  const updated = await fetchRequisitions();
                  setRequisitions(updated);
                }}
              />
            </div>
          </div>
        </main>

        <GridAssistantDock
          isOpen={isChatOpen}
          onClose={() => setIsChatOpen(false)}
          staffRole={claims?.primary_role || "INTERNAL_AUDITOR"}
        />

        {selectedReq && claims && (
          <WorkNotesDrawer
            isOpen={true}
            onClose={() => setSelectedReq(null)}
            entityType="requisition"
            entityId={selectedReq.id}
            authorStaffId={claims.staff_id}
            authorRole={claims.primary_role}
          />
        )}
      </div>
    </>
  );
}

