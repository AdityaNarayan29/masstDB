package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	DefaultDatabase DatabaseConfig `yaml:"default_database"`
	Storage         StorageConfig  `yaml:"storage"`
	Backup          BackupConfig   `yaml:"backup"`
}

// DatabaseConfig holds default database settings
type DatabaseConfig struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// StorageConfig holds storage settings
type StorageConfig struct {
	LocalPath string      `yaml:"local_path"`
	Cloud     CloudConfig `yaml:"cloud"`
}

// CloudConfig holds cloud storage settings
type CloudConfig struct {
	Provider  string `yaml:"provider"` // s3, gcs, azure
	Bucket    string `yaml:"bucket"`
	Region    string `yaml:"region"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
}

// BackupConfig holds backup settings
type BackupConfig struct {
	Compress    bool   `yaml:"compress"`
	DefaultType string `yaml:"default_type"` // full, incremental, differential
}

// DefaultConfig returns a config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		DefaultDatabase: DatabaseConfig{
			Host: "localhost",
		},
		Storage: StorageConfig{
			LocalPath: "./backups",
		},
		Backup: BackupConfig{
			Compress:    true,
			DefaultType: "full",
		},
	}
}

// Load loads configuration from a YAML file
func Load(path string) (*Config, error) {
	config := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil // Return defaults if file doesn't exist
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

// LoadDefault loads configuration from the default location
func LoadDefault() (*Config, error) {
	// Check for config in current directory first
	if _, err := os.Stat(".dbbackup.yaml"); err == nil {
		return Load(".dbbackup.yaml")
	}

	// Check in home directory
	home, err := os.UserHomeDir()
	if err == nil {
		configPath := filepath.Join(home, ".dbbackup.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return Load(configPath)
		}
	}

	return DefaultConfig(), nil
}

// Save saves configuration to a YAML file
func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
