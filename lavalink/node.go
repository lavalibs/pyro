// Package lavalink is a wrapper for the Lavalink audio client for Discord.
package lavalink

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/lavalibs/pyro/lavalink/types"
)

// NodeOptions represents node configuration information
type NodeOptions struct {
	RestEndpoint,
	WSEndpoint,
	Password,
	UserID,
	ShardCount,
	Name string
}

// Node represents a Lavalink node
type Node struct {
	NodeOptions
	Dialer *websocket.Dialer
	Client *http.Client

	conn *websocket.Conn
	rmux sync.Mutex
	wmux sync.Mutex
}

// New makes a new node
func New(options NodeOptions) *Node {
	return &Node{
		NodeOptions: options,
		Dialer:      websocket.DefaultDialer,
		Client:      http.DefaultClient,

		rmux: sync.Mutex{},
		wmux: sync.Mutex{},
	}
}

// LoadTrack from the given endpoint
func (n *Node) LoadTrack(identifier string) (track types.TrackResponse, err error) {
	res, err := n.DoJSON(http.MethodGet, "/loadtracks", nil, func(q url.Values) {
		q.Add("identifier", identifier)
	})
	if err != nil {
		return
	}

	err = json.NewDecoder(res.Body).Decode(&track)
	res.Body.Close()
	return
}

// DecodeTrack from the given endpoint
func (n *Node) DecodeTrack(identifier string) (track types.TrackInfo, err error) {
	res, err := n.DoJSON(http.MethodGet, "/decodetrack", nil, func(q url.Values) {
		q.Add("track", identifier)
	})
	if err != nil {
		return
	}

	err = json.NewDecoder(res.Body).Decode(&track)
	res.Body.Close()
	return
}

// DecodeTracks from the given endpoint
func (n *Node) DecodeTracks(identifiers []string) (tracks []types.TrackInfo, err error) {
	res, err := n.DoJSON(http.MethodPost, "/decodetracks", &identifiers, nil)
	if err != nil {
		return
	}

	err = json.NewDecoder(res.Body).Decode(&tracks)
	res.Body.Close()
	return
}

// Connect this node to the given node
func (n *Node) Connect() error {
	headers := http.Header{}
	headers.Add("Authorization", n.Password)
	headers.Add("Num-Shards", n.ShardCount)
	headers.Add("User-Id", n.UserID)

	conn, _, err := n.Dialer.Dial(n.WSEndpoint, headers)
	n.conn = conn
	return err
}

// Play a track
func (n *Node) Play(guildID uint64, track string, start, end int) error {
	return n.Send(types.Play{
		OP:        "play",
		GuildID:   guildID,
		Track:     track,
		StartTime: start,
		EndTime:   end,
	})
}

// VoiceUpdate sends a voiceUpdate packet to Lavalink
func (n *Node) VoiceUpdate(guildID uint64, sessionID string, event types.VoiceServerUpdate) error {
	return n.Send(types.VoiceUpdate{
		OP:        "voiceUpdate",
		GuildID:   guildID,
		SessionID: sessionID,
		Event:     event,
	})
}

// Stop this guild
func (n *Node) Stop(guildID uint64) error {
	return n.Send(types.Stop{
		OP:      "stop",
		GuildID: guildID,
	})
}

// Pause this guild
func (n *Node) Pause(guildID uint64, pause bool) error {
	return n.Send(types.Pause{
		OP:      "pause",
		GuildID: guildID,
		Pause:   pause,
	})
}

// Seek to a specific point in the current track
func (n *Node) Seek(guildID uint64, position int) error {
	return n.Send(types.Seek{
		OP:       "seek",
		GuildID:  guildID,
		Position: position,
	})
}

// SetVolume sets the volume
func (n *Node) SetVolume(guildID uint64, volume int) error {
	return n.Send(types.Volume{
		OP:      "volume",
		GuildID: guildID,
		Volume:  volume,
	})
}

// SetEqualizer sets the equalizer
func (n *Node) SetEqualizer(guildID uint64, bands []types.EqualizerBand) error {
	return n.Send(types.Equalizer{
		OP:      "equalizer",
		GuildID: guildID,
		Bands:   bands,
	})
}

// Destroy destroys this player
func (n *Node) Destroy(guildID uint64) error {
	return n.Send(types.Destroy{
		OP:      "destroy",
		GuildID: guildID,
	})
}

func (n *Node) Write(p []byte) (int, error) {
	n.wmux.Lock()
	defer n.wmux.Unlock()

	return len(p), n.conn.WriteMessage(websocket.BinaryMessage, p)
}

// Send JSON to this connection
func (n *Node) Send(d interface{}) error {
	n.wmux.Lock()
	defer n.wmux.Unlock()

	return n.conn.WriteJSON(d)
}

func (n *Node) Read(p []byte) (int, error) {
	n.rmux.Lock()
	defer n.rmux.Unlock()

	_, p, err := n.conn.ReadMessage()
	return len(p), err
}

// ReadJSON into the value pointed at by p
func (n *Node) ReadJSON(p interface{}) error {
	n.rmux.Lock()
	defer n.rmux.Unlock()

	return n.conn.ReadJSON(p)
}

// DoJSON performs a JSON HTTP request to the Lavalink node
func (n *Node) DoJSON(method, path string, d interface{}, genQuery func(q url.Values)) (res *http.Response, err error) {
	var (
		b []byte
		r io.Reader
	)

	if d != nil {
		b, err = json.Marshal(d)
		if err != nil {
			return
		}

		r = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, n.RestEndpoint, r)
	if err != nil {
		return
	}

	req.Header.Add("Authorization", n.Password)
	req.URL.Path = path

	if genQuery != nil {
		q := req.URL.Query()
		genQuery(q)
		req.URL.RawQuery = q.Encode()
	}

	res, err = n.Client.Do(req)
	return
}
