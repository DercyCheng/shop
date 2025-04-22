package model

import (
	"context"
	"goods_srv/global"
	"gorm.io/gorm"
	"strconv"
)

/*
类型，这个字段是否能为null，这个字段应该设置为可以为null还是设置为空，0
实际开发过程中尽量设置为不为null
https://zhuanlan.zhihu.com/p/73997266
这些类型我们使用int32还是int
*/
type Category struct {
	BaseModel
	Name  string `gorm:"type:varchar(20);not null;comment:商品名称"`
	Level int32  `gorm:"type:int;not null;default:1;comment:分类级别" json:"level"`
	IsTab bool   `gorm:"not null;default:false;comment:是否在标题栏" json:"is_tab"`
	//外键指向父类别
	ParentCategoryID int32     `json:"parent"`
	ParentCategory   *Category `json:"-"`
	//通过外键和指向的外键反向查询
	SubCategory []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"`
}
type Brands struct {
	BaseModel
	Name string `gorm:"type:varchar(20);not null;comment:品牌名称"`
	Logo string `gorm:"type:varchar(200);default:'';not null;comment:品牌logo"`
}

// 商品分类&品牌  外键
type GoodsCategoryBrand struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand;unique"`
	Category   Category

	BrandsID int32 `gorm:"type:int;index:idx_category_brand;unique"`
	Brands   Brands
}

// 重载多对多表名
func (GoodsCategoryBrand) TableName() string {
	return "goodscategorybrand"
}

type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(200);not null;comment:轮播图"`
	Url   string `gorm:"type:varchar(200);not null;comment:跳转地址"`
	Index int32  `gorm:"type:int;default:1;not null;comment:等级"`
}
type Goods struct {
	BaseModel
	//商品分类&品牌外键
	CategoryID int32 `gorm:"type:int;not null"`
	Category   Category
	BrandsID   int32 `gorm:"type:int;not null"`
	Brands     Brands

	OnSale   bool `gorm:"default:false;not null;comment:是否上架"`
	ShipFree bool `gorm:"default:false;not null;comment:是否免运费"`
	IsNew    bool `gorm:"default:false;not null;comment:是否最新"`
	IsHot    bool `gorm:"default:false;not null;comment:是否热门"`

	Name            string   `gorm:"type:varchar(50);not null;comment:商品名称"`
	GoodsSn         string   `gorm:"type:varchar(50);not null;comment:商品编号"`
	ClickNum        int32    `gorm:"type:int;default:0;not null;comment:点击数"`
	SoldNum         int32    `gorm:"type:int;default:0;not null;comment:销量"`
	FavNum          int32    `gorm:"type:int;default:0;not null;comment:收藏数"`
	MarketPrice     float32  `gorm:"not null;comment:市场价格"`
	ShopPrice       float32  `gorm:"not null;comment:本店价格"`
	GoodsBrief      string   `gorm:"type:varchar(100);not null;comment:商品简介"`
	Images          GormList `gorm:"type:varchar(1000);not null;comment:商品内轮播图"`
	DescImages      GormList `gorm:"type:varchar(1000);not null;comment:详情图片"`
	GoodsFrontImage string   `gorm:"type:varchar(100);not null;comment:封面图"`
}
type GoodsImages struct {
	GoodsID int
	Image   string
}

func (g *Goods) AfterCreate(tx *gorm.DB) (err error) {
	esModel := EsGoods{
		ID:          g.ID,
		CategoryID:  g.CategoryID,
		BrandsID:    g.BrandsID,
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		GoodsBrief:  g.GoodsBrief,
		ShopPrice:   g.ShopPrice,
	}
	_, err = global.EsClient.Index().Index(esModel.GetIndexName()).BodyJson(esModel).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
func (g *Goods) AfterUpdate(tx *gorm.DB) (err error) {
	esModel := EsGoods{
		ID:          g.ID,
		CategoryID:  g.CategoryID,
		BrandsID:    g.BrandsID,
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		GoodsBrief:  g.GoodsBrief,
		ShopPrice:   g.ShopPrice,
	}
	_, err = global.EsClient.Update().Index(esModel.GetIndexName()).Doc(esModel).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
func (g *Goods) AfterDelete(tx *gorm.DB) (err error) {
	_, err = global.EsClient.Delete().Index(EsGoods{}.GetIndexName()).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
