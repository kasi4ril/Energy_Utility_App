/**
 * Next.js Page: Energy Management & Central Dispatch (/portal/energy-management).
 * 4-Zone Enterprise Shell for Grid Control, Load Shedding, and 5-Step Safety PTW.
 */

import React, { useState, useEffect } from "react";
import Head from "next/head";
import { TopCommandBar } from "../../components/shell/TopCommandBar";
import { GridAssistantDock } from "../../components/shell/GridAssistantDock";
import { WorkNotesDrawer } from "../../components/shell/WorkNotesDrawer";
import { ZoneB_DashboardRibbon } from "../../features/energy-management/ZoneB_DashboardRibbon";
import { ZoneC_ActionQueue } from "../../features/energy-management/ZoneC_ActionQueue";
import { ZoneC_IntakeForm } from "../../features/energy-management/ZoneC_IntakeForm";
import {
  fetchFeeders,
  fetchPTWPermits,
  FeederStatusDTO,
  SafetyPermitDTO,
} from "../../features/energy-management/api";
import { parseTokenClaims } from "../../lib/auth";
import { AuthClaims } from "../../types/auth";
import { apiClient } from "../../lib/api-client";

export default function EnergyManagementPage() {
  const [claims, setClaims] = useState<AuthClaims | null>(null);
  const [feeders, setFeeders] = useState<FeederStatusDTO[]>([]);
  const [permits, setPermits] = useState<SafetyPermitDTO[]>([]);
  const [isChatOpen, setIsChatOpen] = useState(false);
  const [selectedEntity, setSelectedEntity] = useState<{ type: string; id: string } | null>(null);

  useEffect(() => {
    const token = apiClient.getToken();
    if (token) setClaims(parseTokenClaims(token));

    const loadData = async () => {
      const f = await fetchFeeders();
      setFeeders(f);
      const p = await fetchPTWPermits();
      setPermits(p);
    };
    loadData();
  }, []);

  return (
    <>
      <Head>
        <title>NESCO Central Dispatch &bull; Energy Management</title>
      </Head>

      <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col">
        <TopCommandBar
          claims={claims}
          departmentTitle="Central Dispatch & Grid Control"
          onToggleChat={() => setIsChatOpen(!isChatOpen)}
        />

        <main className="flex-1 p-6 space-y-6 max-w-7xl w-full mx-auto">
          <ZoneB_DashboardRibbon />

          <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 min-h-[500px]">
            <div className="lg:col-span-8">
              <ZoneC_ActionQueue
                feeders={feeders}
                permits={permits}
                onSelectFeeder={(f) => setSelectedEntity({ type: "feeder", id: f.id })}
                onSelectPermit={(p) => setSelectedEntity({ type: "ptw_permit", id: p.id })}
              />
            </div>
            <div className="lg:col-span-4">
              <ZoneC_IntakeForm
                onSuccess={async () => {
                  const updated = await fetchFeeders();
                  setFeeders(updated);
                }}
              />
            </div>
          </div>
        </main>

        <GridAssistantDock
          isOpen={isChatOpen}
          onClose={() => setIsChatOpen(false)}
          staffRole={claims?.primary_role || "CENTRAL_DISPATCH_OPERATOR"}
        />

        {selectedEntity && claims && (
          <WorkNotesDrawer
            isOpen={true}
            onClose={() => setSelectedEntity(null)}
            entityType={selectedEntity.type}
            entityId={selectedEntity.id}
            authorStaffId={claims.staff_id}
            authorRole={claims.primary_role}
          />
        )}
      </div>
    </>
  );
}

