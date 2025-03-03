package discovery

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/user"
	"sync"
	"time"

	"github.com/inhibitor1217/moru/internal/env"
	"github.com/inhibitor1217/moru/internal/lib/beacon"
	"github.com/inhibitor1217/moru/internal/lib/network"
	"github.com/inhibitor1217/moru/proto/discovery"
	"github.com/samber/lo"
)

const (
	BroadcastInterval = 30 * time.Second
	AnnouncementTTL   = 3 * BroadcastInterval

	HelloInterval = 30 * time.Second
	HelloTTL      = 3 * HelloInterval
)

// localDiscoverySvc is a DiscoverySvc that works in LAN.
type localDiscoverySvc struct {
	me         Peer
	packetSeq  int64
	membership *membership

	mu                 sync.Mutex
	started            bool
	stopped            bool
	stop               chan struct{}
	announcementLoopWg sync.WaitGroup
	listenerLoopWg     sync.WaitGroup

	announcementTick   <-chan time.Time
	forcedAnnouncement chan struct{}

	beacon beacon.Beacon
	log    *slog.Logger
}

// NewLocalDiscoverySvc creates a new DiscoverySvc that works in LAN.
func NewLocalDiscoverySvc(
	beacon beacon.Beacon,
	cfg *env.Config,
) (DiscoverySvc, error) {
	me := Peer{}
	listeningIPs := network.LANAddresses()

	var advertisedIP net.IP
	var advertisedHostURL *string

	if len(listeningIPs) > 0 {
		advertisedIP = listeningIPs[0]
		// TODO support HTTPS via self-signed certificate
		if cfg.Application.Role == env.RoleHost {
			advertisedHostURL = lo.ToPtr(fmt.Sprintf("http://%s:%d", advertisedIP, cfg.HTTP.Port))
		}
	}

	me.ID = randomPeerID()
	me.SessionID = randomSessionID()
	me.Address = advertisedIP
	me.Role = cfg.Application.Role.String()
	if osUser, err := user.Current(); err == nil {
		me.Username = &osUser.Username
	}
	if hostname, err := os.Hostname(); err == nil {
		me.Hostname = &hostname
	}
	me.HostURL = advertisedHostURL
	me.FoundAt = time.Now()
	me.ExpireAt = time.Now().Add(AnnouncementTTL)

	membership := newMembership()
	membership.Discover(me)

	return &localDiscoverySvc{
		me:         me,
		membership: membership,

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
		"me.id", s.me.ID,
		"me.sessionID", s.me.SessionID,
		"me.address", s.me.Address,
		"me.username", lo.FromPtr(s.me.Username),
		"me.hostname", lo.FromPtr(s.me.Hostname),
		"me.role", s.me.Role)

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
	return s.membership.Peers()
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
	for range 5 {
		pkt := announcementPacket(s.me, s.packetSeq)
		s.packetSeq++
		if err := s.beacon.Send(ctx, pkt, beacon.SendBroadcast); err != nil {
			return err
		}
	}
	return nil
}

func (s *localDiscoverySvc) sayHello(ctx context.Context, peer Peer) error {
	pkt := helloRequestPacket(s.me, s.packetSeq)
	s.packetSeq++
	return s.beacon.Send(ctx, pkt, beacon.SendUnicast(peer.Address))
}

func (s *localDiscoverySvc) replyHello(ctx context.Context, peer Peer) error {
	pkt := helloResultPacket(s.me, s.packetSeq)
	s.packetSeq++
	return s.beacon.Send(ctx, pkt, beacon.SendUnicast(peer.Address))
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

		pkt, remoteAddress, err := s.beacon.Receive(ctx)
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

		if err := s.handleMessage(ctx, remoteAddress, msg); err != nil {
			s.log.ErrorContext(ctx, "failed to handle message", "err", err)
			continue
		}
	}
}

func (s *localDiscoverySvc) handleMessage(ctx context.Context, remoteAddr net.Addr, msg *discovery.Message) error {
	remotePeerID, err := PeerIDFromBytes(msg.Id)
	if err != nil {
		return fmt.Errorf("failed to parse peer ID: %w", err)
	}

	s.log.DebugContext(ctx, "received message",
		"remote.address", remoteAddr,
		"remote.peerID", remotePeerID,
		"remote.sessionID", msg.SessionId,
		"seqnum", msg.Seqnum,
		"timestamp", time.UnixMilli(msg.Timestamp))

	// try to parse the remote address (ip:port)
	var remoteIP string
	if ip, _, err := net.SplitHostPort(remoteAddr.String()); err == nil {
		remoteIP = ip
	}

	switch payload := msg.Payload.(type) {
	case *discovery.Message_Announcement:
		peer := Peer{
			ID:        remotePeerID,
			SessionID: msg.SessionId,
			Address:   net.ParseIP(payload.Announcement.Peer.Address),
			Username:  payload.Announcement.Peer.Username,
			Hostname:  payload.Announcement.Peer.Hostname,
			Role:      payload.Announcement.Peer.Role,
			HostURL:   payload.Announcement.Peer.HostUrl,
			FoundAt:   time.Now(),
			ExpireAt:  time.Now().Add(AnnouncementTTL),
		}

		// best effort to fill peer address
		if peer.Address == nil {
			peer.Address = net.ParseIP(remoteIP)
		}

		newPeer := s.membership.Discover(peer)

		s.log.DebugContext(ctx, "discovered peer",
			"peer.id", peer.ID,
			"peer.address", peer.Address,
			"peer.username", lo.FromPtr(peer.Username),
			"peer.hostname", lo.FromPtr(peer.Hostname),
			"peer.role", peer.Role,
			"membership.size", s.membership.Size())

		if newPeer {
			_ = s.Refresh(ctx)
			_ = s.sayHello(ctx, peer)
		}
	case *discovery.Message_HelloRequest:
		if remotePeerID == s.me.ID {
			return nil
		}

		peer := Peer{
			ID:        remotePeerID,
			SessionID: msg.SessionId,
			Address:   net.ParseIP(payload.HelloRequest.Peer.Address),
			Username:  payload.HelloRequest.Peer.Username,
			Hostname:  payload.HelloRequest.Peer.Hostname,
			Role:      payload.HelloRequest.Peer.Role,
			HostURL:   payload.HelloRequest.Peer.HostUrl,
			FoundAt:   time.Now(),
			ExpireAt:  time.Now().Add(HelloTTL),
		}

		_ = s.membership.Discover(peer)

		s.log.DebugContext(ctx, "received hello request",
			"peer.id", peer.ID,
			"peer.address", peer.Address,
			"peer.username", lo.FromPtr(peer.Username),
			"peer.hostname", lo.FromPtr(peer.Hostname),
			"peer.role", peer.Role,
			"membership.size", s.membership.Size())

		_ = s.replyHello(ctx, peer)
	case *discovery.Message_HelloResult:
		if remotePeerID == s.me.ID {
			return nil
		}

		peer := Peer{
			ID:        remotePeerID,
			SessionID: msg.SessionId,
			Address:   net.ParseIP(payload.HelloResult.Peer.Address),
			Username:  payload.HelloResult.Peer.Username,
			Hostname:  payload.HelloResult.Peer.Hostname,
			Role:      payload.HelloResult.Peer.Role,
			HostURL:   payload.HelloResult.Peer.HostUrl,
			FoundAt:   time.Now(),
			ExpireAt:  time.Now().Add(HelloTTL),
		}

		_ = s.membership.Discover(peer)

		s.log.DebugContext(ctx, "received hello reply",
			"peer.id", peer.ID,
			"peer.address", peer.Address,
			"peer.username", lo.FromPtr(peer.Username),
			"peer.hostname", lo.FromPtr(peer.Hostname),
			"peer.role", peer.Role,
			"membership.size", s.membership.Size())
	}

	return nil
}
