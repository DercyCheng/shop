package entity

import "time"

// UserFav represents a user's favorite product
type UserFav struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	GoodsID   int64     `json:"goods_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GoodsInfo represents basic information about a product
// This is populated from the Product service
type GoodsInfo struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	ShopPrice    string `json:"shop_price"`
	Image        string `json:"image"`
	CategoryName string `json:"category_name"`
	BrandName    string `json:"brand_name"`
}

// UserFavWithGoods combines user favorite info with goods details
type UserFavWithGoods struct {
	UserFav
	GoodsInfo *GoodsInfo `json:"goods_info"`
}
