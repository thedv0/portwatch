package snapshot

import (
	"testing"
	"time"
)

func makeCSSnap(t time.Time, port int) Snapshot {
	return Snapshot{
		Timestamp: t,
		Ports:     []Port{{Port: port, Protocol: "tcp"}},
	}
}

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestCursor_InvalidPageSize(t *testing.T) {
	_, err := Cursor(nil, CursorOptions{PageSize: 0})
	if err == nil {
		t.Fatal("expected error for zero PageSize")
	}
}

func TestCursor_EmptyInput(t *testing.T) {
	res, err := Cursor(nil, DefaultCursorOptions())
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Page) != 0 || res.HasMore {
		t.Fatal("expected empty result")
	}
}

func TestCursor_ReturnsPageSizeItems(t *testing.T) {
	snaps := []Snapshot{
		makeCSSnap(epoch.Add(1*time.Minute), 80),
		makeCSSnap(epoch.Add(2*time.Minute), 443),
		makeCSSnap(epoch.Add(3*time.Minute), 8080),
	}
	res, err := Cursor(snaps, CursorOptions{PageSize: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Page) != 2 {
		t.Fatalf("expected 2 items, got %d", len(res.Page))
	}
	if !res.HasMore {
		t.Fatal("expected HasMore=true")
	}
}

func TestCursor_AfterFilters(t *testing.T) {
	snaps := []Snapshot{
		makeCSSnap(epoch.Add(1*time.Minute), 80),
		makeCSSnap(epoch.Add(2*time.Minute), 443),
		makeCSSnap(epoch.Add(3*time.Minute), 8080),
	}
	res, err := Cursor(snaps, CursorOptions{PageSize: 10, After: epoch.Add(1*time.Minute)})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Page) != 2 {
		t.Fatalf("expected 2 items after cursor, got %d", len(res.Page))
	}
}

func TestCursor_NextAfterIsLastTimestamp(t *testing.T) {
	snaps := []Snapshot{
		makeCSSnap(epoch.Add(1*time.Minute), 80),
		makeCSSnap(epoch.Add(2*time.Minute), 443),
	}
	res, _ := Cursor(snaps, CursorOptions{PageSize: 10})
	if !res.NextAfter.Equal(epoch.Add(2 * time.Minute)) {
		t.Fatalf("unexpected NextAfter: %v", res.NextAfter)
	}
}

func TestCursor_OrderedAscending(t *testing.T) {
	snaps := []Snapshot{
		makeCSSnap(epoch.Add(3*time.Minute), 8080),
		makeCSSnap(epoch.Add(1*time.Minute), 80),
		makeCSSnap(epoch.Add(2*time.Minute), 443),
	}
	res, _ := Cursor(snaps, CursorOptions{PageSize: 10})
	for i := 1; i < len(res.Page); i++ {
		if res.Page[i].Timestamp.Before(res.Page[i-1].Timestamp) {
			t.Fatal("page not ordered ascending")
		}
	}
}
