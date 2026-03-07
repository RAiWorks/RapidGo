package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

// Config holds HTTP server settings.
type Config struct {
	Addr            string
	Handler         http.Handler
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// ListenAndServe starts an HTTP server and blocks until SIGINT or SIGTERM is
// received. It then initiates a graceful shutdown, waiting up to
// ShutdownTimeout for active connections to finish.
func ListenAndServe(cfg Config) error {
	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      cfg.Handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	// Create a context that is cancelled on SIGINT/SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start serving in the background.
	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	// Wait for a signal or a server error.
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	slog.Info("shutting down server…")
	stop() // Reset signal handling so a second signal force-kills.

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	slog.Info("server stopped")
	return nil
}

// ServiceConfig identifies a named HTTP service to run on a specific port.
type ServiceConfig struct {
	Name   string // "web", "api", "ws" — for logging
	Config Config // Standard server config (addr, handler, timeouts)
}

// ListenAndServeMulti starts multiple HTTP servers on separate ports and
// blocks until SIGINT/SIGTERM. All servers are shut down gracefully.
func ListenAndServeMulti(services []ServiceConfig) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	servers := make([]*http.Server, len(services))
	errCh := make(chan error, len(services))

	for i, svc := range services {
		srv := &http.Server{
			Addr:         svc.Config.Addr,
			Handler:      svc.Config.Handler,
			ReadTimeout:  svc.Config.ReadTimeout,
			WriteTimeout: svc.Config.WriteTimeout,
			IdleTimeout:  svc.Config.IdleTimeout,
		}
		servers[i] = srv

		go func(name string) {
			slog.Info("service starting", "name", name, "addr", srv.Addr)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				errCh <- fmt.Errorf("service %s: %w", name, err)
			}
		}(svc.Name)
	}

	// Wait for signal or server error
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	slog.Info("shutting down all services…")
	stop()

	shutdownTimeout := services[0].Config.ShutdownTimeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Shutdown all servers
	var firstErr error
	for i, srv := range servers {
		if err := srv.Shutdown(shutdownCtx); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("shutdown %s: %w", services[i].Name, err)
		}
	}
	slog.Info("all services stopped")
	return firstErr
}
