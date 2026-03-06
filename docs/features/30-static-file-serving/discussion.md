# Feature #30 — Static File Serving: Discussion

## What problem does this solve?

Web applications need to serve static assets (CSS, JS, images) and user-uploaded files. Static file serving maps URL paths to filesystem directories so browsers can fetch these resources.

## Why now?

This feature was already fully implemented across Feature #07 (Router) and Feature #17 (Views & Templates). The roadmap listed it as a separate item because the blueprint has a dedicated "Static File Serving" section, but our existing implementation covers 100% of the blueprint specification.

## What does the blueprint specify?

- `r.Static("/static", "./resources/static")` — serve CSS, JS, images.
- `r.Static("/uploads", "./storage/uploads")` — serve user uploads.
- `r.StaticFile("/favicon.ico", "./resources/static/favicon.ico")` — serve favicon.
- `r.LoadHTMLGlob("resources/views/**/*")` — template loading.

## What do we already have?

### Router methods (shipped in #07)
- `Router.Static(urlPath, dirPath)` — wraps `gin.Engine.Static()`.
- `Router.StaticFile(urlPath, filePath)` — wraps `gin.Engine.StaticFile()`.

### RouterProvider.Boot() (shipped in #07, enhanced in #17)
- Conditionally serves `/static` → `resources/static/` (only if dir exists).
- Conditionally serves `/uploads` → `storage/uploads/` (only if dir exists).
- Template engine with `SetFuncMap()` and `LoadTemplates()`.

### Tests (shipped in #07)
- TC-04: `TestStatic_ServesFiles` — verifies static file serving.
- TC-05: `TestStaticFile_ServesSingleFile` — verifies single-file serving.

### Favicon
Intentionally omitted in #17: "no favicon.ico exists yet; can be added when assets are created." Non-blocking.

## Design decision

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Status | Already shipped (#07, #17) | Full blueprint coverage |
| Favicon | Omitted | No favicon.ico file exists; route can be added when one is created |
| Code changes | None needed | Implementation is complete |
