package tools

import (
	"fmt"
	"strings"

	"mcp/internal/models"
	"mcp/internal/utils"

	"github.com/adrg/frontmatter"
	"github.com/golang-nlp/stopwords"
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

// stringHasAllTerms reports whether every term appears in s, case-insensitively.
func stringHasAllTerms(s string, terms []string) bool {
	sLower := strings.ToLower(s)
	for _, term := range terms {
		if !strings.Contains(sLower, strings.ToLower(term)) {
			return false
		}
	}
	return true
}

// stringHasAnyTerms reports whether any term appears in s, case-insensitively.
func stringHasAnyTerms(s string, terms []string) bool {
	sLower := strings.ToLower(s)
	for _, term := range terms {
		if strings.Contains(sLower, strings.ToLower(term)) {
			return true
		}
	}
	return false
}

// removeStopWords drops English stopwords. Comparison is case-insensitive.
func removeStopWords(terms []string) []string {
	var filtered []string
	for _, term := range terms {
		if !stopwords.IsStopWord("en", term) {
			filtered = append(filtered, term)
		}
	}
	return filtered
}

// SearchMetadata searches assets for case-insensitive matches.
// A resource matches when any query term is in the URI, name, or description,
// or when every remaining query term is in the body. Stopwords are removed.
// Only resources whose URI starts with uriPrefix are returned.
func SearchMetadata(finder utils.AssetsFinder, query string, uriPrefix string) ([]ResourceMetadata, error) {
	assetPaths, err := finder.GetAllAssetPaths()
	if err != nil {
		return nil, fmt.Errorf("failed to list asset paths: %w", err)
	}

	queryTerms := removeStopWords(strings.Fields(strings.ToLower(query)))

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
		if stringHasAnyTerms(meta.URI, queryTerms) ||
			stringHasAnyTerms(meta.Name, queryTerms) ||
			stringHasAnyTerms(meta.Description, queryTerms) ||
			stringHasAllTerms(bodyLower, queryTerms) {
			matches = append(matches, ResourceMetadata{
				URI:         meta.URI,
				Name:        meta.Name,
				Description: meta.Description,
			})
		}
	}

	return matches, nil
}
