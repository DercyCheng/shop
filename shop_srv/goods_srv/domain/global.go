package global

import (
	"github.com/olivere/elastic/v7"
	"goods_srv/config"
	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	//NacosConfig  config.NacosConfig
	EsClient *elastic.Client
)
