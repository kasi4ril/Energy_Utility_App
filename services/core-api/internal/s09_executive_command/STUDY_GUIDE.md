# Sprint 9 Study Guide: S09 Executive Command & Predictive Maintenance Engine

## 1. Architectural Scope & Bounded Context
The **Executive Command & Predictive Maintenance Engine** (`s09_executive_command`) forms the strategic nerve center of the NESCO utility platform. It aggregates live operational telemetry across generation, central dispatch, commercial billing, and internal audit into a unified C-Suite dashboard, enforces digital sign-off queues with cryptographic segregation of duties, and executes condition-based maintenance (CBM) prognostics that autonomously raise work orders before catastrophic equipment failures occur.

```
+----------------------------------------------------------------------------------------------------+
|                                    S09 Executive Command Center                                    |
|                                                                                                    |
|  [C-Suite Cross-Departmental Pulse Aggregation]                                                    |
|   ├── S05 Generation Master: Live hydro/thermal MW generation logs                                 |
|   ├── S06 Energy Dispatch: Bulk TCN intake MW, system load, feeder service states                 |
|   ├── S08 Commercial Ops: Revenue inflow, prepaid share %, open theft/churn alerts                 |
|   └── S07 Stores & Audit: Pending MD procurements, unreturned scrap items                          |
|            │                                                                                       |
|            ▼                                                                                       |
|  [Executive Operational Telemetry Snapshots]                                                       |
|   ├── Shift-window captures (MORNING, AFTERNOON, NIGHT, SUMMARY_24H)                              |
|   └── Immutable audit snapshots for board reporting and regulatory filings                         |
|            │                                                                                       |
|  [One-Click Digital Sign-Off Queue (Role-Segregated Governance)]                                    |
|   ├── MD: High-value store procurements & capital asset requisitions                               |
|   ├── HEAD_COMMERCIAL_OPS: Customer care disputed billing reversals & tariff credit sign-offs     |
|   ├── HEAD_INTERNAL_AUDIT: Staff reimbursement & claim release concurrences                        |
|   ├── HEAD_SERVICES: High-risk maintenance Safety Permit-to-Work (PTW) authorizations              |
|   └── GM_GENERATION: Hydro plant turbine overhaul & station spare requisitions                     |
|            │                                                                                       |
|  [Predictive Maintenance Prognostic Center]                                                        |
|   ├── 1. Hydro Reservoir Availability: Head forecast, discharge vs inflow, MWh potential           |
|   ├── 2. Generator Mechanical Health: ISO 10816 vibration, winding temp, oil pressure             |
|   ├── 3. Transformer CBM: Overload %, top oil temperature, repetitive fuse trips                   |
|   └── 4. Feeder Patrol Risk: Fault trip frequency, line density, storm forecast                    |
|            └── Auto-Work-Order Generator -> Dispatches Work Order directly into S04 Services      |
+----------------------------------------------------------------------------------------------------+
```

---

## 2. Core Domain Invariants & Physics Models

### 1. Hydro Water Availability & Generation Prognosis
Predicts available water head and power capacity over a forecast window $D$:
$$\Delta H = (Q_{\text{in}} - Q_{\text{dis}}) \times D \times 0.0000864$$
$$H_{\text{forecast}} = \max(0, H_{\text{current}} + \Delta H)$$
$$P_{\text{avail}} = \frac{9.81 \times Q_{\text{dis}} \times H_{\text{forecast}} \times \eta}{1000} \text{ MW}$$
$$\text{Forecasted MWh} = P_{\text{avail}} \times D \times 24$$
- **Risk Evaluation:**
  - If $H_{\text{forecast}} < 12.0\text{ m}$: Critical water depletion alert; risk score = 95.0.
  - If $\Delta H < -3.0\text{ m}$: High reservoir drawdown warning; risk score = 75.0.

### 2. Turbine Generator Mechanical & Thermal Health
Evaluates mechanical vibration (ISO 10816 criteria), electrical winding temperature, and lubrication pressure:
- **Vibration Severity:** $V \ge 7.1\text{ mm/s}$ (Zone D - Danger) or $T_{\text{winding}} \ge 95^\circ\text{C}$ or $P_{\text{oil}} < 1.5\text{ bar}$:
  - **Risk:** CRITICAL (Score = 95.0).
  - **Auto Action:** Immediately generates a `WOPreventiveMaintenance` work order in S04 flagged with `PriorityCritical` and `PTWRequired = true`.
- **Preventive Maintenance Interval:** Run-hours since overhaul $> 8000\text{ hrs}$:
  - **Risk:** HIGH (Score = 75.0) -> Auto work order for scheduled turbine overhaul.

### 3. Distribution Transformer Condition-Based Maintenance (CBM)
Monitors electrical loading and dielectric thermal envelope:
- **Overload Criteria:** $\text{Load} \ge 120\%$ or $T_{\text{top\_oil}} \ge 90^\circ\text{C}$ or $\text{Fuses Blows (30d)} \ge 4$:
  - **Risk:** CRITICAL (Score = 95.0).
  - **Auto Action:** Dispatches a `WODTInspection` emergency work order to S04 transformer crews.
- **Degradation Criteria:** $\text{Load} \ge 100\%$ or $T_{\text{top\_oil}} \ge 80^\circ\text{C}$ or $\text{Fuses Blows} \ge 2$:
  - **Risk:** HIGH (Score = 75.0) -> Work order for dielectric oil breakdown voltage (BDV) sampling.

### 4. Feeder Reliability & Predictive Vegetation Patrol
Calculates feeder trip vulnerability from transient fault records and vegetation exposure:
$$\text{Risk Score} = (\text{Trips}_{30\text{d}} \times 5.0) + (\text{Density}_{\text{fault/km}} \times 15.0) + (\text{Days}_{\text{patrol}} \times 0.5) + \text{StormBonus}(20.0)$$
- If score $\ge 70.0$: **CRITICAL** -> Dispatches high-priority bush-clearing and line patrol work order into S04.
- If score $\ge 50.0$: **HIGH** -> Schedules preventive foot patrol.

### 5. C-Suite Segregation of Duties Gate
Action items enforce strict cryptographic segregation of duties. An executive can only digitally sign off action items matching their statutory role:
- `RoleMD`: Can only decide `ActionMDProcurementApproval`.
- `RoleHeadCommercialOps`: Can only decide `ActionHeadCommBillingAdjustment`.
- `RoleHeadInternalAudit`: Can only decide `ActionHeadAuditClaimsConsent`.
- `RoleHeadServices`: Can only decide `ActionHeadServicesPTW`.
- `RoleGMGeneration`: Can only decide `ActionGMGenMaterialReq`.
Attempts to cross-approve or sign unauthorized categories trigger `403 Forbidden` (`ErrUnauthorizedExecutiveRole`). Double decisions on non-pending items trigger `409 Conflict` (`ErrIllegalActionItemState`).

---

## 3. Hexagonal Architecture Implementation

```
services/core-api/internal/s09_executive_command/
├── domain/
│   ├── models.go                    # 171 lines: Executive roles, snapshots, triggers, action items
│   ├── rules.go                     # 258 lines: Hydro, generator, transformer, feeder physics & RBAC validation
│   ├── rules_test.go                # 84 lines: Positive & negative unit tests for action gates
│   └── rules_predictive_test.go     # 222 lines: Boundary & threshold tests for 4 predictive models
├── repository/
│   ├── repository.go                # 194 lines: Repository interface & memory snapshot/trigger persistence
│   ├── memory_repo_actions.go       # 79 lines: In-memory executive action items persistence
│   ├── postgres_repo.go             # 321 lines: PostgreSQL storage for snapshots and predictive triggers
│   ├── postgres_repo_actions.go     # 168 lines: PostgreSQL storage for executive action items
│   ├── repository_test.go           # 179 lines: Repository test suite for snapshots & triggers
│   └── repository_actions_test.go   # 113 lines: Repository test suite for action console
├── service/
│   ├── workflow_service.go          # 201 lines: Live dashboard summary & telemetry snapshot captures
│   ├── workflow_service_predictive.go # 320 lines: 4 predictive maintenance evaluations & auto S04 WO creation
│   ├── workflow_service_actions.go  # 99 lines: Action queueing, digital sign-off gates & trigger acks
│   ├── workflow_service_test.go     # 195 lines: Unit tests for dashboard aggregation & snapshots
│   └── workflow_service_predictive_test.go # 191 lines: Unit tests for predictive CBM & auto work orders
└── transport/
    ├── dto.go                       # 178 lines: DTO contracts & input validation routines
    ├── dto_test.go                  # 148 lines: DTO validation test cases
    ├── http_handlers.go             # 241 lines: Dashboard, snapshots, and action console REST handlers
    ├── http_handlers_predictive.go  # 177 lines: Predictive prognostics and trigger acknowledgement endpoints
    ├── http_handlers_test.go        # 326 lines: Integration test suite for dashboard & action items
    └── http_handlers_predictive_test.go # 316 lines: Integration test suite for predictive maintenance
```

---

## 4. REST API Endpoint Specifications

| Method | Endpoint | Description | Auth / Role Gate |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/executive/dashboard` | Aggregates live cross-departmental pulse | Bearer JWT (C-Suite) |
| `POST` | `/api/v1/executive/snapshots` | Captures operational telemetry snapshot | Bearer JWT |
| `GET` | `/api/v1/executive/snapshots/latest` | Retrieves latest operational snapshot | Bearer JWT |
| `GET` | `/api/v1/executive/snapshots` | Lists historical telemetry snapshots | Bearer JWT |
| `POST` | `/api/v1/executive/actions` | Enqueues pending C-Suite action item | Bearer JWT |
| `POST` | `/api/v1/executive/actions/{id}/decide` | Digital approval/rejection sign-off | Bearer JWT (Strict Role Match) |
| `GET` | `/api/v1/executive/actions` | Lists action items by role and status | Bearer JWT |
| `POST` | `/api/v1/executive/predictive/hydro` | Evaluates hydro availability prognosis | Bearer JWT |
| `POST` | `/api/v1/executive/predictive/generator` | Evaluates turbine generator health | Bearer JWT |
| `POST` | `/api/v1/executive/predictive/transformer`| Evaluates transformer CBM metrics | Bearer JWT |
| `POST` | `/api/v1/executive/predictive/feeder` | Evaluates feeder trip probability | Bearer JWT |
| `POST` | `/api/v1/executive/predictive/triggers/{id}/ack` | Acknowledges predictive alarm | Bearer JWT |
| `GET` | `/api/v1/executive/predictive/triggers` | Lists predictive maintenance alarms | Bearer JWT |

---

## 5. Security & Verification Assurance
1. **Reversible Schema Migrations:** `migrations/009_s09_executive_command.up.sql` and `migrations/009_s09_executive_command.down.sql` execute cleanly and reversibly.
2. **File Length Discipline:** Every single Go file in `s09_executive_command` strictly adheres to $\le 350$ lines (highest file: 326 lines).
3. **Statement Test Coverage:**
   - `domain`: **92.6%**
   - `service`: **87.4%**
   - `transport`: **83.6%**
   - All exceeding the strict $\ge 80\%$ standard.
4. **Air-Gapped Sovereign AI:** All physics calculations and prognostic thresholds execute locally on-premise without external telemetry leakage.

