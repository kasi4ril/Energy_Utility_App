# Sprint 7 Study Guide: S07 Stores, Internal Audit Gate & MD Procurement

## 1. Architectural Scope & Bounded Context
The **Stores, Inventory & Internal Audit** bounded context coordinates warehouse stock tracking, live `@GridAssistant` inventory inquiries, dual-path material requisitions, reverse logistics scrap returns, and the automated 3-way staff financial allowance claims verification desk.

```
+----------------------------------------------------------------------------------------------------+
|                                      S07 Stores & Internal Audit                                   |
|                                                                                                    |
|  [Staff / Field Engineers] <--- `@GridAssistant` Chat ---> [Warehouse Inventory Ledger]            |
|                                                              - SKU Balances, Reorder Thresholds    |
|                                                              - Row-Level HMAC Tamper Proofing      |
|                                                                                                    |
|  [Dual-Path Material Requisitions]                                                                 |
|   ├── Path A (In-Stock):                                                                           |
|   │     Unit Head -> Dept/GM Approval -> Audit Consent -> Store Officer Issuance & Digital Waybill|
|   └── Path B (Out-of-Stock / Critical Supplier Purchase):                                          |
|         Unit Head -> Dept/GM Approval -> Audit Pricing Consent -> MD Approval -> Vendor Dispatch  |
|                                                                                                    |
|  [Reverse Logistics Scrap Return]                                                                  |
|   - Replaced/damaged cables, meters, transformer parts returned to Store Scrap Yard                |
|   - Unreturned scrap blocks future material requests                                               |
|                                                                                                    |
|  [Staff Allowance Claims Desk (3-Way Automated Cross-Check)]                                       |
|   - Cross-check 1: Shift Duty Roster / Call Sheets                                                 |
|   - Cross-check 2: High-Voltage Safety PTW Crew Log                                                |
|   - Cross-check 3: Mobile PWA GPS Proof-of-Work Distance Tolerance                                 |
|   - Fraud Scoring (0.00-1.00) & Auditor Sign-Off Desk                                              |
+----------------------------------------------------------------------------------------------------+
```

---

## 2. Core Domain Invariants & Business Logic

### Live `@GridAssistant` Stock Checks & Store Officer Discipline
- Every SKU master record tracks `quantity_available`, `reorder_threshold`, `unit_cost`, and `bin_location`.
- Whenever stock drops to or below `reorder_threshold`, the event `stores.reorder_threshold_breached` is published for procurement.
- Every ledger update re-computes a row-level cryptographic HMAC signature (`crypto.ComputeRowHMAC`).

### Dual-Path Material Requisition Workflow
1. **Path A (In-Stock Requisition):**
   - **Step 1 (`SUBMITTED`):** Unit Head raises material request with required line items.
   - **Step 2 (`DEPT_APPROVED`):** Departmental Head (or Generation GM for Plants 1–6) approves.
   - **Step 3 (`AUDIT_CONSENTED`):** Internal Audit verifies linked work order and checks that no prior scrap returns are delinquent.
   - **Step 4 (`ISSUED`):** Store Officer issues materials, generates digital `WAYBILL-YYYY-XXXX`, and automatically decrements inventory balances.
2. **Path B (Out-of-Stock / Critical Supplier Procurement):**
   - When requested quantity exceeds warehouse balance, `in_stock_flag = false` and status transitions to `MD_PROCUREMENT_PENDING`.
   - **Internal Audit Consent:** Auditor audits supplier quotes and pricing compliance $\rightarrow$ `audit_consent_status = CONSENTED`.
   - **MD Procurement Approval:** Managing Director (MD) reviews the audit trail and authorizes purchase order $\rightarrow$ status `MD_APPROVED`.
   - Vendor materials received $\rightarrow$ Store Officer updates stock balance and generates waybill.

### Reverse Logistics (Scrap Return) Tracker
- High-value materials (cables, copper windings, meters, transformers) require damaged/replaced scrap to be logged into the Store Scrap Yard (`stores_audit.scrap_returns`).
- Logging scrap sets `scrap_returned_flag = true` on the parent requisition, clearing audit holds for subsequent dispatches.

### Staff Allowance Claims Desk (Automated 3-Way Cross-Check)
Staff financial claims (Overtime, Travel Out-of-Station, Hazard) undergo automated 3-way validation:
1. **Call-Collated Shift Roster Check:** Verified against Generation/Dispatch shift roll-calls. Mismatch triggers high-risk fraud flag (`FLAGGED_FRAUD`, score `0.950`).
2. **Safety PTW Crew Log Check:** If linked to a PTW, verifies claimant was named in isolated crew.
3. **Field PWA GPS Proof-of-Work Check:** Verifies work order coordinates match claim date.
- Clean 3-way match: `VERIFIED` (score `0.050`), status `AI_VERIFIED`.
- Mismatch: `FLAGGED_MISMATCH` (score `0.720`), status `AI_FLAGGED`.
- Formal decision by Internal Auditor: `AUDIT_APPROVED` or `AUDIT_REJECTED`.

---

## 3. Hexagonal Symmetry & File Structure

```
services/core-api/internal/s07_stores_audit/
├── domain/
│   ├── models.go                    # StoreItem, MaterialRequisition, ScrapReturn, StaffAllowanceClaim
│   ├── rules.go                     # Dual-path guards, requisition cost calculation, 3-way claim rules
│   └── rules_test.go                # 84.6% pure domain statement coverage
├── repository/
│   ├── repository.go                # Repository interface & thread-safe MemoryRepository
│   ├── postgres_repo.go             # PostgreSQL store inventory and requisition queries
│   ├── postgres_repo_scan.go        # PostgreSQL row scanners
│   ├── postgres_repo_scrap.go       # PostgreSQL scrap return queries
│   ├── postgres_repo_claims.go      # PostgreSQL allowance claims queries
│   └── repository_test.go           # MemoryRepository unit test suite
├── service/
│   ├── workflow_service.go          # Inventory stock, dual-path requisitions, waybills, scrap returns
│   ├── workflow_service_claims.go   # 3-way allowance claims verification & auditor review
│   ├── workflow_service_test.go     # Inventory & requisition unit tests
│   └── workflow_service_claims_test.go # 3-way claims & fraud scoring unit tests
└── transport/
    ├── dto.go                       # DTO structs & Validate() methods
    ├── dto_test.go                  # DTO validation tests
    ├── http_handlers.go             # REST endpoints for inventory & requisitions
    ├── http_handlers_scrap.go       # REST endpoints for reverse logistics scrap returns
    ├── http_handlers_claims.go      # REST endpoints for staff claims & audit review
    ├── http_handlers_test.go        # Inventory & requisition HTTP integration tests
    ├── http_handlers_requisitions_test.go # Path B MD procurement integration tests
    └── http_handlers_claims_test.go # Staff claims HTTP integration tests
```

---

## 4. REST API Endpoint Reference

| Method | Endpoint | Auth | Purpose |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/stores/items` | Bearer JWT | Register warehouse inventory SKU |
| `GET` | `/api/v1/stores/items` | Bearer JWT | List items (`?category=...&low_stock=true`) |
| `GET` | `/api/v1/stores/items/{code}` | Bearer JWT | Live `@GridAssistant` stock balance check |
| `POST` | `/api/v1/stores/items/{id}/stock` | Bearer JWT | Adjust inventory quantity balance |
| `POST` | `/api/v1/stores/requisitions` | Bearer JWT | Unit Head initiates material requisition |
| `GET` | `/api/v1/stores/requisitions` | Bearer JWT | List requisitions |
| `GET` | `/api/v1/stores/requisitions/{id}` | Bearer JWT | Get requisition details |
| `POST` | `/api/v1/stores/requisitions/{id}/dept-approve` | Bearer JWT | Dept Head / Generation GM approval |
| `POST` | `/api/v1/stores/requisitions/{id}/audit-consent` | Bearer JWT | Internal Audit compliance consent |
| `POST` | `/api/v1/stores/requisitions/{id}/md-approve` | Bearer JWT | Managing Director procurement approval |
| `POST` | `/api/v1/stores/requisitions/{id}/issue` | Bearer JWT | Store Officer material issuance & waybill |
| `POST` | `/api/v1/stores/scrap-returns` | Bearer JWT | Log returned scrap into scrap yard |
| `GET` | `/api/v1/stores/scrap-returns` | Bearer JWT | List scrap return records |
| `POST` | `/api/v1/audit/claims` | Bearer JWT | Staff submits allowance claim (3-way check) |
| `GET` | `/api/v1/audit/claims` | Bearer JWT | Query claims (`?status=...&staff_id=...`) |
| `GET` | `/api/v1/audit/claims/{id}` | Bearer JWT | Get claim details and AI audit score |
| `POST` | `/api/v1/audit/claims/{id}/review` | Bearer JWT | Internal Auditor formal approval / query |

---

## 5. Verification & Test Execution

Run test suites for Sprint 7:
```bash
# Domain tests
go test -v -cover ./internal/s07_stores_audit/domain/...

# Service tests
go test -v -cover ./internal/s07_stores_audit/service/...

# Transport tests
go test -v -cover ./internal/s07_stores_audit/transport/...

# Verify file line counts (strict <= 350 lines limit)
wc -l internal/s07_stores_audit/**/*.go
```

