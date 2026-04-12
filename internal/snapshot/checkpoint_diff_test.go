package snapshot_test

import (
	"testing"

	"github.com/user/portwatch/internal/snapshot"
)

func TestCompareToCheckpoint_NoChange(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewCheckpointStore(dir)
	ports := makeCheckpointPorts()
	_ = store.Save(snapshot.Checkpoint{Name: "snap", Ports: ports})

	res, err := snapshot.CompareToCheckpoint(store, "snap", ports)
	if err != nil {
		t.Fatalf("CompareToCheckpoint: %v", err)
	}
	if len(res.Added) != 0 || len(res.Removed) != 0 {
		t.Errorf("expected no diff, got added=%d removed=%d", len(res.Added), len(res.Removed))
	}
	if res.Unchanged != 2 {
		t.Errorf("unchanged: got %d want 2", res.Unchanged)
	}
}

func TestCompareToCheckpoint_Added(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewCheckpointStore(dir)
	_ = store.Save(snapshot.Checkpoint{Name: "snap", Ports: makeCheckpointPorts()})

	current := append(makeCheckpointPorts(), snapshot.Port{Port: 8080, Protocol: "tcp", Process: "app", PID: 200})
	res, err := snapshot.CompareToCheckpoint(store, "snap", current)
	if err != nil {
		t.Fatalf("CompareToCheckpoint: %v", err)
	}
	if len(res.Added) != 1 {
		t.Errorf("added: got %d want 1", len(res.Added))
	}
	if res.Added[0].Port != 8080 {
		t.Errorf("added port: got %d want 8080", res.Added[0].Port)
	}
}

func TestCompareToCheckpoint_Removed(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewCheckpointStore(dir)
	_ = store.Save(snapshot.Checkpoint{Name: "snap", Ports: makeCheckpointPorts()})

	current := []snapshot.Port{{Port: 80, Protocol: "tcp", Process: "nginx", PID: 100}}
	res, err := snapshot.CompareToCheckpoint(store, "snap", current)
	if err != nil {
		t.Fatalf("CompareToCheckpoint: %v", err)
	}
	if len(res.Removed) != 1 {
		t.Errorf("removed: got %d want 1", len(res.Removed))
	}
	if res.Removed[0].Port != 443 {
		t.Errorf("removed port: got %d want 443", res.Removed[0].Port)
	}
}

func TestCompareToCheckpoint_MissingCheckpoint(t *testing.T) {
	store := snapshot.NewCheckpointStore(t.TempDir())
	_, err := snapshot.CompareToCheckpoint(store, "ghost", makeCheckpointPorts())
	if err == nil {
		t.Fatal("expected error for missing checkpoint")
	}
}

func TestCompareToCheckpoint_Name(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewCheckpointStore(dir)
	ports := makeCheckpointPorts()
	_ = store.Save(snapshot.Checkpoint{Name: "mycheck", Ports: ports})

	res, err := snapshot.CompareToCheckpoint(store, "mycheck", ports)
	if err != nil {
		t.Fatalf("CompareToCheckpoint: %v", err)
	}
	if res.CheckpointName != "mycheck" {
		t.Errorf("CheckpointName: got %q want %q", res.CheckpointName, "mycheck")
	}
}
