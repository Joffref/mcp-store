package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mcp-store-importer",
	Short: "Import MCPs from a directory",
	Long:  `mcp-store-importer is a CLI tool to import MCPs from a config file`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
