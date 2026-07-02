package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	configPath := flag.String("config", getDefaultConfigPath(), "Path to config file")
	flag.Parse()

	cfg, err := loadAppConfig(*configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if err := validateAppConfig(cfg); err != nil {
		return err
	}

	application, err := wireApplication(cfg)
	if err != nil {
		return fmt.Errorf("wire application: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg.LogSummary()
	slog.Info("gateway server starting")

	errCh := make(chan error, 1)
	go func() {
		if err := application.Echo.Start(cfg.ListenAddr()); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("server stopped unexpectedly: %w", err)
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown: %w", err)
	}

	slog.Info("server stopped")
	return nil
}
