package queue

import (
	"sort"
	"time"

	"github.com/go-redis/redis"
	"github.com/lavalibs/pyro/store"
)

//go:generate esc -o scripts.go -pkg queue -private redis_scripts

var (
	// LPut puts an element into the list at the specified index
	LPut = redis.NewScript(_escFSMustString(false, "/redis_scripts/lput.lua"))

	// LOverride resets the elements in a list
	LOverride = redis.NewScript(_escFSMustString(false, "/redis_scripts/loverride.lua"))

	// MultiRPopLPush moves multiple elements from the right of one list to the left of another
	MultiRPopLPush = redis.NewScript(_escFSMustString(false, "/redis_scripts/multirpoplpush.lua"))
)

// RedisQueue represents a song queue in redis
type RedisQueue struct {
	c *redis.Client
}

// Add adds songs to the end queue
func (q *RedisQueue) Add(guildID uint64, tracks map[int]string) error {
	return LPut.Run(q.c, []string{
		store.KeyPrefixPlayerQueue.Fmt(guildID),
	}, tracks).Err()
}

// Set overwrites songs in the queue
func (q *RedisQueue) Set(guildID uint64, tracks []string) error {
	intr := make([]interface{}, len(tracks))
	for i, track := range tracks {
		intr[i] = track
	}

	return LOverride.Run(q.c, []string{
		store.KeyPrefixPlayerQueue.Fmt(guildID),
	}, intr...).Err()
}

// Unshift adds songs to the front of the queue
func (q *RedisQueue) Unshift(guildID uint64, tracks ...string) error {
	intr := make([]interface{}, len(tracks))
	for i, t := range tracks {
		intr[i] = t
	}
	return q.c.RPush(store.KeyPrefixPlayerQueue.Fmt(guildID), intr...).Err()
}

// Remove removes a song from the queue at the index
func (q *RedisQueue) Remove(guildID uint64, index int) error {
	// TODO: Redis will remove the first occurrance of the element, not the specific index
	return q.c.Eval(`local index = redis.call('lindex', KEYS[1], ARGV[1])
redis.call('lrem', KEYS[1], index)`, []string{
		store.KeyPrefixPlayerQueue.Fmt(guildID),
	}, index).Err()
}

// Next advances the playlist
func (q *RedisQueue) Next(guildID uint64, count int) (skipped []string, err error) {
	res, err := MultiRPopLPush.Run(q.c, []string{
		store.KeyPrefixPlayerQueue.Fmt(guildID),
		store.KeyPrefixPlayerPrevious.Fmt(guildID),
	}, count).Result()
	skipped = res.([]string)
	return
}

// Sort sorts the list
func (q *RedisQueue) Sort(guildID uint64, predicate func(a, b string) bool) (list []string, err error) {
	list, err = q.c.LRange(store.KeyPrefixPlayerQueue.Fmt(guildID), 0, -1).Result()
	if err != nil {
		return
	}

	sort.Slice(list, func(i, j int) bool {
		return predicate(list[i], list[j])
	})

	intr := make([]interface{}, len(list))
	for i, item := range list {
		intr[i] = item
	}

	err = q.c.Eval(_escFSMustString(false, "loverride.lua"), []string{
		store.KeyPrefixPlayerQueue.Fmt(guildID),
	}, intr...).Err()
	return
}

// Move moves songs in the list by index
func (q *RedisQueue) Move(guildID uint64, from, to int) error {
	return q.c.Eval(_escFSMustString(false, "lmove.lua"), []string{
		store.KeyPrefixPlayerQueue.Fmt(guildID),
	}, from, to).Err()
}

// Shuffle shuffles the queue
func (q *RedisQueue) Shuffle(guildID uint64) ([]string, error) {
	list, err := q.c.Eval(_escFSMustString(false, "lshuffle.lua"), []string{
		store.KeyPrefixPlayerQueue.Fmt(guildID),
	}, time.Now()).Result()
	if err != nil {
		return []string{}, err
	}

	return list.([]string), nil
}

// Splice splices the queue
func (q *RedisQueue) Splice(guildID uint64, start, deleteCount int, tracks ...string) ([]string, error) {
	args := make([]interface{}, len(tracks)+2)
	args[0] = start
	args[1] = deleteCount
	for i, track := range tracks {
		args[i+2] = track
	}

	list, err := q.c.Eval(_escFSMustString(false, "lrevsplice.lua"), []string{
		store.KeyPrefixPlayerQueue.Fmt(guildID),
	}, args...).Result()
	if err != nil {
		return []string{}, err
	}

	return list.([]string), nil
}

// Trim trims the queue
func (q *RedisQueue) Trim(guildID uint64, start, end int) error {
	return q.c.LTrim(store.KeyPrefixPlayerQueue.Fmt(guildID), int64(start), int64(end)).Err()
}

// NowPlaying gets the currently playing track
func (q *RedisQueue) NowPlaying(guildID uint64) (string, error) {
	return q.c.LIndex(store.KeyPrefixPlayerPrevious.Fmt(guildID), 0).Result()
}

// List lists the songs in the queue
func (q *RedisQueue) List(guildID uint64, index int, count uint) ([]string, error) {
	var last int64
	if index < 0 {
		last = -int64(uint(-index) + count)
	} else {
		last = int64(uint(index) + count)
	}

	return q.c.LRange(store.KeyPrefixPlayerQueue.Fmt(guildID), int64(index), last).Result()
}