/**
 * Next.js Page: Services & Maintenance Department Portals (/portal/services/[unit]).
 * Dynamic shell for Feeder, Transformer, and Metering Unit Heads.
 */

import React, { useState, useEffect } from "react";
import Head from "next/head";
import { useRouter } from "next/router";
import { TopCommandBar } from "../../../components/shell/TopCommandBar";
import { GridAssistantDock } from "../../../components/shell/GridAssistantDock";
import { WorkNotesDrawer } from "../../../components/shell/WorkNotesDrawer";
import { ZoneB_DashboardRibbon } from "../../../features/services-maint/ZoneB_DashboardRibbon";
import { ZoneC_ActionQueue } from "../../../features/services-maint/ZoneC_ActionQueue";
import { ZoneC_IntakeForm } from "../../../features/services-maint/ZoneC_IntakeForm";
import { fetchWorkOrders, verifyProofOfWork, WorkOrderDTO } from "../../../features/services-maint/api";
import { parseTokenClaims } from "../../../lib/auth";
import { AuthClaims } from "../../../types/auth";
import { apiClient } from "../../../lib/api-client";

export default function ServicesUnitPage() {
  const router = useRouter();
  const { unit } = router.query;
  const unitStr = typeof unit === "string" ? unit : "feeder";

  const [claims, setClaims] = useState<AuthClaims | null>(null);
  const [orders, setOrders] = useState<WorkOrderDTO[]>([]);
  const [isChatOpen, setIsChatOpen] = useState(false);
  const [selectedOrder, setSelectedOrder] = useState<WorkOrderDTO | null>(null);

  const unitTypeMap: Record<string, string> = {
    feeder: "FEEDER_MAINT",
    transformer: "TRANSFORMER_MAINT",
    metering: "METERING_SERVICES",
  };

  useEffect(() => {
    const token = apiClient.getToken();
    if (token) setClaims(parseTokenClaims(token));

    const loadOrders = async () => {
      const targetUnitType = unitTypeMap[unitStr] || "FEEDER_MAINT";
      const data = await fetchWorkOrders(targetUnitType);
      setOrders(data);
    };
    if (unit) loadOrders();
  }, [unit, unitStr]);

  const handleVerifyProof = async (orderId: string) => {
    if (!claims) return;
    try {
      await verifyProofOfWork(orderId, true, "Unit Head verified GPS within tolerance and approved before/after photos.", claims.staff_id);
      setOrders((prev) =>
        prev.map((o) => (o.id === orderId ? { ...o, status: "ATTENDED_TO" } : o))
      );
      alert("Proof-of-work verified! Work order marked ATTENDED_TO and reverted to Customer Care.");
    } catch (err: any) {
      alert(err.message || "Failed to verify proof");
    }
  };

  const title = `Services & Maintenance &bull; ${unitStr.toUpperCase()} UNIT`;

  return (
    <>
      <Head>
        <title>{title}</title>
      </Head>

      <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col">
        <TopCommandBar
          claims={claims}
          departmentTitle={`Maintenance &bull; ${unitStr.toUpperCase()}`}
          onToggleChat={() => setIsChatOpen(!isChatOpen)}
        />

        <main className="flex-1 p-6 space-y-6 max-w-7xl w-full mx-auto">
          <ZoneB_DashboardRibbon unitTitle={`${unitStr.toUpperCase()} Unit`} />

          <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 min-h-[500px]">
            <div className="lg:col-span-8">
              <ZoneC_ActionQueue
                orders={orders}
                onSelectOrder={(o) => setSelectedOrder(o)}
                onVerifyProof={handleVerifyProof}
              />
            </div>
            <div className="lg:col-span-4">
              <ZoneC_IntakeForm
                staffId={claims?.staff_id || "STAFF-MAINT-01"}
                unitType={unitTypeMap[unitStr] || "FEEDER_MAINT"}
                onSuccess={async () => {
                  const targetUnitType = unitTypeMap[unitStr] || "FEEDER_MAINT";
                  const updated = await fetchWorkOrders(targetUnitType);
                  setOrders(updated);
                }}
              />
            </div>
          </div>
        </main>

        <GridAssistantDock
          isOpen={isChatOpen}
          onClose={() => setIsChatOpen(false)}
          staffRole={claims?.primary_role || "UNIT_HEAD_FEEDER_MAINT"}
        />

        {selectedOrder && claims && (
          <WorkNotesDrawer
            isOpen={true}
            onClose={() => setSelectedOrder(null)}
            entityType="work_order"
            entityId={selectedOrder.id}
            authorStaffId={claims.staff_id}
            authorRole={claims.primary_role}
          />
        )}
      </div>
    </>
  );
}

