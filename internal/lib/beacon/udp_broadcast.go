package beacon

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"time"
)

// UDPBroadcastConfig holds the configuration for UDP broadcast beacon.
type UDPBroadcastConfig struct {
	// Port is the UDP port number to use for broadcasting and receiving.
	Port int

	// Addr is the broadcast address to send messages to.
	// If empty, "255.255.255.255" will be used.
	Addr string

	// WriteTimeout configures the maximum duration for writing a message.
	// Default value is 1 second.
	WriteTimeout time.Duration
}

// udpBroadcast implements the Beacon interface using UDP broadcast.
type udpBroadcast struct {
	cfg UDPBroadcastConfig

	started bool
	stopped bool
	mu      sync.Mutex

	inbox  chan []byte
	outbox chan udpRecv
	stop   chan struct{}

	sendWg sync.WaitGroup
	recvWg sync.WaitGroup

	log *slog.Logger
}

type udpRecv struct {
	data []byte
	addr net.Addr
}

// NewUDPBroadcast creates a new UDP broadcast beacon with the given configuration.
func NewUDPBroadcast(cfg UDPBroadcastConfig) (Beacon, error) {
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return nil, fmt.Errorf("invalid port number: %d", cfg.Port)
	}

	if cfg.Addr == "" {
		cfg.Addr = "255.255.255.255"
	}

	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = 1 * time.Second
	}

	return &udpBroadcast{
		cfg: cfg,

		inbox:  make(chan []byte),
		outbox: make(chan udpRecv, 64),
		stop:   make(chan struct{}),

		log: slog.Default().With("source", "beacon.udpBroadcast"),
	}, nil
}

func (b *udpBroadcast) Start(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.started {
		return fmt.Errorf("beacon already started")
	} else if b.stopped {
		return fmt.Errorf("beacon already stopped")
	}
	b.started = true

	b.log.InfoContext(ctx, "starting UDP broadcast beacon")

	bgCtx := context.WithoutCancel(ctx)

	b.sendWg.Add(1)
	go b.senderLoop(bgCtx)

	b.recvWg.Add(1)
	go b.receiverLoop(bgCtx)

	return nil
}

func (b *udpBroadcast) Stop(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.started {
		return fmt.Errorf("beacon not started")
	} else if b.stopped {
		return fmt.Errorf("beacon already stopped")
	}
	b.stopped = true

	b.log.InfoContext(ctx, "stopping UDP broadcast beacon")

	close(b.stop)

	b.sendWg.Wait()
	b.recvWg.Wait()

	return nil
}

func (b *udpBroadcast) Send(ctx context.Context, msg []byte) error {
	select {
	case b.inbox <- msg:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-b.stop:
		return fmt.Errorf("beacon already stopped")
	}
}

func (b *udpBroadcast) Receive(ctx context.Context) ([]byte, net.Addr, error) {
	select {
	case recv := <-b.outbox:
		return recv.data, recv.addr, nil
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	case <-b.stop:
		return nil, nil, fmt.Errorf("beacon already stopped")
	}
}

func (b *udpBroadcast) senderLoop(ctx context.Context) {
	defer b.sendWg.Done()

	b.log.InfoContext(ctx, "starting UDP broadcast sender loop")
	defer b.log.InfoContext(ctx, "stopping UDP broadcast sender loop")

	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		b.log.ErrorContext(ctx, "failed to listen UDP", "error", err)
		return
	}
	defer conn.Close()

	for {
		select {
		case msg := <-b.inbox:
			dst := &net.UDPAddr{IP: net.ParseIP(b.cfg.Addr), Port: b.cfg.Port}
			_ = conn.SetWriteDeadline(time.Now().Add(b.cfg.WriteTimeout))
			_, err = conn.WriteTo(msg, dst)
			_ = conn.SetWriteDeadline(time.Time{})

			if err != nil {
				b.log.ErrorContext(ctx, "failed to broadcast message", "error", err)
				continue
			}

			b.log.DebugContext(ctx, "sent UDP broadcast message", "addr", dst, "size", len(msg))
		case <-b.stop:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (b *udpBroadcast) receiverLoop(ctx context.Context) {
	defer b.recvWg.Done()

	b.log.InfoContext(ctx, "starting UDP broadcast receiver loop")
	defer b.log.InfoContext(ctx, "stopping UDP broadcast receiver loop")

	conn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: b.cfg.Port})
	if err != nil {
		b.log.ErrorContext(ctx, "failed to listen UDP", "error", err)
		return
	}

	// separate goroutine for stop signal
	go func() {
		<-b.stop
		_ = conn.Close() // this will interrupt conn.ReadFrom
	}()

	bs := make([]byte, 65536)
	for {
		n, addr, err := conn.ReadFrom(bs)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}

			b.log.ErrorContext(ctx, "failed to receive message", "error", err)
			time.Sleep(1 * time.Second)
			continue
		}

		b.log.DebugContext(ctx, "received UDP broadcast message", "addr", addr, "size", n)

		buf := make([]byte, n)
		copy(buf, bs[:n])

		select {
		case b.outbox <- udpRecv{data: buf, addr: addr}:
		case <-b.stop:
			return
		case <-ctx.Done():
			return
		}
	}
}
