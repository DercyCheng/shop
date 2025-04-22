package initialize

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"go.uber.org/zap"
)

func InitSentinel() {
	err := sentinel.InitDefault()
	if err != nil {
		zap.S().Fatalf("初始化异常：%v", err)
	}
	//配置限流规则
	//这种配置应该从nacos中读取，但是官方文档还没有支持热加载
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "goods-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              12,
			StatIntervalInMs:       6000,
		},
	})
	if err != nil {
		zap.S().Fatalf("配置规则出错：%v", err)
	}
}
