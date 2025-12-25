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

	AllowRequestAtomic(ctx context.Context, reqs []RateLimitReq) (bool, error)

	// AllowBurstSmoothing checks if enough time has passed since the last request (Standalone)
	AllowBurstSmoothing(ctx context.Context, key string, minIntervalSeconds float64) (bool, error)

	// RecordFailure increments failure count. If threshhold met, qurantines the job
	// Returns true if the job was just qurantined

	// ClearIdempotency removes the idempotency key (used for retries)
	ClearIdempotency(ctx context.Context, jobID string) error
}

type RateLimitReq struct {
	Key         string
	Capacity    int
	RefillRate  float64
	Cost        int
	MinInterval float64
	WarmupMs    int64
}
