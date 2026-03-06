# Feature #32 — Caching: Test Plan

## Test cases

| TC | Description | Method | Expected |
|----|-------------|--------|----------|
| TC-01 | Set and Get returns value | Set "k" → Get "k" | Returns stored value |
| TC-02 | Get missing key returns empty string | Get "nope" | `""`, nil error |
| TC-03 | Get expired key returns empty string | Set TTL=1ms, sleep, Get | `""`, nil error |
| TC-04 | Delete removes key | Set, Delete, Get | `""` after delete |
| TC-05 | Flush clears all keys | Set 3 keys, Flush, Get all | All return `""` |
| TC-06 | Concurrent access is safe | Parallel Set/Get/Delete | No race (run with -race) |
| TC-07 | NewStore "memory" returns MemoryCache | NewStore("memory", "") | Non-nil Store, nil error |
| TC-08 | NewStore unknown driver returns error | NewStore("redis", "") | nil Store, non-nil error |
| TC-09 | Prefix applied to keys | NewStore(_, "app:"), Set "k" | Internal key is "app:k" |
| TC-10 | Set overwrites existing key | Set "k" twice | Get returns latest value |

## Notes

- TC-03: Use a very short TTL (1ms) and `time.Sleep(5ms)` to test expiry.
- TC-06: Use `t.Parallel()` or a goroutine pool with `sync.WaitGroup`; run tests with `-race`.
- TC-09: Verify via Get("k") returning the value (prefix is internal).
