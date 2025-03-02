package discovery

import (
	"encoding/binary"
	"errors"

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

func announcementPacket(peer Peer) []byte {
	return makePacket(&discovery.Message{
		Id:        peer.ID[:],
		SessionId: peer.SessionID,
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
