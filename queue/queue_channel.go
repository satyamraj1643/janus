package queue

import "github.com/satyamraj1643/janus/spec"

var JobQueue = make(chan spec.Job, 3)

func RemainingCapacity() int {
	return cap(JobQueue) - len(JobQueue)
}