# ADR 001: Cache-Aside Pattern via `GetOrSet`

## Status

Accepted

## Context

Consumers need a standard way to load data on demand into Redis cache. Common pattern: check cache, miss -> query source -> store -> return. Every consumer reimplements this with subtle differences in error handling, TTL strategy, and thundering herd protection.

## Decision

Add `GetOrSet` as a package-level function in `cache/` (not on `Cache` interface). Signature:

```go
func GetOrSet(ctx context.Context, c Cache, key string, ttl time.Duration,
    fn func(context.Context) (string, error), opts ...GetOrSetOptions) (string, error)
```

### Approach

1. **Fast path**: `c.Get(key)`. Hit -> return.
2. **Thundering herd**: `c.SetNX(key:lock)` — first caller acquires lock, computes, stores, releases. Others retry cache briefly, then fall through to `fn`.
3. **Lock TTL**: Configurable via `GetOrSetOptions.LockTTL` (default 10s). Auto-expires if holder crashes.
4. **Nil caching**: `GetOrSetOptions.CacheNil` — stores sentinel for empty fn results to prevent repeat misses.

### Why standalone function, not interface method

`GetOrSet` is a composite operation built on `Get`/`Set`/`SetNX`/`Del`. Every `Cache` impl gets it for free. Keeping it off the interface avoids forcing all implementations to provide it.

## Configuration

| Option | Default | Purpose |
|--------|---------|---------|
| `LockTTL` | 10s | Thundering herd lock TTL |
| `CacheNil` | false | Cache empty fn results |

## Consequences

- Consumers stop reimplementing cache-aside.
- Thundering herd mitigated without external lock service.
- Lock TTL must exceed fn worst-case execution; too-short TTL can cause redundant fn calls.
- `Get` now returns `ErrCacheMiss` instead of `redis.Nil` — test assertions updated accordingly.

## References

- AWS ElastiCache: [Lazy Loading Strategy](https://docs.aws.amazon.com/AmazonElastiCache/latest/red-ug/Strategies.html)
- Azure: [Cache-Aside Pattern](https://learn.microsoft.com/en-us/azure/architecture/patterns/cache-aside)
