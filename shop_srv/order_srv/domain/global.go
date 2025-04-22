package global

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
	"io"
	"order_srv/config"
	v1goods "order_srv/proto/goods"
	v1inventory "order_srv/proto/inventory"
)

var (
	DB *gorm.DB
	//DB2 *sql.DB
	ServerConfig config.ServerConfig
	//NacosConfig  config.NacosConfig
	//RedisRs      *redsync.Redsync

	JaegerTracer opentracing.Tracer
	JaegerCloser io.Closer

	MQPushClient     rocketmq.PushConsumer
	MQSendTranClient rocketmq.TransactionProducer
	MQSendClient     rocketmq.Producer

	GoodsSrvClient     v1goods.GoodsClient
	InventorySrvClient v1inventory.InventoryClient
)
