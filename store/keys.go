package store

import "fmt"

// Key represents a key type
type Key rune

// Fmt creates a key of the given value
func (k Key) Fmt(a interface{}) string {
	return fmt.Sprintf("%v%v", k, a)
}

func (k Key) String() string {
	return string(k)
}

// Constants for Redis keys
const (
	KeyStatsList  Key = iota // single set of node names sorted by CPU usage
	KeyNodes                 // set of node names
	KeyNodeDeaths            // pubsub channel where node deaths are announced
)

// Constants for Redis key prefixes for different data
const (
	KeyPrefixPlayerState    Key = iota // k=guild, v=JSON state
	KeyPrefixVoiceSession              // k=guild, v=session ID
	KeyPrefixVoiceServer               // k=guild, v=JSON voice server
	KeyPrefixNodeStats                 // k=node, v=JSON stats
	KeyPrefixNodePlayers               // k=node, v[set]=guild
	KeyPrefixPlayerNode                // k=guild, v=node
	KeyPrefixNodePackets               // k=node, v[pubsub]=JSON Lavalink ws packet
	KeyPrefixPlayerQueue               // k=guild, v[list]=track identifier
	KeyPrefixPlayerPrevious            // k=guild, v[list]=track identifier
)
