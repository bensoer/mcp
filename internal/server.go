package internal

import (
	"mcp/internal/resources"
	"mcp/internal/utils"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func BoostrapServer() (*mcp.Server, error) {
	assetFolderRoot := "assets"
	aff := utils.NewAssetsFinder(&assetFolderRoot)

	pythonResourceService := resources.NewPythonResourceService(aff)

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "Ben's Standards, Practices, Linting & Expectations For AI Agents",
		Version: "1.0.0",
	}, nil)

	server.AddResource(
		&mcp.Resource{
			URI:         "standards://python/logging",
			Name:        "Python Logging Standards",
			Title:       "Python Logging Standards",
			MIMEType:    "text/markdown",
			Description: "Python Logging Standards, Formats, Syntax Expections and Examples for AI Agents",
		},
		pythonResourceService.GetPythonLoggingResource,
	)

	server.AddResource(
		&mcp.Resource{
			URI:         "standards://python/syntax",
			Name:        "Python Syntax Standards",
			Title:       "Python Syntax Standards",
			MIMEType:    "text/markdown",
			Description: "Standards for how python code should look. Naming conventions, and structure preferences",
		},
		pythonResourceService.GetPythonSyntaxResource,
	)

	server.AddResource(
		&mcp.Resource{
			URI:         "standards://python/architecture",
			Name:        "Python Architecture Standards",
			Title:       "Python Architecture Standards",
			MIMEType:    "text/markdown",
			Description: "Standards for the architectural design of Python applications, including module structure, design patterns, and best practices",
		},
		pythonResourceService.GetPythonArchitectureResource,
	)

	return server, nil
}
