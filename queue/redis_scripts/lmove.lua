--[[
Reused with permission from github.com/lavalibs/lavaqueue

MIT License

Copyright (c) 2018 Will Nelson

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
]]

local KEY = KEYS[1]
local FROM = tonumber(ARGV[1])
local TO = tonumber(ARGV[2])

if FROM == nil then return redis.redis_error('origin must be a number') end
if TO == nil then return redis.redis_error('destination must be a number') end

local list = redis.call('lrange', KEY, 0, -1)

if FROM == TO then return list end
if FROM < 0 then FROM = #list + FROM end
if TO < 0 then TO = #list + TO end

-- provided indexes are 0-based
local val = table.remove(list, FROM + 1)
table.insert(list, TO + 1, val)

redis.call('del', KEY)
redis.call('rpush', KEY, unpack(list))
return list
