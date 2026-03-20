package main

import (
	"context"
	"fmt"
	"goddns/internal/config"
	"goddns/internal/plugin"
	"log/slog"
	"os"
	"sync"
	"time"
)

func main() {
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(logHandler)

	if err := run(context.Background(), logger); err != nil {
		slog.Error("error running ddns updater", "error", err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *slog.Logger) error {
	logger.Info("starting goddns")

	configPath := "config.toml"
	if p := os.Getenv("GODDNS_CONFIG"); p != "" {
		configPath = p
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	logger.Info(
		"loaded config",
		"default_interval", cfg.Interval,
		"retrievers", len(cfg.Retrievers),
		"providers", len(cfg.Providers),
		"instances", len(cfg.Instances),
	)

	var wg sync.WaitGroup
	for name, inst := range cfg.Instances {
		interval := cfg.Interval
		if inst.Interval != nil {
			interval = *inst.Interval
		}

		ret, err := plugin.BuildRetriever(cfg.Retrievers, inst.Retriever)
		if err != nil {
			return fmt.Errorf("instance %q: %w", name, err)
		}

		providers, err := plugin.BuildProviders(cfg.Providers, inst.Providers)
		if err != nil {
			return fmt.Errorf("instance %q: %w", name, err)
		}

		logger.Info(
			"initializing instance",
			"instance", name,
			"interval", interval,
			"retriever", inst.Retriever["name"],
			"providers", len(inst.Providers),
		)

		dur, err := time.ParseDuration(interval)
		if err != nil {
			return fmt.Errorf("invalid duration format %q", interval)
		}

		wg.Add(1)
		go func(name string, ret plugin.Retriever, providers []plugin.Provider, dur time.Duration) {
			defer wg.Done()
			runInstance(ctx, logger.With("instance", name), dur, ret, providers)
		}(name, ret, providers, dur)
	}

	wg.Wait()
	return nil
}

func runInstance(ctx context.Context, logger *slog.Logger, dur time.Duration, ret plugin.Retriever, providers []plugin.Provider) {
	var latestIP string
	ticker := time.NewTicker(dur)
	defer ticker.Stop()

	for {
		if err := fetchAndUpdate(ctx, logger, ret, providers, &latestIP); err != nil {
			logger.Error("fetching and updating ip address", "error", err)
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func fetchAndUpdate(_ context.Context, logger *slog.Logger, ret plugin.Retriever, providers []plugin.Provider, latestIP *string) error {
	logger.Debug("fetching ip address")
	ip, err := ret.GetIPAddress()
	if err != nil {
		return err
	}

	logger.Debug("received ip address", "ip", ip)

	if ip == *latestIP {
		logger.Debug("ip address did not change")
		return nil
	}

	logger.Info("detected new ip address", "new_ip", ip, "old_ip", *latestIP)

	for _, prv := range providers {
		if err := prv.SetIPAddress(ip); err != nil {
			return err
		}
	}

	*latestIP = ip
	logger.Info("updated ip address", "ip", ip)
	return nil
}
