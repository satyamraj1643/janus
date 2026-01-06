package db

import (
	"context"
	"encoding/json"
	"log"

	"github.com/satyamraj1643/janus/spec"
)

func SaveJob(decision *spec.JobDecision) error {
	ctx := context.Background()
	job := decision.Job

	admittedInc := 0
	if decision.Status == "accepted" {
		admittedInc = 1
	}

	// 1. Upsert Batch
	// Check if this is a new batch insertion to update user_association later
	var isNewBatch bool

	// Using xmax = 0 to detect if it was an insert.
	// Note: In some PG versions/configurations this might be tricky, but typically xmax=0 implies insertion.
	// Alternative: we can't easily get 'is_new' from ON CONFLICT without xmax check or semantic difference.

	// We init total_jobs=1 because this current job is part of it.
	err := Pool.QueryRow(ctx,
		`INSERT INTO batch (batch_id, user_id, batch_name, created_at, total_jobs, admitted_jobs) 
		 VALUES ($1, $2, $3, NOW(), 1, $4)
		 ON CONFLICT (batch_id) DO UPDATE 
		 SET total_jobs = batch.total_jobs + 1, 
		     admitted_jobs = batch.admitted_jobs + $4
		 RETURNING (xmax = 0)`,
		decision.BatchID, job.OwnerID, decision.BatchName, admittedInc,
	).Scan(&isNewBatch)

	if err != nil {
		log.Printf("Error upserting batch %s for job %s: %v", decision.BatchID, decision.JobID, err)
		// We continue, but aware that batch stats might be desynced or FK failed.
		// If FK failed (user not found), the next insert (jobs) will also likely fail if it relies on same user.
	}

	batchInc := 0
	if isNewBatch {
		batchInc = 1
	}

	// 2. Insert/Update Job
	payloadBytes, _ := json.Marshal(job.Payload)

	query := `
		INSERT INTO jobs (job_id, user_id, batch_id, job_status, job_payload, created_at, reason, global_config_id)
		VALUES ($1, $2, $3, $4, $5, NOW(), $6, $7)
		ON CONFLICT (job_id) DO UPDATE 
		SET job_status = $4, job_payload = $5, reason = $6
	`

	_, err = Pool.Exec(ctx, query,
		decision.JobID,
		job.OwnerID,
		decision.BatchID, // TEXT type now
		decision.Status,
		payloadBytes,
		decision.Reason,
		job.GlobalConfigID,
	)

	if err != nil {
		log.Printf("Failed to save job %s: %v", decision.JobID, err)
		return err
	}

	// 3. Update User Association Stats
	if job.GlobalConfigID != "" {
		succeededInc := 0
		failedInc := 0
		if decision.Status == "accepted" {
			succeededInc = 1
		} else {
			failedInc = 1
		}

		_, err = Pool.Exec(ctx,
			`INSERT INTO user_association (user_id, config_id, total_jobs, succeeded_jobs, failed_jobs, no_of_jobs, no_of_batches)
			 VALUES ($3, $4, 1, $1, $2, 1, $5)
			 ON CONFLICT (config_id) DO UPDATE 
			 SET total_jobs = COALESCE(user_association.total_jobs, 0) + 1, 
			     succeeded_jobs = COALESCE(user_association.succeeded_jobs, 0) + $1, 
				 failed_jobs = COALESCE(user_association.failed_jobs, 0) + $2,
				 no_of_jobs = COALESCE(user_association.no_of_jobs, 0) + 1,
				 no_of_batches = COALESCE(user_association.no_of_batches, 0) + $5`,
			succeededInc, failedInc, job.OwnerID, job.GlobalConfigID, batchInc,
		)
		if err != nil {
			log.Printf("Error updating user_association for user %s config %s: %v", job.OwnerID, job.GlobalConfigID, err)
		}
	}

	return nil
}
