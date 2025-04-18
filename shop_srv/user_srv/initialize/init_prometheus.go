package initialize

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// InitPrometheusServer 初始化Prometheus指标暴露服务器
func InitPrometheusServer(port int) {
	// 创建HTTP服务来暴露Prometheus指标
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		
		// 添加健康检查端点
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		
		addr := fmt.Sprintf(":%d", port)
		server := &http.Server{
			Addr:         addr,
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		
		zap.S().Infof("启动Prometheus指标暴露服务器在 %s", addr)
		
		// 添加重试机制
		maxRetries := 3
		retryInterval := 5 * time.Second
		
		var err error
		for i := 0; i < maxRetries; i++ {
			err = server.ListenAndServe()
			if err == nil {
				break
			}
			
			zap.S().Warnf("Prometheus指标服务器启动失败(尝试 %d/%d): %s", i+1, maxRetries, err)
			if i < maxRetries-1 {
				zap.S().Infof("将在 %s 后重试...", retryInterval)
				time.Sleep(retryInterval)
			}
		}
		
		if err != nil {
			zap.S().Errorf("Prometheus指标服务器启动失败，所有重试均失败: %s", err)
		}
	}()
}