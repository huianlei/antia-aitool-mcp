#!/usr/bin/env bash
# Test Jenkins MCP plugin by sending JSON-RPC requests via stdio

set -e

SERVER="./antia-aitool-mcp"
CONFIG="configs/config.yaml"

echo "==> Testing Jenkins MCP Plugin"
echo ""

# Start the server
echo "Starting MCP server..."
export JENKINS_PASSWORD="hal@123"

# Test 1: Initialize
echo "Test 1: Initialize connection"
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"0.1.0","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}' | $SERVER --config $CONFIG 2>/dev/null | jq .

# Test 2: List tools
echo ""
echo "Test 2: List available tools"
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' | $SERVER --config $CONFIG 2>/dev/null | jq '.result.tools[] | select(.name | startswith("jenkins"))'

# Test 3: Call jenkins_list_jobs
echo ""
echo "Test 3: Call jenkins_list_jobs"
echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"jenkins_list_jobs","arguments":{}}}' | $SERVER --config $CONFIG 2>/dev/null | jq .

echo ""
echo "==> Tests completed!"
