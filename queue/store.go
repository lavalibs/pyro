package queue

type Store interface {
	AddTracks(guildID uint64, tracks ...string) error
}
