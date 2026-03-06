# Feature #32 — Caching: Tasks

## Prerequisites

- [x] Configuration system shipped (#02)
- [x] Service container shipped (#05)
- [x] `core/cache/` directory exists (stub with `.gitkeep`)

## Implementation tasks

| # | Task | File(s) | Status |
|---|------|---------|--------|
| 1 | Define `Store` interface | `core/cache/cache.go` | ⬜ |
| 2 | Implement `MemoryCache` (Get/Set/Delete/Flush with TTL) | `core/cache/cache.go` | ⬜ |
| 3 | Implement `NewStore()` factory reading driver + prefix | `core/cache/cache.go` | ⬜ |
| 4 | Remove `.gitkeep` | `core/cache/.gitkeep` | ⬜ |
| 5 | Write tests | `core/cache/cache_test.go` | ⬜ |
| 6 | Full regression + `go vet` | — | ⬜ |
| 7 | Commit, merge, review doc, roadmap update | — | ⬜ |

## Acceptance criteria

- `Store` interface defines Get/Set/Delete/Flush.
- `MemoryCache` is thread-safe (sync.RWMutex).
- `Get()` returns `""` for expired entries (lazy cleanup).
- `Set()` stores with absolute expiry time.
- `Flush()` clears all entries.
- `NewStore("memory", prefix)` wraps keys with prefix.
- `NewStore("unknown", ...)` returns an error.
- All existing tests pass (regression).
