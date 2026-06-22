package cmd

import (
	"github.com/spf13/cobra"
)

func Execute() error {
	root := &cobra.Command{
		Use:          "prongs",
		Short:        "Fast, custom security scanner",
		Version:      Version,
		SilenceUsage: true,
	}

	root.AddCommand(newScanCmd())
	root.AddCommand(versionCmd)

	return root.Execute()
}
