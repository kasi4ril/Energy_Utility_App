# Sprint 1 Study Guide: Foundation, Grid Topology & WorkNotes Engine (`s01_foundation`)

## 1. Purpose & Utility Domain Context
The **Foundation Bounded Context (`s01_foundation`)** is the foundational bedrock of the NESCO Enterprise Utility Platform. In an integrated electricity generation and distribution utility operating six remote hydroelectric plants and regional 33kV/11kV distribution corridors:
- **Organizational Truth:** All personnel, departments, and physical substations must adhere to a strict hierarchical tree (`HQ` -> `Districts` -> `Service Units` and `Master Station` -> `Plants 1–6`).
- **Topological Integrity:** Downstream assets cannot exist in isolation. Distribution Transformers (DTs) must feed from designated 11kV Feeders, which must route from 33kV Substations, fed by Generating Plants.
- **Cryptographic Auditability:** Inter-departmental coordination (between Generation, Dispatch, Customer Care, Maintenance, and Stores) cannot rely on mutable rows or ephemeral chats. The `WorkNotes` engine provides an append-only, SHA-256 chained ledger where every operational note is cryptographically bound to its parent entity and previous note.

---

## 2. File-by-File Architectural Breakdown

```text
services/core-api/internal/s01_foundation/
├── domain/
│   ├── models.go       # Structs & Enums: OrgUnit, StaffProfile, GridAsset, WorkNote, GPSPoint
│   ├── rules.go        # Pure business rules: hierarchy checks, grid topology validation, GPS boundaries
│   └── rules_test.go   # Comprehensive domain unit tests (> 98% coverage)
├── repository/
│   ├── repository.go                  # Repository interface & thread-safe in-memory mock implementation
│   ├── postgres_repo.go               # PostgreSQL + PostGIS queries for OrgUnits & StaffProfiles
│   ├── postgres_repo_assets_notes.go  # PostgreSQL + PostGIS queries for GridAssets & WorkNotes
│   └── repository_test.go             # Storage layer unit tests
├── service/
│   ├── workflow_service.go      # Workflow orchestrator: signs HMAC, hashes passwords, computes SHA-256
│   └── workflow_service_test.go # Orchestration tests, negative boundary & tamper detection tests
└── transport/
    ├── dto.go                 # REST request/response validation structs with Validate() methods
    ├── http_handlers.go       # net/http REST endpoints & RBAC middleware wiring
    └── http_handlers_test.go  # End-to-end HTTP integration tests
```

### Shared Packages:
- `pkg/crypto/hmac_chain.go`: SHA-256 WorkNotes hash chaining (`ComputeWorkNoteHash`), HMAC row signature generation/verification (`ComputeRowHMAC`), and RFC 4122 v4 UUID generator.
- `pkg/auth/jwt_middleware.go`: Cryptographic JWT generation, validation, 15-minute token expiry, and RBAC middleware (`RequireAuth`, `RequireRoles`, `IsExecutiveRole`).
- `pkg/events/nats_bus.go`: Thread-safe asynchronous event bus with wildcard matching (`*.transferred`, `grid.*`).
- `pkg/database/postgres.go`: PostgreSQL connection pool wrapper and DSN constructor.

---

## 3. Step-by-Step Data & Execution Flow

```
[ HTTP Client / Web Portal ]
           │
           │  POST /api/v1/notes/GRID_ASSET/{assetId}
           ▼
[ Transport Layer: http_handlers.go ]
   - Extracts URL path variables: entityType="GRID_ASSET", entityId
   - Decodes JSON payload into AppendWorkNoteRequest
   - Runs req.Validate()
           │
           ▼
[ Service Layer: workflow_service.go ]
   - Validates author_staff_id exists in repository
   - Queries repository for latest note on this entity
   - Extracts prevHash (GenesisHash if first note in chain)
   - Computes current_hash = SHA256(prevHash + entityType + entityId + author + noteText + ts)
   - Validates WorkNote domain integrity
           │
           ▼
[ Repository Layer: postgres_repo_assets_notes.go ]
   - Executes INSERT into foundation.work_notes
   - Database trigger trg_immutable_work_notes guarantees append-only immutability
           │
           ▼
[ Asynchronous Event Bus: nats_bus.go ]
   - Publishes "foundation.work_note.appended" event
   - Dispatches non-blockingly to listening departments / Ollama summary sidecar
```

---

## 4. Key Engineering Concepts for Personal Study

1. **SHA-256 Prev-Hash Chaining (Blockchain-Lite in Enterprise SQL):**
   ```go
   currentHash = SHA256(prevHash + ":" + entityType + ":" + entityID + ":" + authorID + ":" + noteText + ":" + timestamp)
   ```
   If any past note in the database is modified, its stored `current_hash` and all subsequent `prev_hash` references fail validation immediately, exposing unauthorized database tampering.

2. **Row HMAC Signatures:**
   Sensitive records (`foundation.staff_profiles`, `foundation.grid_assets`) compute an HMAC-SHA256 signature using a server-side pepper secret. Even if an attacker gains raw SQL access, modifying an employee's role from `STORE_OFFICER` to `MD` invalidates the row signature, blocking authentication and raising an audit alarm.

3. **PostGIS Spatial Geometry:**
   Grid assets and org units store coordinates using `GEOMETRY(Point, 4326)` indexed with GiST (`USING GIST(gps_point)`). This enables lightning-fast spatial bounding-box searches and distance queries across transmission corridors.

4. **Clean Hexagonal Architecture:**
   The `domain` package has zero external dependencies (no SQL, no HTTP, no JSON). The `service` layer coordinates use cases, and `transport` handles protocol translation. This ensures business rules can be verified with high-speed unit tests.

---

## 5. How to Run & Verify Locally

### Run Unit & Integration Tests:
```bash
cd services/core-api
go test -v -cover ./...
```

### Run Static Analysis (`go vet`):
```bash
cd services/core-api
go vet ./...
```

### Start the Core API Server:
```bash
cd services/core-api
go run cmd/server/main.go
```

### Sample Verification via cURL:

1. **Health Check:**
   ```bash
   curl -s http://localhost:8080/healthz
   ```

2. **Login Handshake:**
   ```bash
   curl -s -X POST http://localhost:8080/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"admin@nesco.ng","password":"AdminPassword#123"}'
   ```

3. **Append Chained WorkNote:**
   ```bash
   curl -s -X POST http://localhost:8080/api/v1/notes/GRID_ASSET/11111111-2222-3333-4444-555555555555 \
     -H "Authorization: Bearer <TOKEN>" \
     -H "Content-Type: application/json" \
     -d '{"author_staff_id":"...","note_text":"Feeder inspection completed without alarms."}'
   ```

