package queue

import "github.com/satyamraj1643/janus/spec"

var JobQueue = make(chan spec.Job, 1024)
var ResultQueue = make(chan *spec.JobDecision, 1024)

func RemainingCapacity() int {
	return cap(JobQueue) - len(JobQueue)
}