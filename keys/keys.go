package keys

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
	StatsList  Key = iota // single set of node names sorted by CPU usage
	Nodes                 // set of node names
	NodeDeaths            // pubsub channel where node deaths are announced
)

// Constants for Redis key prefixes for different data
const (
	PrefixPlayerState    Key = iota // k=guild, v=JSON state
	PrefixVoiceSession              // k=guild, v=session ID
	PrefixVoiceServer               // k=guild, v=JSON voice server
	PrefixNodeStats                 // k=node, v=JSON stats
	PrefixNodePlayers               // k=node, v[set]=guild
	PrefixNodePackets               // k=node, v[pubsub]=JSON Lavalink ws packet
	PrefixPlayerQueue               // k=guild, v[list]=track identifier
	PrefixPlayerPrevious            // k=guild, v[list]=track identifier
)
