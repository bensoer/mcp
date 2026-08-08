package utils

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestAssetsFinderImpl_GetAssetFolderRoot(t *testing.T) {
	root := "testdata"
	af := NewAssetsFinder(&root)

	got := af.GetAssetFolderRoot()
	if got != root {
		t.Errorf("GetAssetFolderRoot() = %q, want %q", got, root)
	}
}

func TestAssetsFinderImpl_GetAssetPath(t *testing.T) {
	root := "testdata"
	af := NewAssetsFinder(&root)

	got := af.GetAssetPath("foo/bar.md")
	want := filepath.Join("testdata", "foo/bar.md")
	if got != want {
		t.Errorf("GetAssetPath() = %q, want %q", got, want)
	}
}

func TestAssetsFinderImpl_GetAssetContents(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.md")
	if err := os.WriteFile(file, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	af := NewAssetsFinder(&dir)
	got, err := af.GetAssetContents("test.md")
	if err != nil {
		t.Fatalf("GetAssetContents() error: %v", err)
	}
	if string(got) != "hello" {
		t.Errorf("GetAssetContents() = %q, want %q", string(got), "hello")
	}
}

func TestAssetsFinderImpl_GetAssetContents_NotFound(t *testing.T) {
	dir := t.TempDir()
	af := NewAssetsFinder(&dir)

	_, err := af.GetAssetContents("nonexistent.md")
	if err == nil {
		t.Error("GetAssetContents() expected error for nonexistent file, got nil")
	}
}

func TestAssetsFinderImpl_GetAllAssetPaths(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "sub"), 0o755)
	os.WriteFile(filepath.Join(dir, "a.md"), []byte("a"), 0o644)
	os.WriteFile(filepath.Join(dir, "sub", "b.md"), []byte("b"), 0o644)

	af := NewAssetsFinder(&dir)
	got, err := af.GetAllAssetPaths()
	if err != nil {
		t.Fatalf("GetAllAssetPaths() error: %v", err)
	}

	sort.Strings(got)
	want := []string{"a.md", filepath.Join("sub", "b.md")}
	sort.Strings(want)

	if len(got) != len(want) {
		t.Fatalf("GetAllAssetPaths() returned %d paths, want %d: %v", len(got), len(want), got)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("GetAllAssetPaths()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestAssetsFinderImpl_GetAllAssetPaths_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	af := NewAssetsFinder(&dir)

	got, err := af.GetAllAssetPaths()
	if err != nil {
		t.Fatalf("GetAllAssetPaths() error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("GetAllAssetPaths() = %v, want empty slice", got)
	}
}

func TestAssetsFinderImpl_GetAllAssetPaths_NonexistentRoot(t *testing.T) {
	root := "nonexistent_dir_abc123"
	af := NewAssetsFinder(&root)

	_, err := af.GetAllAssetPaths()
	if err == nil {
		t.Error("GetAllAssetPaths() expected error for nonexistent root, got nil")
	}
}