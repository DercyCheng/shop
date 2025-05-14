package configs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
	
	"gopkg.in/yaml.v3"
)

// Config 配置结构体
type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Database  DatabaseConfig  `yaml:"database"`
	Redis     RedisConfig     `yaml:"redis"`
	Logger    LoggerConfig    `yaml:"logger"`
	Registry  RegistryConfig  `yaml:"registry"`
	Inventory InventoryConfig `yaml:"inventory"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Name         string        `yaml:"name"`
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	GRPC         GRPCConfig    `yaml:"grpc"`
}

// GRPCConfig GRPC配置
type GRPCConfig struct {
	Port int `yaml:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	Charset  string `yaml:"charset"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level      string `yaml:"level"`
	Format     string `yaml:"format"`
	OutputFile string `yaml:"output_file"`
	MaxSize    int    `yaml:"max_size"`    // 日志文件最大大小，单位MB
	MaxBackups int    `yaml:"max_backups"` // 保留的旧文件最大数量
	MaxAge     int    `yaml:"max_age"`     // 保留的旧文件最大天数
	Compress   bool   `yaml:"compress"`    // 是否压缩
}

// RegistryConfig 服务注册配置
type RegistryConfig struct {
	Type     string      `yaml:"type"`
	Consul   ConsulConfig `yaml:"consul"`
	Nacos    NacosConfig  `yaml:"nacos"`
	Interval int         `yaml:"interval"` // 健康检查间隔，单位秒
}

// ConsulConfig Consul配置
type ConsulConfig struct {
	Address string `yaml:"address"`
}

// NacosConfig Nacos配置
type NacosConfig struct {
	Address   string `yaml:"address"`
	Port      uint64 `yaml:"port"`
	Namespace string `yaml:"namespace"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
}

// InventoryConfig 库存服务特定配置
type InventoryConfig struct {
	LockTimeout        int  `yaml:"lock_timeout"`        // 库存锁定超时时间（秒）
	DefaultWarehouseID int  `yaml:"default_warehouse_id"` // 默认仓库ID
	EnableCache        bool `yaml:"enable_cache"`        // 是否启用缓存
	CacheTTL           int  `yaml:"cache_ttl"`          // 缓存过期时间（秒）
}

// LoadConfig 加载配置
func LoadConfig(configFile string) (*Config, error) {
	// 如果配置文件路径为空，则使用默认路径
	if configFile == "" {
		// 获取当前执行文件所在目录
		exePath, err := os.Executable()
		if err != nil {
			return nil, err
		}
		exeDir := filepath.Dir(exePath)
		configFile = filepath.Join(exeDir, "configs", "config.yaml")
	}

	// 读取配置文件
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	// 解析配置
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
