package snapshot

import (
	"testing"
	"time"
)

func fixedCacheTime(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestCache_SetAndGet(t *testing.T) {
	c := NewCache(DefaultCacheOptions())
	snap := Snapshot{Ports: []Port{{Port: 80, Protocol: "tcp"}}}
	c.Set("k1", snap)
	got, ok := c.Get("k1")
	if !ok {
		t.Fatal("expected hit")
	}
	if len(got.Ports) != 1 {
		t.Fatalf("expected 1 port, got %d", len(got.Ports))
	}
}

func TestCache_Get_Missing(t *testing.T) {
	c := NewCache(DefaultCacheOptions())
	_, ok := c.Get("nonexistent")
	if ok {
		t.Fatal("expected miss")
	}
}

func TestCache_Get_Expired(t *testing.T) {
	now := time.Now()
	c := NewCache(CacheOptions{MaxEntries: 10, TTL: time.Minute})
	c.clock = fixedCacheTime(now)
	c.Set("k1", Snapshot{})
	// advance clock past TTL
	c.clock = fixedCacheTime(now.Add(2 * time.Minute))
	_, ok := c.Get("k1")
	if ok {
		t.Fatal("expected expired entry to be a miss")
	}
}

func TestCache_Delete(t *testing.T) {
	c := NewCache(DefaultCacheOptions())
	c.Set("k1", Snapshot{})
	c.Delete("k1")
	_, ok := c.Get("k1")
	if ok {
		t.Fatal("expected miss after delete")
	}
}

func TestCache_Len_ExcludesExpired(t *testing.T) {
	now := time.Now()
	c := NewCache(CacheOptions{MaxEntries: 10, TTL: time.Minute})
	c.clock = fixedCacheTime(now)
	c.Set("a", Snapshot{})
	c.Set("b", Snapshot{})
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
	c.clock = fixedCacheTime(now.Add(2 * time.Minute))
	if c.Len() != 0 {
		t.Fatalf("expected 0 after expiry, got %d", c.Len())
	}
}

func TestCache_MaxEntries_EvictsOldest(t *testing.T) {
	now := time.Now()
	c := NewCache(CacheOptions{MaxEntries: 2, TTL: time.Hour})
	c.clock = fixedCacheTime(now)
	c.Set("first", Snapshot{})
	c.clock = fixedCacheTime(now.Add(time.Second))
	c.Set("second", Snapshot{})
	c.clock = fixedCacheTime(now.Add(2 * time.Second))
	c.Set("third", Snapshot{})
	// total stored entries should not exceed MaxEntries
	count := 0
	for _, k := range []string{"first", "second", "third"} {
		if _, ok := c.Get(k); ok {
			count++
		}
	}
	if count > 2 {
		t.Fatalf("expected at most 2 live entries, got %d", count)
	}
}

func TestDefaultCacheOptions_Values(t *testing.T) {
	opts := DefaultCacheOptions()
	if opts.MaxEntries <= 0 {
		t.Error("MaxEntries should be positive")
	}
	if opts.TTL <= 0 {
		t.Error("TTL should be positive")
	}
}
