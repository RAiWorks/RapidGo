# Feature #37 — Graceful Shutdown: Design

## Package

`core/server/server.go`

## Public API

```go
// Config holds server configuration.
type Config struct {
    Addr             string
    Handler          http.Handler
    ReadTimeout      time.Duration
    WriteTimeout     time.Duration
    IdleTimeout      time.Duration
    ShutdownTimeout  time.Duration
}

// ListenAndServe starts an HTTP server and blocks until a shutdown signal is
// received. It returns nil on clean shutdown.
func ListenAndServe(cfg Config) error
```

## Integration

`core/cli/serve.go` replaces `r.Run()` with:

```go
err := server.ListenAndServe(server.Config{
    Addr:            ":" + port,
    Handler:         r,
    ReadTimeout:     15 * time.Second,
    WriteTimeout:    15 * time.Second,
    IdleTimeout:     60 * time.Second,
    ShutdownTimeout: 30 * time.Second,
})
```

## Shutdown Sequence

1. Receive SIGINT or SIGTERM.
2. Log "shutting down server…".
3. Call `srv.Shutdown(ctx)` with `ShutdownTimeout`.
4. Close DB connection if container has `"db"` binding.
5. Log "server stopped".

## File Layout

```
core/server/
  server.go       — Config, ListenAndServe
  server_test.go  — Tests for startup, shutdown, timeouts
```
