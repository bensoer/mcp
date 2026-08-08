package utils

import (
	"os"
	"path/filepath"
	"strings"
)

type AssetsFinder interface {
	GetAssetPath(assetName string) string
	GetAssetFolderRoot() string
	GetAssetContents(assetPathInsideFolder string) ([]byte, error)
	GetAllAssetPaths() ([]string, error)
}

type AssetsFinderImpl struct {
	assetFolderRoot string
}

func NewAssetsFinder(assetFolderRoot string) AssetsFinder {
	return &AssetsFinderImpl{
		assetFolderRoot: assetFolderRoot,
	}
}

func (a *AssetsFinderImpl) GetAssetPath(assetName string) string {
	return filepath.Join(a.GetAssetFolderRoot(), assetName)
}

func (a *AssetsFinderImpl) GetAssetFolderRoot() string {
	return a.assetFolderRoot
}

func (a *AssetsFinderImpl) GetAssetContents(assetPathInsideFolder string) ([]byte, error) {
	return os.ReadFile(a.GetAssetPath(assetPathInsideFolder))
}

func (a *AssetsFinderImpl) GetAllAssetPaths() ([]string, error) {
	var files []string

	err := filepath.WalkDir(a.GetAssetFolderRoot(), func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}

		rel, err := filepath.Rel(a.GetAssetFolderRoot(), path)
		if err != nil {
			return err
		}

		files = append(files, rel)
		return nil
	})

	return files, err
}
