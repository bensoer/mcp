package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"mcp/internal/utils"
)

const (
	DiscoverStandardsName        = "discover_standards"
	DiscoverStandardsDescription = "Returns all standards:// resources as {uri, name, description} metadata. Use resources/read on a specific URI to fetch the full document."
)

type DiscoverStandardsInput struct{} // empty — no parameters

type DiscoverStandardsOutput struct {
	Resources []ResourceMetadata `json:"resources"`
}

// discoverStandards is an unexported helper that does the actual work. It is
// NOT the tool itself — the closure passed to mcp.AddTool below is. Keeping
// the helper unexported signals that nothing outside this package should
// call it directly.
func discoverStandards(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	_ DiscoverStandardsInput,
	finder utils.AssetsFinder,
) (*mcp.CallToolResult, DiscoverStandardsOutput, error) {
	all, err := ListAllMetadata(finder)
	if err != nil {
		return nil, DiscoverStandardsOutput{}, err
	}
	return nil, DiscoverStandardsOutput{
		Resources: filterByPrefix(all, "standards://"),
	}, nil
}

// RegisterDiscoverStandardsBootstrap wires the helper into the server. The
// closure captures `finder` and is what actually gets registered as the tool.
// The "Bootstrap" suffix makes the intent explicit: this function bootstraps
// the tool into the server.
func RegisterDiscoverStandardsBootstrap(server *mcp.Server, finder utils.AssetsFinder) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        DiscoverStandardsName,
		Description: DiscoverStandardsDescription,
	}, func(ctx context.Context, req *mcp.CallToolRequest, in DiscoverStandardsInput) (*mcp.CallToolResult, DiscoverStandardsOutput, error) {
		return discoverStandards(ctx, req, in, finder)
	})
}
