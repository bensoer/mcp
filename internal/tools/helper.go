package tools

import (
	"fmt"
	"strings"

	"mcp/internal/models"
	"mcp/internal/utils"

	"github.com/adrg/frontmatter"
)

// ListAllMetadata walks all assets and returns metadata for every one.
// Returns all resources (both standards:// and workflows://).
func ListAllMetadata(finder utils.AssetsFinder) ([]ResourceMetadata, error) {
	assetPaths, err := finder.GetAllAssetPaths()
	if err != nil {
		return nil, fmt.Errorf("failed to list asset paths: %w", err)
	}

	var resources []ResourceMetadata
	for _, assetPath := range assetPaths {
		contents, err := finder.GetAssetContents(assetPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read asset %s: %w", assetPath, err)
		}

		var meta models.FrontMatter
		if _, err = frontmatter.Parse(strings.NewReader(string(contents)), &meta); err != nil {
			return nil, fmt.Errorf("failed to parse frontmatter for asset %s: %w", assetPath, err)
		}

		resources = append(resources, ResourceMetadata{
			URI:         meta.URI,
			Name:        meta.Name,
			Description: meta.Description,
		})
	}

	return resources, nil
}

// filterByPrefix filters resources to only those whose URI starts with the given prefix.
func filterByPrefix(resources []ResourceMetadata, prefix string) []ResourceMetadata {
	var filtered []ResourceMetadata
	for _, r := range resources {
		if strings.HasPrefix(r.URI, prefix) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

// SearchMetadata searches all assets for case-insensitive partial matches.
// Searches URI, Name, Description AND the document body.
// Only returns resources whose URI starts with uriPrefix (e.g., "standards://").
func SearchMetadata(finder utils.AssetsFinder, query string, uriPrefix string) ([]ResourceMetadata, error) {
	assetPaths, err := finder.GetAllAssetPaths()
	if err != nil {
		return nil, fmt.Errorf("failed to list asset paths: %w", err)
	}

	queryLower := strings.ToLower(query)

	var matches []ResourceMetadata
	for _, assetPath := range assetPaths {
		contents, err := finder.GetAssetContents(assetPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read asset %s: %w", assetPath, err)
		}

		var meta models.FrontMatter
		body, err := frontmatter.Parse(strings.NewReader(string(contents)), &meta)
		if err != nil {
			return nil, fmt.Errorf("failed to parse frontmatter for asset %s: %w", assetPath, err)
		}

		if !strings.HasPrefix(meta.URI, uriPrefix) {
			continue
		}

		bodyLower := strings.ToLower(string(body))
		if strings.Contains(strings.ToLower(meta.URI), queryLower) ||
			strings.Contains(strings.ToLower(meta.Name), queryLower) ||
			strings.Contains(strings.ToLower(meta.Description), queryLower) ||
			strings.Contains(bodyLower, queryLower) {
			matches = append(matches, ResourceMetadata{
				URI:         meta.URI,
				Name:        meta.Name,
				Description: meta.Description,
			})
		}
	}

	return matches, nil
}
