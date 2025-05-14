package entity

import (
	"time"
)

// Product 商品实体
type Product struct {
	ID              int64     `json:"id"`
	CategoryID      int64     `json:"category_id"`
	BrandsID        int64     `json:"brands_id"`
	OnSale          bool      `json:"on_sale"`
	ShipFree        bool      `json:"ship_free"`
	IsNew           bool      `json:"is_new"`
	IsHot           bool      `json:"is_hot"`
	Name            string    `json:"name"`
	GoodsSN         string    `json:"goods_sn"`
	ClickNum        int       `json:"click_num"`
	SoldNum         int       `json:"sold_num"`
	FavNum          int       `json:"fav_num"`
	MarketPrice     float64   `json:"market_price"`
	ShopPrice       float64   `json:"shop_price"`
	GoodsBrief      string    `json:"goods_brief"`
	GoodsDesc       string    `json:"goods_desc"`
	GoodsFrontImage string    `json:"goods_front_image"`
	IsDeleted       bool      `json:"is_deleted"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
	
	// 关联实体
	Category *Category `json:"category,omitempty"`
	Brand    *Brand    `json:"brand,omitempty"`
	Images   []string  `json:"images,omitempty"`
	
	// SKU相关
	SkuList []*ProductSKU `json:"sku_list,omitempty"`
}

// ProductSKU 商品SKU实体
type ProductSKU struct {
	ID             int64     `json:"id"`
	ProductID      int64     `json:"product_id"`
	SkuName        string    `json:"sku_name"`
	SkuCode        string    `json:"sku_code"`
	BarCode        string    `json:"bar_code"`
	Price          float64   `json:"price"`
	PromotionPrice float64   `json:"promotion_price"`
	Points         int       `json:"points"`
	Stocks         int       `json:"stocks"`
	Image          string    `json:"image"`
	OriginalStock  int       `json:"original_stock"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	
	// 规格值，如 {"颜色": "红色", "尺寸": "XL"}
	SpecValues map[string]string `json:"spec_values,omitempty"`
}

// ProductAttribute 商品属性实体
type ProductAttribute struct {
	ID         int64     `json:"id"`
	ProductID  int64     `json:"product_id"`
	AttrName   string    `json:"attr_name"`
	AttrValue  string    `json:"attr_value"`
	AttrSort   int       `json:"attr_sort"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ProductImage 商品图片实体
type ProductImage struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	ImageURL  string    `json:"image_url"`
	IsMain    bool      `json:"is_main"`
	Sort      int       `json:"sort"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ProductSpec 商品规格定义实体
type ProductSpec struct {
	ID         int64     `json:"id"`
	ProductID  int64     `json:"product_id"`
	SpecName   string    `json:"spec_name"`       // 规格名，如：颜色、尺寸
	SpecValues []string  `json:"spec_values"`     // 规格值列表，如：["红色", "蓝色", "绿色"]
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// IncreaseClickNum 增加商品点击数
func (p *Product) IncreaseClickNum() {
	p.ClickNum++
	p.UpdatedAt = time.Now()
}

// IncreaseSoldNum 增加商品销量
func (p *Product) IncreaseSoldNum(count int) {
	p.SoldNum += count
	p.UpdatedAt = time.Now()
}

// IncreaseFavNum 增加收藏数
func (p *Product) IncreaseFavNum() {
	p.FavNum++
	p.UpdatedAt = time.Now()
}

// DecreaseFavNum 减少收藏数
func (p *Product) DecreaseFavNum() {
	if p.FavNum > 0 {
		p.FavNum--
		p.UpdatedAt = time.Now()
	}
}

// SetOnSale 设置商品上下架状态
func (p *Product) SetOnSale(onSale bool) {
	p.OnSale = onSale
	p.UpdatedAt = time.Now()
}

// SetDeleteStatus 设置删除状态
func (p *Product) SetDeleteStatus(isDeleted bool) {
	p.IsDeleted = isDeleted
	p.UpdatedAt = time.Now()
	if isDeleted {
		now := time.Now()
		p.DeletedAt = &now
	} else {
		p.DeletedAt = nil
	}
}
