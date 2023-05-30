package main

import (
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"mx_shop/user_srv/model"
	"os"
	"strings"
	"time"
)

//func genMd5(code string) string {
//	Md5 := md5.New()
//	_, _ = io.WriteString(Md5, code)
//	return hex.EncodeToString(Md5.Sum(nil))
//}

func main() {
	// Using custom options
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode("generic password", options)
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
	fmt.Println(newPassword)
	passwordInfo := strings.Split(newPassword, "$")
	fmt.Println(passwordInfo) // 从1开始
	check := password.Verify("generic password", passwordInfo[2], passwordInfo[3], options)
	//check := password.Verify("generic password", salt, encodedPwd, options)
	fmt.Println(check) // true

	dsn := "root:root@tcp(172.18.81.229:3306)/mxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
	// 配置日志输出  可查看具体执行的SQL
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 单数表，在建表时不会自动添加s
		},
		Logger: newLogger,
	})
	if err != nil {
		panic("数据库连接失败")
	}
	// 根据定义的表结构建表
	_ = db.AutoMigrate(&model.User{})

}
