# Feature #32 — Caching: Changelog

## [Unreleased]

### Added
- `core/cache/cache.go` — `Store` interface, `MemoryCache`, `NewStore()` factory.
- 10 test cases (TC-01 to TC-10) for CRUD, expiry, concurrency, prefix, and factory.

### Removed
- `core/cache/.gitkeep` — replaced by real implementation.

### Deviation log
| # | Blueprint | Ours | Reason |
|---|-----------|------|--------|
| 1 | `RedisCache` with `go-redis/v9` | Deferred | Avoids heavy external dependency + running Redis requirement |
| 2 | Prefix in each driver | Prefix wrapper in `NewStore()` factory | Keeps driver implementations simpler; prefix is a cross-cutting concern |
