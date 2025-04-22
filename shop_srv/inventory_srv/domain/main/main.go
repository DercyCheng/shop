package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"inventory_srv/model"
	"io"
	"log"
	"os"
	"time"
)

type User struct {
	gorm.Model
	Name string
}

func genMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}
func main() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // 日志级别
			//IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful: true, // 禁用彩色打印
		},
	)
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	dsn := "root:123456@tcp(localhost:3306)/mxshop_inventory_srv2?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	//var users model.User

	//_ = db.AutoMigrate(&model.Inventory{}, &model.StockSellDetail{})
	//orderDetail := model.StockSellDetail{
	//	OrderSn: "imooc-jzin",
	//	Status:  1,
	//	Detail:  []model.GoodsDetail{{1, 2}, {2, 3}},
	//}
	//db.Create(&orderDetail)
	var sellDetail model.StockSellDetail
	db.Where(&model.StockSellDetail{OrderSn: "imooc-jzin"}).First(&sellDetail)
	//fmt.Println(sellDetail.Detail)
	for _, good := range sellDetail.Detail {
		fmt.Println(good.Goods, good.Num)
	}

}
func Paginate(pn, pSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := pn
		pageSize := pSize
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
