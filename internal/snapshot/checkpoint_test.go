package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

func makeCheckpointPorts() []snapshot.Port {
	return []snapshot.Port{
		{Port: 80, Protocol: "tcp", Process: "nginx", PID: 100},
		{Port: 443, Protocol: "tcp", Process: "nginx", PID: 100},
	}
}

func TestCheckpointStore_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewCheckpointStore(dir)

	cp := snapshot.Checkpoint{
		Name:      "before-deploy",
		CreatedAt: time.Now().UTC(),
		Ports:     makeCheckpointPorts(),
		Note:      "pre-release",
	}
	if err := store.Save(cp); err != nil {
		t.Fatalf("Save: %v", err)
	}
	got, err := store.Load("before-deploy")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Name != cp.Name {
		t.Errorf("name: got %q want %q", got.Name, cp.Name)
	}
	if len(got.Ports) != len(cp.Ports) {
		t.Errorf("ports len: got %d want %d", len(got.Ports), len(cp.Ports))
	}
}

func TestCheckpointStore_Load_Missing(t *testing.T) {
	store := snapshot.NewCheckpointStore(t.TempDir())
	_, err := store.Load("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing checkpoint")
	}
}

func TestCheckpointStore_Save_EmptyName(t *testing.T) {
	store := snapshot.NewCheckpointStore(t.TempDir())
	err := store.Save(snapshot.Checkpoint{Name: "", Ports: makeCheckpointPorts()})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestCheckpointStore_List(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewCheckpointStore(dir)

	for _, name := range []string{"alpha", "beta", "gamma"} {
		if err := store.Save(snapshot.Checkpoint{Name: name, Ports: makeCheckpointPorts()}); err != nil {
			t.Fatalf("Save %q: %v", name, err)
		}
	}
	names, err := store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 3 {
		t.Errorf("expected 3 checkpoints, got %d", len(names))
	}
}

func TestCheckpointStore_Delete(t *testing.T) {
	dir := t.TempDir()
	store := snapshot.NewCheckpointStore(dir)

	cp := snapshot.Checkpoint{Name: "to-delete", Ports: makeCheckpointPorts()}
	_ = store.Save(cp)
	if err := store.Delete("to-delete"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "to-delete.json")); !os.IsNotExist(err) {
		t.Error("expected file to be removed")
	}
}

func TestCheckpointStore_List_MissingDir(t *testing.T) {
	store := snapshot.NewCheckpointStore("/nonexistent/path")
	names, err := store.List()
	if err != nil {
		t.Fatalf("List on missing dir should not error: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected empty list, got %v", names)
	}
}
