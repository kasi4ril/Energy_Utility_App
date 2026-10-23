/**
 * Extended Analytics & Telemetry Dataset for Customer Care.
 * Compliant with Observations 1, 2, 3 & 6:
 * - Period-based query metrics (Daily, Weekly, Monthly, Quarterly, Yearly)
 * - Multi-year recorded data (2026, 2025, 2024, 2023)
 * - Quarters Q1 (Jan-Mar), Q2 (Apr-Jun), Q3 (Jul-Sep), Q4 (Oct-Dec)
 * - Full 12-Month spectrum (January - December)
 * - Daily granularity for any given month
 * - Postpaid payment and Prepaid vending analytics by month
 * - Meter installation milestones (14,850 total, category and week/month breakdown)
 * - National Identity (NIN/NIMC) auto-verification directory
 * Strictly <= 350 lines.
 */

import {
  PeriodMetrics,
  AnalyticsPeriod,
  BillingVendingStats,
  MeterCategoryStats,
  NINVerificationResult,
} from "./types";

export const YEARLY_METRICS_DATA: Record<string, PeriodMetrics> = {
  "2026": {
    period: "YEARLY",
    label: "Operational Year 2026 (YTD)",
    lodged_count: 4180,
    resolved_count: 3970,
    unresolved_count: 210,
    new_customers_count: 1420,
    resolved_pct: 95.0,
    unresolved_pct: 5.0,
  },
  "2025": {
    period: "YEARLY",
    label: "Operational Year 2025 (Historical Audit)",
    lodged_count: 4890,
    resolved_count: 4520,
    unresolved_count: 370,
    new_customers_count: 1680,
    resolved_pct: 92.4,
    unresolved_pct: 7.6,
  },
  "2024": {
    period: "YEARLY",
    label: "Operational Year 2024 (Historical Audit)",
    lodged_count: 4620,
    resolved_count: 4190,
    unresolved_count: 430,
    new_customers_count: 1510,
    resolved_pct: 90.7,
    unresolved_pct: 9.3,
  },
  "2023": {
    period: "YEARLY",
    label: "Centenary Year 2023 (Historical Audit)",
    lodged_count: 4310,
    resolved_count: 3820,
    unresolved_count: 490,
    new_customers_count: 1390,
    resolved_pct: 88.6,
    unresolved_pct: 11.4,
  },
};

export function getTicketHistogramData(
  year: string = "2026",
  monthKey: string = "10",
  day: string = "23"
): { label: string; lodged: number; resolved: number }[] {
  if (monthKey === "ALL") {
    const mNames = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
    const yRatio = (YEARLY_METRICS_DATA[year]?.lodged_count || 4180) / 4180;
    return mNames.map((m, idx) => {
      const k = (idx + 1).toString().padStart(2, "0");
      const base = MONTHLY_METRICS_DATA[k] || { lodged_count: 340, resolved_count: 310 };
      return {
        label: m,
        lodged: Math.round(base.lodged_count * yRatio),
        resolved: Math.round(base.resolved_count * yRatio),
      };
    });
  }

  if (day === "ALL") {
    const mData = getMonthMetricsForYear(monthKey, year);
    const avgLodged = Math.round(mData.lodged_count / 4);
    const avgResolved = Math.round(mData.resolved_count / 4);
    return [
      { label: "W1 (1-7)", lodged: avgLodged - 4, resolved: avgResolved - 3 },
      { label: "W2 (8-14)", lodged: avgLodged + 2, resolved: avgResolved + 1 },
      { label: "W3 (15-21)", lodged: avgLodged + 5, resolved: avgResolved + 4 },
      { label: "W4 (22-end)", lodged: avgLodged - 3, resolved: avgResolved - 2 },
    ];
  }

  const dayNum = parseInt(day, 10) || 23;
  const dMetrics = getDailyMetricsForMonth(monthKey, dayNum, year);
  const intervals = ["08h", "10h", "12h", "14h", "16h", "18h"];
  return intervals.map((time, i) => {
    const l = Math.max(1, Math.round(dMetrics.lodged_count / 6) + (i % 2 === 0 ? 1 : -1));
    const r = Math.max(0, l - (i === 4 || i === 5 ? 1 : 0));
    return { label: time, lodged: l, resolved: r };
  });
}

export const MONTHLY_METRICS_DATA: Record<string, PeriodMetrics> = {
  "01": { period: "MONTHLY", label: "January 2026", lodged_count: 310, resolved_count: 288, unresolved_count: 22, new_customers_count: 102, resolved_pct: 92.9, unresolved_pct: 7.1 },
  "02": { period: "MONTHLY", label: "February 2026", lodged_count: 325, resolved_count: 305, unresolved_count: 20, new_customers_count: 106, resolved_pct: 93.8, unresolved_pct: 6.2 },
  "03": { period: "MONTHLY", label: "March 2026", lodged_count: 345, resolved_count: 322, unresolved_count: 23, new_customers_count: 112, resolved_pct: 93.3, unresolved_pct: 6.7 },
  "04": { period: "MONTHLY", label: "April 2026", lodged_count: 340, resolved_count: 318, unresolved_count: 22, new_customers_count: 118, resolved_pct: 93.5, unresolved_pct: 6.5 },
  "05": { period: "MONTHLY", label: "May 2026", lodged_count: 352, resolved_count: 330, unresolved_count: 22, new_customers_count: 122, resolved_pct: 93.8, unresolved_pct: 6.2 },
  "06": { period: "MONTHLY", label: "June 2026", lodged_count: 368, resolved_count: 344, unresolved_count: 24, new_customers_count: 125, resolved_pct: 93.5, unresolved_pct: 6.5 },
  "07": { period: "MONTHLY", label: "July 2026", lodged_count: 360, resolved_count: 340, unresolved_count: 20, new_customers_count: 128, resolved_pct: 94.4, unresolved_pct: 5.6 },
  "08": { period: "MONTHLY", label: "August 2026", lodged_count: 364, resolved_count: 346, unresolved_count: 18, new_customers_count: 132, resolved_pct: 95.1, unresolved_pct: 4.9 },
  "09": { period: "MONTHLY", label: "September 2026", lodged_count: 374, resolved_count: 355, unresolved_count: 19, new_customers_count: 133, resolved_pct: 94.9, unresolved_pct: 5.1 },
  "10": { period: "MONTHLY", label: "October 2026 (Present Month)", lodged_count: 348, resolved_count: 296, unresolved_count: 52, new_customers_count: 115, resolved_pct: 85.1, unresolved_pct: 14.9 },
  "11": { period: "MONTHLY", label: "November 2026", lodged_count: 330, resolved_count: 315, unresolved_count: 15, new_customers_count: 98, resolved_pct: 95.5, unresolved_pct: 4.5 },
  "12": { period: "MONTHLY", label: "December 2026", lodged_count: 364, resolved_count: 345, unresolved_count: 19, new_customers_count: 129, resolved_pct: 94.8, unresolved_pct: 5.2 },
};

export function getMonthMetricsForYear(monthKey: string = "10", year: string = "2026"): PeriodMetrics {
  const base = MONTHLY_METRICS_DATA[monthKey] || MONTHLY_METRICS_DATA["10"];
  const monthNames: Record<string, string> = {
    "01": "January", "02": "February", "03": "March", "04": "April",
    "05": "May", "06": "June", "07": "July", "08": "August",
    "09": "September", "10": "October", "11": "November", "12": "December",
  };
  const mName = monthNames[monthKey] || "October";
  if (year === "2026") {
    return base;
  }
  const yearMetrics = YEARLY_METRICS_DATA[year] || YEARLY_METRICS_DATA["2025"];
  const ratio = yearMetrics.lodged_count / 4180;
  const lodged = Math.round(base.lodged_count * ratio);
  const unresolved = Math.max(5, Math.round(base.unresolved_count * ratio));
  const resolved = lodged - unresolved;
  const newCust = Math.round(base.new_customers_count * ratio);
  const resolvedPct = parseFloat(((resolved / lodged) * 100).toFixed(1));
  const unresolvedPct = parseFloat((100 - resolvedPct).toFixed(1));
  return {
    period: "MONTHLY",
    label: `${mName} ${year} Monthly Report`,
    lodged_count: lodged,
    resolved_count: resolved,
    unresolved_count: unresolved,
    new_customers_count: newCust,
    resolved_pct: resolvedPct,
    unresolved_pct: unresolvedPct,
  };
}

export function getDailyMetricsForMonth(monthKey: string = "10", day: number = 23, year: string = "2026"): PeriodMetrics {
  const monthNames: Record<string, string> = {
    "01": "January", "02": "February", "03": "March", "04": "April",
    "05": "May", "06": "June", "07": "July", "08": "August",
    "09": "September", "10": "October", "11": "November", "12": "December",
  };
  const mName = monthNames[monthKey] || "October";

  // October 23 (Today's live grid cycle)
  if (monthKey === "10" && day === 23) {
    return {
      period: "DAILY",
      label: `Today's Grid Cycle (${day} ${mName} ${year})`,
      lodged_count: 14,
      resolved_count: 12,
      unresolved_count: 2,
      new_customers_count: 4,
      resolved_pct: 85.7,
      unresolved_pct: 14.3,
    };
  }

  // Deterministic daily values for any chosen day & month
  const seed = (parseInt(monthKey, 10) * 7 + day * 13) % 20;
  const lodged = 10 + seed;
  const unresolved = Math.max(1, (seed % 4));
  const resolved = lodged - unresolved;
  const newCust = Math.max(1, ((seed + 2) % 6) + 1);
  const resolvedPct = parseFloat(((resolved / lodged) * 100).toFixed(1));
  const unresolvedPct = parseFloat((100 - resolvedPct).toFixed(1));

  return {
    period: "DAILY",
    label: `Day ${day} of ${mName} ${year}`,
    lodged_count: lodged,
    resolved_count: resolved,
    unresolved_count: unresolved,
    new_customers_count: newCust,
    resolved_pct: resolvedPct,
    unresolved_pct: unresolvedPct,
  };
}

export const PERIOD_METRICS_MAP: Record<AnalyticsPeriod, PeriodMetrics> = {
  DAILY: getDailyMetricsForMonth("10", 23, "2026"),
  WEEKLY: {
    period: "WEEKLY",
    label: "Current Operating Week (Week 43)",
    lodged_count: 86,
    resolved_count: 76,
    unresolved_count: 10,
    new_customers_count: 28,
    resolved_pct: 88.4,
    unresolved_pct: 11.6,
  },
  MONTHLY: MONTHLY_METRICS_DATA["10"],
  QUARTERLY: YEARLY_METRICS_DATA["2026"],
  YEARLY: YEARLY_METRICS_DATA["2026"],
};

export const MONTHLY_BILLING_VENDING_DATA: Record<string, BillingVendingStats> = {
  "2026-10": { month: "October 2026 (Present Month)", postpaid_paid_count: 1420, postpaid_paid_amount: 42600000, postpaid_unpaid_count: 380, postpaid_unpaid_amount: 11400000, prepaid_vended_count: 8450, prepaid_vended_amount: 84500000, prepaid_unvended_count: 1120 },
  "2026-09": { month: "September 2026 (Past Month)", postpaid_paid_count: 1390, postpaid_paid_amount: 41700000, postpaid_unpaid_count: 410, postpaid_unpaid_amount: 12300000, prepaid_vended_count: 8210, prepaid_vended_amount: 82100000, prepaid_unvended_count: 1360 },
  "2026-08": { month: "August 2026 (Past Month)", postpaid_paid_count: 1350, postpaid_paid_amount: 40500000, postpaid_unpaid_count: 450, postpaid_unpaid_amount: 13500000, prepaid_vended_count: 7980, prepaid_vended_amount: 79800000, prepaid_unvended_count: 1590 },
  "2026-07": { month: "July 2026 (Past Month)", postpaid_paid_count: 1310, postpaid_paid_amount: 39300000, postpaid_unpaid_count: 490, postpaid_unpaid_amount: 14700000, prepaid_vended_count: 7800, prepaid_vended_amount: 78000000, prepaid_unvended_count: 1770 },
};

export const METER_INSTALLATION_STATS: MeterCategoryStats = {
  total_installed: 14850,
  residential_count: 11200,
  md_count: 850,
  industrial_count: 2800,
  period_view: "MONTHLY",
  period_installs: [
    { period_label: "Oct 2026 (M10)", residential: 240, md: 18, industrial: 62, total: 320 },
    { period_label: "Sep 2026 (M09)", residential: 280, md: 22, industrial: 75, total: 377 },
    { period_label: "Aug 2026 (M08)", residential: 310, md: 19, industrial: 81, total: 410 },
    { period_label: "Jul 2026 (M07)", residential: 260, md: 15, industrial: 68, total: 343 },
    { period_label: "Week 43 (Current)", residential: 58, md: 4, industrial: 14, total: 76 },
    { period_label: "Week 42 (Past)", residential: 64, md: 5, industrial: 18, total: 87 },
    { period_label: "Week 41 (Past)", residential: 60, md: 4, industrial: 15, total: 79 },
  ],
};

export function getMeterInstallationData(year: string = "2026", monthKey: string = "10") {
  const yRatio = year === "2026" ? 1.0 : year === "2025" ? 0.92 : year === "2024" ? 0.85 : 0.78;
  const mNames: Record<string, string> = { "01": "Jan", "02": "Feb", "03": "Mar", "04": "Apr", "05": "May", "06": "Jun", "07": "Jul", "08": "Aug", "09": "Sep", "10": "Oct", "11": "Nov", "12": "Dec" };

  if (monthKey === "ALL") {
    const total = Math.round(3850 * yRatio);
    const res = Math.round(total * 0.754);
    const md = Math.round(total * 0.057);
    const ind = total - res - md;
    const quarters = [
      { label: `Q1 ${year}`, residential: Math.round(res * 0.23), md: Math.round(md * 0.24), industrial: Math.round(ind * 0.22) },
      { label: `Q2 ${year}`, residential: Math.round(res * 0.26), md: Math.round(md * 0.25), industrial: Math.round(ind * 0.27) },
      { label: `Q3 ${year}`, residential: Math.round(res * 0.28), md: Math.round(md * 0.26), industrial: Math.round(ind * 0.26) },
      { label: `Q4 ${year}`, residential: Math.round(res * 0.23), md: Math.round(md * 0.25), industrial: Math.round(ind * 0.25) },
    ];
    return {
      total_installed: total, residential_count: res, md_count: md, industrial_count: ind,
      period_label: `Annual Rollout ${year}`,
      chart_data: quarters.map(q => ({ ...q, total: q.residential + q.md + q.industrial })),
    };
  }

  const mTotal = Math.round(320 * yRatio);
  const mRes = Math.round(mTotal * 0.75);
  const mMd = Math.round(mTotal * 0.06);
  const mInd = mTotal - mRes - mMd;
  const weeks = [
    { label: "W1 (1-7)", residential: Math.round(mRes * 0.24), md: Math.round(mMd * 0.25), industrial: Math.round(mInd * 0.23) },
    { label: "W2 (8-14)", residential: Math.round(mRes * 0.26), md: Math.round(mMd * 0.25), industrial: Math.round(mInd * 0.27) },
    { label: "W3 (15-21)", residential: Math.round(mRes * 0.27), md: Math.round(mMd * 0.25), industrial: Math.round(mInd * 0.25) },
    { label: "W4 (22+)", residential: Math.round(mRes * 0.23), md: Math.round(mMd * 0.25), industrial: Math.round(mInd * 0.25) },
  ];
  return {
    total_installed: mTotal, residential_count: mRes, md_count: mMd, industrial_count: mInd,
    period_label: `${mNames[monthKey] || "Oct"} ${year} Rollout`,
    chart_data: weeks.map(w => ({ ...w, total: w.residential + w.md + w.industrial })),
  };
}

const NIMC_REGISTRY: Record<string, { name: string; gender: string; dob: string }> = {
  "24891029481": { name: "IBRAHIM DANLADI GYANG", gender: "MALE", dob: "1984-06-12" },
  "88219401923": { name: "GRACE PAM CHOLLOM", gender: "FEMALE", dob: "1991-09-24" },
  "91029481726": { name: "BITRUS EZEKIEL LAR", gender: "MALE", dob: "1978-03-15" },
  "12349018273": { name: "MARYAM SANI ABUBAKAR", gender: "FEMALE", dob: "1987-11-04" },
  "10293847561": { name: "DAUDA YAKUBU SHENDAM", gender: "MALE", dob: "1982-08-19" },
  "55443322110": { name: "CHRISTOPHER CHUKWUMA OKAFOR", gender: "MALE", dob: "1989-12-02" },
};

export function lookupNIMCName(nin: string): NINVerificationResult {
  const clean = nin.replace(/\D/g, "");
  if (clean.length !== 11) {
    return {
      nin: clean,
      verified_name: "",
      gender: "",
      date_of_birth: "",
      status: "INVALID",
    };
  }

  if (NIMC_REGISTRY[clean]) {
    const r = NIMC_REGISTRY[clean];
    return {
      nin: clean,
      verified_name: r.name,
      gender: r.gender,
      date_of_birth: r.dob,
      status: "VERIFIED",
    };
  }

  // Deterministic fallback generation for arbitrary 11-digit NINs
  const surnames = ["DANGOTE", "BELLO", "YUSUF", "AL-HASSAN", "OKONKWO", "OGUNLEWE", "GOWON"];
  const firstnames = ["EMMANUEL", "FATIMA", "SHEHU", "NGOZI", "MUSA", "KELECHI", "BLESSING"];
  const sIdx = (parseInt(clean.slice(0, 3), 10) || 1) % surnames.length;
  const fIdx = (parseInt(clean.slice(3, 6), 10) || 2) % firstnames.length;

  return {
    nin: clean,
    verified_name: `${firstnames[fIdx]} ${surnames[sIdx]}`,
    gender: "VERIFIED",
    date_of_birth: "1986-04-10",
    status: "VERIFIED",
  };
}
