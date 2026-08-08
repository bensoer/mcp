package internal

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"testing"

	"mcp/internal/utils"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type mockAssetsFinder struct {
	paths    []string
	contents map[string]string
	loadErr  error
}

func (m *mockAssetsFinder) GetAssetPath(assetName string) string {
	return assetName
}

func (m *mockAssetsFinder) GetAssetFolderRoot() string {
	return "mock://"
}

func (m *mockAssetsFinder) GetAssetContents(assetPathInsideFolder string) ([]byte, error) {
	if m.loadErr != nil {
		return nil, m.loadErr
	}
	c, ok := m.contents[assetPathInsideFolder]
	if !ok {
		return nil, fmt.Errorf("mock: file not found: %s", assetPathInsideFolder)
	}
	return []byte(c), nil
}

func (m *mockAssetsFinder) GetAllAssetPaths() ([]string, error) {
	if m.loadErr != nil {
		return nil, m.loadErr
	}
	return m.paths, nil
}

var _ utils.AssetsFinder = (*mockAssetsFinder)(nil)

func validAssetYAML(uri, name string) string {
	return fmt.Sprintf(`---
uri: %s
name: %s
description: Test resource
languages:
  - all
file_types:
  - "*.*"
priority: required
---
# %s

Content for %s.
`, uri, name, name, name)
}

func TestBootstrapServer_ValidAssets(t *testing.T) {
	mock := &mockAssetsFinder{
		paths: []string{"a.md", "b.md"},
		contents: map[string]string{
			"a.md": validAssetYAML("test://a", "Resource A"),
			"b.md": validAssetYAML("test://b", "Resource B"),
		},
	}

	server, err := BootstrapServer(mock)
	if err != nil {
		t.Fatalf("BootstrapServer() error: %v", err)
	}
	if server == nil {
		t.Fatal("BootstrapServer() returned nil server")
	}
}

func TestBootstrapServer_MissingURI(t *testing.T) {
	mock := &mockAssetsFinder{
		paths: []string{"bad.md"},
		contents: map[string]string{
			"bad.md": `---
name: Missing URI
---
# Content
`,
		},
	}

	_, err := BootstrapServer(mock)
	if err == nil {
		t.Error("BootstrapServer() expected error for missing uri, got nil")
	}
}

func TestBootstrapServer_MissingName(t *testing.T) {
	mock := &mockAssetsFinder{
		paths: []string{"bad.md"},
		contents: map[string]string{
			"bad.md": `---
uri: test://noname
---
# Content
`,
		},
	}

	_, err := BootstrapServer(mock)
	if err == nil {
		t.Error("BootstrapServer() expected error for missing name, got nil")
	}
}

func TestBootstrapServer_InvalidYAML(t *testing.T) {
	mock := &mockAssetsFinder{
		paths: []string{"bad.md"},
		contents: map[string]string{
			"bad.md": "not valid yaml at all",
		},
	}

	_, err := BootstrapServer(mock)
	if err == nil {
		t.Error("BootstrapServer() expected error for invalid YAML, got nil")
	}
}

func TestBootstrapServer_GetAllAssetPathsError(t *testing.T) {
	mock := &mockAssetsFinder{
		loadErr: errors.New("disk failure"),
	}

	_, err := BootstrapServer(mock)
	if err == nil {
		t.Error("BootstrapServer() expected error from GetAllAssetPaths, got nil")
	}
}

func TestBootstrapServer_GetAssetContentsError(t *testing.T) {
	mock := &mockAssetsFinder{
		paths:    []string{"a.md"},
		contents: map[string]string{},
	}

	_, err := BootstrapServer(mock)
	if err == nil {
		t.Error("BootstrapServer() expected error from GetAssetContents, got nil")
	}
}

func TestBootstrapServer_Integration(t *testing.T) {
	mock := &mockAssetsFinder{
		paths: []string{"greeting.md"},
		contents: map[string]string{
			"greeting.md": validAssetYAML("standards://test/greeting", "Greeting"),
		},
	}

	server, err := BootstrapServer(mock)
	if err != nil {
		t.Fatalf("BootstrapServer() error: %v", err)
	}

	ctx := context.Background()
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "1.0.0"}, nil)

	t1, t2 := mcp.NewInMemoryTransports()
	if _, err := server.Connect(ctx, t1, nil); err != nil {
		t.Fatalf("server.Connect() error: %v", err)
	}
	cs, err := client.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatalf("client.Connect() error: %v", err)
	}
	defer cs.Close()

	resources := make(map[string]mcp.Resource)
	for r, err := range cs.Resources(ctx, nil) {
		if err != nil {
			t.Fatalf("cs.Resources() error: %v", err)
		}
		resources[r.URI] = *r
	}

	if len(resources) != 1 {
		t.Fatalf("expected 1 resource, got %d: %v", len(resources), resources)
	}

	r, ok := resources["standards://test/greeting"]
	if !ok {
		t.Fatalf("resource 'standards://test/greeting' not found in: %v", maps.Keys(resources))
	}
	if r.Name != "Greeting" {
		t.Errorf("resource name = %q, want %q", r.Name, "Greeting")
	}

	result, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "standards://test/greeting"})
	if err != nil {
		t.Fatalf("cs.ReadResource() error: %v", err)
	}
	if len(result.Contents) != 1 {
		t.Fatalf("expected 1 content item, got %d", len(result.Contents))
	}
	if result.Contents[0].MIMEType != "text/markdown" {
		t.Errorf("MIMEType = %q, want %q", result.Contents[0].MIMEType, "text/markdown")
	}

	expectedContent := mock.contents["greeting.md"]
	if result.Contents[0].Text != expectedContent {
		t.Errorf("content mismatch:\ngot:  %q\nwant: %q", result.Contents[0].Text, expectedContent)
	}
}

func TestBootstrapServer_Integration_UnknownResource(t *testing.T) {
	mock := &mockAssetsFinder{
		paths:    []string{"greeting.md"},
		contents: map[string]string{
			"greeting.md": validAssetYAML("standards://test/greeting", "Greeting"),
		},
	}

	server, err := BootstrapServer(mock)
	if err != nil {
		t.Fatalf("BootstrapServer() error: %v", err)
	}

	ctx := context.Background()
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "1.0.0"}, nil)

	t1, t2 := mcp.NewInMemoryTransports()
	if _, err := server.Connect(ctx, t1, nil); err != nil {
		t.Fatalf("server.Connect() error: %v", err)
	}
	cs, err := client.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatalf("client.Connect() error: %v", err)
	}
	defer cs.Close()

	_, err = cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: "standards://test/nonexistent"})
	if err == nil {
		t.Error("cs.ReadResource() expected error for unknown resource, got nil")
	}
}