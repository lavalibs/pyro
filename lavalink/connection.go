package lavalink

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/lavalibs/pyro/lavalink/types"
	"github.com/lavalibs/pyro/store"
)

// Connection represents a WebSocket connection to a Lavalink server
type Connection struct {
	Name string
	ws   *websocket.Conn
	mux  sync.RWMutex
}

// ConnectOptions are the options used when connecting to the Lavalink server
type ConnectOptions struct {
	Name       string
	Endpoint   string
	Password   string
	ShardCount int
	UserID     uint64
	Dialer     *websocket.Dialer
}

// Connect this node to the given node
func Connect(opts ConnectOptions) (*Connection, error) {
	headers := http.Header{}
	headers.Add("Authorization", opts.Password)
	headers.Add("Num-Shards", strconv.FormatInt(int64(opts.ShardCount), 10))
	headers.Add("User-Id", strconv.FormatUint(opts.UserID, 10))

	if opts.Dialer == nil {
		opts.Dialer = websocket.DefaultDialer
	}

	ws, _, err := opts.Dialer.Dial(opts.Endpoint, headers)
	return &Connection{opts.Name, ws, sync.RWMutex{}}, err
}

// Play a track
func (c *Connection) Play(guildID uint64, track string, start, end int) error {
	return c.Send(types.Play{
		Op:        "play",
		GuildID:   guildID,
		Track:     track,
		StartTime: start,
		EndTime:   end,
	})
}

// VoiceUpdate sends a voiceUpdate packet to Lavalink
func (c *Connection) VoiceUpdate(guildID uint64, sessionID string, event types.VoiceServerUpdate) error {
	return c.Send(types.VoiceUpdate{
		Op:        "voiceUpdate",
		GuildID:   guildID,
		SessionID: sessionID,
		Event:     event,
	})
}

// Stop this guild
func (c *Connection) Stop(guildID uint64) error {
	return c.Send(types.Stop{
		Op:      "stop",
		GuildID: guildID,
	})
}

// Pause this guild
func (c *Connection) Pause(guildID uint64, pause bool) error {
	return c.Send(types.Pause{
		Op:      "pause",
		GuildID: guildID,
		Pause:   pause,
	})
}

// Seek to a specific point in the current track
func (c *Connection) Seek(guildID uint64, position int) error {
	return c.Send(types.Seek{
		Op:       "seek",
		GuildID:  guildID,
		Position: position,
	})
}

// SetVolume sets the volume
func (c *Connection) SetVolume(guildID uint64, volume int) error {
	return c.Send(types.Volume{
		Op:      "volume",
		GuildID: guildID,
		Volume:  volume,
	})
}

// SetEqualizer sets the equalizer
func (c *Connection) SetEqualizer(guildID uint64, bands []types.EqualizerBand) error {
	return c.Send(types.Equalizer{
		Op:      "equalizer",
		GuildID: guildID,
		Bands:   bands,
	})
}

// Destroy destroys this player
func (c *Connection) Destroy(guildID uint64) error {
	return c.Send(types.Destroy{
		Op:      "destroy",
		GuildID: guildID,
	})
}

func (c *Connection) Write(p []byte) (int, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	return len(p), c.ws.WriteMessage(websocket.BinaryMessage, p)
}

// Send JSON to this connection
func (c *Connection) Send(d interface{}) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.ws.WriteJSON(d)
}

func (c *Connection) Read(p []byte) (int, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	_, p, err := c.ws.ReadMessage()
	return len(p), err
}

// ReadJSON into the value pointed at by p
func (c *Connection) ReadJSON(p interface{}) error {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.ws.ReadJSON(p)
}

// Store consumes events from Lavalink and stores them in the given cache. Since this reads events
// from the WebSocket, no other goroutines should attempt to read.
func (c *Connection) Store(store store.Lavalink) (err error) {
	d := &types.BasePacket{}
	b := []byte{}
	for {
		_, err = c.Read(b)
		if err != nil {
			return
		}

		err = json.Unmarshal(b, d)
		if err != nil {
			return
		}

		switch d.Op {
		case types.OpPlayerUpdate:
			pk := types.PlayerUpdate{}
			err = json.Unmarshal(b, &pk)
			if err != nil {
				return
			}

			err = store.SetPlayer(pk)
			if err != nil {
				return
			}
		case types.OpStats:
			pk := types.Stats{}
			err = json.Unmarshal(b, &pk)
			if err != nil {
				return
			}

			err = store.SetStats(c.Name, pk)
			if err != nil {
				return
			}
		}
	}
}
