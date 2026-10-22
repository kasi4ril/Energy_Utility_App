# Sprint 2 Study Guide: Ollama Local AI Guardian, Cryptographic Ledger & `@GridAssistant` (`s02_ollama_guardian`)

## 1. Purpose & Utility Domain Context
In a large electrical utility like NESCO Nig Ltd, sensitive operations directly impact revenue, legal liability, and grid integrity:
- **Revenue Leakage:** Unauthorized customer tariff downgrades (e.g. Band A to Band C) or manual debt write-offs can silently bleed millions of Naira if staff can unilaterally update balances.
- **Physical Grid Risks:** Unauthorized meter swaps and store stock alterations obscure energy theft and material shrinkage.
- **Air-Gapped AI Sovereignty:** Zero PII, customer billing ledgers, or grid topologies may be sent to third-party cloud LLM APIs. All intelligence, anomaly detection, and chat queries must run locally on-premise via Ollama (`http://127.0.0.1:11434`).

The **`s02_ollama_guardian`** module establishes a **Dual-Brain Local AI Governance Architecture**:
1. **Brain 1 (Watchdog):** Evaluates Maker-Checker fraud risk scores (`0.00–1.00`), enforces strict segregation of duties (Makers cannot check their own work), and maintains an unalterable `immutable_audit_chain` using SHA-256 prev-hash links.
2. **Brain 2 (Co-Pilot):** Powers `@GridAssistant` in community chat, giving field and office staff instant lookups into warehouse inventory, WorkNotes thread summaries, and feeder statuses.

---

## 2. File-by-File Architectural Breakdown

```text
services/core-api/internal/s02_ollama_guardian/
├── domain/
│   ├── models.go       # Structs & Enums: AuditRecord, MakerCheckerRequest, ChatMessage, RiskEvaluation
│   ├── rules.go        # Pure business rules: maker-checker validation, risk threshold guards, audit hash formulas
│   └── rules_test.go   # Comprehensive domain unit tests (> 90% coverage)
├── repository/
│   ├── repository.go      # Repository interface & thread-safe in-memory mock implementation
│   ├── postgres_repo.go   # PostgreSQL operations for security_ai schema
│   └── repository_test.go # Repository test suite
├── service/
│   ├── workflow_service.go      # Service orchestrator: sidecar communication, audit chaining, chat processing
│   └── workflow_service_test.go # Orchestration tests, negative boundary & tamper detection tests
└── transport/
    ├── dto.go                 # Request/response validation structs with Validate() methods
    ├── http_handlers.go       # net/http REST endpoints & RBAC middleware wiring
    └── http_handlers_test.go  # End-to-end HTTP integration tests

services/ai-guardian/
├── app/
│   ├── main.py                         # FastAPI/HTTP entrypoint (:8001) for internal AI calls
│   ├── core/ollama_client.py           # Air-gapped Ollama client with deterministic fallback
│   ├── agents/
│   │   ├── fraud_and_gps_watchdog.py   # Maker-Checker fraud risk scorer & Haversine GPS proximity auditor
│   │   ├── grid_assistant_chat.py      # @GridAssistant in-app store inventory & workflow co-pilot
│   │   └── report_formatter.py         # Standardized shift handover executive report formatter
│   └── tests/
│       └── test_s02_ollama_guardian.py # 100% passing Python test suite (11 unit tests)
```

---

## 3. Step-by-Step Data & Execution Flow

```
[ Maker (e.g. Customer Care / Store Officer) ]
                     │
                     │  POST /api/v1/maker-checker/request
                     ▼
        [ Go API Gateway (:8080) ]
                     │
                     │  HTTP POST /internal/ai/risk-score
                     ▼
     [ Python Local AI Sidecar (:8001) ]
     (Evaluates category, amount delta, justification length)
                     │
                     ▼
          Returns: Risk Score = 0.88
                     │
                     ▼
        [ Go API Gateway (:8080) ]
   - Saves request as PENDING with risk score 0.88
   - Emits "security.fraud_alert.high_risk" event
                     │
                     ▼
[ Checker (e.g. Managing Director / Department Head) ]
                     │
                     │  POST /api/v1/maker-checker/{id}/approve
                     ▼
        [ Go API Gateway (:8080) ]
   - Enforces Maker != Checker
   - Enforces Executive Tier role (required when risk score >= 0.80)
   - Fetches predecessor audit record from security_ai.immutable_audit_chain
   - Computes: current_hash = SHA256(prev_hash + table + id + action + payload + actor + ts)
   - Appends new row to immutable_audit_chain (guaranteed append-only by DB trigger)
   - Emits "security.maker_checker.approved" to execute downstream mutation
```

---

## 4. Key Engineering Concepts for Personal Study

1. **Maker-Checker Dual Authorization Pattern:**
   Critical financial and grid topology changes cannot execute in a single step. The initiator (Maker) proposes the change, while an authorized reviewer (Checker) approves it. The state machine guarantees that `maker_id != checker_id`.

2. **Dynamic Risk-Tiered Escalation:**
   Instead of static approval gates, requests receive an automated AI Fraud Risk Score from `0.00` to `1.00`. Changes with risk `>= 0.80` strictly reject departmental unit heads and demand C-Suite Executive sign-off (`MD`, `HEAD_SERVICES`, `HEAD_COMMERCIAL_OPS`, etc.).

3. **Tamper-Evident Audit Chain:**
   Each audit entry references the `current_hash` of the prior entry as its `prev_hash`:
   ```go
   currentHash = SHA256(prevHash + ":" + tableName + ":" + recordID + ":" + action + ":" + newJSON + ":" + actorID + ":" + timestamp)
   ```
   If a malicious actor directly updates a record in PostgreSQL, verifying the chain continuity immediately detects the cryptographic discontinuity.

4. **Great-Circle Haversine Calculation:**
   To prevent "ghost" field repairs, the watchdog calculates the spherical distance between the Team Lead's GPS device coordinates and the registered grid asset:
   $$\Delta\sigma = 2 \arcsin\left(\sqrt{\sin^2\left(\frac{\Delta\phi}{2}\right) + \cos\phi_1\cos\phi_2\sin^2\left(\frac{\Delta\lambda}{2}\right)}\right)$$
   $$d = R \cdot \Delta\sigma$$

---

## 5. How to Run & Verify Locally

### Run Go Tests:
```bash
cd services/core-api
go test -v -cover ./internal/s02_ollama_guardian/...
```

### Run Python AI Tests:
```bash
cd services/ai-guardian
python3 -m unittest discover tests
```

### Start the AI Guardian Sidecar:
```bash
cd services/ai-guardian
python3 app/main.py
```

### Verify Endpoints via cURL:

1. **Submit a Maker-Checker Request:**
   ```bash
   curl -s -X POST http://localhost:8080/api/v1/maker-checker/request \
     -H "Authorization: Bearer <TOKEN>" \
     -H "Content-Type: application/json" \
     -d '{
       "module": "customercare",
       "change_category": "TARIFF_BAND_CHANGE",
       "target_record_id": "cust-uuid-1",
       "payload_json": "{\"old_band\":\"A\",\"new_band\":\"C\"}",
       "justification": "Customer night closure verified by field lead"
     }'
   ```

2. **Query the Audit Chain:**
   ```bash
   curl -s http://localhost:8080/api/v1/audit/chain/customercare/cust-uuid-1 \
     -H "Authorization: Bearer <TOKEN>"
   ```

3. **Interact with `@GridAssistant` in Chat:**
   ```bash
   curl -s -X POST http://localhost:8080/api/v1/chat/message \
     -H "Authorization: Bearer <TOKEN>" \
     -H "Content-Type: application/json" \
     -d '{
       "channel_id": "general",
       "message_text": "@GridAssistant do we have three-phase meters in store?"
     }'
   ```

