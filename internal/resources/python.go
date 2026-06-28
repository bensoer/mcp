package resources

import (
	"context"
	"mcp/internal/utils"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/zap"
)

type PythonResourceService interface {
	GetPythonLoggingResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error)
	GetPythonSyntaxResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error)
	GetPythonArchitectureResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error)
}

type PythonResourceServiceImpl struct {
	assetsFinder utils.AssetsFinder
}

func NewPythonResourceService(assetsFinder utils.AssetsFinder) PythonResourceService {
	return &PythonResourceServiceImpl{
		assetsFinder: assetsFinder,
	}
}

func (s *PythonResourceServiceImpl) GetPythonLoggingResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	uri := req.Params.URI

	contents, err := s.assetsFinder.GetAssetContents("python/logging.md")
	if err != nil {
		zap.S().Errorf("Failed to read asset contents for URI %s: %v", uri, err)
		return nil, mcp.ResourceNotFoundError(uri)
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{{URI: uri, Text: string(contents)}},
	}, nil
}

func (s *PythonResourceServiceImpl) GetPythonSyntaxResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	uri := req.Params.URI

	contents, err := s.assetsFinder.GetAssetContents("python/syntax.md")
	if err != nil {
		zap.S().Errorf("Failed to read asset contents for URI %s: %v", uri, err)
		return nil, mcp.ResourceNotFoundError(uri)
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{{URI: uri, Text: string(contents)}},
	}, nil
}

func (s *PythonResourceServiceImpl) GetPythonArchitectureResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	uri := req.Params.URI

	contents, err := s.assetsFinder.GetAssetContents("python/architecture.md")
	if err != nil {
		zap.S().Errorf("Failed to read asset contents for URI %s: %v", uri, err)
		return nil, mcp.ResourceNotFoundError(uri)
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{{URI: uri, Text: string(contents)}},
	}, nil
}
