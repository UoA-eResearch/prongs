package scanner_test

import (
	"testing"

	"github.com/thomaslaurenson/prongs/internal/scanner"
)

func TestByNameCoversAll(t *testing.T) {
	t.Parallel()

	for _, s := range scanner.All {
		t.Run(s.Name(), func(t *testing.T) {
			t.Parallel()
			got, ok := scanner.ByName[s.Name()]
			if !ok {
				t.Errorf("ByName missing entry for %q", s.Name())
			}
			if got.Name() != s.Name() {
				t.Errorf("ByName[%q].Name() = %q, want %q", s.Name(), got.Name(), s.Name())
			}
		})
	}
}

func TestDefaultsAreSubsetOfAll(t *testing.T) {
	t.Parallel()

	defaults := scanner.Defaults()
	for _, d := range defaults {
		if !d.DefaultEnabled() {
			t.Errorf("Defaults() returned %q which has DefaultEnabled() == false", d.Name())
		}
		found := false
		for _, s := range scanner.All {
			if s.Name() == d.Name() {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Defaults() returned %q which is not in All", d.Name())
		}
	}
}

func TestDefaultsExcludesNonDefault(t *testing.T) {
	t.Parallel()

	defaults := scanner.Defaults()
	defaultNames := make(map[string]bool, len(defaults))
	for _, d := range defaults {
		defaultNames[d.Name()] = true
	}

	for _, s := range scanner.All {
		if !s.DefaultEnabled() && defaultNames[s.Name()] {
			t.Errorf("Defaults() included %q but DefaultEnabled() is false", s.Name())
		}
	}
}
