package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type HetznerCloudProvider struct {
	Token  string
	Zone   string
	RRName string
}

func NewHetznerCloudProviderFromConfig(params map[string]any) (Provider, error) {
	apiToken, _ := params["api_token"].(string)
	zone, _ := params["zone"].(string)
	rrName, _ := params["rr_name"].(string)

	if apiToken == "" {
		return nil, fmt.Errorf("api_token is required")
	}

	if zone == "" {
		return nil, fmt.Errorf("zone is required")
	}

	if rrName == "" {
		return nil, fmt.Errorf("rr_name is required")
	}

	return HetznerCloudProvider{
		Token:  apiToken,
		Zone:   zone,
		RRName: rrName,
	}, nil
}

func (p HetznerCloudProvider) SetIPAddress(ip string) error {
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
			{
				Value:   ip,
				Comment: "Updated by goddns",
			},
		},
	}

	reqDataBytes, err := json.Marshal(reqData)
	if err != nil {
		return err
	}

	url := "https://api.hetzner.cloud/v1/zones/" + p.Zone + "/rrsets/" + p.RRName + "/A/actions/set_records"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqDataBytes))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+p.Token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status 200, got %d", res.StatusCode)
	}

	defer res.Body.Close()

	var data map[string]any
	if json.NewDecoder(res.Body).Decode(&data) != nil {
		return err
	}

	return nil
}

var _ Provider = HetznerCloudProvider{}
