package initialize

import (
	"fmt"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"web_golang/order-web/global"
)

func InitRedis() {
	global.RedisClient = goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	pool := goredis.NewPool(global.RedisClient) // or, pool := redigo.NewPool(...)
	global.RedisRs = redsync.New(pool)
}
