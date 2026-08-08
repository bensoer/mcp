package internal

import (
	"bytes"
	"context"
	"fmt"

	"mcp/internal/models"
	"mcp/internal/utils"

	"github.com/adrg/frontmatter"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func BootstrapServer(finder utils.AssetsFinder) (*mcp.Server, error) {
	assetPaths, err := finder.GetAllAssetPaths()
	if err != nil {
		return nil, err
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "Ben's Standards, Practices, Linting & Expectations For AI Agents",
		Version: "1.0.0",
	}, nil)

	for _, assetPath := range assetPaths {
		contents, err := finder.GetAssetContents(assetPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read asset %s: %w", assetPath, err)
		}

		var metaData models.FrontMatter
		_, err = frontmatter.Parse(bytes.NewReader(contents), &metaData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse frontmatter for asset %s: %w", assetPath, err)
		}

		if metaData.URI == "" {
			return nil, fmt.Errorf("asset %s is missing required frontmatter field: uri", assetPath)
		}

		if metaData.Name == "" {
			return nil, fmt.Errorf("asset %s is missing required frontmatter field: name", assetPath)
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
