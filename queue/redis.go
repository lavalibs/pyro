package queue

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/lavalibs/pyro/keys"
)

//go:generate esc -o scripts.go -pkg queue -private redis_scripts

var (
	// LPut puts elements into the list at the specified index
	// Keys:
	// - list to put the elements into
	// Values:
	// - ...json.Marshalled map[int]string where K is position and V is the value to insert
	LPut = redis.NewScript(_escFSMustString(false, "/redis_scripts/lput.lua"))

	// LOverride resets the elements in a list
	// Keys:
	// - list to override
	// Values:
	// - ...values to insert
	LOverride = redis.NewScript(_escFSMustString(false, "/redis_scripts/loverride.lua"))

	// LMove moves an element in a list
	// Keys:
	// - the list to move the elements in
	// Values:
	// - (int) from index
	// - (int) to index
	LMove = redis.NewScript(_escFSMustString(false, "/redis_scripts/lmove.lua"))

	// LShuffle shuffles a list
	// Keys:
	// - the list to shuffle
	// Values:
	// - randomization seed
	LShuffle = redis.NewScript(_escFSMustString(false, "/redis_scripts/lshuffle.lua"))

	// LRevSplice splices a list in reverse
	// Keys:
	// - the list to splice
	// Values:
	// - (int) start: the index to start splicing at
	// - (int) deleteCount: the number of elements to remove
	// - ...elements to insert
	LRevSplice = redis.NewScript(_escFSMustString(false, "/redis_scripts/lrevsplice.lua"))

	// MultiRPopLPush moves multiple elements from the right of one list to the left of another
	// Keys:
	// - RPop list
	// - LPush list
	// Values:
	// - (int) count: number of elements to move
	MultiRPopLPush = redis.NewScript(_escFSMustString(false, "/redis_scripts/multirpoplpush.lua"))
)

// RedisQueue represents a song queue in redis
type RedisQueue struct {
	c *redis.Client
}

// Redis makes a new Redis queue
func Redis(c *redis.Client) *RedisQueue {
	return &RedisQueue{c}
}

// Add adds songs to the end queue
func (q *RedisQueue) Add(guildID uint64, tracks []string) error {
	intr := make([]interface{}, len(tracks))
	for i, track := range tracks {
		intr[i] = track
	}

	return q.c.LPush(keys.PrefixPlayerQueue.Fmt(guildID), intr...).Err()
}

// Set overwrites songs in the queue
func (q *RedisQueue) Set(guildID uint64, tracks []string) error {
	intr := make([]interface{}, len(tracks))
	for i, track := range tracks {
		intr[i] = track
	}

	return LOverride.Run(q.c, []string{
		keys.PrefixPlayerQueue.Fmt(guildID),
	}, intr...).Err()
}

// Put puts the specified elements into the queue at the specified positions
func (q *RedisQueue) Put(guildID uint64, tracks map[int]string) error {
	b, err := json.Marshal(tracks)
	if err != nil {
		return err
	}

	return LPut.Run(q.c, []string{
		keys.PrefixPlayerQueue.Fmt(guildID),
	}, b).Err()
}

// Unshift adds songs to the front of the queue
func (q *RedisQueue) Unshift(guildID uint64, tracks ...string) error {
	l := len(tracks)
	intr := make([]interface{}, l)
	for i, t := range tracks {
		intr[l-i-1] = t
	}
	return q.c.RPush(keys.PrefixPlayerQueue.Fmt(guildID), intr...).Err()
}

// Remove removes a song from the queue at the index
func (q *RedisQueue) Remove(guildID uint64, index int) error {
	return LRevSplice.Run(q.c, []string{
		keys.PrefixPlayerQueue.Fmt(guildID),
	}, index, 1).Err()
}

// Next advances the playlist
func (q *RedisQueue) Next(guildID uint64, count int) (skipped []string, err error) {
	res, err := MultiRPopLPush.Run(q.c, []string{
		keys.PrefixPlayerQueue.Fmt(guildID),
		keys.PrefixPlayerPrevious.Fmt(guildID),
	}, count).Result()
	intr := res.([]interface{})
	skipped = make([]string, len(intr))
	for i, t := range intr {
		skipped[i] = t.(string)
	}
	return
}

// Move moves songs in the list by index
func (q *RedisQueue) Move(guildID uint64, from, to int) error {
	return LMove.Run(q.c, []string{
		keys.PrefixPlayerQueue.Fmt(guildID),
	}, -from-1, -to-1).Err()
}

// Shuffle shuffles the queue
func (q *RedisQueue) Shuffle(guildID uint64) (list []string, err error) {
	res, err := LShuffle.Run(q.c, []string{
		keys.PrefixPlayerQueue.Fmt(guildID),
	}, time.Now().UnixNano()).Result()
	if err != nil {
		return
	}

	intr := res.([]interface{})
	list = make([]string, len(intr))
	for i, t := range intr {
		list[i] = t.(string)
	}
	return
}

// Splice splices the queue
func (q *RedisQueue) Splice(guildID uint64, start, deleteCount int, tracks ...string) (list []string, err error) {
	args := make([]interface{}, len(tracks)+2)
	args[0] = start
	args[1] = deleteCount
	for i, track := range tracks {
		args[i+2] = track
	}

	res, err := LRevSplice.Run(q.c, []string{
		keys.PrefixPlayerQueue.Fmt(guildID),
	}, args...).Result()
	if err != nil {
		return
	}

	intr := res.([]interface{})
	list = make([]string, len(intr))
	for i, t := range intr {
		list[i] = t.(string)
	}
	return
}

// Trim trims the queue
func (q *RedisQueue) Trim(guildID uint64, start, end int) error {
	return q.c.LTrim(keys.PrefixPlayerQueue.Fmt(guildID), -int64(end)-1, -int64(start)-1).Err()
}

// NowPlaying gets the currently playing track
func (q *RedisQueue) NowPlaying(guildID uint64) (np string, err error) {
	np, err = q.c.LIndex(keys.PrefixPlayerPrevious.Fmt(guildID), 0).Result()
	if err != nil && err == redis.Nil {
		err = nil
	}
	return
}

// List lists the songs in the queue
func (q *RedisQueue) List(guildID uint64, index int, count uint) (tracks []string, err error) {
	var last int64
	if index < 0 {
		last = -int64(uint(-index) + count)
	} else if count == 0 {
		last = -1
	} else {
		last = int64(uint(index) + count)
	}

	tracks, err = q.c.LRange(keys.PrefixPlayerQueue.Fmt(guildID), int64(index), last).Result()
	if err != nil {
		return
	}

	l := len(tracks)
	for i := 0; i < l/2; i++ {
		tracks[i], tracks[l-i-1] = tracks[l-i-1], tracks[i]
	}
	return
}
