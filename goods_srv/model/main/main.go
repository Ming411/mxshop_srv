package main

import (
	"log"
	"mx_shop/goods_srv/model"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func main() {
	dsn := "root:root@tcp(172.18.81.229:3306)/mxshop_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"
	// 配置日志输出  可查看具体执行的SQL
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // 输出INFO级别的日志
			Colorful:      true,        // 禁用彩色打印
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
	_ = db.AutoMigrate(
		&model.Category{},
		&model.Brands{},
		&model.Goods{},
		&model.GoodsCategoryBrand{},
		&model.Banner{},
	)

}
