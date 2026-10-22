/**
 * Authentication helper utilities, claims decoding, and role-based landing page resolution.
 * Implements Blueprint Section 10.1 auto-routing logic.
 */

import { AuthClaims, PrimaryRole } from "../types/auth";

/**
 * Returns the exact post-login portal route matching the authenticated staff's primary role.
 */
export function getLandingPageForRole(role: PrimaryRole): string {
  switch (role) {
    case "MD":
    case "HEAD_SERVICES":
    case "HEAD_INTERNAL_AUDIT":
    case "HEAD_COMMERCIAL_OPS":
    case "GM_GENERATION":
    case "HEAD_ENERGY_MGMT":
      return "/executive-command";

    case "CUSTOMER_CARE_OFFICER":
    case "HEAD_CUSTOMER_CARE":
      return "/portal/customer-care";

    case "UNIT_HEAD_FEEDER_MAINT":
      return "/portal/services/feeder";
    case "UNIT_HEAD_TX_MAINT":
      return "/portal/services/transformer";
    case "UNIT_HEAD_METERING":
      return "/portal/services/metering";

    case "TEAM_LEAD_FEEDER":
    case "TEAM_LEAD_TX":
    case "TEAM_LEAD_METERING":
    case "SAFETY_REP_HSE":
      return "/field-pwa/tasks";

    case "MASTER_STATION_OPERATOR":
    case "UNIT_HEAD_GEN_ELECTRICAL":
    case "UNIT_HEAD_GEN_MECHANICAL":
    case "WORKSHOP_LEAD":
    case "HR_OFFICER_PLANT":
      return "/portal/generation-master";

    case "CENTRAL_DISPATCH_OPERATOR":
    case "SUBSTATION_OPERATOR":
      return "/portal/energy-management";

    case "STORE_OFFICER":
      return "/portal/stores";

    case "INTERNAL_AUDITOR":
      return "/portal/internal-audit";

    case "UNIT_HEAD_BILLING":
      return "/portal/billing-commercial";

    default:
      return "/login";
  }
}

/**
 * Decodes base64 URL-encoded claims from an unencrypted JWT token client-side.
 */
export function parseTokenClaims(token: string): AuthClaims | null {
  try {
    const parts = token.split(".");
    if (parts.length !== 3) return null;
    const base64Url = parts[1];
    const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split("")
        .map((c) => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2))
        .join("")
    );
    return JSON.parse(jsonPayload) as AuthClaims;
  } catch {
    return null;
  }
}

/**
 * Evaluates whether a given role belongs to the C-Suite executive tier.
 */
export function isExecutiveRole(role: string): boolean {
  return [
    "MD",
    "HEAD_SERVICES",
    "HEAD_INTERNAL_AUDIT",
    "HEAD_COMMERCIAL_OPS",
    "GM_GENERATION",
    "HEAD_ENERGY_MGMT",
  ].includes(role);
}

