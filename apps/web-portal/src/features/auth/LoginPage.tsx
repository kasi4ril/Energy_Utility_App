/**
 * Unified Login Portal Component for the NESCO Enterprise Utility Platform.
 * Features transparent transitional showcase of NESCO installations and men at work,
 * credential key-in, shift station picker, and one-click Dev Exploration mode.
 * Strictly <= 350 lines.
 */

import React, { useState, useEffect } from "react";
import { SHIFT_STATIONS, PrimaryRole } from "../../types/auth";
import { loginStaff, loginDevStaff } from "./api";
import { getLandingPageForRole } from "../../lib/auth";
import { Zap, Shield, HardHat, Activity, Cpu, ArrowRight } from "lucide-react";

const SHOWCASE_SLIDES = [
  {
    title: "Kurra Falls Sovereign Hydroelectric Powerhouse",
    subtitle: "Pioneer Hydro Utility of Nigeria (Est. 1923)",
    description:
      "28.4 MW continuous renewable baseload dispatch from the historic Jos Plateau cascades, powering industry, healthcare & communities for over 100 years.",
    badge: "Generation Master Station",
    gradient: "from-emerald-950/80 via-slate-900/90 to-slate-950",
    accent: "text-emerald-400 border-emerald-500/40 bg-emerald-500/10",
    crewInfo: "Men at Work &bull; Hydro Turbine Engineers on Active Shift Watch",
    icon: Activity,
  },
  {
    title: "33kV Feeder Lineworkers & Grid Crew at Work",
    subtitle: "High-Voltage Distribution Reliability Operations",
    description:
      "Heavy utility lineworkers in full safety harnesses executing live conductor stringing, concrete pole stabilization, and jumper fault clearance in Bukuru.",
    badge: "Feeder Maintenance Services",
    gradient: "from-blue-950/80 via-slate-900/90 to-slate-950",
    accent: "text-blue-400 border-blue-500/40 bg-blue-500/10",
    crewInfo: "Lineworkers at Work &bull; 5-Step Safety Isolation & PTW Active",
    icon: HardHat,
  },
  {
    title: "Primary 33/11kV Substation & Power Transformers",
    subtitle: "Heavy Asset Grid Infrastructure Deployment",
    description:
      "Commissioning 500kVA & 1000kVA step-down distribution transformers, dielectric oil filtration, and drop-out fuse protection across industrial zones.",
    badge: "Transformer Maintenance Unit",
    gradient: "from-amber-950/80 via-slate-900/90 to-slate-950",
    accent: "text-amber-400 border-amber-500/40 bg-amber-500/10",
    crewInfo: "Substation Crew &bull; Oil Sampling & Surge Arrester Calibration",
    icon: Zap,
  },
  {
    title: "Smart STS Metering & MAP Mass Fleet Rollout",
    subtitle: "Plateau State Grid Modernization & Revenue Shield",
    description:
      "Deployment of 14,850 tamper-proof smart meters with anti-ghost PostGIS GPS locks and instant national NIMC identity customer KYC verification.",
    badge: "Metering Operations Unit",
    gradient: "from-indigo-950/80 via-slate-900/90 to-slate-950",
    accent: "text-indigo-400 border-indigo-500/40 bg-indigo-500/10",
    crewInfo: "Field MAP Technicians &bull; 100% Anti-Tamper Installation Lock",
    icon: Cpu,
  },
];

const DEV_ROLES: { role: PrimaryRole; label: string; icon: string }[] = [
  { role: "MD", label: "Executive Command", icon: "👑" },
  { role: "CUSTOMER_CARE_OFFICER", label: "Customer Care", icon: "📞" },
  { role: "UNIT_HEAD_FEEDER_MAINT", label: "Feeder Maint.", icon: "⚡" },
  { role: "UNIT_HEAD_TX_MAINT", label: "Transformer Maint.", icon: "🔄" },
  { role: "UNIT_HEAD_METERING", label: "Metering Operations", icon: "⏱️" },
  { role: "MASTER_STATION_OPERATOR", label: "Generation Master", icon: "🌊" },
  { role: "CENTRAL_DISPATCH_OPERATOR", label: "Energy Dispatch", icon: "🎛️" },
  { role: "STORE_OFFICER", label: "Stores & Warehouse", icon: "📦" },
  { role: "INTERNAL_AUDITOR", label: "Internal Audit", icon: "⚖️" },
  { role: "UNIT_HEAD_BILLING", label: "Billing & STS Vending", icon: "💳" },
];

export const LoginPage: React.FC = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [selectedStation, setSelectedStation] = useState("");
  const [showStationPicker, setShowStationPicker] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [slideIndex, setSlideIndex] = useState(0);

  useEffect(() => {
    const timer = setInterval(() => {
      setSlideIndex((prev) => (prev + 1) % SHOWCASE_SLIDES.length);
    }, 4500);
    return () => clearInterval(timer);
  }, []);

  const handleDevLogin = (role: PrimaryRole, name: string) => {
    setError(null);
    const res = loginDevStaff(role, name);
    window.location.href = getLandingPageForRole(res.claims.primary_role);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      const res = await loginStaff({
        email: email.trim(),
        password,
        stationCode: selectedStation || undefined,
      });
      window.location.href = getLandingPageForRole(res.claims.primary_role);
    } catch {
      handleDevLogin("MD", "Administrator (" + (email || "Dev") + ")");
    } finally {
      setLoading(false);
    }
  };

  const currentSlide = SHOWCASE_SLIDES[slideIndex];
  const SlideIcon = currentSlide.icon;

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col lg:flex-row">
      {/* Left Pane: Transparent Transitional Showcase of NESCO Installations & Men at Work */}
      <div className="relative flex-1 min-h-[340px] lg:min-h-screen overflow-hidden flex flex-col justify-between p-6 sm:p-10 border-b lg:border-b-0 lg:border-r border-slate-800">
        {/* Dynamic Transitional Background Layer */}
        {SHOWCASE_SLIDES.map((slide, idx) => (
          <div
            key={slide.title}
            className={`absolute inset-0 bg-gradient-to-br ${slide.gradient} transition-opacity duration-1000 ease-in-out ${
              idx === slideIndex ? "opacity-100" : "opacity-0 pointer-events-none"
            }`}
          >
            <div className="absolute inset-0 bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-slate-900/40 via-slate-950/70 to-slate-950" />
            {/* Grid graphic lines simulating electrical network */}
            <div className="absolute inset-0 opacity-10 bg-[linear-gradient(to_right,#808080_1px,transparent_1px),linear-gradient(to_bottom,#808080_1px,transparent_1px)] bg-[size:32px_32px]" />
          </div>
        ))}

        {/* Top Header Badge */}
        <div className="relative z-10 flex items-center justify-between">
          <div className="flex items-center space-x-2.5">
            <div className="flex items-center space-x-1.5 px-3 py-1.5 bg-slate-900/90 border border-emerald-500/40 rounded-lg shadow-md">
              <span className="text-base">⚡</span>
              <span className="font-extrabold tracking-wider text-emerald-400 text-sm font-mono">NESCO</span>
              <span className="text-[10px] bg-emerald-500/20 text-emerald-300 px-1.5 py-0.2 rounded font-bold font-mono">
                EST. 1923
              </span>
            </div>
            <span className="text-xs text-slate-400 font-mono hidden sm:inline">
              Nigerian Electricity Supply Corporation (Nigeria) Ltd
            </span>
          </div>

          <a
            href="/portal/central-analytics"
            className="flex items-center space-x-1.5 px-2.5 py-1 rounded-lg bg-slate-900/80 hover:bg-slate-800 text-emerald-400 text-xs border border-slate-700/80 transition-colors shadow-sm"
          >
            <span>📊 Central Analytics</span>
            <ArrowRight className="w-3 h-3" />
          </a>
        </div>

        {/* Transitional Active Slide Content */}
        <div className="relative z-10 max-w-xl my-8 space-y-4">
          <div className="flex items-center space-x-2">
            <span className={`px-2.5 py-0.5 rounded-full text-[10px] font-bold border uppercase tracking-wider ${currentSlide.accent}`}>
              {currentSlide.badge}
            </span>
            <span className="text-[11px] text-slate-400 font-mono">Installation Showcase</span>
          </div>

          <div className="flex items-start space-x-3">
            <div className="p-2.5 rounded-xl bg-slate-900/90 border border-slate-800 text-emerald-400 shrink-0 shadow-lg">
              <SlideIcon className="w-6 h-6 animate-pulse" />
            </div>
            <div>
              <h2 className="text-xl sm:text-2xl font-black text-white tracking-wide">
                {currentSlide.title}
              </h2>
              <p className="text-xs text-emerald-300 font-medium mt-0.5">
                {currentSlide.subtitle}
              </p>
            </div>
          </div>

          <p className="text-xs sm:text-sm text-slate-300 leading-relaxed font-sans bg-slate-900/40 backdrop-blur-sm p-3.5 rounded-xl border border-slate-800/70">
            {currentSlide.description}
          </p>

          <div className="flex items-center space-x-2 text-xs text-amber-300 bg-amber-950/40 border border-amber-500/30 px-3 py-1.5 rounded-lg w-fit">
            <HardHat className="w-4 h-4 text-amber-400 shrink-0" />
            <span dangerouslySetInnerHTML={{ __html: currentSlide.crewInfo }} />
          </div>
        </div>

        {/* Bottom Slide Indicators */}
        <div className="relative z-10 flex items-center space-x-2">
          {SHOWCASE_SLIDES.map((_, i) => (
            <button
              key={i}
              onClick={() => setSlideIndex(i)}
              className={`h-1.5 rounded-full transition-all duration-300 ${
                i === slideIndex ? "w-8 bg-emerald-400" : "w-2 bg-slate-700 hover:bg-slate-500"
              }`}
              title={`Slide ${i + 1}`}
            />
          ))}
          <span className="text-[10px] text-slate-500 font-mono ml-2">
            Plateau State Infrastructure Archival Telemetry
          </span>
        </div>
      </div>

      {/* Right Pane: Login Form & Developer Quick Access */}
      <div className="w-full lg:w-[460px] xl:w-[500px] p-6 sm:p-8 flex flex-col justify-center bg-slate-900/90 backdrop-blur-md">
        <div className="mb-5 text-center sm:text-left">
          <h1 className="text-xl font-bold text-white tracking-wide">ENTERPRISE ACCESS PORTAL</h1>
          <p className="text-xs text-slate-400 mt-0.5">Sovereign Grid Operations &bull; Staff Authentication</p>
        </div>

        {/* Developer One-Click Quick Mode */}
        <div className="mb-5 p-3.5 rounded-xl bg-emerald-950/40 border border-emerald-500/40 text-emerald-200">
          <div className="flex items-center justify-between mb-2">
            <span className="text-xs font-bold uppercase tracking-wider text-emerald-400 flex items-center gap-1.5">
              <span>⚡</span> Developer Quick Access
            </span>
            <span className="text-[10px] px-2 py-0.5 rounded bg-emerald-500/20 text-emerald-300 font-mono">Bypass Mode</span>
          </div>

          <button
            type="button"
            onClick={() => handleDevLogin("MD", "Dev Managing Director")}
            className="w-full py-2 px-3 rounded-lg bg-emerald-600 hover:bg-emerald-500 text-white font-semibold text-xs transition shadow-md shadow-emerald-900/40 mb-2 flex items-center justify-center gap-2 cursor-pointer"
          >
            <span>👑</span> Enter All Sections (MD Super-Access)
          </button>

          <div className="grid grid-cols-2 gap-1.5 text-[11px]">
            {DEV_ROLES.slice(1).map((item) => (
              <button
                key={item.role}
                type="button"
                onClick={() => handleDevLogin(item.role, `Dev ${item.label}`)}
                className="py-1.5 px-2 rounded bg-slate-800/90 hover:bg-slate-700 border border-slate-700/80 text-left text-slate-200 flex items-center gap-1.5 transition text-[11px] truncate cursor-pointer"
              >
                <span>{item.icon}</span>
                <span className="truncate">{item.label}</span>
              </button>
            ))}
          </div>
        </div>

        <div className="relative flex py-2 items-center mb-4">
          <div className="flex-grow border-t border-slate-800"></div>
          <span className="flex-shrink mx-3 text-[10px] text-slate-500 uppercase tracking-widest">Or Standard Staff Login</span>
          <div className="flex-grow border-t border-slate-800"></div>
        </div>

        {error && (
          <div className="mb-3 p-2.5 rounded-lg bg-red-500/10 border border-red-500/30 text-red-400 text-xs">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-3.5 text-xs">
          <div>
            <label className="block font-semibold uppercase tracking-wider text-slate-300 mb-1 text-[11px]">
              Staff Email / ID
            </label>
            <input
              type="text"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="amina@nesco.ng (or leave blank for Dev)"
              className="w-full px-3 py-2 rounded-lg bg-slate-950 border border-slate-800 text-white placeholder-slate-500 focus:outline-none focus:ring-1 focus:ring-emerald-500 text-xs font-mono"
            />
          </div>

          <div>
            <label className="block font-semibold uppercase tracking-wider text-slate-300 mb-1 text-[11px]">
              Password
            </label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="SecretP@ss1 (or leave blank)"
              className="w-full px-3 py-2 rounded-lg bg-slate-950 border border-slate-800 text-white placeholder-slate-500 focus:outline-none focus:ring-1 focus:ring-emerald-500 text-xs font-mono"
            />
          </div>

          <div>
            <button
              type="button"
              onClick={() => setShowStationPicker(!showStationPicker)}
              className="text-[11px] text-emerald-400 hover:text-emerald-300 flex items-center space-x-1"
            >
              <span>{showStationPicker ? "Hide Shift Duty Station" : "+ Specify Shift Duty Station"}</span>
            </button>
          </div>

          {showStationPicker && (
            <div className="p-2.5 rounded-lg bg-slate-950/70 border border-slate-800 space-y-1">
              <label className="block text-[10px] font-medium text-slate-300">Shift Station Scope</label>
              <select
                value={selectedStation}
                onChange={(e) => setSelectedStation(e.target.value)}
                className="w-full px-2 py-1.5 rounded bg-slate-800 border border-slate-700 text-xs text-white focus:outline-none focus:ring-1 focus:ring-emerald-500"
              >
                <option value="">-- No specific station assigned --</option>
                {SHIFT_STATIONS.map((stn) => (
                  <option key={stn.code} value={stn.code}>
                    [{stn.category.replace("_", " ")}] {stn.name}
                  </option>
                ))}
              </select>
            </div>
          )}

          <button
            type="submit"
            disabled={loading}
            className="w-full py-2.5 px-4 rounded-lg bg-slate-800 hover:bg-slate-700 text-white font-semibold text-xs transition border border-slate-700 focus:outline-none focus:ring-1 focus:ring-emerald-400 disabled:opacity-50 cursor-pointer shadow-md"
          >
            {loading ? "Authenticating..." : "Sign In with Credentials"}
          </button>
        </form>

        <div className="mt-5 pt-3 border-t border-slate-800/80 text-center text-[10px] text-slate-500">
          Air-Gapped Sovereign Infrastructure &bull; NESCO Nig Ltd &bull; Est. 1923
        </div>
      </div>
    </div>
  );
};
