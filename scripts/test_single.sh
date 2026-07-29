#!/bin/bash
# Simple MCP test - send one request and exit

echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | ./antia-aitool-mcp --config configs/config.yaml 2>&1 &

PID=$!

# Wait 2 seconds for response
sleep 2

# Kill the process
kill $PID 2>/dev/null

# Wait for process to finish
wait $PID 2>/dev/null

echo ""
echo "Test complete"
