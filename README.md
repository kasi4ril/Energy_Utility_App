# NESCO Enterprise Utility Platform

> **Mission-Critical Utility Workflow, Grid Intelligence & Air-Gapped Sovereign AI Platform for NESCO Nigeria Ltd.**

The **NESCO Enterprise Utility Platform** is a sovereign, event-driven enterprise application engineered to unify power generation, transmission dispatch, electrical field maintenance, commercial billing, warehouse inventory, internal audit, and executive command under a single cryptographic governance framework.

Designed for zero SaaS dependency and 100% on-premise execution, the platform embeds a local Ollama AI Guardian sidecar for automated fraud risk scoring, Maker-Checker segregation of duties, spherical Haversine GPS field proof-of-work verification, and conversational grid intelligence (`@GridAssistant`).

---

## Architecture Overview

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

## 10-Sprint Module Breakdown

| Sprint | Bounded Context | Description & Capabilities |
| :--- | :--- | :--- |
| **S01** | `s01_foundation` | Multi-tier org hierarchy (HQ $\to$ Area $\to$ Feeder $\to$ Unit), PostGIS grid topology, federated staff IAM/RBAC, append-only cryptographic WorkNotes hash-chain. |
| **S02** | `s02_ollama_guardian` | Air-gapped local AI sidecar, Maker-Checker fraud scoring, SHA-256 prev-hash audit ledger, in-app conversational `@GridAssistant`. |
| **S03** | `s03_customer_care` | Anti-ghost debtor lock, KYC customer onboarding with premises GPS, direct SLA ticket routing, 2-step customer phone confirmation closure gate. |
| **S04** | `s04_services_maint` | Unit maintenance portals (Feeder, Transformer, Metering), Field PWA, mandatory GPS proof-of-work ($\le 500\text{m}$ proximity) & photo verification. |
| **S05** | `s05_generation_master` | Remote plant collation (Plants 1–6) at Master Station, hydrological/thermal telemetry, Google Form A HMAC-SHA256 webhook, GM workshop loop, and Section 13 Left-Hand Nav Hub (Real-Time DB Analytics, Hourly Reading Sheet, Change of Shift with Immutable Notes, Supervisor Requisitions). |
| **S06** | `s06_energy_management` | Central Dispatch, bulk MW intake (6 Hydro/Thermal + TCN Hook), priority feeder rationing (Bands A–E), Google Form B webhook, 5-Step Safety Permit-to-Work (PTW) lock. |
| **S07** | `s07_stores_audit` | Dual-Path material requisitions (Path A In-Stock vs Path B MD Procurement), reverse logistics scrap tracking, automated 3-way staff claims audit desk. |
| **S08** | `s08_billing_commercial` | Zero-estimation metered billing, 20-digit STS token generation, NERC customer meter refund amortization ledger, wrong billing ticket reversion loop, theft/churn AI detection. |
| **S09** | `s09_executive_command` | C-Suite shared pulse dashboard, immutable operational telemetry snapshots, role-segregated digital sign-off queue, 4 condition-based predictive maintenance (CBM) models with autonomous S04 dispatch. |
| **S10** | `s10_fullstack_wiring` | Next.js Unified Login with Shift Station selector, Edge Role Auto-Router, 4-Zone Enterprise UI Shell, Field PWA offline cache, Docker Compose & NATS orchestration. |

---

## Technology Stack

- **Backend Core:** Go 1.22 (Standard Library HTTP, Clean Hexagonal Architecture, zero heavyweight web frameworks)
- **Database & Spatial Engine:** PostgreSQL 16 with PostGIS 3.4 (Spherical Haversine GPS calculation, spatial indexes, append-only trigger rules)
- **Event Streaming & Pub/Sub:** NATS JetStream 2.10
- **AI Governance Sidecar:** Python 3.11 with FastAPI, LangChain, and local air-gapped Ollama (`llama3` / `mistral`)
- **Web Portal:** Next.js 14, React 18, TypeScript, Tailwind CSS, Lucide Icons
- **Field PWA:** HTML5/TypeScript Progressive Web App, Service Worker cache, IndexedDB/localStorage offline queue, Geolocation API
- **Container Orchestration:** Docker Compose multi-stage builds

---

## Repository Structure

```
Energy_Utility_App/
├── .env.example                               # Production environment variable template
├── .gitignore                                 # Git ignore patterns
├── docker-compose.yml                         # Multi-container orchestration topology
├── BLUE_PRINT.md                              # Enterprise blueprint & constitutional requirements
├── IMPLEMENTATION_PLAN.md                     # 10-sprint execution plan & gate protocol
├── MILESTONES.md                              # Living sprint milestone execution ledger
├── PROJECT_CHANGELOG.md                       # Comprehensive release changelog
├── README.md                                  # Executive summary & quick-start guide
├── apps/
│   ├── web-portal/                            # Next.js 14 4-Zone Enterprise Web Portal
│   └── field-pwa/                             # Offline-first Field Crew Progressive Web App
├── docs/
│   └── openapi.yaml                           # OpenAPI 3.0.3 API specification
├── forms-specs/
│   ├── google_form_a_generation_master.js     # Master Station Google Forms HMAC forwarder
│   └── google_form_b_energy_dispatch.js       # Central Dispatch Google Forms HMAC forwarder
├── migrations/                                # Up/Down SQL migrations (001 through 009)
└── services/
    ├── ai-guardian/                           # Python 3.11 Local AI Sovereign Sidecar
    └── core-api/                              # Go 1.22 Modular Monolith Backend
        ├── cmd/server/main.go                 # Enterprise HTTP server & gateway
        ├── internal/                          # S01 through S10 hexagonal bounded contexts
        ├── pkg/                               # Shared packages (auth, crypto, database, events)
        └── test/integration/                  # Multi-sprint end-to-end integration test suites
```

---

## Quick Start & Verification

### 1. Run Complete Platform with Docker Compose
```bash
cp .env.example .env
docker-compose up --build -d
```
- Web Portal: `http://localhost:3000`
- Core API Gateway: `http://localhost:8080/healthz`
- AI Guardian Sidecar: `http://localhost:8001/health`

### 2. Run Go Test Suite (All 40 Packages)
```bash
cd services/core-api
go test -count=1 ./...
```

### 3. Run End-to-End Cross-Sprint Integration Tests
```bash
cd services/core-api
go test -v ./test/integration/...
```

### 4. Run Python AI Guardian Test Suite
```bash
cd services/ai-guardian
python3 -m unittest discover -s tests
```

---

## Architectural Constraints & Compliance

- **Strict File Limit:** Every single code file across all languages is strictly $\le 350$ lines of code (maximum file in repository is 336 lines).
- **Test Coverage:** $\ge 80\%$ statement test coverage with negative security and boundary test cases across all bounded contexts.
- **Data Sovereignty:** 100% on-premise execution with zero external vendor telemetry or public SaaS cloud dependencies.
- **Cryptographic Immutability:** SHA-256 previous-hash pointer verification enforced on all WorkNotes and audit ledger tables.

---

## License
Proprietary & Confidential. All rights reserved by **NESCO Nigeria Ltd.**
