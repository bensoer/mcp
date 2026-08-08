package models

import (
	"strings"
	"testing"

	"github.com/adrg/frontmatter"
)

func TestFrontMatter_ParseAllFields(t *testing.T) {
	input := `---
uri: standards://test/example
name: Test Example
description: A test resource
languages:
  - go
  - python
file_types:
  - "*.go"
  - "*.py"
priority: required
related_resources:
  - standards://test/other
  - standards://test/another
---
# Body content
`

	var fm FrontMatter
	_, err := frontmatter.Parse(strings.NewReader(input), &fm)
	if err != nil {
		t.Fatalf("frontmatter.Parse() error: %v", err)
	}

	if fm.URI != "standards://test/example" {
		t.Errorf("URI = %q, want %q", fm.URI, "standards://test/example")
	}
	if fm.Name != "Test Example" {
		t.Errorf("Name = %q, want %q", fm.Name, "Test Example")
	}
	if fm.Description != "A test resource" {
		t.Errorf("Description = %q, want %q", fm.Description, "A test resource")
	}
	if len(fm.Languages) != 2 || fm.Languages[0] != "go" || fm.Languages[1] != "python" {
		t.Errorf("Languages = %v, want [go python]", fm.Languages)
	}
	if len(fm.FileTypes) != 2 || fm.FileTypes[0] != "*.go" || fm.FileTypes[1] != "*.py" {
		t.Errorf("FileTypes = %v, want [*.go *.py]", fm.FileTypes)
	}
	if fm.Priority != "required" {
		t.Errorf("Priority = %q, want %q", fm.Priority, "required")
	}
	if len(fm.RelatedResources) != 2 || fm.RelatedResources[0] != "standards://test/other" || fm.RelatedResources[1] != "standards://test/another" {
		t.Errorf("RelatedResources = %v, want [standards://test/other standards://test/another]", fm.RelatedResources)
	}
}

func TestFrontMatter_CamelCaseKeysAreDropped(t *testing.T) {
	input := `---
relatedResources:
  - standards://test/camel
related_resources:
  - standards://test/snake
---
`

	var fm FrontMatter
	_, err := frontmatter.Parse(strings.NewReader(input), &fm)
	if err != nil {
		t.Fatalf("frontmatter.Parse() error: %v", err)
	}

	if len(fm.RelatedResources) != 1 {
		t.Fatalf("RelatedResources length = %d, want 1", len(fm.RelatedResources))
	}
	if fm.RelatedResources[0] != "standards://test/snake" {
		t.Errorf("RelatedResources[0] = %q, want %q", fm.RelatedResources[0], "standards://test/snake")
	}
}

func TestFrontMatter_MinimalFields(t *testing.T) {
	input := `---
uri: test://minimal
name: Minimal
---
`

	var fm FrontMatter
	_, err := frontmatter.Parse(strings.NewReader(input), &fm)
	if err != nil {
		t.Fatalf("frontmatter.Parse() error: %v", err)
	}

	if fm.URI != "test://minimal" {
		t.Errorf("URI = %q, want %q", fm.URI, "test://minimal")
	}
	if fm.Name != "Minimal" {
		t.Errorf("Name = %q, want %q", fm.Name, "Minimal")
	}
	if fm.Description != "" {
		t.Errorf("Description = %q, want empty", fm.Description)
	}
	if len(fm.Languages) != 0 {
		t.Errorf("Languages = %v, want empty", fm.Languages)
	}
}

func TestFrontMatter_RestReturnsBodyWithoutFrontmatter(t *testing.T) {
	input := `---
uri: test://f
name: F
---
# Heading

Body content here.
`

	var fm FrontMatter
	rest, err := frontmatter.Parse(strings.NewReader(input), &fm)
	if err != nil {
		t.Fatalf("frontmatter.Parse() error: %v", err)
	}

	if !strings.Contains(string(rest), "# Heading") {
		t.Error("rest does not contain body content")
	}
	if strings.Contains(string(rest), "---") {
		t.Error("rest should not contain frontmatter delimiters")
	}
}
