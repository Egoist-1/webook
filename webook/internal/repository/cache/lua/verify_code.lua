local key = KEYS[1]
local val = ARGV[2]
local cnt = KEYS[1].."cnt"

local code = redis.call("GET",key)

local c = tonumber(redis.call("GET",cnt))
if c <= 0 then
    --验证次数过多
    return -1
end

if str == nil then
    --key 不存在
    return -3
end
if val ~= str then
    -- 验证码错误
    redis.call("DECR",cnt)
    return -2
end
if val == str then
    return 0
end