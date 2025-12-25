package admission

import (
	"context"
	"fmt"
	"time"

	"github.com/satyamraj1643/janus/internal/store"
	"github.com/satyamraj1643/janus/spec"
)

// checkIdempotency verifies if the job has already been admitted duirng that window or not

func (ac *AdmissionController) checkIdempotency(ctx context.Context, job spec.Job) error {
	window := time.Duration(ac.Policy.DefaultJobPolicy.IdempotencyWindowMs) * time.Millisecond

	exists, err := ac.Store.CheckAndMarkAdmitted(ctx, job.ID, window)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("job %s already submitted within window %v", job.ID, window)
	}

	return nil
}

// check external API dependecies limit to determine of per window rate limits are crossed or not

func (ac *AdmissionController) checkDependecyLimit(ctx context.Context, job spec.Job) error {
	for depName, cost := range job.Dependencies {
		// Get the policy

		depPolicy, ok := ac.Policy.Dependencies[depName]
		if !ok || depPolicy.RateLimit == nil {
			continue // No limit provided skip
		}

		// Check limit

		limit := depPolicy.RateLimit.MaxRequests
		window := time.Duration(depPolicy.RateLimit.WindowMs) * time.Millisecond

		allowed, err := ac.Store.AllowRequest(ctx, depName, limit, window, cost)

		if err != nil {
			return err
		}

		if !allowed {
			return fmt.Errorf("dependency '%s' rate limit exceeded (cost %d, limit %d)", depName, cost, limit)
		}
	}

	return nil
}

func (ac *AdmissionController) checkTenantQuota(ctx context.Context, job spec.Job) error {
	limit := ac.Policy.GlobalExecutionLimit.MaxConcurrentPerTenant // Capacity

	windowMs := ac.Policy.GlobalExecutionLimit.WindowMs

	//Calculate refill rate (Tokens per second)
	// Avoid division by zero

	if windowMs == 0 {
		windowMs = 1000
	}

	refillRate := float64(limit) / (float64(windowMs) / 1000.0)

	tenantKey := fmt.Sprintf("tenant:%s", job.TenantID)

	// Call the token bucket method
	// Cost = 1 (assuming 1 slot per job for rn)

	allowed, err := ac.Store.AllowRequestTokenBucket(ctx, tenantKey, limit, refillRate, 1)

	if err != nil {
		return err
	}

	if !allowed {
		return fmt.Errorf("tenant '%s' quota exceeded", job.TenantID)
	}

	return nil

}

func (ac *AdmissionController) checkGlobalLimit(ctx context.Context, job spec.Job) error {
	limit := ac.Policy.GlobalExecutionLimit.MaxJobs
	windowsMs := ac.Policy.GlobalExecutionLimit.WindowMs

	if windowsMs == 0 {
		windowsMs = 1000
	}

	refillRate := float64(limit) / (float64(windowsMs) / 1000.0)

	key := "janus:global_request_quota"

	allowed, err := ac.Store.AllowRequestTokenBucket(ctx, key, limit, refillRate, 1)

	if err != nil {
		return err
	}

	if !allowed {
		return fmt.Errorf("global execution limit exceeded")
	}

	return nil

}

func (ac *AdmissionController) checkPriority(ctx context.Context, job spec.Job) error {
	minPriority := ac.Policy.GlobalExecutionLimit.MinPriority

	if job.Priority < minPriority {
		return fmt.Errorf("job priority %d is below minimum threshold %d", job.Priority, minPriority)
	}
	return nil

}

// 1. Prepare global limit
func (ac *AdmissionController) getGlobalLimitParameters() store.RateLimitReq {
	limit := ac.Policy.GlobalExecutionLimit.MaxJobs
	windowMs := ac.Policy.GlobalExecutionLimit.WindowMs
	if windowMs == 0 {
		windowMs = 1000
	}
	refillRate := float64(limit) / (float64(windowMs) / 1000.0)

	return store.RateLimitReq{
		Key:         "global_request_quota",
		Capacity:    limit,
		RefillRate:  refillRate,
		Cost:        1,
		MinInterval: float64(ac.Policy.GlobalExecutionLimit.MinIntervalMs) / 1000.0,
	}
}

// 2. Prepare tenant limit
func (ac *AdmissionController) getTenanatQuotaParams(job spec.Job) store.RateLimitReq {
	limit := ac.Policy.GlobalExecutionLimit.MaxConcurrentPerTenant
	windowMs := ac.Policy.GlobalExecutionLimit.WindowMs
	if windowMs == 0 {
		windowMs = 1000
	}
	refillRate := float64(limit) / (float64(windowMs) / 1000.0)

	return store.RateLimitReq{
		Key:         fmt.Sprintf("tenant:%s", job.TenantID),
		Capacity:    limit,
		RefillRate:  refillRate,
		Cost:        1,
		MinInterval: float64(ac.Policy.GlobalExecutionLimit.MinIntervalMs) / 1000.0,
	}
}

// 3. Prepare dependency limit
func (ac *AdmissionController) getDependencyParams(job spec.Job) []store.RateLimitReq {
	var reqs []store.RateLimitReq

	for depName, cost := range job.Dependencies {
		policy, exists := ac.Policy.Dependencies[depName]
		if exists && policy.RateLimit != nil {
			limit := policy.RateLimit.MaxRequests
			windowMs := policy.RateLimit.WindowMs
			if windowMs == 0 {
				windowMs = 1000
			}
			refillRate := float64(limit) / (float64(windowMs) / 1000.0)

			reqs = append(reqs, store.RateLimitReq{
				Key:         fmt.Sprintf("dependency:%s", depName),
				Capacity:    limit,
				RefillRate:  refillRate,
				Cost:        cost,
				MinInterval: float64(policy.MinIntervalMs) / 1000.0,
				WarmupMs:    policy.WarmupMs,
			})
		}
	}
	return reqs
}

// Not relevent for any process for janus or jobs, but for standalone key wise burst smoothing.
func (ac *AdmissionController) CheckBurstSmoothing(ctx context.Context, key string, minIntervalSeconds float64) error {
	allowed, err := ac.Store.AllowBurstSmoothing(ctx, key, minIntervalSeconds)
	if err != nil {
		return err
	}
	if !allowed {
		return fmt.Errorf("burst smoothing limit exceeded for key '%s' (min_interval %fs)", key, minIntervalSeconds)
	}
	return nil
}
