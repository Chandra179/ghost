package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache implements Cache backed by a single Redis instance.
// Use for high-throughput key-value cache with TTL expiry.
// Not suitable for: multi-datacenter setups (no replication awareness),
// strongly-consistent reads (Redis async replication), or values >512MB.
type RedisCache struct {
	client *redis.Client
	addr   string
}

// NewRedisCache connects to a Redis instance at addr (host:port).
// Returns a Cache interface.
// No auth or TLS — use redis:// URL for production or extend.
// On connect failure the first op will error, not this call.
func NewRedisCache(addr string) Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisCache{client: rdb, addr: addr}
}

// Addr returns the Redis server address this cache connects to.
func (r *RedisCache) Addr() string { return r.addr }

// Set stores value under key with ttl. Overwrites existing key.
// Use for: fresh writes, cache-warming, TTL refreshes.
// Limitation: no atomic existence check — use SetNX for that.
func (r *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Get retrieves value for key. Returns ErrCacheMiss if missing.
// Use for: cache reads, hot-path lookups.
// Limitation: no bulk/pipe — call N times for N keys.
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrCacheMiss
	}
	return val, err
}

// Del removes one or more keys. No-op if key missing.
// Use for: cache invalidation on write, manual eviction.
// Limitation: single-node — no cross-shard delete in cluster.
func (r *RedisCache) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// SetNX atomically sets key only if it does not exist. Returns true on success.
// Use for: distributed locks, deduplication, first-writer-wins.
// Limitation: not a reentrant lock — no lock-ID or TTL extension built in.
func (r *RedisCache) SetNX(ctx context.Context, key, value string, ttl time.Duration) (bool, error) {
	err := r.client.SetArgs(ctx, key, value, redis.SetArgs{
		Mode: "NX",
		TTL:  ttl,
	}).Err()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// Ping checks Redis connection liveness. Returns nil if reachable.
// Use for: startup health-check, readiness probe.
// Limitation: checks TCP only, not storage — a full Redis may still fail on ops.
func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Close shuts down the underlying Redis connection. Call once on app shutdown.
// Further ops will error. Idempotent.
func (r *RedisCache) Close() error {
	return r.client.Close()
}
