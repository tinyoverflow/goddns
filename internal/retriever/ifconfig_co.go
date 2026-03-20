package retriever

import (
	"encoding/json"
	"net/http"
)

type IfConfigRetriever struct {
	client *http.Client
	url    string
}

func (r IfConfigRetriever) GetIPAddress() (string, error) {
	client := r.client
	if client == nil {
		client = http.DefaultClient
	}

	url := r.url
	if url == "" {
		url = "https://ifconfig.co/json"
	}

	res, err := client.Get(url)
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
