package cmd

import (
	"fmt"
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
		targetFile  string
		scannerArgs []string
		all         bool
		output      string
		concurrency int
	)

	var scannerNames []string
	for _, s := range scanner.All {
		scannerNames = append(scannerNames, s.Name())
	}

	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Run scanners against target CIDRs",
		Long: `Run one or more scanners against one or more target networks.

Targets are CIDR ranges or single IPs, provided via --target (repeatable
and/or comma-separated) or --target-file (a file, one entry per line). The
two flags are mutually exclusive. If neither is provided, the TARGET_CIDRS
environment variable (comma-separated) is used as a fallback.

Examples:
  # Run one scanner against a single network
  prongs scan --scanner password-ssh --target 192.168.0.0/24

  # Run all default scanners against multiple networks
  prongs scan --all --target 192.168.0.0/24 --target 10.0.0.0/24

  # Multiple networks as a single comma-separated value
  prongs scan --all --target 192.168.0.0/24,10.0.0.0/24

  # Load targets from a file
  prongs scan --all --target-file targets.txt

  # Pretty-print output
  prongs scan --all --target 192.168.0.0/24 --output pretty

  # Limit the number of concurrent probes
  prongs scan --all --target 192.168.0.0/24 --concurrency 50

  # Use the TARGET_CIDRS environment variable (comma-separated)
  TARGET_CIDRS=192.168.0.0/24,10.0.0.0/24 prongs scan --all`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if output != "text" && output != "pretty" {
				return fmt.Errorf("--output must be 'text' or 'pretty', got %q", output)
			}

			if !all && len(scannerArgs) == 0 {
				return fmt.Errorf("provide --scanner <name> or --all\nAvailable: %s",
					strings.Join(scannerNames, ", "))
			}

			if len(targetArgs) > 0 && targetFile != "" {
				return fmt.Errorf("--target and --target-file are mutually exclusive")
			}

			var active []scanner.Scanner
			if all {
				active = scanner.Defaults()
			} else {
				for _, name := range scannerArgs {
					s, ok := scanner.ByName[name]
					if !ok {
						return fmt.Errorf("unknown scanner %q - available: %s",
							name, strings.Join(scannerNames, ", "))
					}
					active = append(active, s)
				}
			}

			cidrList, err := target.Resolve(targetArgs, targetFile)
			if err != nil {
				return err
			}
			if len(cidrList) == 0 {
				return fmt.Errorf("no valid targets found")
			}

			hosts, err := target.Expand(cidrList)
			if err != nil {
				return err
			}
			if len(hosts) == 0 {
				return fmt.Errorf("no valid IP addresses found after CIDR expansion")
			}

			engine.Run(active, hosts, concurrency, output == "pretty")
			return nil
		},
	}

	cmd.Flags().StringArrayVar(&targetArgs, "target", nil,
		"CIDR(s) to scan (repeatable, comma-separated)")
	cmd.Flags().StringVar(&targetFile, "target-file", "",
		"Path to a file of CIDRs or IPs, one per line")
	cmd.Flags().StringArrayVar(&scannerArgs, "scanner", nil,
		fmt.Sprintf("Scanner to run (repeatable) - choices: %s", strings.Join(scannerNames, ", ")))
	cmd.Flags().BoolVar(&all, "all", false, "Run all default-enabled scanners")
	cmd.Flags().StringVar(&output, "output", "text",
		"Output format: 'text' (TSV, default) or 'pretty' (human-readable)")
	cmd.Flags().IntVarP(&concurrency, "concurrency", "c", config.DefaultConcurrency,
		"Max concurrent probes")

	return cmd
}
