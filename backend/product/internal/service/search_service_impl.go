package service

import (
	"context"
	"time"
	
	"shop/backend/product/internal/domain/entity"
)

var (
	ErrSearchFailed = "search operation failed"
)

// SearchServiceImpl 搜索服务实现
type SearchServiceImpl struct {
	searchRepo SearchRepository
	productRepo ProductRepository
}

// NewSearchService 创建搜索服务实例
func NewSearchService(
	searchRepo SearchRepository,
	productRepo ProductRepository,
) SearchService {
	return &SearchServiceImpl{
		searchRepo:  searchRepo,
		productRepo: productRepo,
	}
}

// SearchProducts 搜索商品
func (s *SearchServiceImpl) SearchProducts(ctx context.Context, params *SearchParams) (*SearchResult, error) {
	// 记录热搜词（此功能可以异步处理）
	if params.Keyword != "" {
		go s.recordHotKeyword(context.Background(), params.Keyword)
	}
	
	// 设置默认值
	if params.Page <= 0 {
		params.Page = 1
	}
	
	if params.PageSize <= 0 {
		params.PageSize = 10
	}
	
	// 使用搜索仓储执行搜索
	results, err := s.searchRepo.SearchProducts(ctx, *params)
	if err != nil {
		return nil, err
	}
	
	return results, nil
}

// recordHotKeyword 记录热搜词（简化实现，实际项目中可能需要存储到数据库）
func (s *SearchServiceImpl) recordHotKeyword(ctx context.Context, keyword string) {
	// 实际实现中，这里应该更新热搜词统计
	// 例如，将关键词及其搜索次数存入数据库或Redis
}

// IndexProduct 将商品索引到搜索引擎
func (s *SearchServiceImpl) IndexProduct(ctx context.Context, product *entity.Product) error {
	return s.searchRepo.IndexProduct(ctx, product)
}

// BatchIndexProducts 批量索引商品
func (s *SearchServiceImpl) BatchIndexProducts(ctx context.Context, products []*entity.Product) error {
	return s.searchRepo.BatchIndexProducts(ctx, products)
}

// DeleteProductIndex 从搜索引擎中删除商品索引
func (s *SearchServiceImpl) DeleteProductIndex(ctx context.Context, id int64) error {
	return s.searchRepo.DeleteProductIndex(ctx, id)
}

// SyncProductIndex 同步所有商品到搜索引擎
func (s *SearchServiceImpl) SyncProductIndex(ctx context.Context) error {
	startTime := time.Now()
	err := s.searchRepo.SyncProductIndex(ctx)
	
	// 记录同步时间和状态（实际实现中可以记录到数据库）
	_ = time.Since(startTime)
	
	return err
}
