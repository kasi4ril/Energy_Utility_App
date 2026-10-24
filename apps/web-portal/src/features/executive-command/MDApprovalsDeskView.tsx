/**
 * Workspace 1: Managing Director Approvals & Requisitions Desk.
 * Allows MD to access requisitions from Generation Manager, Head Engineering Services,
 * Commercial Ops, Internal Audit, Stores, and HSE.
 * Features:
 * - Review and Technical Directives space
 * - Inter-department minuting acknowledgment
 * - Authoritative digital sign-off and approval
 * Strictly <= 350 lines.
 */

import React, { useState } from "react";
import { ExecutiveRequisition, DepartmentSource, RequisitionStatus } from "./types";
import { MDRequisitionActionCard } from "./MDRequisitionActionCard";
import {
  CheckCircle2,
  Clock,
  Filter,
  Search,
  Zap,
  Building,
  AlertCircle,
  FileCheck,
} from "lucide-react";

interface ApprovalsDeskProps {
  requisitions: ExecutiveRequisition[];
  onApprove: (id: string, notes: string) => void;
  onReview: (id: string, notes: string) => void;
  onMinute: (id: string, departments: DepartmentSource[], notes: string) => void;
  onReject: (id: string, notes: string) => void;
}

export const MDApprovalsDeskView: React.FC<ApprovalsDeskProps> = ({
  requisitions,
  onApprove,
  onReview,
  onMinute,
  onReject,
}) => {
  const [selectedDept, setSelectedDept] = useState<string>("ALL");
  const [selectedStatus, setSelectedStatus] = useState<string>("ALL");
  const [searchQuery, setSearchQuery] = useState<string>("");
  const [selectedId, setSelectedId] = useState<string>(requisitions[0]?.id || "");

  // KPI calculations
  const totalCount = requisitions.length;
  const pendingCount = requisitions.filter((r) => r.status === "PENDING_APPROVAL").length;
  const underReviewCount = requisitions.filter((r) => r.status === "UNDER_REVIEW").length;
  const approvedCount = requisitions.filter((r) => r.status === "APPROVED").length;
  const totalVolumeNGN = requisitions.reduce((sum, r) => sum + r.amount, 0);

  // Filter requisitions
  const filtered = requisitions.filter((req) => {
    if (selectedDept !== "ALL" && req.source_department !== selectedDept) return false;
    if (selectedStatus !== "ALL" && req.status !== selectedStatus) return false;
    if (searchQuery.trim()) {
      const q = searchQuery.toLowerCase();
      const match =
        req.request_no.toLowerCase().includes(q) ||
        req.title.toLowerCase().includes(q) ||
        req.submitting_head_name.toLowerCase().includes(q);
      if (!match) return false;
    }
    return true;
  });

  const activeRequisition = requisitions.find((r) => r.id === selectedId) || filtered[0];

  return (
    <div className="space-y-6 text-xs text-slate-100">
      {/* 1. Executive Top Banner & 4 KPIs */}
      <div className="bg-gradient-to-r from-slate-900 via-amber-950/40 to-slate-900 border border-amber-500/30 rounded-2xl p-5 shadow-xl">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <div className="flex items-center space-x-2 mb-1.5">
              <span className="px-2.5 py-0.5 rounded-full text-[10px] font-bold tracking-wider uppercase bg-amber-500/20 border border-amber-500/40 text-amber-300">
                Departmental Heads Authorizations &bull; MD Command
              </span>
              <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-emerald-500/20 text-emerald-300 border border-emerald-500/30">
                Maker-Checker Verified
              </span>
            </div>
            <h1 className="text-xl sm:text-2xl font-black text-white tracking-tight">
              Executive Approvals &amp; Departmental Requisitions Desk
            </h1>
            <p className="text-xs text-slate-300 mt-1 max-w-2xl leading-relaxed">
              Centralized decision terminal for review and approval of requisitions submitted by the Generation Manager,
              Head Engineering Services, Commercial Operations, Internal Audit, and Stores.
            </p>
          </div>

          <div className="flex items-center space-x-2 shrink-0">
            <div className="bg-slate-950/80 border border-amber-500/40 px-3.5 py-2 rounded-xl text-right">
              <div className="text-[10px] text-amber-400 font-bold uppercase tracking-wider">Gross Requested</div>
              <div className="text-lg font-black font-mono text-emerald-400">
                ₦{(totalVolumeNGN / 1000000).toFixed(1)}M
              </div>
            </div>
          </div>
        </div>

        {/* 4 KPIs */}
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 mt-4">
          <div className="bg-slate-950/80 border border-slate-800 rounded-xl p-3">
            <div className="text-slate-400 font-semibold text-[11px] flex items-center justify-between">
              <span>Total Requisitions</span>
              <Zap className="w-3.5 h-3.5 text-amber-400" />
            </div>
            <div className="text-xl font-black font-mono text-white mt-1">{totalCount}</div>
          </div>

          <div className="bg-slate-950/80 border border-slate-800 rounded-xl p-3">
            <div className="text-slate-400 font-semibold text-[11px] flex items-center justify-between">
              <span>Pending MD Approval</span>
              <Clock className="w-3.5 h-3.5 text-amber-400" />
            </div>
            <div className="text-xl font-black font-mono text-amber-300 mt-1">{pendingCount} Due</div>
          </div>

          <div className="bg-slate-950/80 border border-slate-800 rounded-xl p-3">
            <div className="text-slate-400 font-semibold text-[11px] flex items-center justify-between">
              <span>Under Review / Query</span>
              <AlertCircle className="w-3.5 h-3.5 text-cyan-400" />
            </div>
            <div className="text-xl font-black font-mono text-cyan-300 mt-1">{underReviewCount}</div>
          </div>

          <div className="bg-slate-950/80 border border-slate-800 rounded-xl p-3">
            <div className="text-slate-400 font-semibold text-[11px] flex items-center justify-between">
              <span>Approved &amp; Authorized</span>
              <FileCheck className="w-3.5 h-3.5 text-emerald-400" />
            </div>
            <div className="text-xl font-black font-mono text-emerald-300 mt-1">{approvedCount}</div>
          </div>
        </div>
      </div>

      {/* 2. Filter & Search Bar */}
      <div className="bg-slate-900 border border-slate-800 rounded-xl p-3 flex flex-wrap items-center justify-between gap-3">
        <div className="flex flex-wrap items-center gap-2">
          {/* Department Dropdown with Explicit High-Contrast Styling for Light & Dark Styles */}
          <div className="flex items-center space-x-1.5 bg-slate-950 px-2.5 py-1 rounded-lg border border-slate-700">
            <Filter className="w-3.5 h-3.5 text-amber-400" />
            <select
              value={selectedDept}
              onChange={(e) => setSelectedDept(e.target.value)}
              className="bg-slate-900 text-slate-100 text-xs font-semibold focus:outline-none focus:ring-1 focus:ring-amber-500 rounded px-1.5 py-0.5 border border-slate-700 cursor-pointer"
              style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
            >
              <option value="ALL" className="bg-slate-900 text-slate-100 py-1" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>
                All Departments
              </option>
              <option value="GENERATION" className="bg-slate-900 text-slate-100 py-1" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>
                Generation (Hydro Fleet)
              </option>
              <option value="ENGINEERING_SERVICES" className="bg-slate-900 text-slate-100 py-1" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>
                Engineering &amp; Distribution
              </option>
              <option value="COMMERCIAL_OPS" className="bg-slate-900 text-slate-100 py-1" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>
                Commercial Operations
              </option>
              <option value="INTERNAL_AUDIT" className="bg-slate-900 text-slate-100 py-1" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>
                Internal Audit &amp; Compliance
              </option>
              <option value="STORES_PROCUREMENT" className="bg-slate-900 text-slate-100 py-1" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>
                Stores &amp; Procurement
              </option>
              <option value="HSE_SAFETY" className="bg-slate-900 text-slate-100 py-1" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>
                HSE &amp; Safety Division
              </option>
            </select>
          </div>

          {/* Status Dropdown with Explicit High-Contrast Styling for Light & Dark Styles */}
          <div className="flex items-center space-x-1.5 bg-slate-950 px-2.5 py-1 rounded-lg border border-slate-700">
            <span className="text-slate-300 text-[11px] font-semibold">Status:</span>
            <select
              value={selectedStatus}
              onChange={(e) => setSelectedStatus(e.target.value)}
              className="bg-slate-900 text-slate-100 text-xs font-semibold focus:outline-none focus:ring-1 focus:ring-amber-500 rounded px-1.5 py-0.5 border border-slate-700 cursor-pointer"
              style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}
            >
              <option value="ALL" className="bg-slate-900 text-slate-100 py-1" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>
                All Statuses
              </option>
              <option value="PENDING_APPROVAL" className="bg-slate-900 text-slate-100 py-1" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>
                Pending Approval
              </option>
              <option value="UNDER_REVIEW" className="bg-slate-900 text-slate-100 py-1" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>
                Under Review / Queried
              </option>
              <option value="APPROVED" className="bg-slate-900 text-slate-100 py-1" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>
                Approved &amp; Authorized
              </option>
              <option value="ACKNOWLEDGED_MINUTED" className="bg-slate-900 text-slate-100 py-1" style={{ backgroundColor: "#0f172a", color: "#f8fafc" }}>
                Acknowledged &amp; Minuted
              </option>
            </select>
          </div>
        </div>

        <div className="relative w-full sm:w-64">
          <Search className="w-3.5 h-3.5 absolute left-2.5 top-2.5 text-slate-400" />
          <input
            type="text"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            placeholder="Search by Request #, title, head..."
            className="w-full bg-slate-950 border border-slate-800 rounded-lg pl-8 pr-3 py-1.5 text-xs text-white placeholder-slate-500 focus:outline-none focus:border-amber-500"
          />
        </div>
      </div>

      {/* 3. Master-Detail Workspace Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-5 items-start">
        {/* Left Column: Requisition Selector Cards */}
        <div className="lg:col-span-5 space-y-2.5">
          <div className="text-[11px] font-bold text-slate-400 uppercase tracking-wider flex items-center justify-between px-1">
            <span>Requisitions Queue ({filtered.length})</span>
            <span className="font-mono text-amber-400">{pendingCount} Awaiting Decision</span>
          </div>

          {filtered.length === 0 ? (
            <div className="bg-slate-900 border border-slate-800 rounded-xl p-8 text-center text-slate-500">
              No requisitions match the selected filters.
            </div>
          ) : (
            filtered.map((req) => {
              const isSelected = req.id === activeRequisition?.id;
              return (
                <div
                  key={req.id}
                  onClick={() => setSelectedId(req.id)}
                  className={`p-3.5 rounded-xl border transition-all cursor-pointer space-y-2 ${
                    isSelected
                      ? "bg-slate-800/90 border-amber-500 shadow-md ring-1 ring-amber-500/40"
                      : "bg-slate-900/80 border-slate-800 hover:border-slate-700 hover:bg-slate-850"
                  }`}
                >
                  <div className="flex items-center justify-between">
                    <span className="font-mono font-bold text-[11px] text-amber-400">
                      {req.request_no}
                    </span>
                    <span
                      className={`text-[9px] font-bold px-1.5 py-0.2 rounded uppercase font-mono ${
                        req.status === "APPROVED"
                          ? "bg-emerald-500/20 text-emerald-300"
                          : req.status === "UNDER_REVIEW"
                          ? "bg-cyan-500/20 text-cyan-300"
                          : req.status === "ACKNOWLEDGED_MINUTED"
                          ? "bg-blue-500/20 text-blue-300"
                          : "bg-amber-500/20 text-amber-300"
                      }`}
                    >
                      {req.status.replace(/_/g, " ")}
                    </span>
                  </div>

                  <h3 className="font-bold text-white text-xs line-clamp-2 leading-snug">
                    {req.title}
                  </h3>

                  <div className="flex items-center justify-between text-[11px] pt-1 border-t border-slate-800/80 text-slate-400">
                    <span>{req.submitting_head_name}</span>
                    <strong className="text-emerald-400 font-mono font-bold">
                      ₦{req.amount.toLocaleString()}
                    </strong>
                  </div>
                </div>
              );
            })
          )}
        </div>

        {/* Right Column: Active Requisition Review, Minuting & Approval Card */}
        <div className="lg:col-span-7">
          {activeRequisition ? (
            <MDRequisitionActionCard
              key={activeRequisition.id}
              requisition={activeRequisition}
              onApprove={onApprove}
              onReview={onReview}
              onMinute={onMinute}
              onReject={onReject}
            />
          ) : (
            <div className="bg-slate-900 border border-slate-800 rounded-xl p-12 text-center text-slate-500">
              Select a requisition from the queue to review and authorize.
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

