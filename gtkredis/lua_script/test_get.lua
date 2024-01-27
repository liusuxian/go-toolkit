---@diagnostic disable: undefined-global
local result = redis.call('GET', KEYS[1])
if not result then
  return 1
end
result = redis.call('GET', KEYS[2])
if not result then
  return 2
end
return 3
