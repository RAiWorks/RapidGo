# 📝 Changelog: Importable Library Split

> **Feature**: `57` — Importable Library Split
> **Branch**: `v2` (integration) with sub-branches per step
> **Started**: —
> **Completed**: —

---

## Log

<!-- Add entries as you work. Most recent first. -->

### Gate B — PASSED
- Zero `app/`, `routes/`, `http/`, `plugins/`, `database/models`, `database/seeders` imports in `core/`
- All 40 packages pass, all 7 coupling points resolved

### Step B4 — `feature/v2-06-migrate-decouple` → `v2` (commit `500d9f5`)
- `core/cli/migrate.go`: removed `database/models`, uses `modelRegistryFn()`
- `core/cli/seed.go`: removed `database/seeders`, uses `seederFn()`
- `cmd/main.go`: wired `SetModelRegistry(models.All)` + `SetSeeder()`
- `database/migrations/migrations_test.go`: replaced `models.User/Post/AuditLog` with `testMigrationModel`
- `database/seeders/seeders_test.go`: simplified `setupTestDB`, added `setupUserSeederDB`

### Step B3 — `feature/v2-05-worker-decouple` → `v2` (commit `034d062`)
- `core/cli/work.go`: removed `app/jobs`, `app/providers`, uses `NewApp()` + `jobRegistrar`
- `core/cli/schedule_run.go`: removed `app/providers`, `app/schedule`, uses `NewApp()` + `scheduleRegistrar`
- `cmd/main.go`: wired `SetJobRegistrar(jobs.RegisterJobs)` + `SetScheduleRegistrar(schedule.RegisterSchedule)`

### Step B2 — `feature/v2-04-serve-decouple` → `v2` (commit `d58019c`)
- `core/cli/serve.go`: removed `routes` import, uses `routeRegistrar` callback
- Static file/template setup stays in library
- `cmd/main.go`: wired `SetRoutes()` with mode-conditional Register calls

### Step B1 — `feature/v2-03-root-decouple` → `v2` (commit `a62995c`)
- `core/cli/root.go`: removed `app/providers`, `NewApp()` uses `bootstrapFn`
- `core/cli/cli_test.go`: updated `TestNewApp_ReturnsBootedApp` to use `SetBootstrap` with test bindings
- `cmd/main.go`: wired `SetBootstrap()` with all provider registrations

### Step A2 — `feature/v2-02-audit-decouple` → `v2` (commit `932fdc0`)
- Created `core/audit/model.go` with canonical `AuditLog` struct
- `core/audit/audit.go` + `audit_test.go`: removed `database/models` import
- `database/models/audit_log.go`: replaced struct with `type AuditLog = audit.AuditLog`

### Step A1 — `feature/v2-01-hooks-foundation` → `v2` (commit `2e93063`)
- Created `core/cli/hooks.go`: 6 callback types, 6 package-level vars, 6 `Set*()` functions
- Created `core/cli/hooks_test.go`: 7 tests (nil defaults + setter storage)

---

## Deviations from Plan

| What Changed | Original Plan | What Actually Happened | Why |
|---|---|---|---|

## Key Decisions Made During Build

| Decision | Context | Date |
|---|---|---|
