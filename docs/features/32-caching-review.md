# Feature #32 — Caching: Review

## Summary

Implemented a driver-based caching subsystem with a `Store` interface, in-memory backend, key-prefix wrapper, and factory function.

## Files changed

| File | Action | Purpose |
|------|--------|---------|
| `core/cache/cache.go` | Created | `Store` interface, `MemoryCache`, `prefixStore`, `NewStore()` |
| `core/cache/cache_test.go` | Created | 10 test cases |
| `core/cache/.gitkeep` | Deleted | Replaced by real implementation |

## Blueprint compliance

| Blueprint item | Status | Notes |
|----------------|--------|-------|
| `Store` interface (Get/Set/Delete/Flush) | ✅ | Exact match |
| `MemoryCache` with sync.RWMutex | ✅ | Improved: lazy delete on Get() |
| `RedisCache` with go-redis/v9 | ⏳ Deferred | Avoids heavy dependency |
| `CACHE_DRIVER` / `CACHE_PREFIX` / `CACHE_TTL` env vars | ✅ | Factory reads driver + prefix |
| `NewMemoryCache()` constructor | ✅ | Exact match |

## Deviations

| # | Blueprint | Ours | Reason |
|---|-----------|------|--------|
| 1 | `RedisCache` with `go-redis/v9` | Deferred | Avoids external dependency + running Redis requirement |
| 2 | Prefix baked into each driver | `prefixStore` wrapper in factory | Cleaner separation; prefix is a cross-cutting concern |
| 3 | Get() holds RLock, doesn't clean up | Get() upgrades to write lock for lazy delete | Prevents unbounded memory growth from expired entries |

## Test results

| TC | Description | Result |
|----|-------------|--------|
| TC-01 | Set and Get returns value | ✅ PASS |
| TC-02 | Get missing key returns empty string | ✅ PASS |
| TC-03 | Get expired key returns empty string + lazy delete | ✅ PASS |
| TC-04 | Delete removes key | ✅ PASS |
| TC-05 | Flush clears all keys | ✅ PASS |
| TC-06 | Concurrent access is safe | ✅ PASS |
| TC-07 | NewStore "memory" returns valid Store | ✅ PASS |
| TC-08 | NewStore unknown driver returns error | ✅ PASS |
| TC-09 | Prefix applied to keys | ✅ PASS |
| TC-10 | Set overwrites existing key | ✅ PASS |

## Regression

- All 27 packages pass (`go test ./...`)
- `go vet ./...` clean
- No new dependencies added (stdlib only)
