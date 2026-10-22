# Sprint 10 Study Guide: S10 Full-Stack Wiring, Enterprise UI & Field PWA

## 1. Architectural Scope & Bounded Context
The **Full-Stack Wiring, Enterprise UI Shell & Field PWA** (`s10_fullstack_wiring`) represents the culminating integration tier of the NESCO Enterprise Utility Platform. It connects the 9 sovereign backend bounded contexts (S01 through S09), the local Python AI Guardian sidecar (:8001), and the PostGIS event store into a unified, air-gapped web portal, field technician progressive web app (PWA), and containerized deployment topology.

```
+---------------------------------------------------------------------------------------------------------+
|                                    NESCO Enterprise Full-Stack Topology                                |
|                                                                                                         |
|  [Browser Clients / Desktop Portals]           [Mobile Field Crews (Offline-First)]                     |
|           apps/web-portal (Port 3000)                       apps/field-pwa                              |
|   ├── Next.js Edge Middleware Router               ├── Service Worker & App Manifest Cache              |
|   ├── 4-Zone Enterprise Operating Shell            ├── LocalStorage / IndexedDB Offline Action Queue    |
|   └── Universal WorkNotes Drawer & @GridAssistant  └── Hardware GPS Geolocation & Photo Proof Hash      |
|                        │                                                  │                             |
|                        └─────────────────────────┬────────────────────────┘                             |
|                                                  ▼                                                      |
|                               [API Gateway / Core Backend (:8080)]                                      |
|                                   services/core-api (Go Monolith)                                       |
|   ├── S01 Foundation & RBAC IAM                    ├── S06 Energy Management & 5-Step PTW               |
|   ├── S02 Ollama Guardian & Ledger                 ├── S07 Stores & Internal Audit Desk                 |
|   ├── S03 Customer Care & Anti-Ghost               ├── S08 Billing, Vending & Revenue AI                |
|   ├── S04 Services & Maintenance Portals           └── S09 Executive Command & CBM Prognostics          |
|   └── S05 Generation Master Station Hub                                                                 |
|                        │                                                  │                             |
|                        ├─────────────────────────────────┐                │                             |
|                        ▼                                 ▼                ▼                             |
|             [PostgreSQL 16 + PostGIS]          [Local AI Guardian]   [NATS JetStream]                   |
|                Database Layer (:5432)             Python Sidecar       Event Bus (:4222)                |
|               (Spatial & Relational Data)              (:8001)                                          |
+---------------------------------------------------------------------------------------------------------+
```

---

## 2. Unified Login, Shift Stations & Edge Auto-Router

### 2.1 Unified Login & Station Selection
Staff members authenticate via a single entry point (`/login`) supplying enterprise credentials (`email` and `password`). Operators stationed at physical control centers (Master Station, Central Dispatch, Area Substations) select an active physical `Shift Station` from the dynamic selector.
Upon successful authentication, the backend issues an RS256/HS256 signed JWT bearing the operator's `primary_role`, `staff_id`, `full_name`, and `org_unit_id`.

### 2.2 Role-to-Portal Matrix & Next.js Edge Middleware
The Next.js Edge Middleware (`apps/web-portal/src/middleware.ts`) intercepts incoming requests and deterministically evaluates role-to-portal routing invariants:

| Staff Role / Category | Resolved Landing Route | Primary Responsibility |
| :--- | :--- | :--- |
| `MD`, `BOARD_DIRECTOR` | `/executive-command` | C-Suite Live Pulse, Sign-Off Queue, Board Telemetry |
| `CUSTOMER_CARE_REP`, `CALL_CENTER_AGENT` | `/portal/customer-care` | KYC Customer Onboarding, Ticket Triaging, 2-Step Closure |
| `UNIT_HEAD_FEEDER_MAINT`, `TEAM_LEAD_FEEDER` | `/portal/services/feeder` | Feeder Outage Triage, PWA Proof Reviews, Patrol Logs |
| `UNIT_HEAD_TX_MAINT`, `TEAM_LEAD_TX` | `/portal/services/transformer` | Distribution Transformer Overhauls, Oil Dielectric Checks |
| `UNIT_HEAD_METERING`, `METER_INSPECTOR` | `/portal/services/metering` | Meter Installations, Disconnections, STS Recalibrations |
| `MASTER_STATION_OPERATOR`, `GM_GENERATION` | `/portal/generation-master` | Hydro Plant Telemetry (Plants 1–6), Workshop Reconditioning |
| `SYSTEM_DISPATCHER`, `SUBSTATION_OPERATOR` | `/portal/energy-management` | Grid Balancing, Feeder Rationing, 5-Step PTW Lock |
| `STORE_OFFICER`, `PURCHASING_OFFICER` | `/portal/stores` | Material Ledgers, Waybills, Scrap Return Intakes |
| `HEAD_INTERNAL_AUDIT`, `AUDIT_OFFICER` | `/portal/internal-audit` | 3-Way Staff Claims Verification, Stock Audit Ledger |
| `HEAD_COMMERCIAL_OPS`, `BILLING_OFFICER` | `/portal/billing-commercial` | Zero-Estimation Vending, STS Tokens, Wrong Billing Desk |

---

## 3. The 4-Zone Enterprise UI Shell

Every operational portal implements the unified 4-zone UI structure:

```
+----------------------------------------------------------------------------------------------------+
| ZONE A: Top Command Bar                                                                            |
| [NESCO Logo] | [Active Shift Station] | [Role Switcher] | [WorkNotes Drawer (SHA-256)] | [User Profile] |
+----------------------------------------------------------------------------------------------------+
| ZONE B: Dynamic KPI Dashboard Ribbon                                                               |
| [Telemetry Metric 1]  |  [Telemetry Metric 2]  |  [Telemetry Metric 3]  |  [Critical Action Alert] |
+----------------------------------------------------------------------------------------------------+
| ZONE C: Split Operational Work Area                                                                |
|                                                                                                    |
|  [65% Left Pane: Role Action Queue]               │  [35% Right Pane: Intake / Execution Form]    |
|  - Real-time work items and triage list           │  - Structured task ingestion form              |
|  - Priority sorting & status badge indicators     │  - Form validation matching backend invariants |
|  - Action buttons (Review, Dispatch, Sign-off)    │  - Real-time feedback and submission status    |
|                                                   │                                                |
+----------------------------------------------------------------------------------------------------+
| ZONE D: @GridAssistant Conversational AI Dock (Bottom-Right Floating Pane)                         |
| - On-premise air-gapped LLM copilot                                                               |
| - Query inventory stock balances without leaving the task screen                                  |
| - Summarize long WorkNotes comment threads into bulleted handover briefs                          |
+----------------------------------------------------------------------------------------------------+
```

---

## 4. Field Technician PWA (Offline-First Architecture)

Located in `apps/field-pwa/`, the technician application is built for mobile crews operating in remote, low-bandwidth distribution zones.
- **Service Worker Cache (`public/sw.js`):** Intercepts network calls and serves static assets and recent task manifests offline.
- **Offline Storage (`src/lib/offline-storage.ts`):** Serializes pending proof-of-work submissions into IndexedDB/localStorage when connectivity is unavailable.
- **Hardware GPS Verification:** Leverages `navigator.geolocation` to capture high-precision coordinates (`latitude`, `longitude`) right at the physical distribution asset.
- **Client-Side SHA-256 Photo Hashing:** Computes cryptographic hashes of before-and-after maintenance photos prior to payload transmission, guaranteeing image integrity.
- **Auto-Sync:** Background event listeners monitor `window.addEventListener('online')` and automatically drain the offline task queue into the API gateway once cellular coverage is restored.

---

## 5. Deployment Topology & Docker Orchestration

The platform deploys with zero external SaaS dependencies using `docker-compose.yml`:
1. `nesco-db`: PostgreSQL 16 with PostGIS extension (`postgis/postgis:16-3.4-alpine`) on port `5432`.
2. `nesco-nats`: NATS JetStream messaging bus (`nats:latest`) on port `4222`.
3. `nesco-ai-guardian`: Air-gapped Python AI Guardian on port `8001` with health checks on `/health`.
4. `nesco-core-api`: Go 1.22 multi-stage binary on port `8080` with `/healthz` monitoring.
5. `nesco-web-portal`: Next.js 14 multi-stage production server on port `3000`.

---

## 6. End-to-End Cross-Sprint Verification

The complete platform is validated via five end-to-end integration suites (`services/core-api/test/integration/`):
1. **`TestE2E_Auth_And_CustomerCare_To_Services_Lifecycle`:** Verifies S01 staff login $\to$ S03 KYC onboarding $\to$ customer ticketing $\to$ S04 maintenance dispatch $\to$ field proof verification $\to$ 2-step phone confirmation gate.
2. **`TestE2E_GoogleFormsWebhooks_And_SafetyPTW_Lifecycle`:** Verifies S05 Form A HMAC-SHA256 generation webhook $\to$ S06 Form B dispatch webhook $\to$ 5-Step Safety PTW execution sequence.
3. **`TestE2E_MeteredVending_And_RefundAmortization_Lifecycle`:** Verifies S08 STS token generation $\to$ vending ledger $\to$ NERC meter refund amortization schedule.
4. **`TestE2E_ExecutiveCommand_And_PredictiveMaintenance_To_S04`:** Verifies S09 live telemetry aggregation $\to$ condition-based predictive maintenance $\to$ autonomous preventive work order generation in S04.
5. **`TestE2E_StoreStockLedger_And_DualPathRequisitions`:** Verifies S07 inventory ledger $\to$ Dual-Path Requisitions (Path A in-stock vs Path B MD procurement) $\to$ Reverse logistics scrap returns.

