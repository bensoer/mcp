package utils

import (
	"os"
	"path/filepath"
)

type AssetsFinder interface {
	GetAssetPath(assetName string) string
	GetAssetFolderRoot() string
	GetAssetContents(assetPathInsideFolder string) ([]byte, error)
}

type AssetsFinderImpl struct {
	assetFolderRoot *string
}

func NewAssetsFinder(assetFolderRoot *string) AssetsFinder {
	return &AssetsFinderImpl{
		assetFolderRoot: assetFolderRoot,
	}
}

func (a *AssetsFinderImpl) GetAssetPath(assetName string) string {
	return filepath.Join(a.GetAssetFolderRoot(), assetName)
}

func (a *AssetsFinderImpl) GetAssetFolderRoot() string {
	return *a.assetFolderRoot
}

func (a *AssetsFinderImpl) GetAssetContents(assetPathInsideFolder string) ([]byte, error) {
	return os.ReadFile(a.GetAssetPath(assetPathInsideFolder))
}
