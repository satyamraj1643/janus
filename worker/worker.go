package worker

import (
	"log"
	"github.com/satyamraj1643/janus/queue"
)


func StartJanusService(n int) {
	for i:= range n {
		go func (id int) {
			log.Printf("Processor %d started:", id)
			for job := range queue.JobQueue {
				log.Printf("worker %d processing job %s", id, job.ID)
				// Run the janus logic

				log.Println("Processing job", job)

				
            
				// After processing, send to ResultQueue
				queue.ResultQueue <- job
			}
		}(i)
	}
}


func StartDBWriter(n int){
	for i := range n {
		go func(id int){
			log.Printf("DBWriter %d started", id)

			for job := range queue.ResultQueue{
				log.Printf("DBWriter %d: saving %s to DB", id, job.ID)

				log.Println("Saving job", job)

				// Write to DB here
			}
		}(i)
	}
}