package initialize

import (
	"fmt"
	"log"
	"os"
	"time"

	//_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"inventory_srv/global"
)

func InitDB() {
	//dsn := "root:123456@tcp(localhost:3306)/mxshop_inventory_srv2?charset=utf8mb4&parseTime=True&loc=Local"
	c := global.ServerConfig.MysqlInfo
	//zap.S().Info(c.User, c.Password, c.Host, c.Port, c.Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Name)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // 日志级别
			//LogLevel: logger.Silent, // 日志级别
			//IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful: true, // 禁用彩色打印
		},
	)
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情

	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	//DB2()
}

//func DB2() {
//	var err error
//	c := global.ServerConfig.MysqlInfo
//	//cnnstr := "user:pwd@tcp(127.0.0.1:3306)/db?charset=utf8&parseTime=True"
//	cnnstr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", c.User, c.Password, c.Host, c.Port, c.Name)
//	global.DB2, err = sql.Open("mysql", cnnstr)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	err = global.DB2.Ping()
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	err = worm.InitMysql(global.DB2)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//}
