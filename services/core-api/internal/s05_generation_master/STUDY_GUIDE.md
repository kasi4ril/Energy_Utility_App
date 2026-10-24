# Sprint 5 Study Guide: Generation Department & Master Station Hub (`s05_generation_master`)

## 1. Purpose & Utility Domain Context
The **Generation Department & Master Station Hub Bounded Context (`s05_generation_master`)** manages electricity generation across NESCO's six remote hydroelectric power stations (Plants 1 to 6) and centralizes operational command at the **Master Station**:
- **Remote Plant Call Collation:** Due to the physical remoteness and terrain constraints of the six generating stations, field operators at individual plants do not each run independent web servers or complex local platforms. Instead, the **Operator on Duty at the Master Station** contacts each plant via dedicated voice/radio lines at shift handover and hourly intervals.
- **Dual Intake Channels (Direct REST + Google Forms Webhook):** The Master Operator can log shift data directly via the Web Portal or through **Standardized Google Form A**. Responses submitted via Google Forms are automatically signed with an HMAC-SHA256 signature and forwarded to our `/api/v1/generation/webhook/google-forms` ingestion bridge.
- **Physical Hydrological & Generator Guard:** Telemetry readings are vetted against physical laws: reservoir water levels must be between 0.0m and 250.0m, generation MW cannot be negative or reported on offline plants, and early warning threshold alerts are triggered for high stator temperatures (> 110°C), excessive bearing vibration (> 4.5 mm/s), and low lube oil pressure (< 1.5 bar).
- **GM Approval & Workshop Fabrication Loop:** When a turbine runner, guide bearing, or generator coil suffers damage, the Unit Head assesses the fault and escalates to the **General Manager (Generation)**. Upon GM approval, the job moves to the **Central Workshop Team** for machining, coil rewinding, or welding, and is tracked until re-installation and testing.

---

## 2. File-by-File Architectural Breakdown

```text
services/core-api/internal/s05_generation_master/
├── domain/
│   ├── models.go       # Structs: MasterStationTelemetryLog, RemotePlantShiftRollcall, PlantFaultAndWorkshopJob
│   ├── rules.go        # Physical bounds, early warning alerts, GM approval transition, sequential workshop checks
│   └── rules_test.go   # Pure domain unit tests (90.4% coverage)
├── repository/
│   ├── repository.go           # Storage interface & MemoryRepository mock implementation
│   ├── postgres_repo.go        # PostgreSQL queries for Telemetry & Roll-call
│   ├── postgres_repo_jobs.go   # PostgreSQL queries for Plant Faults & Workshop jobs
│   └── repository_test.go      # Storage unit tests
├── service/
│   ├── workflow_service.go      # Orchestrator: Telemetry logging, HMAC Google Forms ingestion, GM & Workshop loops
│   └── workflow_service_test.go # Orchestration tests, webhook HMAC verification, fault lifecycles (82.8% coverage)
└── transport/
    ├── dto.go                 # DTOs: RecordTelemetryDTO, RecordRollcallDTO, ReportFaultJobDTO, GMApprovalDTO
    ├── dto_test.go            # Comprehensive DTO validation tests
    ├── http_handlers.go       # net/http REST endpoints, RBAC authentication, HMAC webhook handler
    └── http_handlers_test.go  # End-to-end HTTP integration tests (88.1% coverage)
```

### Reference Assets:
- `forms-specs/google_form_a_generation_master.js`: Google Apps Script reference implementation for Google Form A onFormSubmit webhook forwarder.
- `migrations/005_s05_generation_master.up.sql`: Database schema `generation_master` containing telemetry, rollcall, and workshop jobs.
- `migrations/005_s05_generation_master.down.sql`: Reversible database migration.

---

## 3. Step-by-Step Data & Execution Flow

### A. Remote Voice-Call Collation & Google Forms Ingestion
```
[ Generating Plants 1–6 ]
           │
           │  Voice / Radio Call Hourly Roll-Call & Log Sheet Data
           ▼
[ Master Station Operator on Duty ]
           │
           ├── Option A: Web Portal Direct REST Entry
           │             POST /api/v1/generation/telemetry
           │             POST /api/v1/generation/rollcall
           │
           └── Option B: Low-Bandwidth Standardized Google Form A
                         -> Apps Script computes HMAC-SHA256 signature
                         -> POST /api/v1/generation/webhook/google-forms
                         -> Verifies X-Signature-SHA256 header against HMAC_SECRET
                         -> Evaluates water level, MW, winding temperature, vibration
                         -> Persists telemetry & roll-call records
                         -> Emits 'generation.telemetry.recorded' and 'generation.telemetry.alert'
```

### B. Plant Fault & Workshop Fabrication Lifecycle
```
1. Operator Reports Fault:
   POST /api/v1/generation/faults -> Status starts at REPORTED (or UNIT_HEAD_REVIEWED).
   - If requires_workshop is true -> WorkshopStatus begins at PENDING_DISPATCH.

2. General Manager (GM) Review Gate:
   POST /api/v1/generation/faults/{id}/gm-approval (Role: GENERAL_MANAGER_GEN)
   - If Approved: Status -> GM_APPROVED. Authorizes store requisition & workshop labor.
   - If Rejected: Status -> CANCELLED (requires explanation notes).

3. Central Workshop Reconditioning:
   POST /api/v1/generation/faults/{id}/workshop
   - Transitions sequentially: PENDING_DISPATCH -> MACHINING_IN_PROGRESS / WINDING_FABRICATION
     -> QUALITY_CHECK -> COMPLETED_DELIVERED.
   - On delivery, the plant job automatically marks RESOLVED_REINSTALLED!
```

---

## 4. Concrete cURL Examples

### 1. Ingest Google Forms HMAC-Signed Shift Webhook
```bash
curl -X POST http://localhost:8080/api/v1/generation/webhook/google-forms \
  -H "X-Signature-SHA256: <computed-hmac-sha256-hex>" \
  -H "Content-Type: application/json" \
  -d '{
    "shift_id": "SHIFT-20261022-M",
    "shift_window": "MORNING",
    "master_operator_id": "00000000-0000-0000-0000-000000000001",
    "plants": [
      {
        "plant_id": "00000000-0000-0000-0000-000000000002",
        "plant_state": "IN_SERVICE",
        "water_level_m": 88.5,
        "mw_generated": 16.5,
        "stator_temp_c": 72.0,
        "vibration_mm_s": 1.2,
        "oil_pressure_bar": 3.2,
        "operator_names": ["Operator Musa", "Operator Yakubu"],
        "protection_engineers": ["Engr. Pam"],
        "handover_notes": "Turbine 1 running at optimal efficiency."
      }
    ]
  }'
```

### 2. Report Plant Mechanical Fault (Requires Workshop)
```bash
curl -X POST http://localhost:8080/api/v1/generation/faults \
  -H "Authorization: Bearer <TOKEN_OPERATOR>" \
  -H "Content-Type: application/json" \
  -d '{
    "plant_id": "00000000-0000-0000-0000-000000000002",
    "wing": "MECHANICAL",
    "fault_description": "Severe cavitation pitting on turbine runner bucket; dynamic unbalance detected.",
    "severity": "HIGH",
    "requires_workshop": true
  }'
```

### 3. General Manager (GM) Approval
```bash
curl -X POST http://localhost:8080/api/v1/generation/faults/<job-uuid>/gm-approval \
  -H "Authorization: Bearer <TOKEN_GM>" \
  -H "Content-Type: application/json" \
  -d '{
    "approved": true,
    "notes": "Approved for emergency stainless steel weld overlay and dynamic balancing in central workshop."
  }'
```

---

## 5. Security & Domain Invariants
1. **HMAC Webhook Cryptographic Barrier:** Google Forms webhook endpoints are unauthenticated via JWT (since Google Apps Script does not possess a staff user session), but are strictly guarded by HMAC-SHA256 signing using the shared master integrity secret. Unsigned or tampered payloads are immediately rejected with HTTP 401.
2. **Offline Generation Physics Violation:** An offline generating plant cannot generate active power. Submissions reporting active power on an out-of-service plant are rejected by domain validation.
3. **GM Gate Bypass Prevention:** Workshop machinists cannot begin work or charge store materials without explicit digital authorization from the General Manager.
4. **Immutable Shift Handover Notes:** Under NERC Grid Code regulations and Section 13 standards, once an operator on duty submits shift handover notes, the text is cryptographically locked with `is_locked: true` and cannot be edited, altered, or deleted by any user or administrator.

---

## 6. Web Portal Landing Page Architecture (Section 13 Implementation)

The Generation Master Station Portal (`/portal/generation-master`) employs a persistent 2-column layout (`GenerationNav_LeftSidebarPane.tsx`) featuring 4 dedicated operational endpoints:

1. **Activity 1: Real-Time Analytics (`GenerationAnalyticsView.tsx`, `/portal/generation-master/analytics`)**
   - Synthesizes live hourly readings from PostgreSQL `generation_master.master_station_telemetry_logs`.
   - Displays reservoir depths (m), headwater/tailwater elevations, and river inflow rates (m³/s).
   - Compares active MW output vs. installed capacity per station with percentage utilization meters.
   - Evaluates machine state of health (ISO 10816 bearing vibration, stator winding temperature, lube oil pressure).
   - Highlights out-of-service stations with root-cause categorization (`TURBINE_MECH_FAULT`, `HYDRO_LOW_WATER`, etc.).

2. **Activity 2: Hourly Reading Sheet Form (`ReadingSheetForm.tsx`, `/portal/generation-master/reading-sheet`)**
   - Standard Form S05-TEL for the operator on duty to commit readings straight to the database.
   - Captures active MW, reactive MVAr, frequency (50.00 Hz), water level, elevations, inflow, stator temp, bearing vibration (with instant > 4.5 mm/s trip alarm), lube oil pressure, and operator attribution.
   - Immediately commits to DB and updates real-time analytics with confirmation receipt.

3. **Activity 3: Change of Shift Desk (`ChangeOfShiftForm.tsx`, `/portal/generation-master/change-of-shift`)**
   - Captures concluded shift window, outgoing operator, incoming operator, and 6-plant physical staff muster.
   - Short column for handover notes strictly enforced as **immutable**: once submitted, notes are cryptographically signed and cannot be edited.
   - Sealed historical register with SHA-256 tamper-proof lock badges.

4. **Activity 4: Supervisor Requisition Section (`SupervisorRequisitionView.tsx`, `/portal/generation-master/requisition`)**
   - Direct digital channel for operators to write official notes and spare parts requisitions to immediate Unit Heads (`Engr. D. Gyang`, `Engr. S. Lar`), the General Manager (`Engr. M. B. Jauro`), Workshop Lead, and Central Stores.
   - Includes itemized Bill of Materials builder and live queue tracking supervisor acknowledgments and GM store releases.


