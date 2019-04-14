package node

import (
	"github.com/lavalibs/pyro/lavalink"
	"github.com/valyala/fasthttp"
)

// Node represents a proxy between Lavalink and Redis
type Node struct {
	Connection *lavalink.Connection
	Cache      lavalink.ClusterCache
}

// HandleSend handles a request intended to be sent to the Lavalink websocket; clients are
// responsible for type checking the post body before sending it
func (n *Node) HandleSend(ctx *fasthttp.RequestCtx) {
	_, err := n.Connection.Write(ctx.PostBody())
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}
