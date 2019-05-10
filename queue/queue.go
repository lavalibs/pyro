package queue

// Queue represents a store of songs
type Queue interface {
	Add(guildID uint64, tracks []string) error
	Set(guildID uint64, tracks []string) error
	Put(guildID uint64, tracks map[int]string) error
	Unshift(guildID uint64, tracks ...string) error
	Remove(guildID uint64, index int) error
	Next(guildID uint64, count int) ([]string, error)
	Move(guildID uint64, from, to int) error
	Shuffle(guildID uint64) ([]string, error)
	Splice(guildID uint64, start, deleteCount int, tracks ...string) ([]string, error)
	Trim(guildID uint64, start, end int) error
	NowPlaying(guildID uint64) (string, error)
	List(guildID uint64, index int, count uint) ([]string, error)
}
