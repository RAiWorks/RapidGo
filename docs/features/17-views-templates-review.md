# 📋 Review: Views & Templates

> **Feature**: `17` — Views & Templates
> **Branch**: `feature/17-views-templates`
> **Merged**: 2026-03-06 → `main` (`e6744ce`)
> **Status**: ✅ SHIPPED

---

## Summary

Added server-side rendering support using Go's `html/template` package. The Router now supports template loading, template functions (with the `route` helper), and static file serving. A sample `home.html` template replaces the previous JSON response from the Home controller.

## Files Changed

| File | Action | Purpose |
|---|---|---|
| `core/router/view.go` | Created | `DefaultFuncMap()` with `route` template helper |
| `core/router/view_test.go` | Created | 5 tests (FuncMap, template rendering, static serving) |
| `core/router/router.go` | Modified | Added `SetFuncMap`, `LoadTemplates`, `Static`, `StaticFile` methods |
| `app/providers/router_provider.go` | Modified | Template engine + static serving in Boot |
| `http/controllers/home_controller.go` | Modified | `c.JSON()` → `c.HTML()` with `home.html` |
| `http/controllers/controllers_test.go` | Modified | Updated Home test for HTML, added template setup helper |
| `resources/views/home.html` | Created | Sample HTML5 template |
| `docs/project-roadmap.md` | Modified | #17 → ✅ |

## Test Results

- **New tests**: 5 (in `view_test.go`)
- **Updated tests**: 1 (`TestHome_RendersHTML` replaces `TestHome_ReturnsWelcome`)
- **Total tests**: 177
- **Failures**: 0

## Deviations from Architecture

| What Changed | Why |
|---|---|
| Provider guards `LoadTemplates` with `filepath.Glob` check | `LoadHTMLGlob` panics when pattern matches no files; tests run from package dirs without `resources/views/` |
| Provider guards `Static` with `os.Stat` check | Same issue — directories don't exist in test contexts |

## Key Design Points

- `SetFuncMap` must be called before `LoadTemplates` (Gin requirement)
- Template functions include `route` for named route URL generation
- Static serving: `/static` → `resources/static/`, `/uploads` → `storage/uploads/`
- Home controller now renders HTML via `c.HTML()` — proves full SSR pipeline
