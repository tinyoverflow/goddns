package plugins

import (
	"encoding/json"
	"goddns/internal/plugin"
	"net/http"
)

// IfConfigCoConfig holds the configuration parameters for the ifconfigco retriever.
type IfConfigCoConfig struct {
	BaseURL string `json:"base_url" default:"https://ifconfig.co" doc:"API base URL"`
}

type ifConfigRetriever struct {
	client  *http.Client
	baseURL string
}

func init() {
	plugin.RegisterRetriever("ifconfigco", newIfConfigCoFromConfig, IfConfigCoConfig{})
}

func newIfConfigCoFromConfig(params map[string]any) (plugin.Retriever, error) {
	cfg, err := plugin.Decode[IfConfigCoConfig](params)
	if err != nil {
		return nil, err
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://ifconfig.co"
	}

	return ifConfigRetriever{baseURL: cfg.BaseURL}, nil
}

func (r ifConfigRetriever) GetIPAddress() (string, error) {
	client := r.client
	if client == nil {
		client = http.DefaultClient
	}

	res, err := client.Get(r.baseURL + "/json")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var data map[string]any
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return "", err
	}

	return data["ip"].(string), nil
}

var _ plugin.Retriever = ifConfigRetriever{}
