package target_test

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/thomaslaurenson/prongs/internal/target"
)

func TestResolveFromArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		targets []string
		want    []string
	}{
		{
			name:    "single CIDR",
			targets: []string{"192.168.0.0/24"},
			want:    []string{"192.168.0.0/24"},
		},
		{
			name:    "comma-separated CIDRs in one arg",
			targets: []string{"192.168.0.0/24,10.0.0.0/8"},
			want:    []string{"192.168.0.0/24", "10.0.0.0/8"},
		},
		{
			name:    "multiple args",
			targets: []string{"192.168.0.0/24", "10.0.0.0/8"},
			want:    []string{"192.168.0.0/24", "10.0.0.0/8"},
		},
		{
			name:    "whitespace trimmed",
			targets: []string{"192.168.0.0/24, 10.0.0.0/8"},
			want:    []string{"192.168.0.0/24", "10.0.0.0/8"},
		},
		{
			name:    "trailing comma ignored",
			targets: []string{"192.168.0.0/24,"},
			want:    []string{"192.168.0.0/24"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := target.Resolve(tc.targets, "")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !slices.Equal(got, tc.want) {
				t.Errorf("Resolve(%v, \"\") = %v, want %v", tc.targets, got, tc.want)
			}
		})
	}
}

func TestResolveFromFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		contents string
		want     []string
	}{
		{
			name:     "simple file",
			contents: "192.168.0.0/24\n10.0.0.0/8\n",
			want:     []string{"192.168.0.0/24", "10.0.0.0/8"},
		},
		{
			name:     "blank lines skipped",
			contents: "192.168.0.0/24\n\n10.0.0.0/8\n",
			want:     []string{"192.168.0.0/24", "10.0.0.0/8"},
		},
		{
			name:     "whitespace-only lines skipped",
			contents: "192.168.0.0/24\n   \n10.0.0.0/8",
			want:     []string{"192.168.0.0/24", "10.0.0.0/8"},
		},
		{
			name:     "file with single entry and no trailing newline",
			contents: "10.0.0.1",
			want:     []string{"10.0.0.1"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			path := filepath.Join(dir, "targets.txt")
			if err := os.WriteFile(path, []byte(tc.contents), 0o600); err != nil {
				t.Fatalf("WriteFile: %v", err)
			}
			got, err := target.Resolve(nil, path)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !slices.Equal(got, tc.want) {
				t.Errorf("Resolve(nil, file) = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestResolveFromFileMissing(t *testing.T) {
	t.Parallel()
	_, err := target.Resolve(nil, "/nonexistent/targets.txt")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestResolveFromEnv(t *testing.T) {
	t.Setenv("TARGET_CIDRS", "192.168.0.0/24,10.0.0.0/8")
	got, err := target.Resolve(nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"192.168.0.0/24", "10.0.0.0/8"}
	if !slices.Equal(got, want) {
		t.Errorf("Resolve from env = %v, want %v", got, want)
	}
}

func TestResolveNoSource(t *testing.T) {
	t.Setenv("TARGET_CIDRS", "")
	_, err := target.Resolve(nil, "")
	if err == nil {
		t.Fatal("expected error when no targets are provided, got nil")
	}
}
