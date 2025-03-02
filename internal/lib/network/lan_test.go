package network_test

import (
	"testing"

	"github.com/inhibitor1217/moru/internal/lib/network"
	"github.com/stretchr/testify/assert"
)

func TestLANAddresses(t *testing.T) {
	ips := network.LANAddresses()

	assert.NotNil(t, ips)
	t.Logf("LAN addresses: %v", ips)
}
