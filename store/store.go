package store

import "github.com/lavalibs/pyro/lavalink/types"

// Lavalink represents the required interface for a complete Lavalink Store
type Lavalink interface {
	GetPlayer(guildID uint64, state *types.PlayerState) error
	SetPlayer(upd types.PlayerUpdate) error
	GetVoiceUpdate(guildID uint64) (sessionID string, event types.VoiceServerUpdate, err error)
	SetVoiceSession(guildID uint64, sessionID string) error
	SetVoiceServer(pk types.VoiceServerUpdate) error
	GetStats(node string, stats *types.Stats) error
	SetStats(node string, stats types.Stats) error
}

// LavalinkCluster is the interface a backend storage system must implement to be used as a cluster Store
type LavalinkCluster interface {
	Lavalink

	// claim a player for a specified node
	ClaimPlayer(node string, guildID uint64) (bool, error)

	// get the node of a player
	PlayerNode(guildID uint64) (string, error)

	// release a player from any claims on it
	ReleasePlayer(node string, guildID uint64) (bool, error)

	// announce the destruction of a node and the subsequent availability of all its players
	AnnounceDeath(node string) error

	// consume death notifications for a node should pull players from the node's set of
	// players until there are no remaining players
	ConsumeDeaths(node string) error
}
