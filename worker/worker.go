package worker

import (
	"log"

	"github.com/satyamraj1643/janus/db"
	"github.com/satyamraj1643/janus/queue"
)

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
				if err := db.SaveJob(decision); err != nil {
					log.Printf("Error saving job %s: %v", decision.JobID, err)
				}
			}
		}(i)
	}
}
