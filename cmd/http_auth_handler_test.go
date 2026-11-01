package main

import (
	"bytes"
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

func TestHTTPAuthRegistrationAndDuplicate(t *testing.T) {
	server, store, databasePath := newAuthTestServer(t, 100)

	apiKey := registerAPIKey(t, server.URL, "  User@Example.COM  ", http.StatusOK)
	if len(apiKey) != 64 {
		t.Fatalf("API key length: got %d, want 64", len(apiKey))
	}

	email, ok, err := store.Authenticate(apiKey)
	if err != nil {
		t.Fatalf("authenticate API key: %v", err)
	}
	if !ok {
		t.Fatal("registered API key was not found")
	}
	if email != "user@example.com" {
		t.Errorf("registered email: got %q, want %q", email, "user@example.com")
	}

	response := doRequest(t, http.MethodPost, server.URL+"/register", `{"email":"user@example.com"}`, "application/json", "")
	if response.StatusCode != http.StatusForbidden {
		t.Fatalf("duplicate register status: got %d, want %d", response.StatusCode, http.StatusForbidden)
	}
	assertBody(t, response, `{"error":"forbidden"}`)

	databaseBytes, err := os.ReadFile(databasePath)
	if err != nil {
		t.Fatalf("read auth database: %v", err)
	}
	if strings.Contains(string(databaseBytes), apiKey) {
		t.Error("auth database contains the plaintext API key")
	}
}

func TestHTTPAuthMalformedRegistrationDoesNotCreateRow(t *testing.T) {
	server, _, _ := newAuthTestServer(t, 100)

	tests := []struct {
		name        string
		method      string
		contentType string
		body        string
	}{
		{name: "malformed JSON", method: http.MethodPost, contentType: "application/json", body: `{"email":`},
		{name: "wrong content type", method: http.MethodPost, contentType: "text/plain", body: `{"email":"bad@example.com"}`},
		{name: "missing email", method: http.MethodPost, contentType: "application/json", body: `{}`},
		{name: "empty email", method: http.MethodPost, contentType: "application/json", body: `{"email":"  "}`},
		{name: "oversized body", method: http.MethodPost, contentType: "application/json", body: `{"email":"` + strings.Repeat("a", 1<<10) + `"}`},
		{name: "wrong method", method: http.MethodGet, contentType: "application/json", body: `{"email":"bad@example.com"}`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := doRequest(t, test.method, server.URL+"/register", test.body, test.contentType, "")
			if response.StatusCode != http.StatusForbidden {
				t.Fatalf("status: got %d, want %d", response.StatusCode, http.StatusForbidden)
			}
			assertBody(t, response, `{"error":"forbidden"}`)
		})
	}

	apiKey := registerAPIKey(t, server.URL, "bad@example.com", http.StatusOK)
	if apiKey == "" {
		t.Fatal("expected registration to succeed after malformed attempts")
	}

	response := doRequest(t, http.MethodPost, server.URL+"/register", `{"email":"[OTHER]","unknown":true}`, "application/json", "")
	if response.StatusCode != http.StatusOK {
		t.Fatalf("unknown field status: got %d, want %d", response.StatusCode, http.StatusOK)
	}
}

func TestHTTPAuthMCPRequiresValidBearer(t *testing.T) {
	server, _, _ := newAuthTestServer(t, 100)

	response := doRequest(t, http.MethodPost, server.URL+"/mcp", `{}`, "application/json", "")
	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("missing bearer status: got %d, want %d", response.StatusCode, http.StatusUnauthorized)
	}
	assertEmptyBody(t, response)

	response = doRequest(t, http.MethodPost, server.URL+"/mcp", `{}`, "application/json", "not-a-key")
	if response.StatusCode != http.StatusUnauthorized {
		t.Fatalf("wrong bearer status: got %d, want %d", response.StatusCode, http.StatusUnauthorized)
	}
	assertEmptyBody(t, response)

	health := doRequest(t, http.MethodGet, server.URL+"/healthz", "", "", "")
	if health.StatusCode != http.StatusOK {
		t.Fatalf("healthz status: got %d, want %d", health.StatusCode, http.StatusOK)
	}
	assertBody(t, health, "ok\n")
}

func TestHTTPAuthValidBearerListsResources(t *testing.T) {
	server, _, _ := newAuthTestServer(t, 100)
	apiKey := registerAPIKey(t, server.URL, "resources@example.com", http.StatusOK)

	initialize := postMCP(t, server.URL, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"auth-test","version":"1"}}}`, apiKey)
	if initialize["error"] != nil {
		t.Fatalf("initialize returned error: %v", initialize["error"])
	}

	listResponse := postMCP(t, server.URL, `{"jsonrpc":"2.0","id":2,"method":"resources/list","params":{}}`, apiKey)
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

func TestHTTPAuthRateLimitsAuthenticatedAndRegisterTraffic(t *testing.T) {
	server, _, _ := newAuthTestServer(t, 1)
	apiKey := registerAPIKey(t, server.URL, "limited@example.com", http.StatusOK)

	for attempt := 0; attempt < 2; attempt++ {
		response := doRequest(t, http.MethodPost, server.URL+"/mcp", `{}`, "application/json", apiKey)
		if response.StatusCode == http.StatusTooManyRequests {
			t.Fatalf("authenticated attempt %d was rate limited too early", attempt+1)
		}
	}
	response := doRequest(t, http.MethodPost, server.URL+"/mcp", `{}`, "application/json", apiKey)
	if response.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("authenticated rate limit status: got %d, want %d", response.StatusCode, http.StatusTooManyRequests)
	}
	assertRetryAfter(t, response)
	assertEmptyBody(t, response)

	registerServer, _, _ := newAuthTestServer(t, 1)
	for attempt := 0; attempt < 2; attempt++ {
		response = doRequest(t, http.MethodPost, registerServer.URL+"/register", `{`, "application/json", "")
		if response.StatusCode != http.StatusForbidden {
			t.Fatalf("register attempt %d status: got %d, want %d", attempt+1, response.StatusCode, http.StatusForbidden)
		}
	}
	response = doRequest(t, http.MethodPost, registerServer.URL+"/register", `{`, "application/json", "")
	if response.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("register rate limit status: got %d, want %d", response.StatusCode, http.StatusTooManyRequests)
	}
	assertRetryAfter(t, response)
	assertEmptyBody(t, response)

	unauthenticatedServer, _, _ := newAuthTestServer(t, 1)
	for attempt := 0; attempt < 2; attempt++ {
		response = doRequest(t, http.MethodPost, unauthenticatedServer.URL+"/mcp", `{}`, "application/json", "")
		if response.StatusCode != http.StatusUnauthorized {
			t.Fatalf("unauthenticated MCP attempt %d status: got %d, want %d", attempt+1, response.StatusCode, http.StatusUnauthorized)
		}
	}
	response = doRequest(t, http.MethodPost, unauthenticatedServer.URL+"/mcp", `{}`, "application/json", "")
	if response.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("unauthenticated MCP rate limit status: got %d, want %d", response.StatusCode, http.StatusTooManyRequests)
	}
	assertRetryAfter(t, response)
	assertEmptyBody(t, response)

	for attempt := 0; attempt < 3; attempt++ {
		response = doRequest(t, http.MethodGet, unauthenticatedServer.URL+"/healthz", "", "", "")
		if response.StatusCode != http.StatusOK {
			t.Fatalf("healthz attempt %d status: got %d, want %d", attempt+1, response.StatusCode, http.StatusOK)
		}
	}
}

func newAuthTestServer(t *testing.T, rateLimit int) (*httptest.Server, *auth.Store, string) {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "std.md"), []byte(fixtureDoc), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	mcpServer, err := internal.BootstrapServer(utils.NewAssetsFinder(dir))
	if err != nil {
		t.Fatalf("BootstrapServer: %v", err)
	}
	databasePath := filepath.Join(t.TempDir(), "keys.db")
	store, err := auth.OpenDatabase(databasePath)
	if err != nil {
		t.Fatalf("open auth database: %v", err)
	}
	handler := auth.NewHandler(newMCPHandler(mcpServer), store, rateLimit, zap.NewNop().Sugar())
	server := httptest.NewServer(handler)
	t.Cleanup(func() {
		server.Close()
		if err := store.Close(); err != nil {
			t.Errorf("close auth store: %v", err)
		}
	})
	return server, store, databasePath
}

func registerAPIKey(t *testing.T, baseURL, email string, wantStatus int) string {
	t.Helper()
	response := doRequest(t, http.MethodPost, baseURL+"/register", `{"email":"`+email+`"}`, "application/json", "")
	if response.StatusCode != wantStatus {
		t.Fatalf("register status: got %d, want %d", response.StatusCode, wantStatus)
	}
	var payload struct {
		APIKey string `json:"apiKey"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("decode register response: %v", err)
	}
	return payload.APIKey
}

func doRequest(t *testing.T, method, url, body, contentType, bearer string) *http.Response {
	t.Helper()
	request, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}
	if bearer != "" {
		request.Header.Set("Authorization", "Bearer "+bearer)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}
	if err := response.Body.Close(); err != nil {
		t.Fatalf("close response: %v", err)
	}
	response.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	return response
}

func postMCP(t *testing.T, baseURL, payload, apiKey string) map[string]any {
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

func assertBody(t *testing.T, response *http.Response, want string) {
	t.Helper()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("read response: %v", err)
	}
	if string(body) != want {
		t.Errorf("body: got %q, want %q", body, want)
	}
}

func assertEmptyBody(t *testing.T, response *http.Response) {
	t.Helper()
	assertBody(t, response, "")
}

func assertRetryAfter(t *testing.T, response *http.Response) {
	t.Helper()
	value := response.Header.Get("Retry-After")
	if value == "" {
		t.Fatal("Retry-After header is missing")
	}
	if value != "1" {
		t.Errorf("Retry-After: got %q, want %q", value, "1")
	}
}
