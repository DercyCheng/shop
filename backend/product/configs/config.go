package configs

import (
	"os"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	ServiceName string       `mapstructure:"service_name"`
	Server      ServerConfig `mapstructure:"server"`
	Database    DBConfig     `mapstructure:"database"`
}

// ServerConfig represents the server configuration
type ServerConfig struct {
	Port int `mapstructure:"port"`
}

// DBConfig represents the database configuration
type DBConfig struct {
	DSN string `mapstructure:"dsn"`
}

// LoadConfig loads the configuration from file or environment variables
func LoadConfig(configPath string) (*Config, error) {
	// Set default values
	viper.SetDefault("service_name", "product-service")
	viper.SetDefault("server.port", 50051)
	viper.SetDefault("database.dsn", "root:password@tcp(mysql:3306)/shop_product?charset=utf8mb4&parseTime=True&loc=Local")

	// Read config file if exists
	if configPath != "" {
		viper.SetConfigFile(configPath)
		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	// Override with environment variables if they exist
	if port := os.Getenv("PRODUCT_SERVICE_PORT"); port != "" {
		viper.Set("server.port", port)
	}

	if dsn := os.Getenv("PRODUCT_DB_DSN"); dsn != "" {
		viper.Set("database.dsn", dsn)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
