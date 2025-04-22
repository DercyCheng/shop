package global

import (
	"github.com/go-redsync/redsync/v4"
	"gorm.io/gorm"
	"inventory_srv/config"
)

var (
	DB *gorm.DB
	//DB2 *sql.DB
	ServerConfig config.ServerConfig
	//NacosConfig  config.NacosConfig
	RedisRs *redsync.Redsync
)
