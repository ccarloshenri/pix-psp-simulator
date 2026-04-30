# pix-psp-simulator

A PIX Payment Service Provider (PSP) simulator implementing the Itaú/Bacen PIX API spec. Covers immediate charges (cob), due-date charges (cobv), payment simulation, and refunds (devolucao).

---

## What it does

- Exposes a PIX API compatible with the Itaú/Bacen cob, cobv, and devolucao specs
- Stores charges and payments in memory (default) or PostgreSQL, switchable via environment variable
- Simulates PIX payment receipt asynchronously: HTTP returns 202 immediately; a background worker marks the charge as CONCLUIDA
- Processes refunds (devolucao) against received payments identified by endToEndId
- Runs with zero external dependencies in memory mode — a single `go run main.go` is enough

---

## Architecture

```
HTTP Request
     │
     ▼
Controller          — parses HTTP, writes JSON response
     │
     ▼
Processor           — validates input, maps to domain objects
     │
     ▼
BO (Business Object) — applies business rules and orchestrates
     │
     ▼
Repository (interface)
     │
     ├── memory/     — 64-shard FNV-hashed RWMutex maps (default)
     └── sql/        — PostgreSQL via lib/pq, auto-migrated on startup
```

```
src/
├── containers/               — DI wiring, config, server setup
└── layers/main/
    ├── bo/                   — business logic per use case
    ├── controller/           — HTTP handlers, JSON serialization
    ├── enums/                — cob/devolucao status constants
    ├── implementations/
    │   ├── channel/          — buffered channel payment queue
    │   ├── memory/           — in-memory sharded repositories
    │   ├── sql/              — PostgreSQL repositories + migrations
    │   ├── uuid/             — ID generator (txid, e2eid, rtrid)
    │   └── worker/           — background payment worker goroutine
    ├── interfaces/           — repository and service contracts
    ├── models/               — cob, cobv, pix, devolucao structs
    └── processor/            — input validation per use case
```

---

## Concurrency model

**Sharded in-memory repositories.** Each repository (cob, cobv, pix) is backed by 64 independent shards. A key is mapped to a shard using FNV-32a hashing. Each shard has its own `sync.RWMutex`, reducing lock contention approximately 64x compared to a single global mutex.

**Async payment simulation.** `POST /cob/simulate` enqueues a job onto a buffered channel (capacity 256) and returns HTTP 202 immediately. A single background goroutine drains the channel, persists the Pix record, and updates the cob/cobv status to CONCLUIDA.

**ID generation.** TxIDs, endToEndIds, and refund IDs are generated with `math/rand.Uint64()`. The global `math/rand` source in Go 1.20+ is goroutine-safe and avoids the syscall overhead of `crypto/rand`.

**Response serialization.** All JSON responses are written via a `sync.Pool` of `bytes.Buffer` instances, reducing allocations and GC pressure under concurrent load.

---

## Endpoints

### Cob — immediate charge

| Method   | Path                              | Description                                      |
|----------|-----------------------------------|--------------------------------------------------|
| `POST`   | `/cob`                            | Create charge with auto-generated txid           |
| `PUT`    | `/cob/{txid}`                     | Create charge with explicit txid                 |
| `GET`    | `/cob`                            | List charges (query: `status`, `inicio`, `fim`)  |
| `GET`    | `/cob/{txid}`                     | Get charge by txid                               |
| `PATCH`  | `/cob/{txid}`                     | Update charge value or expiration                |
| `DELETE` | `/cob/{txid}`                     | Remove charge                                    |

### CobV — due-date charge

| Method   | Path                              | Description                                                           |
|----------|-----------------------------------|-----------------------------------------------------------------------|
| `PUT`    | `/cobv/{txid}`                    | Create due-date charge                                                |
| `GET`    | `/cobv`                           | List due-date charges (query: `status`, `dataDeVencimento`, `inicio`, `fim`) |
| `GET`    | `/cobv/{txid}`                    | Get due-date charge by txid                                           |
| `PATCH`  | `/cobv/{txid}`                    | Update due-date charge                                                |

### Payment simulation and refunds

| Method   | Path                              | Description                                      |
|----------|-----------------------------------|--------------------------------------------------|
| `POST`   | `/cob/simulate`                   | Simulate PIX payment receipt (async, returns 202) |
| `PUT`    | `/cob/{e2eid}/devolucao/{id}`     | Create refund on a received payment              |
| `GET`    | `/cob/{e2eid}/devolucao/{id}`     | Get refund by id                                 |

---

## Local development

### Prerequisites

- Go 1.22+
- Docker and Docker Compose (SQL mode only)

### Memory mode (no external dependencies)

```bash
go run main.go
# Server starts on :8080
```

### SQL mode (PostgreSQL)

```bash
docker-compose up -d
STORAGE=sql DATABASE_URL=postgres://pix:pix@localhost:5432/pix_simulator?sslmode=disable go run main.go
```

Migrations run automatically on startup.

### Tests

```bash
go test ./...
```

### Load test

```bash
go run tests/load/bomber.go
```

The load test warms up with 300 pre-created charges, sweeps concurrency levels (1, 5, 10, 25, 50, 100, 200) with 8-second runs each, selects the best concurrency, and runs a 30-second final mixed-traffic test. It then runs a second 30-second write-focused benchmark (`POST /cob` + `PUT /cob/{txid}` only, 250 goroutines).

---

## Environment variables

| Variable       | Default    | Description                                      |
|----------------|------------|--------------------------------------------------|
| `PORT`         | `8080`     | TCP port the server listens on                   |
| `STORAGE`      | `memory`   | Storage backend: `memory` or `sql`               |
| `DATABASE_URL` | _(empty)_  | PostgreSQL DSN, required when `STORAGE=sql`      |

---

## Performance

Measured on a single machine running both client and server (memory backend).

**Mixed traffic** — 50% GET /cob/{txid}, 15% POST /cob, 15% PUT /cob/{txid}, 10% GET /cob, 10% POST /cob/simulate — 200 goroutines, 30 seconds:

| Throughput | avg    | p99     | Errors |
|------------|--------|---------|--------|
| 623 req/s  | 172 ms | 1079 ms | 0%     |

**Write-focused** — POST /cob + PUT /cob/{txid} only — 250 goroutines, 30 seconds:

| Throughput | avg    | p99    | Errors |
|------------|--------|--------|--------|
| 809 req/s  | 309 ms | 502 ms | 0%     |

---

## Status flows

**Cob / CobV:**
- `ATIVA` → `CONCLUIDA` after a simulated payment is processed by the worker
- `ATIVA` → `REMOVIDA_PELO_USUARIO_RECEBEDOR` after `DELETE /cob/{txid}`

**Devolucao:**
- `EM_PROCESSAMENTO` → `DEVOLVIDO` on success
- `EM_PROCESSAMENTO` → `NAO_REALIZADO` on failure
