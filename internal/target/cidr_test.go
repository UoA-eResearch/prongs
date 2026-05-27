package target_test

import (
	"testing"

	"github.com/thomaslaurenson/prongs/internal/target"
)

func TestExpandSingleIP(t *testing.T) {
	hosts, err := target.Expand([]string{"45.33.32.156"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(hosts) != 1 {
		t.Fatalf("expected 1 host, got %d", len(hosts))
	}
	if hosts[0].String() != "45.33.32.156" {
		t.Errorf("expected 45.33.32.156, got %s", hosts[0])
	}
}

func TestExpandCIDR32(t *testing.T) {
	// /32 should return exactly 1 host (not 0!)
	hosts, err := target.Expand([]string{"45.33.32.156/32"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(hosts) != 1 {
		t.Fatalf("expected 1 host for /32, got %d", len(hosts))
	}
	if hosts[0].String() != "45.33.32.156" {
		t.Errorf("expected 45.33.32.156, got %s", hosts[0])
	}
}

func TestExpandCIDR31(t *testing.T) {
	// /31 should return 2 hosts (RFC 3021 point-to-point)
	hosts, err := target.Expand([]string{"192.168.1.0/31"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(hosts) != 2 {
		t.Fatalf("expected 2 hosts for /31, got %d", len(hosts))
	}
}

func TestExpandCIDR30(t *testing.T) {
	// /30 has 2 usable hosts (4 total - network - broadcast)
	hosts, err := target.Expand([]string{"192.168.1.0/30"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(hosts) != 2 {
		t.Fatalf("expected 2 usable hosts for /30, got %d", len(hosts))
	}
	// Should be .1 and .2 (not .0 or .3)
	if hosts[0].String() != "192.168.1.1" {
		t.Errorf("expected first host to be 192.168.1.1, got %s", hosts[0])
	}
	if hosts[1].String() != "192.168.1.2" {
		t.Errorf("expected second host to be 192.168.1.2, got %s", hosts[1])
	}
}

func TestExpandMultipleCIDRs(t *testing.T) {
	hosts, err := target.Expand([]string{"192.168.1.0/31", "10.0.0.1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(hosts) != 3 {
		t.Fatalf("expected 3 hosts (2 from /31 + 1 bare IP), got %d", len(hosts))
	}
}

func TestExpandInvalid(t *testing.T) {
	_, err := target.Expand([]string{"not-a-cidr"})
	if err == nil {
		t.Fatal("expected error for invalid input")
	}
}

func TestExpandEmptyString(t *testing.T) {
	hosts, err := target.Expand([]string{""})
	if err != nil {
		t.Fatalf("unexpected error for empty string: %v", err)
	}
	if len(hosts) != 0 {
		t.Fatalf("expected 0 hosts for empty string, got %d", len(hosts))
	}
}

func TestExpandCommentLines(t *testing.T) {
	hosts, err := target.Expand([]string{"# This is a comment", "192.168.1.1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(hosts) != 1 {
		t.Fatalf("expected 1 host (comment should be skipped), got %d", len(hosts))
	}
	if hosts[0].String() != "192.168.1.1" {
		t.Errorf("expected 192.168.1.1, got %s", hosts[0])
	}
}
