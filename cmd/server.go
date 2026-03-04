package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.AddCommand(serverStartCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Command to control the server",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			slog.Error("Failed to call cmd.Help() within the server command")
		}
	},
}

var serverStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Command to start the flick HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		slog.Info("Starting server")
	},
}
