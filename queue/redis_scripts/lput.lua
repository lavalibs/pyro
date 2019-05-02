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
local function reverse(arr)
  if arr == nil then return arr end

  local i, j = 1, #arr
  while i < j do
    arr[i], arr[j] = arr[j], arr[i]

    i = i + 1
    j = j - 1
  end
  return arr
end

local list = redis.call('lrange', KEYS[1], 0, -1)
reverse(list)
for k, v in pairs(cjson.decode(ARGV[1])) do
	local pos = tonumber(k)
	if pos == nil then return redis.error_reply('positions must be numbers') end
	table.insert(list, pos+1, v)
end
redis.call('del', KEYS[1])
return redis.call('lpush', KEYS[1], unpack(list))
