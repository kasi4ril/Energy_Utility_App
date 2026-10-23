/**
 * Customer Care Complaint Feed for Head of Distribution Services.
 * Real-time visibility into customer complaints routed to Feeder, Transformer & Metering Units.
 * Strictly <= 350 lines.
 */

import React, { useState, useEffect } from "react";
import { fetchCustomerTickets } from "../customer-care/api";
import { ServiceTicket } from "../customer-care/types";
import {
  AlertTriangle,
  Flame,
  Clock,
  MapPin,
  Phone,
  User,
  ShieldAlert,
  CheckCircle,
  Filter,
  ArrowUpRight,
  Zap,
} from "lucide-react";

interface FeedProps {
  currentUnit?: string;
  onCreateWorkOrderFromTicket?: (ticket: ServiceTicket) => void;
}

export function CustomerCareComplaintFeed({ currentUnit, onCreateWorkOrderFromTicket }: FeedProps) {
  const [tickets, setTickets] = useState<ServiceTicket[]>([]);
  const [filter, setFilter] = useState<"ALL" | "FEEDER" | "TRANSFORMER" | "METER">("ALL");
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    fetchCustomerTickets()
      .then((data) => {
        // Filter tickets belonging to Distribution Services (Feeder, Transformer, Metering)
        const servicesTickets = (data || []).filter(
          (t) =>
            t.unit_type?.includes("FEEDER") ||
            t.unit_type?.includes("TRANSFORMER") ||
            t.unit_type?.includes("METER") ||
            t.assigned_unit_name?.toLowerCase().includes("feeder") ||
            t.assigned_unit_name?.toLowerCase().includes("transformer") ||
            t.assigned_unit_name?.toLowerCase().includes("meter")
        );
        setTickets(servicesTickets);
      })
      .catch(() => {})
      .finally(() => setIsLoading(false));
  }, []);

  const filteredTickets = tickets.filter((t) => {
    if (filter === "FEEDER") {
      return (
        t.unit_type?.includes("FEEDER") ||
        t.assigned_unit_name?.toLowerCase().includes("feeder") ||
        t.category === "LINE_SNAP" ||
        t.category === "FALLEN_POLE"
      );
    }
    if (filter === "TRANSFORMER") {
      return (
        t.unit_type?.includes("TRANSFORMER") ||
        t.assigned_unit_name?.toLowerCase().includes("transformer") ||
        t.category === "TRANSFORMER_FAULT" ||
        t.category === "OIL_LEAK" ||
        t.category === "FUSE_BLOWN"
      );
    }
    if (filter === "METER") {
      return (
        t.unit_type?.includes("METER") ||
        t.assigned_unit_name?.toLowerCase().includes("meter") ||
        t.category === "FAULTY_METER" ||
        t.category === "NEW_METER_SCHEDULE"
      );
    }
    return true;
  });

  return (
    <div className="bg-slate-900/90 border border-slate-800 rounded-xl p-4 sm:p-5 space-y-4">
      {/* Header with Title and Filter Tabs */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 pb-3 border-b border-slate-800">
        <div>
          <div className="flex items-center space-x-2">
            <span className="p-1.5 rounded-lg bg-blue-500/20 text-blue-400 border border-blue-500/30">
              <Zap className="w-4 h-4" />
            </span>
            <h3 className="text-sm font-bold text-slate-100 uppercase tracking-wide">
              Customer Care Complaint Intake Feed &bull; Distribution Services
            </h3>
          </div>
          <p className="text-xs text-slate-400 mt-0.5">
            Head of Distribution Services cross-departmental oversight for Feeder, Transformer &amp; Metering queues
          </p>
        </div>

        {/* Filter Pills */}
        <div className="flex items-center space-x-1.5 bg-slate-950 p-1 rounded-lg border border-slate-800 text-xs">
          {(["ALL", "FEEDER", "TRANSFORMER", "METER"] as const).map((cat) => (
            <button
              key={cat}
              onClick={() => setFilter(cat)}
              className={`px-2.5 py-1 rounded font-semibold transition-all ${
                filter === cat
                  ? "bg-blue-600 text-white shadow-sm"
                  : "text-slate-400 hover:text-slate-200"
              }`}
            >
              {cat === "ALL" ? `All Services (${tickets.length})` : cat}
            </button>
          ))}
        </div>
      </div>

      {isLoading ? (
        <div className="text-center py-8 text-xs text-slate-400">Loading live customer complaints...</div>
      ) : filteredTickets.length === 0 ? (
        <div className="text-center py-8 text-xs text-slate-400 bg-slate-950/40 rounded-lg border border-slate-800">
          No complaints currently routed for this category.
        </div>
      ) : (
        <div className="space-y-3">
          {filteredTickets.map((ticket) => {
            const isEmergency = ticket.is_safety_emergency || ticket.nerc_escalation_tier === "TIER_3";
            return (
              <div
                key={ticket.id}
                className={`p-3.5 rounded-xl border transition-all text-xs ${
                  isEmergency
                    ? "bg-red-950/40 border-red-500/50 shadow-[0_0_12px_rgba(239,68,68,0.2)]"
                    : "bg-slate-950/70 border-slate-800 hover:border-slate-700"
                }`}
              >
                <div className="flex flex-col md:flex-row md:items-start justify-between gap-3">
                  {/* Left Column: Ticket Info & Category */}
                  <div className="space-y-1.5 flex-1">
                    <div className="flex flex-wrap items-center gap-2">
                      <span className="font-mono font-bold text-slate-100">{ticket.ticket_number}</span>
                      <span
                        className={`px-2 py-0.5 rounded text-[10px] font-bold ${
                          isEmergency
                            ? "bg-red-500/20 text-red-300 border border-red-500/40 animate-pulse"
                            : ticket.priority === "HIGH"
                            ? "bg-amber-500/20 text-amber-300 border border-amber-500/40"
                            : "bg-slate-800 text-slate-300"
                        }`}
                      >
                        {isEmergency ? "SAFETY HAZARD (<2H)" : ticket.priority}
                      </span>
                      <span className="text-[10px] bg-blue-900/40 text-blue-300 px-2 py-0.5 rounded border border-blue-700/50">
                        {ticket.category?.replace(/_/g, " ")}
                      </span>
                      <span className="text-[10px] text-slate-400 font-mono">
                        Target: <strong className="text-slate-300">{ticket.assigned_unit_name}</strong>
                      </span>
                    </div>

                    <p className="text-slate-200 font-sans leading-relaxed text-[11px] font-medium">
                      {ticket.description}
                    </p>

                    {/* Customer & Location Context */}
                    <div className="flex flex-wrap items-center gap-x-4 gap-y-1 text-[11px] text-slate-400">
                      <span className="flex items-center space-x-1">
                        <User className="w-3 h-3 text-slate-500" />
                        <span>{ticket.customer_name}</span>
                      </span>
                      <span className="flex items-center space-x-1">
                        <Phone className="w-3 h-3 text-slate-500" />
                        <span className="font-mono">{ticket.customer_phone}</span>
                      </span>
                      <span className="flex items-center space-x-1">
                        <MapPin className="w-3 h-3 text-slate-500" />
                        <span>{ticket.service_address}</span>
                      </span>
                      {ticket.landmark && (
                        <span className="text-blue-300 font-medium">
                          &bull; Landmark: {ticket.landmark}
                        </span>
                      )}
                    </div>
                  </div>

                  {/* Right Column: Feeder/DT Code & Action */}
                  <div className="flex flex-row md:flex-col items-end justify-between md:justify-start gap-2 shrink-0">
                    <div className="text-right font-mono text-[10px] text-slate-400">
                      <div>Feeder: <span className="text-sky-300 font-bold">{ticket.feeder_line_code || "N/A"}</span></div>
                      <div>DT: <span className="text-amber-300">{ticket.dt_station_code || "N/A"}</span></div>
                    </div>

                    <button
                      onClick={() => {
                        if (onCreateWorkOrderFromTicket) {
                          onCreateWorkOrderFromTicket(ticket);
                        } else {
                          alert(`Assigned ticket ${ticket.ticket_number} to immediate field dispatch queue.`);
                        }
                      }}
                      className="px-2.5 py-1.5 bg-blue-600/80 hover:bg-blue-600 text-white font-semibold rounded-lg flex items-center space-x-1 text-[11px] transition-colors shadow-sm"
                    >
                      <span>Dispatch / Create WO</span>
                      <ArrowUpRight className="w-3 h-3" />
                    </button>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}

