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
	s store.StateStore,
) *AdmissionController {
	return &AdmissionController{
		Policy: nil,
		Store:  s,
	}
}

func (ac *AdmissionController) Reject(
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

func (ac *AdmissionController) Accept(
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

	// Parse per-job config from DB
	jobPolicy, err := policy.ParseConfig(job.Config)
	if err != nil {
		return ac.Reject(job, "invalid_config", err)
	}

	// Create a temporary controller with the job's policy
	tempAC := &AdmissionController{
		Policy: jobPolicy,
		Store:  ac.Store,
	}

	// 0. Priority check
	if err := tempAC.checkPriority(ctx, job); err != nil {
		return ac.Reject(job, "priority_too_low", err)
	}

	// 1. Idempotency check
	if err := tempAC.checkIdempotency(ctx, job); err != nil {
		return ac.Reject(job, "duplicate_request", err)
	}

	// 2. Prepare rate-limit requests
	var reqs []store.RateLimitReq
	reqs = append(reqs, tempAC.getGlobalLimitParameters())
	reqs = append(reqs, tempAC.getTenanatQuotaParams(job))
	reqs = append(reqs, tempAC.getDependencyParams(job)...)

	// 3. Atomic verification
	allowed, err := ac.Store.AllowRequestAtomic(ctx, reqs)
	if err != nil {
		_ = ac.Store.ClearIdempotency(ctx, job.ID)
		return ac.Reject(job, "store_error", err)
	}

	if !allowed {
		_ = ac.Store.ClearIdempotency(ctx, job.ID)
		return ac.Reject(job, "rate_limit_exceeded", fmt.Errorf("quota exceeded"))
	}

	return ac.Accept(job), nil
}

func (ac *AdmissionController) CheckBatchAtomic(
	ctx context.Context,
	jobs []spec.Job,
) ([]*spec.JobDecision, error) {
	if len(jobs) == 0 {
		return nil, nil
	}

	// Parse policy for the first job (assuming batch shares config for now, or mix)
	// Ideally, we respect each job's config.

	decisions := make([]*spec.JobDecision, len(jobs))
	var allReqs []store.RateLimitReq
	var validJobs []spec.Job
	var validIndices []int

	// 1. Pre-validation loop
	for i, job := range jobs {
		jobPolicy, err := policy.ParseConfig(job.Config)
		if err != nil {
			d, _ := ac.Reject(job, "invalid_config", err)
			decisions[i] = d
			continue
		}

		tempAC := &AdmissionController{Policy: jobPolicy, Store: ac.Store} // lightweight

		if err := tempAC.checkPriority(ctx, job); err != nil {
			d, _ := ac.Reject(job, "priority_too_low", err)
			decisions[i] = d
			continue
		}

		if err := tempAC.checkIdempotency(ctx, job); err != nil {
			d, _ := ac.Reject(job, "duplicate_request", err)
			decisions[i] = d
			continue
		}

		// Collect limits provided
		allReqs = append(allReqs, tempAC.getGlobalLimitParameters())
		allReqs = append(allReqs, tempAC.getTenanatQuotaParams(job))
		allReqs = append(allReqs, tempAC.getDependencyParams(job)...)

		validJobs = append(validJobs, job)
		validIndices = append(validIndices, i)
	}

	// If no valid jobs to check against DB, return early
	if len(validJobs) == 0 {
		return decisions, nil
	}

	// 2. Atomic DB Check
	allowed, err := ac.Store.AllowRequestAtomic(ctx, allReqs)

	if err != nil {
		// System error - reject all remaining
		for _, idx := range validIndices {
			_ = ac.Store.ClearIdempotency(ctx, jobs[idx].ID)
			d, _ := ac.Reject(jobs[idx], "store_error", err)
			decisions[idx] = d
		}
		return decisions, nil
	}

	if !allowed {
		// Atomic failure - reject all remaining
		for _, idx := range validIndices {
			_ = ac.Store.ClearIdempotency(ctx, jobs[idx].ID)
			d, _ := ac.Reject(jobs[idx], "batch_quota_exceeded", fmt.Errorf("atomic batch rejected"))
			decisions[idx] = d
		}
		return decisions, nil
	}

	// 3. Accept all valid
	for _, idx := range validIndices {
		decisions[idx] = ac.Accept(jobs[idx])
	}

	return decisions, nil
}
