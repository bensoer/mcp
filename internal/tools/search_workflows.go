package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"mcp/internal/utils"
)

const (
	SearchWorkflowsName        = "search_workflows"
	SearchWorkflowsDescription = "Searches workflows (frontmatter + body) for case-insensitive partial matches of `query`. Returns metadata only ({uri, name, description}). Call resources/read on a matching URI to fetch the full document. The response itself does NOT contain document content."
)

type SearchWorkflowsInput struct {
	Query string `json:"query"`
}

type SearchWorkflowsOutput struct {
	Query     string             `json:"query"`
	Resources []ResourceMetadata `json:"resources"`
	Total     int                `json:"total"`
}

// searchWorkflows is an unexported helper that does the actual work. It is
// NOT the tool itself — the closure passed to mcp.AddTool below is. Keeping
// the helper unexported signals that nothing outside this package should
// call it directly.
func searchWorkflows(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	in SearchWorkflowsInput,
	finder utils.AssetsFinder,
) (*mcp.CallToolResult, SearchWorkflowsOutput, error) {
	matches, err := SearchMetadata(finder, in.Query, "workflows://")
	if err != nil {
		return nil, SearchWorkflowsOutput{}, err
	}
	return nil, SearchWorkflowsOutput{
		Query:     in.Query,
		Resources: matches,
		Total:     len(matches),
	}, nil
}

// RegisterSearchWorkflowsBootstrap wires the helper into the server. The
// closure captures `finder` and is what actually gets registered as the tool.
// The "Bootstrap" suffix makes the intent explicit: this function bootstraps
// the tool into the server.
func RegisterSearchWorkflowsBootstrap(server *mcp.Server, finder utils.AssetsFinder) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        SearchWorkflowsName,
		Description: SearchWorkflowsDescription,
	}, func(ctx context.Context, req *mcp.CallToolRequest, in SearchWorkflowsInput) (*mcp.CallToolResult, SearchWorkflowsOutput, error) {
		return searchWorkflows(ctx, req, in, finder)
	})
}
