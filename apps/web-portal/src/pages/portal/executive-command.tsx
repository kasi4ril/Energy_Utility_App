/**
 * Next.js Page: Executive Command Center (/portal/executive-command).
 * Managing Director & CEO Tier-1 Command Landing Page.
 * Structured with:
 * - Zone A: Thin Top Command Bar (Executive identity, WAT clock, quick links)
 * - Persistent Left-Hand Navigation Sidebar (5 operating workspaces)
 * - 5 Operating Workspaces:
 *   1. Approvals & Requisitions Desk (Dept Heads Requests, Review, Minuting, Approval)
 *   2. Executive Analytics & Oversight Board (Requests, Spend, SLA turnaround)
 *   3. Inter-Departmental Information & Minuting Hub (Cross-Unit Memos)
 *   4. MD Professional Responsibilities & Governance Suite (NERC, Board, Concessions, HSE)
 *   5. Decisions Register & Cryptographic Audit Ledger (SHA-256 Chained Stamps)
 * Strictly <= 350 lines.
 */

import React, { useState, useEffect } from "react";
import Head from "next/head";
import { useRouter } from "next/router";
import { ExecutiveHeaderBar } from "../../features/executive-command/ExecutiveHeaderBar";
import { ExecutiveNav_LeftSidebarPane } from "../../features/executive-command/ExecutiveNav_LeftSidebarPane";
import { MDApprovalsDeskView } from "../../features/executive-command/MDApprovalsDeskView";
import { MDExecutiveAnalyticsView } from "../../features/executive-command/MDExecutiveAnalyticsView";
import { MDInterDepartmentHubView } from "../../features/executive-command/MDInterDepartmentHubView";
import { MDGovernanceSuiteView } from "../../features/executive-command/MDGovernanceSuiteView";
import { MDDecisionsLedgerView } from "../../features/executive-command/MDDecisionsLedgerView";
import { GridAssistantDock } from "../../components/shell/GridAssistantDock";
import { WorkNotesDrawer } from "../../components/shell/WorkNotesDrawer";
import {
  MOCK_EXECUTIVE_REQUISITIONS,
  MOCK_INTER_DEPT_MEMOS,
  MOCK_BOARD_RESOLUTIONS,
  MOCK_REGULATORY_FILINGS,
  MOCK_STRATEGIC_CONCESSIONS,
  MOCK_HEAD_SCORECARDS,
} from "../../features/executive-command/mockData";
import {
  ExecutiveWorkspace,
  ExecutiveRequisition,
  ExecutiveInterDeptMemo,
  DepartmentSource,
} from "../../features/executive-command/types";
import { parseTokenClaims } from "../../lib/auth";
import { AuthClaims } from "../../types/auth";
import { apiClient } from "../../lib/api-client";

interface ExecutiveCommandPageProps {
  defaultWorkspace?: ExecutiveWorkspace;
}

export default function ExecutiveCommandPage({
  defaultWorkspace = "approvals",
}: ExecutiveCommandPageProps) {
  const router = useRouter();
  const [claims, setClaims] = useState<AuthClaims | null>(null);
  const [activeWorkspace, setActiveWorkspace] = useState<ExecutiveWorkspace>(defaultWorkspace);
  const [requisitions, setRequisitions] = useState<ExecutiveRequisition[]>(MOCK_EXECUTIVE_REQUISITIONS);
  const [memos, setMemos] = useState<ExecutiveInterDeptMemo[]>(MOCK_INTER_DEPT_MEMOS);
  const [isChatOpen, setIsChatOpen] = useState(false);
  const [toastMessage, setToastMessage] = useState<string | null>(null);

  useEffect(() => {
    const token = apiClient.getToken();
    if (token) setClaims(parseTokenClaims(token));

    if (router.query.view && typeof router.query.view === "string") {
      setActiveWorkspace(router.query.view as ExecutiveWorkspace);
    }
  }, [router.query.view]);

  const showToast = (msg: string) => {
    setToastMessage(msg);
    setTimeout(() => setToastMessage(null), 4000);
  };

  // Requisition Approval Handler
  const handleApprove = (id: string, notes: string) => {
    const officerName = claims?.full_name || "Engr. Dr. E. K. Dung (MD/CEO)";
    setRequisitions((prev) =>
      prev.map((r) =>
        r.id === id
          ? {
              ...r,
              status: "APPROVED",
              md_review_notes: notes,
              decision_date: new Date().toISOString(),
              decision_by: officerName,
              hash_chain_seal: `0x${Math.random().toString(16).substring(2, 10)}${Math.random().toString(16).substring(2, 10)}`,
            }
          : r
      )
    );
    showToast(`Requisition #${id} officially APPROVED and capital expenditure authorized!`);
  };

  // Requisition Review & Query Handler
  const handleReview = (id: string, notes: string) => {
    setRequisitions((prev) =>
      prev.map((r) =>
        r.id === id
          ? {
              ...r,
              status: "UNDER_REVIEW",
              md_review_notes: notes,
              md_review_date: new Date().toISOString(),
            }
          : r
      )
    );
    showToast(`MD Review Directive recorded. Department head queried under audit trace.`);
  };

  // Requisition Minuting to other Departments
  const handleMinute = (id: string, departments: DepartmentSource[], notes: string) => {
    setRequisitions((prev) =>
      prev.map((r) =>
        r.id === id
          ? {
              ...r,
              status: "ACKNOWLEDGED_MINUTED",
              minuted_to_departments: departments,
              minuted_notes: notes,
              decision_date: new Date().toISOString(),
            }
          : r
      )
    );
    showToast(`Requisition minuted to ${departments.length} concerned departments with executive remarks.`);
  };

  // Requisition Decline Handler
  const handleReject = (id: string, notes: string) => {
    setRequisitions((prev) =>
      prev.map((r) =>
        r.id === id
          ? {
              ...r,
              status: "REJECTED",
              md_review_notes: notes,
              decision_date: new Date().toISOString(),
            }
          : r
      )
    );
    showToast(`Requisition #${id} DECLINED and returned to Department Head for revision.`);
  };

  // Memo Minuting Handler
  const handleMinuteMemo = (memoId: string, remarks: string) => {
    setMemos((prev) =>
      prev.map((m) =>
        m.id === memoId
          ? {
              ...m,
              status: "NOTED_AND_MINUTED",
              md_minute_remarks: remarks,
            }
          : m
      )
    );
    showToast(`Executive Directive minuted to concerned Departmental Heads.`);
  };

  const pendingApprovalsCount = requisitions.filter((r) => r.status === "PENDING_APPROVAL").length;
  const underReviewCount = requisitions.filter((r) => r.status === "UNDER_REVIEW").length;

  return (
    <div className="h-screen overflow-hidden flex flex-col bg-slate-950 font-sans antialiased text-slate-100">
      <Head>
        <title>Managing Director Command Center | NESCO Enterprise Sovereign Grid</title>
      </Head>

      {/* Top Header Bar */}
      <ExecutiveHeaderBar
        claims={claims}
        pendingApprovalsCount={pendingApprovalsCount}
        onToggleChat={() => setIsChatOpen(!isChatOpen)}
      />

      {/* Toast Notification Alert */}
      {toastMessage && (
        <div className="fixed top-16 right-6 z-50 bg-emerald-900 border border-emerald-500 text-white px-4 py-2.5 rounded-xl shadow-2xl flex items-center space-x-2 text-xs font-semibold animate-in fade-in slide-in-from-top-4">
          <span>⚡</span>
          <span>{toastMessage}</span>
        </div>
      )}

      {/* Main 2-Column Body Shell */}
      <div className="flex-1 flex flex-col md:flex-row overflow-hidden relative">
        {/* Left-Hand Navigation Sidebar Pane */}
        <ExecutiveNav_LeftSidebarPane
          activeWorkspace={activeWorkspace}
          onSelectWorkspace={(ws) => {
            setActiveWorkspace(ws);
            router.push(`/portal/executive-command?view=${ws}`, undefined, { shallow: true });
          }}
          pendingApprovalsCount={pendingApprovalsCount}
          underReviewCount={underReviewCount}
        />

        {/* Scrollable Main Operational Workspace */}
        <main className="flex-1 overflow-y-auto p-4 sm:p-6 lg:p-8 bg-gradient-to-b from-slate-950 via-slate-900 to-slate-950">
          <div className="max-w-7xl mx-auto">
            {activeWorkspace === "approvals" && (
              <MDApprovalsDeskView
                requisitions={requisitions}
                onApprove={handleApprove}
                onReview={handleReview}
                onMinute={handleMinute}
                onReject={handleReject}
              />
            )}

            {activeWorkspace === "analytics" && (
              <MDExecutiveAnalyticsView requisitions={requisitions} />
            )}

            {activeWorkspace === "inter-department" && (
              <MDInterDepartmentHubView
                memos={memos}
                onMinuteMemo={handleMinuteMemo}
              />
            )}

            {activeWorkspace === "governance" && (
              <MDGovernanceSuiteView
                boardResolutions={MOCK_BOARD_RESOLUTIONS}
                regulatoryFilings={MOCK_REGULATORY_FILINGS}
                strategicConcessions={MOCK_STRATEGIC_CONCESSIONS}
                headScorecards={MOCK_HEAD_SCORECARDS}
              />
            )}

            {activeWorkspace === "ledger" && (
              <MDDecisionsLedgerView
                requisitions={requisitions}
                memos={memos}
              />
            )}
          </div>
        </main>
      </div>

      {/* Slide-out AI Co-Pilot Dock */}
      <GridAssistantDock
        isOpen={isChatOpen}
        onClose={() => setIsChatOpen(false)}
        staffRole={claims?.primary_role || "MANAGING_DIRECTOR"}
      />
    </div>
  );
}
