package policy

import (
	"time"

	"github.com/satyamraj1643/janus/spec"
)

type Policy struct {
	Version              int                         `json:"version"`
	GlobalExecutionLimit GlobalExecutionLimit        `json:"global_execution_limit"`
	Dependencies         map[string]DependencyPolicy `json:"dependencies"`
	DefaultJobPolicy     JobPolicy                   `json:"default_job_policy"`
}

type GlobalExecutionLimit struct {
	MaxJobs                int `json:"max_jobs"`
	WindowMs               int `json:"window_ms"`
	MaxConcurrentPerTenant int `json:"max_concurrent_per_tenant"` // : Prevent "Noisy Neighbor"
	MinPriority            int `json:"min_priority"`              // : Emergency "Kill Switch" gate
	MinIntervalMs          int `json:"min_interval_ms"`           // : Burst Smoothing
}

type DependencyPolicy struct {
	Type          spec.DependencyType `json:"type"`
	RateLimit     *spec.RateLimit     `json:"rate_limit"`
	Concurrent    *spec.Concurrency   `json:"concurrent"`
	MinIntervalMs int64               `json:"min_interval_ms"`
	WarmupMs      int64               `json:"warmup_ms"` // : Protects cold startups
}

// RateLimiter is a placeholder for token bucket state
type RateLimiter struct {
	Tokens float64
	Last   time.Time
}

type JobPolicy struct {
	Dependencies        map[string]int    `json:"dependencies"`
	IdempotencyWindowMs int64             `json:"idempotency_window_ms"`
	ScopeLimits         map[string]int    `json:"scope_limits"`
	ScopeKeys           []string          `json:"scope_keys"`
	Retry               RetryPolicy       `json:"retry"`
	Execution           ExecutionPolicy   `json:"execution"`
	Quarantine          *QuarantinePolicy `json:"quarantine,omitempty"` // : Poison Pill Protection
}

type QuarantinePolicy struct {
	FailureThreshold     int   `json:"failure_threshold"`      // : Strikes before ban
	QuarantineDurationMs int64 `json:"quarantine_duration_ms"` // : How long to ban
	MonitoringWindowMs   int64 `json:"monitoring_window_ms"`   // : Time window for strikes
}

type RetryPolicy struct {
	MaxAttempts    int    `json:"max_attempts"`
	Backoff        string `json:"backoff"`
	InitialDelayMs int    `json:"initial_delay_ms"`
}

type ExecutionPolicy struct {
	TimeoutMs int `json:"timeout_ms"`
}
