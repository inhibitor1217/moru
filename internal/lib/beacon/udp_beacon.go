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

// UDPBeaconConfig holds the configuration for UDP broadcast beacon.
type UDPBeaconConfig struct {
	// Port is the UDP port number to use for broadcasting and receiving.
	Port int

	// WriteTimeout configures the maximum duration for writing a message.
	// Default value is 1 second.
	WriteTimeout time.Duration
}

// udpBeacon implements the Beacon interface using UDP broadcast.
type udpBeacon struct {
	cfg UDPBeaconConfig

	started bool
	stopped bool
	mu      sync.Mutex

	inbox  chan udpSend
	outbox chan udpRecv
	stop   chan struct{}

	sendWg sync.WaitGroup
	recvWg sync.WaitGroup

	log *slog.Logger
}

type udpSend struct {
	data []byte
	addr net.Addr
}

type udpRecv struct {
	data []byte
	addr net.Addr
}

// NewUDPBeacon creates a new UDP beacon with the given configuration.
func NewUDPBeacon(cfg UDPBeaconConfig) (Beacon, error) {
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return nil, fmt.Errorf("invalid port number: %d", cfg.Port)
	}

	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = 1 * time.Second
	}

	return &udpBeacon{
		cfg: cfg,

		inbox:  make(chan udpSend),
		outbox: make(chan udpRecv, 64),
		stop:   make(chan struct{}),

		log: slog.Default().With("source", "beacon.udpBeacon"),
	}, nil
}

func (b *udpBeacon) Start(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.started {
		return fmt.Errorf("beacon already started")
	} else if b.stopped {
		return fmt.Errorf("beacon already stopped")
	}
	b.started = true

	b.log.InfoContext(ctx, "starting UDP beacon")

	bgCtx := context.WithoutCancel(ctx)

	b.sendWg.Add(1)
	go b.senderLoop(bgCtx)

	b.recvWg.Add(1)
	go b.receiverLoop(bgCtx)

	return nil
}

func (b *udpBeacon) Stop(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.started {
		return fmt.Errorf("beacon not started")
	} else if b.stopped {
		return fmt.Errorf("beacon already stopped")
	}
	b.stopped = true

	b.log.InfoContext(ctx, "stopping UDP beacon")

	close(b.stop)

	b.sendWg.Wait()
	b.recvWg.Wait()

	return nil
}

func (b *udpBeacon) Send(ctx context.Context, msg []byte, opts ...SendOption) error {
	option := sendOpts{broadcast: true}
	for _, opt := range opts {
		opt(&option)
	}

	var addr net.Addr
	if option.broadcast {
		addr = &net.UDPAddr{IP: net.ParseIP("255.255.255.255"), Port: b.cfg.Port}
	} else {
		addr = &net.UDPAddr{IP: option.unicastIP, Port: b.cfg.Port}
	}

	select {
	case b.inbox <- udpSend{data: msg, addr: addr}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-b.stop:
		return ErrBeaconStopped
	}
}

func (b *udpBeacon) Receive(ctx context.Context) ([]byte, net.Addr, error) {
	select {
	case recv := <-b.outbox:
		return recv.data, recv.addr, nil
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	case <-b.stop:
		return nil, nil, ErrBeaconStopped
	}
}

func (b *udpBeacon) senderLoop(ctx context.Context) {
	defer b.sendWg.Done()

	b.log.InfoContext(ctx, "starting UDP sender loop")
	defer b.log.InfoContext(ctx, "stopping UDP sender loop")

	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		b.log.ErrorContext(ctx, "failed to listen UDP", "error", err)
		return
	}
	defer conn.Close()

	if err = conn.SetWriteBuffer(64 * 1024); err != nil {
		b.log.ErrorContext(ctx, "failed to set write buffer", "error", err)
		return
	}

	for {
		select {
		case in := <-b.inbox:
			_ = conn.SetWriteDeadline(time.Now().Add(b.cfg.WriteTimeout))
			_, err = conn.WriteTo(in.data, in.addr)
			_ = conn.SetWriteDeadline(time.Time{})

			if err != nil {
				b.log.ErrorContext(ctx, "failed to message", "error", err)
				continue
			}

			b.log.DebugContext(ctx, "sent UDP message", "addr", in.addr.String(), "size", len(in.data))
		case <-b.stop:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (b *udpBeacon) receiverLoop(ctx context.Context) {
	defer b.recvWg.Done()

	b.log.InfoContext(ctx, "starting UDP receiver loop")
	defer b.log.InfoContext(ctx, "stopping UDP receiver loop")

	conn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: b.cfg.Port})
	if err != nil {
		b.log.ErrorContext(ctx, "failed to listen UDP", "error", err)
		return
	}

	if err = conn.SetReadBuffer(64 * 1024); err != nil {
		b.log.ErrorContext(ctx, "failed to set read buffer", "error", err)
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

		b.log.DebugContext(ctx, "received UDP message", "addr", addr, "size", n)

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
