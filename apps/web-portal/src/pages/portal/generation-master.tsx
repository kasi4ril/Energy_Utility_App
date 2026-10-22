/**
 * Next.js Page: Generation Master Station Hub (/portal/generation-master).
 * 4-Zone Enterprise Shell for Plants 1–6 Operations & Hydrological Monitoring.
 */

import React, { useState, useEffect } from "react";
import Head from "next/head";
import { TopCommandBar } from "../../components/shell/TopCommandBar";
import { GridAssistantDock } from "../../components/shell/GridAssistantDock";
import { WorkNotesDrawer } from "../../components/shell/WorkNotesDrawer";
import { ZoneB_DashboardRibbon } from "../../features/generation-master/ZoneB_DashboardRibbon";
import { ZoneC_ActionQueue } from "../../features/generation-master/ZoneC_ActionQueue";
import { ZoneC_IntakeForm } from "../../features/generation-master/ZoneC_IntakeForm";
import { fetchPlantTelemetry, PlantTelemetryDTO } from "../../features/generation-master/api";
import { parseTokenClaims } from "../../lib/auth";
import { AuthClaims } from "../../types/auth";
import { apiClient } from "../../lib/api-client";

export default function GenerationMasterPage() {
  const [claims, setClaims] = useState<AuthClaims | null>(null);
  const [logs, setLogs] = useState<PlantTelemetryDTO[]>([]);
  const [isChatOpen, setIsChatOpen] = useState(false);
  const [selectedLog, setSelectedLog] = useState<PlantTelemetryDTO | null>(null);

  useEffect(() => {
    const token = apiClient.getToken();
    if (token) setClaims(parseTokenClaims(token));

    const loadTelemetry = async () => {
      const data = await fetchPlantTelemetry();
      setLogs(data);
    };
    loadTelemetry();
  }, []);

  return (
    <>
      <Head>
        <title>NESCO Generation Master Hub &bull; Plants 1–6</title>
      </Head>

      <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col">
        <TopCommandBar
          claims={claims}
          departmentTitle="Generation Master Station (Plants 1–6)"
          onToggleChat={() => setIsChatOpen(!isChatOpen)}
        />

        <main className="flex-1 p-6 space-y-6 max-w-7xl w-full mx-auto">
          <ZoneB_DashboardRibbon />

          <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 min-h-[500px]">
            <div className="lg:col-span-8">
              <ZoneC_ActionQueue
                logs={logs}
                onSelectLog={(log) => setSelectedLog(log)}
              />
            </div>
            <div className="lg:col-span-4">
              <ZoneC_IntakeForm
                onSuccess={async () => {
                  const updated = await fetchPlantTelemetry();
                  setLogs(updated);
                }}
              />
            </div>
          </div>
        </main>

        <GridAssistantDock
          isOpen={isChatOpen}
          onClose={() => setIsChatOpen(false)}
          staffRole={claims?.primary_role || "MASTER_STATION_OPERATOR"}
        />

        {selectedLog && claims && (
          <WorkNotesDrawer
            isOpen={true}
            onClose={() => setSelectedLog(null)}
            entityType="plant_log"
            entityId={selectedLog.id}
            authorStaffId={claims.staff_id}
            authorRole={claims.primary_role}
          />
        )}
      </div>
    </>
  );
}

