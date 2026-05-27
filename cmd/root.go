package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thomaslaurenson/prongs/internal/config"
)

func Execute() error {
	root := &cobra.Command{
		Use:          "prongs",
		Short:        "Fast, custom security scanner",
		Version:      config.Version,
		SilenceUsage: true,
	}

	root.AddCommand(newScanCmd())

	return root.Execute()
}
