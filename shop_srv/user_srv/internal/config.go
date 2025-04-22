package initialize

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"user_srv/global"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}
func InitConfig() {
	//从配置文件中读取对应的配置
	pro := GetEnvInfo("JZIN_PRO")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("%s-debug.yaml", configFilePrefix)
	if pro {
		configFileName = fmt.Sprintf("%s-pro.yaml", configFilePrefix)
	}
	v := viper.New()
	//文件的路径如何设置
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	//全局变量
	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		panic(err)
	}
	zap.S().Infof("配置信息：%v", global.ServerConfig)

}
