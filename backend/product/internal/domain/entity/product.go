package entity

import (
	"time"
)

// Product represents a product entity in the system
type Product struct {
	ID              int64      `gorm:"primaryKey;column:id"`
	CategoryID      int64      `gorm:"column:category_id;not null;index:idx_category_id"`
	BrandID         int64      `gorm:"column:brands_id;not null;index:idx_brands_id"`
	Name            string     `gorm:"column:name;type:varchar(100);not null"`
	GoodsSN         string     `gorm:"column:goods_sn;type:varchar(50)"`
	OnSale          bool       `gorm:"column:on_sale;default:1"`
	ShipFree        bool       `gorm:"column:ship_free;default:1"`
	IsNew           bool       `gorm:"column:is_new;default:0"`
	IsHot           bool       `gorm:"column:is_hot;default:0"`
	ClickNum        int        `gorm:"column:click_num;default:0"`
	SoldNum         int        `gorm:"column:sold_num;default:0"`
	FavNum          int        `gorm:"column:fav_num;default:0"`
	MarketPrice     float64    `gorm:"column:market_price;default:0"`
	ShopPrice       float64    `gorm:"column:shop_price;default:0"`
	GoodsBrief      string     `gorm:"column:goods_brief;type:varchar(255)"`
	GoodsDesc       string     `gorm:"column:goods_desc;type:text"`
	GoodsFrontImage string     `gorm:"column:goods_front_image;type:varchar(255)"`
	IsDeleted       bool       `gorm:"column:is_deleted;default:0"`
	CreatedAt       time.Time  `gorm:"column:created_at"`
	UpdatedAt       time.Time  `gorm:"column:updated_at"`
	DeletedAt       *time.Time `gorm:"column:deleted_at;index"`

	// Relations
	Category *Category `gorm:"-"`
	Brand    *Brand    `gorm:"-"`
	Images   []string  `gorm:"-"`
}

// TableName returns the table name for the Product entity
func (Product) TableName() string {
	return "goods"
}

// ProductImage represents a product image in the system
type ProductImage struct {
	ID        int64      `gorm:"primaryKey;column:id"`
	GoodsID   int64      `gorm:"column:goods_id;not null;index"`
	Image     string     `gorm:"column:image;type:varchar(255);not null"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index"`
}

// TableName returns the table name for the ProductImage entity
func (ProductImage) TableName() string {
	return "goods_image"
}
