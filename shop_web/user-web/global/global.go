package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis/v8"
	"web_golang/user-web/config"
	"web_golang/user-web/proto"
)

var (
	Trans         ut.Translator
	ServerConfig  *config.ServerConfig = &config.ServerConfig{}
	UserSrvClient proto.UserClient
	NacosConfig   *config.NacosConfig = &config.NacosConfig{}
	RedisClient   *redis.Client
)
