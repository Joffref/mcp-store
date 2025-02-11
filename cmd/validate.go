package cmd

import (
	"fmt"

	"github.com/beamlit/mcp-store/internal/store"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the MCP store",
	Run:   runValidate,
}

func init() {
	validateCmd.Flags().StringVarP(&configPath, "config", "c", "", "The path to the config file")
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) {
	if configPath == "" {
		cmd.Help()
		return
	}

	store := store.Store{}
	handleError("read config file", store.Read(configPath))
	handleError("validate config file", store.ValidateWithDefaultValues())

	// Print validation success message
	fmt.Println("Configuration validated successfully")
}
