package configs

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents the application configuration
type Config struct {
	Environment string         `yaml:"environment"`
	Server      ServerConfig   `yaml:"server"`
	Database    DatabaseConfig `yaml:"database"`
}

// ServerConfig contains server-related configuration
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// DatabaseConfig contains database-related configuration
type DatabaseConfig struct {
	Host                   string `yaml:"host"`
	Port                   int    `yaml:"port"`
	User                   string `yaml:"user"`
	Password               string `yaml:"password"`
	Name                   string `yaml:"name"`
	MaxIdleConns           int    `yaml:"maxIdleConns"`
	MaxOpenConns           int    `yaml:"maxOpenConns"`
	ConnMaxLifetimeMinutes int    `yaml:"connMaxLifetimeMinutes"`
}

// LoadConfig loads the configuration from the specified YAML file
func LoadConfig(filePath string) (*Config, error) {
	// Read the config file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal the YAML into the Config struct
	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply environment variable overrides if they exist
	applyEnvironmentOverrides(config)

	return config, nil
}

// applyEnvironmentOverrides applies environment variable overrides to the config values
func applyEnvironmentOverrides(config *Config) {
	// Server settings
	if port := os.Getenv("INVENTORY_SERVER_PORT"); port != "" {
		// Try to parse the port as an integer
		var portInt int
		if _, err := fmt.Sscanf(port, "%d", &portInt); err == nil {
			config.Server.Port = portInt
		}
	}

	// Database settings
	if dbHost := os.Getenv("INVENTORY_DB_HOST"); dbHost != "" {
		config.Database.Host = dbHost
	}
	if dbPort := os.Getenv("INVENTORY_DB_PORT"); dbPort != "" {
		var portInt int
		if _, err := fmt.Sscanf(dbPort, "%d", &portInt); err == nil {
			config.Database.Port = portInt
		}
	}
	if dbUser := os.Getenv("INVENTORY_DB_USER"); dbUser != "" {
		config.Database.User = dbUser
	}
	if dbPassword := os.Getenv("INVENTORY_DB_PASSWORD"); dbPassword != "" {
		config.Database.Password = dbPassword
	}
	if dbName := os.Getenv("INVENTORY_DB_NAME"); dbName != "" {
		config.Database.Name = dbName
	}
}
