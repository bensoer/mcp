package tools

import (
	"mcp/internal/utils"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterAll is the canonical list of tools. Adding a new tool:
//   1. Create the tool file with Name, Description, types, handler, and Register<Name>Bootstrap
//   2. Append Register<Name>Bootstrap to this slice
var RegisterAll = []func(*mcp.Server, utils.AssetsFinder){
	RegisterDiscoverStandardsBootstrap,
	RegisterDiscoverWorkflowsBootstrap,
	RegisterSearchStandardsBootstrap,
	RegisterSearchWorkflowsBootstrap,
}
