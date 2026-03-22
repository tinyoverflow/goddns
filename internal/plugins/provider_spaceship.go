package plugins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goddns/internal/plugin"
	"net/http"
)

type spaceshipConfig struct {
	APIKey    string `json:"api_key"    required:"true" doc:"Spaceship API key"`
	APISecret string `json:"api_secret" required:"true" doc:"Spaceship API secret"`
	Domain    string `json:"domain"     required:"true" doc:"The domain name (e.g. example.com)"`
	Subdomain string `json:"subdomain"       required:"true" doc:"The (sub)domain name"`
	TTL       int    `json:"ttl"        default:"3600" doc:"The TTL (Time to live) for the entry"`
}

type spaceshipProvider struct {
	apiKey    string
	apiSecret string
	subdomain string
	domain    string
	ttl       int
}

func init() {
	plugin.RegisterProvider("spaceship", newSpaceshipFromConfig, spaceshipConfig{})
}

func newSpaceshipFromConfig(params map[string]any) (plugin.Provider, error) {
	cfg, err := plugin.Decode[spaceshipConfig](params)
	if err != nil {
		return nil, fmt.Errorf("spaceship: %w", err)
	}

	if cfg.APIKey == "" {
		return nil, fmt.Errorf("spaceship: api_key is required")
	}
	if cfg.APISecret == "" {
		return nil, fmt.Errorf("spaceship: api_secret is required")
	}
	if cfg.Domain == "" {
		return nil, fmt.Errorf("spaceship: domain is required")
	}
	if cfg.Subdomain == "" {
		return nil, fmt.Errorf("spaceship: subdomain is required")
	}
	if cfg.TTL == 0 {
		cfg.TTL = 3600
	}

	return spaceshipProvider{
		apiKey:    cfg.APIKey,
		apiSecret: cfg.APISecret,
		domain:    cfg.Domain,
		subdomain: cfg.Subdomain,
		ttl:       cfg.TTL,
	}, nil
}

func (p spaceshipProvider) SetIPAddress(ip string) error {
	client := http.DefaultClient

	reqData := struct {
		Force bool `json:"force"`
		Items []struct {
			Type    string `json:"type"`
			Name    string `json:"name"`
			TTL     int    `json:"ttl"`
			Address string `json:"address"`
		} `json:"items"`
	}{
		Items: []struct {
			Type    string `json:"type"`
			Name    string `json:"name"`
			TTL     int    `json:"ttl"`
			Address string `json:"address"`
		}{
			{Type: "A", Name: p.subdomain, TTL: p.ttl, Address: ip},
		},
	}

	reqDataBytes, err := json.Marshal(reqData)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://spaceship.dev/api/v1/dns/records/%v", p.domain)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqDataBytes))
	if err != nil {
		return err
	}

	req.Header.Add("X-API-Key", p.apiKey)
	req.Header.Add("X-API-Secret", p.apiSecret)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("expected status 204, got %d", res.StatusCode)
	}

	var data map[string]any
	if json.NewDecoder(res.Body).Decode(&data) != nil {
		return err
	}

	return nil
}

var _ plugin.Provider = spaceshipProvider{}
