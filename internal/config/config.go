// Package config handles application configuration loading and management.
package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config 应用配置.
type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	Sandbox SandboxConfig `mapstructure:"sandbox"`
	Log     LogConfig     `mapstructure:"log"`
}

// ServerConfig 服务器配置.
type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

// SandboxConfig 沙箱配置.
type SandboxConfig struct {
	WorkspaceDir string `mapstructure:"workspace_dir"`
	MaxFileSize  int64  `mapstructure:"max_file_size"`
	ShellTimeout int    `mapstructure:"shell_timeout"`
}

// LogConfig 日志配置.
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// Load 加载配置文件.
func Load(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 设置默认值
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("sandbox.workspace_dir", "/tmp/agent-sandbox")
	viper.SetDefault("sandbox.max_file_size", 104857600)
	viper.SetDefault("sandbox.shell_timeout", 300)
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// Address 返回服务器监听地址.
func (c *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
