# 🔍 Review: Logging

> **Feature**: `03` — Logging
> **Status**: ✅ Complete

---

## What Went Well

- Zero new external dependencies — slog is stdlib since Go 1.21
- Config-driven setup worked seamlessly with Feature #02's `config.Env()` helpers
- JSON and text format switching works correctly
- File output with `storage/logs/app.log` creates directories and appends properly
- All 9 unit tests passed (after fixing file handle cleanup)

## What Could Be Improved

- File output test uses `os.Chdir` which is process-wide — same limitation as Feature #02's `TestLoad_NoEnvFile`
- No log rotation built in — acceptable for now, handled by external tools

## Lessons Learned

- slog handlers write immediately — file handles must be closed properly or `t.TempDir()` cleanup fails on Windows
- Adding `Close()` for file handle cleanup was necessary but wasn't in the original architecture doc — discovered during testing
- `slog.SetDefault()` makes the global `slog.Info()` etc. work without passing logger instances — clean pattern for frameworks

## Deviations from Plan

1. **Added `Close()` function** — Not in original architecture doc. Needed to properly close file handles when `LOG_OUTPUT=file`. Without it, Windows tests fail because TempDir cleanup can't delete open files. This is a good addition for production use too (graceful shutdown).
