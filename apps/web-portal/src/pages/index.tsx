/**
 * NESCO Enterprise Utility Platform - Public Home & Landing Page (/).
 * Sovereign Grid Operating System for Nigeria Electricity Supply Corporation (NESCO) Nig Ltd.
 */

import React, { useState, useEffect } from "react";
import Head from "next/head";
import Link from "next/link";
import { LandingGridTicker } from "../components/landing/LandingGridTicker";
import { LandingDepartmentGrid } from "../components/landing/LandingDepartmentGrid";
import { LandingArchitecturePillars } from "../components/landing/LandingArchitecturePillars";
import { apiClient } from "../lib/api-client";
import { parseTokenClaims } from "../lib/auth";
import { AuthClaims } from "../types/auth";

export default function HomePage() {
  const [claims, setClaims] = useState<AuthClaims | null>(null);
  const [watTime, setWatTime] = useState<string>("");

  useEffect(() => {
    const token = apiClient.getToken();
    if (token) setClaims(parseTokenClaims(token));

    const updateTime = () => {
      const now = new Date();
      const timeStr = now.toLocaleTimeString("en-GB", {
        timeZone: "Africa/Lagos",
        hour12: false,
        hour: "2-digit",
        minute: "2-digit",
        second: "2-digit",
      });
      const dateStr = now.toLocaleDateString("en-GB", {
        timeZone: "Africa/Lagos",
        day: "2-digit",
        month: "short",
        year: "numeric",
      });
      setWatTime(`${dateStr} &bull; ${timeStr} WAT`);
    };

    updateTime();
    const timer = setInterval(updateTime, 1000);
    return () => clearInterval(timer);
  }, []);

  return (
    <>
      <Head>
        <title>NESCO Enterprise Utility Platform &bull; Sovereign Grid OS (Est. 1923)</title>
        <meta
          name="description"
          content="Enterprise Utility Workflow, Grid Intelligence & Local AI Governance for NESCO Nigeria Ltd."
        />
      </Head>

      <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans selection:bg-cyan-500/30">
        {/* Top Ticker Announcement Bar */}
        <div className="bg-slate-950 border-b border-slate-800/80 py-2 px-6 flex items-center justify-between text-xs select-none">
          <div className="flex items-center space-x-3">
            <span className="flex items-center space-x-1.5 font-mono text-cyan-400">
              <span className="w-2 h-2 rounded-full bg-cyan-400 animate-pulse" />
              <strong className="text-white">GRID FREQ:</strong> 50.02 Hz (Nominal Stability)
            </span>
            <span className="text-slate-700 hidden sm:inline">|</span>
            <span className="text-slate-400 hidden sm:inline font-mono">
              BULK POWER: <strong className="text-emerald-400">63.6 MW Total</strong>
            </span>
            <span className="text-slate-700 hidden md:inline">|</span>
            <span
              className="text-slate-400 hidden md:inline font-mono text-[11px]"
              dangerouslySetInnerHTML={{ __html: watTime || "Loading WAT..." }}
            />
          </div>

          <div className="flex items-center space-x-3 text-xs">
            {claims ? (
              <span className="text-slate-300">
                Logged in as: <strong className="text-white">{claims.full_name}</strong> (
                <span className="text-cyan-400 font-mono text-[11px]">{claims.primary_role}</span>)
              </span>
            ) : (
              <Link
                href="/login"
                className="text-cyan-400 hover:text-cyan-300 font-bold transition flex items-center space-x-1"
              >
                <span>Staff Authentication</span>
                <span>&rarr;</span>
              </Link>
            )}
          </div>
        </div>

        {/* Global Navigation Header */}
        <nav className="h-16 border-b border-slate-800/80 px-6 flex items-center justify-between max-w-7xl w-full mx-auto select-none">
          <div className="flex items-center space-x-3">
            <div className="w-9 h-9 rounded-xl bg-cyan-950/80 border border-cyan-500/40 flex items-center justify-center text-cyan-400 font-black shadow-inner shadow-cyan-900/40 text-lg">
              ⚡
            </div>
            <div>
              <div className="flex items-center space-x-2">
                <span className="text-base font-black text-white tracking-wider">NESCO</span>
                <span className="text-[10px] font-mono text-cyan-400 bg-cyan-950/60 border border-cyan-800/60 px-1.5 py-0.2 rounded font-bold">
                  EST. 1923
                </span>
              </div>
              <div className="text-[10px] text-slate-400 font-mono tracking-tight">
                SOVEREIGN GRID OPERATING SYSTEM
              </div>
            </div>
          </div>

          <div className="flex items-center space-x-4">
            <Link
              href="/portal/central-analytics"
              className="text-xs text-slate-300 hover:text-cyan-300 transition hidden md:flex items-center space-x-1 font-medium"
            >
              <span>📊</span>
              <span>Central Analytics</span>
            </Link>
            <Link
              href="/field-pwa/tasks"
              className="text-xs text-slate-400 hover:text-white transition hidden md:flex items-center space-x-1 font-medium"
            >
              <span>📱</span>
              <span>Field PWA</span>
            </Link>
            <Link
              href={claims ? "/portal/executive-command" : "/login"}
              className="px-4 py-2 rounded-lg bg-cyan-600 hover:bg-cyan-500 text-white font-bold text-xs transition shadow-lg shadow-cyan-950 flex items-center space-x-1.5"
            >
              <span>{claims ? "👑 My Portal" : "Enterprise Sign In"}</span>
            </Link>
          </div>
        </nav>

        {/* Hero Section */}
        <header className="py-16 px-6 text-center max-w-4xl mx-auto space-y-6 select-none">
          <div className="inline-flex items-center space-x-2 px-3.5 py-1 rounded-full bg-slate-900 border border-slate-800 text-slate-300 text-xs shadow-inner">
            <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse" />
            <span className="font-medium">
              Air-Gapped Sovereign Utility Architecture &bull; 100% On-Premise
            </span>
          </div>

          <h1 className="text-4xl sm:text-5xl font-black text-white tracking-tight leading-tight">
            Next-Generation Grid Intelligence &amp; Operational Governance
          </h1>

          <p className="text-base sm:text-lg text-slate-400 max-w-2xl mx-auto leading-relaxed">
            Unifying clean hydroelectric generation (Kurra Falls Plants 1–6), transmission dispatch,
            substation maintenance, commercial STS vending, and C-Suite command under an immutable
            SHA-256 hash-chain ledger.
          </p>

          <div className="pt-2 flex flex-wrap items-center justify-center gap-3 sm:gap-4">
            <Link
              href="/portal/executive-command"
              className="px-5 py-3 rounded-xl bg-cyan-600 hover:bg-cyan-500 text-white font-bold text-xs sm:text-sm transition shadow-xl shadow-cyan-950 flex items-center space-x-2"
            >
              <span>👑</span>
              <span>Launch Executive Command</span>
            </Link>
            <Link
              href="/portal/central-analytics"
              className="px-5 py-3 rounded-xl bg-slate-900 hover:bg-slate-800 border border-slate-700 text-slate-200 font-bold text-xs sm:text-sm transition flex items-center space-x-2"
            >
              <span>📊</span>
              <span>Sovereign Central Analytics</span>
            </Link>
            <Link
              href="/field-pwa/tasks"
              className="px-5 py-3 rounded-xl bg-slate-950 hover:bg-slate-900 border border-slate-800 text-slate-400 hover:text-white font-semibold text-xs sm:text-sm transition flex items-center space-x-2"
            >
              <span>📱</span>
              <span>Technician PWA</span>
            </Link>
          </div>
        </header>

        {/* Live Grid Metrics Pulse Ribbon */}
        <LandingGridTicker
          frequencyHz={50.02}
          hydroMw={48.6}
          tcnMw={15.0}
          activeFeeders="28 / 32"
          collectionsNgn="₦4.85M"
          compliancePct={96.8}
        />

        {/* 10 Operational Department Portals Launchpad */}
        <LandingDepartmentGrid />

        {/* 4 Architectural Engineering Pillars */}
        <LandingArchitecturePillars />

        {/* Historical Heritage Banner */}
        <section className="border-t border-slate-800/80 bg-gradient-to-r from-slate-950 via-slate-900/60 to-slate-950 py-8 px-6 text-center select-none">
          <div className="max-w-4xl mx-auto space-y-2">
            <div className="text-xs font-mono font-bold text-amber-400 uppercase tracking-widest">
              Over A Century of Pioneer Clean Energy (1923 – 2026)
            </div>
            <p className="text-xs text-slate-400 leading-relaxed max-w-2xl mx-auto">
              NESCO (Nigeria Electricity Supply Corporation) established Nigeria's first hydroelectric
              power network in 1923 at Kurra Falls on the Jos Plateau. Operating six synchronous
              hydro generation stations with 100% sovereign air-gapped grid intelligence.
            </p>
          </div>
        </section>

        {/* Global Footer */}
        <footer className="border-t border-slate-800/80 py-6 px-6 text-xs text-slate-500 select-none bg-slate-950">
          <div className="max-w-7xl mx-auto flex flex-col sm:flex-row items-center justify-between gap-3">
            <div>
              &copy; {new Date().getFullYear()} <strong>NESCO Nigeria Ltd</strong> &bull; Established 1923
              &bull; Pioneer Hydroelectric Power Utility in Nigeria
            </div>
            <div className="flex items-center space-x-4 font-mono text-[11px]">
              <span className="text-cyan-400 font-bold">Platform v1.5.0</span>
              <span>&bull;</span>
              <span>Air-Gapped Sovereign Deployment</span>
              <span>&bull;</span>
              <span className="text-emerald-400">NERC &amp; NEMSA Certified</span>
            </div>
          </div>
        </footer>
      </div>
    </>
  );
}
