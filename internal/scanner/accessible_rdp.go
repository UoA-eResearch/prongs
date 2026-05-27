package scanner

import (
	"net"
	"time"

	"github.com/thomaslaurenson/prongs/internal/config"
)

// AccessibleRDP checks if TCP 3389 (RDP) is reachable.
type AccessibleRDP struct{}

func (s *AccessibleRDP) Name() string { return "accessible-rdp" }

// DefaultEnabled returns false - mirrors the Python behaviour where accessible-rdp
// was silently excluded from --enable-all. Now it's explicit.
func (s *AccessibleRDP) DefaultEnabled() bool { return false }

func (s *AccessibleRDP) Run(ip net.IP) (Result, bool) {
	addr := net.JoinHostPort(ip.String(), "3389")
	conn, err := net.DialTimeout("tcp", addr, time.Duration(config.DefaultTimeout)*time.Second)
	if err != nil {
		return Result{}, false
	}
	conn.Close()
	return Result{
		Timestamp: time.Now().UTC(),
		IP:        ip,
		ScanType:  s.Name(),
		Port:      3389,
	}, true
}
