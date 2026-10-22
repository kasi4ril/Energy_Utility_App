# Sprint 4 Study Guide: Services & Maintenance Department (`s04_services_maint`)

## 1. Purpose & Utility Domain Context
The **Services & Maintenance Bounded Context (`s04_services_maint`)** coordinates physical field engineering operations across three vital utility sub-units under the **Head of Services**:
1. **Feeder Maintenance Unit (`FEEDER_MAINT`):** Responsible for 33kV and 11kV transmission/distribution corridors, overhead line spans, crossarms, disc/pin insulators, and pole erection.
2. **Transformer Maintenance Unit (`TRANSFORMER_MAINT`):** Responsible for 33/11kV injection substations, distribution transformers (DTs), HT/LT dropout fuses, silica gel breathers, and transformer oil dielectric testing/topping.
3. **Metering Unit (`METERING`):** Responsible for single-phase, three-phase, and Maximum Demand (MD) customer meter installations, tamper inspections, and cryptographic seal integrity.

### The Closed-Loop Field Governance Problem
Historically in electrical utilities, field operations suffer from phantom repairs: technicians claiming to have visited a remote fault site or swapped a transformer fuse without actually deploying, or submitting previously taken photos.
Sprint 4 implements **closed-loop field governance**:
- **Haversine GPS Locking:** Device GPS captured by the Team Lead's mobile PWA must strictly fall within permissible radius of the registered asset (Feeder: 150m, Transformer: 75m, Meter: 50m).
- **Anti-Photo-Reuse Gate:** Before and after photos are hashed with SHA-256. Identical before/after hashes or duplicate reuse of historic photos are rejected as fraud.
- **PTW Substation Isolation Lock:** High-voltage work orders marked `ptw_required: true` cannot proceed without a valid Permit-to-Work (PTW) safety permit issued by Central Dispatch.
- **Unit Head Verification & Customer Care Reversion:** The field technician cannot mark a customer ticket attended directly. The Team Lead submits proof -> Unit Head verifies technical validity -> Unit Head approves and automatically transitions the linked Customer Care ticket to `ATTENDED_TO`. Customer Care then performs the final customer confirmation call.

---

## 2. File-by-File Architectural Breakdown

```text
services/core-api/internal/s04_services_maint/
├── domain/
│   ├── models.go       # MaintenanceWorkOrder, ProofOfWorkLog, UnitType, WorkOrderType, WorkOrderStatus
│   ├── rules.go        # Haversine proximity calculations, unit GPS tolerances, anti-photo-reuse check, PTW lock
│   └── rules_test.go   # Pure domain tests (100.0% coverage)
├── repository/
│   ├── repository.go             # Storage interface & MemoryRepository mock implementation
│   ├── postgres_repo.go          # Postgres queries for MaintenanceWorkOrders
│   ├── postgres_repo_proofs.go   # Postgres queries for ProofOfWorkLogs (split to satisfy <= 350 lines)
│   └── repository_test.go        # Storage tests
├── service/
│   ├── workflow_service.go       # Orchestrator: Work order dispatch, proof submission, AI watchdog, ticket reversion
│   └── workflow_service_test.go  # Orchestration tests, mock AI sidecar, customer ticket linkage (80.0% coverage)
└── transport/
    ├── dto.go                 # DTOs: CreateWorkOrderDTO, SubmitProofDTO, VerifyProofDTO with Validate() methods
    ├── http_handlers.go       # net/http REST endpoints & RBAC authentication wiring
    └── http_handlers_test.go  # End-to-end HTTP integration tests (82.5% coverage)
```

---

## 3. Step-by-Step Data & Execution Flow

```
[ Step 1: Work Order Dispatch ]
Unit Head receives complaint from Customer Care or grid inspection:
  POST /api/v1/services/work-orders
  -> Service validates target Grid Asset exists in GIS topology.
  -> Generates order number: WO-YYYYMMDD-XXXX.
  -> Persists work order with status DISPATCHED_TO_FIELD.
  -> Publishes 'services.work_order.dispatched' event to notify Team Lead on Mobile PWA.

[ Step 2: Field Crew Execution & Proof Ingestion ]
Team Lead arrives on site, performs repair, and submits proof via Mobile PWA:
  POST /api/v1/services/work-orders/{id}/proof
  -> Validates non-empty before & after photos and SHA-256 hashes.
  -> Anti-Photo-Reuse check: beforeHash != afterHash.
  -> Enforces meter seal number for Metering tasks.
  -> Enforces PTW permit if ptw_required is true.
  -> Calls Python AI Watchdog (:8001 /internal/ai/verify-gps):
     - Computes Haversine distance between device GPS and registered asset GPS.
     - Compares against unit tolerance (Feeder: 150m, Transformer: 75m, Meter: 50m).
  -> Persists ProofOfWorkLog and transitions Work Order to PROOF_SUBMITTED.

[ Step 3: Unit Head Verification & Customer Care Reversion ]
Unit Head reviews proof log, photos, and AI risk score:
  POST /api/v1/services/work-orders/{id}/verify
  -> If Approved:
     - Work Order status transitions to ATTENDED_VERIFIED.
     - Linked Customer Care ticket automatically updates to ATTENDED_TO with repair summary.
     - Customer Care Officer now receives ticket on their confirmation closure queue!
  -> If Rejected:
     - Work Order status transitions to REJECTED_REWORK with remedial instructions.
     - Field crew must remediate and re-upload evidence.
```

---

## 4. Concrete cURL Examples

### 1. Dispatch Maintenance Work Order (Unit Head)
```bash
curl -X POST http://localhost:8080/api/v1/services/work-orders \
  -H "Authorization: Bearer <TOKEN_UNIT_HEAD>" \
  -H "Content-Type: application/json" \
  -d '{
    "ticket_id": "00000000-0000-0000-0000-000000000001",
    "unit_type": "TRANSFORMER_MAINT",
    "work_order_type": "CORRECTIVE_REPAIR",
    "priority": "HIGH",
    "team_lead_id": "00000000-0000-0000-0000-000000000002",
    "asset_id": "00000000-0000-0000-0000-000000000003",
    "ptw_required": false,
    "instructions": "Replace blown 200A HRC fuse on red phase of Tundun Wada 500kVA DT."
  }'
```

### 2. Submit Mobile PWA Field Proof of Work (Team Lead)
```bash
curl -X POST http://localhost:8080/api/v1/services/work-orders/<wo-uuid>/proof \
  -H "Authorization: Bearer <TOKEN_TEAM_LEAD>" \
  -H "Content-Type: application/json" \
  -d '{
    "latitude": 9.9321,
    "longitude": 8.8920,
    "before_photo_url": "https://storage.nesco.ng/proofs/20261022_before.jpg",
    "after_photo_url": "https://storage.nesco.ng/proofs/20261022_after.jpg",
    "photo_hash_before": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "photo_hash_after": "ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb",
    "work_details": "Replaced blown 200A fuse, tightened loose busbar connection, and tested voltage at 230V."
  }'
```

### 3. Verify Proof of Work & Revert to Customer Care (Unit Head)
```bash
curl -X POST http://localhost:8080/api/v1/services/work-orders/<wo-uuid>/verify \
  -H "Authorization: Bearer <TOKEN_UNIT_HEAD>" \
  -H "Content-Type: application/json" \
  -d '{
    "approved": true,
    "notes": "Verified fuse replacement and terminal voltage balanced across all 3 phases."
  }'
```

---

## 5. Security & Domain Invariants
1. **Air-Gapped AI Watchdog Graceful Degradation:** If the Python AI sidecar is restarting or unreachable, the Go workflow service applies local spherical Haversine calculations directly without failing the request, preserving field team availability in remote transmission corridors.
2. **Actor Authorization:** Only the assigned `team_lead_id` can upload proof for a work order. Non-assigned technicians or impersonators receive HTTP 422.
3. **Closed-Loop Ticket Integration:** Tickets are strictly transitioned to `ATTENDED_TO` only after Unit Head sign-off, ensuring that Customer Care never calls a customer prematurely before engineering verification.

