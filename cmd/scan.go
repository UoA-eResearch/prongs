package cmd

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thomaslaurenson/prongs/internal/config"
	"github.com/thomaslaurenson/prongs/internal/engine"
	"github.com/thomaslaurenson/prongs/internal/scanner"
	"github.com/thomaslaurenson/prongs/internal/target"
)

func newScanCmd() *cobra.Command {
	var (
		targetArgs  []string
		scannerArgs []string
		all         bool
		output      string
		concurrency int
	)

	// Build scanner name list for help text
	var scannerNames []string
	for _, s := range scanner.All {
		scannerNames = append(scannerNames, s.Name())
	}

	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Run scanners against target CIDRs",
		Long: `Run one or more scanners against one or more target networks.

Targets can be provided as CIDR ranges or single IPs, either directly
via --target or in a file (one per line). If --target is omitted,
the TARGET_CIDRS environment variable is used as a fallback.

Examples:
  prongs scan --scanner password-ssh --target 192.168.0.0/24
  prongs scan --scanner password-ssh --target targets.txt
  prongs scan --all --output pretty
  TARGET_CIDRS=10.0.0.0/8 prongs scan --all`,
		RunE: func(cmd *cobra.Command, args []string) error {

			// Validate output flag
			if output != "text" && output != "pretty" {
				return fmt.Errorf("--output must be 'text' or 'pretty', got %q", output)
			}

			// Validate: at least one scanner selected
			if !all && len(scannerArgs) == 0 {
				return fmt.Errorf("provide --scanner <name> or --all\nAvailable: %s",
					strings.Join(scannerNames, ", "))
			}

			// Resolve active scanners
			var active []scanner.Scanner
			if all {
				for _, s := range scanner.All {
					if s.DefaultEnabled() {
						active = append(active, s)
					}
				}
			} else {
				for _, name := range scannerArgs {
					found := false
					for _, s := range scanner.All {
						if s.Name() == name {
							active = append(active, s)
							found = true
							break
						}
					}
					if !found {
						return fmt.Errorf("unknown scanner %q - available: %s",
							name, strings.Join(scannerNames, ", "))
					}
				}
			}

			// Resolve target CIDRs
			// Priority: --target flag > TARGET_CIDRS env var
			var cidrList []string
			if len(targetArgs) > 0 {
				for _, t := range targetArgs {
					// Auto-detect file vs CIDR: check if it exists as a file first
					if looksLikeFile(t) {
						data, err := os.ReadFile(t)
						if err != nil {
							return fmt.Errorf("reading targets file %q: %w", t, err)
						}
						for _, line := range strings.Split(string(data), "\n") {
							if line = strings.TrimSpace(line); line != "" {
								cidrList = append(cidrList, line)
							}
						}
					} else {
						// Comma-separated inline CIDRs also accepted
						for _, cidr := range strings.Split(t, ",") {
							if cidr = strings.TrimSpace(cidr); cidr != "" {
								cidrList = append(cidrList, cidr)
							}
						}
					}
				}
			} else {
				env := os.Getenv("TARGET_CIDRS")
				if env == "" {
					return fmt.Errorf("no targets: provide --target or set TARGET_CIDRS env var")
				}
				for _, cidr := range strings.Split(env, ",") {
					if cidr = strings.TrimSpace(cidr); cidr != "" {
						cidrList = append(cidrList, cidr)
					}
				}
			}

			// Validate we have at least one target before expanding
			if len(cidrList) == 0 {
				return fmt.Errorf("no valid targets found (empty file or environment variable)")
			}

			hosts, err := target.Expand(cidrList)
			if err != nil {
				return err
			}

			// Validate we got at least one host after expansion
			if len(hosts) == 0 {
				return fmt.Errorf("no valid IP addresses found after CIDR expansion (check for comments or invalid CIDRs)")
			}

			// Run the scanners concurrently
			engine.Run(active, hosts, concurrency, output == "pretty")
			return nil
		},
	}

	cmd.Flags().StringArrayVar(&targetArgs, "target", nil,
		"Target CIDR(s) or path to a file of CIDRs (repeatable, comma-separated)")
	cmd.Flags().StringArrayVar(&scannerArgs, "scanner", nil,
		fmt.Sprintf("Scanner to run (repeatable) - choices: %s", strings.Join(scannerNames, ", ")))
	cmd.Flags().BoolVar(&all, "all", false, "Run all default-enabled scanners")
	cmd.Flags().StringVar(&output, "output", "text",
		"Output format: 'text' (TSV, default) or 'pretty' (human-readable)")
	cmd.Flags().IntVarP(&concurrency, "concurrency", "c", config.DefaultConcurrency,
		"Max concurrent probes")

	return cmd
}

// looksLikeFile returns true when s is more likely a file path than a CIDR/IP.
// An existing file always wins; for non-existent paths containing a separator,
// we use net.ParseIP and net.ParseCIDR to decide.
func looksLikeFile(s string) bool {
	// Existing file always wins
	info, err := os.Stat(s)
	if err == nil && !info.IsDir() {
		return true
	}

	// If it contains a path separator but isn't a valid IP or CIDR, treat as file
	if strings.ContainsAny(s, "/\\") {
		if net.ParseIP(s) != nil {
			return false
		}
		_, _, err := net.ParseCIDR(s)
		return err != nil
	}

	return false
}
