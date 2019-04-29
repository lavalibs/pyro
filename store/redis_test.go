package store

import (
	// "os"
	"testing"

	"github.com/go-redis/redis"
	"github.com/lavalibs/pyro/lavalink/types"
	"github.com/stretchr/testify/assert"
)

var c *Redis

func makeClient(t *testing.T) *Redis {
	if c == nil {
		c = NewRedis(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
		c.c.FlushAll()
	}

	return c
}

func TestRedis(t *testing.T) {
	makeClient(t)
	pong, err := c.c.Ping().Result()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(pong)
}

func TestSetGetPlayer(t *testing.T) {
	r := makeClient(t)
	upd := types.PlayerUpdate{
		GuildID: 1,
		State: types.PlayerState{
			Time:     1,
			Position: 1,
		},
	}

	if err := r.SetPlayer(upd); err != nil {
		t.Fatal(err)
	}

	st := &types.PlayerState{}
	err := r.GetPlayer(1, st)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, upd.State, *st, "store stated not equal to saved state")
}

func TestSetGetVoice(t *testing.T) {
	r := makeClient(t)
	if err := r.SetVoiceSession(1, "session"); err != nil {
		t.Fatal(err)
	}

	server := types.VoiceServerUpdate{
		GuildID:  1,
		Token:    "token",
		Endpoint: "endpoint",
	}

	if err := r.SetVoiceServer(server); err != nil {
		t.Fatal(err)
	}

	sess, fetchedServer, err := r.GetVoiceUpdate(1)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, sess, "session", "unexpected voice session ID")
	assert.Equal(t, server, fetchedServer, "unexpected voice server info")
}

func TestSetGetStats(t *testing.T) {
	r := makeClient(t)
	stats := types.Stats{
		Players:        1,
		PlayingPlayers: 1,
		Uptime:         1,
		Memory: types.StatsMemory{
			Free:       1,
			Used:       1,
			Allocated:  1,
			Reservable: 1,
		},
		CPU: types.StatsCPU{
			Cores:        1,
			SystemLoad:   1,
			LavalinkLoad: 1,
		},
		FrameStats: types.StatsFrames{
			Sent:    1,
			Nulled:  1,
			Deficit: 1,
		},
	}

	if err := r.SetStats("node", stats); err != nil {
		t.Fatal(err)
	}

	fetchedStats := &types.Stats{}
	if err := r.GetStats("node", fetchedStats); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, stats, *fetchedStats, "unexpected stats values")
}

func TestClaimPlayer(t *testing.T) {
	r := makeClient(t)
	ok, err := r.ClaimPlayer("node", 1)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Error("attempt to claim node failed")
	}

	node, err := r.PlayerNode(1)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "node", node, "unexpected node name")
}

func TestReleasePlayer(t *testing.T) {
	r := makeClient(t)
	ok, err := r.ClaimPlayer("node", 1)
	assert.NoError(t, err)
	assert.True(t, ok, "attempt to claim node failed")

	node, err := r.PlayerNode(1)
	assert.NoError(t, err)
	assert.Equal(t, "node", node, "unexpected node name")

	ok, err = r.ReleasePlayer("node", 1)
	assert.NoError(t, err)
	assert.True(t, ok, "unable to release player")

	node, err = r.PlayerNode(1)
	assert.NoError(t, err)
	assert.Equal(t, "", node, "unexpected node name after release")
}

func TestDeathAnnouncement(t *testing.T) {
	r := makeClient(t)
	ready := make(chan struct{})
	go func() {
		err := r.ConsumeDeaths("other node", ready)
		assert.NoError(t, err)
	}()
	<-ready

	ok, err := r.ClaimPlayer("node", 1)
	assert.NoError(t, err)
	assert.True(t, ok)

	err = r.AnnounceDeath("node")
	assert.NoError(t, err)

	node, err := r.PlayerNode(1)
	assert.NoError(t, err)
	assert.Equal(t, "other node", node)
}
