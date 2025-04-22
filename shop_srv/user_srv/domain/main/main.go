package main

import (
	"crypto/md5"
	"encoding/hex"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"os"
	"time"
	"user_srv/model"
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
	dsn := "root:123456@tcp(localhost:3306)/mxshop_user_srv2?charset=utf8mb4&parseTime=True&loc=Local"
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
	_ = db.AutoMigrate(&model.User{})
	//var users model.User
	//for i := 0; i < 10; i++ {
	//	user := model.User{
	//		NickName: fmt.Sprintf("jzin%d", i),
	//		Mobile:   fmt.Sprintf("1878222222%d", i),
	//		Password: "password",
	//	}
	//	db.Save(&user)
}

//_ = db.AutoMigrate(&model.User{})
//fmt.Println(genMd5("123456"))
// Using the default options
//salt, encodedPwd := password.Encode("generic password", nil)
//fmt.Println(salt)
//fmt.Println(encodedPwd)
//check := password.Verify("generic password", salt, encodedPwd, nil)
//fmt.Println(check) // true

//Using custom options
//options := &password.Options{16, 100, 32, sha512.New}
//salt, encodedPwd := password.Encode("generic password", options)
//newPassword:=fmt.Sprintf("$pbkdf2-sha512$%s$%s",salt,encodedPwd)
//fmt.Println(len(newPassword))
//passwordInfo:=strings.Split(newPassword,"$")
//fmt.Println(passwordInfo)
//check := password.Verify("generic password", passwordInfo[2], passwordInfo[3], options)
//fmt.Println(check) // true

// }
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
