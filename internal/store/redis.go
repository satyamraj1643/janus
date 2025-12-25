package store

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// -- Embedding Lua Scripts Start --

//go:embed tenant_starvation.lua
var tenantStarvationScriptContent string
var tenantStarvationScript = redis.NewScript(tenantStarvationScriptContent)

//go:embed atomic_token_bucket.lua
var atomicTokenBucketScriptContent string
var atomicTokenBucketScript = redis.NewScript(atomicTokenBucketScriptContent)

//go:embed burst_smoothing.lua
var burstSmoothingScriptContent string
var burstSmoothingScript = redis.NewScript(burstSmoothingScriptContent)

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

// Standalone Idempotency logic

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

// Standalone external API Rate limiting logic
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

// Standalone tenanat-starvation logic

func (r *RedisStore) AllowRequestTokenBucket(ctx context.Context, key string, capacity int, refillRate float64, cost int) (bool, error) {
	tokensKey := fmt.Sprintf("janus:quota:%s:tokens", key)
	timestampKey := fmt.Sprintf("janus:quota:%s:ts", key)

	now := float64(time.Now().UnixNano()) / 1e9 // Current time in seconds

	//Keys : [tokensKey, timestampKey]
	//Args : [capacity, refill_rate, cost, now]

	result, err := tenantStarvationScript.Run(ctx, r.client, []string{tokensKey, timestampKey}, capacity, refillRate, cost, now).Result()
	if err != nil {
		return false, err
	}

	if res, ok := result.(int64); ok {
		return res == 1, nil
	}

	return false, nil
}

// One single handler for API Rate Limiting, Tenanat Starvation and global execution limit.
func (r *RedisStore) AllowRequestAtomic(ctx context.Context, reqs []RateLimitReq) (bool, error) {
	if len(reqs) == 0 {
		return true, nil
	}

	keys := make([]string, 0, len(reqs)*2)
	args := make([]any, 0, 2+(len(reqs)*3))

	now := float64(time.Now().UnixNano()) / 1e9
	args = append(args, now, len(reqs))

	for _, req := range reqs {
		keys = append(keys, fmt.Sprintf("janus:quota:%s:tokens", req.Key))
		keys = append(keys, fmt.Sprintf("janus:quota:%s:ts", req.Key))
		keys = append(keys, fmt.Sprintf("janus:quota:%s:created", req.Key))
		args = append(args, req.Capacity, req.RefillRate, req.Cost, req.MinInterval, req.WarmupMs)
	}

	res, err := atomicTokenBucketScript.Run(ctx, r.client, keys, args...).Result()
	if err != nil {
		return false, err
	}

	return res.(int64) == 1, nil
}

// AllowBurstSmoothing implements [StateStore].
func (r *RedisStore) AllowBurstSmoothing(ctx context.Context, key string, minIntervalSeconds float64) (bool, error) {
	tsKey := fmt.Sprintf("janus:smoothing:%s:ts", key)
	now := float64(time.Now().UnixNano()) / 1e9

	res, err := burstSmoothingScript.Run(ctx, r.client, []string{tsKey}, now, minIntervalSeconds).Result()
	if err != nil {
		return false, err
	}

	return res.(int64) == 1, nil

}

func (r *RedisStore) ClearIdempotency(ctx context.Context, jobID string) error {
	key := fmt.Sprintf("janus:idempotency:%s", jobID)
	return r.client.Del(ctx, key).Err()
}

func (r *RedisStore) Flush(ctx context.Context) error {
	return r.client.FlushDB(ctx).Err()
}
