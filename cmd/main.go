package main

import (
	"context"
	"mcp/internal"
	"mcp/internal/logger"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

func main() {

	logger.InitLogger(logger.DEVELOPMENT)
	defer zap.L().Sync()

	//zap.S().Info("Starting MCP Server...")

	server, err := internal.BoostrapServer()

	if err != nil {
		zap.S().Fatalf("Failed to bootstrap server: %v", err)
	}

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		zap.S().Fatal(err)
	}
}
