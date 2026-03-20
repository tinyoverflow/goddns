package plugins

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIfConfigRetriever_GetIPAddress(t *testing.T) {
	const expectedIP = "203.0.113.42"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/json" {
			t.Fatalf("unexpected request path: %v", r.URL.Path)
		}

		json.NewEncoder(w).Encode(map[string]any{"ip": expectedIP})
	}))
	defer server.Close()

	ret := ifConfigRetriever{client: server.Client(), baseURL: server.URL}

	ip, err := ret.GetIPAddress()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ip != expectedIP {
		t.Errorf("expected IP %q, got %q", expectedIP, ip)
	}
}

func TestIfConfigRetriever_GetIPAddress_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	ret := ifConfigRetriever{client: server.Client(), baseURL: server.URL}

	_, err := ret.GetIPAddress()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
