package lavalink

import "github.com/lavalibs/pyro/lavalink/types"

// Cache represents the required interface for a complete Lavalink cache
type Cache interface {
	PlayerCache
	VoiceUpdateCache
	StatsCache
}

// PlayerCache represents a cache of player state, keyed by guild ID
type PlayerCache interface {
	GetPlayer(guildID uint64) (types.PlayerState, error)
	SetPlayer(upd types.PlayerUpdate) error
}

// VoiceUpdateCache represents a cache of voice data to be potentially sent to Lavalink
type VoiceUpdateCache interface {
	GetVoiceUpdate(guildID uint64) (sessionID string, event types.VoiceServerUpdate, err error)
	SetVoiceState(pk types.VoiceStateUpdate) error
	SetVoiceServer(pk types.VoiceServerUpdate) error
}

// StatsCache represents a cache of stats, keyed by an arbitrary node ID
type StatsCache interface {
	GetStats(node string) (types.Stats, error)
	SetStats(node string, stats types.Stats) error
}

// ClusterCache is the interface a backend storage system must implement to be used as a cluster cache
type ClusterCache interface {
	Cache

	// claim a player for a specified node
	ClaimPlayer(node string, guildID uint64) (bool, error)

	// release a player from any claims on it
	ReleasePlayer(node string, guildID uint64) error

	// announce the destruction of a node and the subsequent availability of all its players
	AnnounceDeath(node string) error

	// consume death notifications for a node should pull players from the node's set of
	// players until there are no remaining players
	ConsumeDeaths(node string) error
}
