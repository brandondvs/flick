package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/brandondvs/flick/internal/config"
)

var configFileFlag string

func init() {
	rootCmd.PersistentFlags().StringVar(&configFileFlag, "config-file", "config.yaml", "path to configuration file")
}

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
	cobra.OnInitialize(func() {
		if err := config.Load(configFileFlag); err != nil {
			slog.Error("Failed to read configuration file at path", "error", err, "config-file-path", configFileFlag)
			os.Exit(1)
		}
	})

	if err := rootCmd.Execute(); err != nil {
		slog.Error("Root command failed to execute", "error", err)
		os.Exit(1)
	}
}
