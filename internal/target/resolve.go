package target

import (
	"fmt"
	"os"
	"strings"
)

// Resolve returns a flat list of CIDR/IP strings from the given inputs.
// Priority: inline targets > file > TARGET_CIDRS env var.
// If all inputs are empty, an error is returned.
func Resolve(targets []string, file string) ([]string, error) {
	if len(targets) > 0 {
		return splitCIDRArgs(targets), nil
	}
	if file != "" {
		return readCIDRFile(file)
	}
	return fromEnv()
}

func splitCIDRArgs(args []string) []string {
	var out []string
	for _, a := range args {
		for _, cidr := range strings.Split(a, ",") {
			if cidr = strings.TrimSpace(cidr); cidr != "" {
				out = append(out, cidr)
			}
		}
	}
	return out
}

func readCIDRFile(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading targets file %q: %w", path, err)
	}
	var out []string
	for _, line := range strings.Split(string(data), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			out = append(out, line)
		}
	}
	return out, nil
}

func fromEnv() ([]string, error) {
	env := os.Getenv("TARGET_CIDRS")
	if env == "" {
		return nil, fmt.Errorf("no targets: provide --target, --target-file, or set TARGET_CIDRS")
	}
	var out []string
	for _, cidr := range strings.Split(env, ",") {
		if cidr = strings.TrimSpace(cidr); cidr != "" {
			out = append(out, cidr)
		}
	}
	return out, nil
}
