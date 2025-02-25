package discovery

import "time"

// PeerID is a unique identifier for a peer on the network.
//
// PeerID should be not forgeable by other peers; in specific, the peer should own the certificate
// with private key which SHA-256 hash is equal to PeerID.
// For now, we skip the certificate part and use a random string as PeerID.
type PeerID string

// Peer is a metadata of a peer on the network.
type Peer struct {
	ID        PeerID    // Unique identifier of the peer device.
	SessionID string    // Randomly generated session ID. We can identify whether the peer process restarted by this.
	Address   string    // IP address of the peer device.
	FoundAt   time.Time // Time when the peer was discovered.
	ExpireAt  time.Time // Time when the peer will be removed from the known peers list.
}
