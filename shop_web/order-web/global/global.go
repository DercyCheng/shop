package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"web_golang/order-web/config"
	"web_golang/order-web/proto"
)

var (
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
	NacosConfig  *config.NacosConfig  = &config.NacosConfig{}

	GoodsSrvClient     proto.GoodsClient
	OrderSrvClient     proto.OrderClient
	InventorySrvClient proto.InventoryClient

	Trans       ut.Translator
	RedisRs     *redsync.Redsync
	RedisClient *redis.Client
)
