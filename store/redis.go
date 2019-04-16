package store

import (
	"encoding/json"

	"github.com/go-redis/redis"
	"github.com/lavalibs/pyro/lavalink/types"
)

//go:generate esc -o scripts.go -pkg store -private redis_scripts

// Redis represents a clustering client. Used 1:1 with Lavalink nodes.
type Redis struct {
	*redis.Client
}

// NewRedis makes a new cluster client
func NewRedis(r *redis.Client) *Redis {
	return &Redis{r}
}

// GetPlayer gets the cached player state for the given guild
func (r *Redis) GetPlayer(guildID uint64) (st types.PlayerState, err error) {
	state, err := r.Get(KeyPrefixPlayerState.Fmt(guildID)).Bytes()
	if err != nil || len(state) == 0 {
		return
	}

	err = json.Unmarshal(state, &st)
	return
}

// SetPlayer sets player info
func (r *Redis) SetPlayer(upd types.PlayerUpdate) error {
	b, err := json.Marshal(upd.State)
	if err != nil {
		return err
	}

	err = r.Set(KeyPrefixPlayerState.Fmt(upd.GuildID), b, 0).Err()
	return err
}

// GetVoiceUpdate gets voice update info.
func (r *Redis) GetVoiceUpdate(guildID uint64) (session string, server types.VoiceServerUpdate, err error) {
	session, err = r.Get(KeyPrefixVoiceSession.Fmt(guildID)).Result()
	if err != nil {
		return
	}

	serverBytes, err := r.Get(KeyPrefixVoiceServer.Fmt(guildID)).Bytes()
	if err != nil {
		return
	}

	err = json.Unmarshal(serverBytes, &server)
	return
}

// SetVoiceState sets the voice state update in Redis
func (r *Redis) SetVoiceState(pk types.VoiceStateUpdate) error {
	return r.Set(KeyPrefixVoiceSession.Fmt(pk.GuildID), pk.SessionID, 0).Err()
}

// SetVoiceServer sets the voice server information in Redis
func (r *Redis) SetVoiceServer(pk types.VoiceServerUpdate) error {
	b, err := json.Marshal(&pk)
	if err != nil {
		return err
	}

	return r.Set(KeyPrefixVoiceServer.Fmt(pk.GuildID), b, 0).Err()
}

// GetStats gets the stats for a node
func (r *Redis) GetStats(node string) (stats types.Stats, err error) {
	b, err := r.Get(KeyPrefixNodeStats.Fmt(node)).Bytes()
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &stats)
	return
}

// SetStats sets the stats for a node. It also updates a set of node names sorted by system CPU
// load.
func (r *Redis) SetStats(node string, stats types.Stats) error {
	err := r.ZAdd(string(KeyStatsList), redis.Z{
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

	return r.Set(KeyPrefixNodeStats.Fmt(node), b, 0).Err()
}

// CreateNode creates a node
func (r *Redis) CreateNode(name string) error {
	return r.SAdd(string(KeyNodes), name).Err()
}

// DeleteNode deletes a node
func (r *Redis) DeleteNode(name string) error {
	return r.SRem(string(KeyNodes), name).Err()
}

// ClaimPlayer claims a player for the node
func (r *Redis) ClaimPlayer(node string, guildID uint64) (bool, error) {
	nodes, err := r.SMembers(string(KeyNodes)).Result()
	if err != nil {
		return false, err
	}

	return r.Eval(`for local i = 1, #KEYS-1 do
	local has = redis.call("sismember", KEYS[i], ARGV[1])
	if has then return KEYS[i] == KEYS[#KEYS] end
end
redis.call("sadd", KEYS[#KEYS], ARGV[1])
return true`, append(nodes, node), guildID).Bool()
}

// ReleasePlayer releases a player from a node
func (r *Redis) ReleasePlayer(node string, guildID uint64) error {
	return r.SRem(node, guildID).Err()
}

// AnnounceDeath destroys a node
func (r *Redis) AnnounceDeath(node string) error {
	return r.Publish(string(KeyNodeDeaths), node).Err()
}

// ConsumeDeaths consumes death notifications, ignoring the specified node
func (r *Redis) ConsumeDeaths(node string) error {
	pubsub := r.Subscribe(string(KeyNodeDeaths))
	defer pubsub.Close()

	for msg := range pubsub.Channel() {
		if msg.Payload == node {
			continue
		}

		for {
			player, err := r.SPop(KeyPrefixNodePlayers.Fmt(msg.Payload)).Result()
			if err != nil {
				return err
			}
			if player == "" {
				break
			}

			err = r.Eval(`redis.call('sadd', KEYS[1], ARGV[1])
redis.call('set', KEYS[2], ARGV[2])`, []string{
				KeyPrefixNodePlayers.Fmt(node),
				KeyPrefixPlayerNode.Fmt(player),
			}, player, node).Err()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
