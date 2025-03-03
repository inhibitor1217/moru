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

func (m *membership) Discover(peer Peer) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	var newPeer = true
	if old, ok := m.peers[peer.ID]; ok {
		if old.SessionID == peer.SessionID &&
			old.ExpireAt.After(time.Now()) {
			newPeer = false
		}
	}

	m.peers[peer.ID] = peer
	return newPeer
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
	for id, peer := range m.peers {
		if peer.ExpireAt.Before(time.Now()) {
			delete(m.peers, id)
		}
	}
}
