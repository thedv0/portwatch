package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/daemon"
	"github.com/user/portwatch/internal/rules"
)

func main() {
	configPath := flag.String("config", "configs/portwatch.yaml", "path to config file")
	flag.Parse()

	cfg, err := rules.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("portwatch: failed to load config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("portwatch: invalid config: %v", err)
	}

	notifier := alert.NewNotifier(os.Stderr)
	d := daemon.New(cfg, notifier)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := d.Run(ctx); err != nil && err != context.Canceled {
		log.Fatalf("portwatch: daemon exited with error: %v", err)
	}
}
