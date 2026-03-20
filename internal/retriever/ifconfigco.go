package retriever

import (
	"encoding/json"
	"net/http"
)

type IfConfigRetriever struct {
	client  *http.Client
	BaseURL string
}

func NewIfConfigCoRetrieverFromConfig(params map[string]any) (Retriever, error) {
	baseURL, _ := params["base_url"].(string)

	if baseURL == "" {
		baseURL = "https://ifconfig.co"
	}

	return IfConfigRetriever{BaseURL: baseURL}, nil
}

func (r IfConfigRetriever) GetIPAddress() (string, error) {
	client := r.client
	if client == nil {
		client = http.DefaultClient
	}

	res, err := client.Get(r.BaseURL + "/json")
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

var _ Retriever = IfConfigRetriever{}
