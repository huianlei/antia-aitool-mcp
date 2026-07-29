// Package plugins imports all plugin implementations to register them
package plugins

import (
	// Import plugins to trigger their init() registration
	_ "github.com/huianlei/antia-aitool-mcp/internal/plugins/jenkins"
	_ "github.com/huianlei/antia-aitool-mcp/internal/plugins/mock"
	// Future plugins:
	// _ "github.com/huianlei/antia-aitool-mcp/internal/plugins/redis"
	// _ "github.com/huianlei/antia-aitool-mcp/internal/plugins/elasticsearch"
)
