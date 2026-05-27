package engine_test

import (
	"bytes"
	"io"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/thomaslaurenson/prongs/internal/engine"
	"github.com/thomaslaurenson/prongs/internal/scanner"
	"github.com/thomaslaurenson/prongs/internal/target"
)

// captureStdout redirects os.Stdout for the duration of fn, then returns
// everything that was printed.
func captureStdout(fn func()) string {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	old := os.Stdout
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()
	return buf.String()
}

// TestEngineSSHOnScanme is the end-to-end smoke test equivalent of Python's
// test_scanners_smoke. It exercises CIDR expansion + engine + SSH scanner
// against the real scanme.nmap.org host.
func TestEngineSSHOnScanme(t *testing.T) {
	hosts, err := target.Expand([]string{"45.33.32.156/32"})
	if err != nil {
		t.Fatalf("CIDR expand failed: %v", err)
	}

	out := captureStdout(func() {
		engine.Run(
			[]scanner.Scanner{&scanner.PasswordSSH{}},
			hosts,
			10,    // low concurrency - it's one host
			false, // machine-readable output
		)
	})

	lines := []string{}
	for _, l := range strings.Split(strings.TrimSpace(out), "\n") {
		if strings.TrimSpace(l) != "" {
			lines = append(lines, l)
		}
	}

	if len(lines) == 0 {
		t.Fatal("expected at least one output line, got none")
	}

	// Validate output format: timestamp\tip\tscanner\tport
	for _, line := range lines {
		parts := strings.Split(line, "\t")
		if len(parts) != 4 {
			t.Errorf("expected 4 tab-separated fields, got %d: %q", len(parts), line)
			continue
		}
		if _, err := time.Parse(time.RFC3339, parts[0]); err != nil {
			t.Errorf("field 0 is not RFC3339 timestamp: %q", parts[0])
		}
		if net.ParseIP(parts[1]) == nil {
			t.Errorf("field 1 is not an IP: %q", parts[1])
		}
		if parts[2] == "" {
			t.Errorf("field 2 (scanner name) is empty")
		}
		if parts[3] == "" {
			t.Errorf("field 3 (port) is empty")
		}
	}

	// Verify the expected finding is present
	found := false
	for _, line := range lines {
		if strings.Contains(line, "45.33.32.156") &&
			strings.Contains(line, "password-ssh") &&
			strings.Contains(line, "\t22") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected a password-ssh finding for 45.33.32.156:22, got output:\n%s", out)
	}
}

// TestEnginePrettyPrint exercises the pretty-print code path: the progress
// ticker goroutine, printResult with prettyPrint=true, and the summary line.
func TestEnginePrettyPrint(t *testing.T) {
	hosts, err := target.Expand([]string{"45.33.32.156/32"})
	if err != nil {
		t.Fatalf("CIDR expand failed: %v", err)
	}

	out := captureStdout(func() {
		engine.Run(
			[]scanner.Scanner{&scanner.PasswordSSH{}},
			hosts,
			10,
			true,
		)
	})

	if !strings.Contains(out, "Total hosts/checks:") {
		t.Errorf("expected summary line in pretty output, got:\n%s", out)
	}
}

// TestEngineEmptyHosts verifies a run with zero hosts produces no output and
// does not deadlock.
func TestEngineEmptyHosts(t *testing.T) {
	out := captureStdout(func() {
		engine.Run([]scanner.Scanner{&scanner.PasswordSSH{}}, nil, 4, false)
	})
	if strings.TrimSpace(out) != "" {
		t.Errorf("expected no output with no hosts, got:\n%s", out)
	}
}

// TestEngineEmptyScanners verifies a run with zero scanners produces no output
// and does not deadlock.
func TestEngineEmptyScanners(t *testing.T) {
	hosts, err := target.Expand([]string{"45.33.32.156/32"})
	if err != nil {
		t.Fatalf("CIDR expand failed: %v", err)
	}

	out := captureStdout(func() {
		engine.Run(nil, hosts, 4, false)
	})
	if strings.TrimSpace(out) != "" {
		t.Errorf("expected no output with no scanners, got:\n%s", out)
	}
}

// TestEngineNoFindings verifies that the engine produces no output when
// scanners find nothing - i.e. no phantom results, no extra lines.
func TestEngineNoFindings(t *testing.T) {
	hosts, err := target.Expand([]string{"45.33.32.156/32"})
	if err != nil {
		t.Fatalf("CIDR expand failed: %v", err)
	}

	// RDP and DB are both expected to be closed on scanme.nmap.org
	out := captureStdout(func() {
		engine.Run(
			[]scanner.Scanner{&scanner.AccessibleRDP{}, &scanner.AccessibleDB{}},
			hosts,
			10,
			false,
		)
	})

	if strings.TrimSpace(out) != "" {
		t.Errorf("expected no output for RDP/DB scan on scanme.nmap.org, got:\n%s", out)
	}
}
