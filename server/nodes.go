package server

import (
	"fmt"
	"net/http"

	"github.com/appellation/fasthttputil"
	"github.com/valyala/fasthttp"
)

// GetNode gets the given node
func (s *Server) GetNode(ctx *fasthttp.RequestCtx) {
	nodeID := ctx.UserValue("nodeID").(string)
	node, err := s.Cache.GetStats(nodeID)
	if err != nil {
		ctx.Error(err.Error(), http.StatusInternalServerError)
		return
	}

	if node == nil {
		ctx.SetStatusCode(http.StatusNotFound)
		return
	}

	err = fasthttputil.SendJSON(ctx, node)
	if err != nil {
		fmt.Println(err)
	}
}
