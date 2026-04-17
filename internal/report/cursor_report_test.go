package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

func makeCursorResult(n int, hasMore bool) snapshot.CursorResult {
	base := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	var page []snapshot.Snapshot
	for i := 0; i < n; i++ {
		page = append(page, snapshot.Snapshot{
			Timestamp: base.Add(time.Duration(i) * time.Minute),
			Ports:     []snapshot.Port{{Port: 80 + i, Protocol: "tcp"}},
		})
	}
	var next time.Time
	if n > 0 {
		next = page[n-1].Timestamp
	}
	return snapshot.CursorResult{Page: page, HasMore: hasMore, NextAfter: next}
}

func TestBuildCursorReport_Returned(t *testing.T) {
	r := BuildCursorReport(makeCursorResult(3, true))
	if r.Returned != 3 {
		t.Fatalf("expected Returned=3, got %d", r.Returned)
	}
}

func TestBuildCursorReport_HasMore(t *testing.T) {
	r := BuildCursorReport(makeCursorResult(2, true))
	if !r.HasMore {
		t.Fatal("expected HasMore=true")
	}
}

func TestBuildCursorReport_TimestampSet(t *testing.T) {
	r := BuildCursorReport(makeCursorResult(1, false))
	if r.Timestamp.IsZero() {
		t.Fatal("expected non-zero timestamp")
	}
}

func TestWriteCursorText_ContainsHeaders(t *testing.T) {
	r := BuildCursorReport(makeCursorResult(2, false))
	var buf bytes.Buffer
	if err := WriteCursorText(&buf, r); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{"Cursor Report", "Returned", "Has More"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in output", want)
		}
	}
}

func TestWriteCursorJSON_ValidJSON(t *testing.T) {
	r := BuildCursorReport(makeCursorResult(2, true))
	var buf bytes.Buffer
	if err := WriteCursorJSON(&buf, r); err != nil {
		t.Fatal(err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := out["has_more"]; !ok {
		t.Error("missing has_more field")
	}
}
