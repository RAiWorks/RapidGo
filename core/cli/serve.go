package cli

import (
	"fmt"
	"log/slog"

	"github.com/RAiWorks/RGo/core/config"
	"github.com/RAiWorks/RGo/core/container"
	"github.com/RAiWorks/RGo/core/router"
	"github.com/spf13/cobra"
)

var servePort string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		application := NewApp()

		port := config.Env("APP_PORT", "8080")
		if servePort != "" {
			port = servePort
		}

		appName := config.Env("APP_NAME", "RGo")
		appEnv := config.AppEnv()

		fmt.Println("=================================")
		fmt.Printf("  %s Framework\n", appName)
		fmt.Println("  github.com/RAiWorks/RGo")
		fmt.Println("=================================")
		fmt.Printf("  Environment: %s\n", appEnv)
		fmt.Printf("  Port: %s\n", port)
		fmt.Printf("  Debug: %v\n", config.IsDebug())
		fmt.Println("=================================")

		slog.Info("server starting",
			"app", appName,
			"port", port,
			"env", appEnv,
		)

		r := container.MustMake[*router.Router](application.Container, "router")
		if err := r.Run(":" + port); err != nil {
			slog.Error("server failed to start", "err", err)
		}
	},
}

func init() {
	serveCmd.Flags().StringVarP(&servePort, "port", "p", "", "port to listen on (overrides APP_PORT)")
}
