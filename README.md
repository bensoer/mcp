# mcp
Personal MCP Server with practices, preferences and validation tooling for AI Agents

# Prerequisits
- Golang 1.24+
- Makefile

# Setup
1. Clone repository and `cd` to project root
2. Run `make` to build the binary
3. Run the mcp server by executing `./bin/mcp`


# Ideas / Designs / Plans
- Create resource endpoints containing documentation of all standards, practices, preferences etc
- Create tool endpoints that allow LLM to send code samples and have it validated
against standards, practices and preferences - acting as a guardrail that it followed everything

# Developer Notes
- https://github.com/golang-standards/project-layout
- https://pkg.go.dev/github.com/modelcontextprotocol/go-sdk@v1.6.1/mcp#pkg-overview
- https://modelcontextprotocol.io/docs/getting-started/intro
