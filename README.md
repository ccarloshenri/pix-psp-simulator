# PIX PSP Simulator

A simulator for a Brazilian PIX Payment Service Provider (PSP), built with Go and the standard library. Covers immediate charges (Cob), due-date charges (CobV), payments (Pix), and refunds (Devolução).

## Running

```bash
go run main.go
# Server starts on :8080 (override with PORT env var)
```

## Architecture

```
Controller → Processor → BO → Repository (in-memory)
```

Each layer has a single responsibility: controllers parse HTTP, processors validate input, BOs run business logic, repositories store data.

---

## Endpoints

### Cob (Immediate Charge)

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/cob` | Create charge (auto-generated txid) |
| `PUT` | `/cob/{txid}` | Create charge with explicit txid |
| `GET` | `/cob` | List charges (filters: `status`, `inicio`, `fim`) |
| `GET` | `/cob/{txid}` | Get charge by txid |
| `PATCH` | `/cob/{txid}` | Update charge value/expiration |
| `DELETE` | `/cob/{txid}` | Cancel charge |
| `POST` | `/cob/simulate` | Simulate a PIX payment for a txid |
| `PUT` | `/cob/{e2eid}/devolucao/{id}` | Create refund on a received payment |
| `GET` | `/cob/{e2eid}/devolucao/{id}` | Get refund |

### CobV (Charge with Due Date)

| Method | Path | Description |
|--------|------|-------------|
| `PUT` | `/cobv/{txid}` | Create charge with due date |
| `GET` | `/cobv` | List charges (filters: `status`, `dataDeVencimento`, `inicio`, `fim`) |
| `GET` | `/cobv/{txid}` | Get charge by txid |
| `PATCH` | `/cobv/{txid}` | Update charge |

---

## Testing with curl

### 1. Create an immediate charge (auto txid)

```bash
curl -s -X POST http://localhost:8080/cob \
  -H "Content-Type: application/json" \
  -d '{
    "chave": "+5511999998888",
    "expiracao": 3600,
    "valor": { "original": "100.00" },
    "devedor": { "cpf": "12345678909", "nome": "João Silva" }
  }' | jq .
```

### 2. Create a charge with explicit txid

```bash
curl -s -X PUT http://localhost:8080/cob/mytxid123 \
  -H "Content-Type: application/json" \
  -d '{
    "chave": "+5511999998888",
    "valor": { "original": "50.00" }
  }' | jq .
```

### 3. Get a charge

```bash
curl -s http://localhost:8080/cob/mytxid123 | jq .
```

### 4. Update a charge

```bash
curl -s -X PATCH http://localhost:8080/cob/mytxid123 \
  -H "Content-Type: application/json" \
  -d '{ "valor": "75.00", "expiracao": 7200 }' | jq .
```

### 5. Create a charge with due date

```bash
curl -s -X PUT http://localhost:8080/cobv/cobvtxid001 \
  -H "Content-Type: application/json" \
  -d '{
    "chave": "+5511999998888",
    "valor": { "original": "200.00" },
    "devedor": { "cpf": "12345678909", "nome": "Maria Souza" },
    "calendario": {
      "dataDeVencimento": "2025-12-31",
      "validadeAposVencimento": 30
    }
  }' | jq .
```

### 6. List charges

```bash
# All charges
curl -s "http://localhost:8080/cob" | jq .

# Filter by status
curl -s "http://localhost:8080/cob?status=ATIVA" | jq .

# Filter by date range (RFC3339)
curl -s "http://localhost:8080/cob?inicio=2025-01-01T00:00:00Z&fim=2025-12-31T23:59:59Z" | jq .
```

### 7. Simulate a PIX payment

```bash
# Use the txid from a previously created cob/cobv
curl -s -X POST http://localhost:8080/cob/simulate \
  -H "Content-Type: application/json" \
  -d '{
    "txid": "mytxid123",
    "valor": "100.00",
    "infopagador": "Pagamento referente a fatura"
  }' | jq .
# Response includes the endToEndId of the received payment
```

### 8. Create a refund

```bash
# Use the endToEndId returned by the simulate endpoint
curl -s -X PUT http://localhost:8080/cob/E607469482025010112345678901/devolucao/dev001 \
  -H "Content-Type: application/json" \
  -d '{
    "valor": "50.00",
    "natureza": "ORIGINAL",
    "descricaoDevolucao": "Partial refund"
  }' | jq .
```

### 9. Get a refund

```bash
curl -s http://localhost:8080/cob/E607469482025010112345678901/devolucao/dev001 | jq .
```

### 10. List cobv charges

```bash
curl -s "http://localhost:8080/cobv" | jq .
curl -s "http://localhost:8080/cobv?dataDeVencimento=2025-12-31" | jq .
```

### 11. Cancel a charge

```bash
curl -s -X DELETE http://localhost:8080/cob/mytxid123 | jq .
```

---

## Running Tests

```bash
go test ./...
```

---

## Status Flows

**Cob / CobV status:**
- `ATIVA` → `CONCLUIDA` (after payment)
- `ATIVA` → `REMOVIDA_PELO_USUARIO_RECEBEDOR` (after DELETE)

**Devolução status:**
- `EM_PROCESSAMENTO` → `DEVOLVIDO` (simulated immediately)
- `EM_PROCESSAMENTO` → `NAO_REALIZADO` (on failure)
