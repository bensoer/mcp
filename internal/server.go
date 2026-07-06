package internal

import (
	"bytes"
	"context"
	"mcp/internal/models"
	"mcp/internal/utils"

	"github.com/adrg/frontmatter"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func BoostrapServer() (*mcp.Server, error) {
	assetFolderRoot := "assets"
	aff := utils.NewAssetsFinder(&assetFolderRoot)
	assetPaths, err := aff.GetAllAssetPaths()
	if err != nil {
		return nil, err
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "Ben's Standards, Practices, Linting & Expectations For AI Agents",
		Version: "1.0.0",
	}, nil)

	for _, assetPath := range assetPaths {
		contents, err := aff.GetAssetContents(assetPath)
		if err != nil {
			return nil, err
		}

		var metaData models.FrontMatter
		_, err = frontmatter.Parse(bytes.NewReader(contents), &metaData)
		if err != nil {
			return nil, err
		}

		server.AddResource(
			&mcp.Resource{
				URI:         metaData.URI,
				Name:        metaData.Name,
				Title:       metaData.Name,
				MIMEType:    "text/markdown",
				Description: metaData.Description,
				Size:        int64(len(contents)),
			},
			func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
				return &mcp.ReadResourceResult{
					Contents: []*mcp.ResourceContents{
						{
							URI:      metaData.URI,
							Text:     string(contents),
							MIMEType: "text/markdown",
						},
					},
				}, nil
			},
		)
	}

	return server, nil
}

// 	pythonResourceService := resources.NewPythonResourceService(aff)

// 	server.AddResource(
// 		&mcp.Resource{
// 			URI:         "standards://python/logging",
// 			Name:        "Python Logging Standards",
// 			Title:       "Python Logging Standards",
// 			MIMEType:    "text/markdown",
// 			Description: "Python Logging Standards, Formats, Syntax Expections and Examples for AI Agents",
// 		},
// 		pythonResourceService.GetPythonLoggingResource,
// 	)

// 	server.AddResource(
// 		&mcp.Resource{
// 			URI:         "standards://python/syntax",
// 			Name:        "Python Syntax Standards",
// 			Title:       "Python Syntax Standards",
// 			MIMEType:    "text/markdown",
// 			Description: "Standards for how python code should look. Naming conventions, and structure preferences",
// 		},
// 		pythonResourceService.GetPythonSyntaxResource,
// 	)

// 	server.AddResource(
// 		&mcp.Resource{
// 			URI:         "standards://python/architecture",
// 			Name:        "Python Architecture Standards",
// 			Title:       "Python Architecture Standards",
// 			MIMEType:    "text/markdown",
// 			Description: "Standards for the architectural design of Python applications, including module structure, design patterns, and best practices",
// 		},
// 		pythonResourceService.GetPythonArchitectureResource,
// 	)

// 	return server, nil
// }
