//go:build e2e

package e2e

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	commitMessagesURI = "standards://git/commit-messages"
	serverName        = "Ben's Standards, Practices, Linting & Expectations For AI Agents"
)

func TestMain(m *testing.M) {
	if os.Getenv("MCP_E2E_BASE_URL") == "" {
		fmt.Fprintln(os.Stderr, "MCP_E2E_BASE_URL is required")
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestMCPHTTPService(t *testing.T) {
	baseURL := strings.TrimRight(os.Getenv("MCP_E2E_BASE_URL"), "/")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	t.Cleanup(cancel)

	t.Run("healthz", func(t *testing.T) {
		waitForHealth(t, ctx, baseURL+"/healthz")
	})

	t.Run("root is not health endpoint", func(t *testing.T) {
		status, body := get(t, ctx, baseURL+"/")
		if status == http.StatusOK && body == "ok\n" {
			t.Fatal("GET / returned the health response")
		}
	})

	authEnabled := os.Getenv("MCP_E2E_AUTH_ENABLED") == "true"
	apiKey := ""
	if authEnabled {
		t.Run("mcp requires bearer authentication", func(t *testing.T) {
			status, body := post(t, ctx, baseURL+"/mcp", "{}", "application/json", "")
			if status != http.StatusUnauthorized {
				t.Fatalf("unauthenticated MCP status = %d, want %d; body %q", status, http.StatusUnauthorized, body)
			}
		})
		t.Run("register API key", func(t *testing.T) {
			var err error
			apiKey, err = registerAPIKey(t, ctx, baseURL)
			if err != nil {
				t.Fatalf("register API key: %v", err)
			}
		})
	}

	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-e2e", Version: "1.0.0"}, nil)
	session, err := client.Connect(ctx, streamableTransport(baseURL+"/mcp", apiKey), nil)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(func() { _ = session.Close() })

	t.Run("initialize", func(t *testing.T) {
		result := session.InitializeResult()
		if result == nil || result.ServerInfo == nil {
			t.Fatal("initialize result has no server info")
		}
		if result.ServerInfo.Name != serverName {
			t.Fatalf("server name = %q, want %q", result.ServerInfo.Name, serverName)
		}
	})

	t.Run("list resources", func(t *testing.T) {
		result, err := session.ListResources(ctx, nil)
		if err != nil {
			t.Fatalf("ListResources: %v", err)
		}
		if len(result.Resources) == 0 {
			t.Fatal("ListResources returned no resources")
		}
		found := false
		for _, resource := range result.Resources {
			if resource.URI == commitMessagesURI {
				found = true
				if resource.MIMEType != "text/markdown" {
					t.Errorf("MIME type = %q, want %q", resource.MIMEType, "text/markdown")
				}
			}
		}
		if !found {
			t.Fatalf("resource %q not found", commitMessagesURI)
		}
	})

	t.Run("read resource", func(t *testing.T) {
		result, err := session.ReadResource(ctx, &mcp.ReadResourceParams{URI: commitMessagesURI})
		if err != nil {
			t.Fatalf("ReadResource: %v", err)
		}
		if len(result.Contents) != 1 {
			t.Fatalf("content count = %d, want 1", len(result.Contents))
		}
		content := result.Contents[0]
		if content.MIMEType != "text/markdown" {
			t.Errorf("MIME type = %q, want %q", content.MIMEType, "text/markdown")
		}
		if !strings.Contains(content.Text, "uri: "+commitMessagesURI) {
			t.Errorf("content does not contain frontmatter URI %q", commitMessagesURI)
		}
	})

	t.Run("unknown resource errors", func(t *testing.T) {
		_, err := session.ReadResource(ctx, &mcp.ReadResourceParams{URI: "standards://test/does-not-exist"})
		if err == nil {
			t.Fatal("ReadResource returned nil error for unknown resource")
		}
	})

	t.Run("list tools", func(t *testing.T) {
		result, err := session.ListTools(ctx, nil)
		if err != nil {
			t.Fatalf("ListTools: %v", err)
		}
		got := make([]string, 0, len(result.Tools))
		for _, tool := range result.Tools {
			got = append(got, tool.Name)
		}
		slices.Sort(got)
		want := []string{
			"discover_standards",
			"discover_workflows",
			"search_standards",
			"search_workflows",
		}
		if !slices.Equal(got, want) {
			t.Fatalf("tool names = %v, want %v", got, want)
		}
	})

	t.Run("discover standards", func(t *testing.T) {
		var output resourceListOutput
		callTool(t, ctx, session, "discover_standards", map[string]any{}, &output)
		if len(output.Resources) == 0 {
			t.Fatal("discover_standards returned no resources")
		}
		found := false
		for _, resource := range output.Resources {
			if !strings.HasPrefix(resource.URI, "standards://") {
				t.Errorf("URI = %q, want standards:// prefix", resource.URI)
			}
			found = found || resource.URI == commitMessagesURI
		}
		if !found {
			t.Fatalf("discover_standards did not return %q", commitMessagesURI)
		}
	})

	t.Run("discover workflows", func(t *testing.T) {
		var output resourceListOutput
		callTool(t, ctx, session, "discover_workflows", map[string]any{}, &output)
		if len(output.Resources) == 0 {
			t.Fatal("discover_workflows returned no resources")
		}
		for _, resource := range output.Resources {
			if !strings.HasPrefix(resource.URI, "workflows://") {
				t.Errorf("URI = %q, want workflows:// prefix", resource.URI)
			}
		}
	})

	t.Run("search standards", func(t *testing.T) {
		var output resourceListOutput
		callTool(t, ctx, session, "search_standards", map[string]any{"query": "commit"}, &output)
		if len(output.Resources) == 0 {
			t.Fatal("search_standards returned no matches for commit")
		}
		for _, resource := range output.Resources {
			if !strings.HasPrefix(resource.URI, "standards://") {
				t.Errorf("URI = %q, want standards:// prefix", resource.URI)
			}
		}
	})

	t.Run("search standards empty result", func(t *testing.T) {
		var output resourceListOutput
		callTool(t, ctx, session, "search_standards", map[string]any{"query": "nonexistent-token-xyz"}, &output)
		if len(output.Resources) != 0 {
			t.Fatalf("search_standards returned %d resources, want 0", len(output.Resources))
		}
	})

	t.Run("search workflows empty query", func(t *testing.T) {
		var output resourceListOutput
		callTool(t, ctx, session, "search_workflows", map[string]any{"query": ""}, &output)
		if len(output.Resources) == 0 {
			t.Fatal("search_workflows returned no workflows")
		}
	})

	t.Run("second client connects", func(t *testing.T) {
		secondClient := mcp.NewClient(&mcp.Implementation{Name: "mcp-e2e-second", Version: "1.0.0"}, nil)
		secondSession, err := secondClient.Connect(ctx, streamableTransport(baseURL+"/mcp", apiKey), nil)
		if err != nil {
			t.Fatalf("second connect: %v", err)
		}
		if err := secondSession.Close(); err != nil {
			t.Errorf("close second session: %v", err)
		}
	})

	t.Run("garbage post does not stop server", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/mcp", strings.NewReader("not valid json"))
		if err != nil {
			t.Fatalf("create request: %v", err)
		}
		if apiKey != "" {
			req.Header.Set("Authorization", "Bearer "+apiKey)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("POST /mcp: %v", err)
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()

		status, body := get(t, ctx, baseURL+"/healthz")
		if status != http.StatusOK || body != "ok\n" {
			t.Fatalf("health after garbage POST = (%d, %q), want (%d, %q)", status, body, http.StatusOK, "ok\n")
		}
	})
}

type resourceMetadata struct {
	URI string `json:"uri"`
}

type resourceListOutput struct {
	Resources []resourceMetadata `json:"resources"`
}

type bearerRoundTripper struct {
	token string
	base  http.RoundTripper
}

func (t bearerRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	cloned := request.Clone(request.Context())
	cloned.Header.Set("Authorization", "Bearer "+t.token)
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(cloned)
}

func streamableTransport(endpoint, apiKey string) *mcp.StreamableClientTransport {
	transport := &mcp.StreamableClientTransport{Endpoint: endpoint}
	if apiKey != "" {
		transport.HTTPClient = &http.Client{Transport: bearerRoundTripper{token: apiKey}}
	}
	return transport
}

func registerAPIKey(t *testing.T, ctx context.Context, baseURL string) (string, error) {
	t.Helper()
	status, body := post(t, ctx, baseURL+"/register", `{"email":"[EMAIL]"}`, "application/json", "")
	if status != http.StatusOK {
		return "", fmt.Errorf("register status = %d, want %d; body %q", status, http.StatusOK, body)
	}
	var payload struct {
		APIKey string `json:"apiKey"`
	}
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		return "", fmt.Errorf("decode register response: %w", err)
	}
	if payload.APIKey == "" {
		return "", errors.New("register response has no API key")
	}
	return payload.APIKey, nil
}

func post(t *testing.T, ctx context.Context, url, body, contentType, apiKey string) (int, string) {
	t.Helper()
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		t.Fatalf("create POST %s: %v", url, err)
	}
	request.Header.Set("Content-Type", contentType)
	if apiKey != "" {
		request.Header.Set("Authorization", "Bearer "+apiKey)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("POST %s: %v", url, err)
	}
	defer func() { _ = response.Body.Close() }()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("read POST %s: %v", url, err)
	}
	return response.StatusCode, string(responseBody)
}

func callTool(t *testing.T, ctx context.Context, session *mcp.ClientSession, name string, args map[string]any, output *resourceListOutput) {
	t.Helper()
	result, err := session.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		t.Fatalf("CallTool(%q): %v", name, err)
	}
	content, ok := result.StructuredContent.(map[string]any)
	if !ok {
		t.Fatalf("CallTool(%q) structured content has type %T, want map[string]any", name, result.StructuredContent)
	}
	encoded, err := json.Marshal(content)
	if err != nil {
		t.Fatalf("marshal CallTool(%q) result: %v", name, err)
	}
	if err := json.Unmarshal(encoded, output); err != nil {
		t.Fatalf("unmarshal CallTool(%q) result: %v", name, err)
	}
}

func waitForHealth(t *testing.T, ctx context.Context, url string) {
	t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	for {
		status, body, err := request(ctx, url)
		if err == nil && status == http.StatusOK && body == "ok\n" {
			return
		}
		if time.Now().After(deadline) {
			if err != nil {
				t.Fatalf("GET %s did not become healthy: %v", url, err)
			}
			t.Fatalf("GET %s did not become healthy: status = %d, body = %q", url, status, body)
		}
		select {
		case <-ctx.Done():
			t.Fatalf("GET %s did not become healthy: %v", url, ctx.Err())
		case <-time.After(250 * time.Millisecond):
		}
	}
}

func get(t *testing.T, ctx context.Context, url string) (int, string) {
	t.Helper()
	status, body, err := request(ctx, url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	return status, body
}

func request(ctx context.Context, url string) (int, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, "", fmt.Errorf("create GET %s: %w", url, err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, "", err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", fmt.Errorf("read GET %s: %w", url, err)
	}
	return resp.StatusCode, string(body), nil
}
