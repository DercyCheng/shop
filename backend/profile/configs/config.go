package configs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"gopkg.in/yaml.v2"
)

// Config holds all configuration for the service
type Config struct {
	Server  ServerConfig  `yaml:"server"`
	MySQL   MySQLConfig   `yaml:"mysql"`
	MongoDB MongoDBConfig `yaml:"mongodb"`
	Redis   RedisConfig   `yaml:"redis"`
	Consul  ConsulConfig  `yaml:"consul"`
	Nacos   NacosConfig   `yaml:"nacos"`
	Log     LogConfig     `yaml:"log"`
	JWT     JWTConfig     `yaml:"jwt"`
	GRPC    GRPCConfig    `yaml:"grpc"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// MySQLConfig holds MySQL configuration
type MySQLConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

// MongoDBConfig holds MongoDB configuration
type MongoDBConfig struct {
	URI         string `yaml:"uri"`
	Database    string `yaml:"database"`
	MaxPoolSize uint64 `yaml:"max_pool_size"`
	MinPoolSize uint64 `yaml:"min_pool_size"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// ConsulConfig holds Consul configuration
type ConsulConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// NacosConfig holds Nacos configuration
type NacosConfig struct {
	Host        string `yaml:"host"`
	Port        uint64 `yaml:"port"`
	NamespaceID string `yaml:"namespace_id"`
	Group       string `yaml:"group"`
	DataID      string `yaml:"data_id"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level     string `yaml:"level"`
	Path      string `yaml:"path"`
	ErrorPath string `yaml:"error_path"`
	Format    string `yaml:"format"`
}

// JWTConfig holds JWT authentication configuration
type JWTConfig struct {
	SecretKey   string `yaml:"secret_key"`
	ExpiresTime int    `yaml:"expires_time"`
	RefreshTime int    `yaml:"refresh_time"`
}

// ServiceConfig holds gRPC service connection configuration
type ServiceConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// GRPCConfig holds all gRPC connection configurations
type GRPCConfig struct {
	ProductService ServiceConfig `yaml:"product_service"`
	UserService    ServiceConfig `yaml:"user_service"`
}

// LoadConfig loads configuration from file
func LoadConfig(configPath string) (*Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("read config file failed: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("unmarshal config failed: %w", err)
	}

	return &config, nil
}

// LoadConfigFromNacos loads configuration from Nacos config center
func LoadConfigFromNacos() (*Config, error) {
	// This function needs the basic Nacos configuration
	// We'll load it from a local file first
	configPath := "configs/config.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Try to find relative to the executable
		execPath, err := os.Executable()
		if err != nil {
			return nil, fmt.Errorf("get executable path failed: %w", err)
		}
		configPath = path.Join(path.Dir(execPath), "configs/config.yaml")
	}

	localConfig, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("load local config failed: %w", err)
	}

	// Create ServerConfig
	sc := []constant.ServerConfig{
		{
			IpAddr: localConfig.Nacos.Host,
			Port:   localConfig.Nacos.Port,
		},
	}

	// Create ClientConfig
	cc := constant.ClientConfig{
		NamespaceId:         localConfig.Nacos.NamespaceID,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogLevel:            "error",
	}

	// Create dynamic config client
	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		return nil, fmt.Errorf("create nacos config client failed: %w", err)
	}

	// Get config
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: localConfig.Nacos.DataID,
		Group:  localConfig.Nacos.Group,
	})
	if err != nil {
		return nil, fmt.Errorf("get config from nacos failed: %w", err)
	}

	// Parse config
	var config Config
	if err := yaml.Unmarshal([]byte(content), &config); err != nil {
		return nil, fmt.Errorf("unmarshal config failed: %w", err)
	}

	// Listen for config changes
	client.ListenConfig(vo.ConfigParam{
		DataId: localConfig.Nacos.DataID,
		Group:  localConfig.Nacos.Group,
		OnChange: func(namespace, group, dataId, data string) {
			var newConfig Config
			if err := yaml.Unmarshal([]byte(data), &newConfig); err != nil {
				fmt.Printf("Failed to parse updated config: %v\n", err)
				return
			}
			// Actually apply the new config in a production system
			// This would need a proper concurrency-safe mechanism
			fmt.Println("Config updated")
		},
	})

	return &config, nil
}
