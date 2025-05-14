package configs

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config 商品服务配置结构
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
	
	ElasticSearch struct {
		Addresses []string `yaml:"addresses"`
		Username  string   `yaml:"username"`
		Password  string   `yaml:"password"`
	} `yaml:"elasticsearch"`
	
	Consul struct {
		Address string `yaml:"address"`
	} `yaml:"consul"`
	
	OSS struct {
		Endpoint  string `yaml:"endpoint"`
		AccessKey string `yaml:"accessKey"`
		SecretKey string `yaml:"secretKey"`
		Bucket    string `yaml:"bucket"`
		URLPrefix string `yaml:"urlPrefix"`
	} `yaml:"oss"`
	
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
