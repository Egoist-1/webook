local key = KEYS[1]
local val = ARGV[1]
local cnt = KEYS[1].."cnt"

local ttl = tonumber(redis.call("ttl",key))
if ttl == -1 then
    --没有过期时间
    return -1
elseif ttl == -2 or ttl < 540 then
    redis.call("set",key,val)
    redis.call("expire",600)
    redis.call("set",cnt,3)
    redis.call("expire",600)
    return 0
else
    --频繁发送
    return -2
end