package config

import (
	"os"
	"strconv"
)

type Config struct {
	Interval    int
	HCloudToken string
	HCloudZone  string
}

func Load() Config {
	return Config{
		Interval:    getEnvAsInt("GODDNS_INTERVAL", 300),
		HCloudToken: getEnvAsString("GODDNS_HCLOUD_TOKEN", ""),
		HCloudZone:  getEnvAsString("GODDNS_HCLOUD_ZONE", ""),
	}
}

func getEnvAsString(key string, defaultValue string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	return val
}

func getEnvAsInt(key string, fallback int) int {
	val, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}

	valInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return valInt
}
