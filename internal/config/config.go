package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Server  ServerConfig             `mapstructure:"server"`
	Plugins map[string]PluginConfig  `mapstructure:"plugins"`
}

// ServerConfig holds server-level settings
type ServerConfig struct {
	Mode string     `mapstructure:"mode"` // stdio or http
	HTTP HTTPConfig `mapstructure:"http"`
	Log  LogConfig  `mapstructure:"log"`
}

// HTTPConfig holds HTTP server settings
type HTTPConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string `mapstructure:"level"`  // debug, info, warn, error
	Format string `mapstructure:"format"` // json, console
	File   string `mapstructure:"file"`   // empty = stderr
}

// PluginConfig holds plugin-specific configuration
type PluginConfig struct {
	Enabled bool                   `mapstructure:"enabled"`
	Raw     map[string]interface{} `mapstructure:",remain"`
}

// Load loads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Read config file
	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			return nil, errors.Wrap(err, "failed to read config file")
		}
	}

	// Environment variables override config file
	// ANTIA_SERVER_LOG_LEVEL=debug -> server.log.level
	v.SetEnvPrefix("ANTIA")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Unmarshal config
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal config")
	}

	// Expand environment variables in plugin configs
	expandEnvVars(&cfg)

	// Validate config
	if err := validate(&cfg); err != nil {
		return nil, errors.Wrap(err, "config validation failed")
	}

	return &cfg, nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.mode", "stdio")
	v.SetDefault("server.http.enabled", false)
	v.SetDefault("server.http.host", "0.0.0.0")
	v.SetDefault("server.http.port", 8080)
	v.SetDefault("server.log.level", "info")
	v.SetDefault("server.log.format", "json")
	v.SetDefault("server.log.file", "")
}

// validate validates the configuration
func validate(cfg *Config) error {
	// Validate server mode
	if cfg.Server.Mode != "stdio" && cfg.Server.Mode != "http" {
		return fmt.Errorf("invalid server mode: %s (must be 'stdio' or 'http')", cfg.Server.Mode)
	}

	// Validate log level
	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[cfg.Server.Log.Level] {
		return fmt.Errorf("invalid log level: %s", cfg.Server.Log.Level)
	}

	// Validate log format
	if cfg.Server.Log.Format != "json" && cfg.Server.Log.Format != "console" {
		return fmt.Errorf("invalid log format: %s (must be 'json' or 'console')", cfg.Server.Log.Format)
	}

	// Validate HTTP config if enabled
	if cfg.Server.HTTP.Enabled {
		if cfg.Server.HTTP.Port < 1 || cfg.Server.HTTP.Port > 65535 {
			return fmt.Errorf("invalid HTTP port: %d", cfg.Server.HTTP.Port)
		}
	}

	return nil
}

// expandEnvVars expands environment variables in the format ${VAR} or $VAR
func expandEnvVars(cfg *Config) {
	envVarRegex := regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Z_][A-Z0-9_]*)`)

	// Expand environment variables in all plugin configs
	for pluginName, pluginCfg := range cfg.Plugins {
		expandMap(pluginCfg.Raw, envVarRegex)
		cfg.Plugins[pluginName] = pluginCfg
	}
}

// expandMap recursively expands environment variables in a map
func expandMap(m map[string]interface{}, regex *regexp.Regexp) {
	for key, val := range m {
		switch v := val.(type) {
		case string:
			m[key] = expandString(v, regex)
		case map[string]interface{}:
			expandMap(v, regex)
		case []interface{}:
			for i, item := range v {
				if s, ok := item.(string); ok {
					v[i] = expandString(s, regex)
				} else if subMap, ok := item.(map[string]interface{}); ok {
					expandMap(subMap, regex)
				}
			}
		}
	}
}

// expandString expands environment variables in a string
func expandString(s string, regex *regexp.Regexp) string {
	return regex.ReplaceAllStringFunc(s, func(match string) string {
		// Extract variable name from ${VAR} or $VAR
		varName := strings.TrimPrefix(strings.TrimSuffix(strings.TrimPrefix(match, "${"), "}"), "$")
		if envVal := os.Getenv(varName); envVal != "" {
			return envVal
		}
		return match // Keep original if env var not found
	})
}

// GetPluginConfig retrieves configuration for a specific plugin
func (c *Config) GetPluginConfig(pluginName string) (PluginConfig, bool) {
	cfg, exists := c.Plugins[pluginName]
	return cfg, exists
}

// GetString retrieves a string value from plugin config
func (pc PluginConfig) GetString(key string) string {
	if val, ok := pc.Raw[key].(string); ok {
		return val
	}
	return ""
}

// GetInt retrieves an int value from plugin config
func (pc PluginConfig) GetInt(key string) int {
	if val, ok := pc.Raw[key].(int); ok {
		return val
	}
	return 0
}

// GetBool retrieves a bool value from plugin config
func (pc PluginConfig) GetBool(key string) bool {
	if val, ok := pc.Raw[key].(bool); ok {
		return val
	}
	return false
}

// GetDuration retrieves a duration value from plugin config
func (pc PluginConfig) GetDuration(key string) time.Duration {
	if val, ok := pc.Raw[key].(string); ok {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return 0
}

// GetMap retrieves a map from plugin config
func (pc PluginConfig) GetMap(key string) map[string]interface{} {
	if val, ok := pc.Raw[key].(map[string]interface{}); ok {
		return val
	}
	return nil
}

// GetStringSlice retrieves a string slice from plugin config
func (pc PluginConfig) GetStringSlice(key string) []string {
	if val, ok := pc.Raw[key].([]interface{}); ok {
		result := make([]string, 0, len(val))
		for _, v := range val {
			if s, ok := v.(string); ok {
				result = append(result, s)
			}
		}
		return result
	}
	return nil
}
