package config

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}
type JaegerConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

//	type RedisConfig struct {
//		Host string `mapstructure:"host" json:"host"`
//		Port int    `mapstructure:"port" json:"port"`
//		Db   int    `mapstructure:"db" json:"db"`
//	}
type GoodsSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}
type InventorySrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}
type RocketMQConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}
type ServerConfig struct {
	Name       string       `mapstructure:"name" json:"name"`
	Tags       []string     `mapstructure:"tags" json:"tags"`
	Host       string       `mapstructure:"host" json:"host"`
	MysqlInfo  MysqlConfig  `mapstructure:"mysql" json:"mysql"`
	ConsulInfo ConsulConfig `mapstructure:"consul" json:"consul"`
	//RedisInfo  RedisConfig  `mapstructure:"redis" json:"redis"`
	RocketMQConfig RocketMQConfig `mapstructure:"rocketmq" json:"rocketmq"`
	//商品微服务的配置
	GoodsSrvInfo GoodsSrvConfig `mapstructure:"v1goods_srv" json:"v1goods_srv"`
	//库存微服务的配置
	InventorySrvInfo InventorySrvConfig `mapstructure:"v1inventory_srv" json:"v1inventory_srv"`
	JaegerInfo       JaegerConfig       `mapstructure:"jaeger" json:"jaeger"`
}
