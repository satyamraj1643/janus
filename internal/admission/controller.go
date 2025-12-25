package admission

import (
	"context"
	"fmt"
	"sync"
	"time" // Added import

	"github.com/satyamraj1643/janus/internal/policy"
	"github.com/satyamraj1643/janus/internal/store"
	"github.com/satyamraj1643/janus/spec"
)

type Decision struct {
	JobID     string    `json:"job_id"`
	Admitted  bool      `json:"admitted"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
}

type Stats struct {
	TotalRequests    int64            `json:"total_requests"`
	AdmittedRequests int64            `json:"admitted_requests"`
	RejectedRequests int64            `json:"rejected_requests"`
	RejectionReasons map[string]int64 `json:"rejection_reasons"`
	RecentDecisions  []Decision       `json:"recent_decisions"` // <-- Log
}

type AdmissionController struct {
	Policy *policy.Policy
	Store  store.StateStore
	stats  Stats
	mu     sync.RWMutex
}

func NewAdmissionController(p *policy.Policy, s store.StateStore) *AdmissionController {
	return &AdmissionController{
		Policy: p,
		Store:  s,
		stats: Stats{
			RejectionReasons: make(map[string]int64),
			RecentDecisions:  make([]Decision, 0), // Init
		},
	}
}

func (ac *AdmissionController) GetStats() Stats {
	ac.mu.RLock()
	defer ac.mu.RUnlock()
	return ac.stats
}

func (ac *AdmissionController) UpdatePolicy(newPolicy *policy.Policy) {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.Policy = newPolicy
}

func (ac *AdmissionController) recordDecision(jobID string, admitted bool, reason string) {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	// Update Counters
	ac.stats.TotalRequests++
	if admitted {
		ac.stats.AdmittedRequests++
	} else {
		ac.stats.RejectedRequests++
		if ac.stats.RejectionReasons == nil {
			ac.stats.RejectionReasons = make(map[string]int64)
		}
		ac.stats.RejectionReasons[reason]++
	}

	// Update Log (Keep last 50)
	decision := Decision{
		JobID:     jobID,
		Admitted:  admitted,
		Reason:    reason,
		Timestamp: time.Now(),
	}

	// Prepend or Append? Append is easier, Frontend can reverse.
	ac.stats.RecentDecisions = append(ac.stats.RecentDecisions, decision)
	if len(ac.stats.RecentDecisions) > 50 {
		ac.stats.RecentDecisions = ac.stats.RecentDecisions[1:] // Remove oldest
	}
}

func (ac *AdmissionController) Check(ctx context.Context, job spec.Job) error {
	// 0. Priority Check
	if err := ac.checkPriority(ctx, job); err != nil {
		ac.recordDecision(job.ID, false, "priority_too_low")
		return err
	}

	// 1. Idempotency Check
	if err := ac.checkIdempotency(ctx, job); err != nil {
		ac.recordDecision(job.ID, false, "duplicate_request")
		return err
	}

	// 2. Prepare Limits
	var reqs []store.RateLimitReq
	reqs = append(reqs, ac.getGlobalLimitParameters())
	reqs = append(reqs, ac.getTenanatQuotaParams(job))
	depReqs := ac.getDependencyParams(job)
	reqs = append(reqs, depReqs...)

	// 3. Atomic Verification
	allowed, err := ac.Store.AllowRequestAtomic(ctx, reqs)
	if err != nil {
		_ = ac.Store.ClearIdempotency(ctx, job.ID)
		ac.recordDecision(job.ID, false, "store_error")
		return fmt.Errorf("admission check failed: %v", err)
	}

	if !allowed {
		_ = ac.Store.ClearIdempotency(ctx, job.ID)
		ac.recordDecision(job.ID, false, "rate_limit_exceeded")
		return fmt.Errorf("job rejected due to quota limits (atomic check)")
	}

	// Success
	ac.recordDecision(job.ID, true, "admitted")
	return nil
}
