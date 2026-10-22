# Sprint 6 Study Guide: S06 Energy Management & Transmission Dispatch

## 1. Architectural Scope & Bounded Context
The **Energy Management & Transmission Dispatch** bounded context coordinates real-time bulk energy allocation, automated priority feeder load shedding, and human-life-critical high-voltage **Permit-to-Work (PTW)** isolations.

```
+-----------------------------------------------------------------------------------------+
|                               S06 Energy Management                                     |
|                                                                                         |
|  [Google Form B Webhook] ---> [Bulk Intake Telemetry] ---> [Dual-Source Power Balance]  |
|                                  - Hydro Plants 1-6            Generation + TCN (MW)    |
|                                  - TCN Grid Intake                                      |
|                                                                     |                   |
|                                                                     v                   |
|  [5-Step Safety PTW Protocol] <---> [Feeder Status Ledger] <--- [Priority Rationing]    |
|   1. Request (Unit Head)             - FDR-IND-01 (Band A)       Rank 1 to 5 Allocation |
|   2. Dual Approval (Dept + Safety)   - FDR-RES-01 (Band D)       Load Shedding Guard    |
|   3. Substation Isolation (Operator) - Live State (PTW/Shed)                            |
|   4. Execution & Monitoring                                                             |
|   5. Clearance & Restoration                                                            |
+-----------------------------------------------------------------------------------------+
```

---

## 2. Core Domain Invariants & Physical Rules

### Dual-Source Bulk MW Intake
- **Formula:** $\text{Bulk Intake (MW)} = \sum_{p=1}^6 \text{Plant Generation MW} + \text{TCN Bulk Intake MW}$
- Physical Sanity: Bulk generation and TCN intake must be non-negative ($\ge 0.0\text{ MW}$).
- TCN Intake is initialized to $0.00\text{ MW}$ until grid interconnection is commissioned.

### Priority Feeder Rationing Algorithm
When total electrical demand exceeds bulk power intake:
1. Feeders are ranked by priority (Priority Rank 1: Band A / High Revenue down to Rank 5: Band E).
2. The scheduler iterates through feeders in strict rank ascending order:
   - Feeders within cumulative intake capacity remain energized (`IN_SERVICE`).
   - Once cumulative demand exceeds available MW, subsequent feeders transition to `OUT_ON_LOAD_SHEDDING`.
   - Feeders currently under safety isolation (`OUT_ON_PERMIT_PTW`) or tripped on fault (`OUT_ON_FAULT`) are preserved in their isolated states and cannot be energized.

### The 5-Step Safety Permit-to-Work (PTW) Protocol
High-voltage line maintenance is strictly guarded by the 5-step life cycle:
1. **Request (Step 1):** Unit Head submits work description and asset UUID $\rightarrow$ Status `REQUESTED`.
2. **Dual Approval (Step 2):** Requires digital approval from **BOTH**:
   - `HEAD_SERVICES` (Head of Services)
   - `SAFETY_REPRESENTATIVE_HSE` (Safety Representative)
   - Status transitions to `APPROVED_PENDING_ISOLATION` only when **both** sign.
3. **Substation Isolation (Step 3):** Substation Operator racks out breaker, applies padlocks, and earths line with minimum 10-character confirmation $\rightarrow$ Status `ISOLATED_EARTHED`.
   - Feeder ledger state shifts automatically to `OUT_ON_PERMIT_PTW`.
4. **Execution (Step 4):** Maintenance proceeds safely under HSE oversight.
5. **Clearance & Restoration (Step 5):**
   - Step 5a: HSE Safety Rep certifies all crew clear and grounds removed $\rightarrow$ Status `SAFETY_CLEARED`.
   - Step 5b: Substation Operator closes breaker and restores feeder $\rightarrow$ Status `RESTORED_CLOSED`.
   - Feeder ledger state reverts automatically to `IN_SERVICE`.

---

## 3. Hexagonal Symmetry & File Structure

```
services/core-api/internal/s06_energy_management/
├── domain/
│   ├── models.go               # FeederState, PTWStatus, DispatchTelemetryLog, SafetyPermitPTW
│   ├── rules.go                # Dual intake calculation, priority rationing, PTW sequence guards
│   └── rules_test.go           # 91.2% pure domain statement coverage
├── repository/
│   ├── repository.go           # Repository interface & thread-safe MemoryRepository
│   ├── postgres_repo.go        # PostgreSQL dispatch logs & feeder status queries
│   ├── postgres_repo_ptw.go    # PostgreSQL 5-step PTW queries & transaction guards
│   └── repository_test.go      # MemoryRepository comprehensive test suite
├── service/
│   ├── workflow_service.go     # Bulk intake, Google Forms B webhook, priority rationing
│   ├── workflow_service_ptw.go # The 5-step Safety PTW state machine coordinator
│   ├── workflow_service_test.go# Telemetry, webhook, and rationing unit tests
│   └── workflow_service_ptw_test.go # PTW 5-step lifecycle & negative path unit tests
└── transport/
    ├── dto.go                  # DTO structs & Validate() methods
    ├── dto_test.go             # DTO validation tests
    ├── http_handlers.go        # REST endpoints & HMAC-SHA256 Google Forms webhook receiver
    ├── http_handlers_ptw.go    # REST endpoints for PTW request, approve, isolate, clear, restore
    ├── http_handlers_test.go   # HTTP integration tests
    └── http_handlers_ptw_test.go # PTW HTTP lifecycle integration tests
```

---

## 4. API Endpoint Reference

| Method | Endpoint | Auth | Purpose |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/energy/bulk-intake` | Bearer JWT | Log hourly generation and TCN intake |
| `GET` | `/api/v1/energy/bulk-intake/latest` | Bearer JWT | Fetch latest recorded intake |
| `GET` | `/api/v1/energy/bulk-intake` | Bearer JWT | Query historical intake logs |
| `POST` | `/api/v1/energy/rationing/calculate` | Bearer JWT | Execute priority feeder rationing |
| `GET` | `/api/v1/energy/feeders` | Bearer JWT | List all feeder states and loads |
| `GET` | `/api/v1/energy/feeders/{assetId}` | Bearer JWT | Get single feeder telemetry |
| `POST` | `/api/v1/energy/webhook/google-forms` | HMAC-SHA256 | Form B ingestion webhook |
| `POST` | `/api/v1/energy/ptw/request` | Bearer JWT | Step 1: Unit Head initiates PTW |
| `POST` | `/api/v1/energy/ptw/{id}/approve` | Bearer JWT | Step 2: Dept Head / HSE dual approval |
| `POST` | `/api/v1/energy/ptw/{id}/isolate` | Bearer JWT | Step 3: Substation Operator racks out & earths |
| `POST` | `/api/v1/energy/ptw/{id}/clear-safety` | Bearer JWT | Step 5a: HSE signs crew clear |
| `POST` | `/api/v1/energy/ptw/{id}/restore` | Bearer JWT | Step 5b: Substation Operator re-energizes line |
| `GET` | `/api/v1/energy/ptw/{id}` | Bearer JWT | Get permit by UUID or Permit No |
| `GET` | `/api/v1/energy/ptw` | Bearer JWT | Filter permits by status and asset ID |

---

## 5. Verification & Test Execution

Run all test suites for Sprint 6:
```bash
# Domain tests
go test -v -cover ./internal/s06_energy_management/domain/...

# Service tests
go test -v -cover ./internal/s06_energy_management/service/...

# Transport tests
go test -v -cover ./internal/s06_energy_management/transport/...

# Verify file line counts (strict <= 350 lines limit)
wc -l internal/s06_energy_management/**/*.go
```

