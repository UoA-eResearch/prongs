package target

import (
	"fmt"
	"net"
)

// Expand parses one or more CIDR strings and returns every host IP.
// Single IPs (e.g. "45.33.32.156") are accepted without a prefix length.
// Comment lines (starting with #) are skipped.
func Expand(cidrs []string) ([]net.IP, error) {
	var hosts []net.IP
	for _, cidr := range cidrs {
		// Skip empty lines and comments
		if cidr == "" || len(cidr) > 0 && cidr[0] == '#' {
			continue
		}

		// Accept bare IPs by treating them as /32
		if ip := net.ParseIP(cidr); ip != nil {
			hosts = append(hosts, ip)
			continue
		}

		// Parse as CIDR
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, fmt.Errorf("invalid CIDR or IP: %s", cidr)
		}

		// Determine network size
		ones, bits := network.Mask.Size()

		// Special case: /32 (single host) or /31 (point-to-point per RFC 3021)
		if ones >= bits-1 {
			// For /32: just the single IP
			// For /31: both IPs are usable
			for ip := cloneIP(network.IP); network.Contains(ip); inc(ip) {
				hosts = append(hosts, cloneIP(ip))
			}
			continue
		}

		// For /30 and larger: skip network and broadcast addresses
		for ip := cloneIP(network.IP); network.Contains(ip); inc(ip) {
			if !isNetworkOrBroadcast(ip, network) {
				hosts = append(hosts, cloneIP(ip))
			}
		}
	}
	return hosts, nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func cloneIP(ip net.IP) net.IP {
	clone := make(net.IP, len(ip))
	copy(clone, ip)
	return clone
}

func isNetworkOrBroadcast(ip net.IP, network *net.IPNet) bool {
	// Network address: all host bits zero
	if ip.Equal(network.IP) {
		return true
	}
	// Broadcast: all host bits one
	broadcast := make(net.IP, len(network.IP))
	for i := range network.IP {
		broadcast[i] = network.IP[i] | ^network.Mask[i]
	}
	return ip.Equal(broadcast)
}
