package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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
