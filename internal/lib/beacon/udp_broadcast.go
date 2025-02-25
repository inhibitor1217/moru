package beacon

import (
	"context"
	"fmt"
	"net"
)

// UDPBroadcastConfig holds the configuration for UDP broadcast beacon.
type UDPBroadcastConfig struct {
	// Port is the UDP port number to use for broadcasting and receiving.
	Port int

	// Addr is the broadcast address to send messages to.
	// If empty, "255.255.255.255" will be used.
	Addr string
}

// udpBroadcast implements the Beacon interface using UDP broadcast.
type udpBroadcast struct {
	cfg UDPBroadcastConfig
}

// NewUDPBroadcast creates a new UDP broadcast beacon with the given configuration.
func NewUDPBroadcast(cfg UDPBroadcastConfig) (Beacon, error) {
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return nil, fmt.Errorf("invalid port number: %d", cfg.Port)
	}

	if cfg.Addr == "" {
		cfg.Addr = "255.255.255.255"
	}

	return &udpBroadcast{
		cfg: cfg,
	}, nil
}

func (b *udpBroadcast) Start(ctx context.Context) error {
	// TODO implement
	return nil
}

func (b *udpBroadcast) Stop(ctx context.Context) error {
	// TODO implement
	return nil
}

func (b *udpBroadcast) Send(ctx context.Context, msg []byte) error {
	// TODO implement
	return nil
}

func (b *udpBroadcast) Receive(ctx context.Context) ([]byte, net.Addr, error) {
	// TODO implement
	return nil, nil, nil
}
