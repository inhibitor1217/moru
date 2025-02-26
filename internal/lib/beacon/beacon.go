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
	Send(ctx context.Context, msg []byte) error

	// Receive blocks until a message is received or the context is cancelled.
	Receive(ctx context.Context) ([]byte, net.Addr, error)
}
