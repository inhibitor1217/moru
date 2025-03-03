package discovery

import (
	"encoding/binary"
	"errors"
	"time"

	"github.com/inhibitor1217/moru/proto/discovery"
	"google.golang.org/protobuf/proto"
)

const (
	// Magic number used for UDP broadcast message header.
	magic = uint32(0xa39713dd)
)

func makePacket(msg proto.Message) []byte {
	bs, _ := proto.Marshal(msg)
	buf := make([]byte, 4+len(bs))
	binary.BigEndian.PutUint32(buf, magic)
	copy(buf[4:], bs)
	return buf
}

func parsePacket(buf []byte) (*discovery.Message, error) {
	if len(buf) < 4 {
		return nil, errors.New("packet too short")
	}
	if binary.BigEndian.Uint32(buf) != magic {
		return nil, errors.New("invalid magic number")
	}
	msg := &discovery.Message{}
	if err := proto.Unmarshal(buf[4:], msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func announcementPacket(peer Peer, seqnum int64) []byte {
	return makePacket(&discovery.Message{
		Id:        peer.ID[:],
		SessionId: peer.SessionID,
		Seqnum:    seqnum,
		Timestamp: time.Now().UnixMilli(),
		Payload: &discovery.Message_Announcement{
			Announcement: &discovery.Announcement{
				Peer: &discovery.Peer{
					Id:        peer.ID[:],
					SessionId: peer.SessionID,
					Address:   peer.Address,
					Username:  peer.Username,
					Hostname:  peer.Hostname,
					Role:      peer.Role,
				},
			},
		},
	})
}

func helloRequestPacket(peer Peer, seqnum int64) []byte {
	return makePacket(&discovery.Message{
		Id:        peer.ID[:],
		SessionId: peer.SessionID,
		Seqnum:    seqnum,
		Timestamp: time.Now().UnixMilli(),
		Payload: &discovery.Message_HelloRequest{
			HelloRequest: &discovery.HelloRequest{
				Peer: &discovery.Peer{
					Id:        peer.ID[:],
					SessionId: peer.SessionID,
					Address:   peer.Address,
					Username:  peer.Username,
					Hostname:  peer.Hostname,
					Role:      peer.Role,
				},
			},
		},
	})
}

func helloResultPacket(peer Peer, seqnum int64) []byte {
	return makePacket(&discovery.Message{
		Id:        peer.ID[:],
		SessionId: peer.SessionID,
		Seqnum:    seqnum,
		Timestamp: time.Now().UnixMilli(),
		Payload: &discovery.Message_HelloResult{
			HelloResult: &discovery.HelloResult{
				Peer: &discovery.Peer{
					Id:        peer.ID[:],
					SessionId: peer.SessionID,
					Address:   peer.Address,
					Username:  peer.Username,
					Hostname:  peer.Hostname,
					Role:      peer.Role,
				},
			},
		},
	})
}
