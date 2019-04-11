package lavalink

import (
	"encoding/json"
	"net/http"

	"github.com/lavalibs/pyro/lavalink/types"
)

// ClusterNode represents a node in a cluster of Lavalink nodes
type ClusterNode struct {
	*Node
	*http.ServeMux
	cache ClusterCache
}

// NewClusterNode makes a new cluster node and connects it to a cluster cache
func NewClusterNode(cache ClusterCache, opts NodeOptions) *ClusterNode {
	n := &ClusterNode{New(opts), http.NewServeMux(), cache}
	go cache.ConsumeDeaths(n.Name)
	go n.read()

	return n
}

func (n *ClusterNode) read() (err error) {
	for {
		b := []byte{}
		_, err = n.Read(b)
		if err != nil {
			return
		}

		d := types.BasePacket{}
		err = json.Unmarshal(b, &d)
		if err != nil {
			return
		}

		switch d.OP {
		case "playerUpdate":
			pk := types.PlayerUpdate{}
			err = json.Unmarshal(b, &pk)
			if err != nil {
				return
			}

			err = n.cache.SetPlayer(pk)
			if err != nil {
				return
			}
		case "stats":
			pk := types.Stats{}
			err = json.Unmarshal(b, &pk)
			if err != nil {
				return
			}

			err = n.cache.SetStats(n.Name, pk)
			if err != nil {
				return
			}
		}
	}
}
