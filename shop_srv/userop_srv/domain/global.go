package global

import (
	"gorm.io/gorm"
	"userop_srv/config"
)

var (
	DB *gorm.DB
	//DB2 *sql.DB
	ServerConfig config.ServerConfig
	//NacosConfig  config.NacosConfig
	//RedisRs      *redsync.Redsync
)
