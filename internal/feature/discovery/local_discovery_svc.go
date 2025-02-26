package discovery

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/inhibitor1217/moru/internal/lib/beacon"
)

const (
	BroadcastInterval = 30 * time.Second
)

// localDiscoverySvc is a DiscoverySvc that works in LAN.
type localDiscoverySvc struct {
	myID        PeerID
	mySessionID int64

	started            bool
	stopped            bool
	mu                 sync.Mutex
	stop               chan struct{}
	announcementLoopWg sync.WaitGroup

	announcementTick   <-chan time.Time
	forcedAnnouncement chan struct{}

	beacon beacon.Beacon
	log    *slog.Logger
}

// NewLocalDiscoverySvc creates a new DiscoverySvc that works in LAN.
func NewLocalDiscoverySvc(beacon beacon.Beacon) (DiscoverySvc, error) {
	myID := randomPeerID()
	mySessionID := randomSessionID()

	return &localDiscoverySvc{
		myID:        myID,
		mySessionID: mySessionID,

		stop: make(chan struct{}),

		announcementTick:   time.Tick(BroadcastInterval),
		forcedAnnouncement: make(chan struct{}),

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

	s.log.Info("starting LAN discovery service",
		"myID", s.myID,
		"mySessionID", s.mySessionID)

	bgCtx := context.WithoutCancel(ctx)

	s.announcementLoopWg.Add(1)
	go s.announcementLoop(bgCtx)

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
	close(s.stop)

	s.announcementLoopWg.Wait()

	return nil
}

func (s *localDiscoverySvc) KnownPeers() []Peer {
	return nil
}

func (s *localDiscoverySvc) Refresh(ctx context.Context) error {
	s.forcedAnnouncement <- struct{}{}
	return nil
}

func (s *localDiscoverySvc) announcementLoop(ctx context.Context) {
	s.log.InfoContext(ctx, "starting discovery announcement loop")
	defer s.log.InfoContext(ctx, "stopping discovery announcement loop")

	defer s.announcementLoopWg.Done()

	for {
		s.announce(ctx)

		select {
		case <-s.announcementTick:
		case <-s.forcedAnnouncement:
		case <-s.stop:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (s *localDiscoverySvc) announce(ctx context.Context) {
	pkt := announcementPacket(s.myID, s.mySessionID, "") // TODO fill out TCP address
	if err := s.beacon.Send(ctx, pkt); err != nil {
		s.log.ErrorContext(ctx, "failed to send announcement packet", "err", err)
	}
}
