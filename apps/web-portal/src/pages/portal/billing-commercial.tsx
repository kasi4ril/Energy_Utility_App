/**
 * Next.js Page: Billing & Commercial Operations (/portal/billing-commercial).
 * 4-Zone Enterprise Shell for Zero-Estimation Vending and Revenue Protection AI.
 */

import React, { useState, useEffect } from "react";
import Head from "next/head";
import { TopCommandBar } from "../../components/shell/TopCommandBar";
import { GridAssistantDock } from "../../components/shell/GridAssistantDock";
import { WorkNotesDrawer } from "../../components/shell/WorkNotesDrawer";
import { ZoneB_DashboardRibbon } from "../../features/billing-commercial/ZoneB_DashboardRibbon";
import { ZoneC_ActionQueue } from "../../features/billing-commercial/ZoneC_ActionQueue";
import { ZoneC_IntakeForm } from "../../features/billing-commercial/ZoneC_IntakeForm";
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
  const [transactions, setTransactions] = useState<VendingTransactionDTO[]>([]);
  const [alerts, setAlerts] = useState<TheftAlertDTO[]>([]);
  const [isChatOpen, setIsChatOpen] = useState(false);
  const [selectedEntity, setSelectedEntity] = useState<{ type: string; id: string } | null>(null);

  useEffect(() => {
    const token = apiClient.getToken();
    if (token) setClaims(parseTokenClaims(token));

    const loadData = async () => {
      const txs = await fetchTransactions();
      setTransactions(txs);
      const als = await fetchTheftAlerts();
      setAlerts(als);
    };
    loadData();
  }, []);

  return (
    <>
      <Head>
        <title>NESCO Commercial Operations &bull; Zero-Estimation Vending</title>
      </Head>

      <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col">
        <TopCommandBar
          claims={claims}
          departmentTitle="Billing, Commercial Operations & Revenue Protection"
          onToggleChat={() => setIsChatOpen(!isChatOpen)}
        />

        <main className="flex-1 p-6 space-y-6 max-w-7xl w-full mx-auto">
          <ZoneB_DashboardRibbon />

          <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 min-h-[500px]">
            <div className="lg:col-span-8">
              <ZoneC_ActionQueue
                transactions={transactions}
                alerts={alerts}
                onSelectTx={(tx) => setSelectedEntity({ type: "vending_tx", id: tx.id })}
                onSelectAlert={(a) => setSelectedEntity({ type: "theft_alert", id: a.id })}
              />
            </div>
            <div className="lg:col-span-4">
              <ZoneC_IntakeForm
                onSuccess={async () => {
                  const updated = await fetchTransactions();
                  setTransactions(updated);
                }}
              />
            </div>
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

