#!/bin/bash
# Test MCP Server with mock plugin

set -e

echo "==> Testing MCP Server"
echo ""

# Test 1: Initialize
echo "Test 1: Initialize"
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | ./antia-aitool-mcp --config configs/config.yaml

echo ""
echo "---"
echo ""

# Test 2: List tools
echo "Test 2: List tools"
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' | ./antia-aitool-mcp --config configs/config.yaml

echo ""
echo "---"
echo ""

# Test 3: Call mock_echo
echo "Test 3: Call mock_echo"
echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"mock_echo","arguments":{"message":"Hello MCP!"}}}' | ./antia-aitool-mcp --config configs/config.yaml

echo ""
echo "---"
echo ""

# Test 4: Call mock_time
echo "Test 4: Call mock_time"
echo '{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"mock_time"}}' | ./antia-aitool-mcp --config configs/config.yaml

echo ""
echo "---"
echo ""

# Test 5: Call mock_add
echo "Test 5: Call mock_add"
echo '{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"mock_add","arguments":{"a":10,"b":32}}}' | ./antia-aitool-mcp --config configs/config.yaml

echo ""
echo ""
echo "==> All tests completed!"
