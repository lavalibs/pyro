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

local SOURCE = KEYS[1]
local DESTINATION = KEYS[2]
local COUNT = tonumber(ARGV[1])

if COUNT == 0 then return {} end

if COUNT == 1 then -- if there's only one, redis has a built-in command for this
  local key = redis.call('rpoplpush', SOURCE, DESTINATION)

  if key then return {key} end
  return {}
end

local elems = {}
if COUNT < 0 then -- negative numbers mean we need to reverse direction
  for i = 1, COUNT * -1 do
    elems[i] = redis.call('lpop', DESTINATION)
    if not elems[i] then break end
  end

  if #elems > 0 then redis.call('rpush', SOURCE, unpack(elems)) end
else
  for i = 1, COUNT do
    elems[i] = redis.call('rpop', SOURCE)
    if not elems[i] then break end
  end

  if #elems > 0 then redis.call('lpush', DESTINATION, unpack(elems)) end
end

return elems
