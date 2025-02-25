package discovery

import "context"

// DiscoverySvc manages the known peers on the network.
//
// DiscoverySvc is responsible for discovering new peers and keeping track of known peers.
// It broadcasts and listens for peer discovery messages on the local network.
type DiscoverySvc interface {
	// Start begins discovering new peers and keeping track of known peers.
	// It returns an error if the service fails to start.
	Start(ctx context.Context) error

	// Stop gracefully shuts down the discovery service.
	Stop(ctx context.Context) error

	// KnownPeers returns the list of known peers.
	KnownPeers() []Peer

	// Refresh tries to refresh the known peers list up-to-date.
	// Other peers will echo the announcement message with broadcast, and hopefully
	// this peer will receive the message and update the known peers list.
	Refresh(ctx context.Context) error
}
