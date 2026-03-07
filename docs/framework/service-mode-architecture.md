# RGo Framework — Service Mode Architecture Discussion

> **Status**: DISCUSSION / RFC  
> **Date**: 2026-03-07  
> **Author**: Architecture Review  
> **Scope**: Can RGo run as monolith, split microservices, or combined service subsets?

---

## 1. The Question

Can we build a full application (Web SSR + APIs + WebSocket + Workers) with RGo and later:

1. Run **everything in one binary** (monolith)
2. Split and run **each as its own service** (microservice)
3. Run **combined subsets** (e.g., API + WebSocket together, Web SSR separately)

**Answer: YES — but the framework needs targeted changes first.**

The underlying technology (Go, Gin, GORM) fully supports this. The blockers are in RGo's bootstrapping layer, not in the core components.

---

## 2. Current Architecture

### How the app starts today

```
cmd/main.go → cli.Execute() → serve command → NewApp()
    ↓
NewApp() registers ALL 6 providers (hardcoded, no choice):
    1. ConfigProvider      → loads .env
    2. LoggerProvider      → sets up slog
    3. DatabaseProvider    → registers "db" singleton
    4. SessionProvider     → registers "session" (REQUIRES "db")
    5. MiddlewareProvider  → registers global middleware maps
    6. RouterProvider      → creates ONE router, registers ALL routes
    ↓
application.Boot() → calls Boot() on all 6 in order
    ↓
serve command → extracts "router" → server.ListenAndServe(:8080)
```

**Result**: Every `rgo serve` starts the FULL app — Web + API + WebSocket — on ONE port. No way to customize.

### Provider dependency graph

```
ConfigProvider ─────────────────────────────────────────┐
    │                                                    │
    ├──→ LoggerProvider (reads LOG_* from config)        │
    │                                                    │
    ├──→ DatabaseProvider (reads DB_* from config)       │
    │         │                                          │
    │         └──→ SessionProvider (HARD depends on DB)  │
    │                                                    │
    ├──→ MiddlewareProvider (independent)                │
    │                                                    │
    └──→ RouterProvider                                  │
              ├── RegisterWeb()  ← always called         │
              ├── RegisterAPI()  ← always called         │
              └── health.Routes() ← if DB exists         │
```

### What blocks service splitting today

| # | Blocker | Location | Why it matters |
|---|---------|----------|----------------|
| 1 | **Hardcoded provider list** | `core/cli/root.go` NewApp() | All 6 providers ALWAYS register — can't skip DB for a stateless API service |
| 2 | **SessionProvider requires DB** | `app/providers/session_provider.go` | Calls `MustMake[*gorm.DB](c, "db")` — panics if no DB registered |
| 3 | **Single router key** | RouterProvider registers `"router"` | Only one router per container — can't have `"router:api"` and `"router:web"` |
| 4 | **Routes always load together** | RouterProvider calls both `RegisterWeb()` AND `RegisterAPI()` | Can't run API-only or Web-only |
| 5 | **Global middleware registry** | `core/middleware/registry.go` uses package-level maps | All services in same process share middleware namespace |
| 6 | **Global named routes** | `core/router/named.go` uses package-level map | Route names collide across services in same process |
| 7 | **Single serve command** | `core/cli/serve.go` starts one server on one port | Can't run Web on :8080 and API on :3001 simultaneously |
| 8 | **No worker/queue mode** | No CLI command for background processing | Can't run `rgo worker` for job processing |

---

## 3. What the Technology Supports

### Can Gin serve a subset of routes?

**YES.** Gin engines are independent route trees. You can create multiple engines with different routes:

```go
// API service — only API routes
apiEngine := gin.New()
apiEngine.GET("/api/users", handlers.ListUsers)
apiEngine.POST("/api/users", handlers.CreateUser)

// Web service — only web routes  
webEngine := gin.New()
webEngine.GET("/", controllers.Home)
webEngine.GET("/about", controllers.About)
```

### Can we run multiple Gin engines in one process?

**YES.** Each engine is fully independent. Run them on different ports:

```go
go server.ListenAndServe(server.Config{Addr: ":8080", Handler: webEngine})
go server.ListenAndServe(server.Config{Addr: ":3001", Handler: apiEngine})
go server.ListenAndServe(server.Config{Addr: ":3002", Handler: wsEngine})
select {} // block forever
```

### Can WebSocket run on a different port?

**YES.** WebSocket is just an HTTP upgrade on a route. Put it on any engine:

```go
wsEngine := gin.New()
wsEngine.GET("/ws", websocket.Upgrader(chatHandler, nil))
server.ListenAndServe(server.Config{Addr: ":3002", Handler: wsEngine})
```

### Can we have isolated containers?

**YES.** The container is a plain struct — create multiple:

```go
apiContainer := container.New()
apiContainer.Singleton("db", apiDBFactory)

webContainer := container.New()
// No DB needed for web — skip it
```

---

## 4. Proposed Solution: Service Modes

### The concept

Add a `RGO_MODE` environment variable (or `--mode` flag) that controls which services start:

```
RGO_MODE=all        →  Web + API + WebSocket (monolith, current behavior)
RGO_MODE=web        →  Web SSR only (templates, static files)
RGO_MODE=api        →  API only (JSON endpoints)
RGO_MODE=ws         →  WebSocket only
RGO_MODE=worker     →  Background job processor only
RGO_MODE=api,ws     →  API + WebSocket combined
RGO_MODE=web,api    →  Web + API (no WebSocket)
```

### How it would work

```
rgo serve                     →  starts in RGO_MODE (default: all)
rgo serve --mode=api          →  API-only service on API_PORT
rgo serve --mode=web          →  Web-only service on WEB_PORT
rgo serve --mode=api,ws       →  API + WebSocket on separate ports
rgo worker                    →  Job processing (no HTTP)
```

### Architecture for multi-mode support

```
┌─────────────────────────────────────────────────────────┐
│                    rgo serve --mode=all                  │
│                                                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────┐ │
│  │ Web SSR  │  │   API    │  │WebSocket │  │ Worker │ │
│  │  :8080   │  │  :3001   │  │  :3002   │  │  (bg)  │ │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └───┬────┘ │
│       │              │              │             │      │
│  ┌────┴──────────────┴──────────────┴─────────────┴───┐ │
│  │              Shared Container                      │ │
│  │   "db"  "cache"  "session"  "events"  "queue"     │ │
│  └────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘

─── OR ───

┌─────────────┐   ┌─────────────┐   ┌─────────────┐
│  rgo serve  │   │  rgo serve  │   │  rgo serve  │
│  --mode=web │   │  --mode=api │   │  --mode=ws  │
│    :8080    │   │    :3001    │   │    :3002    │
│             │   │             │   │             │
│  Own .env   │   │  Own .env   │   │  Own .env   │
│  Own DB?    │   │  Own DB     │   │  No DB      │
│  Templates  │   │  JSON only  │   │  WebSocket  │
└──────┬──────┘   └──────┬──────┘   └──────┬──────┘
       │                  │                  │
       └──────── Shared Redis/DB ────────────┘
```

---

## 5. Changes Required

### Phase 1: Optional Providers (SMALL — ~1 day)

**Goal**: Allow providers to be skipped based on configuration.

**Change 1**: Add `Enabled()` method to Provider interface:

```go
// core/container/provider.go
type Provider interface {
    Register(c *Container)
    Boot(c *Container)
}

// New optional interface:
type ConditionalProvider interface {
    Provider
    Enabled() bool  // Return false to skip this provider
}
```

**Change 2**: Update `App.Register()` to check:

```go
func (a *App) Register(provider Provider) {
    if cp, ok := provider.(ConditionalProvider); ok && !cp.Enabled() {
        return // Skip disabled provider
    }
    a.providers = append(a.providers, provider)
    provider.Register(a.Container)
}
```

**Change 3**: Fix SessionProvider to not panic without DB:

```go
func (p *SessionProvider) Register(c *container.Container) {
    c.Singleton("session", func(c *container.Container) interface{} {
        var db *gorm.DB
        if c.Has("db") {
            db, _ = container.Make[*gorm.DB](c, "db")
        }
        store, _ := session.NewStore(db)
        return session.NewManager(store)
    })
}
```

### Phase 2: Service Mode in Serve Command (SMALL — ~1 day)

**Goal**: `rgo serve --mode=api` only loads API routes.

**Change 1**: Add `--mode` flag to serve command:

```go
var serveCmd = &cobra.Command{
    Use:   "serve",
    Short: "Start the HTTP server",
    Run: func(cmd *cobra.Command, args []string) {
        mode, _ := cmd.Flags().GetString("mode")
        if mode == "" {
            mode = os.Getenv("RGO_MODE")
        }
        if mode == "" {
            mode = "all"
        }
        // ...
    },
}

func init() {
    serveCmd.Flags().String("mode", "", "Service mode: all, web, api, ws, worker")
}
```

**Change 2**: Make RouterProvider mode-aware:

```go
func (p *RouterProvider) Boot(c *container.Container) {
    r := container.MustMake[*router.Router](c, "router")
    mode := config.Get("RGO_MODE", "all")

    r.SetFuncMap(router.DefaultFuncMap())

    if mode == "all" || strings.Contains(mode, "web") {
        // Load templates and static files
        viewsDir := filepath.Join("resources", "views")
        if info, err := os.Stat(viewsDir); err == nil && info.IsDir() {
            r.LoadTemplates(viewsDir)
        }
        if info, err := os.Stat("resources/static"); err == nil && info.IsDir() {
            r.Static("/static", "./resources/static")
        }
        routes.RegisterWeb(r)
    }

    if mode == "all" || strings.Contains(mode, "api") {
        routes.RegisterAPI(r)
    }

    if mode == "all" || strings.Contains(mode, "ws") {
        routes.RegisterWS(r) // New: WebSocket routes
    }
}
```

### Phase 3: Multi-Port Serving (MEDIUM — ~2 days)

**Goal**: Run different services on different ports from one binary.

**Change**: Extend serve command to start multiple servers:

```go
// Service definitions from env or config
// WEB_PORT=8080  API_PORT=3001  WS_PORT=3002

func runMultiMode(app *app.App) {
    modes := strings.Split(config.Get("RGO_MODE", "all"), ",")
    
    var wg sync.WaitGroup
    
    for _, mode := range modes {
        switch strings.TrimSpace(mode) {
        case "web":
            wg.Add(1)
            go func() {
                defer wg.Done()
                webRouter := buildWebRouter(app)
                server.ListenAndServe(server.Config{
                    Addr:    ":" + config.Get("WEB_PORT", "8080"),
                    Handler: webRouter,
                })
            }()
        case "api":
            wg.Add(1)
            go func() {
                defer wg.Done()
                apiRouter := buildAPIRouter(app)
                server.ListenAndServe(server.Config{
                    Addr:    ":" + config.Get("API_PORT", "3001"),
                    Handler: apiRouter,
                })
            }()
        case "ws":
            wg.Add(1)
            go func() {
                defer wg.Done()
                wsRouter := buildWSRouter(app)
                server.ListenAndServe(server.Config{
                    Addr:    ":" + config.Get("WS_PORT", "3002"),
                    Handler: wsRouter,
                })
            }()
        }
    }
    
    wg.Wait()
}
```

### Phase 4: Remove Global State (MEDIUM — ~1 day)

**Goal**: Allow multiple services in same process without conflicts.

**Change 1**: Move middleware registry into container:

```go
// Instead of package-level maps:
type MiddlewareRegistry struct {
    aliases map[string]gin.HandlerFunc
    groups  map[string][]gin.HandlerFunc
}

// Register in container:
c.Instance("middleware", NewMiddlewareRegistry())
```

**Change 2**: Move named routes into Router:

```go
// Instead of package-level map:
type Router struct {
    engine      *gin.Engine
    namedRoutes map[string]string  // Per-router named routes
}
```

### Phase 5: Worker/Queue Mode (LARGE — ~3-5 days)

**Goal**: `rgo worker` processes background jobs.

This requires a new subsystem:

```go
// core/queue/queue.go
type Job interface {
    Handle() error
}

type Queue interface {
    Push(job Job) error
    Pop() (Job, error)
    Listen(handler func(Job)) error
}

// Backends: memory, redis, database
```

---

## 6. Deployment Scenarios

### Scenario A: Monolith (Today + Phase 1-2)

```
.env:
  RGO_MODE=all
  APP_PORT=8080

$ rgo serve
→ Starts everything on :8080
→ Web SSR + API + WebSocket all served from one binary
→ Perfect for: MVP, small teams, simple deployments
```

### Scenario B: Split Services (Phase 2-3)

```
# Service 1: Web Frontend
.env:
  RGO_MODE=web
  WEB_PORT=8080
  DB_HOST=shared-db.internal

# Service 2: API Backend
.env:
  RGO_MODE=api
  API_PORT=3001
  DB_HOST=shared-db.internal

# Service 3: WebSocket Server
.env:
  RGO_MODE=ws
  WS_PORT=3002
  REDIS_HOST=shared-redis.internal

→ Each runs independently
→ Scale API separately from Web
→ Deploy WebSocket on different infrastructure
```

### Scenario C: Combined Subsets (Phase 3)

```
# Service 1: Web + API together
.env:
  RGO_MODE=web,api
  WEB_PORT=8080
  API_PORT=3001

# Service 2: WebSocket + Worker
.env:
  RGO_MODE=ws,worker
  WS_PORT=3002

→ Flexible grouping based on traffic patterns
→ Co-locate related services
```

### Scenario D: Separate Binaries (Phase 2 + build tags)

```
# Build API-only binary (smaller, faster startup)
$ RGO_MODE=api go build -o api-server ./cmd

# Build Web-only binary
$ RGO_MODE=web go build -o web-server ./cmd

→ Each binary only contains needed code
→ Smaller container images
→ Independent deployment pipelines
```

---

## 7. What We DON'T Need to Build

Some things are better solved by infrastructure, not framework:

| Concern | Framework solution? | Better solution |
|---------|-------------------|-----------------|
| Service discovery | ❌ No | Kubernetes / Consul / DNS |
| Load balancing | ❌ No | Nginx / HAProxy / K8s Service |
| Service-to-service auth | ❌ No | mTLS / service mesh (Istio) |
| Distributed tracing | ❌ No | OpenTelemetry / Jaeger |
| Config management | ❌ No | Vault / ConfigMap / .env per service |
| Container orchestration | ❌ No | Docker Compose / Kubernetes |
| gRPC between services | ❌ No | Use gRPC directly if needed |
| API Gateway | ❌ No | Kong / Traefik / Caddy |

RGo should focus on making the **application code** splittable, not reimplementing infrastructure.

---

## 8. Comparison with Other Frameworks

| Capability | Laravel | Spring Boot | NestJS | RGo (Today) | RGo (Proposed) |
|------------|---------|-------------|--------|-------------|----------------|
| Monolith mode | ✅ | ✅ | ✅ | ✅ | ✅ |
| Optional providers | ✅ | ✅ | ✅ | ❌ | ✅ Phase 1 |
| Service modes | ✅ (artisan) | ✅ (profiles) | ✅ (modules) | ❌ | ✅ Phase 2 |
| Multi-port | ❌ | ✅ | ✅ | ❌ | ✅ Phase 3 |
| Background workers | ✅ (queues) | ✅ | ✅ (bull) | ❌ | ✅ Phase 5 |
| Microservice toolkit | ⚠️ Lumen | ✅ (cloud) | ✅ | ❌ | ✅ Phase 2-3 |

---

## 9. Implementation Priority

| Phase | Effort | Impact | Recommendation |
|-------|--------|--------|----------------|
| Phase 1: Optional Providers | ~1 day | HIGH | **Do first** — unblocks everything |
| Phase 2: Service Mode Flag | ~1 day | HIGH | **Do second** — enables API-only / Web-only |
| Phase 3: Multi-Port Serving | ~2 days | MEDIUM | Do when actually splitting services |
| Phase 4: Remove Global State | ~1 day | LOW | Only needed for multi-service in same process |
| Phase 5: Worker/Queue | ~3-5 days | HIGH | Do when background processing needed |

**Total for Phase 1-3**: ~4 days of work to go from monolith-only to fully splittable.

---

## 10. Decision

**Is it possible?** YES — the Go/Gin/GORM stack fully supports this pattern.

**Is it hard?** NO — the changes are surgical, not architectural rewrites. The container/provider pattern is already the right foundation.

**Should we do it now?** Phase 1-2 (2 days) should be done proactively. Phase 3-5 can wait until actually needed.

**The key insight**: RGo's architecture is already 80% there. The container, providers, router, and server are all properly decoupled. The remaining 20% is just removing a few hardcoded assumptions in the bootstrap layer.
