package plugins

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"goddns/internal/plugin"
	"net/http"
)

// UnifiConfig holds the configuration parameters for the unifi retriever.
type UnifiConfig struct {
	BaseURL   string `json:"base_url"   required:"true"   doc:"Unifi controller base URL (e.g. https://192.168.1.1)"`
	APIToken  string `json:"api_token"  required:"true"   doc:"API authentication token"`
	SiteID    string `json:"site_id"    default:"default" doc:"Site ID"`
	VerifyTLS bool   `json:"verify_tls" default:"false"   doc:"Enable TLS certificate verification"`
}

type unifiRetriever struct {
	baseURL  string
	apiToken string
	siteID   string
	client   *http.Client
}

func init() {
	plugin.RegisterRetriever("unifi", newUnifiFromConfig, UnifiConfig{})
}

func newUnifiFromConfig(params map[string]any) (plugin.Retriever, error) {
	cfg, err := plugin.Decode[UnifiConfig](params)
	if err != nil {
		return nil, fmt.Errorf("unifi retriever: %w", err)
	}

	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("unifi retriever: base_url is required")
	}
	if cfg.APIToken == "" {
		return nil, fmt.Errorf("unifi retriever: api_token is required")
	}
	if cfg.SiteID == "" {
		cfg.SiteID = "default"
	}

	tlsConfig := &tls.Config{}
	if !cfg.VerifyTLS {
		tlsConfig.InsecureSkipVerify = true
	}

	return unifiRetriever{
		baseURL:  cfg.BaseURL,
		apiToken: cfg.APIToken,
		siteID:   cfg.SiteID,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		},
	}, nil
}

func (r unifiRetriever) GetIPAddress() (string, error) {
	url := fmt.Sprintf("%s/proxy/network/api/s/%s/stat/health", r.baseURL, r.siteID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("X-API-KEY", r.apiToken)
	req.Header.Set("Accept", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("health request returned status %d", resp.StatusCode)
	}

	var data struct {
		Data []struct {
			WanIP string `json:"wan_ip"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	for _, entry := range data.Data {
		if entry.WanIP != "" {
			return entry.WanIP, nil
		}
	}

	return "", fmt.Errorf("wan_ip not found in health response")
}

var _ plugin.Retriever = unifiRetriever{}
