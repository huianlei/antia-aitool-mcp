package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/huianlei/antia-aitool-mcp/internal/config"
	"github.com/huianlei/antia-aitool-mcp/internal/mcp"
	"github.com/huianlei/antia-aitool-mcp/internal/plugin"
	_ "github.com/huianlei/antia-aitool-mcp/internal/plugins" // Import to register plugins
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Version information (set by build flags)
	Version   = "dev"
	BuildTime = "unknown"

	// Config file path
	configPath string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "antia-aitool-mcp",
		Short: "Antia AI Tool MCP Server",
		Long:  `A universal MCP server framework for integrating internal services with Claude.`,
		RunE:  run,
	}

	rootCmd.Flags().StringVar(&configPath, "config", "configs/config.yaml", "Path to configuration file")
	rootCmd.AddCommand(versionCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	fmt.Fprintf(os.Stderr, "Antia AI Tool MCP Server %s (built at %s)\n", Version, BuildTime)
	fmt.Fprintf(os.Stderr, "Starting server with config: %s\n", configPath)

	// 1. Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		return errors.Wrap(err, "failed to load configuration")
	}

	fmt.Fprintf(os.Stderr, "Configuration loaded: mode=%s, log_level=%s\n",
		cfg.Server.Mode, cfg.Server.Log.Level)

	// 2. Initialize logger
	logger, err := initLogger(cfg)
	if err != nil {
		return errors.Wrap(err, "failed to initialize logger")
	}
	defer logger.Sync()

	logger.Info("server starting",
		zap.String("version", Version),
		zap.String("build_time", BuildTime),
	)

	// 3. Create plugin manager
	pluginManager := plugin.NewManager(logger)

	// 4. Load enabled plugins
	enabledCount := 0
	for pluginName, pluginConfig := range cfg.Plugins {
		if !pluginConfig.Enabled {
			logger.Debug("plugin disabled", zap.String("plugin", pluginName))
			continue
		}

		logger.Info("loading plugin", zap.String("plugin", pluginName))
		if err := pluginManager.LoadPlugin(pluginName, pluginConfig); err != nil {
			logger.Error("failed to load plugin",
				zap.String("plugin", pluginName),
				zap.Error(err),
			)
			return errors.Wrapf(err, "failed to load plugin: %s", pluginName)
		}
		enabledCount++
	}

	logger.Info("plugins loaded", zap.Int("count", enabledCount))

	if enabledCount == 0 {
		logger.Warn("no plugins enabled")
	}

	// 5. Create MCP server
	mcpServer := mcp.NewServer(pluginManager, logger)

	// 6. Start server based on mode
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	errChan := make(chan error, 1)

	// Start transport
	switch cfg.Server.Mode {
	case "stdio":
		logger.Info("starting stdio transport")
		transport := mcp.NewStdioTransport(mcpServer, logger)
		go func() {
			errChan <- transport.Start(ctx)
		}()

	case "http":
		logger.Info("starting HTTP transport",
			zap.String("host", cfg.Server.HTTP.Host),
			zap.Int("port", cfg.Server.HTTP.Port))
		transport := mcp.NewHTTPTransport(cfg.Server.HTTP.Host, cfg.Server.HTTP.Port, mcpServer, logger)
		go func() {
			errChan <- transport.Start(ctx)
		}()

	default:
		return fmt.Errorf("invalid server mode: %s", cfg.Server.Mode)
	}

	// Wait for shutdown signal or error
	select {
	case sig := <-sigChan:
		logger.Info("received signal", zap.String("signal", sig.String()))
		cancel()
	case err := <-errChan:
		if err != nil && err != context.Canceled {
			logger.Error("transport error", zap.Error(err))
		}
	}

	// 7. Graceful shutdown
	logger.Info("shutting down server")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10)
	defer shutdownCancel()

	if err := mcpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown error", zap.Error(err))
		return err
	}

	logger.Info("server stopped")
	return nil
}

func initLogger(cfg *config.Config) (*zap.Logger, error) {
	// Parse log level
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(cfg.Server.Log.Level)); err != nil {
		return nil, errors.Wrap(err, "invalid log level")
	}

	// Create logger config
	logConfig := zap.NewProductionConfig()
	logConfig.Level = zap.NewAtomicLevelAt(level)

	// Set encoding format
	if cfg.Server.Log.Format == "console" {
		logConfig.Encoding = "console"
		logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Set output paths
	if cfg.Server.Log.File != "" {
		logConfig.OutputPaths = []string{cfg.Server.Log.File, "stderr"}
	} else {
		logConfig.OutputPaths = []string{"stderr"}
	}

	return logConfig.Build()
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version:    %s\n", Version)
			fmt.Printf("Build Time: %s\n", BuildTime)
		},
	}
}
