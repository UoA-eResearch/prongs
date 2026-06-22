package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is set at build time using:
// -ldflags "-X github.com/thomaslaurenson/prongs/cmd.Version=...".
// It falls back to "dev" for local builds that do not inject a value.
var Version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the prongs version",
	Args:  cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("prongs version %s\n", Version)
	},
}
