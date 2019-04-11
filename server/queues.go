package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/appellation/fasthttputil"
	"github.com/valyala/fasthttp"
)

// GetSongs gets songs for the given guild
func (s *Server) GetSongs(ctx *fasthttp.RequestCtx) {
	tracks, err := s.Queue.List(ctx.UserValue("guildID").(string), 0, -1)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	err = fasthttputil.SendJSON(ctx, tracks)
}

// PatchSongs adds songs to the given queue
func (s *Server) PatchSongs(ctx *fasthttp.RequestCtx) {
	var tracks map[int]string
	err := json.Unmarshal(ctx.PostBody(), &tracks)
	if err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}

	queue, err := s.Queue.Add(ctx.UserValue("guildID").(string), tracks)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	err = fasthttputil.SendJSON(ctx, queue)
}

// PutSongs sets the songs for a given queue
func (s *Server) PutSongs(ctx *fasthttp.RequestCtx) {
	var tracks []string
	err := json.Unmarshal(ctx.PostBody(), &tracks)
	if err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}

	queue, err := s.Queue.Set(ctx.UserValue("guildID").(string), tracks)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(ctx).Encode(queue)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
}

// GetSong gets the given song from the given queue
func (s *Server) GetSong(ctx *fasthttp.RequestCtx) {
	var (
		pos   = ctx.UserValue("songPosition").(string)
		guild = ctx.UserValue("guildID").(string)
	)

	posConv, err := strconv.ParseInt(pos, 10, 32)
	if err != nil {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	track, err := s.Queue.Get(guild, int(posConv))
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	if track == "" {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	song, err := s.Node.DecodeTrack(track)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(ctx).Encode(song)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
}
