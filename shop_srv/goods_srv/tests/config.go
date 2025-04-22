package main

import (
	"github.com/spf13/viper"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}
func InitConfig() {
	//从配置文件中读取对应的配置
	//debug := GetEnvInfo("MXSHOP_DEBUG")
	//configFilePrefix := "config"
	//configFileName := fmt.Sprintf("goods_srv/%s-pro.yaml", configFilePrefix)
	//if debug {
	//	configFileName = fmt.Sprintf("goods_srv/%s-debug.yaml", configFilePrefix)
	//}
	//v := viper.New()
	////文件的路径如何设置
	//v.SetConfigFile(configFileName)
	//if err := v.ReadInConfig(); err != nil {
	//	panic(err)
	//}
	////全局变量
	//if err := v.Unmarshal(&global.NacosConfig); err != nil {
	//	panic(err)
	//}
	//zap.S().Infof("配置信息：%v", global.NacosConfig)
	//sc := []constant.ServerConfig{
	//	*constant.NewServerConfig(global.NacosConfig.Host, global.NacosConfig.Port, constant.WithContextPath("/nacos")),
	//}
	//cc := *constant.NewClientConfig(
	//	constant.WithNamespaceId(global.NacosConfig.Namespace),
	//	constant.WithTimeoutMs(5000),
	//	constant.WithNotLoadCacheAtStart(true),
	//	constant.WithLogDir("tmp/nacos/log"),
	//	constant.WithCacheDir("tmp/nacos/cache"),
	//	constant.WithLogLevel("debug"),
	//)
	//client, err := clients.NewConfigClient(
	//	vo.NacosClientParam{
	//		ClientConfig:  &cc,
	//		ServerConfigs: sc,
	//	},
	//)
	//if err != nil {
	//	fmt.Printf("PublishConfig err:%+v \n", err)
	//}
	//content, err := client.GetConfig(vo.ConfigParam{
	//	DataId: global.NacosConfig.DataId,
	//	Group:  global.NacosConfig.Group,
	//})
	//if err != nil {
	//	zap.S().Fatalf("读取nacos配置失败： %s", err.Error())
	//}
	////fmt.Println(content)
	////serverConfig := config.ServerConfig{}
	////想要将一个字符串转换成struct需要去设置这个struct的tag
	//err = json.Unmarshal([]byte(content), &global.ServerConfig)
	//if err != nil {
	//	zap.S().Fatalf("读取nacos配置失败： %s", err.Error())
	//}
	//zap.S().Info(global.NacosConfig)
	//fmt.Println(serverConfig)

}
