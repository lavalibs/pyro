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

math.randomseed(tonumber(ARGV[1]))
local function shuffle(t)
  for i = #t, 1, -1 do
    local rand = math.random(i)
    t[i], t[rand] = t[rand], t[i]
  end
  return t
end

local KEY = KEYS[1]
local list = redis.call('lrange', KEY, 0, -1)

if #list > 0 then
  shuffle(list)
  redis.call('del', KEY)
  redis.call('lpush', KEY, unpack(list))
end

return list
