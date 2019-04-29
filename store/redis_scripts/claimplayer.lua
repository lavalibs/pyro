for i = 1, #KEYS-1 do
	local has = redis.call("sismember", KEYS[i], ARGV[1])
	if has == 1 then return KEYS[i] == KEYS[#KEYS] end
end
redis.call("sadd", KEYS[#KEYS], ARGV[1])
return true
