# Project Context — mcp

## Purpose

`mcp` is a personal Model Context Protocol server written in Go. It serves engineering standards, workflows, and search tools to AI coding agents.

## What is in this repo

1. **Server (Go 1.26)** — `cmd/`, `internal/`. One binary: `bin/mcp`. Speaks stdio JSON-RPC.
2. **Payload (`assets/`)** — markdown documents registered at startup as `standards://` and `workflows://` resources. This is data the server serves. It is not instructions for building this repo.
3. **Agent guidance** — `AGENTS.md` is the index. Rules, skills, and commands for maintaining this repository live in `.agents/`.

## Self-reference

This server is often connected to the agent editing it (`ben-mcp-server`). A broken build or bad asset frontmatter can break later standards lookups.
