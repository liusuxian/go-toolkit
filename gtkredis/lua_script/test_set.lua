---@diagnostic disable: undefined-global
local result = redis.call('SETEX', KEYS[1], 10, tonumber(ARGV[1], 10))
if result['ok'] then
  return 1
end
return 2
