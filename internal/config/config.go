package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Interval   string                    `toml:"interval"`
	Retrievers map[string]map[string]any `toml:"retriever"`
	Providers  map[string]map[string]any `toml:"provider"`
	Instances  map[string]InstanceConfig `toml:"instance"`
}

type InstanceConfig struct {
	Interval  *string          `toml:"interval"`
	Retriever map[string]any   `toml:"retriever"`
	Providers []map[string]any `toml:"provider"`
}

func Load(path string) (Config, error) {
	cfg := Config{
		Interval: "5m",
	}

	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return Config{}, fmt.Errorf("decoding config file: %w", err)
	}

	return cfg, nil
}
