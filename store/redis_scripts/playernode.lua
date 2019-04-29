for i = 1, #KEYS do
	local has = redis.call("sismember", KEYS[i], ARGV[#ARGV])
	if has == 1 then return ARGV[i] end
end
