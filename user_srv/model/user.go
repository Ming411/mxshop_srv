package model

import (
	"gorm.io/gorm"
	"time"
)

// 自定义表结构

type BaseModel struct {
	ID        int32          `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:"column:add_time"`
	UpdatedAt time.Time      `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt // 软删除字段
	IsDeleted bool           `gorm:"column:is_deleted"`
}
type User struct {
	BaseModel // 继承
	// 索引  唯一号码  11长度  不为空
	// 这里索引有啥用？
	Mobile   string `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	Password string `gorm:"type:varchar(100);not null"`
	NickName string `gorm:"type:varchar(20)"`
	// 此处为啥要使用指针？？？
	Birthday *time.Time `gorm:"type:datetime"`
	Gender   string     `gorm:"column:gender;default:male;type:varchar(6) comment 'female表示女，male表示男'"`
	Role     int        `gorm:"column:role;default:1;type:int comment '1表示普通用户，2表示管理员'"`
}
