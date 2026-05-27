package scanner

import (
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/thomaslaurenson/prongs/internal/config"
)

// PasswordSSH checks whether SSH password authentication is enabled on port 22.
// It does NOT attempt to brute-force - it only probes which auth methods the
// server advertises, mirroring the Python paramiko auth_none technique.
type PasswordSSH struct{}

func (s *PasswordSSH) Name() string         { return "password-ssh" }
func (s *PasswordSSH) DefaultEnabled() bool { return true }

func (s *PasswordSSH) Run(ip net.IP) (Result, bool) {
	addr := net.JoinHostPort(ip.String(), "22")
	timeout := time.Duration(config.DefaultTimeout) * time.Second

	// Quick TCP check first - avoids SSH handshake overhead on closed ports.
	// Mirrors the Python socket pre-check before invoking paramiko.
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return Result{}, false
	}
	conn.Close()

	// Probe using an empty password attempt. x/crypto/ssh internally sends a
	// "none" auth request first (RFC 4252 §5.2) to get the server's supported
	// method list, then only attempts our supplied methods that appear in that
	// list. The resulting error reports what was actually tried:
	//
	//   "ssh: unable to authenticate, attempted methods [none password], ..."
	//
	// "password" appearing in that list means the server offered it - i.e.
	// password auth is enabled. If the server does not support password auth,
	// the library never tries it and "password" is absent from the error.
	cfg := &ssh.ClientConfig{
		User: "cats_are_mythical", // same probe username as the Python version
		Auth: []ssh.AuthMethod{
			ssh.Password(""), // will be rejected; presence in error = server supports it
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // scanner context - host key unknown
		Timeout:         timeout,
	}

	_, err = ssh.Dial("tcp", addr, cfg)
	if err == nil {
		// Empty password accepted - that's a finding too.
		return Result{
			Timestamp: time.Now().UTC(),
			IP:        ip,
			ScanType:  s.Name(),
			Port:      22,
		}, true
	}

	// Check whether "password" appears in the list of attempted methods inside
	// the error string. We look inside brackets specifically to avoid matching
	// the word "password" elsewhere in unrelated error messages.
	// e.g. "attempted methods [none password], no supported methods remain"
	if methodInError(err.Error(), "password") {
		return Result{
			Timestamp: time.Now().UTC(),
			IP:        ip,
			ScanType:  s.Name(),
			Port:      22,
		}, true
	}

	return Result{}, false
}

// methodInError reports whether the given auth method name appears inside a
// bracket-delimited method list in an x/crypto/ssh error string, e.g.
// "attempted methods [none password], no supported methods remain".
func methodInError(errStr, method string) bool {
	inBracket := false
	start := 0
	for i, c := range errStr {
		switch c {
		case '[':
			inBracket = true
			start = i + 1
		case ']':
			if inBracket {
				for _, m := range strings.Fields(errStr[start:i]) {
					if strings.TrimRight(m, ",") == method {
						return true
					}
				}
			}
			inBracket = false
		}
	}
	return false
}
