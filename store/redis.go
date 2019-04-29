package store

import (
	"encoding/json"
	"errors"

	"github.com/go-redis/redis"
	"github.com/lavalibs/pyro/lavalink/types"
)

//go:generate esc -o scripts.go -pkg store -private redis_scripts

// ErrUnknownRedisResponse occurs when the library is expecting a redis response and receives one
// of incorrect type
var ErrUnknownRedisResponse = errors.New("unknown redis response")

var (
	// ClaimPlayer claims the player in redis
	ClaimPlayer = redis.NewScript(_escFSMustString(false, "/redis_scripts/claimplayer.lua"))

	// PlayerNode gets the player of the node
	PlayerNode = redis.NewScript(_escFSMustString(false, "/redis_scripts/playernode.lua"))
)

// Redis represents a clustering client. Used 1:1 with Lavalink nodes.
type Redis struct {
	c    *redis.Client
	opts *redis.Options
}

// NewRedis makes a new cluster client
func NewRedis(opts *redis.Options) *Redis {
	return &Redis{redis.NewClient(opts), opts}
}

// GetPlayer gets the cached player state for the given guild
func (r *Redis) GetPlayer(guildID uint64, state *types.PlayerState) (err error) {
	res, err := r.c.Get(KeyPrefixPlayerState.Fmt(guildID)).Bytes()
	if err != nil {
		if err == redis.Nil {
			err = nil
		}
		return
	}

	err = json.Unmarshal(res, &state)
	return
}

// SetPlayer sets player info
func (r *Redis) SetPlayer(upd types.PlayerUpdate) error {
	b, err := json.Marshal(upd.State)
	if err != nil {
		return err
	}

	return r.c.Set(KeyPrefixPlayerState.Fmt(upd.GuildID), b, 0).Err()
}

// GetVoiceUpdate gets voice update info.
func (r *Redis) GetVoiceUpdate(guildID uint64) (session string, server types.VoiceServerUpdate, err error) {
	session, err = r.c.Get(KeyPrefixVoiceSession.Fmt(guildID)).Result()
	if err != nil {
		if err == redis.Nil {
			err = nil
		}
		return
	}

	serverBytes, err := r.c.Get(KeyPrefixVoiceServer.Fmt(guildID)).Bytes()
	if err != nil {
		if err == redis.Nil {
			err = nil
		}
		return
	}

	err = json.Unmarshal(serverBytes, &server)
	return
}

// SetVoiceSession sets the voice session in Redis
func (r *Redis) SetVoiceSession(guildID uint64, sessionID string) error {
	return r.c.Set(KeyPrefixVoiceSession.Fmt(guildID), sessionID, 0).Err()
}

// SetVoiceServer sets the voice server information in Redis
func (r *Redis) SetVoiceServer(pk types.VoiceServerUpdate) error {
	b, err := json.Marshal(&pk)
	if err != nil {
		return err
	}

	return r.c.Set(KeyPrefixVoiceServer.Fmt(pk.GuildID), b, 0).Err()
}

// GetStats gets the stats for a node
func (r *Redis) GetStats(node string, stats *types.Stats) (err error) {
	b, err := r.c.Get(KeyPrefixNodeStats.Fmt(node)).Bytes()
	if err != nil {
		return
	}

	err = json.Unmarshal(b, stats)
	return
}

// SetStats sets the stats for a node. It also updates a set of node names sorted by system CPU
// load.
func (r *Redis) SetStats(node string, stats types.Stats) error {
	err := r.c.ZAdd(string(KeyStatsList), redis.Z{
		Member: node,
		Score:  stats.CPU.SystemLoad / float64(stats.CPU.Cores),
	}).Err()
	if err != nil {
		return err
	}

	b, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	return r.c.Set(KeyPrefixNodeStats.Fmt(node), b, 0).Err()
}

// CreateNode creates a node
func (r *Redis) CreateNode(name string) error {
	return r.c.SAdd(string(KeyNodes), name).Err()
}

// DeleteNode deletes a node
func (r *Redis) DeleteNode(name string) error {
	return r.c.SRem(string(KeyNodes), name).Err()
}

// ClaimPlayer claims a player for the node
func (r *Redis) ClaimPlayer(node string, guildID uint64) (bool, error) {
	err := r.CreateNode(node)
	if err != nil {
		return false, err
	}

	nodes, err := r.c.SMembers(string(KeyNodes)).Result()
	if err != nil {
		return false, err
	}

	nodeKeys := make([]string, len(nodes))
	for i, n := range nodes {
		nodeKeys[i] = KeyPrefixNodePlayers.Fmt(n)
	}

	return ClaimPlayer.Run(r.c, append(nodeKeys, KeyPrefixNodePlayers.Fmt(node)), guildID).Bool()
}

// PlayerNode gets the node that the player is running on
func (r *Redis) PlayerNode(guildID uint64) (string, error) {
	nodes, err := r.c.SMembers(string(KeyNodes)).Result()
	if err != nil {
		return "", err
	}

	args := make([]interface{}, len(nodes))
	nodeKeys := make([]string, len(nodes))
	for i, node := range nodes {
		nodeKeys[i] = KeyPrefixNodePlayers.Fmt(node)
		args[i] = node
	}

	node, err := PlayerNode.Run(r.c, nodeKeys, append(args, guildID)...).String()
	if err == redis.Nil {
		return "", nil
	}
	return node, err
}

// ReleasePlayer releases a player from a node
func (r *Redis) ReleasePlayer(node string, guildID uint64) (bool, error) {
	count, err := r.c.SRem(KeyPrefixNodePlayers.Fmt(node), guildID).Result()
	return count > 0, err
}

// AnnounceDeath destroys a node
func (r *Redis) AnnounceDeath(node string) error {
	return r.c.Publish(string(KeyNodeDeaths), node).Err()
}

// ConsumeDeaths consumes death notifications, ignoring the specified node
func (r *Redis) ConsumeDeaths(node string, deaths chan string) error {
	err := r.CreateNode(node)
	if err != nil {
		return err
	}

	c := redis.NewClient(r.opts)
	pubsub := c.Subscribe(string(KeyNodeDeaths))
	defer c.Close()
	defer pubsub.Close()

	rcv, err := pubsub.Receive()
	if err != nil {
		return err
	}

	switch rcv.(type) {
	case *redis.Subscription:
		if deaths != nil {
			deaths <- ""
		}
	case *redis.Message:
		r.handleDeath(rcv.(*redis.Message).Payload, node)
	case *redis.Pong:
	default:
		return ErrUnknownRedisResponse
	}

	for msg := range pubsub.Channel() {
		err = r.handleDeath(msg.Payload, node)
		if err != nil {
			return err
		}

		if deaths != nil {
			deaths <- msg.Payload
		}
	}

	return nil
}

func (r *Redis) handleDeath(from, to string) error {
	if from == to {
		return nil
	}

	for {
		player, err := r.c.SPop(KeyPrefixNodePlayers.Fmt(from)).Result()
		if err != nil {
			if err == redis.Nil {
				break
			}
			return err
		}

		err = r.c.SAdd(KeyPrefixNodePlayers.Fmt(to), player).Err()
		if err != nil {
			return err
		}
	}

	return nil
}
