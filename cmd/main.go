package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"mcp/internal"
	"mcp/internal/auth"
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
		Stateless: true,
	})
	mux := http.NewServeMux()
	mux.Handle("/mcp", mcpHandler)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	return mux
}

type serverConfig struct {
	httpAddr  string
	authPath  string
	rateLimit int
}

func parseServerConfig(args []string) (serverConfig, error) {
	// A local FlagSet avoids colliding with the global set that `go test` parses,
	// which a process-wide `-http` flag would break.
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	httpAddr := fs.String("http", "", "host:port for streamable HTTP transport; empty means stdio")
	authPath := fs.String("auth", "", "Bolt database path for HTTP API key authentication; empty means authentication is off")
	rateLimit := fs.Int("rate-limit", 0, "sustained HTTP requests per second per identity; zero uses defaults")
	if err := fs.Parse(args); err != nil {
		return serverConfig{}, err
	}

	authSet := false
	rateLimitSet := false
	fs.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "auth":
			authSet = true
		case "rate-limit":
			rateLimitSet = true
		}
	})

	config := serverConfig{
		httpAddr:  *httpAddr,
		authPath:  *authPath,
		rateLimit: *rateLimit,
	}
	// Whitespace-only is a fatal error; unset (empty) means stdio.
	if config.httpAddr != "" && strings.TrimSpace(config.httpAddr) == "" {
		return serverConfig{}, errors.New("-http flag value is whitespace-only")
	}
	if config.authPath != "" && strings.TrimSpace(config.authPath) == "" {
		return serverConfig{}, errors.New("-auth flag value is whitespace-only")
	}
	if config.rateLimit < 0 {
		return serverConfig{}, errors.New("-rate-limit must not be negative")
	}
	if config.httpAddr == "" {
		if authSet {
			return serverConfig{}, errors.New("-auth requires -http")
		}
		if rateLimitSet {
			return serverConfig{}, errors.New("-rate-limit requires -http")
		}
	}
	return config, nil
}

func main() {
	config, err := parseServerConfig(os.Args[1:])
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		logger.InitLogger(logger.PRODUCTION)
		zap.S().Fatalf("failed to parse flags: %v", err)
	}

	if config.httpAddr == "" {
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

	if config.httpAddr == "" {
		if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
			zap.S().Fatal(err)
		}
		return
	}

	var store *auth.Store
	if config.authPath != "" {
		store, err = auth.OpenDatabase(config.authPath)
		if err != nil {
			zap.S().Fatalf("Failed to open auth database: %v", err)
		}
	}
	closeStore := func() {
		if store == nil {
			return
		}
		if err := store.Close(); err != nil {
			zap.S().Errorf("Failed to close auth database: %v", err)
		}
		store = nil
	}
	defer closeStore()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	handler := newMCPHandler(server)
	if store != nil {
		handler = auth.NewHandler(handler, store, config.rateLimit, zap.S())
	}
	srv := &http.Server{
		Addr:              config.httpAddr,
		Handler:           handler,
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
			closeStore()
			zap.S().Fatalf("http shutdown: %v", err)
		}
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			closeStore()
			zap.S().Fatal(err)
		}
	}
	closeStore()
}
