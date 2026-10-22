/**
 * Google Apps Script Webhook Forwarder: Google Form B (Energy Dispatch & Substation Roll-Call)
 * 
 * Bounded Context: S06 Energy Management & Transmission Dispatch
 * Deployment: Install on Google Form onFormSubmit trigger in Google Apps Script editor.
 * Purpose: Signs submitted form responses with HMAC-SHA256 and POSTs to Core API Gateway.
 */

const CORE_API_URL = "https://api.nesco.ng/api/v1/energy/webhook/google-forms";
const HMAC_SECRET = "nesco-dispatch-hmac-integrity-key-32b!";

function onFormSubmit(e) {
  const itemResponses = e.response.getItemResponses();
  const formData = {};

  for (let i = 0; i < itemResponses.length; i++) {
    const itemResponse = itemResponses[i];
    formData[itemResponse.getItem().getTitle()] = itemResponse.getResponse();
  }

  // Construct Standardized JSON Payload for Central Dispatch & Substation Ingestion
  const payload = {
    shift_id: formData["Shift Identifier"] || ("SHIFT-" + Utilities.formatDate(new Date(), "GMT", "yyyyMMdd-HH")),
    shift_window: formData["Shift Window (MORNING/AFTERNOON/NIGHT)"] || "MORNING",
    operator_id: formData["Dispatch Operator Staff ID"] || "OP-DISPATCH-01",
    generation_mw_in: parseFloat(formData["Total Generation Intake (MW)"] || "0.0"),
    tcn_mw_in: parseFloat(formData["TCN Bulk Intake (MW)"] || "0.0"),
    substation_rollcall: JSON.stringify(parseSubstationRollcall(formData)),
    handover_notes: formData["Dispatch & Grid Handover Notes"] || "",
    feeders: parseFeederReadings(formData)
  };

  const payloadString = JSON.stringify(payload);

  // Compute HMAC SHA-256 Signature
  const signatureBytes = Utilities.computeHmacSha256Signature(payloadString, HMAC_SECRET);
  const signatureHex = signatureBytes.reduce((str, chr) => {
    const hex = (chr < 0 ? chr + 256 : chr).toString(16);
    return str + (hex.length === 1 ? '0' : '') + hex;
  }, '');

  // Dispatch HTTP POST to NESCO Core API Gateway
  const options = {
    method: "post",
    contentType: "application/json",
    headers: {
      "X-Signature-SHA256": signatureHex
    },
    payload: payloadString,
    muteHttpExceptions: true
  };

  try {
    const response = UrlFetchApp.fetch(CORE_API_URL, options);
    Logger.log("Core API Response: " + response.getResponseCode() + " " + response.getContentText());
  } catch (err) {
    Logger.log("Webhook forwarding failed: " + err.toString());
  }
}

function parseSubstationRollcall(formData) {
  // Parses substation duty operators from Form sections
  return [
    {
      substation: "Bukuru 33kV Injection",
      operator_name: formData["Bukuru Substation Duty Operator"] || "Unspecified",
      status: formData["Bukuru Substation Operational State"] || "NORMAL"
    },
    {
      substation: "Jos Central 33kV",
      operator_name: formData["Jos Central Duty Operator"] || "Unspecified",
      status: formData["Jos Central Operational State"] || "NORMAL"
    }
  ];
}

function parseFeederReadings(formData) {
  // Extracts feeder line telemetry for priority dispatch ledger
  const feeders = [];
  const defaultFeeders = [
    { code: "FDR-IND-01", name: "Industrial Line 1", rank: 1, band: "A" },
    { code: "FDR-COM-01", name: "Commercial Hub 1", rank: 2, band: "B" },
    { code: "FDR-RES-01", name: "Residential Line 1", rank: 4, band: "D" }
  ];

  for (let i = 0; i < defaultFeeders.length; i++) {
    const f = defaultFeeders[i];
    const loadMW = parseFloat(formData[f.name + " Load (MW)"] || "3.5");
    const currentAmps = parseFloat(formData[f.name + " Current (Amps)"] || "180.0");
    const state = formData[f.name + " State (IN_SERVICE/OUT_ON_FAULT/OUT_ON_LOAD_SHEDDING/OUT_ON_PERMIT_PTW)"] || "IN_SERVICE";

    feeders.push({
      feeder_asset_id: "asset-" + f.code.toLowerCase(),
      feeder_code: f.code,
      current_load_mw: loadMW,
      current_amps: currentAmps,
      state: state,
      priority_rank: f.rank,
      tariff_band: f.band
    });
  }

  return feeders;
}

