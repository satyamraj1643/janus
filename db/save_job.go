package db

import (
	"context"
	"encoding/json"
	"log"

	"github.com/satyamraj1643/janus/spec"
)

func SaveJob(job spec.Job) error {
	ctx := context.Background()

	// 1. Ensure Batch Exists
	// We use ON CONFLICT (batch_id) DO NOTHING to efficiently ensure the batch row is present.
	// Since batch_id is now TEXT, we can pass our custom string IDs directly.
	_, err := Pool.Exec(ctx,
		"INSERT INTO batch (batch_id, user_id, batch_name, created_at, total_jobs, admitted_jobs) VALUES ($1, $2, $3, NOW(), 0, 0) ON CONFLICT (batch_id) DO NOTHING",
		job.BatchID, job.TenantID, job.BatchName,
	)
	if err != nil {
		log.Printf("Error ensuring batch %s for job %s: %v", job.BatchID, job.ID, err)
		// Proceeding might fail FK constraint if inserted failed, but we log it.
	}

	// 2. Insert/Update Job
	payloadBytes, _ := json.Marshal(job.Payload)

	// WARNING: 'completed' status must exist in enum. If not, use 'accepted'.
	status := "accepted"

	query := `
		INSERT INTO jobs (job_id, user_id, batch_id, job_status, job_payload, created_at, reason)
		VALUES ($1, $2, $3, $4, $5, NOW(), '')
		ON CONFLICT (job_id) DO UPDATE 
		SET job_status = $4, job_payload = $5
	`

	_, err = Pool.Exec(ctx, query,
		job.ID,
		job.TenantID,
		job.BatchID, // TEXT type now
		status,
		payloadBytes,
	)

	if err != nil {
		log.Printf("Failed to save job %s: %v", job.ID, err)
		return err
	}
	return nil
}
