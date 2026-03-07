# 📝 Changelog: Plugin / Module System

> **Feature**: `44` — Plugin / Module System
> **Status**: � SHIPPED
> **Date**: 2026-03-07
> **Commit**: `0905794` (merge to main)

---

## Added

- `core/plugin/plugin.go` — `Plugin` interface (embeds `container.Provider` + `Name()`), optional interfaces (`RouteRegistrar`, `CommandRegistrar`, `EventRegistrar`), `PluginManager` with `Add()`, `RegisterAll()`, `BootAll()`, `RegisterRoutes()`, `RegisterCommands()`, `RegisterEvents()`
- `core/plugin/plugin_test.go` — 18 tests covering manager construction, plugin registration, duplicate detection, lifecycle ordering, service binding, type-asserted route/command/event wiring
- `plugins/example/example.go` — `ExamplePlugin` implementing Plugin, RouteRegistrar, CommandRegistrar with demo route (`/example/hello`) and CLI command (`example:greet`)
- `app/plugins.go` — `RegisterPlugins()` application entry point

## Changed

- `core/cli/root.go` — added `RootCmd()` accessor function for plugin CLI command registration

## Files

| File | Action |
|------|--------|
| `core/plugin/plugin.go` | NEW |
| `core/plugin/plugin_test.go` | NEW |
| `plugins/example/example.go` | NEW |
| `app/plugins.go` | NEW |
| `core/cli/root.go` | MODIFIED |

## Migration Guide

No migration needed. This is a new feature with no breaking changes.
