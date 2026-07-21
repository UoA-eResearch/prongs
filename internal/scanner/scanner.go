package scanner

import (
	"net"
	"time"
)

// Result is a single confirmed finding.
type Result struct {
	Timestamp time.Time
	IP        net.IP
	ScanType  string
	Port      int
}

// Scanner is implemented by every scan module.
type Scanner interface {
	// Name returns the scanner identifier used in --scanner flags and output.
	Name() string

	// DefaultEnabled returns false for scanners excluded from --all.
	// Mirrors the Python hack in run.py that excluded accessible-rdp from --enable-all.
	DefaultEnabled() bool

	// Run probes a single IP. Returns (result, true) on a finding, (zero, false) otherwise.
	Run(ip net.IP) (Result, bool)
}
