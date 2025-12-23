-- keys: [tokens_key, timestamp_key]
-- argv: [capacity, refill_rate, cost, now_time]

local tokens_key = KEYS[1]
local timestamp_key = KEYS[2]

local capacity = tonumber(ARGV[1])
local refill_rate = tonumber(ARGV[2])
local cost = tonumber(ARGV[3])
local now_time = tonumber(ARGV[4])

--1. Get current tokens
local last_tokens = tonumber(redis.call("get", tokens_key))
if last_tokens == nil then
    last_tokens = capacity
end

--2. Get last refill time
local last_refill = tonumber(redis.call("get", timestamp_key))
if last_refill == nil then
    last_refill = 0
end

--3. Calculate refill 
local delta = math.max(0, now_time-last_refill)
local filled_tokens = math.min(capacity, last_tokens+ (delta*refill_rate))

--4. Check if allowed
local allowed = filled_tokens >= cost
local new_tokens = filled_tokens

if allowed then
    new_tokens = filled_tokens - cost
    --Update state
    redis.call("set", tokens_key,new_tokens)
    redis.call("set", timestamp_key, now_time)

    return 1 --Allowed
else
    return 0 --Not allowed
end


