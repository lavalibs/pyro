package store

import "github.com/lavalibs/pyro/lavalink/types"

// Lavalink represents the required interface for a complete Lavalink Store
type Lavalink interface {
	GetPlayer(guildID uint64, state *types.PlayerState) error
	SetPlayer(upd types.PlayerUpdate) error
	GetVoiceUpdate(guildID uint64) (sessionID string, event types.VoiceServerUpdate, err error)
	SetVoiceState(pk types.VoiceStateUpdate) error
	SetVoiceServer(pk types.VoiceServerUpdate) error
	GetStats(node string, stats *types.Stats) error
	SetStats(node string, stats types.Stats) error
}

// Queue represents a store of songs
type Queue interface {
	Add(guildID uint64, tracks map[int]string) error
	Set(guildID uint64, tracks []string) error
	Unshift(guildID uint64, tracks ...string) error
	Remove(guildID uint64, index int) error
	Next(guildID uint64, count int) ([]string, error)
	Sort(uint64, func(a, b string) int) ([]string, error)
	Move(guildID uint64, from, to int) error
	Shuffle(guildID uint64) ([]string, error)
	Splice(guildID uint64, start, deleteCount int, tracks ...string) ([]string, error)
	Trim(guildID uint64, start, end int) error
	NowPlaying(guildID uint64) (string, error)
	List(guildID uint64, index int, count uint) ([]string, error)
}

// LavalinkCluster is the interface a backend storage system must implement to be used as a cluster Store
type LavalinkCluster interface {
	Lavalink

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
