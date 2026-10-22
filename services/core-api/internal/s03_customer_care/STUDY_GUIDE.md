# Sprint 3 Study Guide: Customer Care, Anti-Ghosting & Ticketing Engine (`s03_customer_care`)

## 1. Purpose & Utility Domain Context
The **Customer Care & Omnichannel Ticketing Bounded Context (`s03_customer_care`)** serves as the operational frontline and anti-fraud checkpoint of NESCO. In Nigerian electricity distribution (DisCo) operations:
- **Anti-Ghost KYC Registration:** Power theft and billing leakage often arise from "ghost" accounts created without physical connection to a Distribution Transformer (DT) or created by delinquent debtors evading disconnection. All new customer premises GPS coordinates are strictly enforced to reside within 500 meters of their assigned DT, and prior debtor histories are locked.
- **Direct SLA Routing:** To eliminate internal ticket bouncing and bureaucratic delays, complaints are deterministically routed at ingestion based on fault classification:
  - `FEEDER_LINE_FAULT` -> `UNIT_HEAD_FEEDER_MAINT` (Supervised by `HEAD_SERVICES`)
  - `TRANSFORMER_FUSE_BURNT` -> `UNIT_HEAD_TX_MAINT` (Supervised by `HEAD_SERVICES`)
  - `METER_FAULT` / `NEW_METER_INSTALL` -> `UNIT_HEAD_METERING` (Supervised by `HEAD_SERVICES`)
  - `BILLING_DISPUTE` / `OVERBILLING` -> `UNIT_HEAD_BILLING` (Supervised by `HEAD_COMMERCIAL_OPS`)
- **Strict Customer-Confirmation Closure Gate:** Technical field crews can *only* mark a ticket `ATTENDED_TO` with verifiable proof notes. A ticket can **never** be marked `CLOSED_RESOLVED` by field technicians. Only a `CUSTOMER_CARE_OFFICER` can close the ticket after conducting an outbound customer verification call. If the customer indicates the issue persists, the ticket immediately transitions to `REOPENED_ESCALATED`.
- **Meter Operations Dispatch:** Formalized scheduling pipeline (`NEW_INSTALL`, `REPLACEMENT`, `RECERTIFICATION`, `TAMPER_INSPECTION`) linking customer tickets directly to metering teams.

---

## 2. File-by-File Architectural Breakdown

```text
services/core-api/internal/s03_customer_care/
├── domain/
│   ├── models.go       # Structs & Enums: Customer, ServiceTicket, MeterSchedule, TicketCategory, TicketStatus
│   ├── rules.go        # Pure business rules: Haversine <= 500m DT distance, Direct Routing map, Closure Gate validation
│   └── rules_test.go   # Comprehensive domain unit tests (86.1% statement coverage)
├── repository/
│   ├── repository.go                 # Interface & thread-safe in-memory mock repository
│   ├── postgres_repo.go              # PostgreSQL + PostGIS queries for Customers
│   ├── postgres_repo_tickets.go      # PostgreSQL queries for Tickets & Meter Schedules
│   └── repository_test.go            # Storage layer unit tests
├── service/
│   ├── workflow_service.go      # Business orchestration: KYC checks, SLA computation, Ticket transitions, Events
│   └── workflow_service_test.go # Orchestration tests, anti-ghost checks, closure gate enforcement (81.0% coverage)
└── transport/
    ├── dto.go                 # DTOs with strict input validation for customer, ticket, and schedule endpoints
    ├── http_handlers.go       # net/http REST endpoints & RBAC middleware wiring (84.5% coverage)
    └── http_handlers_test.go  # End-to-end HTTP integration & role enforcement tests
```

### Database Migrations:
- `migrations/003_s03_customer_care.up.sql`: Creates `customercare.customers`, `customercare.service_tickets`, and `customercare.meter_schedules` with PostGIS geometry, foreign key constraints to `grid_assets`, and index optimizations.
- `migrations/003_s03_customer_care.down.sql`: Reversible rollback dropping the schema cleanly.

---

## 3. Step-by-Step Data & Execution Flow

### A. Anti-Ghost Customer Onboarding Flow
```
[ Field Officer / Web Portal ]
           │
           │  POST /api/v1/customers/onboard
           ▼
[ Transport Layer: http_handlers.go ]
   - Validates FullName, Phone, TransformerID, BuildingClass, TariffBand
           │
           ▼
[ Service Layer: workflow_service.go ]
   - Fetches assigned Distribution Transformer (DT) coordinates from Grid Repo
   - Calculates Haversine distance between customer premises GPS and DT GPS
   - Enforces Anti-Ghost Lock: If distance > 500m -> rejects with ErrDTDistanceExceeded
   - Checks customer phone/NIN against prior debtor balance records
   - Generates unique Account Number (`NES-CUST-...`) and CIN (`CIN-...`)
   - Computes SHA-256 HMAC row signature for data integrity
   - Persists to customercare.customers
   - Publishes `customer.onboarded` asynchronous event
```

### B. Direct Routing & Two-Step Closure Gate Flow
```
1. Customer / Call Center Logs Ticket:
   POST /api/v1/tickets -> Service automatically resolves Unit Head and Dept Head Supervisor:
   - Sets SLA Due Date (e.g., 24h for Transformer Fault, 72h for Billing Dispute)
   - Status starts at OPEN_PENDING_DISPATCH.

2. Technical Crew Attends Ticket:
   POST /api/v1/tickets/{id}/attend (Role: FIELD_TECHNICIAN, LINEMAN, METER_INSTALLER)
   - Must submit non-empty AttendedProofNote.
   - Status transitions from OPEN to ATTENDED_TO.
   - Field crew CANNOT close ticket.

3. Customer Care Outbound Verification:
   POST /api/v1/tickets/{id}/confirm-close (Role: CUSTOMER_CARE_OFFICER)
   - Conducts phone call with consumer to verify electrical restoration.
   - If Customer Confirmed == true:
     -> Status transitions to CLOSED_RESOLVED.
   - If Customer Confirmed == false:
     -> Status transitions to REOPENED_ESCALATED.
```

---

## 4. Concrete cURL Examples

### 1. Onboard New Customer (Anti-Ghost Verification)
```bash
curl -X POST http://localhost:8080/api/v1/customers/onboard \
  -H "Authorization: Bearer <TOKEN_CUSTOMER_CARE>" \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Alhaji Musa Ibrahim",
    "phone": "+2348031234567",
    "nin": "12345678901",
    "building_class": "COMMERCIAL",
    "latitude": 9.9320,
    "longitude": 8.8920,
    "transformer_id": "00000000-0000-0000-0000-000000000002",
    "meter_serial": "NESCO-MTR-9988",
    "tariff_band": "BAND_A"
  }'
```

### 2. Log Service Complaint Ticket (Auto-Routing)
```bash
curl -X POST http://localhost:8080/api/v1/tickets \
  -H "Authorization: Bearer <TOKEN_CUSTOMER_CARE>" \
  -H "Content-Type: application/json" \
  -d '{
    "channel": "CALL_CENTER",
    "category": "TRANSFORMER_FUSE_BURNT",
    "description": "Explosion sound on 500kVA DT feeder pillar, phase 2 dead.",
    "assigned_unit_id": "00000000-0000-0000-0000-000000000001"
  }'
```

### 3. Technical Crew Submits Proof of Work
```bash
curl -X POST http://localhost:8080/api/v1/tickets/<ticket-uuid>/attend \
  -H "Authorization: Bearer <TOKEN_FIELD_CREW>" \
  -H "Content-Type: application/json" \
  -d '{
    "proof_note": "Replaced blown 200A HRC fuse on red phase. Voltage verified at 230V phase-to-neutral."
  }'
```

### 4. Customer Care Verifies & Confirms Closure
```bash
curl -X POST http://localhost:8080/api/v1/tickets/<ticket-uuid>/confirm-close \
  -H "Authorization: Bearer <TOKEN_CUSTOMER_CARE>" \
  -H "Content-Type: application/json" \
  -d '{
    "confirmed": true,
    "feedback": "Spoke with Alhaji Musa; confirmed shop power restored and voltage stable."
  }'
```

---

## 5. Security & Domain Invariants
1. **Haversine Proximity Guarantee:** Distance calculation accounts for Earth curvature. A distance exceeding 500.0 meters from the assigned DT causes immediate rejection to halt unauthorized load allocations.
2. **Technician Closure Prevention:** Even administrators or super-users cannot bypass the requirement that only a `CUSTOMER_CARE_OFFICER` can mark a ticket `CLOSED_RESOLVED` after live customer contact.
3. **Audit Immutability:** All ticket transitions and proof notes are recorded with timestamps, actor IDs, and emitted to the NATS event bus for historical audit trail logging.

