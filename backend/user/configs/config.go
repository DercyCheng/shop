package configs

import (
	"time"

	"github.com/spf13/viper"
)

// Config 配置结构
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	MongoDB  MongoDBConfig // 添加MongoDB配置
	JWT      JWTConfig
	Log      LogConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	HTTPPort string
	GRPCPort string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	DSN             string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
}

// MongoDBConfig MongoDB配置
type MongoDBConfig struct {
	URI      string
	Database string
	Timeout  time.Duration
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret          string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string
	Path       string
	ErrorPath  string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Compress   bool
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../configs")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	// 设置默认值
	if config.Server.HTTPPort == "" {
		config.Server.HTTPPort = "8080"
	}
	if config.Server.GRPCPort == "" {
		config.Server.GRPCPort = "9090"
	}

	return config, nil
}
