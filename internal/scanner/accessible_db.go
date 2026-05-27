package scanner

import (
	"fmt"
	"net"
	"time"

	"github.com/thomaslaurenson/prongs/internal/config"
)

// dbPorts are the ports checked per host. Each is tried independently.
var dbPorts = []int{3306, 5432} // MySQL, Postgres

// AccessibleDB checks whether any common database port is reachable.
type AccessibleDB struct{}

func (s *AccessibleDB) Name() string         { return "accessible-db" }
func (s *AccessibleDB) DefaultEnabled() bool { return true }

func (s *AccessibleDB) Run(ip net.IP) (Result, bool) {
	for _, port := range dbPorts {
		addr := net.JoinHostPort(ip.String(), fmt.Sprintf("%d", port))
		conn, err := net.DialTimeout("tcp", addr, time.Duration(config.DefaultTimeout)*time.Second)
		if err != nil {
			continue
		}
		conn.Close()
		return Result{
			Timestamp: time.Now().UTC(),
			IP:        ip,
			ScanType:  s.Name(),
			Port:      port,
		}, true
	}
	return Result{}, false
	// Note: this fixes the Python bug where the final result_queue.put used
	// whichever port value the loop happened to end on (always 5432).
}
