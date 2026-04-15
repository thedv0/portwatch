package snapshot

import (
	"sync"
	"time"
)

// CacheOptions configures the in-memory snapshot cache.
type CacheOptions struct {
	MaxEntries int
	TTL        time.Duration
}

// DefaultCacheOptions returns sensible defaults.
func DefaultCacheOptions() CacheOptions {
	return CacheOptions{
		MaxEntries: 64,
		TTL:        5 * time.Minute,
	}
}

type cacheEntry struct {
	snap      Snapshot
	ExpiresAt time.Time
}

// Cache is a thread-safe in-memory store for recent snapshots.
type Cache struct {
	mu      sync.Mutex
	entries map[string]cacheEntry
	opts    CacheOptions
	clock   func() time.Time
}

// NewCache creates a Cache with the given options.
func NewCache(opts CacheOptions) *Cache {
	if opts.MaxEntries <= 0 {
		opts.MaxEntries = DefaultCacheOptions().MaxEntries
	}
	if opts.TTL <= 0 {
		opts.TTL = DefaultCacheOptions().TTL
	}
	return &Cache{
		entries: make(map[string]cacheEntry),
		opts:    opts,
		clock:   time.Now,
	}
}

// Set stores a snapshot under the given key.
func (c *Cache) Set(key string, snap Snapshot) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.entries) >= c.opts.MaxEntries {
		c.evictOldest()
	}
	c.entries[key] = cacheEntry{
		snap:      snap,
		ExpiresAt: c.clock().Add(c.opts.TTL),
	}
}

// Get retrieves a snapshot by key. Returns false if missing or expired.
func (c *Cache) Get(key string) (Snapshot, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entries[key]
	if !ok {
		return Snapshot{}, false
	}
	if c.clock().After(e.ExpiresAt) {
		delete(c.entries, key)
		return Snapshot{}, false
	}
	return e.snap, true
}

// Delete removes a key from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Len returns the number of live (non-expired) entries.
func (c *Cache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.clock()
	count := 0
	for _, e := range c.entries {
		if !now.After(e.ExpiresAt) {
			count++
		}
	}
	return count
}

// evictOldest removes the entry with the earliest expiry. Must be called with lock held.
func (c *Cache) evictOldest() {
	var oldest string
	var oldestTime time.Time
	for k, e := range c.entries {
		if oldest == "" || e.ExpiresAt.Before(oldestTime) {
			oldest = k
			oldestTime = e.ExpiresAt
		}
	}
	if oldest != "" {
		delete(c.entries, oldest)
	}
}
