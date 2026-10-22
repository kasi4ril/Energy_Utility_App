/**
 * Next.js Page: Customer Care & KYC Portal (/portal/customer-care).
 * 4-Zone Enterprise Shell for Customer Care Officers & Triage.
 */

import React, { useState, useEffect } from "react";
import Head from "next/head";
import { TopCommandBar } from "../../components/shell/TopCommandBar";
import { GridAssistantDock } from "../../components/shell/GridAssistantDock";
import { WorkNotesDrawer } from "../../components/shell/WorkNotesDrawer";
import { ZoneB_DashboardRibbon } from "../../features/customer-care/ZoneB_DashboardRibbon";
import { ZoneC_ActionQueue } from "../../features/customer-care/ZoneC_ActionQueue";
import { ZoneC_IntakeForm } from "../../features/customer-care/ZoneC_IntakeForm";
import { fetchCustomerTickets, confirmTicketClosure, ServiceTicketDTO } from "../../features/customer-care/api";
import { parseTokenClaims } from "../../lib/auth";
import { AuthClaims } from "../../types/auth";
import { apiClient } from "../../lib/api-client";

export default function CustomerCarePage() {
  const [claims, setClaims] = useState<AuthClaims | null>(null);
  const [tickets, setTickets] = useState<ServiceTicketDTO[]>([]);
  const [isChatOpen, setIsChatOpen] = useState(false);
  const [selectedTicket, setSelectedTicket] = useState<ServiceTicketDTO | null>(null);

  useEffect(() => {
    const token = apiClient.getToken();
    if (token) setClaims(parseTokenClaims(token));

    const loadTickets = async () => {
      const data = await fetchCustomerTickets();
      setTickets(data);
    };
    loadTickets();
  }, []);

  const handleConfirmClose = async (ticketId: string) => {
    if (!claims) return;
    try {
      await confirmTicketClosure(ticketId, "Customer verbally confirmed fault resolution.", claims.staff_id);
      setTickets((prev) =>
        prev.map((t) => (t.id === ticketId ? { ...t, status: "CLOSED" } : t))
      );
      alert("Ticket closed following 2-step customer confirmation gate!");
    } catch (err: any) {
      alert(err.message || "Failed to confirm ticket closure");
    }
  };

  return (
    <>
      <Head>
        <title>NESCO Customer Care &bull; KYC Onboarding & SLA Triage</title>
      </Head>

      <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col">
        <TopCommandBar
          claims={claims}
          departmentTitle="Customer Care & KYC Triage"
          onToggleChat={() => setIsChatOpen(!isChatOpen)}
        />

        <main className="flex-1 p-6 space-y-6 max-w-7xl w-full mx-auto">
          <ZoneB_DashboardRibbon />

          <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 min-h-[500px]">
            <div className="lg:col-span-8">
              <ZoneC_ActionQueue
                tickets={tickets}
                onSelectTicket={(t) => setSelectedTicket(t)}
                onConfirmClose={handleConfirmClose}
              />
            </div>
            <div className="lg:col-span-4">
              <ZoneC_IntakeForm
                staffId={claims?.staff_id || "STAFF-CC-01"}
                onSuccess={async () => {
                  const updated = await fetchCustomerTickets();
                  setTickets(updated);
                }}
              />
            </div>
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

