package discovery

import (
	"sync"
	"time"
)

// membership manages a list of peers on the network with expiration.
type membership struct {
	peers map[PeerID]Peer
	mu    sync.Mutex
}

// newMembership creates a new membership.
func newMembership() *membership {
	return &membership{
		peers: make(map[PeerID]Peer),
	}
}

func (m *membership) Discover(peer Peer) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.peers[peer.ID] = peer
}

func (m *membership) Peers() []Peer {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.flush()

	peers := make([]Peer, 0, len(m.peers))
	for _, peer := range m.peers {
		peers = append(peers, peer)
	}
	return peers
}

func (m *membership) Size() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.flush()

	return len(m.peers)
}

func (m *membership) flush() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, peer := range m.peers {
		if peer.ExpireAt.Before(time.Now()) {
			delete(m.peers, id)
		}
	}
}
