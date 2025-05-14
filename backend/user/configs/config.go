package configs

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config 用户服务配置结构
type Config struct {
	Server struct {
		Name string `yaml:"name"`
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	
	MySQL struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"mysql"`
	
	Redis struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`
	
	Consul struct {
		Address string `yaml:"address"`
	} `yaml:"consul"`
	
	JWT struct {
		Secret     string `yaml:"secret"`
		ExpireTime int    `yaml:"expireTime"` // Token过期时间（小时）
		Issuer     string `yaml:"issuer"`     // 签发者
	} `yaml:"jwt"`
	
	LogLevel string `yaml:"logLevel"`
	LogFile  string `yaml:"logFile"`
}

// LoadConfig 从配置文件加载配置
func LoadConfig() (*Config, error) {
	config := &Config{}
	
	// 配置文件路径可以通过环境变量指定，默认为当前目录下的config.yaml
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}
	
	// 读取配置文件
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	
	// 解析YAML
	if err := yaml.Unmarshal(file, config); err != nil {
		return nil, err
	}
	
	return config, nil
}
