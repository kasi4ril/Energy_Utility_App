/**
 * Next.js Root Page (/): Directs staff to login or their role-specific portal.
 */

import { useEffect } from "react";
import { useRouter } from "next/router";
import { apiClient } from "../lib/api-client";
import { parseTokenClaims, getLandingPageForRole } from "../lib/auth";

export default function IndexPage() {
  const router = useRouter();

  useEffect(() => {
    const token = apiClient.getToken();
    if (!token) {
      router.replace("/login");
      return;
    }

    const claims = parseTokenClaims(token);
    if (!claims) {
      router.replace("/login");
      return;
    }

    const target = getLandingPageForRole(claims.primary_role);
    router.replace(target);
  }, [router]);

  return (
    <div className="min-h-screen bg-slate-950 flex items-center justify-center text-slate-400 text-xs font-mono">
      Routing to assigned departmental workspace...
    </div>
  );
}

