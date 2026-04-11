package snapshot

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

type mockScanner struct {
	mu      sync.Mutex
	calls   int
	results [][]scanner.PortState
	err     error
}

func (m *mockScanner) Scan() ([]scanner.PortState, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.err != nil {
		return nil, m.err
	}
	if m.calls >= len(m.results) {
		return m.results[len(m.results)-1], nil
	}
	res := m.results[m.calls]
	m.calls++
	return res, nil
}

func TestWatcher_CallsOnChange(t *testing.T) {
	p1 := scanner.PortState{Port: 80, Protocol: "tcp", Process: "nginx", PID: 1}
	p2 := scanner.PortState{Port: 443, Protocol: "tcp", Process: "nginx", PID: 2}

	ms := &mockScanner{
		results: [][]scanner.PortState{{p1}, {p1, p2}},
	}

	var got DiffResult
	var mu sync.Mutex
	changed := make(chan struct{}, 1)

	w := NewWatcher(ms, WatchOptions{
		Interval: 10 * time.Millisecond,
		OnChange: func(d DiffResult) {
			mu.Lock()
			g d
			mu.Unlock()
			changed <- struct{}{}
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	go w.Run(ctx)

	select {
	case <-changed:
	case <-ctx.Done():
		t.Fatal("timed out waiting for OnChange")
	}

	mu.Lock()
	defer mu.Unlock()
	if len(got.Added) != 1 || got.Added[0].Port != 443 {
		t.Errorf("expected port 443 added, got %+v", got.Added)
	}
}

func TestWatcher_StopsOnCancel(t *testing.T) {
	ms := &mockScanner{
		results: [][]scanner.PortState{{}},
	}
	w := NewWatcher(ms, WatchOptions{
		Interval: 10 * time.Millisecond,
	})
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		w.Run(ctx)
		close(done)
	}()
	cancel()
	select {
	case <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("watcher did not stop after cancel")
	}
}

func TestWatcher_DefaultInterval(t *testing.T) {
	ms := &mockScanner{results: [][]scanner.PortState{{}}}
	w := NewWatcher(ms, WatchOptions{})
	if w.opts.Interval != 30*time.Second {
		t.Errorf("expected default 30s, got %v", w.opts.Interval)
	}
}
