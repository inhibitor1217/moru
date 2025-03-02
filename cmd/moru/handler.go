package main

import (
	"context"
	"log/slog"

	discoverypb "github.com/inhibitor1217/moru/proto/discovery"
	"google.golang.org/protobuf/proto"
)

func (m *moru) knownPeers(
	_ context.Context,
	log *slog.Logger,
	reqBuf []byte,
) []byte {
	req := &discoverypb.KnownPeersRequest{}
	if err := proto.Unmarshal(reqBuf, req); err != nil {
		log.Error("failed to unmarshal known peers request",
			"error", err)
		return nil
	}

	// TODO implement
	res := &discoverypb.KnownPeersResult{
		Peers: []*discoverypb.Peer{
			{
				Id:        []byte("peer1"),
				SessionId: 42,
				Address:   "",
				Username:  nil,
				Hostname:  nil,
				Role:      "peer",
			},
		},
	}
	resBuf, err := proto.Marshal(res)
	if err != nil {
		log.Error("failed to marshal known peers result",
			"error", err)
	}

	return resBuf
}
