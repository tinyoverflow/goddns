package retriever

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

type UnifiRetriever struct {
	baseURL string
	apiKey  string
	siteID  string
	client  *http.Client
}

func NewUnifiRetriever(host, apiKey, siteID string) UnifiRetriever {
	if siteID == "" {
		siteID = "default"
	}
	
	return UnifiRetriever{
		baseURL: "https://" + host,
		apiKey:  apiKey,
		siteID:  siteID,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // UDM uses self-signed certs
			},
		},
	}
}

func (r UnifiRetriever) GetIPAddress() (string, error) {
	url := fmt.Sprintf("%s/proxy/network/api/s/%s/stat/health", r.baseURL, r.siteID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("X-API-KEY", r.apiKey)
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

var _ Retriever = UnifiRetriever{}
