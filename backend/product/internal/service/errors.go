package service

import (
	"errors"
)

var (
	// ErrProductNotFound 商品未找到错误
	ErrProductNotFound = errors.New("product not found")
	
	// ErrCategoryNotFound 分类未找到错误
	ErrCategoryNotFound = errors.New("category not found")
	
	// ErrBrandNotFound 品牌未找到错误
	ErrBrandNotFound = errors.New("brand not found")
	
	// ErrBannerNotFound 轮播图未找到错误
	ErrBannerNotFound = errors.New("banner not found")
	
	// ErrDatabaseOperation 数据库操作错误
	ErrDatabaseOperation = errors.New("database operation failed")
	
	// ErrInvalidParameter 参数无效错误
	ErrInvalidParameter = errors.New("invalid parameter")
	
	// ErrSearchEngine 搜索引擎错误
	ErrSearchEngine = errors.New("search engine operation failed")
	
	// ErrProductExists 商品已存在错误
	ErrProductExists = errors.New("product already exists")
	
	// ErrCategoryExists 分类已存在错误
	ErrCategoryExists = errors.New("category already exists")
	
	// ErrBrandExists 品牌已存在错误
	ErrBrandExists = errors.New("brand already exists")
	
	// ErrCategoryHasChildren 分类有子分类错误
	ErrCategoryHasChildren = errors.New("category has children")
	
	// ErrCategoryHasProducts 分类有商品错误
	ErrCategoryHasProducts = errors.New("category has products")
	
	// ErrBrandHasProducts 品牌有商品错误
	ErrBrandHasProducts = errors.New("brand has products")
)
