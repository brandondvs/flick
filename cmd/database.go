package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/brandondvs/flick/internal/config"
	"github.com/brandondvs/flick/internal/database"
)

func init() {
	rootCmd.AddCommand(databaseCmd)

	databaseCmd.AddCommand(databaseValidateCmd)
}

var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Database commands for testing and managing database schema",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var databaseValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate that there's a healthy database connection and the database schema is correct",
	Run: func(cmd *cobra.Command, args []string) {
		db, err := database.New()
		if err != nil {
			slog.Error("Failed to create database connection", "error", err)
			os.Exit(1)
		}

		slog.Info("Checking database connection", "connection_string", config.DatabaseConnectionString())
		if err := db.Ping(); err != nil {
			slog.Error("Failed to ping database", "error", err)
			os.Exit(1)
		}
		slog.Info("Connection valid!")
	},
}
