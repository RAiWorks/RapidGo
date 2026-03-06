# Feature #30 — Static File Serving: Architecture

## Status: Already implemented in Features #07 and #17

No new files or modifications required.

## Existing implementation

| File | Purpose | Shipped in |
|------|---------|------------|
| `core/router/router.go` | `Static()`, `StaticFile()` methods | #07 |
| `app/providers/router_provider.go` | `/static` + `/uploads` route setup | #07, #17 |
| `core/router/view_test.go` | TC-04, TC-05 — static serving tests | #07 |
| `resources/static/.gitkeep` | Static assets directory | #07 |
| `storage/uploads/.gitkeep` | User uploads directory | #05 |

## Request flow

```
Browser → GET /static/style.css
    │
    ▼
RouterProvider.Boot()
    └── r.Static("/static", "./resources/static")
        └── gin.Engine.Static() → serves resources/static/style.css

Browser → GET /uploads/photo.jpg
    │
    ▼
RouterProvider.Boot()
    └── r.Static("/uploads", "./storage/uploads")
        └── gin.Engine.Static() → serves storage/uploads/photo.jpg
```

## Blueprint comparison

| Blueprint spec | Our implementation | Match |
|---|---|---|
| `r.Static("/static", ...)` | `r.Static("/static", "./resources/static")` | ✅ |
| `r.Static("/uploads", ...)` | `r.Static("/uploads", "./storage/uploads")` | ✅ |
| `r.StaticFile("/favicon.ico", ...)` | Omitted (no favicon.ico exists) | ⚠️ Intentional |
| `r.LoadHTMLGlob(...)` | `r.LoadTemplates(pattern)` | ✅ (equivalent) |
