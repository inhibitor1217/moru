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

	knownPeers := m.discoverySvc.KnownPeers()

	res := &discoverypb.KnownPeersResult{
		Peers: make([]*discoverypb.Peer, 0, len(knownPeers)),
	}
	for _, peer := range knownPeers {
		res.Peers = append(res.Peers, &discoverypb.Peer{
			Id:        peer.ID.Bytes(),
			SessionId: peer.SessionID,
			Address:   peer.Address,
			Username:  peer.Username,
			Hostname:  peer.Hostname,
			Role:      peer.Role,
		})
	}

	resBuf, err := proto.Marshal(res)
	if err != nil {
		log.Error("failed to marshal known peers result",
			"error", err)
	}

	return resBuf
}
