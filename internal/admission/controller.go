package admission

import (
	"context"
	"github.com/satyamraj1643/janus/internal/policy"
	"github.com/satyamraj1643/janus/internal/store"
	"github.com/satyamraj1643/janus/spec"
)

type AdmissionController struct {
	Policy *policy.Policy
	Store  store.StateStore

	// In-memory counters removed for skeleton state
}

func NewAdmissionController(p *policy.Policy, s store.StateStore) *AdmissionController {
	return &AdmissionController{
		Policy: p,
		Store:  s,
	}
}

// Check evaluates if a job can be admitted. Returns nil if admitted, error if rejected.
func (ac *AdmissionController) Check(ctx context.Context, job spec.Job) error {

	 // 1 Idempotency check 

	//  if err := ac.checkIdempotency(ctx, job); err != nil {
	// 	return err
	//  }

	//  if err := ac.checkDependecyLimit(ctx, job); err != nil {
	// 	return err
	//  }

	//  if err := ac.checkTenantQuota(ctx, job); err != nil {
	// 	return err
	//  }

	if err := ac.checkGlobalLimit(ctx, job); err != nil {
		return err
	}



	 return nil
}
