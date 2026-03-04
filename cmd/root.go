package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "flick",
	Short: "Flick is a feature flag HTTP service",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			slog.Error("Failed to call cmd.Help() within the root command", "error", err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		slog.Error("Root command failed to execute", "error", err)
		os.Exit(1)
	}
}
