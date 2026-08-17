package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"mcp/internal/utils"
)

const (
	SearchStandardsName        = "search_standards"
	SearchStandardsDescription = "Searches standards (frontmatter + body) for case-insensitive partial matches of `query`. Returns metadata only ({uri, name, description}). Call resources/read on a matching URI to fetch the full document. The response itself does NOT contain document content."
)

type SearchStandardsInput struct {
	Query string `json:"query"`
}

type SearchStandardsOutput struct {
	Query     string             `json:"query"`
	Resources []ResourceMetadata `json:"resources"`
	Total     int                `json:"total"`
}

// searchStandards is an unexported helper that does the actual work. It is
// NOT the tool itself — the closure passed to mcp.AddTool below is. Keeping
// the helper unexported signals that nothing outside this package should
// call it directly.
func searchStandards(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	in SearchStandardsInput,
	finder utils.AssetsFinder,
) (*mcp.CallToolResult, SearchStandardsOutput, error) {
	matches, err := SearchMetadata(finder, in.Query, "standards://")
	if err != nil {
		return nil, SearchStandardsOutput{}, err
	}
	return nil, SearchStandardsOutput{
		Query:     in.Query,
		Resources: matches,
		Total:     len(matches),
	}, nil
}

// RegisterSearchStandardsBootstrap wires the helper into the server. The
// closure captures `finder` and is what actually gets registered as the tool.
// The "Bootstrap" suffix makes the intent explicit: this function bootstraps
// the tool into the server.
func RegisterSearchStandardsBootstrap(server *mcp.Server, finder utils.AssetsFinder) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        SearchStandardsName,
		Description: SearchStandardsDescription,
	}, func(ctx context.Context, req *mcp.CallToolRequest, in SearchStandardsInput) (*mcp.CallToolResult, SearchStandardsOutput, error) {
		return searchStandards(ctx, req, in, finder)
	})
}
