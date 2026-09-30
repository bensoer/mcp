package tools

import (
	"errors"
	"strings"
	"testing"
)

func TestStringHasAllTerms(t *testing.T) {
	tests := []struct {
		s     string
		terms []string
		want  bool
		name  string
	}{
		// Empty terms slice → true (vacuous truth)
		{`hello`, []string{}, true, "empty terms"},
		// All terms present (lowercase mix)
		{`hello world`, []string{"hello", "world"}, true, "all present"},
		{`Hello World`, []string{"hello", "world"}, true, "case insensitive"},
		// One term missing
		{`hello world`, []string{"hello", "earth"}, false, "one missing"},
		// Terms in different cases
		{`HeLLo WoRlD`, []string{"hello", "world"}, true, "mixed case"},
		// Empty string being searched
		{``, []string{"hello"}, false, "empty string"},
		{``, []string{}, true, "empty string and empty terms"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringHasAllTerms(tt.s, tt.terms); got != tt.want {
				t.Errorf("stringHasAllTerms(%q, %v) = %v, want %v", tt.s, tt.terms, got, tt.want)
			}
		})
	}
}

// TestStringHasAnyTerms tests the stringHasAnyTerms helper function.
func TestStringHasAnyTerms(t *testing.T) {
	tests := []struct {
		s     string
		terms []string
		want  bool
		name  string
	}{
		// Empty terms slice → false
		{`hello`, []string{}, false, "empty terms"},
		// One term matches
		{`hello world`, []string{"world"}, true, "one term matches"},
		{`hello world`, []string{"hello"}, true, "one term matches (first)"},
		// No terms match
		{`hello world`, []string{"foo", "bar"}, false, "no terms match"},
		// Multiple terms, one matches
		{`hello world`, []string{"foo", "world", "bar"}, true, "multiple terms, one matches"},
		// Case insensitivity
		{`Hello World`, []string{"hello"}, true, "case insensitive match"},
		{`Hello World`, []string{"WORLD"}, true, "case insensitive match (upper)"},
		// Empty string being searched
		{``, []string{"hello"}, false, "empty string"},
		{``, []string{}, false, "empty string and empty terms"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringHasAnyTerms(tt.s, tt.terms); got != tt.want {
				t.Errorf("stringHasAnyTerms(%q, %v) = %v, want %v", tt.s, tt.terms, got, tt.want)
			}
		})
	}
}

// TestRemoveStopWords tests the removeStopWords helper function.
func TestRemoveStopWords(t *testing.T) {
	tests := []struct {
		input []string
		want  []string
		name  string
	}{
		// Common English stopwords removed: "the", "is", "a", "of", "and", "to", "in"
		{[]string{"the", "is", "a", "of", "and", "to", "in"}, []string{}, "all stopwords"},
		// Non-stopwords retained: "git", "commit", "python", "workflow"
		{[]string{"git", "commit", "python", "workflow"}, []string{"git", "commit", "python", "workflow"}, "no stopwords"},
		// Mixed list
		{[]string{"the", "git", "is", "commit", "a", "python", "of", "workflow", "and", "to", "in"}, []string{"git", "commit", "python", "workflow"}, "mixed"},
		// Empty list
		{[]string{}, []string{}, "empty"},
		// Single stopword
		{[]string{"the"}, []string{}, "single stopword"},
		// All stopwords → empty result (already covered)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeStopWords(tt.input); !sliceEqual(got, tt.want) {
				t.Errorf("removeStopWords(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// sliceEqual is a helper to compare two string slices.
func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

// TestSearchMetadata_StopwordOnlyQuery tests that a query containing only stopwords
// returns all standards resources (since the term list is empty after stopword removal).
func TestSearchMetadata_StopwordOnlyQuery(t *testing.T) {
	mock := &mockAssetsFinder{
		paths: []string{
			"a.md",
			"b.md",
		},
		contents: map[string]string{
			"a.md": validAssetYAML("standards://test/a", "Resource A"),
			"b.md": validAssetYAML("standards://test/b", "Resource B"),
		},
	}

	// Query with only stopwords
	matches, err := SearchMetadata(mock, "the is a of and to in", "standards://")
	if err != nil {
		t.Fatalf("SearchMetadata() error: %v", err)
	}
	// After stopword removal, queryTerms is empty -> stringHasAllTerms returns true for any body
	// So both resources should match (assuming they are standards://)
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches for stopword-only query, got %d", len(matches))
	}
	// Verify both resources are returned
	foundA := false
	foundB := false
	for _, m := range matches {
		if m.URI == "standards://test/a" {
			foundA = true
		}
		if m.URI == "standards://test/b" {
			foundB = true
		}
	}
	if !foundA {
		t.Error("missing resource A")
	}
	if !foundB {
		t.Error("missing resource B")
	}
}

// TestSearchMetadata_BodyRequiresAllTerms tests that the body requires all terms
// (using stringHasAllTerms) while URI/Name/Description use any terms (stringHasAnyTerms).
func TestSearchMetadata_BodyRequiresAllTerms(t *testing.T) {
	// Create an asset where:
	//   URI: "standards://test/uri"
	//   Name: "Name"
	//   Description: "Description"
	//   Body: "term1 term2"  (missing term3)
	// Query: "term1 term2 term3"
	// Expect: match because URI/Name/Description have at least one term? Actually, let's design:
	//   We want the body to have only 2 of 3 terms, but the URI has the third term.
	//   Since URI uses ANY semantics, having one term in URI is enough for that field.
	//   But we also need the body to have ALL terms? No, the condition is:
	//   (URI has ANY term) OR (Name has ANY term) OR (Description has ANY term) OR (Body has ALL terms)
	//   So if the body does not have all terms, we rely on the other fields having at least one term.
	//
	// Let's set:
	//   URI: "standards://test/term3"   -> contains "term3"
	//   Name: "Name"
	//   Description: "Description"
	//   Body: "term1 term2"
	//   Query: "term1 term2 term3"
	//
	//   URI: has "term3" -> ANY term matches -> true
	//   So overall should match.
	mock := &mockAssetsFinder{
		paths: []string{"asset.md"},
		contents: map[string]string{
			"asset.md": validAssetYAML("standards://test/term3", "Name") + "\n\nterm1 term2\n",
		},
	}

	matches, err := SearchMetadata(mock, "term1 term2 term3", "standards://")
	if err != nil {
		t.Fatalf("SearchMetadata() error: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].URI != "standards://test/term3" {
		t.Errorf("unexpected URI: %q", matches[0].URI)
	}
}

// TestSearchMetadata_OnlyFrontmatterMatch tests that a query matching only in
// URI/Name/Description (with ANY semantics) returns the resource.
func TestSearchMetadata_OnlyFrontmatterMatch(t *testing.T) {
	// Body has none of the terms, but URI has one term.
	mock := &mockAssetsFinder{
		paths: []string{"asset.md"},
		contents: map[string]string{
			"asset.md": validAssetYAML("standards://test/hello", "Test") + "\n\nno match here\n",
		},
	}

	matches, err := SearchMetadata(mock, "hello", "standards://")
	if err != nil {
		t.Fatalf("SearchMetadata() error: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].URI != "standards://test/hello" {
		t.Errorf("unexpected URI: %q", matches[0].URI)
	}
}

// TestSearchMetadata_BodyMatchAllTerms tests that a query matching all terms in the body
// (with ALL semantics) returns the resource even if frontmatter has no matches.
func TestSearchMetadata_BodyMatchAllTerms(t *testing.T) {
	// Frontmatter has none of the terms, but body has all terms.
	mock := &mockAssetsFinder{
		paths: []string{"asset.md"},
		contents: map[string]string{
			"asset.md": validAssetYAML("standards://test/foo", "Test") + "\n\nhello world test\n",
		},
	}

	matches, err := SearchMetadata(mock, "hello world test", "standards://")
	if err != nil {
		t.Fatalf("SearchMetadata() error: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].URI != "standards://test/foo" {
		t.Errorf("unexpected URI: %q", matches[0].URI)
	}
}

// TestSearchMetadata_GetAssetContentsError tests that GetAssetContents errors are wrapped.
func TestSearchMetadata_GetAssetContentsError(t *testing.T) {
	mock := &mockAssetsFinder{
		paths: []string{"good.md", "bad.md"},
		contents: map[string]string{
			"good.md": validAssetYAML("standards://test/good", "Good Resource"),
			"bad.md":  validAssetYAML("standards://test/bad", "Bad Resource"),
		},
		loadErr: errors.New("boom"),
	}
	_, err := SearchMetadata(mock, "test", "standards://")
	if err == nil {
		t.Fatal("SearchMetadata() expected error from GetAssetContents, got nil")
	}
	if !strings.Contains(err.Error(), "boom") {
		t.Errorf("error = %q, want error containing 'boom'", err.Error())
	}
}

// TestSearchMetadata_FrontmatterParseError tests that invalid frontmatter returns a wrapped error.
func TestSearchMetadata_FrontmatterParseError(t *testing.T) {
	mock := &mockAssetsFinder{
		paths: []string{"asset.md"},
		contents: map[string]string{
			"asset.md": `---
uri: [not valid yaml
name: Test
---
body
`},
	}

	_, err := SearchMetadata(mock, "test", "standards://")
	if err == nil {
		t.Error("SearchMetadata() expected error from frontmatter parse, got nil")
	}
	if !strings.Contains(err.Error(), "failed to parse frontmatter for asset asset.md") {
		t.Errorf("error = %q, want error containing 'failed to parse frontmatter'", err)
	}
}
