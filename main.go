package main

import (
	"context"
	"goddns/internal/config"
	"goddns/internal/provider"
	"goddns/internal/retriever"
	"log/slog"
	"os"
	"time"
)

func main() {
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
	logger := slog.New(logHandler)

	if err := run(context.Background(), logger); err != nil {
		slog.Error("error running ddns updater", "error", err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *slog.Logger) error {
	logger.Info("starting goddns")

	cfg := config.Load()
	logger.Info(
		"loaded config",
		"interval", cfg.Interval,
		"zone", cfg.HCloudZone,
	)

	ret := retriever.IfConfigRetriever{}
	prv := provider.HetznerCloudProvider{
		Token: cfg.HCloudToken,
		Zone:  cfg.HCloudZone,
	}

	tickDuration := time.Duration(cfg.Interval) * time.Second
	ticker := time.NewTicker(tickDuration)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			err := fetchAndUpdate(ctx, logger, ret, prv)
			if err != nil {
				slog.Error("error fetching and updating", "error", err.Error())
			}
		}
	}
}

func fetchAndUpdate(_ context.Context, logger *slog.Logger, ret retriever.Retriever, prv provider.Provider) error {
	logger.Debug("fetching ip address")
	ip, err := ret.GetIPAddress()
	if err != nil {
		return err
	}

	logger.Debug("updating ip address", "ip", ip)
	if err := prv.SetIPAddress(ip); err != nil {
		return err
	}

	logger.Info("updated ip address", "ip", ip)
	return nil
}
