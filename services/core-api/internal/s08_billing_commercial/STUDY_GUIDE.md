# Sprint 8 Study Guide: S08 Billing, Commercial Operations & Revenue Protection AI

## 1. Architectural Scope & Bounded Context
The **Billing, Commercial Operations & Revenue Protection AI** bounded context coordinates the utility's commercial lifecycle, eliminating estimated billing through pure metered vending, providing automated customer meter-purchase refund ledgers, implementing a strict Head of Commercial Operations gate for wrong billing disputes that re-engages Customer Care, and executing local AI revenue protection scans for theft/bypass and off-grid churn risks.

```
+----------------------------------------------------------------------------------------------------+
|                                    S08 Billing & Commercial Operations                             |
|                                                                                                    |
|  [Payment Hub / Multi-Channel Ingestion]                                                           |
|   ├── Online Web Portal (Paystack, Flutterwave)                                                    |
|   ├── Fintech Agency Networks (OPay, Moniepoint, PalmPay, Paga)                                    |
|   └── Commercial Bank Apps / USSD (Zenith, GTBank, FirstBank)                                      |
|            │                                                                                       |
|            ▼                                                                                       |
|  [Zero-Estimation Metered Billing Engine]                                                           |
|   ├── Formula: kWh = (Amount Paid + Meter Refund Credited) / Tariff Rate Per kWh                   |
|   └── 20-Digit STS Token Derivation: XXXX-XXXX-XXXX-XXXX-XXXX                                      |
|            │                                                                                       |
|            ├── Active Meter Refund Ledger Check                                                    |
|            │    └── 10% Vend Amortization: Installment Credited -> Balance Decremented            |
|            │                                                                                       |
|  [Wrong Billing Resolution Closed Loop]                                                            |
|   ├── Step 1: Customer Care logs WRONG_BILLING complaint ticket                                    |
|   ├── Step 2: Billing Unit audits tariffs/multipliers -> submits Adjustment Proposal               |
|   ├── Step 3: Head of Commercial Operations reviews audit evidence & signs off                     |
|   └── Step 4: System automatically reverts ticket to ATTENDED_TO for Customer Care closure call    |
|                                                                                                    |
|  [Revenue Protection Local AI Watchdog]                                                            |
|   ├── Model A: Energy Theft / Bypass Detection                                                     |
|   │     - Feeder energized + vending drop >= 50% OR DT bulk vs metered loss >= 30%                 |
|   └── Model B: Customer Churn / Off-Grid Solar Migration Risk                                      |
|         - High outage count (> 10) + sustained vending collapse >= 40%                             |
+----------------------------------------------------------------------------------------------------+
```

---

## 2. Core Domain Invariants & Business Logic

### Zero-Estimation Metered Billing
Estimated billing is architecturally impossible in the platform. Vended energy kilowatt-hours (kWh) is mathematically derived strictly from tariff rates and net credited values:
$$\text{kWh} = \frac{\text{Amount Paid} + \text{Meter Refund Credited}}{\text{Tariff Rate Per kWh}}$$
Every prepaid transaction generates a 20-digit Standard Transfer Specification (STS) token formatted into 5 blocks of 4 digits: `XXXX-XXXX-XXXX-XXXX-XXXX`.

### Customer Meter-Purchase Refund Amortization Ledger
In accordance with Nigerian Electricity Regulatory Commission (NERC) regulations, customers who finance their meter installation receive periodic amortization reimbursements:
1. When registered, the ledger tracks `initial_meter_cost`, `refund_percentage_per_tx` (defaulting to 10%), and `remaining_refund_balance`.
2. On every vending transaction:
   $$\text{Installment} = \min(\text{Vend Amount} \times \text{Percentage}, \text{Remaining Balance})$$
3. The installment is credited as additional vending energy without out-of-pocket customer expense.
4. When `remaining_refund_balance` reaches zero, the ledger status automatically transitions to `AMORTIZED_COMPLETED`.

### Wrong Billing Commercial Resolution Loop
1. **Ticket Initiation:** Customer Care logs a dispute ticket (`WRONG_BILLING` or `UNCREDITED_PAYMENT`).
2. **Investigation & Proposal:** Billing Officer investigates multiplier errors, meter reading typos, or tariff band mismatches and submits a `BillingAdjustment` proposal.
3. **Head of Commercial Operations Approval Gate:** The Head of Commercial Operations conducts digital sign-off (`ValidateCommercialHeadApproval`).
4. **Closed-Loop Ticket Reversion:** Upon approval, the status transitions to `APPLIED_REVERTED_TO_CARE`, and the platform automatically reverts the linked Customer Care ticket to `ATTENDED_TO` with audit notes. This enforces the invariant that Customer Care officers must call the customer to verify satisfaction before closing the ticket.

### Revenue Protection AI Anomaly Watchdog
1. **Energy Theft Bypass Detector:**
   - Evaluates distribution transformer bulk intake against sum of metered customer sales. Losses exceeding 30% trigger transformer-level bypass flags.
   - Monitors physical customer vending cadence: if a feeder is energized while customer vending drops $\ge 50\%$, an immediate line bypass alarm is raised with high confidence (0.880).
2. **Customer Churn Risk Forecaster:**
   - Identifies commercial and industrial customers facing recurrent feeder outages ($> 10$) combined with vending drops $\ge 40\%$, alerting commercial teams of imminent off-grid solar/generator migration.

---

## 3. REST API Route Registry

| Method | Endpoint | Description | Role / Auth |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/billing/vending/vend` | Execute metered vending payment and generate STS token | Bearer JWT (Staff/Fintech) |
| `GET` | `/api/v1/billing/vending/{id}` | Query vending payment receipt and energy summary | Bearer JWT |
| `GET` | `/api/v1/billing/vending` | List vending transactions by customer/meter | Bearer JWT |
| `POST` | `/api/v1/billing/refunds` | Register customer meter-purchase refund ledger | Bearer JWT (Commercial Admin) |
| `GET` | `/api/v1/billing/refunds/{customer_id}` | Query remaining meter refund balance and history | Bearer JWT |
| `POST` | `/api/v1/billing/adjustments` | Submit proposed billing adjustment for dispute | Bearer JWT (Billing Officer) |
| `POST` | `/api/v1/billing/adjustments/{id}/decide` | Approve/reject adjustment and revert ticket to Care | Bearer JWT (Head Commercial Ops) |
| `GET` | `/api/v1/billing/adjustments/{id}` | Query billing adjustment details | Bearer JWT |
| `GET` | `/api/v1/billing/adjustments` | List billing adjustments by status or customer | Bearer JWT |
| `POST` | `/api/v1/billing/revenue-protection/scan` | Trigger local AI telemetry scan for theft/churn | Bearer JWT (Revenue Protection) |
| `POST` | `/api/v1/billing/revenue-protection/alerts/{id}/review` | Update investigation status and notes on AI alert | Bearer JWT (Revenue Protection) |
| `GET` | `/api/v1/billing/revenue-protection/alerts/{id}` | Query AI alert details and evidence | Bearer JWT |
| `GET` | `/api/v1/billing/revenue-protection/alerts` | List revenue protection AI alerts | Bearer JWT |

---

## 4. Operational Runbook & Disaster Recovery

### Vending Gateway Timeout Recovery
If an external fintech (OPay, Moniepoint) or bank gateway encounters network drops:
1. Re-query using unique idempotency key `tx_reference`.
2. If transaction was already persisted, return the existing STS token and energy allocation without debiting or re-amortizing refund ledgers.
3. Unverified settlements flag transactions for internal audit reconciliation.

### Customer Care Callback Discipline
When a billing adjustment is approved by the Head of Commercial Operations:
1. Customer ticket status updates to `ATTENDED_TO`.
2. Customer Care dashboard receives prompt to conduct outbound verification call.
3. Ticket can only be marked `CLOSED` after customer confirms satisfaction.

