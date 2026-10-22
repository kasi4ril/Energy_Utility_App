/**
 * Next.js Page: Customer Care & Consumer Relations Department (/portal/customer-care).
 * Full Enterprise Portal with Mission Header, Task Navigation Pane, and Modular Workspaces.
 */

import React, { useState, useEffect } from "react";
import Head from "next/head";
import { TopCommandBar } from "../../components/shell/TopCommandBar";
import { GridAssistantDock } from "../../components/shell/GridAssistantDock";
import { WorkNotesDrawer } from "../../components/shell/WorkNotesDrawer";
import { CustomerCareHeader } from "../../features/customer-care/CustomerCareHeader";
import { CustomerCareNavigation } from "../../features/customer-care/CustomerCareNavigation";
import { TriageQueueView } from "../../features/customer-care/TriageQueueView";
import { LogComplaintForm } from "../../features/customer-care/LogComplaintForm";
import { KYCOnboardingForm } from "../../features/customer-care/KYCOnboardingForm";
import { CustomerCRMView } from "../../features/customer-care/CustomerCRMView";
import { CustomerAnalyticsView } from "../../features/customer-care/CustomerAnalyticsView";
import {
  fetchCustomerTickets,
  fetchCustomerDirectory,
  confirmTicketClosure,
  fetchCustomerAnalytics,
} from "../../features/customer-care/api";
import {
  CustomerCareTask,
  ServiceTicket,
  CustomerRecord,
  CustomerCareAnalytics,
} from "../../features/customer-care/types";
import { parseTokenClaims } from "../../lib/auth";
import { AuthClaims } from "../../types/auth";
import { apiClient } from "../../lib/api-client";

export default function CustomerCarePage() {
  const [claims, setClaims] = useState<AuthClaims | null>(null);
  const [activeTask, setActiveTask] = useState<CustomerCareTask>("triage-queue");
  const [tickets, setTickets] = useState<ServiceTicket[]>([]);
  const [customers, setCustomers] = useState<CustomerRecord[]>([]);
  const [analytics, setAnalytics] = useState<CustomerCareAnalytics | null>(null);
  const [selectedTicket, setSelectedTicket] = useState<ServiceTicket | null>(null);
  const [isChatOpen, setIsChatOpen] = useState(false);

  const loadData = async () => {
    const [tList, cList, aData] = await Promise.all([
      fetchCustomerTickets(),
      fetchCustomerDirectory(),
      fetchCustomerAnalytics(),
    ]);
    setTickets(tList);
    setCustomers(cList);
    setAnalytics(aData);
  };

  useEffect(() => {
    const token = apiClient.getToken();
    if (token) setClaims(parseTokenClaims(token));
    loadData();
  }, []);

  const handleConfirmClose = async (ticketId: string, notes: string) => {
    const staffId = claims?.staff_id || "STAFF-CC-01";
    try {
      await confirmTicketClosure(ticketId, notes, staffId);
      setTickets((prev) =>
        prev.map((t) =>
          t.id === ticketId
            ? {
                ...t,
                status: "CLOSED",
                customer_call_verified: true,
                customer_feedback_notes: notes,
                verified_by_cc_officer_id: staffId,
              }
            : t
        )
      );
      alert("Ticket closed successfully under NERC 2-step customer call verification protocol!");
    } catch (err: any) {
      alert(err.message || "Failed to confirm ticket closure");
    }
  };

  const openTicketsCount = tickets.filter(
    (t) => t.status === "ROUTED" || t.status === "LOGGED"
  ).length;

  const awaitingConfirmationCount = tickets.filter(
    (t) => t.status === "ATTENDED_TO"
  ).length;

  return (
    <>
      <Head>
        <title>NESCO Customer Care &bull; Consumer Relations &amp; CRM Workspace</title>
      </Head>

      <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans">
        <TopCommandBar
          claims={claims}
          departmentTitle="Consumer Relations & Customer Care"
          onToggleChat={() => setIsChatOpen(!isChatOpen)}
        />

        <main className="flex-1 p-4 sm:p-6 space-y-6 max-w-7xl w-full mx-auto">
          {/* Department Executive Overview & Regulatory Mission */}
          <CustomerCareHeader
            activeTicketsCount={openTicketsCount}
            attendedAwaitingCloseCount={awaitingConfirmationCount}
            onRefresh={loadData}
          />

          {/* Operational Task Navigation Pane */}
          <CustomerCareNavigation
            activeTask={activeTask}
            onSelectTask={(task) => setActiveTask(task)}
            openTicketsCount={openTicketsCount}
            awaitingConfirmationCount={awaitingConfirmationCount}
            totalCustomersCount={customers.length}
          />

          {/* Dedicated Task Workspace Views */}
          <div className="transition-all duration-200">
            {activeTask === "triage-queue" && (
              <TriageQueueView
                tickets={tickets}
                onSelectTicket={(t) => setSelectedTicket(t)}
                onConfirmClose={handleConfirmClose}
                onNavigateToLogComplaint={() => setActiveTask("log-complaint")}
              />
            )}

            {activeTask === "log-complaint" && (
              <LogComplaintForm
                customers={customers}
                onTicketCreated={loadData}
                onNavigateToQueue={() => setActiveTask("triage-queue")}
              />
            )}

            {activeTask === "kyc-onboarding" && (
              <KYCOnboardingForm
                staffId={claims?.staff_id || "STAFF-CC-01"}
                onCustomerOnboarded={loadData}
                onNavigateToCRM={() => setActiveTask("crm-directory")}
              />
            )}

            {activeTask === "crm-directory" && (
              <CustomerCRMView
                customers={customers}
                onLogComplaintForCustomer={() => setActiveTask("log-complaint")}
                onNavigateToKYC={() => setActiveTask("kyc-onboarding")}
              />
            )}

            {activeTask === "analytics" && analytics && (
              <CustomerAnalyticsView
                analytics={analytics}
                onRefresh={loadData}
              />
            )}
          </div>
        </main>

        <GridAssistantDock
          isOpen={isChatOpen}
          onClose={() => setIsChatOpen(false)}
          staffRole={claims?.primary_role || "CUSTOMER_CARE_OFFICER"}
        />

        {selectedTicket && claims && (
          <WorkNotesDrawer
            isOpen={true}
            onClose={() => setSelectedTicket(null)}
            entityType="ticket"
            entityId={selectedTicket.id}
            authorStaffId={claims.staff_id}
            authorRole={claims.primary_role}
          />
        )}
      </div>
    </>
  );
}
