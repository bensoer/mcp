package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"mcp/internal"
	"mcp/internal/utils"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// fixtureDoc is a minimal resource document that BootstrapServer accepts:
// non-empty uri (with a scheme), non-empty name, valid frontmatter.
const fixtureDoc = `---
uri: standards://plan/http
name: HTTP Plan Fixture
description: fixture
---

fixture
`

func TestHTTPHandlerLegacyInitializeWithoutSession(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "std.md"), []byte(fixtureDoc), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	server, err := internal.BootstrapServer(utils.NewAssetsFinder(dir))
	if err != nil {
		t.Fatalf("BootstrapServer: %v", err)
	}

	ts := httptest.NewServer(newMCPHandler(server))
	t.Cleanup(ts.Close)

	post := func(payload string) map[string]any {
		t.Helper()
		req, err := http.NewRequest(http.MethodPost, ts.URL+"/mcp", strings.NewReader(payload))
		if err != nil {
			t.Fatalf("new request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json, text/event-stream")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("post: %v", err)
		}
		defer func() { _ = resp.Body.Close() }()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("read response: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("post status: got %d, want %d; body %s", resp.StatusCode, http.StatusOK, body)
		}
		var response map[string]any
		for _, line := range strings.Split(string(body), "\n") {
			data, ok := strings.CutPrefix(line, "data: ")
			if !ok {
				continue
			}
			if err := json.Unmarshal([]byte(data), &response); err != nil {
				t.Fatalf("decode SSE data: %v; body %s", err, body)
			}
			break
		}
		if response == nil {
			t.Fatalf("no JSON-RPC response in SSE body %s", body)
		}
		return response
	}

	initResponse := post(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"legacy-test","version":"1"}}}`)
	if initResponse["error"] != nil {
		t.Fatalf("legacy initialize returned error: %v", initResponse["error"])
	}
	if got := initResponse["result"].(map[string]any)["protocolVersion"]; got != "2024-11-05" {
		t.Fatalf("negotiated protocol version: got %v, want %q", got, "2024-11-05")
	}

	listResponse := post(`{"jsonrpc":"2.0","id":2,"method":"resources/list","params":{}}`)
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
		t.Errorf("expected resource %q in legacy resources/list; got %#v", "standards://plan/http", resources)
	}

	getResponse, err := http.Get(ts.URL + "/mcp")
	if err != nil {
		t.Fatalf("GET /mcp: %v", err)
	}
	defer func() { _ = getResponse.Body.Close() }()
	if getResponse.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("GET /mcp status: got %d, want %d", getResponse.StatusCode, http.StatusMethodNotAllowed)
	}
}

func TestHTTPHandlerListsResourcesAndHealthz(t *testing.T) {
	// Build a throwaway assets tree in a temp dir; do not touch repo assets/.
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "std.md"), []byte(fixtureDoc), 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	server, err := internal.BootstrapServer(utils.NewAssetsFinder(dir))
	if err != nil {
		t.Fatalf("BootstrapServer: %v", err)
	}

	ts := httptest.NewServer(newMCPHandler(server))
	t.Cleanup(ts.Close)

	// httptest listens on loopback; the client Host is loopback, so DNS
	// rebinding protection must not 403 this request.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)

	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "0"}, nil)
	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: ts.URL + "/mcp"}, nil)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(func() { _ = session.Close() })

	res, err := session.ListResources(ctx, nil)
	if err != nil {
		t.Fatalf("ListResources: %v", err)
	}

	found := false
	for _, r := range res.Resources {
		if r.URI == "standards://plan/http" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected resource %q in ListResources; got %d resources", "standards://plan/http", len(res.Resources))
	}

	// GET /healthz is a plain HTTP probe, independent of the MCP client.
	resp, err := http.Get(ts.URL + "/healthz")
	if err != nil {
		t.Fatalf("healthz get: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("healthz status: got %d, want %d", resp.StatusCode, http.StatusOK)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("healthz read: %v", err)
	}
	if string(body) != "ok\n" {
		t.Errorf("healthz body: got %q, want %q", string(body), "ok\n")
	}
}
