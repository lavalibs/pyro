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
	guildID, err := strconv.ParseUint(ctx.UserValue("guildID").(string), 10, 64)
	if err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}

	tracks, err := s.Queue.List(guildID, 0, 0)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	err = fasthttputil.SendJSON(ctx, tracks)
}

// PatchSongs adds songs to the given queue
func (s *Server) PatchSongs(ctx *fasthttp.RequestCtx) {
	guildID, err := strconv.ParseUint(ctx.UserValue("guildID").(string), 10, 64)
	if err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}

	var tracks map[int]string
	err = json.Unmarshal(ctx.PostBody(), &tracks)
	if err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}

	err = s.Queue.Put(guildID, tracks)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusNoContent)
}

// PutSongs sets the songs for a given queue
func (s *Server) PutSongs(ctx *fasthttp.RequestCtx) {
	guildID, err := strconv.ParseUint(ctx.UserValue("guildID").(string), 10, 64)
	if err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}

	var tracks []string
	err = json.Unmarshal(ctx.PostBody(), &tracks)
	if err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}

	err = s.Queue.Set(guildID, tracks)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(http.StatusNoContent)
}

// GetSong gets the given song from the given queue
func (s *Server) GetSong(ctx *fasthttp.RequestCtx) {
	guildID, err := strconv.ParseUint(ctx.UserValue("guildID").(string), 10, 64)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	pos, err := strconv.ParseInt(ctx.UserValue("songPosition").(string), 10, 0)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	tracks, err := s.Queue.List(guildID, int(pos), 1)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	if len(tracks) != 1 || tracks[0] == "" {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	song, err := s.HTTP.DecodeTrack(tracks[0])
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
