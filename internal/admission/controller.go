package admission

import (
	"context"
	"fmt"
	"time"
	"github.com/satyamraj1643/janus/internal/policy"
	"github.com/satyamraj1643/janus/internal/store"
	"github.com/satyamraj1643/janus/spec"
)



type AdmissionController struct {
	Policy *policy.Policy
	Store  store.StateStore
}

/*
	NewAdmissionController is PURE.
	All wiring must happen in main().
*/
func NewAdmissionController(
	p *policy.Policy,
	s store.StateStore,
) *AdmissionController {
	return &AdmissionController{
		Policy: p,
		Store:  s,
	}
}

func (ac *AdmissionController) reject(
	job spec.Job,
	reason string,
	err error,
) (*spec.JobDecision, error) {

	return &spec.JobDecision{
		JobID:     job.ID,
		BatchID:   job.BatchID,
		BatchName: job.BatchName,
		Status:    "rejected",
		Reason:    reason,
		Timestamp: time.Now(),
		Job:       job,
	}, err
}

func (ac *AdmissionController) accept(
	job spec.Job,
) *spec.JobDecision {

	return &spec.JobDecision{
		JobID:     job.ID,
		BatchID:   job.BatchID,
		BatchName: job.BatchName,
		Status:    "accepted",
		Timestamp: time.Now(),
		Job:       job,
	}
}

func (ac *AdmissionController) Check(
	ctx context.Context,
	job spec.Job,
) (*spec.JobDecision, error) {

	// 0. Priority check
	if err := ac.checkPriority(ctx, job); err != nil {
		return ac.reject(job, "priority_too_low", err)
	}

	// 1. Idempotency check
	if err := ac.checkIdempotency(ctx, job); err != nil {
		return ac.reject(job, "duplicate_request", err)
	}

	// 2. Prepare rate-limit requests
	var reqs []store.RateLimitReq
	reqs = append(reqs, ac.getGlobalLimitParameters())
	reqs = append(reqs, ac.getTenanatQuotaParams(job))
	reqs = append(reqs, ac.getDependencyParams(job)...)

	// 3. Atomic verification
	allowed, err := ac.Store.AllowRequestAtomic(ctx, reqs)
	if err != nil {
		_ = ac.Store.ClearIdempotency(ctx, job.ID)
		return ac.reject(job, "store_error", err)
	}

	if !allowed {
		_ = ac.Store.ClearIdempotency(ctx, job.ID)
		return ac.reject(job, "rate_limit_exceeded", fmt.Errorf("quota exceeded"))
	}

	return ac.accept(job), nil
}
