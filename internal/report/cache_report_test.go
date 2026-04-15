package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func makeCacheReport() CacheReport {
	return BuildCacheReport(12, 64, 5*time.Minute, 300, 45)
}

func TestBuildCacheReport_Fields(t *testing.T) {
	r := makeCacheReport()
	if r.LiveEntries != 12 {
		t.Errorf("LiveEntries: want 12, got %d", r.LiveEntries)
	}
	if r.MaxEntries != 64 {
		t.Errorf("MaxEntries: want 64, got %d", r.MaxEntries)
	}
	if r.TTLSeconds != 300 {
		t.Errorf("TTLSeconds: want 300, got %f", r.TTLSeconds)
	}
	if r.Hits != 300 {
		t.Errorf("Hits: want 300, got %d", r.Hits)
	}
	if r.Misses != 45 {
		t.Errorf("Misses: want 45, got %d", r.Misses)
	}
}

func TestBuildCacheReport_TimestampSet(t *testing.T) {
	before := time.Now()
	r := makeCacheReport()
	after := time.Now()
	if r.Timestamp.Before(before) || r.Timestamp.After(after) {
		t.Error("Timestamp not within expected range")
	}
}

func TestWriteCacheText_ContainsHeaders(t *testing.T) {
	var buf bytes.Buffer
	r := makeCacheReport()
	if err := WriteCacheText(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"Cache Report", "Live Entries", "TTL", "Hits", "Misses"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestWriteCacheJSON_ValidJSON(t *testing.T) {
	var buf bytes.Buffer
	r := makeCacheReport()
	if err := WriteCacheJSON(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out CacheReport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.LiveEntries != r.LiveEntries {
		t.Errorf("LiveEntries mismatch: want %d, got %d", r.LiveEntries, out.LiveEntries)
	}
}

func TestWriteCacheJSON_HitsMissesPreserved(t *testing.T) {
	var buf bytes.Buffer
	r := BuildCacheReport(5, 32, time.Minute, 99, 1)
	_ = WriteCacheJSON(&buf, r)
	var out CacheReport
	_ = json.Unmarshal(buf.Bytes(), &out)
	if out.Hits != 99 || out.Misses != 1 {
		t.Errorf("hits/misses not preserved: %d/%d", out.Hits, out.Misses)
	}
}
