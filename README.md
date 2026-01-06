
# Janus Service

A high-performance, distributed **Job Admission & Rate Limiting Service** written in Go. Janus acts as a smart gatekeeper for downstream systems, enforcing granular quotas and policies at the edge.

## 🚀 Key Features

*   **Synchronous Admission Control**: Provides immediate `Accepted`/`Rejected` feedback to clients.
*   **Distributed Rate Limiting**: Uses **Redis Lua Scripts** for atomic, high-performance Token Bucket rate limiting.
*   **Atomic Batch Processing**: "All-or-Nothing" semantics for job batches—if one job fails admission, the entire batch is rejected.
*   **Dynamic Reconfiguration**: Updates policies in real-time without downtime using PostgreSQL `LISTEN/NOTIFY`.
*   **Multi-Level Quotas**: Enforces limits at Global, Tenant (User), and Dependency levels.

## 🏗 Architecture

### 1. Ingestion Layer (API)
*   **Endpoints**: JSON-based REST APIs for Single Jobs and Batches.
*   **Authentication**: Identifying users via `X-User-ID`.
*   **Middleware**: Efficient **Read-Through Caching** serves active configurations from memory (with Redis fallback/flush), minimizing DB hits.

### 2. Admission Controller (Logic Core)
Policies are strictly enforced in the following order:
1.  **Priority Check**: Drops low-priority jobs during load.
2.  **Idempotency**: Prevents duplicate processing within a time window.
3.  **Dependency Rate Limits**: Token Bucket check for external resource usage (unified for single & atomic jobs).
4.  **Tenant Quotas**: Fair usage limits per user.
5.  **Global Limits**: Safety valve for total system throughput.

### 3. Persistence Layer
*   **Async Writer**: Admitted jobs are queued and asynchronously persisted to **PostgreSQL**.
*   **State Store**: **Redis** maintains high-speed counters and token buckets for distributed state.

## 🛠 Tech Stack
*   **Language**: Go (Golang)
*   **Datastores**: PostgreSQL (Persistent Data), Redis (Ephemeral State/Rate Limiting)
*   **Deployment**: Dockerized Microservice

## ⚡ Quick Start
```bash
# Run locally
go run cmd/api/main.go
```
*Requires `DB_URL` (PostgreSQL) and `REDIS_ADDR` (Redis) environment variables.*
