package retriever

import (
	"encoding/json"
	"net/http"
)

type IfConfigRetriever struct{}

func (r IfConfigRetriever) GetIPAddress() (string, error) {
	client := http.DefaultClient

	res, err := client.Get("https://ifconfig.co/json")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var data map[string]any
	if json.NewDecoder(res.Body).Decode(&data) != nil {
		return "", err
	}

	return data["ip"].(string), nil
}

var _ Retriever = IfConfigRetriever{}
