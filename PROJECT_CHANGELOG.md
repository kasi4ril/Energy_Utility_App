# Project Changelog & Architecture Audit Log

All notable changes to the Enterprise Utility Platform will be documented in this file.
Format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), adhering to Conventional Commits.

---

## [1.1.0] - 2026-10-23

### Added - Enterprise Observations & Blueprint Customer Care Harmonization
- **Complaints & Customer Growth Query Engine (Option A Executive Suite - `AnalyticsPeriodQueryCard.tsx`, `analyticsData.ts`):**
  - **Year Dropdown**: Hierarchical dropdown selector for recorded database years (`2026 Current Operating Year`, `2025 Audit`, `2024 Audit`, `2023 Centenary`).
  - **Month Dropdown**: Full Year Overview (`ALL`) or individual months (`01 - January` through `12 - December`).
  - **Day Dropdown**: Full Month Summary (`ALL`) or individual day reports (`Day 01` through `Day 31`).
  - **Removed Quarterly Filter**: Streamlined hierarchical Year -> Month -> Day drilldown.
  - **Donut Chart Gauge**: Recharts-powered interactive donut chart visualizing Tickets Resolved (Emerald) vs Active Backlog (Amber) with dynamic center percentage and NERC regulatory compliance badge.
  - **Dual Grouped Bar Histogram**: Recharts-powered inflow vs resolution pace histogram dynamically adapting across months (for year query), weeks (for month query), or hourly shift intervals (for daily query).
  - **Quick Shortcuts**: One-click preset pills (`⚡ Today 23 Oct 2026`, `📅 Oct 2026 Month`, `📊 2026 Full Year`).
  - **Core Metric Display**: Total Complaints Lodged, Tickets Resolved, Active Unresolved, and NIN-verified New Customers registered.
- **Customer Revenue, Vending & Meter Fleet Analytics (`BillingVendingMeterAnalyticsCard.tsx`, `analyticsData.ts`):**
  - **Interactive Meters Installed Query**: Added Year dropdown (`2023-2026`) and Month dropdown (`All Months`, `01-12`) to the Meters Installed side.
  - **Stacked Category Bar Chart**: Multi-color stacked bar chart visualizing deployment across Residential (R2 - Teal), MD Customers (Indigo), and Industrial (Amber) for any queried year/month.
  - Postpaid Paid vs Unpaid/Arrears breakdown with monetary totals.
  - Prepaid STS Vended vs Zero-Vend/Flagged accounts with token issuance tallies.
  - Total installed fleet tracking (14,850 meters) with dynamic period energized totals.
- **Distribution Services Cross-Departmental Visibility (`CustomerCareComplaintFeed.tsx`):**
  - Head of Distribution Services live intake feed embedded into Services & Maintenance portal (`/portal/services/[unit]`).
  - Categorized oversight into tickets routed to Feeder Maintenance, Transformer Maintenance, and Metering Operations.
  - One-click work order dispatch action button from incoming customer complaints.
- **Smart Complaint Auto-Population & Landmark Columns (`View3_SmartComplaintDesk.tsx`):**
  - Keying in or searching customer account or meter number instantly auto-populates Full Name, Phone (+234), Service Address, and Major Landmark.
  - Distinct Major Landmark column and input integrated into complaint intake and backend dispatch payload.
- **National NIN Identity Verification & Name Auto-Population (`View2_NewCustomerKYC.tsx`):**
  - 11-digit NIN input automatically triggers national identity lookup via NIMC API (`verifyNIN`).
  - Auto-populates legal verified full customer name and displays a green verified badge.
  - Dedicated "Major Landmark" input positioned before Premises Physical Address.
- **Unified Department Header & Central Daily Analytics (`ZoneA_ThinTopHeaderBar.tsx`, `TopCommandBar.tsx`, `central-analytics.tsx`):**
  - Official NESCO crest logo expanded into a prominent, commanding brand mark (`w-11 h-11` glowing emblem, `text-xl` font-mono, `Sovereign Utility 1923`) spanning the full height down on the far-left.
  - "📊 Central Daily Analytics" button shifted towards the right side alongside operational action controls.
  - Renamed "Landing Page" to "Back" button (`ArrowLeft` icon) and shifted towards the right side alongside Central Analytics, giving the prominent NESCO crest and brand identity full space on the left.
- **Static Top Header Bar & Fixed Left Navigation Pane (`customer-care.tsx`, `ZoneA_ThinTopHeaderBar.tsx`, `ZoneNav_LeftSidebarPane.tsx`, `TopCommandBar.tsx`):**
  - Pinned the Top Command Bar and Thin Top Header Bar as static/sticky headers (`sticky top-0 z-40 bg-slate-900/95 backdrop-blur-md shrink-0`).
  - Pinned the persistent left-hand navigation pane (`ZoneNav_LeftSidebarPane`) as a static sidebar (`w-full md:w-64 shrink-0 h-full md:sticky md:top-14 md:h-[calc(100vh-3.5rem)] overflow-y-auto`).
  - Configured the Customer Care layout shell (`h-screen overflow-hidden`) with a dedicated independently-scrollable main workspace (`flex-1 overflow-y-auto min-w-0`), guaranteeing that both the top header and the left navigation pane stand completely still while scrolling up and down.
- **Omniscient Multi-Database Co-Pilot (`grid_assistant_chat.py`, `GridAssistantDock.tsx`):**
  - Connected `@GridAssistant` across Stores stockpile (conductors, transformers, fuses, oil), Hydro generation (Kurra, Jekko MW outputs), Feeder outages, Revenue collections, and automated Daily Shift Handover report generation.
  - Added one-click quick intelligence prompt chips to the slide-out chat dock.
- **Transitional Visual Login Showcase (`LoginPage.tsx`):**
  - Dynamic transparent auto-advancing slideshow featuring NESCO historical hydro installations and field lineworkers at work.

---

## [1.0.0] - 2026-10-22

### Added - Sprint 10: Full-Stack Wiring, Enterprise UI Shell & Field PWA (`s10_fullstack_wiring`)
- **Frontend Web Portal (`apps/web-portal/`):**
  - Next.js 14 Web Portal configured with TypeScript, Tailwind CSS, and strict typing.
  - `src/types/`: Domain models and auth contracts mirroring Go backend JWT claims and role matrices (`auth.ts`, `domain.ts`, `index.ts`).
  - `src/lib/api-client.ts` & `src/lib/auth.ts`: Central typed HTTP fetcher with Bearer injection, JWT base64URL claims decoder, and deterministic role-to-portal resolution (`getLandingPageForRole`).
  - `src/middleware.ts`: Next.js Edge Middleware enforcing role auto-routing to prevent unauthorized cross-departmental access.
  - `src/features/auth/`: Unified login interface with dynamic physical Shift Station selector (`LoginPage.tsx`, `api.ts`).
  - `src/components/shell/`: 4-Zone Enterprise UI Shell components including Zone A Top Command Bar (`TopCommandBar.tsx`), Zone D Floating Conversational AI Dock (`GridAssistantDock.tsx`), and universal cryptographic WorkNotes drawer (`WorkNotesDrawer.tsx`).
  - `src/components/ui/`: Reusable Tailwind UI building blocks (`Card.tsx`, `Badge.tsx`, `Table.tsx`).
  - Feature Modules (`ZoneB_DashboardRibbon.tsx`, `ZoneC_ActionQueue.tsx`, `ZoneC_IntakeForm.tsx`, `api.ts`):
    - `customer-care`: KYC onboarding form, ticket action queue, SLA health indicators.
    - `services-maint`: Unit work order queue (Feeder, Transformer, Metering), proof-of-work review, new order intake.
    - `generation-master`: Plant 1–6 telemetry ribbon, plant intake & workshop dispatch forms.
    - `energy-management`: Grid intake metrics, 5-Step Safety PTW approval queue & request form.
    - `stores-audit`: Stock ledger ribbon, requisition review queue, scrap & claim intake forms.
    - `billing-commercial`: STS token vending form, refund amortization tracker, billing adjustment review queue.
    - `executive-command`: C-Suite live pulse, predictive maintenance CBM alerts, role-segregated sign-off console.
  - `src/pages/`: Next.js page routes for all operational departments (`login.tsx`, `executive-command.tsx`, `portal/customer-care.tsx`, `portal/services/[unit].tsx`, `portal/generation-master.tsx`, `portal/energy-management.tsx`, `portal/stores.tsx`, `portal/internal-audit.tsx`, `portal/billing-commercial.tsx`, `index.tsx`).
- **Field Technician PWA (`apps/field-pwa/`):**
  - Offline-first progressive web app with ServiceWorker caching (`public/sw.js`, `public/manifest.json`).
  - Hardware GPS location acquisition (`navigator.geolocation`) and client-side photo SHA-256 hash generation for tamper-proof field proof of work.
  - Offline task queue using IndexedDB/localStorage with automatic background drain on connectivity restoration (`src/lib/offline-storage.ts`, `src/pages/tasks.tsx`).
- **Deployment Topology & Docker Orchestration:**
  - `docker-compose.yml`: Multi-container production deployment topology running PostgreSQL 16 PostGIS, NATS JetStream, Python AI Guardian (:8001), Go Core API (:8080), and Next.js Web Portal (:3000).
  - `.env.example`: Master environment configuration template with air-gapped sovereign defaults.
  - Multi-stage Dockerfiles: `services/core-api/Dockerfile`, `services/ai-guardian/Dockerfile`, and `apps/web-portal/Dockerfile`.
- **End-to-End Cross-Sprint Verification:**
  - `services/core-api/test/integration/`: 5 automated end-to-end multi-sprint integration test suites (`setup_test.go`, `e2e_customer_services_test.go`, `e2e_generation_ptw_test.go`, `e2e_commercial_executive_test.go`) exercising the entire utility workflow lifecycle with 100% pass rate.
- **Documentation & Study Guides:**
  - `services/core-api/internal/s10_fullstack_wiring/STUDY_GUIDE.md`: Comprehensive living study guide detailing full-stack topology, 4-zone UI patterns, PWA offline mechanics, and container orchestration.
  - `MILESTONES.md`: Updated master sprint roadmap to reflect 100% completion across all 10 sprints.

---

## [0.9.0] - 2026-10-22

### Added - Sprint 9: Executive Command & Predictive Maintenance Engine (`s09_executive_command`)
- **Database Migrations:**
  - `migrations/009_s09_executive_command.up.sql`: Schema `executive_command` containing `executive_snapshots` (shift-window MW generation, TCN intake, system load, feeder health tallies, revenue inflow, theft/churn alerts, audit/procurement tallies), `predictive_maintenance_triggers` (hydro water availability, generator vibration/thermal trip, transformer CBM, feeder trip/storm risk, linked automated work orders), and `executive_action_items` (role-specific digital sign-off queue for MD, Head of Commercial Ops, Head of Internal Audit, Head of Services, GM Generation) with btree indexes and uniqueness constraints.
  - `migrations/009_s09_executive_command.down.sql`: Safe reversible rollback dropping tables and schema cleanly.
- **Go Backend Core (`internal/s09_executive_command/`):**
  - `domain/models.go`, `rules.go`, `rules_test.go`, `rules_predictive_test.go`: Hydro water balance & MWh prognosis, ISO 10816 generator mechanical health evaluation, transformer thermal/overload CBM evaluation, feeder vegetation & storm patrol risk calculation, and role segregation enforcement rules (**92.6% statement coverage**).
  - `repository/repository.go`, `memory_repo_actions.go`, `postgres_repo.go`, `postgres_repo_actions.go`, `repository_test.go`, `repository_actions_test.go`: Hexagonal persistence contract, complete PostgreSQL storage engine, and thread-safe in-memory test mocks strictly $\le 321$ lines (**100.0% in-memory repo test coverage**).
  - `service/workflow_service.go`, `workflow_service_predictive.go`, `workflow_service_actions.go`, `workflow_service_test.go`, `workflow_service_predictive_test.go`: Business workflow orchestrators for C-Suite live telemetry aggregation, operational snapshots, 4 condition-based maintenance evaluations with automated dispatch into `s04repo.Repository`, and role-segregated action consoles (**87.4% statement coverage**).
  - `transport/dto.go`, `dto_test.go`, `http_handlers.go`, `http_handlers_predictive.go`, `http_handlers_test.go`, `http_handlers_predictive_test.go`: REST endpoints for dashboard pulse, snapshots, digital sign-off console, and predictive maintenance prognostics with JWT authentication and segregation of duties (**83.6% statement coverage**).
  - `cmd/server/main.go`: Gateway entrypoint updated to mount S09 routes, wire cross-departmental repositories, and update `/healthz` to `s09_executive_command`.
- **Documentation & Specs:**
  - `services/core-api/internal/s09_executive_command/STUDY_GUIDE.md`: Comprehensive living study guide detailing C-Suite intelligence architecture, physics equations for predictive maintenance models, segregation of duties matrix, and API contracts.
  - `docs/openapi.yaml`: OpenAPI 3.0.3 specification updated with 12 new Executive Command and Predictive Maintenance endpoints.
  - `MILESTONES.md`: Milestone roadmap and detailed execution ledger updated.

---

## [0.8.0] - 2026-10-22

### Added - Sprint 8: Billing, Commercial Operations & Revenue Protection AI (`s08_billing_commercial`)
- **Database Migrations:**
  - `migrations/008_s08_billing_commercial.up.sql`: Schema `billing_commercial` containing `vending_and_payment_transactions` (multi-platform STS token generation, tariff rate per kWh, meter refunds credited, settlement verification), `meter_refund_ledgers` (customer meter amortization, percentage per vend, remaining balance), `billing_adjustments` (wrong billing root causes, old vs adjusted balance, Head of Commercial Ops approval), and `revenue_protection_ai_alerts` (theft bypass, churn risk, confidence score, evidence summary) with btree indexes.
  - `migrations/008_s08_billing_commercial.down.sql`: Safe reversible rollback dropping tables and schema cleanly.
- **Go Backend Core (`internal/s08_billing_commercial/`):**
  - `domain/models.go`, `rules.go`, `rules_test.go`, `rules_adjustments_test.go`: Zero-estimation metered billing calculations, 20-digit STS token generation (`XXXX-XXXX-XXXX-XXXX-XXXX`), customer meter refund amortization rules, wrong billing adjustment validations, and revenue protection AI anomaly algorithms (theft bypass and customer churn risk) (**100.0% statement coverage**).
  - `repository/repository.go`, `memory_repo_adjustments.go`, `postgres_repo.go`, `postgres_repo_adjustments.go`, `repository_test.go`, `repository_adjustments_test.go`: Hexagonal persistence contract, full PostgreSQL storage implementation, and thread-safe in-memory test mocks strictly $\le 332$ lines (**100.0% memory repo test coverage**).
  - `service/workflow_service.go`, `workflow_service_adjustments.go`, `workflow_service_test.go`, `workflow_service_adjustments_test.go`: Workflow orchestrators for metered vending, refund ledger amortization, wrong billing disputes with automatic Customer Care ticket reversion to `ATTENDED_TO`, and local AI revenue protection telemetry scans (**89.3% statement coverage**).
  - `transport/dto.go`, `dto_test.go`, `http_handlers.go`, `http_handlers_adjustments.go`, `http_handlers_test.go`, `http_handlers_adjustments_test.go`: REST endpoints for metered vending, STS tokens, refund ledgers, wrong billing adjustments, and revenue protection AI anomaly scans with JWT security (**81.4% statement coverage**).
  - `cmd/server/main.go`: Gateway entrypoint updated to mount S08 routes and updated `/healthz` to `s08_billing_commercial`.
- **Documentation & Specs:**
  - `services/core-api/internal/s08_billing_commercial/STUDY_GUIDE.md`: Comprehensive living study guide detailing zero-estimation billing invariants, NERC meter refund amortization mechanics, the closed-loop wrong billing resolution gate, and AI revenue protection models.
  - `docs/openapi.yaml`: OpenAPI 3.0.3 specification updated with 11 new Billing, Commercial Operations, and Revenue Protection AI endpoints.
  - `MILESTONES.md`: Milestone roadmap and detailed execution ledger updated.

---

## [0.7.0] - 2026-10-22

### Added - Sprint 7: Stores, Internal Audit Gate & MD Procurement (`s07_stores_audit`)
- **Database Migrations:**
  - `migrations/007_s07_stores_audit.up.sql`: Schema `stores_audit` containing `store_inventory` (SKU master data, categories, balances, reorder thresholds, row HMACs), `material_requisitions` (dual-path in-stock vs out-of-stock MD purchase, unit/dept approvals, audit consent, digital waybills), `scrap_returns` (reverse logistics scrap yard tracker), and `staff_allowance_claims` (overtime, travel, hazard claims with 3-way verification and fraud scores) with btree indexes.
  - `migrations/007_s07_stores_audit.down.sql`: Safe reversible rollback dropping tables and schema cleanly.
- **Go Backend Core (`internal/s07_stores_audit/`):**
  - `domain/models.go` & `rules.go`: Store item sanity rules, reorder threshold evaluator, line item cost calculation, dual-path procurement workflow state guards, reverse logistics scrap return rules, and 3-way allowance claims verification formula (**84.6% test coverage**).
  - `repository/repository.go`, `postgres_repo.go`, `postgres_repo_scan.go`, `postgres_repo_scrap.go`, `postgres_repo_claims.go`: Decoupled persistence contract and PostgreSQL storage engine split into modular files strictly $\le 306$ lines.
  - `service/workflow_service.go`, `workflow_service_claims.go`: Business orchestrator managing warehouse stock adjustments, live `@GridAssistant` balance inquiries, dual-path material requisitions, store waybills, scrap returns, and the automated 3-way claims desk (**83.9% test coverage**).
  - `transport/dto.go`, `dto_test.go`, `http_handlers.go`, `http_handlers_scrap.go`, `http_handlers_claims.go`, `http_handlers_test.go`, `http_handlers_requisitions_test.go`, `http_handlers_claims_test.go`: REST endpoints for inventory SKUs, stock adjustments, dual-path requisitions, MD approval, material waybill dispatch, scrap returns, and staff allowance claims (**80.1% test coverage**).
  - `cmd/server/main.go`: Gateway entrypoint updated to mount S07 routes and updated `/healthz`.
- **Documentation & Specs:**
  - `services/core-api/internal/s07_stores_audit/STUDY_GUIDE.md`: Comprehensive living study guide detailing `@GridAssistant` inventory integration, dual-path procurement, reverse logistics, and the 3-way claims verification desk.
  - `docs/openapi.yaml`: OpenAPI 3.0.3 specification updated with 13 new Stores and Internal Audit endpoints.
  - `MILESTONES.md`: Milestone roadmap and detailed execution ledger updated.

---

## [0.6.0] - 2026-10-22

### Added - Sprint 6: Energy Management & Transmission Dispatch (`s06_energy_management`)
- **Database Migrations:**
  - `migrations/006_s06_energy_management.up.sql`: Schema `energy_management` containing `dispatch_telemetry_logs` (dual-source bulk intake, Generation Plants 1-6 + TCN, roll-call and handover notes), `feeder_statuses` (33kV and 11kV lines, real-time load MW, amperage, outage reason, priority ranks 1-5, NERC tariff bands A-E), and `safety_permits_ptw` (5-step high-voltage life cycle, dual approval timestamps, substation isolation proof, and safety clearance declarations) with btree indexes.
  - `migrations/006_s06_energy_management.down.sql`: Safe reversible rollback dropping tables and schema cleanly.
- **Go Backend Core (`internal/s06_energy_management/`):**
  - `domain/models.go` & `rules.go`: Dual-source bulk intake calculation ($\ge 0.0\text{ MW}$), priority feeder rationing algorithm (Bands A-E, preserving PTW/fault lines from energization), and 5-step Safety PTW state machine sequence guards (**91.2% test coverage**).
  - `repository/repository.go`, `postgres_repo.go`, `postgres_repo_ptw.go`: PostgreSQL persistence contract and thread-safe MemoryRepository split into clean modular files strictly $\le 350$ lines.
  - `service/workflow_service.go`, `workflow_service_ptw.go`: Orchestrator managing hourly bulk intake logs, low-bandwidth Google Forms B HMAC-SHA256 webhook ingestion, priority load shedding execution, and the 5-step Safety Permit-to-Work protocol (**88.4% test coverage**).
  - `transport/dto.go`, `dto_test.go`, `http_handlers.go`, `http_handlers_ptw.go`, `http_handlers_test.go`, `http_handlers_ptw_test.go`: REST endpoints for intake telemetry, priority load shedding calculations, feeder status telemetry, Google Forms webhook receiver, and PTW request/dual-approval/isolation/clearance/restoration (**80.2% test coverage**).
  - `cmd/server/main.go`: Gateway entrypoint updated to mount S06 routes and updated `/healthz`.
- **Forms & Integration Specs:**
  - `forms-specs/google_form_b_energy_dispatch.js`: Production-ready Google Apps Script webhook bridge computing HMAC-SHA256 signatures for Google Form B onFormSubmit trigger.
- **Documentation & Specs:**
  - `services/core-api/internal/s06_energy_management/STUDY_GUIDE.md`: Comprehensive living study guide detailing dual-source power balance, priority rationing mechanics, the 5-step PTW safety lifecycle, and API references.
  - `docs/openapi.yaml`: OpenAPI 3.0.3 specification updated with 13 new Energy Management and PTW endpoints.
  - `MILESTONES.md`: Milestone roadmap and detailed execution ledger updated.

---

## [0.5.0] - 2026-10-22

### Added - Sprint 5: Generation Department & Master Station Hub (`s05_generation_master`)
- **Database Migrations:**
  - `migrations/005_s05_generation_master.up.sql`: Schema `generation_master` containing `master_station_telemetry_logs` (hydrological water levels, active/reactive MW, stator temperature, vibration, oil pressure), `remote_plant_shift_rollcall` (operator, protection engineer, and tech rosters), and `plant_fault_and_workshop_jobs` (GM approval, workshop reconditioning states, store requisition links) with btree indexes.
  - `migrations/005_s05_generation_master.down.sql`: Safe reversible rollback dropping tables and schema cleanly.
- **Go Backend Core (`internal/s05_generation_master/`):**
  - `domain/models.go` & `rules.go`: Physical hydrological and generator bounds, early predictive maintenance warning alerts (stator temp > 110°C, vibration > 4.5 mm/s, low oil pressure < 1.5 bar), GM authorization gates, and sequential workshop lifecycle validation (**90.4% test coverage**).
  - `repository/repository.go`, `postgres_repo.go`, `postgres_repo_jobs.go`: Decoupled persistence contract, PostgreSQL storage engine, and thread-safe in-memory mock repository split across files strictly <= 350 lines.
  - `service/workflow_service.go`: Orchestrator handling telemetry logging, remote call roll-call roster capture, low-bandwidth Google Forms webhook ingestion with cryptographic HMAC-SHA256 signature verification, GM authorization, and central workshop reconditioning (**82.8% test coverage**).
  - `transport/dto.go`, `dto_test.go`, `http_handlers.go`, `http_handlers_test.go`: REST endpoints for telemetry, roll-calls, faults, GM approvals, workshop updates, and the Google Forms webhook receiver (**88.1% test coverage**).
  - `cmd/server/main.go`: Gateway entrypoint updated to register S05 routes and updated `/healthz`.
- **Forms & Integration Specs:**
  - `forms-specs/google_form_a_generation_master.js`: Production-ready Google Apps Script webhook bridge computing HMAC-SHA256 signatures for Google Form A onFormSubmit trigger.
- **Documentation & Specs:**
  - `services/core-api/internal/s05_generation_master/STUDY_GUIDE.md`: Comprehensive living study guide covering call collation mechanics, Google Forms HMAC ingestion, and cURL examples.
  - `docs/openapi.yaml`: OpenAPI 3.0.3 specification updated with 10 new Generation endpoints.
  - `MILESTONES.md`: Milestone roadmap and detailed execution ledger updated.

---

## [0.4.0] - 2026-10-22

### Added - Sprint 4: Services & Maintenance Department (`s04_services_maint`)
- **Database Migrations:**
  - `migrations/004_s04_services_maint.up.sql`: Schema `services_maint` containing `maintenance_work_orders` (linking tickets, units, priorities, assets, PTW safety permits, and store requisition refs) and `proof_of_work_logs` (spatial PostGIS `captured_gps`, distance from asset, SHA-256 before/after photo hashes, meter seals, and AI verification scores) with spatial and btree indices.
  - `migrations/004_s04_services_maint.down.sql`: Safe reversible rollback dropping tables and schema cleanly.
- **Go Backend Core (`internal/s04_services_maint/`):**
  - `domain/models.go` & `rules.go`: Unit-specific Haversine GPS tolerances (Feeder: 150m, Transformer: 75m, Meter: 50m), anti-photo-reuse gate (beforeHash != afterHash), PTW permit enforcement for isolated jobs, and Unit Head proof verification transitions (**100.0% test coverage**).
  - `repository/repository.go`, `postgres_repo.go`, `postgres_repo_proofs.go`: PostgreSQL + PostGIS queries and thread-safe in-memory repository split to maintain strictly <= 350 lines per file.
  - `service/workflow_service.go`: Orchestrator managing work order dispatch, mobile PWA field proof ingestion with air-gapped Python AI watchdog bridge and graceful fallback, Unit Head proof review, and automated reversion of linked customer tickets to `ATTENDED_TO` (**80.0% test coverage**).
  - `transport/dto.go`, `http_handlers.go`: REST endpoints for dispatching work orders, listing/getting orders, submitting PWA field proofs, and Unit Head verification (**82.5% test coverage**).
  - `cmd/server/main.go`: Gateway entrypoint updated to register S04 routes and updated `/healthz`.
- **Documentation & Specs:**
  - `services/core-api/internal/s04_services_maint/STUDY_GUIDE.md`: Comprehensive living study guide detailing architectural roles, closed-loop mechanics, and cURL examples.
  - `docs/openapi.yaml`: OpenAPI 3.0.3 specification updated with 5 new endpoints under the `ServicesMaintenance` tag.
  - `MILESTONES.md`: Milestone roadmap and detailed execution ledger updated.

---

## [0.3.0] - 2026-10-22

### Added - Sprint 3: Customer Care, Anti-Ghosting & Ticketing Engine (`s03_customer_care`)
- **Database Migrations:**
  - `migrations/003_s03_customer_care.up.sql`: Schema `customercare` containing `customers` (with PostGIS `premises_gps` point and FK to `grid_assets`), `service_tickets` (with SLA timestamps, direct routing unit heads, attended notes, and customer verification flags), and `meter_schedules` (activity types, scheduled dates, and team assignments), complete with spatial and btree indexes.
  - `migrations/003_s03_customer_care.down.sql`: Safe reversible rollback dropping tables and schema cleanly.
- **Go Backend Core (`internal/s03_customer_care/`):**
  - `domain/models.go` & `rules.go`: Spherical Haversine GPS proximity check (`<= 500m` of assigned DT), Anti-Ghost debtor prevention, deterministic Direct SLA Routing mapping ticket category to Unit Head & Dept Head, and strict two-step closure gate (**86.1% test coverage**).
  - `repository/repository.go`, `postgres_repo.go`, `postgres_repo_tickets.go`: PostgreSQL + PostGIS queries and thread-safe in-memory repository for rapid testing and edge fallback.
  - `service/workflow_service.go`: Business orchestrator verifying customer DT proximity, generating unique account and CIN numbers, enforcing technician proof of work, conducting customer closure gate transitions, scheduling meter activities, and publishing asynchronous events (**81.0% test coverage**).
  - `transport/dto.go`, `http_handlers.go`: REST endpoints for customer onboarding, listing, ticket logging, attending, two-step customer confirmation closure, and meter activity scheduling with RBAC enforcement (**84.5% test coverage**).
  - `cmd/server/main.go`: Gateway entrypoint updated to register S03 customer care, ticketing, and metering routes.
- **Documentation & Specs:**
  - `services/core-api/internal/s03_customer_care/STUDY_GUIDE.md`: Comprehensive living study guide covering architectural breakdown, execution flows, and cURL examples.
  - `docs/openapi.yaml`: OpenAPI 3.0.3 specification updated with 8 new endpoints across CustomerCare, Tickets, and Metering tags.
  - `MILESTONES.md`: Updated milestone status and execution log.

---

## [0.2.0] - 2026-10-22

### Added - Sprint 2: Ollama Local AI Guardian, Cryptographic Ledger & @GridAssistant (`s02_ollama_guardian`)
- **Database Migrations:**
  - `migrations/002_s02_ollama_guardian.up.sql`: Schema `security_ai` containing `immutable_audit_chain`, `maker_checker_requests`, `community_chat_messages`, indexes, and database trigger `trg_immutable_audit_chain` forbidding raw `UPDATE` or `DELETE` on the audit ledger.
  - `migrations/002_s02_ollama_guardian.down.sql`: Safe reversible rollback dropping triggers, tables, and schema in dependency order.
- **Python Local AI Sidecar (`services/ai-guardian/`):**
  - `app/core/ollama_client.py`: Air-gapped on-premise connection to local Ollama with deterministic heuristic fallback when offline.
  - `app/agents/fraud_and_gps_watchdog.py`: AI Fraud Risk Scorer (`0.00–1.00`) evaluating tariff band changes, debt adjustments, meter swaps, store stock edits, and justification quality, plus spherical Haversine GPS proximity verification.
  - `app/agents/grid_assistant_chat.py`: In-app conversational co-pilot for warehouse stock lookups, WorkNotes summaries, and feeder statuses.
  - `app/agents/report_formatter.py`: Automated shift handover Markdown report compiler.
  - `app/main.py`: Internal HTTP REST service (:8001) for risk scoring, GPS checks, chat queries, and health inspection.
  - `tests/test_s02_ollama_guardian.py`: 11 automated unit and boundary tests passing 100%.
- **Go Backend Core (`internal/s02_ollama_guardian/`):**
  - `domain/models.go` & `rules.go`: Maker-Checker segregation of duties, risk threshold escalation (`>= 0.80` demanding Executive tier), SHA-256 prev-hash audit chain calculation and validation (**90.9% test coverage**).
  - `repository/repository.go` & `postgres_repo.go`: Storage contract for audit chains, maker-checker requests, and chat records with thread-safe in-memory mock.
  - `service/workflow_service.go`: Orchestrator connecting to AI sidecar, executing maker-checker state transitions, appending SHA-256 chained audit entries, and managing `@GridAssistant` chat interactions (**89.4% test coverage**).
  - `transport/dto.go` & `http_handlers.go`: REST endpoints for request submission, approval, rejection, audit verification, and chat (**81.0% test coverage**).
  - `cmd/server/main.go`: Gateway entrypoint updated to register S02 endpoints alongside S01.
- **Documentation & Specs:**
  - `services/core-api/internal/s02_ollama_guardian/STUDY_GUIDE.md`: Comprehensive study guide.
  - `docs/openapi.yaml`: OpenAPI spec expanded with all Maker-Checker, Audit Chain, and Chat endpoints.
  - `MILESTONES.md`: Milestone progress ledger updated.

---

## [0.1.0] - 2026-10-22

### Added - Sprint 1: Core Foundation & Grid Topology (`s01_foundation`)
- Database migrations (`001_s01_foundation.up.sql`, `001_s01_foundation.down.sql`) with PostGIS geometry and append-only WorkNotes trigger.
- Foundational packages: `pkg/crypto`, `pkg/auth`, `pkg/events`, `pkg/database`.
- Hexagonal module for org units, grid assets, federated staff, and chained WorkNotes.
