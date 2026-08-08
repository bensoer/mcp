package main

import (
	"context"
	"os"
	"mcp/internal"
	"mcp/internal/logger"
	"mcp/internal/utils"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

func main() {

	logger.InitLogger(logger.DEVELOPMENT)
	defer zap.L().Sync()

	//zap.S().Info("Starting MCP Server...")

	assetFolderRoot := os.Getenv("MCP_ASSET_ROOT")
	if assetFolderRoot == "" {
		assetFolderRoot = "assets"
	}

	server, err := internal.BootstrapServer(utils.NewAssetsFinder(&assetFolderRoot))

	if err != nil {
		zap.S().Fatalf("Failed to bootstrap server: %v", err)
	}

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		zap.S().Fatal(err)
	}
}
