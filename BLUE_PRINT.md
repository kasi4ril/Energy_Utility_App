# MASTER ARCHITECTURAL & FUNCTIONAL SPECIFICATION DOCUMENT (MAFSD)

**Project Title:** Enterprise Utility Workflow, Grid Intelligence & Local AI Governance Platform  
**Target Domain:** Integrated Electricity Generation (6 Remote Plants) & Distribution Company (DisCo)  
**Architectural Pattern:** Event-Driven Modular Monolith (Microservices-Ready) with Air-Gapped Local LLM Security Sidecar  
**Document Purpose:** Definitive Source-of-Truth Implementation Blueprint for Autonomous Coding Agents & Engineering Teams  
**Date:** October 2026  

---

## SECTION 0: AGENT EXECUTION DIRECTIVES & GLOBAL ARCHITECTURE

### 0.1 Core Engineering Rules for the Coding Agent
1. **Strict Sequential MVP Execution:** Build the platform strictly in the **9-Sprint order** defined below. Never build a downstream workflow before its foundational asset registry, IAM role matrix, and cryptographic audit hooks exist.
2. **Zero Direct Database Mutation on Sensitive Records:** No API endpoint may perform a raw `UPDATE` on customer balances, tariff bands, meter serial numbers, store stock counts, or approved permits without passing through the **Maker-Checker State Machine** and appending a **SHA-256 Hash-Chain Audit Row**.
3. **100% On-Premise AI Data Sovereignty:** All AI inference, natural-language-to-SQL queries, anomaly detection, and `@GridAssistant` chat interactions must execute locally via **Ollama** inside the internal network. Zero PII, financial ledgers, or grid topologies may be transmitted to external cloud LLM APIs.
4. **Universal Inter-Departmental `WorkNotes` Engine:** Every transactional entity (Ticket, Work Order, Permit-to-Work, Material Requisition, Store Waybill, Billing Adjustment, Audit Claim, and Shift Handover) must attach to the unified, append-only `WorkNotes` table allowing cross-departmental `@mention` tagging and automated Ollama executive summarization.

### 0.2 Approved World-Standard Technology Stack
| Architectural Layer | Technology Specification | Purpose & Operational Justification |
| :--- | :--- | :--- |
| **Core Backend & API Gateway** | **Go (Golang 1.23+)** (`Chi` or `Fiber` + `gRPC`) | High-concurrency execution for real-time prepaid vending webhooks, field ticket routing, and Google Forms webhook ingestion with minimal memory overhead. |
| **AI, ML & Analytics Sidecar** | **Python 3.12+** (`FastAPI` + `Scikit-learn` + `PyTorch`) | Hosts predictive maintenance pipelines (water level, generator health, DT load), energy theft/churn models, and Ollama agent orchestration. |
| **Local LLM & RAG Engine** | **Ollama** (`Llama-3.2` / `Qwen-2.5` + `nomic-embed-text`) | Powers the **Data Integrity Guardian**, **Automated Report Formatter**, and **`@GridAssistant`** live store/workflow chat bot. |
| **Web Portal Frontend** | **Next.js (React + TypeScript)** + `Tailwind CSS` + `TanStack Query` | Role-scoped Departmental Landing Pages, Intelligence Dashboards, and the Shared Executive Command Center. |
| **Field & Station PWA** | **Offline-First Progressive Web App (PWA)** + `IndexedDB` | Allows Feeder/Metering Team Leads to capture GPS coordinates, photos, and work proofs offline in remote corridors and auto-sync when online. |
| **Unified Database Cluster** | **PostgreSQL 16+** (`PostGIS` + `TimescaleDB` + `pgvector`) | Schema-per-module isolation; `PostGIS` for spatial asset/GPS validation; `TimescaleDB` for hourly plant/water/feeder telemetry; `pgvector` for local RAG. |
| **Event Bus & Cache** | **NATS JetStream** + **Redis** | Asynchronous cross-department event triggers, real-time dashboard WebSocket pushes, and idempotency locks for vending/webhooks. |
| **Low-Bandwidth Field Ingestion** | **Standardized Google Forms** -> **HMAC-Signed Webhook Bridge** | Enables Master Station (Generation) and Central Dispatch (Energy Management) operators to log remote call-collated shift data reliably over low bandwidth. |

### 0.3 Global Monorepo Folder Blueprint
- `apps/web-portal/` — Next.js Departmental Portals & Executive Command Center
- `apps/field-pwa/` — Offline-first PWA for Team Leads & Safety Reps
- `services/core-api/cmd/server/main.go` — Go Backend Entrypoint
- `services/core-api/internal/s01_foundation/` — Org Hierarchy, Grid Assets, IAM, Federated HR, WorkNotes
- `services/core-api/internal/s02_ollama_guardian/` — Hash-Chain Ledger, Maker-Checker, Chat Router, Webhook Guard
- `services/core-api/internal/s03_customer_care/` — KYC Onboarding, Ticketing, Direct Routing, Closure Gate
- `services/core-api/internal/s04_services_maint/` — Feeder, Transformer & Metering Units, GPS Proof-of-Work
- `services/core-api/internal/s05_generation_master/` — Master Station Hub (6 Plants), Elec/Mech/Workshop, GM Gate
- `services/core-api/internal/s06_energy_management/` — Central Dispatch, 33kV/11kV Rationing, TCN Hook, 5-Step PTW
- `services/core-api/internal/s07_stores_audit/` — Store Ledger, @GridAssistant Stock Bot, Audit Gate, Claims, MD Gate
- `services/core-api/internal/s08_billing_commercial/` — Metered Billing, Meter Refund Ledger, Multi-Gateway Vending
- `services/core-api/internal/s09_executive_command/` — Shared C-Suite Aggregator & Predictive Maintenance Triggers
- `services/ai-guardian/` — Python FastAPI + Ollama Local LLM & Predictive ML Service (`agents/`, `models/`, `webhooks/`)
- `forms-specs/` — Google Forms Apps Script bridges & JSON schemas
- `migrations/` — Sequential SQL migrations (`001_foundation` to `009_executive`)
- `deploy/` — Docker Compose, Nginx TLS Gateway, Backup scripts

### 0.4 Strict Internal Folder & File Naming Standard (Predictable Symmetry Rule)
To guarantee maintainability, easy debugging, and effortless future upgrades, the Coding Agent MUST NOT create monolithic files or arbitrary file names. Every module across the Backend, Frontend, AI Sidecar, and Database must strictly obey this symmetrical file tree:

```text
disco-enterprise-platform/
│
├── 1. DATABASE MIGRATIONS (1-to-1 Sprint Mapping)
├── migrations/
│   ├── 001_s01_foundation.up.sql            # Org Units, Grid Assets, IAM, HR, WorkNotes
│   ├── 002_s02_ollama_guardian.up.sql       # Audit Hash-Chain, Maker-Checker, Chat Logs
│   ├── 003_s03_customer_care.up.sql         # Customers, Tickets, Meter Schedules
│   ├── 004_s04_services_maint.up.sql        # Work Orders, GPS & Photo Proof-of-Work
│   ├── 005_s05_generation_master.up.sql     # Master Station Telemetry, 6-Plant Rollcall, Workshop
│   ├── 006_s06_energy_management.up.sql     # Dispatch Logs, Feeder States, 5-Step Safety PTW
│   ├── 007_s07_stores_audit.up.sql          # Store Stock, Requisitions, Audit Consent, Staff Claims
│   ├── 008_s08_billing_commercial.up.sql    # Metered Billing, Refund Ledger, Multi-Gateway Vending
│   └── 009_s09_executive_command.up.sql     # Executive Views & Predictive Maintenance Alerts
│
├── 2. GO CORE BACKEND (Clean Hexagonal Architecture — Identical Inside Every s01..s09 Module)
├── services/core-api/
│   ├── cmd/server/main.go                   # Wires config, DB pool, NATS, and module routers
│   ├── pkg/                                 # Shared utilities (NEVER put department logic here)
│   │   ├── auth/jwt_middleware.go           # Login token verifier & RBAC role guard
│   │   ├── crypto/hmac_chain.go             # SHA-256 & HMAC signature helpers
│   │   ├── database/postgres.go             # PgBouncer/PostGIS connection pool
│   │   └── events/nats_bus.go               # Cross-department asynchronous event bus
│   └── internal/
│       └── sXX_[module_name]/               # <-- EVERY Sprint module (s01 to s09) has these exact 4 folders:
│           ├── domain/
│           │   ├── models.go                # Struct definitions & ENUM constants
│           │   └── rules.go                 # Pure business validation & state transition rules
│           ├── repository/
│           │   └── postgres_repo.go         # All SQL & PostGIS queries for this module ONLY
│           ├── service/
│           │   └── workflow_service.go      # Orchestrates domain rules, repo, NATS events & AI calls
│           └── transport/
│               ├── http_handlers.go         # REST API & Webhook endpoints for this module
│               └── dto.go                   # Request/Response JSON validation structs
│
├── 3. PYTHON + OLLAMA LOCAL AI SIDECAR (Modular Agent & ML Pipeline Layout)
├── services/ai-guardian/
│   ├── app/
│   │   ├── main.py                          # FastAPI entrypoint (:8001)
│   │   ├── core/ollama_client.py            # Local Ollama connection & prompt templates
│   │   ├── agents/                          # Local LLM Agents
│   │   │   ├── grid_assistant_chat.py       # @GridAssistant live Store Stock & workflow bot
│   │   │   ├── fraud_and_gps_watchdog.py    # Maker-Checker scorer, Anti-Ghost & GPS verifier
│   │   │   └── report_formatter.py          # Shift & Monthly executive PDF/Docx generator
│   │   └── models/                          # Predictive Maintenance & Revenue Protection ML
│   │       ├── hydro_water_forecast.py      # 7-30 day reservoir water & MW predictor
│   │       ├── generator_cbm_model.py       # 6-Plant vibration/temp trip predictor
│   │       ├── transformer_feeder_cbm.py    # DT overload & feeder fault predictor
│   │       └── energy_theft_and_churn.py    # DT energy balance, bypass & churn probability model
│
├── 4. NEXT.JS WEB PORTAL FRONTEND (Feature-Sliced 4-Zone Layout)
├── apps/web-portal/src/
│   ├── middleware.ts                        # Post-login Role-to-Landing-Page auto-router
│   ├── app/                                 # Next.js App Router (Thin route wrappers ONLY)
│   │   ├── login/page.tsx                   # Unified Login Page
│   │   ├── executive-command/page.tsx       # Shared C-Suite & Dept Head Command Center
│   │   └── portal/[department]/page.tsx     # Departmental Landing Page route
│   ├── components/shared/                   # Global Reusable UI Components
│   │   ├── ZoneA_TopCommandBar.tsx          # Global header, role badge, shift status
│   │   ├── ZoneD_GridAssistantDock.tsx      # Slide-out @GridAssistant & community chat
│   │   └── WorkNotesDrawer.tsx              # Universal inter-departmental notes panel
│   └── features/                            # <-- EVERY Department has its own isolated feature folder:
│       ├── customer-care/                   # (Same 3 files inside services, generation, energy-mgmt,
│       │   ├── ZoneB_DashboardRibbon.tsx    #  stores, internal-audit, and billing-commercial)
│       │   ├── ZoneC_ActionQueue.tsx        # Left 65% live ticket/approval table
│       │   ├── ZoneC_IntakeForm.tsx         # Right 35% data entry form
│       │   └── api.ts                       # Typed TanStack Query hooks calling Go Backend
│       └── executive-command/               # Shared MD & Department Heads Dashboard Panels
│           ├── GenerationWaterPanel.tsx
│           ├── FeederStatusMatrixPanel.tsx
│           ├── RevenueByPlatformPanel.tsx
│           ├── TheftAndAuditRadarPanel.tsx
│           └── PredictiveMaintPanel.tsx
│
└── 5. OFFLINE FIELD PWA & GOOGLE FORMS BRIDGES
    ├── apps/field-pwa/src/                  # Offline-first mobile UI for Team Leads & Safety Reps
    │   ├── db/indexeddb_sync.ts             # Offline queue for GPS + Photo proofs & PTWs
    │   └── views/FieldProofForm.tsx         # Mobile GPS capture & before/after photo uploader
    └── forms-specs/                         # Google Forms Webhook Bridges
        ├── form_a_generation_master.gs      # Apps Script HMAC webhook for 6-Plant Master Station
        └── form_b_energy_dispatch.gs        # Apps Script HMAC webhook for 33kV Substations/Feeders


---

## SPRINT 1: CORE FOUNDATION, GRID TOPOLOGY, FEDERATED HR & NOTES ENGINE (`s01_foundation`)

### 1.1 Organizational & Administrative Hierarchy
1. **Administrative Hierarchy:** `Headquarters (HQ)` -> `Districts` -> `Service Units` alongside the remote `Generation Master Station` -> `Generating Plants 1 to 6`.
2. **Electrical Grid Topology:** `Generating Plants (1–6) & Future TCN Bulk Point` -> `33kV Injection Substations` -> `11kV Feeders` -> `Distribution Transformers (DTs)` -> `Service Poles` -> `Customer Meters`.

### 1.2 Role-Based Access Control (RBAC) Matrix
- **Executive Command Tier (Shared C-Suite Dashboard Access):** Managing Director (MD), Head of Services, Head of Internal Audit, Head of Commercial Operations, Head of Generation (General Manager), Head of Energy Management.
- **Departmental & Unit Heads:** Customer Care Head, Feeder Maintenance Unit Head, Transformer Maintenance Unit Head, Metering Unit Head, Electrical Unit Head (Gen), Mechanical Unit Head (Gen), Workshop Lead (Gen), Billing Unit Head, Store Officer, Safety Representative (HSE), Plant HR Officer.
- **Operational & Field Execution Tier:** Master Station Operator (Generation), Central Dispatch Operator (Energy Mgmt), 33kV Substation Operator on Duty, Feeder/Metering Team Leads, Customer Care Officers, Internal Auditors.

### 1.3 Landing Page & Data Intake Specification
- **Grid Asset Registration Workspace:** Captures Asset Code, Asset Type (`PLANT`, `SUBSTATION_33KV`, `FEEDER_11KV`, `TRANSFORMER_DT`), Parent Asset ID, Capacity (MVA/kVA), Approved Tariff Band (`A`–`E`), and mandatory PostGIS GPS Coordinates (`lat/long`).
- **Federated HR & Shift Roster Workspace:** Allows the local **Generation HR Officer** and **HQ HR** to maintain staff profiles, skill certifications, and 24/7 shift rosters, syncing automatically between Generation and HQ via NATS JetStream.
- **Standardized Inter-Departmental `WorkNotes` Component:** Embedded on every module's screen. Captures `Author_ID`, `Department`, `Linked_Entity_Type`, `Linked_Entity_ID`, `Note_Body`, `@Tagged_Units`, and `Timestamp`.

### 1.4 Intelligence Dashboard Metrics
- Active Grid Topology summary (Total MW capacity, active 33kV Substations, 11kV Feeders, DTs, and customer saturation ratio per DT).
- Federated Workforce Roster compliance across HQ, Districts, Service Units, and Generation.

### 1.5 Logical Data Entities (`schema: foundation`)
- `org_units` (`id`, `parent_id`, `unit_code`, `name`, `unit_type`, `location_gps`)
- `staff_profiles` (`id`, `staff_number`, `full_name`, `role`, `org_unit_id`, `certification_status`, `active_shift_state`)
- `grid_assets` (`id`, `parent_asset_id`, `org_unit_id`, `asset_code`, `asset_type`, `rated_capacity_kva`, `tariff_band`, `gps_point`, `operational_state`)
- `work_notes` (`id`, `entity_type`, `entity_id`, `author_staff_id`, `tagged_unit_ids[]`, `note_text`, `ollama_summary`, `created_at`)

---

## SPRINT 2: OLLAMA LOCAL AI GUARDIAN, CRYPTOGRAPHIC LEDGER & `@GridAssistant` (`s02_ollama_guardian`)

### 2.1 Dual-Brain Local AI Architecture
1. **Brain 1 — The Security, Integrity & Anti-Fraud Watchdog:**
   - **Cryptographic Hash-Chain Ledger:** Every sensitive table row includes a `row_hmac_signature` and logs mutations to `immutable_audit_chain` where `current_hash = SHA256(prev_hash + payload + actor_id + timestamp)`. Any direct SQL tampering breaks the chain, freezes the record, and alerts the **Head of Internal Audit** and **MD**.
   - **Maker-Checker Gatekeeper:** Intercepts all high-risk changes (Tariff Band changes, Debt adjustments, Meter swaps, Store stock edits), computes an **AI Fraud Risk Score (`0.00–1.00`)**, and blocks execution until the designated Checker approves.
2. **Brain 2 — `@GridAssistant` Productivity & Community Co-Pilot:**
   - **Real-Time Store Stock Lookup Bot:** Staff tag `@GridAssistant` in the in-app community chat to check exact available quantities of conductors, fuses, meters, or spares in the Store Database before raising a request.
   - **Automated Report & Document Formatter:** Converts shift logs, fault summaries, and audit findings into standardized company reports (PDF/Docx/Markdown) with one click.
   - **Contextual Workflow Helper:** Summarizes multi-department `WorkNotes` threads and answers role-permitted operational queries in natural language.

### 2.2 Logical Data Entities (`schema: security_ai`)
- `immutable_audit_chain` (`id`, `table_name`, `record_id`, `action_type`, `old_json`, `new_json`, `maker_id`, `checker_id`, `prev_hash`, `current_hash`, `timestamp`)
- `maker_checker_requests` (`id`, `module`, `change_category`, `target_record_id`, `payload_json`, `justification`, `ollama_risk_score`, `ollama_reasoning`, `status`)
- `community_chat_messages` (`id`, `channel_id`, `sender_id`, `message_text`, `is_ai_response`, `referenced_entity_id`, `created_at`)

---

## SPRINT 3: CUSTOMER CARE, DIRECT ROUTING & MANDATORY CLOSURE GATE (`s03_customer_care`)

### 3.1 Core Operational Rules
- **Direct-to-Unit-Head Routing with Departmental Head Supervisory Radar:**
  - *Fallen Pole / Line Snap* -> **Feeder Maintenance Unit Head** *(Watched by Head of Services)*.
  - *Bad Transformer / Oil Leak / Fuse* -> **Transformer Maintenance Unit Head** *(Watched by Head of Services)*.
  - *Faulty Meter / New Meter Installation Schedule* -> **Metering Unit Head** *(Watched by Head of Services)*.
  - *Wrong Billing / Uncredited Payment* -> **Billing Unit Head** *(Watched by Head of Commercial Operations)*.
- **Exclusive Customer Care Closure Gate:** Technical units and billing units **cannot close a ticket**. They can only mark a ticket `ATTENDED_TO` with mandatory proof. The ticket returns to Customer Care, who must **call the reporting customer(s) to confirm resolution** before marking it `CLOSED` (or `REOPENED_ESCALATED` if the customer denies fix).

### 3.2 Landing Page & Data Intake Specification
1. **New Customer Key-In Form:** Captures Full Name, Phone, NIN/ID, Premises GPS (`lat/long`), Building Class, Distribution Transformer (DT) Code, and Meter Procurement Application reference.
2. **Omnichannel Complaint Intake Form:** Captures Source (`WALK_IN`, `ONLINE_PORTAL`, `CALL_CENTER`), Account/Meter Number, Category, Description, and Photo/Receipt upload.
3. **Meter Installation & Repair Scheduler:** Allocates customer installation/repair dates against Metering Unit daily capacity.
4. **Customer Confirmation & Closure Desk:** Queue of all `ATTENDED_TO` tickets showing customer phone number, technical proof uploaded by the unit, call outcome log (`CONFIRMED_BY_CUSTOMER` vs. `REJECTED_BY_CUSTOMER`), and final `CLOSE` button.

### 3.3 Intelligence Dashboard Metrics
- First-Contact Billing Resolution Rate, Unit-by-Unit Mean Time to Resolve (MTTR), Customer Confirmation Rejection Rate, and Pending Meter Installation Backlog.

### 3.4 Ollama AI Guardian Rules
- **Anti-Ghost Customer Lock:** Verifies new customer GPS falls within the assigned DT's service radius and checks that no prior disconnected debtor account exists at the same GPS coordinate or phone/ID hash.
- **Duplicate Outage Grouping:** Automatically links multiple customer "No Power" complaints under the same DT/Feeder into a single Master Fault Ticket while requiring Customer Care to confirm restoration with callers before final closure.

### 3.5 Logical Data Entities (`schema: customercare`)
- `customers` (`id`, `account_no`, `cin`, `full_name`, `encrypted_pii`, `transformer_id`, `meter_serial`, `tariff_band`, `ledger_balance`, `kyc_status`, `row_hmac_signature`)
- `service_tickets` (`id`, `ticket_no`, `customer_id`, `channel`, `category`, `assigned_unit_id`, `unit_head_id`, `dept_head_supervisor_id`, `status`, `customer_call_verified`, `verified_by_cc_officer_id`, `sla_due_at`)
- `meter_schedules` (`id`, `customer_id`, `ticket_id`, `schedule_type`, `scheduled_date`, `assigned_metering_team_id`, `status`)

---

## SPRINT 4: SERVICES & MAINTENANCE DEPARTMENT (`s04_services_maint`)

### 4.1 Sub-Units & Closed-Loop Field Workflow
Covers the **Feeder Maintenance Unit**, **Transformer Maintenance Unit**, and **Metering Unit** under the **Head of Services**:
1. **Unit Head Dispatch:** Unit Head receives routed ticket on their Landing Page, attaches technical instructions in `WorkNotes`, and assigns it to the specific **Feeder Team Lead** (or Metering/Transformer Team Lead).
2. **Safety & Store Hooks:** If line isolation is needed, triggers a **Permit-to-Work (PTW)** request to Energy Management (Sprint 6). If materials are needed, queries `@GridAssistant` and raises a **Material Requisition** (Sprint 7).
3. **Mandatory GPS & Photo Proof-of-Work:** When the field repair or meter work is completed, the **Team Lead** submits the **Field Proof Form** on the PWA:
   - Device-captured live **GPS Coordinates** at the exact repair spot.
   - **Before & After Timestamped Photos** (and new meter seal numbers for Metering).
   - Detailed technical work summary.
4. **Unit Head Verification & Reversion:** The Unit Head verifies the Team Lead's proof and clicks **"Mark Attended To & Revert to Customer Care"**.

### 4.2 Landing Page & Dashboard Specification
- **Landing Page:** Unit Head Dispatch Queue, Team Lead Mobile Proof-of-Work Upload Form, Preventive Patrol & DT Load Inspection Logger, and Unit Material Request initiator.
- **Dashboard Metrics:** Active Faults per Feeder/DT, Team Lead Turnaround Time, GPS Proof Mismatch Alerts, and Preventive vs. Corrective Maintenance Ratio.
- **Ollama AI Guardian Rule:** Compares the Team Lead's Proof-of-Work GPS coordinates against the customer meter or pole/DT GIS coordinates. Blocks submission if GPS distance exceeds tolerance or if an uploaded photo hash matches a previously used photo.

### 4.3 Logical Data Entities (`schema: services_maint`)
- `maintenance_work_orders` (`id`, `ticket_id`, `unit_type`, `unit_head_id`, `team_lead_id`, `asset_id`, `ptw_permit_id`, `status`)
- `proof_of_work_logs` (`id`, `work_order_id`, `team_lead_id`, `captured_gps`, `distance_from_asset_meters`, `before_photo_url`, `after_photo_url`, `work_details`, `ollama_gps_verified`, `submitted_at`)

---

## SPRINT 5: GENERATION DEPARTMENT & MASTER STATION HUB (`s05_generation_master`)

### 5.1 Organizational Structure & Remote Call Collation Architecture
- **Leadership:** **General Manager (Generation)** supported by a local **Plant HR Officer** (synced with HQ HR).
- **Technical Wings across the 6 Plants:**
  - *Electrical Wing:* Protection Engineers, System Operators, Electrical Maintenance Unit.
  - *Mechanical Wing:* Mechanical Maintenance Team, Workshop Team.
- **Centralized Master Station Login (Remote Plant Collation):** Due to the remoteness of the **6 Generating Plants**, only the **Operator on Duty at the Master Station** logs directly into the platform. At shift handover and hourly intervals, the Master Station Operator contacts all 6 plants via voice/radio call and submits the **Standardized Generation Master Station Google Form**.

### 5.2 Fault, GM Material Approval & Workshop Fabrication Loop
1. **Master Station Fault Intake:** When any of the 6 plants reports an electrical or mechanical fault via call, the Master Station Operator logs it and routes it to the **Electrical or Mechanical Unit Head**.
2. **GM Escalation & Store Requisition:** The Unit Head assesses the fault and reports to the **General Manager (GM)**. If spares/materials are required, the Unit Head submits a request to the **GM**, who approves and forwards it to **Stores & Internal Audit**.
3. **Workshop Integration:** If a part requires machining, winding, or fabrication, the Unit Head issues a **Workshop Job Order** to the **Workshop Team**, tracking raw materials and turnaround time until re-installation.

### 5.3 Standardized Google Form A: Generation Master Station Collation & Handover
- **Section 1 — Shift & Handover Header:** Timestamp, Shift Window (`Morning 08:00-16:00`, `Afternoon 16:00-00:00`, `Night 00:00-08:00`), Outgoing Master Operator ID, Incoming Master Operator ID.
- **Section 2 — Remote Plant Roll-Call (Plants 1 to 6 via Call):** Names & Staff IDs of System Operators, Protection Engineers, and Electrical/Mechanical technicians physically on duty at each plant.
- **Section 3 — Plant Operational State (Plants 1 to 6):** Status (`IN_SERVICE`, `STANDBY`, `OUT_OF_SERVICE`), Outage Reason Category (`HYDRO_LOW_WATER`, `TURBINE_MECH_FAULT`, `GENERATOR_ELEC_FAULT`, `PROTECTION_RELAY_TRIP`, `PLANNED_MAINTENANCE`), and Detailed Root-Cause Note.
- **Section 4 — Hydrological & Generator Health Telemetry:** Reservoir Water Level (m), Headwater/Tailwater elevation, Inflow/Spillway state, Active Power Generated (MW), Reactive Power (MVAr), Stator Winding Temp (°C), Bearing Temp/Vibration, Lube Oil Pressure.
- **Section 5 — Inter-Shift Handover Notes:** Active faults under GM review, pending Workshop jobs, and safety advisories.

### 5.4 Intelligence Dashboard Metrics
- Live & Daily MW Output across Plants 1–6, Water Level & Reservoir Head curve, Plant Availability Factor (%), Active Staff on Duty across all 6 stations, Workshop Fabrication backlog, and Early Fault Prediction alerts.

### 5.5 Logical Data Entities (`schema: generation_master`)
- `master_station_telemetry_logs` (`id`, `shift_id`, `plant_id`, `plant_state`, `outage_reason`, `water_level_m`, `inflow_rate`, `mw_generated`, `winding_temp_c`, `vibration_mm_s`, `oil_pressure_bar`, `logged_by_operator_id`, `recorded_at`)
- `remote_plant_shift_rollcall` (`id`, `shift_id`, `plant_id`, `on_duty_staff_ids[]`, `handover_notes`, `synced_to_hq_hr`)
- `plant_fault_and_workshop_jobs` (`id`, `plant_id`, `reported_by_operator_id`, `unit_head_id`, `gm_approved`, `requires_workshop`, `workshop_status`, `store_requisition_id`, `status`)

---

## SPRINT 6: ENERGY MANAGEMENT, PRIORITY DISPATCH & 5-STEP SAFETY PTW (`s06_energy_management`)

### 6.1 Dual-Source Bulk Intake & Priority Feeder Rationing
- **Power Sources:** Receives real-time available MW from the **6 Generating Plants (via Master Station)** plus a pre-built, toggleable **TCN Bulk Off-Take Interface** (ready for future National Grid integration).
- **Priority Feeder Rationing Engine:** Allocates available MW across 33kV Injection Substations and 11kV Feeders ranked by: (1) **Real-Time Revenue Inflow Efficiency** and (2) **Regulatory Supply Class / Tariff Band Commitment**.
- **Centralized Call Collation for Remote Substations:** The **Central Dispatch Operator** logs in at the Energy Management hub, collates remote **33kV Injection Substation** staff-on-duty and feeder states via call/radio, and submits the **Standardized Energy Management Google Form**.

### 6.2 The 5-Step Safety Permit-to-Work (PTW) & Substation Isolation Lock
1. **Step 1 (Request):** Feeder or Transformer Maintenance Unit Head raises a **PTW Request** on a specific Feeder/Transformer.
2. **Step 2 (Dual Approval):** Approved digitally by both the **Head of Services** AND the **Safety Representative (HSE)**.
3. **Step 3 (Substation Grant & Isolation):** Presented to the **Substation Staff on Duty**, who isolates the breaker, opens line isolators, applies earth switches, and clicks **"Grant Permit — Line Isolated & Earthed"**.
4. **Step 4 (Safe Execution):** Maintenance proceeds under the **Safety Representative's** watch.
5. **Step 5 (Safety Clearance & Restoration):** Work completes -> **Safety Representative** signs **"Crew Clear & Safe to Energize"** -> **Substation Staff** removes earth/isolation and restores power.

### 6.3 Standardized Google Form B: Energy Management Dispatch & Substation Handover
- **Section 1 — Shift & Substation Roll-Call:** Shift Window, Central Dispatch Operator ID, Remote 33kV Substation Staff on Duty (collated via call).
- **Section 2 — Bulk Energy Intake:** Total MW received from Generation Plants (1–6) + Total MW received from TCN Bulk Supply (`0.00` default until active).
- **Section 3 — Feeder State & Rationing Log (Per 33kV/11kV Feeder):** Feeder Code, Current Load (MW/Amps), Operational State (`IN_SERVICE`, `OUT_ON_FAULT`, `OUT_ON_PERMIT_PTW`, `OUT_ON_LOAD_SHEDDING`), Specific Fault/Tripping Relay Indication, Active PTW Reference, and Scheduled Restoration Time.
- **Section 4 — Dispatch Handover Notes:** Critical grid constraints, active safety locks, and priority switching instructions.

### 6.4 Logical Data Entities (`schema: energy_management`)
- `dispatch_telemetry_logs` (`id`, `shift_id`, `generation_mw_in`, `tcn_mw_in`, `substation_staff_rollcall_json`, `handover_notes`, `logged_at`)
- `feeder_status_ledger` (`id`, `feeder_asset_id`, `state`, `outage_reason_code`, `outage_narrative`, `load_mw`, `priority_rank`, `active_ptw_id`, `updated_at`)
- `safety_permits_ptw` (`id`, `permit_no`, `work_order_id`, `asset_id`, `unit_head_id`, `dept_head_approved`, `safety_rep_id`, `safety_rep_approved`, `substation_operator_id`, `isolation_granted_at`, `safety_clearance_signed_at`, `restored_at`, `status`)

---

## SPRINT 7: STORES, INTERNAL AUDIT GATE, STAFF CLAIMS & MD PROCUREMENT (`s07_stores_audit`)

### 7.1 Live `@GridAssistant` Stock Check & Store Officer Discipline
- Staff query `@GridAssistant` in chat to get instant live balances before initiating requests.
- The **Store Officer** maintains the **Standardized Store Stock Ledger Form** for every receipt, issuance, and scrap return so the database is perpetually accurate.

### 7.2 Dual-Path Store & MD Procurement Workflow (With Mandatory Audit Gate)
- **Path A (Item In Stock):** Unit Head initiates **Material Requisition Form** -> Approved by **Departmental Head / Generation GM** -> Routed to **Internal Audit** for compliance verification (checking active ticket & prior scrap return) -> **Store Officer** issues material and generates digital Waybill.
- **Path B (Item Out of Stock — Critical Supplier Procurement):** System alerts **Internal Audit** and the concerned **Departmental Head** -> Procurement/Store raises **Restocking Purchase Request** -> Validated by **HQ Departmental Head** -> **Internal Audit** verifies pricing/vendor compliance and grants **Audit Consent** -> Routed to the **Managing Director (MD)** for final **MD Procurement Approval**.

### 7.3 Internal Audit Department Landing Page & Staff Claims Verification
- **Audit Store & Procurement Queue:** Dedicated workspace to review, query, or grant **Audit Consent** on all store requisitions and MD restocking requests.
- **Staff Allowance & Claims Verification Desk (Overtime, Travel, Out-of-Station, Hazard):**
  - All staff financial claims pass through Internal Audit.
  - **Automated 3-Way Ollama Verification:** Cross-checks every claim against: (1) **Call-Collated Shift Roster / HR Attendance**, (2) **Safety PTW Crew Log**, and (3) **Field Team Lead's GPS Proof-of-Work timestamp**.
- **Reverse Logistics (Scrap Return) Tracker:** Flags any team that was issued new cables, meters, or transformer parts but has not logged the replaced/damaged scrap back into the Store Scrap Yard.

### 7.4 Logical Data Entities (`schema: stores_audit`)
- `store_inventory` (`id`, `item_code`, `item_name`, `category`, `quantity_available`, `reorder_threshold`, `unit_cost`, `bin_location`, `last_updated_by_store_officer_id`)
- `material_requisitions` (`id`, `req_no`, `origin_dept_id`, `linked_ticket_or_job_id`, `items_json`, `in_stock_flag`, `unit_head_id`, `dept_or_gm_approved`, `audit_officer_id`, `audit_consent_status`, `md_procurement_approved`, `scrap_returned_flag`, `status`)
- `staff_allowance_claims` (`id`, `staff_id`, `claim_type`, `amount`, `claimed_date_window`, `linked_ptw_or_ticket_id`, `roster_match_verified`, `gps_match_verified`, `ollama_audit_flag`, `auditor_id`, `status`)

---

## SPRINT 8: BILLING, MULTI-PLATFORM VENDING & REVENUE PROTECTION AI (`s08_billing_commercial`)

### 8.1 Zero-Estimation Metered Billing & Meter Procurement Refund Engine
- **100% Metered Billing:** Estimated billing is architecturally disabled. Postpaid bills are generated strictly from verified meter readings; Prepaid customers vend via STS tokens.
- **Customer Meter-Purchase Refund Ledger:** When a new customer procures a meter, the cost is logged in `meter_refund_ledgers`. Every subsequent prepaid vend or postpaid payment automatically credits the customer's regulatory refund installment until the meter cost is fully amortized.

### 8.2 Multi-Platform Vending & Payment Aggregator Gateway
- Standardized API/Webhook ingestion hub capturing payments across **all approved revenue channels**, tagging every transaction with its exact `payment_platform_source`:
  - **Online Web & Self-Service Vending Portals** (Company Web Portal, Third-Party Online Utility Gateways).
  - **Fintech & Agency Networks** (OPay, Moniepoint, PalmPay, Paga, POS Agents).
  - **Commercial Bank Apps & USSD** (GTBank, Zenith, FirstBank, Access, NIBSS).
  - **Postpaid Corporate/Direct Bank Remittances**.

### 8.3 Wrong Billing Resolution Loop (Head of Commercial Operations Gate)
1. **Customer Care** logs a **"Wrong Billing"** ticket and routes it directly to the **Billing Unit Head**.
2. **Billing Unit** investigates the root cause and prepares a **Billing Adjustment Proposal**.
3. **Head of Commercial Operations (Commercial Services)** reviews the audit trail and approves the adjustment.
4. Ticket reverts (`ATTENDED_TO`) to **Customer Care**, who calls the customer to confirm satisfaction and marks the ticket **`CLOSED`**.

### 8.4 Energy Theft, Customer Churn & Revenue Leak AI Models
1. **Energy Theft & Meter Bypass Detector:** Compares Distribution Transformer (DT) bulk energy delivered against the sum of all vended/billed kWh under that DT, and flags individual meters whose vending cadence drops >50% below their 6-month moving average while their feeder is energized.
2. **Customer Churn & Off-Grid Migration Predictor:** Computes a `0.00–1.00` Churn Probability Score for high-value commercial/residential customers showing declining purchase trends correlated with feeder outage frequency.
3. **Gateway Settlement & Ledger Leak Watchdog:** Reconciles every generated token and postpaid receipt against bank/fintech settlement webhooks, alerting **Internal Audit** immediately on any unremitted collection.

### 8.5 Logical Data Entities (`schema: billing_commercial`)
- `vending_and_payment_transactions` (`id`, `tx_reference`, `customer_id`, `feeder_id`, `dt_id`, `payment_channel_category`, `payment_platform_source`, `amount_paid`, `kwh_vended`, `meter_refund_credited`, `sts_token`, `settlement_verified`, `created_at`)
- `meter_refund_ledgers` (`id`, `customer_id`, `meter_serial`, `initial_meter_cost`, `total_refunded_to_date`, `remaining_refund_balance`, `status`)
- `billing_adjustments` (`id`, `ticket_id`, `customer_id`, `billing_officer_id`, `old_balance`, `adjusted_balance`, `root_cause_note`, `head_commercial_ops_approved`, `status`)
- `revenue_protection_ai_alerts` (`id`, `customer_or_dt_id`, `alert_type`, `confidence_score`, `evidence_summary`, `assigned_investigation_unit_id`, `status`)

---

## SPRINT 9: SHARED EXECUTIVE COMMAND CENTER & PREDICTIVE MAINTENANCE SUITE (`s09_executive_command`)

### 9.1 Shared Executive Command Dashboard (Role-Authorized Access)
Accessible to the **Managing Director (MD)**, **Head of Services**, **Head of Internal Audit**, **Head of Commercial Operations**, **Head of Generation (GM)**, and **Head of Energy Management**:

| Executive Command Module | Live Intelligence & Decision Metrics Displayed |
| :--- | :--- |
| **1. Generation & Hydro Command (6 Plants)** | Live MW generated per plant; **Plant State (`IN` vs. `OUT` + exact reason)**; Reservoir Water Levels; **Call-Collated Shift Staff on Duty** across Plants 1–6. |
| **2. Grid & Feeder Status Matrix** | Total 33kV/11kV Feeders **`IN SERVICE`** vs. **`OUT OF SERVICE`** categorized by exact cause: **Out on Fault**, **Out on Safety PTW**, or **Out on Load Rationing**; plus **Call-Collated Substation Staff on Duty**. |
| **3. Aggregated & Source-Segmented Revenue Pulse** | Total live revenue inflow (Prepaid + Postpaid) broken down by **Exact Payment Platform Source** (*Online Web Gateways, OPay, Moniepoint, GTBank, Other Banks, Postpaid Direct*) showing volume, % share, gateway uptime, and net settlement. |
| **4. Energy Theft, Churn & Audit Integrity Radar** | High-priority **Energy Theft/Bypass hotspots by DT**, **Customer Churn Risk alerts**, pending **Internal Audit Store/Allowance flags**, and unreturned scrap counts. |
| **5. Predictive Maintenance & Water Forecast Center** | **7–30 Day Hydro Water Availability Forecast**; AI failure-risk countdowns and automated **Preventive Maintenance triggers** for **Generators (Plants 1–6)**, **33kV/11kV Transformers**, and **Feeders**. |
| **6. Executive Action & Approval Console** | Role-specific one-click sign-off queue: **MD** (Audit-cleared Out-of-Stock Supplier Procurements), **Head Commercial Ops** (Billing Adjustments), **Head Audit** (Store & Staff Claims Consent), **Head Services / Safety** (PTW Approvals), and **GM Generation** (Plant Material Requisitions). |

### 9.2 The 4 Predictive Models (`services/ai-guardian/models/`)
1. **Hydrological Water Availability Model:** Forecasting reservoir head and available MWh capacity from Master Station inflow/discharge telemetry.
2. **Generator Preventive Maintenance Model:** Predicting mechanical/electrical trips across the 6 plants from bearing vibration, winding temperature, oil pressure, and cumulative run-hours.
3. **Transformer Condition-Based Maintenance (CBM) Model:** Predicting DT and power transformer overload/insulation failure from billing load aggregation, fuse-blow frequency, and thermal logs.
4. **Feeder Preventive Patrol Model:** Predicting line trips and pole failures from historical fault density, weather patterns, and days since last bush-clearing patrol.

---

## Sprint 10: FULL-STACK WIRING — LOGIN FLOW, FRONTEND ROUTES, API CONTRACTS & DB TOPOLOGY

### 10.1 Unified Authentication, Login Page & Role-Based Auto-Routing
1. **Single Unified Login Portal (`/login`):**
   - **UI Fields:** `Staff Number / Official Email`, `Password`, and `Shift Station Selector` *(shown dynamically if the user is a Master Station Operator, Central Dispatch Operator, or 33kV Substation Operator)*.
   - **Authentication Handshake (`POST /api/v1/auth/login`):**
     - Frontend submits credentials over TLS to the Go Core API.
     - Go Core API verifies `password_hash` and `row_hmac_signature` in `foundation.staff_profiles`.
     - Returns an **HTTP-Only Secure JWT Access Token (15-min expiry)** + **Refresh Token** containing claims: `{ "staff_id", "primary_role", "org_unit_id", "can_view_executive_dashboard", "assigned_asset_id" }`.
2. **Automatic Post-Login Landing Page Router (Next.js Middleware):**
   - Immediately upon login, Next.js inspects `primary_role` and routes the user to their exact workspace:
     - `MD`, `HEAD_SERVICES`, `HEAD_INTERNAL_AUDIT`, `HEAD_COMMERCIAL_OPS`, `GM_GENERATION`, `HEAD_ENERGY_MGMT` -> **`/executive-command`** *(with toggle tab to their specific Departmental Landing Page)*.
     - `CUSTOMER_CARE_OFFICER`, `HEAD_CUSTOMER_CARE` -> **`/portal/customer-care`**
     - `UNIT_HEAD_FEEDER_MAINT`, `UNIT_HEAD_TX_MAINT`, `UNIT_HEAD_METERING` -> **`/portal/services/[unit-type]`**
     - `TEAM_LEAD_FEEDER`, `TEAM_LEAD_TX`, `TEAM_LEAD_METERING`, `SAFETY_REP_HSE` -> **`/field-pwa/tasks`** *(Offline-First Mobile View)*.
     - `MASTER_STATION_OPERATOR`, `UNIT_HEAD_GEN_ELECTRICAL`, `UNIT_HEAD_GEN_MECHANICAL`, `WORKSHOP_LEAD`, `HR_OFFICER_PLANT` -> **`/portal/generation-master`**
     - `CENTRAL_DISPATCH_OPERATOR`, `SUBSTATION_OPERATOR` -> **`/portal/energy-management`**
     - `STORE_OFFICER` -> **`/portal/stores`**
     - `INTERNAL_AUDITOR` -> **`/portal/internal-audit`**
     - `UNIT_HEAD_BILLING` -> **`/portal/billing-commercial`**

---

### 10.2 Standardized Departmental Landing Page & Dashboard UI Layout
Every departmental portal in `apps/web-portal/` MUST follow a standardized **4-Zone Enterprise Shell** so staff experience a unified interface across the company:

1. **Zone A — Top Command Bar (Global):**
   - Displays Logged-in Staff Name, Role Badge, Active Org Unit / Feeder / Plant Scope, Shift Handover Status Indicator, and Notification Bell (for SLA breaches and `@mention` tags).
2. **Zone B — Upper Intelligence Dashboard Ribbon (4 Summary KPI Cards + 2 Live Charts):**
   - Renders real-time departmental metrics via WebSocket/SSE (`GET /api/v1/dashboards/[dept-slug]/metrics`).
   - Includes a **"Generate Shift/Monthly Report"** button powered by the local Ollama Report Formatter.
3. **Zone C — Main Operational Workspace (Split 65% Action Queue / 35% Data Intake Form):**
   - **Left Pane (65%):** Filterable, live-updating Action Queue (e.g., Incoming Routed Tickets, Pending PTW Permits, Pending Audit Clearances, or `ATTENDED_TO` Customer Confirmation Queue). Clicking any row opens the **Universal Inter-Departmental `WorkNotes` Drawer**.
   - **Right Pane (35%):** Context-aware Data Intake Form (e.g., New Customer Key-In, Fault Dispatch Form, Material Requisition Form, or Store Stock Ledger Form).
4. **Zone D — Collapsible Right-Hand `@GridAssistant` & Community Chat Dock:**
   - Persistent slide-out chat panel connected to `ws://core-api/api/v1/chat/stream`.
   - Allows staff to message other departments or tag `@GridAssistant` to query live Store Stock availability, summarize ticket notes, or check feeder status.

---

### 10.3 Master Backend REST API & Webhook Endpoint Contract
All frontend-to-backend communication routes through the Go API Gateway (`/api/v1/*`), while AI/ML tasks route internally from Go to the Python FastAPI/Ollama Sidecar (`http://ai-guardian:8001/internal/ai/*`).

| Sprint / Module | HTTP Method & Endpoint Path | Allowed Roles / Auth | Frontend Source -> Backend Action -> DB Target |
| :--- | :--- | :--- | :--- |
| **S01: Auth & Foundation** | `POST /api/v1/auth/login`<br>`GET /api/v1/foundation/assets`<br>`POST /api/v1/foundation/assets`<br>`POST /api/v1/notes/{entityType}/{entityId}` | Public (Login)<br>All Authenticated<br>Admin / Energy Mgmt<br>All Authenticated | Authenticates staff against `foundation.staff_profiles`.<br>Fetches PostGIS asset tree (`foundation.grid_assets`).<br>Appends SHA-256 chained note to `foundation.work_notes`. |
| **S02: Ollama & Security** | `POST /api/v1/maker-checker/request`<br>`POST /api/v1/maker-checker/{id}/approve`<br>`WS /api/v1/chat/stream` | Maker Roles<br>Checker (Unit/Dept Head)<br>All Authenticated | Proxies to Python `ai-guardian` for Ollama risk score; writes to `security_ai.maker_checker_requests` and `immutable_audit_chain`.<br>Handles `@GridAssistant` live store stock queries. |
| **S03: Customer Care** | `POST /api/v1/customers/onboard`<br>`POST /api/v1/tickets`<br>`POST /api/v1/tickets/{id}/confirm-close`<br>`POST /api/v1/meters/schedule` | Customer Care Officer<br>Customer Care / Portal<br>**Customer Care ONLY**<br>Customer Care Officer | Verifies PostGIS DT radius via Ollama; writes to `customercare.customers`.<br>Auto-routes ticket to Unit Head + Dept Head radar (`customercare.service_tickets`).<br>Logs customer phone confirmation & transitions `ATTENDED_TO` -> `CLOSED`. |
| **S04: Services & Maint** | `POST /api/v1/maintenance/orders/{id}/assign`<br>`POST /api/v1/maintenance/orders/{id}/proof`<br>`POST /api/v1/maintenance/orders/{id}/attend` | Unit Heads (Feeder/TX/Meter)<br>Field Team Leads (PWA)<br>Unit Heads | Assigns ticket to Feeder/Metering Team Lead (`services_maint.maintenance_work_orders`).<br>Uploads live GPS + Before/After photos; Ollama verifies GPS proximity.<br>Marks `ATTENDED_TO` and reverts ticket to Customer Care queue. |
| **S05: Generation Hub** | `POST /api/v1/webhooks/google-forms/gen-master`<br>`POST /api/v1/generation/faults`<br>`POST /api/v1/generation/faults/{id}/gm-approve`<br>`POST /api/v1/generation/hr-roster/sync` | HMAC Webhook Secret<br>Master Station Operator<br>GM Generation<br>Plant HR Officer | Ingests Google Form A (Plants 1–6 call roll-call, IN/OUT state, water level, MW) into `generation_master.master_station_telemetry_logs`.<br>Routes plant fault to Elec/Mech Unit Head & Workshop; GM approves Store request. |
| **S06: Energy Mgmt & PTW** | `POST /api/v1/webhooks/google-forms/dispatch`<br>`POST /api/v1/ptw/request`<br>`POST /api/v1/ptw/{id}/safety-approve`<br>`POST /api/v1/ptw/{id}/substation-isolate`<br>`POST /api/v1/ptw/{id}/restore` | HMAC Webhook Secret<br>Maint Unit Head<br>Head Services + Safety Rep<br>Substation Operator<br>Safety Rep + Substation Op | Ingests Google Form B (Substation roll-call, 6-Plant + TCN MW, Feeder states) into `energy_management.dispatch_telemetry_logs`.<br>Executes the 5-Step Safety Permit-to-Work digital lock on `safety_permits_ptw`. |
| **S07: Stores & Audit** | `GET /api/v1/stores/inventory/search`<br>`POST /api/v1/stores/requisitions`<br>`POST /api/v1/audit/requisitions/{id}/consent`<br>`POST /api/v1/audit/claims/{id}/verify`<br>`POST /api/v1/md/procurements/{id}/approve` | All (`@GridAssistant`)<br>Unit Heads / Store Officer<br>Internal Auditor<br>Internal Auditor<br>**Managing Director (MD)** | Queries live `stores_audit.store_inventory`.<br>Routes in-stock requests through Internal Audit to Store Officer.<br>Routes out-of-stock requests through Audit Consent to MD Approval.<br>Runs 3-way GPS + PTW + Shift Roster check on staff overtime/travel claims. |
| **S08: Billing & Vending** | `POST /api/v1/webhooks/vending/{platform}`<br>`POST /api/v1/billing/postpaid/readings`<br>`POST /api/v1/billing/adjustments`<br>`POST /api/v1/billing/adjustments/{id}/approve` | Mutual-TLS / HMAC Gateway<br>Billing / Metering Unit<br>Billing Unit Head<br>**Head Commercial Ops** | Ingests OPay, Moniepoint, GTBank, and Online Web vending; deducts meter-purchase refund installment (`billing_commercial.meter_refund_ledgers`).<br>Approves Wrong Billing adjustments and reverts to Customer Care to close. |
| **S09: Executive Command** | `GET /api/v1/executive/overview`<br>`GET /api/v1/executive/revenue-by-source`<br>`GET /api/v1/executive/predictive-models` | MD, Head Services, Head Audit, Head Commercial, GM Gen, Head Energy Mgmt | Streams aggregated MW, Feeder In/Out reasons, Call-Collated Shift Staff, Source-Segmented Revenue, Theft/Churn flags, and Water/Generator/TX/Feeder forecasts. |

---

### 10.4 How the Frontend, Backend, AI Sidecar & Database Connect (Connection Topology)
1. **Frontend to Go Core API:**
   - Next.js Web Portal (`apps/web-portal`) and Offline Field PWA (`apps/field-pwa`) communicate exclusively with the **Go Core API (`:8080`)** via REST JSON (`https://api.disco.local/api/v1/*`) and WebSockets (`wss://api.disco.local/api/v1/chat/stream`).
   - The Frontend **never** talks directly to PostgreSQL or Ollama.
2. **Google Forms & Payment Gateways to Go Core API:**
   - **Google Forms A (Generation) & B (Energy Management)** trigger an `onFormSubmit` Google Apps Script that signs the JSON payload with `HMAC-SHA256` and posts to `/api/v1/webhooks/google-forms/*`.
   - **External Vending Platforms (OPay, Moniepoint, GTBank, Online Portals)** post signed transaction payloads to `/api/v1/webhooks/vending/{platform}`.
3. **Go Core API to Python/Ollama Sidecar (`ai-guardian:8001`):**
   - Connected over an isolated internal Docker/Kubernetes bridge network via high-speed gRPC / Internal REST.
   - Go calls Python synchronously for **Maker-Checker Fraud Scoring**, **GPS Proof Validation**, and **`@GridAssistant` Chat Queries**, and asynchronously via **NATS JetStream** for **Report Generation**, **Energy Theft/Churn Scoring**, and **Predictive Maintenance (Water/Generators/Transformers/Feeders)**.
4. **Database Schema Isolation & Connection Pooling (`PgBouncer -> PostgreSQL 16`):**
   - Single unified PostgreSQL 16 cluster partitioned into **8 strict logical schemas** (`foundation`, `security_ai`, `customercare`, `services_maint`, `generation_master`, `energy_management`, `stores_audit`, `billing_commercial`).
   - Each Go internal module connects using a schema-scoped database role, while cross-module workflows communicate via **Foreign Keys on `foundation.grid_assets` / `foundation.staff_profiles`** and asynchronous **NATS JetStream events**.

