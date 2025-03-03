package beacon

import (
	"context"
	"errors"
	"net"
)

var (
	// ErrBeaconStopped is returned when the beacon is stopped.
	ErrBeaconStopped = errors.New("beacon stopped")
)

// Beacon is a service that broadcasts and receives UDP messages on a local network.
type Beacon interface {
	// Start begins broadcasting and receiving messages.
	// It returns an error if the service fails to start.
	Start(ctx context.Context) error

	// Stop gracefully shuts down the beacon service.
	Stop(ctx context.Context) error

	// Send broadcasts a message to the local network.
	Send(ctx context.Context, msg []byte, opts ...SendOption) error

	// Receive blocks until a message is received or the context is cancelled.
	Receive(ctx context.Context) ([]byte, net.Addr, error)
}

type sendOpts struct {
	broadcast bool
	unicastIP net.IP
}

type SendOption func(*sendOpts)

var SendBroadcast = func(opts *sendOpts) {
	opts.broadcast = true
}

var SendUnicast = func(ip net.IP) SendOption {
	return func(opts *sendOpts) {
		opts.broadcast = false
		opts.unicastIP = ip
	}
}
