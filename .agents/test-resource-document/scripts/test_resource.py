#!/usr/bin/env python3
# /// script
# ///
"""
Test harness: verify a resource document is registered and served by the mcp server.

List all registered resources, or read the content of a specific resource URI.
All structured output goes to stdout (JSON); diagnostics go to stderr.

Usage:
    python3 scripts/test_resource.py                  # list all resources
    python3 scripts/test_resource.py --read <URI>     # read one resource
    python3 scripts/test_resource.py --build          # rebuild binary first, then list

Examples:
    python3 scripts/test_resource.py
    python3 scripts/test_resource.py --read standards://git/pull-requests
    python3 scripts/test_resource.py --read standards://git/pull-requests --build
"""

import argparse
import json
import os
import select
import subprocess
import sys
import time
from pathlib import Path


def build_binary(project_root: str) -> str:
    """Rebuild the mcp binary. Returns the binary path on success."""
    binary = os.path.join(project_root, "bin", "mcp")
    result = subprocess.run(
        ["go", "build", "-o", "bin/mcp", "./cmd"],
        cwd=project_root,
        capture_output=True,
        text=True,
    )
    if result.returncode != 0:
        print(
            "Error: build failed. See stderr from 'go build' below.\n"
            f"stdout:\n{result.stdout}\n"
            f"stderr:\n{result.stderr}",
            file=sys.stderr,
        )
        sys.exit(1)
    print(f"Built binary: {binary}", file=sys.stderr)
    return binary


def run_server(binary: str, project_root: str, timeout: float) -> dict:
    """Start the mcp server over stdio, send MCP handshake + a batch of JSON-RPC
    requests, and return the parsed responses.

    The server reads stdin to EOF then exits, so we write all requests before
    closing the pipe.
    """
    r, w = os.pipe()
    proc = subprocess.Popen(
        [binary],
        stdin=r,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        cwd=project_root,
    )
    os.close(r)

    requests = [
        {
            "jsonrpc": "2.0",
            "id": 1,
            "method": "initialize",
            "params": {
                "protocolVersion": "2024-11-05",
                "capabilities": {"resources": {}},
            },
        },
        {"jsonrpc": "2.0", "method": "notifications/initialized", "params": {}},
        {"jsonrpc": "2.0", "id": 2, "method": "resources/list", "params": {}},
    ]

    for req in requests:
        os.write(w, (json.dumps(req) + "\n").encode())

    # Wait before closing stdin so the server has time to process the buffered
    # requests. Closing stdin triggers EOF which causes the stdio transport to
    # shut down immediately — if done too fast the server exits before
    # responding.
    time.sleep(timeout)
    os.close(w)

    # Read stdout via a select loop (the server exits when stdin hits EOF,
    # logging "server is closing: EOF" to stderr — that is EXPECTED).
    stdout_bytes = b""
    while select.select([proc.stdout], [], [], 0.5)[0]:
        chunk = os.read(proc.stdout.fileno(), 65536)
        if not chunk:
            break
        stdout_bytes += chunk

    # Give the process a moment to fully exit, then reap.
    try:
        proc.wait(timeout=3)
    except subprocess.TimeoutExpired:
        proc.terminate()
        proc.wait(timeout=3)

    stderr_text = proc.stderr.read().decode(errors="replace") if proc.stderr else ""
    stdout_text = stdout_bytes.decode(errors="replace")

    messages = []
    for line in stdout_text.splitlines():
        line = line.strip()
        if not line:
            continue
        try:
            messages.append(json.loads(line))
        except json.JSONDecodeError:
            pass

    # "server is closing: EOF" is the EXPECTED shutdown when stdin closes.
    # Only treat FATAL as a real error if it's a bootstrap failure (e.g.
    # missing assets dir, bad frontmatter) — which means no responses were
    # produced. If we got valid JSON-RPC messages, the server worked fine.
    real_fatal = [
        ln for ln in stderr_text.splitlines()
        if "FATAL" in ln and "server is closing" not in ln
    ]
    if real_fatal and not messages:
        print(
            f"Error: server failed to start:\n  {real_fatal[0]}",
            file=sys.stderr,
        )
        print(
            "Likely causes: missing assets/ directory, malformed frontmatter "
            "(empty uri/name), or invalid YAML.",
            file=sys.stderr,
        )
        sys.exit(1)

    return {"messages": messages, "stderr": stderr_text}


def handle_list(args) -> None:
    """Send resources/list and print the resource catalog as JSON."""
    result = run_server(args.binary, args.project_root, args.timeout)
    resources = []
    for msg in result["messages"]:
        if msg.get("id") == 1 and "result" in msg:
            # initialize response — can skip
            pass
        if "result" in msg and "resources" in msg["result"]:
            resources = msg["result"]["resources"]

    print(
        json.dumps(
            {
                "count": len(resources),
                "resources": [
                    {
                        "uri": r.get("uri", ""),
                        "name": r.get("name", ""),
                        "description": r.get("description", ""),
                    }
                    for r in resources
                ],
            },
            indent=2,
        )
    )

    # Diagnostics to stderr
    uri_list = [r.get("uri", "") for r in resources]
    if args.read_target and args.read_target not in uri_list:
        print(
            f"Warning: requested URI '{args.read_target}' not found in resources/list.",
            file=sys.stderr,
        )
    print(f"Found {len(resources)} resource(s).", file=sys.stderr)


def handle_read(args) -> None:
    """Send resources/read for a specific URI and print the content as JSON."""
    uri = args.read_target
    if not uri:
        print(
            "Error: --read requires a URI argument.\n"
            "Usage: test_resource.py --read standards://git/pull-requests",
            file=sys.stderr,
        )
        sys.exit(2)

    binary = args.binary
    project_root = args.project_root

    r, w = os.pipe()
    proc = subprocess.Popen(
        [binary],
        stdin=r,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        cwd=project_root,
    )
    os.close(r)

    batch = [
        {
            "jsonrpc": "2.0",
            "id": 1,
            "method": "initialize",
            "params": {
                "protocolVersion": "2024-11-05",
                "capabilities": {"resources": {}},
            },
        },
        {"jsonrpc": "2.0", "method": "notifications/initialized", "params": {}},
        {"jsonrpc": "2.0", "id": 2, "method": "resources/read", "params": {"uri": uri}},
    ]
    for req in batch:
        os.write(w, (json.dumps(req) + "\n").encode())

    # Wait before closing stdin so the server has time to process requests
    # before EOF shuts it down.
    time.sleep(args.timeout)
    os.close(w)

    # Read stdout via a select loop (the server exits when stdin hits EOF,
    # logging "server is closing: EOF" to stderr — that is EXPECTED).
    stdout_bytes = b""
    while select.select([proc.stdout], [], [], 0.5)[0]:
        chunk = os.read(proc.stdout.fileno(), 65536)
        if not chunk:
            break
        stdout_bytes += chunk

    try:
        proc.wait(timeout=3)
    except subprocess.TimeoutExpired:
        proc.terminate()
        proc.wait(timeout=3)

    stderr_text = proc.stderr.read().decode(errors="replace") if proc.stderr else ""
    stdout_text = stdout_bytes.decode(errors="replace")

    # Parse stdout first. "server is closing: EOF" is the expected normal
    # shutdown when we close stdin — only treat other FATALs (bootstrap
    # failures) as real errors when no responses were produced.
    real_fatal = [
        ln for ln in stderr_text.splitlines()
        if "FATAL" in ln and "server is closing" not in ln
    ]

    text_content = None
    mime_type = None
    for line in stdout_text.splitlines():
        line = line.strip()
        if not line:
            continue
        try:
            msg = json.loads(line)
        except json.JSONDecodeError:
            continue
        if msg.get("id") == 2 and "result" in msg:
            contents = msg["result"].get("contents", [])
            if contents:
                text_content = contents[0].get("text", "")
                mime_type = contents[0].get("mimeType", "")

    if real_fatal and text_content is None:
        print(
            f"Error: server failed to start:\n  {real_fatal[0]}",
            file=sys.stderr,
        )
        print(
            "Likely causes: malformed frontmatter, invalid YAML, or a bad URI.",
            file=sys.stderr,
        )
        sys.exit(1)

    if text_content is None:
        print(
            f"Error: no content returned for URI '{uri}'.\n"
            "Possible causes: the URI does not match any registered resource, "
            "or the server has not been rebuilt since the document was added.",
            file=sys.stderr,
        )
        sys.exit(1)

    # Check that the frontmatter (if present in content) was parsed without issues
    print(
        json.dumps(
            {
                "uri": uri,
                "mime_type": mime_type,
                "size": len(text_content),
                "frontmatter_present": text_content.lstrip().startswith("---"),
            },
            indent=2,
        )
    )
    # Print the full text content to stderr so stdout stays machine-readable JSON
    print("--- content (stderr) ---", file=sys.stderr)
    print(text_content, file=sys.stderr)


def main():
    parser = argparse.ArgumentParser(
        description=(
            "Verify a resource document is discovered and served by the mcp MCP server. "
            "With no flags, lists all resources. Use --read <URI> to fetch one."
        ),
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=(
            "Examples:\n"
            "  test_resource.py\n"
            "  test_resource.py --read standards://git/pull-requests\n"
            "  test_resource.py --read standards://git/pull-requests --build\n"
        ),
    )
    parser.add_argument(
        "--read",
        metavar="URI",
        dest="read_target",
        default=None,
        help="Read and display the content of the given resource URI.",
    )
    parser.add_argument(
        "--build",
        action="store_true",
        help="Rebuild the binary (go build -o bin/mcp ./cmd) before testing.",
    )
    parser.add_argument(
        "--binary",
        default=None,
        help="Path to the mcp binary. Defaults to <project_root>/bin/mcp.",
    )
    parser.add_argument(
        "--project-root",
        default=None,
        help="Project root directory (where assets/ lives). "
        "Defaults to three directories above this script.",
    )
    parser.add_argument(
        "--timeout",
        type=float,
        default=3.0,
        help="Seconds to wait for the server to process buffered requests before "
        "closing stdin / sending EOF (default: 3.0).",
    )
    args = parser.parse_args()

    # Resolve project root: scripts/ lives at .agents/test-resource-document/scripts/,
    # so the project root is 3 parents up.
    if args.project_root is None:
        args.project_root = str(Path(__file__).resolve().parents[3])

    if not os.path.isdir(args.project_root):
        print(
            f"Error: project root '{args.project_root}' does not exist or is not a directory.\n"
            "Use --project-root to specify the correct path.",
            file=sys.stderr,
        )
        sys.exit(1)

    if args.binary is None:
        args.binary = os.path.join(args.project_root, "bin", "mcp")

    if args.build:
        args.binary = build_binary(args.project_root)

    if not os.path.isfile(args.binary):
        print(
            f"Error: binary not found at '{args.binary}'.\n"
            "Run with --build to compile it, or use --binary <path> to specify an existing binary.",
            file=sys.stderr,
        )
        sys.exit(1)

    if args.read_target:
        handle_read(args)
    else:
        handle_list(args)


if __name__ == "__main__":
    main()
