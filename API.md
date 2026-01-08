# Janus API Documentation

**Base URL:** `https://janus-microservice.onrender.com`

---

## Authentication

All job routes require an `X-User-ID` header containing the authenticated Janus user's UUID.

```
X-User-ID: 2dad64a8-3f87-4d6e-9b4c-1cfa5917fd4b
```

---

## Endpoints

### Health Check

| Method | Route | Auth Required |
|--------|-------|---------------|
| GET | `/health` | No |

**Response:**
```json
{"status": "ok"}
```

---

### Single Job Creation

| Method | Route | Auth Required |
|--------|-------|---------------|
| POST | `/dashboard/jobs` | Yes |
| POST | `/system/jobs` | Yes |

**Request Body:**
```json
{
  "batch_name": "my-batch",
  "tenant_id": "tenant-abc",
  "priority": 8,
  "dependencies": {
    "openai": 2,
    "stripe": 1
  },
  "payload": {
    "custom_key": "custom_value"
  }
}
```

**Response (Accepted):** `HTTP 202`
```json
{
  "job_id": "uuid-here",
  "status": "Accepted",
  "reason": ""
}
```

**Response (Rejected):** `HTTP 429`
```json
{
  "job_id": "uuid-here",
  "status": "Rejected",
  "reason": "dependency 'openai' rate limit exceeded"
}
```

---

### Batch Job Creation (Partial)

| Method | Route | Auth Required |
|--------|-------|---------------|
| POST | `/dashboard/jobs/batch` | Yes |
| POST | `/system/jobs/batch` | Yes |

Jobs are evaluated individually. Some may be accepted, some rejected.

**Request Body:**
```json
{
  "batch_name": "my-batch",
  "jobs": [
    {
      "tenant_id": "tenant-abc",
      "priority": 8,
      "dependencies": {"openai": 1}
    },
    {
      "tenant_id": "tenant-xyz",
      "priority": 3,
      "dependencies": {"stripe": 5}
    }
  ]
}
```

**Response:** `HTTP 200`
```json
{
  "total": 2,
  "accepted": 1,
  "rejected": 1
}
```

---

### Batch Job Creation (Atomic)

| Method | Route | Auth Required |
|--------|-------|---------------|
| POST | `/dashboard/jobs/batch/atomic` | Yes |
| POST | `/system/jobs/batch/atomic` | Yes |

**All-or-Nothing:** If ANY job fails validation, the ENTIRE batch is rejected and no side effects occur.

**Request Body:** Same as Partial Batch.

**Response (All Accepted):** `HTTP 202`
```json
[
  {"job_id": "uuid-1", "status": "Accepted", "reason": ""},
  {"job_id": "uuid-2", "status": "Accepted", "reason": ""}
]
```

**Response (Any Rejected):** `HTTP 207`
```json
[
  {"job_id": "uuid-1", "status": "Rejected", "reason": "priority too low"},
  {"job_id": "uuid-2", "status": "Rejected", "reason": "batch rejected atomically"}
]
```

---

## Field Descriptions

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `batch_name` | string | Yes | Name to group jobs under |
| `tenant_id` | string | Yes | Identifier for the tenant/customer |
| `priority` | int | Yes | 1-10, higher = more likely to be admitted |
| `dependencies` | map[string]int | No | External service name → cost (tokens consumed) |
| `payload` | object | No | Custom data passed through to workers |

---

## Error Codes

| HTTP Code | Meaning |
|-----------|---------|
| 200 | Success (Batch Summary) |
| 202 | Accepted |
| 400 | Bad Request (Invalid JSON) |
| 403 | Service Paused / No Active Config |
| 429 | Rate Limited / Rejected |
| 207 | Multi-Status (Atomic batch partial info) |
| 500 | Internal Server Error |
