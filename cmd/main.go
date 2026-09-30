package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"mcp/internal"
	"mcp/internal/logger"
	"mcp/internal/utils"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

// newMCPHandler wires the single bootstrapped server behind a streamable HTTP
// transport on "/mcp" with a minimal health endpoint. It is shared by main
// (the -http path) and the package test.
func newMCPHandler(server *mcp.Server) http.Handler {
	mcpHandler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{
		SessionTimeout: 30 * time.Minute,
	})
	mux := http.NewServeMux()
	mux.Handle("/mcp", mcpHandler)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	return mux
}

func main() {
	// A local FlagSet avoids colliding with the global set that `go test` parses,
	// which a process-wide `-http` flag would break.
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	httpAddr := fs.String("http", "", "host:port for streamable HTTP transport; empty means stdio")
	if err := fs.Parse(os.Args[1:]); err != nil {
		if err != flag.ErrHelp {
			zap.S().Fatalf("failed to parse flags: %v", err)
		}
		return
	}

	// Whitespace-only is a fatal error; unset (empty) means stdio.
	if *httpAddr != "" && strings.TrimSpace(*httpAddr) == "" {
		logger.InitLogger(logger.PRODUCTION)
		zap.S().Fatalf("-http flag value is whitespace-only")
	}

	if *httpAddr == "" {
		logger.InitLogger(logger.DEVELOPMENT)
	} else {
		logger.InitLogger(logger.PRODUCTION)
	}
	defer zap.L().Sync() //nolint:errcheck

	//zap.S().Info("Starting MCP Server...")

	assetFolderRoot := os.Getenv("MCP_ASSET_ROOT")
	if assetFolderRoot == "" {
		assetFolderRoot = "assets"
	}

	server, err := internal.BootstrapServer(utils.NewAssetsFinder(assetFolderRoot))
	if err != nil {
		zap.S().Fatalf("Failed to bootstrap server: %v", err)
	}

	if *httpAddr == "" {
		if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
			zap.S().Fatal(err)
		}
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:              *httpAddr,
		Handler:           newMCPHandler(server),
		ReadHeaderTimeout: 10 * time.Second,
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			zap.S().Fatalf("http shutdown: %v", err)
		}
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			zap.S().Fatal(err)
		}
	}
}
