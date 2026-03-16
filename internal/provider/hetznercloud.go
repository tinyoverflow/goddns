package provider

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type HetznerCloudProvider struct {
	Token string
	Zone  string
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

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.hetzner.cloud/v1/zones/"+p.Zone+"/rrsets/@/A/actions/set_records",
		bytes.NewBuffer(reqDataBytes),
	)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+p.Token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var data map[string]interface{}
	if json.NewDecoder(res.Body).Decode(&data) != nil {
		return err
	}

	return nil
}

var _ Provider = HetznerCloudProvider{}
