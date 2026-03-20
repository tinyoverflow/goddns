package retriever

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

type UnifiRetriever struct {
	baseURL  string
	apiToken string
	siteID   string
	client   *http.Client
}

func NewUnifiRetrieverFromConfig(params map[string]any) (Retriever, error) {
	baseURL, _ := params["base_url"].(string)
	apiToken, _ := params["api_token"].(string)
	siteID, _ := params["site_id"].(string)
	verifyTLS, _ := params["verify_tls"].(bool)

	if baseURL == "" {
		return nil, fmt.Errorf("unifi retriever: base_url is required")
	}

	if apiToken == "" {
		return nil, fmt.Errorf("unifi retriever: api_token is required")
	}

	if siteID == "" {
		siteID = "default"
	}

	tlsConfig := &tls.Config{}
	if !verifyTLS {
		tlsConfig.InsecureSkipVerify = true
	}

	return UnifiRetriever{
		baseURL:  baseURL,
		apiToken: apiToken,
		siteID:   siteID,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		},
	}, nil
}

func (r UnifiRetriever) GetIPAddress() (string, error) {
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

var _ Retriever = UnifiRetriever{}
