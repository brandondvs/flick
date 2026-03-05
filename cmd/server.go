package cmd

import (
	"log/slog"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/brandondvs/flick/internal/server"
	"github.com/brandondvs/flick/internal/store"
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
			slog.Error("Failed to call cmd.Help() within the server command", "error", err)
		}
	},
}

var serverStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Command to start the flick HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		s := store.New()
		srv := server.New(s)

		slog.Info("Server running")
		if err := http.ListenAndServe(":8080", srv); err != nil {
			slog.Error("Failed to start listener", "error", err)
		}
	},
}
