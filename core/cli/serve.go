package cli

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/RAiWorks/RapidGo/v2/core/app"
	"github.com/RAiWorks/RapidGo/v2/core/config"
	"github.com/RAiWorks/RapidGo/v2/core/container"
	"github.com/RAiWorks/RapidGo/v2/core/health"
	"github.com/RAiWorks/RapidGo/v2/core/router"
	"github.com/RAiWorks/RapidGo/v2/core/server"
	"github.com/RAiWorks/RapidGo/v2/core/service"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var servePort string
var serveMode string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load .env so RAPIDGO_MODE can be read from .env file
		config.Load()

		// Resolve mode: CLI flag > env var > default "all"
		modeStr := config.Env("RAPIDGO_MODE", "all")
		if serveMode != "" {
			modeStr = serveMode
		}

		mode, err := service.ParseMode(modeStr)
		if err != nil {
			return fmt.Errorf("invalid service mode: %w", err)
		}

		application := NewApp(mode)

		appName := config.Env("APP_NAME", "RapidGo")
		appEnv := config.AppEnv()

		// Delegate to single-port or multi-port based on active services
		services := mode.Services()
		if len(services) <= 1 || allSamePort(services) {
			port := resolvePort(mode)

			fmt.Println("=================================")
			fmt.Printf("  %s Framework\n", appName)
			fmt.Println("  github.com/RAiWorks/RapidGo")
			fmt.Println("=================================")
			fmt.Printf("  Environment: %s\n", appEnv)
			fmt.Printf("  Mode: %s\n", mode)
			fmt.Printf("  Port: %s\n", port)
			fmt.Printf("  Debug: %v\n", config.IsDebug())
			fmt.Println("=================================")

			slog.Info("server starting",
				"app", appName,
				"mode", mode.String(),
				"port", port,
				"env", appEnv,
			)

			return serveSingle(application, mode)
		}

		// Multi-port — one server per service
		fmt.Println("=================================")
		fmt.Printf("  %s Framework\n", appName)
		fmt.Println("  github.com/RAiWorks/RapidGo")
		fmt.Println("=================================")
		fmt.Printf("  Environment: %s\n", appEnv)
		fmt.Printf("  Mode: %s\n", mode)
		for _, svc := range services {
			fmt.Printf("  %s → :%s\n", svc.String(), resolvePortForMode(svc))
		}
		fmt.Printf("  Debug: %v\n", config.IsDebug())
		fmt.Println("=================================")

		slog.Info("server starting (multi-port)",
			"app", appName,
			"mode", mode.String(),
			"env", appEnv,
		)

		return serveMulti(application, mode)
	},
}

func init() {
	serveCmd.Flags().StringVarP(&servePort, "port", "p", "", "port to listen on (overrides APP_PORT)")
	serveCmd.Flags().StringVarP(&serveMode, "mode", "m", "", "service mode: all, web, api, ws, or comma-separated (overrides RAPIDGO_MODE)")
}

// serveSingle starts one HTTP server on a single port (backward compatible).
func serveSingle(application *app.App, mode service.Mode) error {
	port := resolvePort(mode)
	r := container.MustMake[*router.Router](application.Container, "router")
	return server.ListenAndServe(server.Config{
		Addr:            ":" + port,
		Handler:         r,
		ReadTimeout:     15 * time.Second,
		WriteTimeout:    15 * time.Second,
		IdleTimeout:     60 * time.Second,
		ShutdownTimeout: 30 * time.Second,
	})
}

// serveMulti starts separate HTTP servers per service on different ports.
func serveMulti(application *app.App, mode service.Mode) error {
	var services []server.ServiceConfig

	for _, svc := range mode.Services() {
		r := router.New()
		applyRoutesForMode(r, application.Container, svc)
		port := resolvePortForMode(svc)

		services = append(services, server.ServiceConfig{
			Name: svc.String(),
			Config: server.Config{
				Addr:            ":" + port,
				Handler:         r,
				ReadTimeout:     15 * time.Second,
				WriteTimeout:    15 * time.Second,
				IdleTimeout:     60 * time.Second,
				ShutdownTimeout: 30 * time.Second,
			},
		})
	}

	return server.ListenAndServeMulti(services)
}

// resolvePort returns the port for the active mode.
// Single-mode uses mode-specific port env var, else APP_PORT.
func resolvePort(mode service.Mode) string {
	if servePort != "" {
		return servePort
	}
	services := mode.Services()
	if len(services) == 1 {
		return config.Env(services[0].PortEnvKey(), config.Env("APP_PORT", "8080"))
	}
	return config.Env("APP_PORT", "8080")
}

// resolvePortForMode returns the port for a specific service mode.
func resolvePortForMode(m service.Mode) string {
	return config.Env(m.PortEnvKey(), config.Env("APP_PORT", "8080"))
}

// applyRoutesForMode registers routes on a router for a specific mode.
func applyRoutesForMode(r *router.Router, c *container.Container, m service.Mode) {
	if m.Has(service.ModeWeb) {
		r.SetFuncMap(router.DefaultFuncMap())
		viewsDir := filepath.Join("resources", "views")
		if info, err := os.Stat(viewsDir); err == nil && info.IsDir() {
			r.LoadTemplates(viewsDir)
		}
		if info, err := os.Stat("resources/static"); err == nil && info.IsDir() {
			r.Static("/static", "./resources/static")
		}
		if info, err := os.Stat("storage/uploads"); err == nil && info.IsDir() {
			r.Static("/uploads", "./storage/uploads")
		}
	}

	// Delegate route registration to the application callback
	if routeRegistrar != nil {
		routeRegistrar(r, c, m)
	}

	// Health check — each per-service router gets its own health endpoints
	if c.Has("db") {
		health.Routes(r, func() *gorm.DB {
			return container.MustMake[*gorm.DB](c, "db")
		})
	}
}

// allSamePort returns true if all services resolve to the same port.
func allSamePort(services []service.Mode) bool {
	if len(services) <= 1 {
		return true
	}
	port := config.Env(services[0].PortEnvKey(), config.Env("APP_PORT", "8080"))
	for _, s := range services[1:] {
		if config.Env(s.PortEnvKey(), config.Env("APP_PORT", "8080")) != port {
			return false
		}
	}
	return true
}
