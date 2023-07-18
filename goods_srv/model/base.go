package model

import (
	"database/sql/driver"
	"time"

	"github.com/goccy/go-json"
	"gorm.io/gorm"
)

type BaseModel struct {
	// type:int 声明其在数据库中的数据类型为 int
	ID        int32          `gorm:"primarykey;type:int"` // 主键
	CreatedAt time.Time      `gorm:"column:add_time"`     // 创建时间
	UpdatedAt time.Time      `gorm:"column:update_time"`  // 更新时间
	DeletedAt gorm.DeletedAt // 软删除
	IsDeleted bool           `gorm:"column:is_deleted"` // 是否删除
}

// 用于定义表结构
type GormList []string

func (g GormList) Value() (driver.Value, error) {
	return json.Marshal(g)
}

// 加不加 * 的区别在于方法内部对接收者的修改是否会影响到原始的调用对象。
// 如果你需要在方法内部修改原始的调用对象，可以使用指针作为接收者；
// 如果不需要修改原始的调用对象，可以使用值类型作为接收者。

func (g *GormList) Scan(input interface{}) error {
	// input.([]byte) 将 input 转换为 []byte 类型
	return json.Unmarshal(input.([]byte), &g)
}
