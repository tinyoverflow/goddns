package plugins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goddns/internal/plugin"
	"net/http"
)

// HetznerCloudConfig holds the configuration parameters for the hetzner_cloud provider.
type HetznerCloudConfig struct {
	APIToken string `json:"api_token" required:"true" doc:"Hetzner Cloud API token"`
	Zone     string `json:"zone"      required:"true" doc:"DNS zone (e.g. homelab.com)"`
	RRName   string `json:"rr_name"   required:"true" doc:"Resource record name (e.g. @)"`
}

type hetznerCloudProvider struct {
	token  string
	zone   string
	rrName string
}

func init() {
	plugin.RegisterProvider("hetzner_cloud", newHetznerCloudFromConfig, HetznerCloudConfig{})
}

func newHetznerCloudFromConfig(params map[string]any) (plugin.Provider, error) {
	cfg, err := plugin.Decode[HetznerCloudConfig](params)
	if err != nil {
		return nil, fmt.Errorf("hetzner_cloud: %w", err)
	}

	if cfg.APIToken == "" {
		return nil, fmt.Errorf("hetzner_cloud: api_token is required")
	}
	if cfg.Zone == "" {
		return nil, fmt.Errorf("hetzner_cloud: zone is required")
	}
	if cfg.RRName == "" {
		return nil, fmt.Errorf("hetzner_cloud: rr_name is required")
	}

	return hetznerCloudProvider{
		token:  cfg.APIToken,
		zone:   cfg.Zone,
		rrName: cfg.RRName,
	}, nil
}

func (p hetznerCloudProvider) SetIPAddress(ip string) error {
	client := http.DefaultClient

	reqData := struct {
		Records []struct {
			Value   string `json:"value"`
			Comment string `json:"comment"`
		} `json:"records"`
	}{
		Records: []struct {
			Value   string `json:"value"`
			Comment string `json:"comment"`
		}{
			{Value: ip, Comment: "Updated by goddns"},
		},
	}

	reqDataBytes, err := json.Marshal(reqData)
	if err != nil {
		return err
	}

	url := "https://api.hetzner.cloud/v1/zones/" + p.zone + "/rrsets/" + p.rrName + "/A/actions/set_records"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqDataBytes))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+p.token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status 200, got %d", res.StatusCode)
	}

	var data map[string]any
	if json.NewDecoder(res.Body).Decode(&data) != nil {
		return err
	}

	return nil
}

var _ plugin.Provider = hetznerCloudProvider{}
