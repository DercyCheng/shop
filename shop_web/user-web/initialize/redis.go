package initialize

import (
	"fmt"
	goredislib "github.com/go-redis/redis/v8"
	"web_golang/user-web/global"
)

func InitRedis() {
	global.RedisClient = goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
		DB:   0,
	})
}
