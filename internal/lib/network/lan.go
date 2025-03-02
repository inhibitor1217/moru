package network

import "net"

// LANAddresses returns the local network addresses of the host.
func LANAddresses() []net.IP {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	ips := make([]net.IP, 0)
	for _, iface := range netInterfaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.IsLoopback() {
					continue
				}
				if v.IP.To4() != nil {
					ips = append(ips, v.IP)
				}
			}
		}
	}

	return ips
}
