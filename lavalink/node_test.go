package lavalink

import (
	"testing"
)

func TestNode(t *testing.T) {
	n := New(NodeOptions{
		RestEndpoint: "http://localhost:8081",
		WSEndpoint:   "ws://localhost:8081",
		Password:     "youshallnotpass",
		UserID:       "218844420613734401",
		ShardCount:   "1",
	})

	err := n.Connect()
	if err != nil {
		t.Fatal(err)
		return
	}

	var res interface{}
	res, err = n.LoadTrack("ytsearch:monstercat")
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log(res)

	res, err = n.DecodeTrack("QAAAogIANVtFbGVjdHJvXSAtIE5pdHJvIEZ1biAtIE5ldyBHYW1lIFtNb25zdGVyY2F0IFJlbGVhc2VdABNNb25zdGVyY2F0OiBVbmNhZ2VkAAAAAAAD96AACzZ5X05KZy14b2VFAAEAK2h0dHBzOi8vd3d3LnlvdXR1YmUuY29tL3dhdGNoP3Y9NnlfTkpnLXhvZUUAB3lvdXR1YmUAAAAAAAAAAA==")
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log(res)
}
