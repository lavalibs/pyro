package queue

import (
	"testing"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

var q *RedisQueue

func init() {
	c := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	c.FlushDb()
	q = &RedisQueue{c}
}

func TestInterface(t *testing.T) {
	assert.Implements(t, (*Queue)(nil), q)
}

func TestAdd(t *testing.T) {
	tracks := []string{"a", "b", "c"}
	err := q.Add(1, tracks)
	assert.NoError(t, err)

	list, err := q.List(1, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, tracks, list)
	q.c.FlushDB()
}

func TestSet(t *testing.T) {
	tracks := []string{"a", "b", "c"}
	err := q.Set(1, tracks)
	assert.NoError(t, err)

	list, err := q.List(1, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, tracks, list)
	q.c.FlushDB()
}

func TestPut(t *testing.T) {
	tracks := []string{"a", "c", "d", "e", "f"}
	err := q.Set(1, tracks)
	assert.NoError(t, err)

	err = q.Put(1, map[int]string{
		1: "b",
	})
	assert.NoError(t, err)

	list, err := q.List(1, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "c", "d", "e", "f"}, list)
	q.c.FlushDB()
}

func TestUnshift(t *testing.T) {
	tracks := []string{"a", "b", "c"}
	err := q.Unshift(1, tracks...)
	assert.NoError(t, err)

	list, err := q.List(1, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, tracks, list)

	newTracks := []string{"d", "e", "f", "g"}
	err = q.Unshift(1, newTracks...)
	assert.NoError(t, err)

	expected := append(newTracks, tracks...)
	list, err = q.List(1, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, expected, list)
}

func TestRemove(t *testing.T) {
	tracks := []string{"a", "b", "c", "d", "e"}
	err := q.Set(1, tracks)
	assert.NoError(t, err)

	err = q.Remove(1, 1)
	assert.NoError(t, err)

	list, err := q.List(1, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "c", "d", "e"}, list)
}
