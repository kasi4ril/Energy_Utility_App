/**
 * Next.js Page: Services & Maintenance Department Portals (/portal/services/[unit]).
 * Enterprise Multi-Task Shell for Feeder, Transformer, and Metering Unit Heads.
 */

import React, { useState, useEffect } from "react";
import Head from "next/head";
import { useRouter } from "next/router";
import { TopCommandBar } from "../../../components/shell/TopCommandBar";
import { GridAssistantDock } from "../../../components/shell/GridAssistantDock";
import { WorkNotesDrawer } from "../../../components/shell/WorkNotesDrawer";
import { ServicesMaintHeader } from "../../../features/services-maint/ServicesMaintHeader";
import {
  ServicesMaintNavigation,
  ServicesTask,
} from "../../../features/services-maint/ServicesMaintNavigation";
import { ZoneC_ActionQueue } from "../../../features/services-maint/ZoneC_ActionQueue";
import { ZoneC_IntakeForm } from "../../../features/services-maint/ZoneC_IntakeForm";
import { CustomerCareComplaintFeed } from "../../../features/services-maint/CustomerCareComplaintFeed";
import { fetchWorkOrders, verifyProofOfWork, WorkOrderDTO } from "../../../features/services-maint/api";
import { parseTokenClaims } from "../../../lib/auth";
import { AuthClaims } from "../../../types/auth";
import { apiClient } from "../../../lib/api-client";

export default function ServicesUnitPage() {
  const router = useRouter();
  const { unit } = router.query;
  const unitStr = typeof unit === "string" ? unit : "feeder";

  const [claims, setClaims] = useState<AuthClaims | null>(null);
  const [activeTask, setActiveTask] = useState<ServicesTask>("work-orders");
  const [orders, setOrders] = useState<WorkOrderDTO[]>([]);
  const [isChatOpen, setIsChatOpen] = useState(false);
  const [selectedOrder, setSelectedOrder] = useState<WorkOrderDTO | null>(null);

  const unitTypeMap: Record<string, string> = {
    feeder: "FEEDER_MAINT",
    transformer: "TRANSFORMER_MAINT",
    metering: "METERING_SERVICES",
  };

  const loadOrders = async () => {
    const targetUnitType = unitTypeMap[unitStr] || "FEEDER_MAINT";
    const data = await fetchWorkOrders(targetUnitType);
    setOrders(data);
  };

  useEffect(() => {
    const token = apiClient.getToken();
    if (token) setClaims(parseTokenClaims(token));
    if (unit) loadOrders();
  }, [unit, unitStr]);

  const handleVerifyProof = async (orderId: string) => {
    if (!claims) return;
    try {
      await verifyProofOfWork(
        orderId,
        true,
        "Unit Head verified GPS within tolerance and approved before/after photos.",
        claims.staff_id
      );
      setOrders((prev) =>
        prev.map((o) => (o.id === orderId ? { ...o, status: "ATTENDED_TO" } : o))
      );
      alert("Proof-of-work verified! Work order marked ATTENDED_TO and reverted to Customer Care.");
    } catch (err: any) {
      alert(err.message || "Failed to verify proof");
    }
  };

  const unitTitle = `${unitStr.toUpperCase()} UNIT`;

  return (
    <>
      <Head>
        <title>NESCO Services &bull; {unitTitle}</title>
      </Head>

      <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans">
        <TopCommandBar
          claims={claims}
          departmentTitle={`Field Engineering & Maintenance &bull; ${unitTitle}`}
          onToggleChat={() => setIsChatOpen(!isChatOpen)}
        />

        <main className="flex-1 p-4 sm:p-6 space-y-6 max-w-7xl w-full mx-auto">
          {/* Executive Overview & Regulatory Mission */}
          <ServicesMaintHeader
            unitTitle={unitTitle}
            orderCount={orders.length}
            onRefresh={loadOrders}
          />

          {/* Operational Task Navigation Pane */}
          <ServicesMaintNavigation
            activeTask={activeTask}
            onSelectTask={(task) => setActiveTask(task)}
            orderCount={orders.length}
          />

          {/* Dedicated Task Workspace Views */}
          <div className="transition-all duration-200">
            {activeTask === "work-orders" && (
              <ZoneC_ActionQueue
                orders={orders}
                onSelectOrder={(o) => setSelectedOrder(o)}
                onVerifyProof={handleVerifyProof}
              />
            )}

            {activeTask === "create-order" && (
              <div className="max-w-2xl mx-auto">
                <ZoneC_IntakeForm
                  staffId={claims?.staff_id || "STAFF-MAINT-01"}
                  unitType={unitTypeMap[unitStr] || "FEEDER_MAINT"}
                  onSuccess={loadOrders}
                />
              </div>
            )}

            {activeTask === "proof-verification" && (
              <ZoneC_ActionQueue
                orders={orders}
                onSelectOrder={(o) => setSelectedOrder(o)}
                onVerifyProof={handleVerifyProof}
              />
            )}

            {activeTask === "unit-analytics" && (
              <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-6 text-center space-y-3">
                <h3 className="text-sm font-bold text-white">Unit Maintenance Performance &amp; MTTR</h3>
                <p className="text-xs text-slate-400">
                  Average Time to Repair: <span className="font-bold text-emerald-400">2.8 hrs</span> &bull; 
                  Active Crews on Duty: <span className="font-bold text-sky-400">4 Crews</span> &bull;
                  PTW Safety Compliance: <span className="font-bold text-emerald-400">100%</span>
                </p>
              </div>
            )}

            {activeTask === "customer-complaints" && (
              <CustomerCareComplaintFeed
                currentUnit={unitStr}
                onCreateWorkOrderFromTicket={() => setActiveTask("create-order")}
              />
            )}
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
