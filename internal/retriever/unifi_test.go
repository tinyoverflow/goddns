package retriever

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUnifiRetriever_GetIPAddress(t *testing.T) {
	const testAPIKey = "test-api-key"
	const expectedIP = "203.0.113.42"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-KEY") != testAPIKey {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"subsystem": "www", "wan_ip": expectedIP},
			},
		})
	}))
	defer server.Close()

	ret := UnifiRetriever{baseURL: server.URL, apiToken: testAPIKey, siteID: "default", client: server.Client()}

	ip, err := ret.GetIPAddress()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ip != expectedIP {
		t.Errorf("expected IP %q, got %q", expectedIP, ip)
	}
}

func TestUnifiRetriever_GetIPAddress_Unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	ret := UnifiRetriever{baseURL: server.URL, apiToken: "wrong-key", siteID: "default", client: server.Client()}

	_, err := ret.GetIPAddress()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUnifiRetriever_GetIPAddress_MissingWanIP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"subsystem": "lan"},
			},
		})
	}))
	defer server.Close()

	ret := UnifiRetriever{baseURL: server.URL, apiToken: "key", siteID: "default", client: server.Client()}

	_, err := ret.GetIPAddress()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUnifiRetriever_GetIPAddress_CustomSiteID(t *testing.T) {
	const customSite = "my-site"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/proxy/network/api/s/" + customSite + "/stat/health"
		if r.URL.Path != expected {
			t.Errorf("expected path %q, got %q", expected, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []map[string]any{
				{"wan_ip": "1.2.3.4"},
			},
		})
	}))
	defer server.Close()

	ret := UnifiRetriever{baseURL: server.URL, apiToken: "key", siteID: customSite, client: server.Client()}

	_, err := ret.GetIPAddress()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
