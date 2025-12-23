package store

import (
	"context"
	"time"
)

// StateStore defines the interface for persisting Job state
type StateStore interface {
	// Ping checks the connection to the store
	Ping(ctx context.Context) error

	// CheckAndMarkAdmitted return true if the jobID was already seen within the window
	CheckAndMarkAdmitted(ctx context.Context, jobID string, window time.Duration) (bool, error)

	// If a job requiring a certain external dependency is submitted, can it run or not based on how many jobs already queued for that service per second.

	AllowRequest(ctx context.Context, key string, limit int, window time.Duration, cost int) (bool, error)


	//Flush the datastore - USE WITH CAUTION
	Flush(ctx context.Context) error

	//AllowRequestTokenBucket checks usage against a refillable quota (Token Bucket)

	AllowRequestTokenBucket(ctx context.Context, key string, capacity int, refillRate float64, cost int) (bool, error)

}
