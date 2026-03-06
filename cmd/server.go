package cmd

import (
	"log/slog"
	"net"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/brandondvs/flick/internal/config"
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

		port := strconv.Itoa(config.ServerPort())
		hostPort := net.JoinHostPort(config.ServerHost(), port)

		slog.Info("Server running", "host", hostPort)

		if err := http.ListenAndServe(hostPort, srv); err != nil {
			slog.Error("Failed to start listener", "error", err)
		}
	},
}
