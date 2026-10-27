/**
 * Dynamic Unit Delegator: Services & Maintenance Department (/portal/services/[unit]).
 * Routes to Feeder, Transformer, and Metering Unit Landing Pages.
 */

import React from "react";
import { useRouter } from "next/router";
import FeederUnitLandingPage from "./feeder";
import TransformerUnitLandingPage from "../../../features/services-maint/TransformerUnitLandingPage";
import MeteringUnitLandingPage from "../../../features/services-maint/MeteringUnitLandingPage";

export default function DynamicServicesUnitPage() {
  const router = useRouter();
  const { unit } = router.query;
  const unitStr = typeof unit === "string" ? unit.toLowerCase() : "feeder";

  if (unitStr === "transformer") {
    return <TransformerUnitLandingPage />;
  }

  if (unitStr === "metering") {
    return <MeteringUnitLandingPage />;
  }

  return <FeederUnitLandingPage />;
}
