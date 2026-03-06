# Feature #32 — Caching: Discussion

## What problem does this solve?

Applications frequently re-compute or re-fetch expensive data (DB queries, API calls, template fragments). A caching layer stores results with a TTL so subsequent requests are served from fast in-process memory instead of repeating costly operations.

## Why now?

Configuration (#02), service container (#05), and service providers (#06) are shipped. The cache store can be registered in the container and consumed by any service or controller.

## What does the blueprint specify?

- `Store` interface with `Get(key) (string, error)`, `Set(key, value, ttl)`, `Delete(key)`, `Flush()`.
- `MemoryCache` — in-process `sync.RWMutex` + `map[string]memCacheEntry` with TTL-based expiry.
- `RedisCache` — wraps `github.com/redis/go-redis/v9`.
- Env vars: `CACHE_DRIVER`, `CACHE_PREFIX`, `CACHE_TTL`.
- `CacheProvider` that registers a `"cache"` singleton in the container using `cache.NewMemoryCache()`.

## Design decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Interface | `Store` with Get/Set/Delete/Flush | Blueprint-specified; clean contract for swappable backends |
| MemoryCache | Implement (stdlib only) | No external deps; covers dev/test/single-instance production |
| RedisCache | Defer | Requires `go-redis/v9` + running Redis; add when actually needed |
| Factory function | `NewStore()` reads `CACHE_DRIVER` env | Consistent with storage (#28) and other driver-based subsystems |
| Key prefixing | Applied in `NewStore()` wrapper or per-driver | Keeps keys namespaced; reads `CACHE_PREFIX` env |
| Default TTL | `CACHE_TTL` env (seconds), default 3600 | Blueprint-specified |
| Expired entry cleanup | Lazy on `Get()` | Blueprint approach; avoids background goroutine complexity |
| Package location | `core/cache/` | Existing stub directory |

## What is out of scope?

- Redis driver (deferred — heavy dependency).
- File-based cache driver (roadmap mentions it, but blueprint only specifies Redis + memory).
- Cache tags / invalidation groups (app-level concern).
- Distributed cache locking (app-level concern).
