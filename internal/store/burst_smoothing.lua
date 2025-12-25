-- KEYS: [timestamp_key]
-- ARGV: [now, min_interval]

local ts_key = KEYS[1]
local now_time = tonumber(ARGV[1])
local min_interval = tonumber(ARGV[2])

local last_ts = tonumber(redis.call("get", ts_key))
if last_ts == nil then
    last_ts = 0
end

local delta = math.max(0, now_time - last_ts)

if delta < min_interval then 
    return 0  -- Not enough time has passed to admit this job yet
end

-- Allowed: Update timestamp
redis.call("set", ts_key, now_time)
return 1
