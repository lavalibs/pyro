package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lavalibs/pyro/lavalink/types"
	"github.com/valyala/fasthttp"
)

// PlayerAction represents an action to be taken when modifying the playing state of a player
type PlayerAction struct {
	Action string `json:"action"`
	Track  string `json:"track,omitempty"`
	Paused bool   `json:"paused,omitempty"`
}

// GetPlayer gets the player for the given guild
func (s *Server) GetPlayer(ctx *fasthttp.RequestCtx) {
	guildID, err := strconv.ParseUint(ctx.UserValue("guildID").(string), 64, 10)
	player := &types.PlayerState{}
	err = s.Cache.GetPlayer(guildID, player)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	if player == nil {
		ctx.NotFound()
		return
	}

	err = json.NewEncoder(ctx).Encode(player)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
	}

	ctx.SetContentType("application/json")
}

// PutPlayer modifies the playing state of the player
func (s *Server) PutPlayer(ctx *fasthttp.RequestCtx) {
	var (
		data     = &PlayerAction{}
		guildStr = ctx.UserValue("guildID").(string)
		err      error
		guild    uint64
	)

	guild, err = strconv.ParseUint(guildStr, 10, 64)
	if err != nil {
		ctx.NotFound()
		return
	}

	err = json.Unmarshal(ctx.PostBody(), data)
	if err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}

	switch data.Action {
	case "play":
		if data.Track == "" {
			data.Track, err = s.Queue.NowPlaying(guild)
			if err != nil {
				ctx.Error(err.Error(), http.StatusInternalServerError)
				return
			}
		}

		err = s.Conn.Play(guild, data.Track, 0, 0)
		if err != nil {
			ctx.Error(err.Error(), http.StatusInternalServerError)
			return
		}
	case "pause":
		err = s.Conn.Pause(guild, data.Paused)
		if err != nil {
			ctx.Error(err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
