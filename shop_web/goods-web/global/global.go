package global

import (
	ut "github.com/go-playground/universal-translator"
	"web_golang/goods-web/config"
	"web_golang/goods-web/proto"
)

var (
	Trans          ut.Translator
	ServerConfig   *config.ServerConfig = &config.ServerConfig{}
	GoodsSrvClient proto.GoodsClient
	InvSrvClient   proto.InventoryClient
	NacosConfig    *config.NacosConfig = &config.NacosConfig{}
)
