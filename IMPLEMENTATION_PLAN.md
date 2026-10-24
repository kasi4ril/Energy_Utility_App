# AGENTIC IMPLEMENTATION RUNBOOK & CONSTITUTION
**Target Project:** Enterprise Utility Workflow, Grid Intelligence & Local AI Governance Platform (DisCo & Generation)

## 1. Role & Operating Directives for the Agent
You are an expert Senior Systems Architect and Go/Python Software Engineer with 25 years of experience building mission-critical utility software. You are executing a strict **9-Sprint Modular Monolith Build**. Never skip architectural layers, never bypass safety/audit triggers, and never introduce cloud-dependent LLM calls.

## 2. Non-Negotiable Core Engineering Rules
1. **Append-Only & Cryptographic Integrity:** 
   - Every sensitive table (`customers`, `staff_profiles`, `grid_assets`, `vending_transactions`) must compute a `row_hmac_signature`.
   - All inter-departmental interactions must use `foundation.work_notes` with cryptographic SHA-256 prev-hash chaining. WorkNotes are immutable at the database trigger level.
2. **Maker-Checker Security:** No direct SQL `UPDATE` is permitted on sensitive financial, customer, or inventory states without passing through a Maker-Checker approval workflow and Ollama risk-scoring.
3. **100% Air-Gapped AI:** All LLM calls, Text-to-SQL, anomaly detection, and `@GridAssistant` chat lookups must route locally to the `ai-guardian` FastAPI/Ollama container (`http://ollama:11434`).
4. **Offline-First Field PWA & Remote Call Collation:** 
   - Field Team Leads (Feeder/Metering) log proof-of-work via PWA with live GPS verification against asset coordinates.
   - Generation Master Station (Plants 1–6) and Energy Management (33kV Substations) ingest remote shift data via standardized Google Forms Webhooks (`HMAC-signed`).

## 3. Sprint Execution Roadmap (Current Progress)
- [x] **Sprint 1 (`s01_foundation`):** Org Hierarchy, Grid Topology, IAM, Federated HR, PostGIS Spatial Indexing, and Immutable WorkNotes. *(Completed - Sign-off Granted)*
- [x] **Sprint 2 (`s02_ollama_guardian`):** Hash-Chain Ledger, Maker-Checker Engine, `@GridAssistant` Chat Bot & Report Formatter. *(Completed - Sign-off Granted)*
- [ ] **Sprint 3 (`s03_customer_care`):** KYC Onboarding, Anti-Ghost GPS Lock, Omnichannel Ticketing, Direct Routing, and Mandatory Customer-Confirmation Closure Gate. *(In Progress)*
- [x] **Sprint 4 (`s04_services_maint`):** Feeder, Transformer & Metering Unit Portals, Field Team Lead PWA, Feeder 5-Workspace Landing Page Architecture (`analytics`, `outage-request`, `proof-of-work`, `standards`, `work-orders`), and Anti-Ghost GPS Proof-of-Work Verification. *(Feeder Landing Page Architecture Implemented)*
- [x] **Sprint 5 (`s05_generation_master`):** Master Station Hub (Plants 1–6), Water Level & Generator Telemetry, Google Forms Webhook Ingestion, Workshop Loop, GM Material Requisition Gate, and 4-Activity Left-Hand Navigation Architecture (`analytics`, `reading-sheet`, `change-of-shift` with immutable notes, `requisition`). *(Completed - Sign-off Granted)*
- [ ] **Sprint 6 (`s06_energy_management`):** Central Dispatch Hub, Bulk MW Intake (Plants 1–6 + Dormant TCN Hook), Priority Feeder Rationing, and 5-Step Safety PTW Substation Isolation.
- [ ] **Sprint 7 (`s07_stores_audit`):** Store Ledger, `@GridAssistant` Stock Bot, Internal Audit Gate, Staff Allowance Claims (Overtime/Travel 3-way GPS/Roster cross-check), and MD Out-of-Stock Supplier Approval.
- [ ] **Sprint 8 (`s08_billing_commercial`):** Zero-Estimation Billing, Meter-Purchase Refund Ledger, Multi-Platform Vending Hub (Online, OPay, Moniepoint, GTBank), Wrong Billing Commercial Services Head Gate, and Energy Theft/Churn AI Models.
- [x] **Sprint 9 (`s09_executive_command`):** Managing Director Tier-1 Command Center & 5-Workspace Left-Nav Architecture (`approvals`, `analytics`, `inter-department`, `governance`, `ledger`), Requisition Reviews, Inter-Department Minuting, and C-Suite Governance Suite. *(UI Architecture & Executive Landing Page Implemented)*

## 4. How to Prompt the Agent for Subsequent Sprints
When moving to the next sprint, provide your coding agent with this prompt pattern:
> *"Read `AGENTS_MASTER_RUNBOOK.md` and implement **Sprint X (Module Name)**. Ensure all database schemas, Go backend domain/repo/handlers, and integration hooks align perfectly with the bounded context rules and cryptographic audit guarantees defined in the runbook."*

## 5. Mandatory TDD, Study Documentation & Code Quality Standards

The Coding Agent MUST strictly enforce the following 6 software engineering standards on every single Sprint and task:

1. **Test-Driven Development (TDD — Red-Green-Refactor):**
   - No feature or endpoint is complete without automated tests.
   - For every Go (`*_test.go`), Python (`test_*.py`), and Frontend (`*.test.tsx`) module, write unit and integration tests covering both valid workflows and security/boundary rejections (e.g., testing that unauthorized roles, invalid HMAC signatures, or out-of-range PostGIS GPS coordinates are strictly rejected).
   - Minimum test coverage target per module is **80%**.

2. **Mandatory Module Study Guide (`STUDY_GUIDE.md` in Every Folder):**
   - For every sprint module (`s01` to `s09`) and frontend feature built, the Agent MUST create and update a `STUDY_GUIDE.md` file inside that module's folder (plus a root `PROJECT_CHANGELOG.md`).
   - Each `STUDY_GUIDE.md` must clearly document for the Engineer's learning and review:
     - **Purpose & Architectural Role:** Why this module exists and what utility problem it solves.
     - **File-by-File Breakdown:** What each file and function in the folder does.
     - **Step-by-Step Data & Execution Flow:** How a request enters, gets validated, hits the DB/AI, and returns.
     - **Key Engineering Patterns Used:** Explanation of Go concurrency, PostGIS spatial logic, cryptography, or Python ML concepts applied so the engineer can study and master the implementation.
     - **How to Run & Test:** Exact terminal commands to run the tests and verify the module locally.

3. **Strict 350-Line File Size Limit:**
   - No single `.go`, `.py`, or `.tsx` file may exceed **350 lines of code**. If a file approaches 350 lines, split it into logical, well-named sub-files within the same package.

4. **Reversible Database Migrations (`.up.sql` and `.down.sql`):**
   - Every database migration in `migrations/` must include both an `.up.sql` file and a corresponding `.down.sql` rollback file.

5. **Zero-Warning Static Analysis & Linting:**
   - All Go code must pass `go vet` and `golangci-lint`.
   - All Python code must use strict type hints and pass `ruff` and `mypy`.
   - All Next.js TypeScript code must pass `tsc --noEmit` with strict mode enabled (`no any` types).

6. **Conventional Commits & Interactive API Docs:**
   - Maintain an updated OpenAPI/Swagger specification (`docs/openapi.yaml`) for all REST endpoints.
   - Structure all work into atomic steps using Conventional Commit messages (`test(s01): ...`, `feat(s01): ...`, `docs(s01): ...`).

7. **Professional Inline Code Documentation & Commenting:**
   - Every Go package, Python module, and TypeScript component written by the Agent MUST include thorough, professional documentation comments.
   - **Package & File Headers:** Every file must begin with a docstring or package comment explaining its purpose, architecture layer, and how it interacts with other modules.
   - **Public Functions & Structs:** Every exported function, method, struct, and interface must have a clear GoDoc/Docstring comment explaining what it does, its parameters, return values, and any security/cryptographic side effects (e.g., HMAC signature generation or PostGIS spatial checks).
   - **Complex Business Logic:** Complex algorithms (such as the 5-Step Safety PTW lock, the 3-Way Overtime/Travel claims cross-check, or the Energy Theft/Bypass detection formulas) must include step-by-step explanatory comments so the Engineer can easily study and understand the implementation.
---

## SECTION 11: ADMINISTRATIVE CONTROLS, STAFF TRANSFERS & NIGERIAN GRID COMPLIANCE

### 11.1 Staff Transfers & Organizational Re-Assignment Workflow
- **Initiation:** HQ HR or Plant HR Officer opens the Staff Directory, selects the staff member, and submits a **Transfer Request Form** specifying the new `org_unit_id`, new `primary_role`, and effective date.
- **Backend Enforcement:** The Go Core API updates `foundation.staff_profiles`, re-computes the cryptographic `row_hmac_signature`, publishes a `staff.transferred` event to NATS JetStream, and invalidates active JWT sessions via Redis.
- **Frontend Effect:** Upon next login, Next.js middleware automatically updates the user's scope, removing access to their old district/plant data and granting immediate operational access to their new deployment.

### 11.2 Nigerian Electricity Regulatory Commission (NERC) & NEMSA Compliance Standards
- **Grid Code & Safety (NEMSA):** All switching operations, 33kV/11kV outages, and maintenance procedures strictly require the 5-step Permit-to-Work (PTW) signed by a certified Safety Representative (HSE) before line clearance is granted.
- **Customer Service Standards (NERC):** Service ticket SLAs, wrong-billing resolution loops, and customer confirmation call-backs adhere strictly to regulatory timelines, generating automated compliance breach reports for internal audit and executive review.
- **Tariff & Metering Regulations:** Tariff bands (A–E) are structurally tied to Feeder asset registers, and customer-financed meter installations automatically populate the regulatory Meter-Purchase Refund Ledger for transparent amortization.

---

## 6. Mandatory Sprint Gate Verification Protocol (Definition of Done)

**CRITICAL DIRECTIVE FOR THE CODING AGENT:** 
You are strictly FORBIDDEN from marking any Sprint (`s01` through `s09`) as `[x] Completed` or generating code for the next Sprint until you have executed the **6-Point Sprint Gate Verification** below and presented the passing output to the Lead Engineer for sign-off.

### Gate 1: Reversible Database Migration & Schema Verification
Before writing or finalizing application logic, verify that the Sprint's SQL migrations apply cleanly, enforce constraints, and roll back without orphan artifacts:
1. **Apply Up Migration:** Run `migrate -path ./migrations -database "$DB_URL" up 1` -> Must succeed with `0 errors`.
2. **Test Rollback (Down):** Run `migrate -path ./migrations -database "$DB_URL" down 1` -> Must cleanly drop only the current Sprint's objects.
3. **Re-Apply Up Migration:** Run `migrate -path ./migrations -database "$DB_URL" up 1` -> Must restore the schema cleanly.
4. **Trigger & Extension Check:** Verify that PostGIS spatial indexes (`GIST`), TimescaleDB hypertables (where applicable), and immutable audit/HMAC triggers reject raw unauthorized `UPDATE` or `DELETE` statements.

## 6. Mandatory Sprint Gate Verification Protocol (Definition of Done)

**CRITICAL DIRECTIVE FOR THE CODING AGENT:** 
You are strictly FORBIDDEN from marking any Sprint (`s01` through `s09`) as `[x] Completed` or generating code for the next Sprint until you have executed all 6 checks below and received explicit approval from the Lead Engineer.

### Gate 1: Reversible Database Migration & Schema Verification
1. **Apply Up Migration:** Run `migrate -path ./migrations -database "$DB_URL" up 1` -> Must succeed with `0 errors`.
2. **Test Rollback (Down):** Run `migrate -path ./migrations -database "$DB_URL" down 1` -> Must cleanly drop only the current Sprint's objects.
3. **Re-Apply Up Migration:** Run `migrate -path ./migrations -database "$DB_URL" up 1` -> Must restore the schema cleanly.
4. **Trigger & Extension Check:** Verify that PostGIS spatial indexes (`GIST`), TimescaleDB hypertables, and immutable SHA-256 / HMAC database triggers reject raw unauthorized `UPDATE` or `DELETE` statements.

### Gate 2: Automated TDD Suite & Security Boundary Verification (Minimum 80% Coverage)
1. **Go Core Backend Tests:**
   - Command: `cd services/core-api && go test -v -race -coverprofile=coverage.out ./internal/sXX_[module_name]/... && go tool cover -func=coverage.out`
   - Must pass 100% of tests with zero race conditions and `>= 80.0%` coverage.
   - Must include negative security tests proving that unauthorized RBAC roles, tampered `row_hmac_signature` values, broken SHA-256 `prev_hash` chains, out-of-radius PostGIS GPS coordinates, and bypassed Maker-Checker states are strictly rejected.
2. **Python / Ollama AI Sidecar Tests:**
   - Command: `cd services/ai-guardian && pytest --cov=app tests/test_sXX_[module_name].py`
   - Must verify air-gapped routing to local Ollama (`http://ollama:11434`) and deterministic fallback handling.
3. **Next.js Frontend & PWA Tests:**
   - Command: `cd apps/web-portal && npm run test -- --watchAll=false src/features/[dept-name]`

### Gate 3: Zero-Warning Linting, Type Safety & 350-Line File Limit
1. **Go Static Analysis:** `cd services/core-api && go vet ./... && golangci-lint run` -> Must return `0 issues`.
2. **Python Type & Lint Check:** `cd services/ai-guardian && ruff check app/ && mypy app/` -> Must return `0 issues`.
3. **Frontend TypeScript Check:** `cd apps/web-portal && npx tsc --noEmit && npm run lint` -> Must return `0 errors` (Zero `any` types allowed).
4. **350-Line Modularity Cap:** Every `.go`, `.py`, and `.tsx` file must be **350 lines or fewer**. Any file exceeding 350 lines must be refactored into smaller sub-files before passing Gate 3.

### Gate 4: Predictable Symmetry & Inline Code Commenting Check
- Confirm the Go module contains `domain/models.go`, `domain/rules.go`, `repository/postgres_repo.go`, `service/workflow_service.go`, `transport/http_handlers.go`, and `transport/dto.go`.
- Confirm the Frontend feature folder contains `ZoneB_DashboardRibbon.tsx`, `ZoneC_ActionQueue.tsx`, `ZoneC_IntakeForm.tsx`, and `api.ts`.
- Confirm every file starts with an architectural header comment, every exported struct/function has a complete GoDoc/Docstring comment, and complex utility/regulatory algorithms (NERC, NEMSA 5-Step PTW, MAP Meter Refund, Energy Theft, 3-Way Claims Check) have step-by-step inline comments.

### Gate 5: Living Study Guide (`STUDY_GUIDE.md`) Verification
Confirm `services/core-api/internal/sXX_[module_name]/STUDY_GUIDE.md` (and root `PROJECT_CHANGELOG.md`) has been created and includes:
1. **Purpose & Utility Domain Context:** Why this module exists and what operational problem it solves.
2. **File-by-File Breakdown:** What every file and key function in the module does.
3. **Step-by-Step Execution Flow:** How data flows from the UI/Webhook through the Go API, Ollama AI Sidecar, and PostgreSQL schema.
4. **Engineering Concepts for Personal Study:** Clear explanation of the Go, SQL/PostGIS, Cryptography, or Python ML patterns used.
5. **How to Run & Verify Locally:** Exact terminal and `curl` commands to test the module.

### Gate 6: Mandatory Pause & Sign-Off Scorecard
At the end of every Sprint, the Coding Agent MUST pause, print the following completed scorecard in the chat, and wait for the Engineer to reply **"APPROVED — PROCEED TO NEXT SPRINT"** before writing any code for the next Sprint:

================================================================================
SPRINT GATE VERIFICATION REPORT: [Sprint Number & Module Name]
================================================================================
1. Database Migrations (.up.sql / .down.sql):   [PASS / FAIL]
2. TDD Unit & Security Tests (Coverage %):      [PASS / FAIL - XX.X% Coverage]
3. Linter, Type Check & <= 350 Line File Cap:   [PASS / FAIL - Max lines: XXX]
4. Section 0.4 Folder Symmetry & Code Comments: [PASS / FAIL]
5. STUDY_GUIDE.md & OpenAPI Spec Updated:       [PASS / FAIL]
================================================================================
Status: Waiting for Lead Engineer review and sign-off before proceeding.

