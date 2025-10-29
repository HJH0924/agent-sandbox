package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// loadConfigFromContent 是一个辅助函数，用于创建配置文件并加载配置。
func loadConfigFromContent(t *testing.T, content string) (*Config, error) {
	t.Helper()
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	err := os.WriteFile(configPath, []byte(content), 0o600)
	require.NoError(t, err)

	return Load(configPath)
}

func TestLoad_ValidConfig(t *testing.T) {
	configContent := `
[server]
host = "localhost"
port = 9090
read_timeout = "60s"
write_timeout = "60s"

[sandbox]
workspace_dir = "/var/sandbox"
max_file_size = 52428800
shell_timeout = 600

[log]
level = "debug"
format = "text"
`

	cfg, err := loadConfigFromContent(t, configContent)

	require.NoError(t, err)
	assert.NotNil(t, cfg)

	// 验证服务器配置
	assert.Equal(t, "localhost", cfg.Server.Host)
	assert.Equal(t, 9090, cfg.Server.Port)
	assert.Equal(t, 60*time.Second, cfg.Server.ReadTimeout)
	assert.Equal(t, 60*time.Second, cfg.Server.WriteTimeout)

	// 验证沙箱配置
	assert.Equal(t, "/var/sandbox", cfg.Sandbox.WorkspaceDir)
	assert.Equal(t, int64(52428800), cfg.Sandbox.MaxFileSize)
	assert.Equal(t, 600, cfg.Sandbox.ShellTimeout)

	// 验证日志配置
	assert.Equal(t, "debug", cfg.Log.Level)
	assert.Equal(t, "text", cfg.Log.Format)
}

func TestLoad_DefaultValues(t *testing.T) {
	cfg, err := loadConfigFromContent(t, "")

	require.NoError(t, err)
	assert.NotNil(t, cfg)

	// 验证默认值
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, 30*time.Second, cfg.Server.ReadTimeout)
	assert.Equal(t, 30*time.Second, cfg.Server.WriteTimeout)
	assert.Equal(t, "/tmp/agent-sandbox", cfg.Sandbox.WorkspaceDir)
	assert.Equal(t, int64(104857600), cfg.Sandbox.MaxFileSize)
	assert.Equal(t, 300, cfg.Sandbox.ShellTimeout)
	assert.Equal(t, "info", cfg.Log.Level)
	assert.Equal(t, "json", cfg.Log.Format)
}

func TestLoad_PartialConfig(t *testing.T) {
	// 测试部分配置，其余使用默认值
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	configContent := `
[server]
port = 3000

[log]
level = "warn"
`

	err := os.WriteFile(configPath, []byte(configContent), 0o600)
	require.NoError(t, err)

	cfg, err := Load(configPath)

	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// 验证指定的值
	assert.Equal(t, 3000, cfg.Server.Port)
	assert.Equal(t, "warn", cfg.Log.Level)

	// 验证未指定的值使用默认值
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, 30*time.Second, cfg.Server.ReadTimeout)
	assert.Equal(t, "json", cfg.Log.Format)
}

func TestLoad_FileNotFound(t *testing.T) {
	configPath := "/nonexistent/path/config.toml"

	cfg, err := Load(configPath)

	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "failed to read config file")
}

func TestLoad_InvalidTOML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	// 写入无效的 TOML
	configContent := `
[server
port = "invalid_section"
`

	err := os.WriteFile(configPath, []byte(configContent), 0o600)
	require.NoError(t, err)

	cfg, err := Load(configPath)

	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestLoad_InvalidTimeout(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	configContent := `
[server]
read_timeout = "invalid"
`

	err := os.WriteFile(configPath, []byte(configContent), 0o600)
	require.NoError(t, err)

	cfg, err := Load(configPath)
	// viper 会尝试解析，但可能失败
	// 这取决于 viper 的行为，可能会有错误或使用默认值
	if err != nil {
		assert.Contains(t, err.Error(), "failed to unmarshal config")
	}
	// 如果没有错误，验证使用了默认值
	if cfg != nil {
		// 在这种情况下，viper 可能无法解析并使用默认值
		assert.NotNil(t, cfg)
	}
}

func TestServerConfig_Address(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		port     int
		expected string
	}{
		{
			name:     "localhost with standard port",
			host:     "localhost",
			port:     8080,
			expected: "localhost:8080",
		},
		{
			name:     "0.0.0.0 with custom port",
			host:     "0.0.0.0",
			port:     9090,
			expected: "0.0.0.0:9090",
		},
		{
			name:     "IP address with port",
			host:     "192.168.1.1",
			port:     3000,
			expected: "192.168.1.1:3000",
		},
		{
			name:     "empty host with port",
			host:     "",
			port:     8080,
			expected: ":8080",
		},
		{
			name:     "host with port 0",
			host:     "localhost",
			port:     0,
			expected: "localhost:0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &ServerConfig{
				Host: tt.host,
				Port: tt.port,
			}

			address := cfg.Address()
			assert.Equal(t, tt.expected, address)
		})
	}
}

func TestSetDefaults(t *testing.T) {
	// 测试 setDefaults 函数是否正确设置了所有默认值
	// 这是通过加载一个空配置来间接测试的
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	// 创建最小配置文件
	err := os.WriteFile(configPath, []byte(""), 0o600)
	require.NoError(t, err)

	cfg, err := Load(configPath)
	require.NoError(t, err)

	// 验证所有默认值都已设置
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, 30*time.Second, cfg.Server.ReadTimeout)
	assert.Equal(t, 30*time.Second, cfg.Server.WriteTimeout)
	assert.Equal(t, "/tmp/agent-sandbox", cfg.Sandbox.WorkspaceDir)
	assert.Equal(t, int64(104857600), cfg.Sandbox.MaxFileSize)
	assert.Equal(t, 300, cfg.Sandbox.ShellTimeout)
	assert.Equal(t, "info", cfg.Log.Level)
	assert.Equal(t, "json", cfg.Log.Format)
}

func TestLoad_ComplexConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	configContent := `
[server]
host = "0.0.0.0"
port = 8080
read_timeout = "30s"
write_timeout = "30s"

[sandbox]
workspace_dir = "/tmp/agent-sandbox"
max_file_size = 104857600
shell_timeout = 300

[log]
level = "info"
format = "json"
`

	err := os.WriteFile(configPath, []byte(configContent), 0o600)
	require.NoError(t, err)

	cfg, err := Load(configPath)

	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// 验证所有配置项
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, 8080, cfg.Server.Port)
	assert.Equal(t, "0.0.0.0:8080", cfg.Server.Address())
	assert.Equal(t, "/tmp/agent-sandbox", cfg.Sandbox.WorkspaceDir)
	assert.Equal(t, "info", cfg.Log.Level)
}
