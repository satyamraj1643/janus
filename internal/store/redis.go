package store

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:embed token_bucket.lua

var tokenBucketScriptContent string

var tokenBucketScript = redis.NewScript(tokenBucketScriptContent)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(addr string) *RedisStore {

	return &RedisStore{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func (r *RedisStore) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Idempotency core logic

func (r *RedisStore) CheckAndMarkAdmitted(ctx context.Context, jobID string, window time.Duration) (bool, error) {
	// Construct a unique key for job's idempotency

	key := fmt.Sprintf("janus:idempotency:%s", jobID)

	// SETNX (Set if Not exists)

	// If the key is set successfully, it returns true (meaning it is a new job)
	// If the key already exists, it returns false (meaning it is a duplicate job)

	isNew, err := r.client.SetNX(ctx, key, true, window).Result()
	if err != nil {
		return false, err
	}

	// We return "true" if it WAS admitted previously (ie isNew is false)

	return !isNew, nil
}

// Rate limiting core logic
func (r *RedisStore) AllowRequest(ctx context.Context, key string, limit int, window time.Duration, cost int) (bool, error) {
	// 1. Create a unique key for rate limit

	// format : janus:ratelimit:<resource_name>

	redisKey := fmt.Sprintf("janus:ratelimit:%s", key)

	count, err := r.client.IncrBy(ctx, redisKey, int64(cost)).Result()

	if err != nil {
		return false, err
	}

	// 3. If this is the first request (count == 1) set the expiration window

	if count == int64(cost) {
		r.client.Expire(ctx, redisKey, window)
	}

	if count > int64(limit) {
		return false, nil // Rejected
	}

	return true, nil // Allowed

}

func (r *RedisStore) AllowRequestTokenBucket(ctx context.Context, key string, capacity int, refillRate float64, cost int) (bool, error) {
	tokensKey := fmt.Sprintf("janus:quota:%s:tokens", key)
	timestampKey := fmt.Sprintf("janus:quota:%s:ts", key)

	now := float64(time.Now().UnixNano()) / 1e9 // Current time in seconds

	//Keys : [tokensKey, timestampKey]
	//Args : [capacity, refill_rate, cost, now]

	result, err := tokenBucketScript.Run(ctx, r.client, []string{tokensKey, timestampKey}, capacity, refillRate, cost, now).Result()
	if err != nil {
		return false, err
	}

	if res, ok := result.(int64); ok {
		return res == 1, nil
	}

	return false, nil
}

func (r *RedisStore) Flush(ctx context.Context) error {
	return r.client.FlushDB(ctx).Err()
}
