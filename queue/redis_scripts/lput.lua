local list = redis.call('lrange', KEYS[1], 0, -1)
for k, v in pairs(cjson.decode(ARGV[1])) do
	local pos = tonumber(k)
	if pos == nil then return redis.error_reply('positions must be numbers') end
	table.insert(list, pos+1, v)
end
redis.call('del', KEYS[1])
return redis.call('lpush', KEYS[1], unpack(list))
