package scanner_test

import (
	"net"
	"testing"

	"github.com/thomaslaurenson/prongs/internal/scanner"
)

var testIP = net.ParseIP("45.33.32.156") // scanme.nmap.org

func TestScannerMetadata(t *testing.T) {
	tests := []struct {
		s             scanner.Scanner
		wantName      string
		wantDefaultOn bool
	}{
		{&scanner.PasswordSSH{}, "password-ssh", true},
		{&scanner.AccessibleRDP{}, "accessible-rdp", false},
		{&scanner.AccessibleDB{}, "accessible-db", true},
	}
	for _, tt := range tests {
		t.Run(tt.wantName, func(t *testing.T) {
			if got := tt.s.Name(); got != tt.wantName {
				t.Errorf("Name() = %q, want %q", got, tt.wantName)
			}
			if got := tt.s.DefaultEnabled(); got != tt.wantDefaultOn {
				t.Errorf("DefaultEnabled() = %v, want %v", got, tt.wantDefaultOn)
			}
		})
	}
}

// TestRDPNotExposed checks that scanme.nmap.org does NOT have RDP open.
// This is a negative test - validates the scanner doesn't false-positive.
func TestRDPNotExposed(t *testing.T) {
	s := &scanner.AccessibleRDP{}
	_, found := s.Run(testIP)
	if found {
		t.Errorf("RDP reported open on scanme.nmap.org - unexpected")
	}
}

// TestDBNotExposed checks that scanme.nmap.org does NOT have DB ports open.
func TestDBNotExposed(t *testing.T) {
	s := &scanner.AccessibleDB{}
	_, found := s.Run(testIP)
	if found {
		t.Errorf("DB port reported open on scanme.nmap.org - unexpected")
	}
}

// TestSSHPasswordAuthEnabled checks scanme.nmap.org, which is known to have
// password auth enabled. This is the primary regression test - it catches
// any breakage in the SSH probing logic.
func TestSSHPasswordAuthEnabled(t *testing.T) {
	s := &scanner.PasswordSSH{}
	r, found := s.Run(testIP)
	if !found {
		t.Error("expected password-ssh to be reported open on scanme.nmap.org")
	}
	if r.Port != 22 {
		t.Errorf("expected port 22, got %d", r.Port)
	}
	if r.ScanType != "password-ssh" {
		t.Errorf("expected scan type 'password-ssh', got %s", r.ScanType)
	}
}
