package worker

import (
	"log"
	"github.com/satyamraj1643/janus/queue"
)


func Start(n int) {
	for i:= range n {
		go func (id int) {
			log.Println("worker started:", id)
			for job := range queue.JobQueue {
				log.Printf("worker %d processing job %s", id, job.ID)
				// will implement this later...
			}
		}(i)
	}
}