package initialize

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
	"goods_srv/global"
	"goods_srv/model"
)

func InitEs() {
	//初始化连接
	host := fmt.Sprintf("http://%s:%d", global.ServerConfig.EsInfo.Host, global.ServerConfig.EsInfo.Port)
	var err error
	//logger := log.New(os.Stdout, "mxshop", log.LstdFlags)
	//global.EsClient, err = elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false),
	//	elastic.SetTraceLog(logger))
	global.EsClient, err = elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}
	//新建mapping和index
	//查询index是否存在
	exists, err := global.EsClient.IndexExists(model.EsGoods{}.GetIndexName()).Do(context.Background())
	if err != nil {
		panic(err)
	}
	//新建mapping
	if !exists {
		_, err2 := global.EsClient.CreateIndex(model.EsGoods{}.GetIndexName()).BodyString(model.EsGoods{}.GetMapping()).Do(context.Background())
		if err2 != nil {
			zap.S().Fatalf("创建索引%s失败：%s", model.EsGoods{}.GetIndexName(), err2.Error())
		}
	}
}
