package types

// LoadType is the type of track that was loaded
type LoadType = string

// Valid load types
const (
	LoadTypeTrackLoaded    LoadType = "TRACK_LOADED"
	LoadTypePlaylistLoaded LoadType = "PLAYLIST_LOADED"
	LoadTypeSearchResult   LoadType = "SEARCH_RESULT"
	LoadTypeNoMatches      LoadType = "NO_MATCHES"
	LoadTypeLoadFailed     LoadType = "LOAD_FAILED"
)

// Op is the type of packet received/to send
type Op = string

// Valid Op codes
const (
	OpVoiceUpdate  Op = "voiceUpdate"
	OpPlay            = "play"
	OpStop            = "stop"
	OpPause           = "pause"
	OpSeek            = "seek"
	OpVolume          = "volume"
	OpEqualizer       = "equalizer"
	OpDestroy         = "destroy"
	OpPlayerUpdate    = "playerUpdate"
	OpStats           = "stats"
	OpEvent           = "event"
)

// TrackResponse represents a response from the Lavalink rest API
type TrackResponse struct {
	LoadType     LoadType     `json:"loadType"`
	PlaylistInfo PlaylistInfo `json:"playlistInfo"`
	Tracks       []Track      `json:"tracks"`
}

// PlaylistInfo represents playlist info loaded by Lavalink
type PlaylistInfo struct {
	Name          string `json:"name"`
	SelectedTrack int    `json:"selectedTrack"`
}

// Track represents a track as sent from the Lavalink rest API
type Track struct {
	Track string    `json:"track"`
	Info  TrackInfo `json:"info"`
}

// TrackInfo represents information about a track
type TrackInfo struct {
	Identifier string `json:"identifier"`
	IsSeekable bool   `json:"isSeekable"`
	Author     string `json:"author"`
	Length     int    `json:"length"`
	IsStream   bool   `json:"isStream"`
	Position   int    `json:"position"`
	Title      string `json:"title"`
	URI        string `json:"uri"`
}

// BasePacket represents a basic Lavalink packet with no Op-specific data
type BasePacket struct {
	Op string `json:"op"`
}

// VoiceUpdate represents a voiceUpdate packet sent to Lavalink
type VoiceUpdate struct {
	Op        string            `json:"op"`
	GuildID   uint64            `json:"guildId,string"`
	SessionID string            `json:"session_id"`
	Event     VoiceServerUpdate `json:"event"`
}

// Play represents a play packet sent to Lavalink
type Play struct {
	Op        string `json:"op"`
	GuildID   uint64 `json:"guildId,string"`
	Track     string `json:"track"`
	StartTime int    `json:"startTime,omitempty"`
	EndTime   int    `json:"endTime,omitempty"`
}

// Stop represents a stop packet sent to Lavalink
type Stop struct {
	Op      string `json:"op"`
	GuildID uint64 `json:"guildId,string"`
}

// Pause represents a pause packet sent to Lavalink
type Pause struct {
	Op      string `json:"op"`
	GuildID uint64 `json:"guildId,string"`
	Pause   bool   `json:"pause"`
}

// Seek represents a seek packet sent to Lavalink
type Seek struct {
	Op       string `json:"op"`
	GuildID  uint64 `json:"guildId,string"`
	Position int    `json:"position"`
}

// Volume represents a volume packet sent to Lavalink
type Volume struct {
	Op      string `json:"op"`
	GuildID uint64 `json:"guildId,string"`
	Volume  int    `json:"volume"`
}

// Equalizer represents an equalizer packet sent to Lavalink
type Equalizer struct {
	Op      string          `json:"op"`
	GuildID uint64          `json:"guildId,string"`
	Bands   []EqualizerBand `json:"bands"`
}

// EqualizerBand describes the format of bands sent to lavalink
type EqualizerBand struct {
	Band int     `json:"band"`
	Gain float64 `json:"gain"`
}

// Destroy represents a destroy packet sent to Lavalink
type Destroy struct {
	Op      string `json:"op"`
	GuildID uint64 `json:"guildId,string"`
}

// Stats represents node stats received from Lavalink
type Stats struct {
	Players        int         `json:"players"`
	PlayingPlayers int         `json:"playingPlayers"`
	Uptime         int         `json:"uptime"`
	Memory         StatsMemory `json:"memory"`
	CPU            StatsCPU    `json:"cpu"`
	FrameStats     StatsFrames `json:"frameStats"`
}

type StatsMemory struct {
	Free       int `json:"free"`
	Used       int `json:"used"`
	Allocated  int `json:"allocated"`
	Reservable int `json:"reservable"`
}

type StatsCPU struct {
	Cores        int     `json:"cores"`
	SystemLoad   float64 `json:"systemLoad"`
	LavalinkLoad float64 `json:"lavalinkLoad"`
}

type StatsFrames struct {
	Sent    int `json:"sent"`
	Nulled  int `json:"nulled"`
	Deficit int `json:"deficit"`
}

// PlayerUpdate represents a player update received from Lavalink
type PlayerUpdate struct {
	GuildID uint64      `json:"guildId,string"`
	State   PlayerState `json:"state"`
}

// PlayerState represents the state of a player
type PlayerState struct {
	Time     int `json:"time"`
	Position int `json:"position"`
}

// Event represents a player event emitted from lavalink
type Event struct {
	Type    string `json:"type"`
	GuildID uint64 `json:"guildId,string"`
	Track   string `json:"track"`
	Reason  string `json:"reason"`
	Error   string `json:"error"`
}
