package daemon

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/rules"
)

func defaultConfig() *rules.Config {
	cfg := rules.DefaultConfig()
	cfg.Interval = "50ms"
	return cfg
}

func TestNew_InitialState(t *testing.T) {
	cfg := defaultConfig()
	d := New(cfg, alert.NewNotifier(nil))

	if d.cfg == nil {
		t.Fatal("expected cfg to be set")
	}
	if d.scanner == nil {
		t.Fatal("expected scanner to be set")
	}
	if d.prev == nil {
		t.Fatal("expected prev map to be initialised")
	}
}

func TestRun_CancelStops(t *testing.T) {
	cfg := defaultConfig()
	var buf strings.Builder
	n := alert.NewNotifier(&buf)
	d := New(cfg, n)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	err := d.Run(ctx)
	if err != context.DeadlineExceeded {
		// context.Canceled is also acceptable depending on who cancels first.
		if err != context.Canceled {
			t.Fatalf("expected context error, got %v", err)
		}
	}
}

func TestRun_InvalidInterval(t *testing.T) {
	cfg := defaultConfig()
	cfg.Interval = "not-a-duration"
	d := New(cfg, alert.NewNotifier(nil))

	err := d.Run(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid interval")
	}
}
