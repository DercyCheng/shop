package global

import (
	"gorm.io/gorm"
	"nd/user_srv/config"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig
	
	// Prometheus 指标
	RequestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_request_total",
			Help: "用户服务请求总数",
		},
		[]string{"method", "endpoint", "status"},
	)
	
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "user_request_duration_seconds",
			Help:    "用户服务请求处理时间",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
	
	ActiveConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "user_active_connections",
			Help: "用户服务当前活跃连接数",
		},
	)
)

func init() {
	// 注册Prometheus指标
	prometheus.MustRegister(RequestCounter)
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(ActiveConnections)
}
