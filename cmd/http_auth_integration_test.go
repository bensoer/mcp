package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"mcp/internal"
	"mcp/internal/auth"
	"mcp/internal/utils"

	"go.uber.org/zap"
)

func TestHTTPAuthValidBearerListsMCPResources(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "std.md"), []byte(fixtureDoc), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	mcpServer, err := internal.BootstrapServer(utils.NewAssetsFinder(dir))
	if err != nil {
		t.Fatalf("BootstrapServer: %v", err)
	}
	store, err := auth.OpenDatabase(filepath.Join(t.TempDir(), "keys.db"))
	if err != nil {
		t.Fatalf("open auth database: %v", err)
	}
	t.Cleanup(func() {
		if err := store.Close(); err != nil {
			t.Errorf("close auth store: %v", err)
		}
	})
	handler := auth.NewHandler(newMCPHandler(mcpServer), store, 100, zap.NewNop().Sugar())
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	apiKey := registerAPIKey(t, server.URL)
	initialize := postAuthenticatedMCP(t, server.URL, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"auth-test","version":"1"}}}`, apiKey)
	if initialize["error"] != nil {
		t.Fatalf("initialize returned error: %v", initialize["error"])
	}

	listResponse := postAuthenticatedMCP(t, server.URL, `{"jsonrpc":"2.0","id":2,"method":"resources/list","params":{}}`, apiKey)
	if listResponse["error"] != nil {
		t.Fatalf("resources/list returned error: %v", listResponse["error"])
	}
	result, ok := listResponse["result"].(map[string]any)
	if !ok {
		t.Fatalf("resources/list result: got %#v", listResponse["result"])
	}
	resources, ok := result["resources"].([]any)
	if !ok {
		t.Fatalf("resources/list resources: got %#v", result["resources"])
	}
	found := false
	for _, resource := range resources {
		resourceMap, _ := resource.(map[string]any)
		if resourceMap["uri"] == "standards://plan/http" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected resource %q in resources/list; got %#v", "standards://plan/http", resources)
	}
}

func registerAPIKey(t *testing.T, baseURL string) string {
	t.Helper()
	request, err := http.NewRequest(http.MethodPost, baseURL+"/register", strings.NewReader(`{"email":"[EMAIL]"}`))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	defer func() { _ = response.Body.Close() }()
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		t.Fatalf("register status: got %d, want %d; body %s", response.StatusCode, http.StatusOK, body)
	}
	var payload struct {
		APIKey string `json:"apiKey"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("decode register response: %v", err)
	}
	return payload.APIKey
}

func postAuthenticatedMCP(t *testing.T, baseURL, payload, apiKey string) map[string]any {
	t.Helper()
	request, err := http.NewRequest(http.MethodPost, baseURL+"/mcp", strings.NewReader(payload))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json, text/event-stream")
	request.Header.Set("Authorization", "Bearer "+apiKey)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer func() { _ = response.Body.Close() }()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Fatalf("post status: got %d, want %d; body %s", response.StatusCode, http.StatusOK, body)
	}
	var decoded map[string]any
	for _, line := range strings.Split(string(body), "\n") {
		data, ok := strings.CutPrefix(line, "data: ")
		if !ok {
			continue
		}
		if err := json.Unmarshal([]byte(data), &decoded); err != nil {
			t.Fatalf("decode SSE data: %v; body %s", err, body)
		}
		break
	}
	if decoded == nil {
		t.Fatalf("no JSON-RPC response in SSE body %s", body)
	}
	return decoded
}
