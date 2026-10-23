/**
 * Next.js Page: Customer Care & Service Delivery Department (/portal/customer-care).
 * Compliant with Section 12: Thin Top Header Bar (ZoneA) + Persistent Left-Hand Navigation Sidebar
 * (ZoneNav) + 5-Point Operating Views (Analytics, KYC, Complaint Desk, KCT & Closure Gate).
 * Strictly <= 350 lines.
 */

import React, { useState, useEffect } from "react";
import Head from "next/head";
import { useRouter } from "next/router";
import { ZoneA_ThinTopHeaderBar } from "../../features/customer-care/ZoneA_ThinTopHeaderBar";
import { ZoneNav_LeftSidebarPane } from "../../features/customer-care/ZoneNav_LeftSidebarPane";
import { View1_AnalyticsStory } from "../../features/customer-care/View1_AnalyticsStory";
import { View2_NewCustomerKYC } from "../../features/customer-care/View2_NewCustomerKYC";
import { View3_SmartComplaintDesk } from "../../features/customer-care/View3_SmartComplaintDesk";
import { View4_TokenAndKCTDesk } from "../../features/customer-care/View4_TokenAndKCTDesk";
import { View5_ClosureGateDesk } from "../../features/customer-care/View5_ClosureGateDesk";
import { GridAssistantDock } from "../../components/shell/GridAssistantDock";
import { WorkNotesDrawer } from "../../components/shell/WorkNotesDrawer";
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

interface CustomerCarePageProps {
  defaultView?: CustomerCareTask;
}

export default function CustomerCarePage({ defaultView = "analytics" }: CustomerCarePageProps) {
  const router = useRouter();
  const [claims, setClaims] = useState<AuthClaims | null>(null);
  const [activeTask, setActiveTask] = useState<CustomerCareTask>(defaultView);
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

    // Check query param for initial view
    if (router.query.view && typeof router.query.view === "string") {
      setActiveTask(router.query.view as CustomerCareTask);
    }
  }, [router.query.view]);

  const handleConfirmClose = async (ticketId: string, notes: string, confirmed: boolean) => {
    const staffId = claims?.staff_id || "STAFF-CC-02";
    const staffName = claims?.full_name || "Amina Yusuf (CC Officer)";
    try {
      await confirmTicketClosure(ticketId, notes, staffId, staffName, confirmed);
      setTickets((prev) =>
        prev.map((t) =>
          t.id === ticketId
            ? {
                ...t,
                status: confirmed ? "CLOSED" : "REOPENED_ESCALATED",
                customer_call_verified: true,
                customer_feedback_notes: notes,
                closed_by_cc_staff_id: staffId,
                closed_by_cc_staff_name: staffName,
              }
            : t
        )
      );
      alert(
        confirmed
          ? "Ticket successfully CLOSED under NERC 2-Step Telephone Confirmation Gate!"
          : "Ticket REOPENED and escalated to Technical Department Head!"
      );
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

  const emergencyCount = tickets.filter(
    (t) => t.is_safety_emergency || t.category === "LINE_SNAP"
  ).length;

  const staffId = claims?.staff_id || "STAFF-CC-02";
  const staffName = claims?.full_name || "Amina Yusuf (CC Officer)";

  return (
    <>
      <Head>
        <title>NESCO Customer Care &bull; Service Delivery Department</title>
      </Head>

      <div className="h-screen bg-slate-950 text-slate-100 flex flex-col font-sans overflow-hidden">
        {/* 1. Thin Top Header Bar (Zone A) - Stand Still / Static */}
        <ZoneA_ThinTopHeaderBar
          claims={claims}
          breachCount={emergencyCount}
          onToggleChat={() => setIsChatOpen(!isChatOpen)}
        />

        {/* 2. Main 2-Column Operational Layout */}
        <div className="flex-1 flex flex-col md:flex-row w-full overflow-hidden">
          {/* Left-Hand Navigation Sidebar (ZoneNav) - Stand Still / Static */}
          <ZoneNav_LeftSidebarPane
            activeTask={activeTask}
            onSelectTask={(task) => {
              setActiveTask(task);
              router.push({ pathname: "/portal/customer-care", query: { view: task } }, undefined, {
                shallow: true,
              });
            }}
            openTicketsCount={openTicketsCount}
            awaitingClosureCount={awaitingConfirmationCount}
            emergencyCount={emergencyCount}
          />

          {/* Active Operating View Workspace - Independent Scrollable Area */}
          <main className="flex-1 p-4 sm:p-6 overflow-y-auto min-w-0">
            {activeTask === "analytics" && analytics && (
              <View1_AnalyticsStory analytics={analytics} onRefresh={loadData} />
            )}

            {(activeTask === "register" || activeTask === "kyc-onboarding") && (
              <View2_NewCustomerKYC
                staffId={staffId}
                staffName={staffName}
                onCustomerOnboarded={loadData}
              />
            )}

            {(activeTask === "complain" || activeTask === "triage-queue" || activeTask === "log-complaint") && (
              <View3_SmartComplaintDesk
                customers={customers}
                staffId={staffId}
                staffName={staffName}
                onTicketLogged={loadData}
              />
            )}

            {activeTask === "metering-kct" && (
              <View4_TokenAndKCTDesk staffId={staffId} staffName={staffName} />
            )}

            {(activeTask === "closure-desk" || activeTask === "crm-directory") && (
              <View5_ClosureGateDesk
                tickets={tickets}
                staffId={staffId}
                staffName={staffName}
                onConfirmClose={handleConfirmClose}
                onRefresh={loadData}
              />
            )}
          </main>
        </div>

        {/* Co-Pilot AI & WorkNotes Drawer */}
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
