package tools

import (
	"context"
	"encoding/json"
	"testing"

	"mcp/internal/utils"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// newTestServer creates a server with only tools registered (no resources).
func newTestServer(finder utils.AssetsFinder) (*mcp.Server, error) {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "test-server",
		Version: "1.0.0",
	}, nil)
	for _, register := range RegisterAll {
		register(server, finder)
	}
	return server, nil
}

// newConnectedClient wires up an in-memory MCP client+server pair with the
// given assets and returns the connected client session.
func newConnectedClient(t *testing.T, mock *mockAssetsFinder) (context.Context, *mcp.ClientSession) {
	t.Helper()
	server, err := newTestServer(mock)
	if err != nil {
		t.Fatalf("newTestServer() error: %v", err)
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
	t.Cleanup(func() {
		if err := cs.Close(); err != nil {
			t.Logf("error closing client session: %v", err)
		}
	})
	return ctx, cs
}

// unwrapToolResult calls a tool and unmarshals its StructuredContent into out.
// It fails the test on any error or unexpected type, so callers can use out
// directly without nil/reassignment tracking issues.
func unwrapToolResult(t *testing.T, ctx context.Context, cs *mcp.ClientSession, name string, args map[string]any, out any) {
	t.Helper()
	result, err := cs.CallTool(ctx, &mcp.CallToolParams{
		Name:      name,
		Arguments: args,
	})
	if err != nil {
		t.Fatalf("cs.CallTool(%q) error: %v", name, err)
	}
	if result == nil {
		t.Fatalf("cs.CallTool(%q) returned nil result", name)
	}
	outputMap, ok := result.StructuredContent.(map[string]any)
	if !ok {
		t.Fatalf("expected map for StructuredContent, got %T", result.StructuredContent)
	}
	b, err := json.Marshal(outputMap)
	if err != nil {
		t.Fatalf("failed to marshal StructuredContent: %v", err)
	}
	if err := json.Unmarshal(b, out); err != nil {
		t.Fatalf("failed to unmarshal StructuredContent: %v", err)
	}
}

// TestBootstrapServer_CallTool_DiscoverStandards tests the discover_standards tool.
func TestBootstrapServer_CallTool_DiscoverStandards(t *testing.T) {
	mock := &mockAssetsFinder{
		paths: []string{"a.md", "b.md"},
		contents: map[string]string{
			"a.md": validAssetYAML("standards://test/a", "Resource A"),
			"b.md": validAssetYAML("standards://test/b", "Resource B"),
		},
	}
	ctx, cs := newConnectedClient(t, mock)

	var output struct {
		Resources []ResourceMetadata `json:"resources"`
	}
	unwrapToolResult(t, ctx, cs, "discover_standards", map[string]any{}, &output)

	if len(output.Resources) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(output.Resources))
	}
	for _, r := range output.Resources {
		if !hasPrefix(r.URI, "standards://") {
			t.Errorf("resource URI = %q, want prefix standards://", r.URI)
		}
	}
}

// TestBootstrapServer_CallTool_SearchStandards tests the search_standards tool.
func TestBootstrapServer_CallTool_SearchStandards(t *testing.T) {
	mock := &mockAssetsFinder{
		paths: []string{"commit.md", "workflow.md", "other.md"},
		contents: map[string]string{
			"commit.md":   validAssetYAML("standards://git/commit-messages", "Git Commit Messages"),
			"workflow.md": validAssetYAML("workflows://test/foo", "Test Workflow"),
			"other.md":    validAssetYAML("standards://docs/guide", "Documentation Guide"),
		},
	}
	ctx, cs := newConnectedClient(t, mock)

	var output struct {
		Query     string             `json:"query"`
		Resources []ResourceMetadata `json:"resources"`
		Total     int                `json:"total"`
	}

	// Test 1: Search for "commit" (should match the commit messages resource)
	unwrapToolResult(t, ctx, cs, "search_standards", map[string]any{"query": "commit"}, &output)
	if len(output.Resources) == 0 {
		t.Fatalf("expected at least one match for 'commit', got 0")
	}
	for _, r := range output.Resources {
		if !hasPrefix(r.URI, "standards://") {
			t.Errorf("resource URI = %q, want prefix standards://", r.URI)
		}
	}

	// Test 2: Search for "TEST" (case-insensitive, should match workflow and others)
	unwrapToolResult(t, ctx, cs, "search_standards", map[string]any{"query": "TEST"}, &output)
	if len(output.Resources) == 0 {
		t.Fatalf("expected matches for 'TEST', got 0")
	}

	// Test 3: Search for "content" in body (should match resources whose body has it)
	unwrapToolResult(t, ctx, cs, "search_standards", map[string]any{"query": "content"}, &output)
	if len(output.Resources) == 0 {
		t.Fatalf("expected matches for 'content', got 0")
	}

	// Test 4: Search for "nonexistent" (should return empty array, not error)
	unwrapToolResult(t, ctx, cs, "search_standards", map[string]any{"query": "nonexistent"}, &output)
	if len(output.Resources) != 0 {
		t.Fatalf("expected zero matches for 'nonexistent', got %d", len(output.Resources))
	}
}

// TestBootstrapServer_CallTool_SearchWorkflows_EmptyQuery tests that calling
// search_workflows with an empty query does not cause a handler error.
// The SDK does not enforce the "required" tag for an empty string, so an
// empty query is valid and returns all workflows.
func TestBootstrapServer_CallTool_SearchWorkflows_EmptyQuery(t *testing.T) {
	mock := &mockAssetsFinder{
		paths: []string{"workflow.md"},
		contents: map[string]string{
			"workflow.md": validAssetYAML("workflows://test/foo", "Test Workflow"),
		},
	}
	ctx, cs := newConnectedClient(t, mock)

	t.Run("MissingQuery", func(t *testing.T) {
		var output struct {
			Query     string             `json:"query"`
			Resources []ResourceMetadata `json:"resources"`
			Total     int                `json:"total"`
		}
		unwrapToolResult(t, ctx, cs, "search_workflows", map[string]any{"query": ""}, &output)
		t.Logf("search_workflows with empty query returned %d resources", output.Total)
	})
}

// TestRegisterAll smoke tests that all tools register without panic.
func TestRegisterAll(t *testing.T) {
	mock := &mockAssetsFinder{
		paths: []string{"test.md"},
		contents: map[string]string{
			"test.md": validAssetYAML("standards://test", "Test"),
		},
	}
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "1.0.0"}, nil)
	for _, register := range RegisterAll {
		// Should not panic
		register(server, mock)
	}
}

// Helper function to check if a string has a prefix.
func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
