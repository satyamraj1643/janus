# 📄 2. DESIGN DOCUMENT — Janus-V1

## Purpose of This Document

This document explains **how Janus is conceptually structured** and **how responsibilities are divided**.
It does **not** define inputs or configuration formats.

---

## 1. High-Level System Model

Janus is a **control plane** that mediates execution between producers and workers.

Conceptual flow:

```
Producer → Janus → Execution Queue → Worker → Janus
```

Janus appears:

* before execution (admission)
* after execution (retry/termination decision)

---

## 2. Responsibility Separation

### 2.1 Producers

* Request work
* Provide job identity and scope
* Never control timing or retries

### 2.2 Janus

* Decides execution timing
* Enforces limits and fairness
* Coordinates retries

### 2.3 Workers

* Execute exactly once
* Report outcome
* Never retry or delay

---

## 3. Design Principles

### 3.1 Single Scheduler Principle

Only one component (Janus) may schedule execution or retries.

### 3.2 Declarative Policy

Execution rules are declared statically, not embedded in jobs.

### 3.3 Fail-Safe Defaults

If execution safety cannot be determined, the job does not run.

### 3.4 Worker Minimalism

Workers are stateless executors, not decision-makers.

---

## 4. Core Concepts (Conceptual)

### 4.1 Dependency

A shared system that must be protected (API, service, DB).

### 4.2 Scope

A fairness boundary defining who shares limits.

### 4.3 Admission

The act of allowing a job to execute at a specific time.

### 4.4 Retry Coordination

Retries are treated as new execution attempts governed by the same rules.

---

## 5. Failure Model

Janus assumes:

* workers may crash
* executions may fail
* dependencies may be unavailable

Janus remains correct because:

* decisions are centralized
* retries are coordinated
* limits are enforced pre-execution

---

## 6. Design Guarantees

This design guarantees:

* no retry storms
* no dependency overload
* deterministic behavior
* isolation between tenants

---

## 7. Explicit Design Constraints

* Workers must not retry
* Janus must not interpret business logic
* Execution rules must be centralized
* Inputs must be explicit and validated

---

## 8. Design Summary

Janus transforms execution from:

> uncontrolled parallelism
> into
> centrally governed admission

This is the core architectural value.
