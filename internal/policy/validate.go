package policy

import "fmt"

func (p *Policy) Validate() error {
	if p.Version != 1 {
		return fmt.Errorf("unsupported policy version: %d", p.Version)
	}

	for depName, dep := range p.Dependencies {
		if dep.RateLimit == nil && dep.Concurrent == nil {
			return fmt.Errorf("dependency '%s' must define rate_limit or concurrent", depName)
		}
		if dep.MinIntervalMs < 0 {
			return fmt.Errorf("dependency '%s' min_interval_ms cannot be negative", depName)
		}
	}

	if p.GlobalExecutionLimit.MaxConcurrentPerTenant < 0 {
		return fmt.Errorf("global_execution_limit max_concurrent_per_tenant cannot be negative")
	}

	if p.DefaultJobPolicy.IdempotencyWindowMs < 0 {
		return fmt.Errorf("default_job_policy idempotency_window_ms cannot be negative")
	}

	if p.DefaultJobPolicy.Quarantine != nil {
		if p.DefaultJobPolicy.Quarantine.FailureThreshold <= 0 {
			return fmt.Errorf("default_job_policy quarantine failure_threshold must be > 0")
		}
		if p.DefaultJobPolicy.Quarantine.QuarantineDurationMs <= 0 {
			return fmt.Errorf("default_job_policy quarantine quarantine_duration_ms must be > 0")
		}
		if p.DefaultJobPolicy.Quarantine.MonitoringWindowMs <= 0 {
			return fmt.Errorf("default_job_policy quarantine monitoring_window_ms must be > 0")
		}
	}

	return nil
}
