-- KEYS: [tokens_key_1, ts_key_1, created_key_1, tokens_key_2, ts_key_2, created_key_2, ...]
-- ARGV: [now, count, cap1, rate1, cost1, min_int1, warmup1, cap2, rate2, cost2, min_int2, warmup2, ...]

local now_time = tonumber(ARGV[1])
local count = tonumber(ARGV[2])

--Tables to hold intermediate results so we do not query twice or calcualte twicw

local new_token_list = {}

--1. CHECK PHASE (Read-Only Logic) (Modified to write created_key for initialization)

for i = 0, count-1 do
    local base_arg = 3 + (i * 5) -- Stride 5
    local base_key = 1 + (i * 3) -- Stride 3

    local tokens_key = KEYS[base_key]
    local ts_key = KEYS[base_key + 1]
    local created_key = KEYS[base_key + 2]

    local capacity = tonumber(ARGV[base_arg])
    local refill_rate = tonumber(ARGV[base_arg + 1])
    local cost = tonumber(ARGV[base_arg + 2])
    local min_interval = tonumber(ARGV[base_arg + 3])
    local warmup_ms = tonumber(ARGV[base_arg + 4])

    -- a. Get/Set Created Time
    local created_at = tonumber(redis.call("get", created_key))
    if created_at == nil then
        created_at = now_time
        redis.call("set", created_key, now_time)
    end

    -- b. Apply Warm-up Scaling
    local effective_capacity = capacity
    local effective_rate = refill_rate

    if warmup_ms > 0 then
        local age = now_time - created_at
        local warmup_sec = warmup_ms / 1000.0
        if age < warmup_sec then
             local factor = 0.1 + (0.9 * (age / warmup_sec)) -- Start at 10%
             effective_capacity = capacity * factor
             effective_rate = refill_rate * factor
        end
    end

    -- c. Get current tokens
    local last_tokens = tonumber(redis.call("get", tokens_key))
    if last_tokens == nil then
        last_tokens = effective_capacity -- Start full (relative to effective)
    end

    -- d. Get current timestamp
    local last_ts = tonumber(redis.call("get", ts_key))
    if last_ts == nil then
        last_ts = 0 -- Burst Smoothing: Allow first request
    end

    -- e. Calculate refill
    local delta = math.max(0, now_time - last_ts)

    -- Burst Smoothing Check
    if delta < min_interval then
        return 0 
    end

    local filled = math.min(effective_capacity, last_tokens + (delta * effective_rate))

    -- f. Check cost
    if filled < cost then 
        return 0 
    end

    new_token_list[i+1] = filled - cost

end

--2. COMMIT PHASE (Write logic)

for i = 0, count - 1 do
    local base_key = 1 + (i * 3)
    local tokens_key = KEYS[base_key]
    local ts_key = KEYS[base_key + 1]

    redis.call("set", tokens_key, new_token_list[i+1])
    redis.call("set", ts_key, now_time)
end

return 1
    
    
    


