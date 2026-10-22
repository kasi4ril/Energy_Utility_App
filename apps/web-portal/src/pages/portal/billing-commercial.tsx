/**
 * Next.js Page: Billing & Commercial Operations (/portal/billing-commercial).
 * Enterprise Multi-Task Portal with Mission Header, Task Navigation, and Modular Workspaces.
 */

import React, { useState, useEffect } from "react";
import Head from "next/head";
import { TopCommandBar } from "../../components/shell/TopCommandBar";
import { GridAssistantDock } from "../../components/shell/GridAssistantDock";
import { WorkNotesDrawer } from "../../components/shell/WorkNotesDrawer";
import { BillingCommercialHeader } from "../../features/billing-commercial/BillingCommercialHeader";
import {
  BillingCommercialNavigation,
  BillingTask,
} from "../../features/billing-commercial/BillingCommercialNavigation";
import { ZoneC_ActionQueue } from "../../features/billing-commercial/ZoneC_ActionQueue";
import { ZoneC_IntakeForm } from "../../features/billing-commercial/ZoneC_IntakeForm";
import { CommercialRevenueAnalyticsView } from "../../features/billing-commercial/CommercialRevenueAnalyticsView";
import {
  fetchTransactions,
  fetchTheftAlerts,
  VendingTransactionDTO,
  TheftAlertDTO,
} from "../../features/billing-commercial/api";
import { parseTokenClaims } from "../../lib/auth";
import { AuthClaims } from "../../types/auth";
import { apiClient } from "../../lib/api-client";

export default function BillingCommercialPage() {
  const [claims, setClaims] = useState<AuthClaims | null>(null);
  const [activeTask, setActiveTask] = useState<BillingTask>("vending-desk");
  const [transactions, setTransactions] = useState<VendingTransactionDTO[]>([]);
  const [alerts, setAlerts] = useState<TheftAlertDTO[]>([]);
  const [isChatOpen, setIsChatOpen] = useState(false);
  const [selectedEntity, setSelectedEntity] = useState<{ type: string; id: string } | null>(null);

  const loadData = async () => {
    const [txs, als] = await Promise.all([fetchTransactions(), fetchTheftAlerts()]);
    setTransactions(txs);
    setAlerts(als);
  };

  useEffect(() => {
    const token = apiClient.getToken();
    if (token) setClaims(parseTokenClaims(token));
    loadData();
  }, []);

  return (
    <>
      <Head>
        <title>NESCO Commercial Operations &bull; Zero-Estimation Vending &amp; Revenue AI</title>
      </Head>

      <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans">
        <TopCommandBar
          claims={claims}
          departmentTitle="Commercial Operations & Revenue Protection"
          onToggleChat={() => setIsChatOpen(!isChatOpen)}
        />

        <main className="flex-1 p-4 sm:p-6 space-y-6 max-w-7xl w-full mx-auto">
          {/* Executive Overview & Regulatory Mission */}
          <BillingCommercialHeader
            txCount={transactions.length}
            alertCount={alerts.length}
            onRefresh={loadData}
          />

          {/* Operational Task Navigation Pane */}
          <BillingCommercialNavigation
            activeTask={activeTask}
            onSelectTask={(task) => setActiveTask(task)}
            txCount={transactions.length}
            alertCount={alerts.length}
          />

          {/* Dedicated Task Workspace Views */}
          <div className="transition-all duration-200">
            {activeTask === "vending-desk" && (
              <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">
                <div className="lg:col-span-5">
                  <ZoneC_IntakeForm onSuccess={loadData} />
                </div>
                <div className="lg:col-span-7">
                  <ZoneC_ActionQueue
                    transactions={transactions}
                    alerts={alerts}
                    onSelectTx={(tx) => setSelectedEntity({ type: "vending_tx", id: tx.id })}
                    onSelectAlert={(a) => setSelectedEntity({ type: "theft_alert", id: a.id })}
                  />
                </div>
              </div>
            )}

            {activeTask === "theft-radar" && (
              <ZoneC_ActionQueue
                transactions={transactions}
                alerts={alerts}
                onSelectTx={(tx) => setSelectedEntity({ type: "vending_tx", id: tx.id })}
                onSelectAlert={(a) => setSelectedEntity({ type: "theft_alert", id: a.id })}
              />
            )}

            {activeTask === "vending-ledger" && (
              <ZoneC_ActionQueue
                transactions={transactions}
                alerts={alerts}
                onSelectTx={(tx) => setSelectedEntity({ type: "vending_tx", id: tx.id })}
                onSelectAlert={(a) => setSelectedEntity({ type: "theft_alert", id: a.id })}
              />
            )}

            {activeTask === "revenue-analytics" && (
              <CommercialRevenueAnalyticsView />
            )}
          </div>
        </main>

        <GridAssistantDock
          isOpen={isChatOpen}
          onClose={() => setIsChatOpen(false)}
          staffRole={claims?.primary_role || "UNIT_HEAD_BILLING"}
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
