// Package server provides the stateless container for Lavalink.
package server

import (
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/lavalibs/pyro/cache"
	"github.com/lavalibs/pyro/lavalink"
	"github.com/valyala/fasthttp"
)

// Server represents a server to front a single Lavalink node. Multiple servers must be launched to
// handle multiple nodes.
type Server struct {
	Node *lavalink.Node
	Cache cache.Cache
	Queue cache.Queue
}

// Serve data from this server
func (s *Server) Serve(addr string) {
	router := fasthttprouter.New()

	router.GET("/nodes/:nodeID", s.GetNode)

	router.GET("/players/:guildID", s.GetPlayer)
	router.PUT("/players/:guildID", s.PutPlayer)

	router.GET("/queues/:guildID/songs", s.GetSongs)
	router.PUT("/queues/:guildID/songs", s.PutSongs)
	router.PATCH("/queues/:guildID/songs", s.PatchSongs)

	router.GET("/queues/:guildID/songs/:songPosition", s.GetSong)

	log.Fatal(fasthttp.ListenAndServe(addr, router.Handler))
}
