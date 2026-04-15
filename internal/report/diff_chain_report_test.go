package report

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

func makeChainEntries() []snapshot.ChainEntry {
	now := time.Now().UTC()
	return []snapshot.ChainEntry{
		{
			Index: 0,
			Snap:  snapshot.Snapshot{Timestamp: now},
			Diff:  snapshot.DiffResult{},
		},
		{
			Index: 1,
			Snap:  snapshot.Snapshot{Timestamp: now.Add(time.Minute)},
			Diff: snapshot.DiffResult{
				Added:   []snapshot.Port{{Protocol: "tcp", Port: 443}},
				Removed: []snapshot.Port{},
			},
		},
	}
}

func TestBuildDiffChainReport_EntryCount(t *testing.T) {
	r := BuildDiffChainReport(makeChainEntries())
	if r.EntryCount != 2 {
		t.Errorf("expected 2, got %d", r.EntryCount)
	}
}

func TestBuildDiffChainReport_TotalAdded(t *testing.T) {
	r := BuildDiffChainReport(makeChainEntries())
	if r.TotalAdded != 1 {
		t.Errorf("expected TotalAdded=1, got %d", r.TotalAdded)
	}
}

func TestBuildDiffChainReport_TimestampSet(t *testing.T) {
	r := BuildDiffChainReport(makeChainEntries())
	if r.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestBuildDiffChainReport_EmptyInput(t *testing.T) {
	r := BuildDiffChainReport(nil)
	if r.EntryCount != 0 || r.TotalAdded != 0 || r.TotalRemoved != 0 {
		t.Error("expected all-zero report for empty input")
	}
}

func TestWriteDiffChainText_ContainsHeaders(t *testing.T) {
	r := BuildDiffChainReport(makeChainEntries())
	var buf bytes.Buffer
	if err := WriteDiffChainText(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Diff Chain Report") {
		t.Error("expected 'Diff Chain Report' header")
	}
	if !strings.Contains(out, "Total Added") {
		t.Error("expected 'Total Added' in output")
	}
}

func TestWriteDiffChainJSON_ValidJSON(t *testing.T) {
	r := BuildDiffChainReport(makeChainEntries())
	var buf bytes.Buffer
	if err := WriteDiffChainJSON(&buf, r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out DiffChainReport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out.EntryCount != 2 {
		t.Errorf("expected EntryCount=2 in JSON, got %d", out.EntryCount)
	}
}
