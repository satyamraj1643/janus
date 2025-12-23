# 📄 3. INPUT ARCHITECTURE — Janus-V1

## Purpose of This Document

This document defines **exactly what data Janus consumes** and **from whom**.
It does **not** describe behavior, design philosophy, or internals.

---

## 1. Input Sources (Authoritative)

Janus consumes inputs from **exactly three sources**:

1. Platform configuration
2. Job producers
3. Workers

No other inputs are accepted.

---

## 2. Platform Input: Static Configuration

### Role

Defines execution policy.

### Categories

* global limits
* dependency definitions
* job type definitions

---

### Dependency Definition (Final Schema)

```yaml
dependencies:
  <dependency_id>:
    type: external_api | internal_service | database
    rate_limit: NA | {
      max_requests: <int>
      window_ms: <int>
    }
    concurrent: NA | {
      max_inflight: <int>
    }
```

---

### Job Type Definition

```yaml
job_types:
  <job_type>:
    dependencies: [<dependency_id>]
    scope_keys: [<string>]
    retry:
      max_attempts: <int>
      backoff: fixed | exponential
      initial_delay_ms: <int>
    execution:
      timeout_ms: <int>
```

---

## 3. Producer Input: Job Payload

### Fixed Format

```json
{
  "job_id": "string",
  "job_type": "string",
  "scope": { "key": "value" },
  "payload": { }
}
```

### Rules

* `job_type` must exist
* `scope` must include all required scope keys
* `payload` is opaque

---

## 4. Worker Input: Execution Outcome

### Fixed Format

```json
{
  "job_id": "string",
  "status": "SUCCESS | FAILURE"
}
```

### Rules

* exactly one outcome
* no retries
* no backoff
* no error metadata

---

## 6. Input Validation Guarantees

Janus rejects inputs if:

* required fields are missing
* unknown job types are used
* scope keys are incomplete
* config schema is violated

---

## 7. Input Architecture Summary

| Source   | Provides          |
| -------- | ----------------- |
| Platform | Policy            |
| Producer | Execution request |
| Worker   | Outcome           |

Janus derives all decisions from these inputs alone.
