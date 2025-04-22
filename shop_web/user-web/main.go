package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"web_golang/user-web/global"
	"web_golang/user-web/initialize"
	"web_golang/user-web/utils"
	"web_golang/user-web/utils/register/consul"
	myvalidator "web_golang/user-web/validator"
)

func main() {
	//1.初始化logger
	initialize.InitLogger()
	//2.初始化配置文件
	initialize.InitConfig()
	//3.初始化routers
	Router := initialize.Routers()
	//4.初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}
	//5.初始化srv的连接
	initialize.InitSrcConn()
	//6.初始化redis
	initialize.InitRedis()

	//如果是本地开发环境端口号固定，线上环境启动获取端口号
	viper.AutomaticEnv()
	debug := viper.GetBool("MXSHOP_DEBUG")
	if !debug {
		port, err := utils.GetFreePort()
		if err == nil {
			global.ServerConfig.Port = port
		}
	}

	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		_ = v.RegisterTranslation("mobile", global.Trans, func(ut ut.Translator) error {
			return ut.Add("mobile", "{0} must have a value!", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mobile", fe.Field())
			return t
		})
	}
	//服务注册
	register_client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err := register_client.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败:", err.Error())
	}
	/*
		S()可以获取一个全局的sugar，可以让我们自己设置一个全局的logger
		日志是分级别的，debug, info , warn,error， fetal
	*/
	zap.S().Debugf("启动服务器,端口：%d", global.ServerConfig.Port)
	go func() {
		if err = Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
			zap.S().Panic("启动失败：", err.Error())
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	err = register_client.DeRegister(serviceId)
	if err != nil {
		zap.S().Info("注销失败：", err.Error())
	} else {
		zap.S().Info("注销成功")
	}
}
