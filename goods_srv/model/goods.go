// 用于定义表结构
package model

// 商品类型 表
type Category struct {
	BaseModel                  // 同一个包下可以不用引用
	Name             string    `gorm:"type:varchar(20);not null"`
	ParentCategoryID int32     `gorm:"type:int;not null"` // 父级ID
	ParentCategory   *Category // 父级分类，这里指向自己
	Level            int32     `gorm:"type:int;not null;default:1"`      // 列表级别
	IsTab            bool      `gorm:"type:bool;not null;default:false"` // 是否拥有tab
}

// 商品品牌 表
type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(20);not null"`
	Logo string `gorm:"type:varchar(200);default:'';not null"` // 品牌logo 以URL形式存储
}

/*
由于 类型 和 品牌 是多对多的关系，所以需要一个中间表进行关联
*/
type GoodsCategoryBrand struct {
	BaseModel
	// index 用于创建唯一索引
	CategoryID int32    `gorm:"type:int;index:idx_category_brand,unique;not null"` // 商品分类ID
	Category   Category // 用于关联查询 创建外键

	BrandID int32 `gorm:"type:int;index:idx_category_brand;not null"`
	Brand   Brands
}

// 设置表名 —— 如果不设置，gorm会默认将结构体名称转换为小写通过_连接并且加上s
func (GoodsCategoryBrand) TableName() string {
	return "goods_category_brand"
}

// 轮播图
type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null"`  // 图片地址
	Url   string `gorm:"type:varchar(200);not null"`  // 跳转商品详情地址
	Index int32  `gorm:"type:int;default:1;not null"` // 排序位置
}

// 商品表
type Goods struct {
	BaseModel
	CategoryID      int32    `gorm:"type:int;not null"`
	Category        Category // 用于关联查询 创建外键
	BrandID         int32    `gorm:"type:int;not null"`
	Brand           Brands
	OnSale          bool     `gorm:"type:bool;default:false;not null"` // 是否上架
	ShipFree        bool     `gorm:"type:bool;default:false;not null"` // 是否包邮
	IsNew           bool     `gorm:"type:bool;default:false;not null"` // 是否新品
	IsHot           bool     `gorm:"type:bool;default:false;not null"` // 是否热销
	Name            string   `gorm:"type:varchar(50);not null"`        // 商品名称
	GoodsSn         string   `gorm:"type:varchar(50);not null"`        // 商品编号
	ClickNum        int32    `gorm:"type:int;default:0;not null"`      // 点击数
	SoldNum         int32    `gorm:"type:int;default:0;not null"`      // 销售数量
	FavNum          int32    `gorm:"type:int;default:0;not null"`      // 收藏数量
	MarketPrice     float32  `gorm:"not null"`                         // 市场价
	ShopPrice       float32  `gorm:"not null"`                         // 本店价
	GoodsBrief      string   `gorm:"type:varchar(100);not null"`       // 商品简介
	Images          GormList `gorm:"type:varchar(1000);not null"`      // 商品图片
	DescImages      GormList `gorm:"type:varchar(1000);not null"`      // 商品描述图片
	GoodsFrontImage string   `gorm:"type:varchar(200);not null"`       // 商品封面图片
}
