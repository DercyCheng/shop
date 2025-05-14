package entity

import (
	"time"
)

// Category 商品分类实体
type Category struct {
	ID               int64      `json:"id"`
	Name             string     `json:"name"`
	ParentCategoryID int64      `json:"parent_category_id"`
	Level            int        `json:"level"`
	IsTab            bool       `json:"is_tab"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
	
	// 关联实体
	ParentCategory *Category   `json:"parent_category,omitempty"`
	SubCategories  []*Category `json:"sub_categories,omitempty"`
	Brands         []*Brand    `json:"brands,omitempty"`
}

// IsRootCategory 判断是否为根分类
func (c *Category) IsRootCategory() bool {
	return c.ParentCategoryID == 0
}

// IsLeafCategory 判断是否为叶子分类（无子分类）
func (c *Category) IsLeafCategory() bool {
	return len(c.SubCategories) == 0
}

// AddSubCategory 添加子分类
func (c *Category) AddSubCategory(sub *Category) {
	if c.SubCategories == nil {
		c.SubCategories = make([]*Category, 0)
	}
	sub.ParentCategoryID = c.ID
	sub.Level = c.Level + 1
	sub.UpdatedAt = time.Now()
	c.SubCategories = append(c.SubCategories, sub)
}

// RemoveSubCategory 移除子分类
func (c *Category) RemoveSubCategory(categoryID int64) bool {
	for i, sub := range c.SubCategories {
		if sub.ID == categoryID {
			// 移除指定子分类
			c.SubCategories = append(c.SubCategories[:i], c.SubCategories[i+1:]...)
			return true
		}
	}
	return false
}

// CategoryBrand 分类品牌关系实体
type CategoryBrand struct {
	ID         int64     `json:"id"`
	CategoryID int64     `json:"category_id"`
	BrandID    int64     `json:"brand_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	
	// 关联实体
	Category *Category `json:"category,omitempty"`
	Brand    *Brand    `json:"brand,omitempty"`
}
