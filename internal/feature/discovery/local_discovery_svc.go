package discovery

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/inhibitor1217/moru/internal/lib/beacon"
	"github.com/inhibitor1217/moru/proto/discovery"
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
	listenerLoopWg     sync.WaitGroup

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

	s.listenerLoopWg.Add(1)
	go s.listenerLoop(bgCtx)

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
	s.listenerLoopWg.Wait()

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
		if err := s.announce(ctx); err != nil {
			if errors.Is(err, beacon.ErrBeaconStopped) {
				return
			}

			s.log.ErrorContext(ctx, "failed to announce", "err", err)
			time.Sleep(1 * time.Second)
			continue
		}

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

func (s *localDiscoverySvc) announce(ctx context.Context) error {
	pkt := announcementPacket(s.myID, s.mySessionID, "") // TODO fill out TCP address
	return s.beacon.Send(ctx, pkt)
}

func (s *localDiscoverySvc) listenerLoop(ctx context.Context) {
	s.log.InfoContext(ctx, "starting discovery listener loop")
	defer s.log.InfoContext(ctx, "stopping discovery listener loop")

	defer s.listenerLoopWg.Done()

	for {
		select {
		case <-s.stop:
			return
		case <-ctx.Done():
			return
		default:
		}

		cancelCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()

		pkt, remoteAddress, err := s.beacon.Receive(cancelCtx)
		if err != nil {
			if errors.Is(err, beacon.ErrBeaconStopped) {
				return
			}

			if errors.Is(err, context.DeadlineExceeded) {
				continue
			}

			s.log.ErrorContext(ctx, "failed to receive packet", "err", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if remoteAddress == nil {
			continue
		}

		msg, err := parsePacket(pkt)
		if err != nil {
			s.log.WarnContext(ctx, "failed to parse packet", "err", err)
			continue
		}

		if err := s.handleMessage(ctx, msg); err != nil {
			s.log.ErrorContext(ctx, "failed to handle message", "err", err)
			continue
		}
	}
}

func (s *localDiscoverySvc) handleMessage(ctx context.Context, msg *discovery.Message) error {
	remotePeerID, err := PeerIDFromBytes(msg.Id)
	if err != nil {
		return fmt.Errorf("failed to parse peer ID: %w", err)
	}

	s.log.DebugContext(ctx, "received message",
		"remote.peerID", remotePeerID,
		"remote.sessionID", msg.SessionId)

	// skip broadcasted message from myself
	if remotePeerID == s.myID {
		return nil
	}

	return nil
}
