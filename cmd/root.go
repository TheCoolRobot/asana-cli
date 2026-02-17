package cmd

import (
	

	"github.com/TheCoolRobot/asana-cli/internal/config"
	"github.com/spf13/cobra"
)
// "fmt"
// 	"os"
var (
	jsonOutput  bool
	token       string
	workspace   string
	project     string
)

var rootCmd = &cobra.Command{
	Use:     "asana-cli",
	Short:   "Asana CLI - Beautiful task management",
	Long:    "A feature-rich CLI for managing Asana tasks with TUI and sync daemon",
	Version: "0.1.0",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if token == "" {
			token = config.GetAPIToken()
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "Asana API token (or set ASANA_TOKEN)")
	rootCmd.PersistentFlags().StringVar(&workspace, "workspace", "", "Default workspace ID")
	rootCmd.PersistentFlags().StringVar(&project, "project", "", "Default project ID")

	// Add all commands
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(completeCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(viewCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(meCmd)
}

func Execute() error {
	return rootCmd.Execute()
}