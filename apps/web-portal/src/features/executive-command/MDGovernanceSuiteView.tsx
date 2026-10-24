/**
 * Workspace 4: Managing Director Professional Responsibilities & Governance Suite.
 * Covers:
 * 1. NERC Regulatory & Statutory License Governance (MYTO, NEMSA, FMWR)
 * 2. Board of Directors Governance & Capital Resolutions (> ₦50M)
 * 3. Strategic Grid Expansion & Public-Private Partnerships (PRIDA, PPAs)
 * 4. High-Risk HSE Emergency Command Authority (Level 3 Grid Disaster)
 * 5. Departmental Heads Performance Scorecards & Executive Ratings
 * Strictly <= 350 lines.
 */

import React, { useState } from "react";
import {
  BoardResolution,
  RegulatoryFiling,
  StrategicConcession,
  DepartmentHeadScorecard,
} from "./types";
import {
  Landmark,
  ShieldCheck,
  FileCheck2,
  Award,
  Users,
  Flame,
  AlertTriangle,
  Building,
  CheckCircle,
} from "lucide-react";

interface GovernanceSuiteProps {
  boardResolutions: BoardResolution[];
  regulatoryFilings: RegulatoryFiling[];
  strategicConcessions: StrategicConcession[];
  headScorecards: DepartmentHeadScorecard[];
}

export const MDGovernanceSuiteView: React.FC<GovernanceSuiteProps> = ({
  boardResolutions,
  regulatoryFilings,
  strategicConcessions,
  headScorecards,
}) => {
  const [activeTab, setActiveTab] = useState<"REGULATORY" | "BOARD" | "CONCESSIONS" | "HSE_COMMAND" | "SCORECARDS">("REGULATORY");

  return (
    <div className="space-y-6 text-xs text-slate-100">
      {/* 1. Executive Banner */}
      <div className="bg-gradient-to-r from-slate-900 via-amber-950/40 to-slate-900 border border-amber-500/30 rounded-2xl p-5 shadow-xl">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <div className="flex items-center space-x-2 mb-1.5">
              <span className="px-2.5 py-0.5 rounded-full text-[10px] font-bold tracking-wider uppercase bg-amber-500/20 border border-amber-500/40 text-amber-300">
                Corporate Governance &bull; Statutory Executive Authority
              </span>
              <span className="px-2 py-0.5 rounded text-[10px] font-mono bg-emerald-500/20 text-emerald-300 border border-emerald-500/30">
                MD / CEO Mandate
              </span>
            </div>
            <h1 className="text-xl sm:text-2xl font-black text-white tracking-tight">
              Managing Director Professional Responsibilities &amp; Governance Suite
            </h1>
            <p className="text-xs text-slate-300 mt-1 max-w-2xl leading-relaxed">
              Statutory oversight of NERC grid licenses, Board of Directors resolutions, bilateral power purchase
              concessions, Level-3 HSE disaster command, and Departmental Heads executive scorecards.
            </p>
          </div>

          <div className="flex items-center space-x-2 bg-slate-950/80 border border-amber-500/40 px-3.5 py-2 rounded-xl text-right shrink-0">
            <div>
              <div className="text-[10px] text-amber-400 font-bold uppercase tracking-wider">Concession License</div>
              <div className="text-sm font-black font-mono text-white">NERC/IP/DISCO/042</div>
            </div>
          </div>
        </div>

        {/* Governance Navigation Tabs */}
        <div className="flex flex-wrap gap-2 mt-5 pt-3 border-t border-slate-800">
          <button
            onClick={() => setActiveTab("REGULATORY")}
            className={`px-3 py-1.5 rounded-lg text-xs font-semibold flex items-center space-x-1.5 transition ${
              activeTab === "REGULATORY"
                ? "bg-amber-600 text-white shadow"
                : "bg-slate-950 text-slate-400 hover:text-white"
            }`}
          >
            <ShieldCheck className="w-3.5 h-3.5" />
            <span>NERC &amp; Statutory Regulations</span>
          </button>

          <button
            onClick={() => setActiveTab("BOARD")}
            className={`px-3 py-1.5 rounded-lg text-xs font-semibold flex items-center space-x-1.5 transition ${
              activeTab === "BOARD"
                ? "bg-amber-600 text-white shadow"
                : "bg-slate-950 text-slate-400 hover:text-white"
            }`}
          >
            <Landmark className="w-3.5 h-3.5" />
            <span>Board Resolutions (&gt; ₦50M)</span>
          </button>

          <button
            onClick={() => setActiveTab("CONCESSIONS")}
            className={`px-3 py-1.5 rounded-lg text-xs font-semibold flex items-center space-x-1.5 transition ${
              activeTab === "CONCESSIONS"
                ? "bg-amber-600 text-white shadow"
                : "bg-slate-950 text-slate-400 hover:text-white"
            }`}
          >
            <Building className="w-3.5 h-3.5" />
            <span>Strategic Concessions &amp; PPAs</span>
          </button>

          <button
            onClick={() => setActiveTab("HSE_COMMAND")}
            className={`px-3 py-1.5 rounded-lg text-xs font-semibold flex items-center space-x-1.5 transition ${
              activeTab === "HSE_COMMAND"
                ? "bg-amber-600 text-white shadow"
                : "bg-slate-950 text-slate-400 hover:text-white"
            }`}
          >
            <Flame className="w-3.5 h-3.5" />
            <span>HSE Emergency Command</span>
          </button>

          <button
            onClick={() => setActiveTab("SCORECARDS")}
            className={`px-3 py-1.5 rounded-lg text-xs font-semibold flex items-center space-x-1.5 transition ${
              activeTab === "SCORECARDS"
                ? "bg-amber-600 text-white shadow"
                : "bg-slate-950 text-slate-400 hover:text-white"
            }`}
          >
            <Award className="w-3.5 h-3.5" />
            <span>Dept Heads Scorecards</span>
          </button>
        </div>
      </div>

      {/* 2. Tab Content Panels */}
      {activeTab === "REGULATORY" && (
        <div className="bg-slate-900 border border-slate-800 rounded-xl p-5 shadow-lg space-y-4">
          <div className="flex items-center justify-between pb-2 border-b border-slate-800">
            <div>
              <h2 className="text-sm font-bold text-white flex items-center space-x-2">
                <ShieldCheck className="w-4 h-4 text-amber-400" />
                <span>NERC, NEMSA &amp; Statutory Regulatory Filings</span>
              </h2>
              <p className="text-[11px] text-slate-400 mt-0.5">
                Statutory tariff compliance, NEMSA safety certifications, and river basin hydrological permits.
              </p>
            </div>
            <span className="text-[10px] font-mono text-emerald-400 bg-emerald-500/10 px-2 py-0.5 rounded border border-emerald-500/20">
              100% Audit Adherence
            </span>
          </div>

          <div className="space-y-3">
            {regulatoryFilings.map((filing) => (
              <div
                key={filing.id}
                className="p-3.5 bg-slate-950/70 border border-slate-800 rounded-lg flex flex-col sm:flex-row sm:items-center justify-between gap-3"
              >
                <div>
                  <div className="flex items-center space-x-2">
                    <span className="font-mono font-bold text-amber-400 text-xs">{filing.filing_code}</span>
                    <span className="text-[9px] font-bold px-1.5 py-0.2 rounded bg-slate-800 text-slate-300 font-mono">
                      {filing.regulatory_body}
                    </span>
                  </div>
                  <h3 className="font-bold text-white text-xs mt-1">{filing.title}</h3>
                  <div className="text-[11px] text-slate-400 mt-0.5">
                    Benchmark: <span className="text-slate-300">{filing.license_benchmark}</span>
                  </div>
                </div>

                <div className="text-left sm:text-right shrink-0">
                  <span className="text-[10px] font-bold px-2 py-0.5 rounded bg-emerald-500/20 text-emerald-300 border border-emerald-500/30">
                    {filing.status.replace(/_/g, " ")}
                  </span>
                  <div className="text-[10px] text-slate-400 font-mono mt-1">Due: {filing.due_date}</div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {activeTab === "BOARD" && (
        <div className="bg-slate-900 border border-slate-800 rounded-xl p-5 shadow-lg space-y-4">
          <div className="flex items-center justify-between pb-2 border-b border-slate-800">
            <div>
              <h2 className="text-sm font-bold text-white flex items-center space-x-2">
                <Landmark className="w-4 h-4 text-amber-400" />
                <span>Board of Directors Strategic Resolutions (&gt; ₦50M Mandates)</span>
              </h2>
              <p className="text-[11px] text-slate-400 mt-0.5">
                Executive capital projects approved by Board of Directors and implemented under MD direct oversight.
              </p>
            </div>
          </div>

          <div className="space-y-3">
            {boardResolutions.map((res) => (
              <div key={res.id} className="p-4 bg-slate-950/70 border border-slate-800 rounded-lg space-y-2">
                <div className="flex items-center justify-between">
                  <span className="font-mono font-bold text-amber-400 text-xs">{res.resolution_no}</span>
                  <span className="text-[10px] font-mono font-bold px-2 py-0.5 rounded bg-cyan-500/20 text-cyan-300 border border-cyan-500/30">
                    {res.status.replace(/_/g, " ")}
                  </span>
                </div>
                <h3 className="font-bold text-white text-sm">{res.title}</h3>
                <p className="text-xs text-slate-300">{res.mandate_summary}</p>
                <div className="flex items-center justify-between text-[11px] pt-2 border-t border-slate-800 text-slate-400 font-mono">
                  <span>Authorized Budget: <strong className="text-emerald-400">₦{(res.allocated_budget_ngn / 1000000).toFixed(1)}M</strong></span>
                  <span>Meeting: {res.meeting_date}</span>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {activeTab === "CONCESSIONS" && (
        <div className="bg-slate-900 border border-slate-800 rounded-xl p-5 shadow-lg space-y-4">
          <div className="flex items-center justify-between pb-2 border-b border-slate-800">
            <div>
              <h2 className="text-sm font-bold text-white flex items-center space-x-2">
                <Building className="w-4 h-4 text-amber-400" />
                <span>Strategic Concessions &amp; Bilateral Power Purchase Agreements</span>
              </h2>
              <p className="text-[11px] text-slate-400 mt-0.5">
                Long-term utility concessions, State Government PPPs, and dedicated industrial direct-feed PPAs.
              </p>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-3.5">
            {strategicConcessions.map((con) => (
              <div key={con.id} className="p-4 bg-slate-950/70 border border-slate-800 rounded-lg space-y-2.5">
                <div className="flex items-center justify-between">
                  <h3 className="font-bold text-white text-xs">{con.partner_name}</h3>
                  <span className="text-[9px] font-bold px-1.5 py-0.2 rounded bg-emerald-500/20 text-emerald-300 border border-emerald-500/30">
                    {con.status.replace(/_/g, " ")}
                  </span>
                </div>
                <p className="text-xs text-slate-300">{con.concession_scope}</p>
                <div className="pt-2 border-t border-slate-800 grid grid-cols-2 gap-2 text-[10px] text-slate-400 font-mono">
                  <div>Committed: <strong className="text-amber-400">{con.capacity_commitment_mw} MW</strong></div>
                  <div>Quarterly Rev: <strong className="text-emerald-400">₦{(con.quarterly_revenue_est_ngn / 1000000).toFixed(1)}M</strong></div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {activeTab === "HSE_COMMAND" && (
        <div className="bg-slate-900 border border-rose-800/80 rounded-xl p-5 shadow-lg space-y-4">
          <div className="flex items-center justify-between pb-2 border-b border-slate-800">
            <div>
              <h2 className="text-sm font-bold text-rose-200 flex items-center space-x-2">
                <Flame className="w-4 h-4 text-rose-400" />
                <span>High-Risk HSE Emergency Command Authority</span>
              </h2>
              <p className="text-[11px] text-slate-400 mt-0.5">
                Sole executive power to declare Level-3 Grid Emergencies and convene statutory Boards of Inquiry.
              </p>
            </div>
            <span className="text-[10px] font-bold px-2 py-0.5 rounded bg-rose-500/20 text-rose-300 border border-rose-500/40">
              Commander Authority Active
            </span>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-3.5">
            <div className="bg-slate-950/80 border border-slate-800 rounded-lg p-3.5 space-y-2">
              <h4 className="font-bold text-white text-xs">Level-3 Emergency Grid Disaster Declaration</h4>
              <p className="text-[11px] text-slate-300">
                Authorizes instantaneous isolation of 33kV backbone circuits, emergency spillway discharges at Kurra Reservoir, and mutual aid mobilization.
              </p>
              <button className="px-3 py-1.5 rounded bg-rose-900/60 hover:bg-rose-800 text-rose-200 text-xs font-semibold border border-rose-700">
                View Emergency Action Protocols
              </button>
            </div>

            <div className="bg-slate-950/80 border border-slate-800 rounded-lg p-3.5 space-y-2">
              <h4 className="font-bold text-white text-xs">Statutory Board of Inquiry Commissioning</h4>
              <p className="text-[11px] text-slate-300">
                Mandatory convening of independent technical inquiry for severe electrical accidents, arc flash incidents, or infrastructure sabotage.
              </p>
              <button className="px-3 py-1.5 rounded bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-semibold border border-slate-700">
                Commission Inquiry Panel
              </button>
            </div>
          </div>
        </div>
      )}

      {activeTab === "SCORECARDS" && (
        <div className="bg-slate-900 border border-slate-800 rounded-xl p-5 shadow-lg space-y-4">
          <div className="flex items-center justify-between pb-2 border-b border-slate-800">
            <div>
              <h2 className="text-sm font-bold text-white flex items-center space-x-2">
                <Award className="w-4 h-4 text-amber-400" />
                <span>Departmental Heads Executive Performance Scorecards</span>
              </h2>
              <p className="text-[11px] text-slate-400 mt-0.5">
                Monthly executive review rating for Generation Manager, Head Engineering Services, and C-Suite officers.
              </p>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-3.5">
            {headScorecards.map((sc) => (
              <div key={sc.department} className="p-3.5 bg-slate-950/70 border border-slate-800 rounded-lg space-y-2">
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="font-bold text-white text-xs">{sc.head_name}</h3>
                    <div className="text-[10px] text-slate-400">{sc.role}</div>
                  </div>
                  <span className={`text-[10px] font-bold px-2 py-0.5 rounded font-mono ${
                    sc.status === "EXEMPLARY" ? "bg-emerald-500/20 text-emerald-300" : "bg-cyan-500/20 text-cyan-300"
                  }`}>
                    {sc.kpi_score_pct}% Score
                  </span>
                </div>
                <div className="pt-2 border-t border-slate-800/80 flex items-center justify-between text-[10px] text-slate-400 font-mono">
                  <span>SLA Adherence: <strong className="text-emerald-400">{sc.sla_adherence_pct}%</strong></span>
                  <span>Budget Used: <strong className="text-cyan-300">{sc.budget_utilization_pct}%</strong></span>
                  <span>Pending: <strong className="text-amber-400">{sc.pending_tasks_count}</strong></span>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

