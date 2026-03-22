package plugins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goddns/internal/plugin"
	"net/http"
)

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

type hetznerCloudProviderRequest struct {
	Records []hetznerCloudProviderRequestRecord `json:"records"`
}

type hetznerCloudProviderRequestRecord struct {
	Value   string `json:"value"`
	Comment string `json:"comment"`
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

	reqData := hetznerCloudProviderRequest{
		Records: []hetznerCloudProviderRequestRecord{
			{Value: ip},
		},
	}

	reqDataBytes, err := json.Marshal(reqData)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.hetzner.cloud/v1/zones/%v/rrsets/%v/A/actions/set_records", p.zone, p.rrName)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqDataBytes))
	if err != nil {
		return err
	}

	authorizationHeaderValue := fmt.Sprintf("Bearer %v", p.token)
	req.Header.Add("Authorization", authorizationHeaderValue)
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
