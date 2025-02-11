package cmd

import (
	"fmt"
	"log"

	"github.com/beamlit/mcp-store/internal/store"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the MCP store",
	RunE: func(cmd *cobra.Command, args []string) error {
		store := store.Store{}
		err := store.Read(configPath)
		if err != nil {
			log.Fatalf("Failed to read config file: %v", err)
		}

		err = store.ValidateWithDefaultValues()
		if err != nil {
			log.Fatalf("Failed to validate config file: %v", err)
		}
		fmt.Println(store.Repositories["brave-search"].Branch)
		return nil
	},
}

func init() {
	validateCmd.Flags().StringVarP(&configPath, "config", "c", "", "The path to the config file")
	rootCmd.AddCommand(validateCmd)
}
