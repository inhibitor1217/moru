package discovery

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/inhibitor1217/moru/internal/lib/beacon"
)

// localDiscoverySvc is a DiscoverySvc that works in LAN.
type localDiscoverySvc struct {
	started bool
	stopped bool
	mu      sync.Mutex

	beacon beacon.Beacon
	log    *slog.Logger
}

// NewLocalDiscoverySvc creates a new DiscoverySvc that works in LAN.
func NewLocalDiscoverySvc(beacon beacon.Beacon) (DiscoverySvc, error) {
	return &localDiscoverySvc{
		beacon: beacon,
		log:    slog.Default().With("source", "discovery.localDiscoverySvc"),
	}, nil
}

func (s *localDiscoverySvc) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		return fmt.Errorf("discovery service already started")
	} else if s.stopped {
		return fmt.Errorf("discovery service already stopped")
	}
	s.started = true

	s.log.Info("starting LAN discovery service")

	return nil
}

func (s *localDiscoverySvc) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.started {
		return fmt.Errorf("discovery service not started")
	} else if s.stopped {
		return fmt.Errorf("discovery service already stopped")
	}
	s.stopped = true

	s.log.Info("stopping LAN discovery service")

	return nil
}

func (s *localDiscoverySvc) KnownPeers() []Peer {
	return nil
}

func (s *localDiscoverySvc) Refresh(ctx context.Context) error {
	return nil
}
