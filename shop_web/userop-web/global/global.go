package global

import (
	ut "github.com/go-playground/universal-translator"
	"web_golang/userop-web/config"
	"web_golang/userop-web/proto"
)

var (
	Trans          ut.Translator
	ServerConfig   *config.ServerConfig = &config.ServerConfig{}
	NacosConfig    *config.NacosConfig  = &config.NacosConfig{}
	GoodsSrvClient proto.GoodsClient
	MessageClient  proto.MessageClient
	AddressClient  proto.AddressClient
	UserFavClient  proto.UserFavClient
)
