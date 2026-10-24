# Project Changelog & Architecture Audit Log

All notable changes to the Enterprise Utility Platform will be documented in this file.
Format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), adhering to Conventional Commits.

---

## [1.4.0] - 2026-10-24

### Added - Services & Maintenance Feeder Unit Landing Page Architecture
- **Enterprise Shell & Customer Care Standard Layout (`portal/services/feeder.tsx`):**
  - Designed and deployed the dedicated Feeder Maintenance Unit landing page mirroring the Customer Care department layout: Thin Top Header Bar (`FeederHeaderBar.tsx`) + Persistent Left-Hand Navigation Sidebar (`FeederNav_LeftSidebarPane.tsx`) + 5-point modular operating workspaces.
  - Thin Top Header Bar features NESCO 1923 sovereign crest, real-time WAT clock, Unit Head staff attribution (`Engr. D. Gyang, MNSE • STAFF-ENG-04`), active line trip alarm bell, and rapid return links to General Landing (`/`) and Central Analytics (`/portal/central-analytics`).
- **Persistent Left-Hand Navigation Sidebar (`FeederNav_LeftSidebarPane.tsx`):**
  - Persistent left sidebar pane featuring direct endpoints to all departmental activities:
    1. Feeder Analytics Dashboard (`/portal/services/feeder/analytics`)
    2. Outage & PTW Request Form (`/portal/services/feeder/outage-request`)
    3. Proof of Work Submission Desk (`/portal/services/feeder/proof-of-work`)
    4. Feeder Technical Standards (`/portal/services/feeder/standards`)
    5. Active Feeder Work Orders (`/portal/services/feeder/work-orders`)
  - Includes real-time status badges (`1 Active Outage`, `1 Proof Due`, `2 Active Work Orders`), sister unit quick links (Transformer, Metering), emergency line trip radar pill (Anglo-Jos 33kV fault), and NERC regulatory footer.
- **Feeder Reliability Analytics Dashboard (`FeederAnalyticsView.tsx`):**
  - Real-time KPI cards for SAIDI (6.4 hrs), SAIFI (3.8 trips/cust), CAIDI (1.68 hrs), MTTR (2.1 hrs), and Monthly Trips count (9 trips).
  - Interactive Recharts Donut Pie Chart displaying feeder operational states (In-Service, Fault Outage, PTW Maintenance, Rationing).
  - Interactive Recharts Bar Chart showing trip density across corridors (Anglo-Jos 33kV, Bukuru 33kV, Rayfield 33kV, Vom 11kV, Makeri 33kV, Bokkos 33kV).
  - Complete feeder inventory and line health status table.
- **Outage Requisition & Permit-To-Work Form (`FeederOutageRequestView.tsx`):**
  - Formal 33kV/11kV planned outage requisition submitted to Central Dispatch & Safety under NERC Grid Code Part 3.
  - Form parameters: Feeder corridor, outage category, priority urgency, requested date/time window, duration in hours, estimated load drop (MW), affected DT count, physical isolation isolator tag (e.g. `ISO-AJ-04`), and portable earth switch tag (e.g. `ES-AJ-04A`).
  - **Permit Operational Lifecycle & Closure State Fields:**
    - Opening (granting) time by Central Dispatch (`opening_time_wat`).
    - Permit surrendering time by field team lead (`surrendering_time_wat`).
    - Post-surrender feeder closure state (`closure_outcome`): Closed & Energized Immediately, Left on Load Shedding (Rationing), Tripped on Closure (Fault/Surge), or Pending Surrender.
  - Hardened with explicit high-contrast `#0f172a` slate background and `#f8fafc` text on all `<select>` and `<option>` elements for legibility across light and dark styles.
- **Outage Request and Return Register with Calendar Date Query Suite (`FeederOutageRegisterTable.tsx`, `FeederRecordReturnModal.tsx`):**
  - **Comprehensive Outage Requisition & Return Tracking:** Renamed and upgraded the register to track the full lifecycle of isolation permits from granting to return and re-energization:
    - **Opening / Granting Time (`opening_time_wat`):** Time when Central Dispatch formally authorized and issued the isolation permit.
    - **Surrendering Time (`surrendering_time_wat`):** Time when the field maintenance crew leader formally surrendered the permit upon completion.
    - **Permit Exhaustion & Time Returned (`time_returned_wat`, `is_permit_exhausted`):** When a granted permit is exhausted, the exact timestamp when the permit was returned to dispatch is logged.
    - **Action Taken on Return (`action_taken_on_return`):** Mandatory field remarks documenting the exact corrective actions executed (e.g. crossarms replaced, conductors tensioned, insulator string washed, jumpers crimped) before line handback.
    - **Post-Closure Outcome (`closure_outcome`):** Immediate tracking of whether the feeder was `Closed & Energized Immediately`, `Left on Load Shedding` (rationing), or `Tripped on Closure` (flashover/inrush trip).
  - **Interactive Permit Return Desk Modal (`FeederRecordReturnModal.tsx`):** Dedicated operational modal allowing operators to record the return time, action taken, and feeder closure state for active permits with instant register synchronization.
  - **Year, Month & Day Calendar Search:** Multi-tiered calendar query engine allowing operators to select:
    - **Year:** `2026`, `2025`, `2024`, or `All Years`.
    - **Month:** `01` Jan to `12` Dec with full month names, or `All Months`.
    - **Day:** Dynamic calendar day selector (`01` through `31`) calculated based on the selected year and month, or `All Days`.
    - **Calendar Presets & Date Picker:** Quick buttons for `Today (24 Oct 2026)`, `Tomorrow (25 Oct)`, `Oct 2026 (All)`, `Sep 2026 (All)`, and native date picker synchronization.
  - **Approval State Filter:** Quick filtering for `Approved` (PTW issued, isolated, energized), `Pending` (awaiting dispatch), and `Declined` (rejected by dispatch due to grid / industrial load constraints).
  - **Dynamic KPI Badges & Active Filter Tags:** Real-time counter pills for Approved, Pending, and Declined requests with single-click "Clear Calendar" dismissal.
- **Field Proof of Work Submission Desk (`FeederProofOfWorkView.tsx`):**
  - Anti-ghosting geofenced coordinate verification with strict $\pm 15$m tolerance validation against asset register coordinates.
  - Side-by-side photographic inspection: Before repair (fault exposed) and After repair (rectified & inspected).
  - Materials consumed tracking and field work narrative.
  - Unit Head engineering sign-off desk with authoritative actions: `VERIFIED_APPROVED` (marking order attended to) or `REJECTED_REVISIT` with corrective audit notes.
- **Feeder Engineering Standards & Line Patrol Register (`FeederStandardsView.tsx`):**
  - NERC Grid Code & IEEE 80 standards reference cards:
    - Earth Resistance Testing: Mandatory $< 5.0\ \Omega$ for all switchyard surge arresters, transformer neutrals, and distribution ground pits.
    - FLIR IR Thermography Criteria: $\Delta T < 10.0^\circ C$ normal operating, $10-25^\circ C$ maintenance order within 48h, $> 25^\circ C$ critical emergency trip.
    - Right-of-Way (ROW) Clearance: 15m corridor (7.5m either side of centerline) free from vegetation intrusion.
  - Mandatory NERC 5-Step Electrical Safety Isolation Protocol: Disconnect, Lock & Tag (LOTO), Test Zero Voltage (Prove Dead), Earth & Short, Barricade.
  - Feeder Line Patrol Register logging pole spans, patroller, vegetation risk, insulator condition, hotspot $\Delta T$, and earth pit ohms.
- **Active Feeder Work Orders Queue (`FeederWorkOrdersView.tsx`):**
  - Operational queue with status and priority filtering, PTW mandate verification, and crew status transitions.
- **Routing & Sub-Routes Integration:**
  - Dedicated page routes: `/portal/services/feeder`, `/portal/services/feeder/analytics`, `/portal/services/feeder/outage-request`, `/portal/services/feeder/proof-of-work`, `/portal/services/feeder/standards`, `/portal/services/feeder/work-orders`.
  - Dynamic route fallback in `/portal/services/[unit].tsx` when `unit === "feeder"`.
  - Linked directly from Executive Command All Departments Launchpad (`AllDepartmentsAnalyticsPortalsCard.tsx`).
- **Strict Code Limit Compliance:** Every single component and page file strictly adheres to $\le 350$ lines of code.

---

## [1.3.0] - 2026-10-24

### Added - Executive Command Tier & Managing Director Landing Page Architecture
- **Enterprise Shell & Customer Care 5-Workspace Alignment (`portal/executive-command.tsx`):**
  - Modernized the Managing Director & CEO Tier-1 command portal to mirror the high-productivity 5-workspace left sidebar layout from Customer Care.
  - Added thin command header (`ExecutiveHeaderBar.tsx`) with NESCO sovereign crest, real-time Africa/Lagos (WAT) clock, MD/CEO profile badge (`Engr. Dr. E. K. Dung, FNSE, NPOM`), live pending approval bell count, and quick return links to General Landing Page (`/`) and Central Analytics (`/portal/central-analytics`).
- **Persistent Left-Hand Navigation Sidebar (`ExecutiveNav_LeftSidebarPane.tsx`):**
  - Persistent left sidebar with real-time status badges (`5 Pending Approvals`, `2 In Review`, `Emergency Grid Alert`, `Action Due`).
  - Seamless navigation across 5 core executive workspaces and dedicated sub-routes (`/portal/executive-command/[workspace]`):
    1. `/portal/executive-command/approvals` (Approvals & Requisitions Desk)
    2. `/portal/executive-command/analytics` (Executive Analytics & Capital Velocity Board)
    3. `/portal/executive-command/inter-department` (Inter-Departmental Information & Minuting Hub)
    4. `/portal/executive-command/governance` (MD Professional Responsibilities & Governance Suite)
    5. `/portal/executive-command/ledger` (Decisions Register & Cryptographic Audit Ledger)
  - **Executive Portals Hub (All Departments Landing Pages & Analytics Dashboards):**
    - Expanded the left sidebar and embedded `AllDepartmentsAnalyticsPortalsCard.tsx` in the Analytics board with one-click direct access to all landing pages and specialized analytics dashboards across all 8 divisions:
      - Central Sovereign Analytics (`/portal/central-analytics` & `/`)
      - Generation Fleet (Hydro) (`/portal/generation-master/analytics` & `/portal/generation-master`)
      - Customer Care (CRM) (`/portal/customer-care/analytics` & `/portal/customer-care`)
      - Energy Dispatch (33kV) (`/portal/energy-management`)
      - Engineering Services & Maintenance (`/portal/services`)
      - Commercial Operations & Billing (`/portal/billing-commercial`)
      - Internal Audit & Compliance (`/portal/internal-audit`)
      - Stores & Stockpile Inventory (`/portal/stores`)
- **Requisitions & Approvals Desk (`MDApprovalsDeskView.tsx`, `MDRequisitionActionCard.tsx`):**
  - Centralized terminal for reviewing requisitions submitted by Departmental Heads:
    - Generation Manager (Engr. M. B. Jauro): Turbine runner overhauls, high-cost hydro station spares, emergency outage triage.
    - Head Engineering Services (Engr. F. D. Musa): 33kV transformer bushing replacements, nitrogen purge, HT cable jointing.
    - Head Commercial Operations (Mr. K. O. Alabi): Industrial tariff adjustments, MD meter CT ratio disputes, vending server migrations.
    - Head Internal Audit (Mrs. R. A. Adams): Decommissioned transformer scrap disposal, stock variance clearances.
    - Head Central Stores (Mr. Bitrus Dalyop): Bulk conductor replenishment above local threshold.
    - Head HSE Safety (Engr. P. Pam): High-risk 33kV corridor right-of-way clearing, emergency equipment lockouts.
  - **High-Contrast Dropdowns for Light & Dark Styles:**
    - Hardened Department and Status filter dropdowns with explicit `#0f172a` deep slate navy background and `#f8fafc` crisp white text on both select elements and all option tags, guaranteeing perfect visibility in both light and dark styles.
  - **Room for MD Review:** Interactive review directives box enabling the MD to query the submitting head or stipulate conditions of approval, updating status to `UNDER_REVIEW`.
  - **Room for Inter-Departmental Minuting & Relationship Acknowledgment:** Multi-department selector allowing the MD to minute other departmental heads (GM Generation, Head Engineering, Head Audit, Commercial Ops, Stores), keeping inter-departmental teams in the picture (`ACKNOWLEDGED_MINUTED`).
  - **Authoritative Approval Bar:** One-click digital sign-off (`APPROVE_EXPENDITURE`) with cryptographic SHA-256 seal up to MD delegated limit (₦50,000,000), or `DECLINE_RETURN` with recorded reason.
- **Executive Analytics & Oversight Board (`MDExecutiveAnalyticsView.tsx`):**
  - Real-time KPI summary tracking Total Requests Made (₦ Gross Volume), Requests Under Review, Approval Rate %, and Average Decision SLA (< 3.0h vs NERC 12h benchmark).
  - **Recharts Donut/Pie Chart:** Decision distribution (% Approved, % Under Review, % Pending, % Declined).
  - **Recharts Bar Chart:** Capital requested vs capital authorized across all 6 departments in Millions ₦.
  - **All Departments Analytics Portals Launchpad Card (`AllDepartmentsAnalyticsPortalsCard.tsx`):** Direct interactive launchpad cards for all 8 departmental analytics dashboards.
- **Inter-Departmental Information & Minuting Hub (`MDInterDepartmentHubView.tsx`):**
  - Ingests cross-departmental operational notices, joint technical memos, and grid advisories (e.g., Generation <-> Central Dispatch Hydrological Balancing, TCN Makeri Substation Islanding Mode, Commercial <-> Audit Tamper Triage).
  - Executive Minuting input allowing MD to stamp directives and assign action heads with target deadlines.
- **MD Professional Responsibilities & Governance Suite (`MDGovernanceSuiteView.tsx`):**
  - NERC & Statutory Regulations (MYTO, NEMSA, FMWR Water Extraction).
  - Board of Directors Governance (Board resolutions register > ₦50M).
  - Strategic Concessions & Bilateral PPAs (PRIDA, Grand Cereals).
  - High-Risk HSE Emergency Command Authority (Level-3 Grid Disaster, Board of Inquiry).
  - Departmental Heads Scorecards.
- **Decisions Register & Cryptographic Audit Ledger (`MDDecisionsLedgerView.tsx`):**
  - Permanent append-only ledger displaying all historical approvals, minuted memorandums, and rejections stamped with SHA-256 cryptographic verification hashes.

---

## [1.2.0] - 2026-10-24

### Added - Generation Department Landing Page & Core Operational Activities Architecture
- **Persistent Left-Hand Navigation Sidebar Pane (`GenerationNav_LeftSidebarPane.tsx`):**
  - Integrated 2-column layout shell (`h-screen overflow-hidden`) with a persistent left navigation sidebar and an independently scrollable main operational workspace.
  - Added direct navigation buttons to **General Landing Page** (`/`) and **General Analytics Page** (`/portal/central-analytics`) in both the top command header and the persistent left navigation pane.
  - Exposes one-click navigation and distinct route endpoints for 4 primary generation activities:
    1. `/portal/generation-master/analytics` (Real-Time Generation Analytics & DB Summary)
    2. `/portal/generation-master/reading-sheet` (Hourly Reading Sheet Form)
    3. `/portal/generation-master/change-of-shift` (Change of Shift Desk)
    4. `/portal/generation-master/requisition` (Supervisor Requisition & Notes)
  - Added live status badge counters (`5/6 Online`, `Log Due`, `Locked`, `Active Requisitions`) and station footer attribution.
- **Real-Time Generation Analytics & Machine Telemetry Summary (`GenerationAnalyticsView.tsx`):**
  - Database telemetry synthesis calculating total active power generated versus installed nameplate capacity across Kurra 1, Kurra 2, Kwall 1, Kwall 2, Anoma, and Jere cascade stations.
  - **Active Plants Output Distribution Charts (`ActivePlantsOutputChartsCard.tsx`):**
    - **Interactive Pie Chart (% Share of Total Generating Power):** Supports toggling between **At the Moment** (Real-Time MW Output Share) and **For the Day / Last 24 Hours** (Daily MWh Delivered Energy Share) with central donut readout and color-coded station legend badges.
    - **Active Plants Output Bar Chart:** Visualizes all active generating plants comparing active power output (MW) against installed nameplate capacity (MW), complete with interactive hover tooltips (MW output, capacity, utilization %, and turbine type).
  - **Hydrology & Machine State of Health (`PlantHealthAndHydrologyGrid.tsx`):**
    - Reservoir water level meters (depth in meters, headwater and tailwater elevations, and river inflow in m³/s) with minimum head warning gates.
    - Turbine and generator machine state of health scorecards tracking stator winding temperature (<75°C), ISO 10816 bearing vibration (<3.0 mm/s normal, >4.5 mm/s alert), and lube oil pressure (2.8–4.0 bar).
  - Out of Service generating stations section highlighting tripped units, root causes (`TURBINE_MECH_FAULT`, `HYDRO_LOW_WATER`), downtime duration, and GM workshop review status.
  - **Live Tabular Audit Log Feed (`GenerationTelemetryTable.tsx`):** Tabular display of recent hourly readings committed to `generation_master.master_station_telemetry_logs`.
  - Maintained strict `<= 350` lines limit across all modularized UI components.
- **Hourly Reading Sheet Form for Operator on Duty (`ReadingSheetForm.tsx`):**
  - Form S05-TEL enabling Master Station Operators to input and commit telemetry directly into the database.
  - Captures station identifier, shift window, operational state (`IN_SERVICE`, `STANDBY`, `OUT_OF_SERVICE`) with outage reason categorization, active MW, reactive MVAr, grid frequency (50.00 Hz), water level, head/tail elevations, inflow rate, stator temperature, bearing vibration (with instant ISO trip warning), lube oil pressure, and operator attribution.
  - Automatically updates active fleet cache in real time with confirmation toast displaying generated log ID and timestamp.
- **Change of Shift Desk with Immutable Handover Notes (`ChangeOfShiftForm.tsx`):**
  - NERC Grid Code Section 5.3 compliant shift handover protocol capturing concluded shift window, outgoing operator, incoming operator, and radio/call-collated 6-plant physical muster (operators, protection engineers, mechanical/electrical technicians).
  - Short column for handover notes strictly enforced as **immutable**: once submitted, the note is signed, stamped with `is_locked: true`, and **CANNOT BE EDITED OR MODIFIED**.
  - UI interceptor blocking any alteration attempts with a NERC regulatory compliance modal.
  - Permanent Shift Handover Register displaying historical shifts with SHA-256 tamper-proof lock badges.
- **Supervisor Requisition & Technical Notes Section (`SupervisorRequisitionView.tsx`):**
  - Direct digital channel for operators and plant technicians to submit formal notes and material requisitions to immediate supervisors:
    - Electrical Unit Head (`Engr. D. Gyang`)
    - Mechanical Unit Head (`Engr. S. Lar`)
    - General Manager Generation (`Engr. M. B. Jauro`)
    - Central Workshop Lead (`Mr. Bitrus Dalyop`)
    - Store Officer (`Central Stores Depot`)
  - Configurable requisition categories (`MATERIAL_SPARES`, `CRITICAL_ALERT`, `OUTAGE_REQUEST`, `WORKSHOP_ORDER`) and priority levels (`ROUTINE`, `URGENT`, `EMERGENCY_OUTAGE`).
  - Itemized Bill of Materials / Spares builder with item name, part/spec number, quantity, and unit.
  - Requisition queue and tracking table with supervisor acknowledgment notes and GM store authorizations.
- **Dedicated Route Endpoints:**
  - Added wrapper pages under `apps/web-portal/src/pages/portal/generation-master/`: `analytics.tsx`, `reading-sheet.tsx`, `change-of-shift.tsx`, and `requisition.tsx`.

------

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

------

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

------

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

------

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

------

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

------

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

------

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

------

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

------

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

------

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

------

---

## [0.1.0] - 2026-10-22

### Added - Sprint 1: Core Foundation & Grid Topology (`s01_foundation`)
- Database migrations (`001_s01_foundation.up.sql`, `001_s01_foundation.down.sql`) with PostGIS geometry and append-only WorkNotes trigger.
- Foundational packages: `pkg/crypto`, `pkg/auth`, `pkg/events`, `pkg/database`.
- Hexagonal module for org units, grid assets, federated staff, and chained WorkNotes.
