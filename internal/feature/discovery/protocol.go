package discovery

import (
	"encoding/binary"

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

func announcementPacket(
	id PeerID,
	sessionID int64,
	address string,
) []byte {
	return makePacket(&discovery.Message{
		Id:        id[:],
		SessionId: sessionID,
		Payload: &discovery.Message_Announcement{
			Announcement: &discovery.Announcement{
				Address: address,
			},
		},
	})
}
