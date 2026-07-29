#!/bin/bash
# Comprehensive MCP Server test

echo "==> MCP Server Comprehensive Test"
echo ""

# Test 1: Initialize
echo "Test 1: Initialize"
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | ./antia-aitool-mcp --config configs/config.yaml 2>&1 | grep "jsonrpc" | jq .
echo ""

# Test 2: List tools
echo "Test 2: List tools"
(echo '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}'; sleep 0.1) | ./antia-aitool-mcp --config configs/config.yaml 2>&1 | grep "jsonrpc" | jq .
echo ""

# Test 3: Call mock_echo
echo "Test 3: Call mock_echo"
(echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"mock_echo","arguments":{"message":"Hello MCP!"}}}'; sleep 0.1) | ./antia-aitool-mcp --config configs/config.yaml 2>&1 | grep "jsonrpc" | jq .
echo ""

# Test 4: Call mock_time
echo "Test 4: Call mock_time"
(echo '{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"mock_time"}}'; sleep 0.1) | ./antia-aitool-mcp --config configs/config.yaml 2>&1 | grep "jsonrpc" | jq .
echo ""

# Test 5: Call mock_add
echo "Test 5: Call mock_add"
(echo '{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"mock_add","arguments":{"a":10,"b":32}}}'; sleep 0.1) | ./antia-aitool-mcp --config configs/config.yaml 2>&1 | grep "jsonrpc" | jq .
echo ""

echo "==> All tests completed!"
