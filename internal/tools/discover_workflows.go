package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"mcp/internal/utils"
)

const (
	DiscoverWorkflowsName        = "discover_workflows"
	DiscoverWorkflowsDescription = "Returns all workflows:// resources as {uri, name, description} metadata. Use resources/read on a specific URI to fetch the full document."
)

type DiscoverWorkflowsInput struct{} // empty — no parameters

type DiscoverWorkflowsOutput struct {
	Resources []ResourceMetadata `json:"resources"`
}

// discoverWorkflows is an unexported helper that does the actual work. It is
// NOT the tool itself — the closure passed to mcp.AddTool below is. Keeping
// the helper unexported signals that nothing outside this package should
// call it directly.
func discoverWorkflows(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	_ DiscoverWorkflowsInput,
	finder utils.AssetsFinder,
) (*mcp.CallToolResult, DiscoverWorkflowsOutput, error) {
	all, err := ListAllMetadata(finder)
	if err != nil {
		return nil, DiscoverWorkflowsOutput{}, err
	}
	return nil, DiscoverWorkflowsOutput{
		Resources: filterByPrefix(all, "workflows://"),
	}, nil
}

// RegisterDiscoverWorkflowsBootstrap wires the helper into the server. The
// closure captures `finder` and is what actually gets registered as the tool.
// The "Bootstrap" suffix makes the intent explicit: this function bootstraps
// the tool into the server.
func RegisterDiscoverWorkflowsBootstrap(server *mcp.Server, finder utils.AssetsFinder) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        DiscoverWorkflowsName,
		Description: DiscoverWorkflowsDescription,
	}, func(ctx context.Context, req *mcp.CallToolRequest, in DiscoverWorkflowsInput) (*mcp.CallToolResult, DiscoverWorkflowsOutput, error) {
		return discoverWorkflows(ctx, req, in, finder)
	})
}
