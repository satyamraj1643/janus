package queue

import (
	"errors"

	"github.com/satyamraj1643/janus/spec"
)

var ErrQueueFull = errors.New("queue full")

func Admit(job spec.Job) error {
	select {
	case JobQueue <- job:
		return nil
	default:
		return ErrQueueFull
	}
}


// for atomic batch addition
func AdmitBatch(jobs []spec.Job) error {
	if len(jobs) == 0 {
		return nil
	}

	// Capacity check

	if RemainingCapacity() < len(jobs) {
		return ErrQueueFull
	}

	// Commit

	for _, job := range jobs {
		JobQueue <- job
	}

	return nil
}