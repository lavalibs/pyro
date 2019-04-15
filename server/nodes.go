package server

import (
	"fmt"
	"net/http"

	"github.com/appellation/fasthttputil"
	"github.com/lavalibs/pyro/lavalink/types"
	"github.com/valyala/fasthttp"
)

// GetNode gets the given node
func (s *Server) GetNode(ctx *fasthttp.RequestCtx) {
	nodeID := ctx.UserValue("nodeID").(string)
	stats := &types.Stats{}
	err := s.Cache.GetStats(nodeID, stats)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	if stats == nil {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	err = fasthttputil.SendJSON(ctx, stats)
	if err != nil {
		fmt.Println(err)
	}
}
