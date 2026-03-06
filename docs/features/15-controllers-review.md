# 📋 Review: Controllers

> **Feature**: `15` — Controllers
> **Branch**: `feature/15-controllers`
> **Merged**: 2026-03-06
> **Commit**: `92afdd3` (main)

---

## Summary

Feature #15 adds MVC controllers to `http/controllers/`. A `Home` function controller handles `GET /` with a JSON welcome message, and a `PostController` struct implements all 7 `ResourceController` methods for RESTful CRUD. Routes are wired in `routes/web.go` and `routes/api.go`.

## Files Changed

| File | Type | Description |
|---|---|---|
| `http/controllers/home_controller.go` | Created | `Home()` — JSON welcome on `GET /` |
| `http/controllers/post_controller.go` | Created | `PostController` — 7 CRUD methods |
| `routes/web.go` | Modified | Register `GET /` → `controllers.Home` |
| `routes/api.go` | Modified | Register `APIResource("/posts", &PostController{})` |
| `http/controllers/controllers_test.go` | Created | 9 tests (unit + integration) |

## Dependencies Added

None — all dependencies already present from Features #07 and #08.

## Test Results

| Package | Tests | Status |
|---|---|---|
| `http/controllers` | 9 | ✅ PASS |
| **Full regression** | **164** | **✅ PASS** |

## Deviations from Plan

| What Changed | Why |
|---|---|
| External test package (`controllers_test`) | Import cycle: `controllers_test` → `routes` → `controllers`; external package breaks the cycle |
