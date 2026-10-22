/**
 * Google Apps Script Webhook Forwarder: Google Form A (Generation Master Station)
 * 
 * Bounded Context: S05 Generation Department & Master Station Hub
 * Deployment: Install on Google Form onFormSubmit trigger in Google Apps Script editor.
 * Purpose: Signs submitted form responses with HMAC-SHA256 and POSTs to Core API Gateway.
 */

const CORE_API_URL = "https://api.nesco.ng/api/v1/generation/webhook/google-forms";
const HMAC_SECRET = "nesco-master-hmac-integrity-key-32b!";

function onFormSubmit(e) {
  const itemResponses = e.response.getItemResponses();
  const formData = {};

  for (let i = 0; i < itemResponses.length; i++) {
    const itemResponse = itemResponses[i];
    formData[itemResponse.getItem().getTitle()] = itemResponse.getResponse();
  }

  // Construct Standardized JSON Payload for Plants 1 to 6
  const payload = {
    shift_id: formData["Shift Identifier"] || ("SHIFT-" + Utilities.formatDate(new Date(), "GMT", "yyyyMMdd-HH")),
    shift_window: formData["Shift Window (Morning/Afternoon/Night)"] || "MORNING",
    master_operator_id: formData["Master Operator Staff ID"],
    plants: parsePlantSections(formData)
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

function parsePlantSections(formData) {
  // Parses plant fields for Plants 1 through 6
  const plants = [];
  for (let p = 1; p <= 6; p++) {
    const prefix = "Plant " + p + " - ";
    if (formData[prefix + "Asset ID"]) {
      plants.push({
        plant_id: formData[prefix + "Asset ID"],
        plant_state: formData[prefix + "Operational State"] || "IN_SERVICE",
        water_level_m: parseFloat(formData[prefix + "Water Level (m)"] || "0.0"),
        mw_generated: parseFloat(formData[prefix + "Active Power (MW)"] || "0.0"),
        stator_temp_c: parseFloat(formData[prefix + "Stator Winding Temp (C)"] || "0.0"),
        vibration_mm_s: parseFloat(formData[prefix + "Bearing Vibration (mm/s)"] || "0.0"),
        oil_pressure_bar: parseFloat(formData[prefix + "Lube Oil Pressure (bar)"] || "0.0"),
        operator_names: (formData[prefix + "Operators on Duty"] || "").split(",").map(s => s.trim()).filter(Boolean),
        protection_engineers: (formData[prefix + "Protection Engineers"] || "").split(",").map(s => s.trim()).filter(Boolean),
        technicians: (formData[prefix + "Technicians"] || "").split(",").map(s => s.trim()).filter(Boolean),
        handover_notes: formData[prefix + "Handover Notes"] || "Shift shift operating normally."
      });
    }
  }
  return plants;
}

