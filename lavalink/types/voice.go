package types

// VoiceStateUpdate represents a voice state update packet from Discord
type VoiceStateUpdate struct {
	GuildID   uint64 `json:"guild_id,string"`
	ChannelID uint64 `json:"channel_id,string"`
	UserID    uint64 `json:"user_id,string"`
	SessionID string `json:"session_id,string"`
	Deaf      bool   `json:"deaf"`
	Mute      bool   `json:"mute"`
	SelfDeaf  bool   `json:"self_deaf"`
	SelfMute  bool   `json:"self_mute"`
	Suppress  bool   `json:"suppress"`
}

// VoiceServerUpdate represnts a voice server update packet from Discord
type VoiceServerUpdate struct {
	GuildID  uint64 `json:"guild_id,string"`
	Token    string `json:"token"`
	Endpoint string `json:"endpoint"`
}
