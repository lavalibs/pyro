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
	q = &RedisQueue{c}
}

func TestInterface(t *testing.T) {
	assert.Implements(t, (*Queue)(nil), q)
}

func TestAdd(t *testing.T) {
	tracks := map[int]string{
		0: "a",
		1: "b",
		2: "c",
	}
	err := q.Add(1, tracks)
	assert.NoError(t, err)
}
