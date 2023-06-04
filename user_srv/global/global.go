package global

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var (
	DB *gorm.DB
)

// 如果函数名为 init 那么当导入时 该函数自动执行
func init() {
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
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 单数表，在建表时不会自动添加s
		},
		Logger: newLogger,
	})
	if err != nil {
		panic("数据库连接失败")
	}
}
