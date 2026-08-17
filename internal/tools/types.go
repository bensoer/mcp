package tools

// ResourceMetadata represents the frontmatter metadata of an asset.
// This is what gets returned by discover_* and search_* tools.
type ResourceMetadata struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
