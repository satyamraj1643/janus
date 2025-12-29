package worker

import (
	"context"
	"log"

	"github.com/satyamraj1643/janus/internal/admission"
	"github.com/satyamraj1643/janus/queue"
)

func StartJanusService(n int, ac *admission.AdmissionController) {
	for i := range n {
		go func(id int) {
			log.Printf("Processor %d started:", id)
			for job := range queue.JobQueue {
				log.Printf("worker %d processing job %s", id, job.ID)
				// Run the janus logic

				log.Println("Processing job", job)

				ctx := context.Background()
				returnedJob, err := ac.Check(ctx, job)

				if err != nil {
					log.Printf("worker %d: error processing job %s: %v", id, job.ID, err)
					continue
				}

				// After processing, send to ResultQueue
				queue.ResultQueue <- returnedJob
			}
		}(i)
	}
}

func StartDBWriter(n int) {
	for i := range n {
		go func(id int) {
			log.Printf("DBWriter %d started", id)

			for decision := range queue.ResultQueue {
				log.Printf("DBWriter %d: saving job %s to DB", id, decision.JobID)
				log.Printf("  → Status:    %s", decision.Status)
				log.Printf("  → BatchName: %s", decision.BatchName)
				log.Printf("  → BatchID:   %s", decision.BatchID)
				log.Printf("  → Reason:    %s", decision.Reason)
				log.Printf("  → TenantID:  %s", decision.Job.TenantID)
				log.Printf("  → Priority:  %d", decision.Job.Priority)
				log.Printf("  → Timestamp: %v", decision.Timestamp)

				// Write to DB here
			}
		}(i)
	}
}
