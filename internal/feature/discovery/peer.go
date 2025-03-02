package discovery

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"strings"
	"time"
	"unsafe"
)

// PeerID is a unique identifier for a peer on the network.
//
// PeerID should be not forgeable by other peers; in specific, the peer should own the certificate
// with private key which SHA-256 hash is equal to PeerID.
// For now, we skip the certificate part and use a random string as PeerID.
type PeerID [32]byte

// PeerIDFromBytes creates a PeerID from the given byte slice.
func PeerIDFromBytes(bs []byte) (PeerID, error) {
	if len(bs) != 32 {
		return PeerID{}, errors.New("invalid PeerID length")
	}
	var id PeerID
	copy(id[:], bs)
	return id, nil
}

func (id PeerID) String() string {
	s := base32.StdEncoding.EncodeToString(id[:])
	s = strings.Trim(s, "=")
	chunked := ""
	for i := 0; i < len(s); i += 8 {
		if i+8 > len(s) {
			chunked += s[i:]
		} else {
			chunked += s[i : i+8]
		}
		if i+8 < len(s) {
			chunked += "-"
		}
	}
	return chunked
}

// Peer is a metadata of a peer on the network.
type Peer struct {
	ID        PeerID    // Unique identifier of the peer device.
	SessionID int64     // Randomly generated session ID. We can identify whether the peer process restarted by this.
	Address   string    // IP address of the peer device.
	Username  *string   // Username of the peer device.
	Hostname  *string   // Hostname of the peer device.
	Role      string    // Role of the peer device.
	FoundAt   time.Time // Time when the peer was discovered.
	ExpireAt  time.Time // Time when the peer will be removed from the known peers list.
}

func randomPeerID() PeerID {
	id := PeerID{}
	_, _ = rand.Read(id[:])
	return id
}

func randomSessionID() int64 {
	var id int64
	_, _ = rand.Read((*[8]byte)(unsafe.Pointer(&id))[:])
	return id
}
