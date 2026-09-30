package tools

import (
	"errors"
	"fmt"
	"testing"
)

// mockAssetsFinder is a test double for utils.AssetsFinder.
type mockAssetsFinder struct {
	paths    []string
	contents map[string]string
	loadErr  error
	errPath  string
	errMsg   string
}

// GetAssetPath returns the asset name as the path.
func (m *mockAssetsFinder) GetAssetPath(assetName string) string {
	return assetName
}

// GetAssetFolderRoot returns a mock folder root.
func (m *mockAssetsFinder) GetAssetFolderRoot() string {
	return "mock://"
}

// GetAssetContents returns the content for the given asset path.
// If loadErr is set, it returns that error.
// If errPath is set and matches assetPathInsideFolder, it returns errMsg.
// If the assetPathInsideFolder is not found in contents, it returns a not-found error.
func (m *mockAssetsFinder) GetAssetContents(assetPathInsideFolder string) ([]byte, error) {
	if m.loadErr != nil {
		return nil, m.loadErr
	}
	if m.errPath != "" && assetPathInsideFolder == m.errPath {
		return nil, errors.New(m.errMsg)
	}
	c, ok := m.contents[assetPathInsideFolder]
	if !ok {
		return nil, fmt.Errorf("mock: file not found: %s", assetPathInsideFolder)
	}
	return []byte(c), nil
}

// GetAllAssetPaths returns the list of asset paths.
// If loadErr is set, it returns that error.
func (m *mockAssetsFinder) GetAllAssetPaths() ([]string, error) {
	if m.loadErr != nil {
		return nil, m.loadErr
	}
	return m.paths, nil
}

// validAssetYAML returns a valid YAML frontmatter string for an asset.
func validAssetYAML(uri, name string) string {
	return fmt.Sprintf(`---
uri: %s
name: %s
description: Test resource
languages:
  - all
file_types:
  - "*.*"
priority: required
related_resources:
  - standards://git/commit-messages
---
# %s

Content for %s.
`, uri, name, name, name)
}

// TestListAllMetadata tests the ListAllMetadata helper function.
func TestListAllMetadata(t *testing.T) {
	mock := &mockAssetsFinder{
		paths: []string{"a.md", "b.md"},
		contents: map[string]string{
			"a.md": validAssetYAML("standards://test/a", "Resource A"),
			"b.md": validAssetYAML("workflows://test/b", "Resource B"),
		},
	}

	metadata, err := ListAllMetadata(mock)
	if err != nil {
		t.Fatalf("ListAllMetadata() error: %v", err)
	}
	if len(metadata) != 2 {
		t.Fatalf("expected 2 metadata items, got %d", len(metadata))
	}
	// Check first item
	if metadata[0].URI != "standards://test/a" {
		t.Errorf("first URI = %q, want %q", metadata[0].URI, "standards://test/a")
	}
	if metadata[0].Name != "Resource A" {
		t.Errorf("first name = %q, want %q", metadata[0].Name, "Resource A")
	}
	// Check second item
	if metadata[1].URI != "workflows://test/b" {
		t.Errorf("second URI = %q, want %q", metadata[1].URI, "workflows://test/b")
	}
	if metadata[1].Name != "Resource B" {
		t.Errorf("second name = %q, want %q", metadata[1].Name, "Resource B")
	}
}

// TestListAllMetadata_Error tests error handling in ListAllMetadata.
func TestListAllMetadata_Error(t *testing.T) {
	mock := &mockAssetsFinder{
		loadErr: errors.New("disk failure"),
	}

	_, err := ListAllMetadata(mock)
	if err == nil {
		t.Error("ListAllMetadata() expected error from GetAllAssetPaths, got nil")
	}
	if err.Error() != "failed to list asset paths: disk failure" {
		t.Errorf("error = %q, want error containing 'failed to list asset paths'", err)
	}
}

// TestFilterByPrefix tests the filterByPrefix helper function.
func TestFilterByPrefix(t *testing.T) {
	items := []ResourceMetadata{
		{URI: "standards://git/commit", Name: "Git Commit", Description: "Commit message standards"},
		{URI: "standards://git/branch", Name: "Git Branch", Description: "Branching standards"},
		{URI: "workflows://test/foo", Name: "Test Workflow", Description: "A test workflow"},
		{URI: "standards://python/import", Name: "Python Import", Description: "Import standards"},
	}

	// Filter for standards://git
	filtered := filterByPrefix(items, "standards://git")
	if len(filtered) != 2 {
		t.Fatalf("expected 2 items after filtering, got %d", len(filtered))
	}
	if filtered[0].URI != "standards://git/commit" {
		t.Errorf("first filtered URI = %q, want %q", filtered[0].URI, "standards://git/commit")
	}
	if filtered[1].URI != "standards://git/branch" {
		t.Errorf("second filtered URI = %q, want %q", filtered[1].URI, "standards://git/branch")
	}

	// Filter for workflows://
	filtered = filterByPrefix(items, "workflows://")
	if len(filtered) != 1 {
		t.Fatalf("expected 1 item after filtering workflows, got %d", len(filtered))
	}
	if filtered[0].URI != "workflows://test/foo" {
		t.Errorf("workflow URI = %q, want %q", filtered[0].URI, "workflows://test/foo")
	}

	// Filter for non-existent prefix - should return empty (not error)
	filtered = filterByPrefix(items, "nonexistent://")
	if len(filtered) != 0 {
		t.Fatalf("expected 0 items after filtering for non-existent prefix, got %d", len(filtered))
	}
}

// TestSearchMetadata tests the SearchMetadata helper function.
func TestSearchMetadata(t *testing.T) {
	mock := &mockAssetsFinder{
		paths: []string{
			"standards.md",
			"workflow.md",
			"other.md",
		},
		contents: map[string]string{
			"standards.md": validAssetYAML("standards://test/alpha", "Alpha Resource"),
			"workflow.md":  validAssetYAML("workflows://test/beta", "Beta Workflow"),
			"other.md":     validAssetYAML("standards://test/gamma", "Gamma Resource"),
		},
	}

	// Search for "alpha" (should match in URI and name)
	matches, err := SearchMetadata(mock, "alpha", "standards://")
	if err != nil {
		t.Fatalf("SearchMetadata() error: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("expected 1 match for 'alpha', got %d", len(matches))
	}
	if matches[0].URI != "standards://test/alpha" {
		t.Errorf("match URI = %q, want %q", matches[0].URI, "standards://test/alpha")
	}
	if matches[0].Name != "Alpha Resource" {
		t.Errorf("match name = %q, want %q", matches[0].Name, "Alpha Resource")
	}

	// Search for "TEST" (case-insensitive, should match all three because description has "Test resource")
	matches, err = SearchMetadata(mock, "TEST", "standards://")
	if err != nil {
		t.Fatalf("SearchMetadata() error: %v", err)
	}
	// All three have "Test resource" in description, so all should match
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches for 'TEST' (case-insensitive), got %d", len(matches))
	}

	// Search for "content" in body (should match all three because body has "Content for X")
	matches, err = SearchMetadata(mock, "content", "standards://")
	if err != nil {
		t.Fatalf("SearchMetadata() error: %v", err)
	}
	if len(matches) != 2 {
		t.Fatalf("expected 2 standards matches for 'content' in body, got %d", len(matches))
	}

	// Search for "delta" (no matches)
	matches, err = SearchMetadata(mock, "delta", "standards://")
	if err != nil {
		t.Fatalf("SearchMetadata() error: %v", err)
	}
	if len(matches) != 0 {
		t.Fatalf("expected 0 matches for 'delta', got %d", len(matches))
	}

	// Empty query should match everything (because empty string is contained in any string)
	matches, err = SearchMetadata(mock, "", "standards://")
	if err != nil {
		t.Fatalf("SearchMetadata() error: %v", err)
	}
	if len(matches) != 2 {
		t.Fatalf("expected 2 standards matches for empty query, got %d", len(matches))
	}
}

// TestSearchMetadata_Error tests error handling in SearchMetadata.
func TestSearchMetadata_Error(t *testing.T) {
	mock := &mockAssetsFinder{
		loadErr: errors.New("disk failure"),
	}

	_, err := SearchMetadata(mock, "test", "standards://")
	if err == nil {
		t.Error("SearchMetadata() expected error from GetAllAssetPaths, got nil")
	}
	if err.Error() != "failed to list asset paths: disk failure" {
		t.Errorf("error = %q, want error containing 'failed to list asset paths'", err)
	}
}

// TestStringHasAllTerms tests the stringHasAllTerms helper function.
