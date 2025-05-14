package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	
	"shop/backend/product/internal/domain/entity"
	
	"github.com/olivere/elastic/v7"
)

const (
	// ElasticSearchProductIndex 商品索引名
	ElasticSearchProductIndex = "shop_products"
	
	// ElasticSearchBatchSize 批量操作大小
	ElasticSearchBatchSize = 100
)

// ElasticSearchProductDoc 商品搜索文档结构
type ElasticSearchProductDoc struct {
	ID              int64     `json:"id"`
	CategoryID      int64     `json:"category_id"`
	CategoryName    string    `json:"category_name"`
	BrandID         int64     `json:"brand_id"`
	BrandName       string    `json:"brand_name"`
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
	Keywords        []string  `json:"keywords"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ElasticSearchRepository ElasticSearch仓储实现
type ElasticSearchRepository struct {
	client          *elastic.Client
	productRepo     ProductRepository
	categoryRepo    CategoryRepository
	brandRepo       BrandRepository
	indexName       string
	indexDefinition string
}

// NewElasticSearchRepository 创建ElasticSearch仓储实例
func NewElasticSearchRepository(
	client *elastic.Client,
	productRepo ProductRepository,
	categoryRepo CategoryRepository,
	brandRepo BrandRepository,
) SearchRepository {
	return &ElasticSearchRepository{
		client:       client,
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		brandRepo:    brandRepo,
		indexName:    ElasticSearchProductIndex,
		indexDefinition: `
{
  "settings": {
    "number_of_shards": 3,
    "number_of_replicas": 1,
    "analysis": {
      "analyzer": {
        "ik_smart_pinyin": {
          "type": "custom",
          "tokenizer": "ik_smart",
          "filter": ["pinyin_filter"]
        }
      },
      "filter": {
        "pinyin_filter": {
          "type": "pinyin",
          "keep_original": true,
          "keep_full_pinyin": true,
          "keep_joined_full_pinyin": true,
          "keep_first_letter": true,
          "keep_separate_first_letter": true
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "id": { "type": "long" },
      "category_id": { "type": "long" },
      "category_name": { 
        "type": "text", 
        "analyzer": "ik_smart_pinyin",
        "search_analyzer": "ik_smart" 
      },
      "brand_id": { "type": "long" },
      "brand_name": { 
        "type": "text", 
        "analyzer": "ik_smart_pinyin",
        "search_analyzer": "ik_smart"
      },
      "on_sale": { "type": "boolean" },
      "ship_free": { "type": "boolean" },
      "is_new": { "type": "boolean" },
      "is_hot": { "type": "boolean" },
      "name": { 
        "type": "text", 
        "analyzer": "ik_smart_pinyin",
        "search_analyzer": "ik_smart",
        "fields": {
          "keyword": { "type": "keyword" }
        }
      },
      "goods_sn": { "type": "keyword" },
      "click_num": { "type": "integer" },
      "sold_num": { "type": "integer" },
      "fav_num": { "type": "integer" },
      "market_price": { "type": "float" },
      "shop_price": { "type": "float" },
      "goods_brief": { 
        "type": "text", 
        "analyzer": "ik_smart_pinyin",
        "search_analyzer": "ik_smart" 
      },
      "goods_desc": { 
        "type": "text", 
        "analyzer": "ik_smart_pinyin",
        "search_analyzer": "ik_smart" 
      },
      "goods_front_image": { "type": "keyword", "index": false },
      "keywords": { 
        "type": "text", 
        "analyzer": "ik_smart_pinyin",
        "search_analyzer": "ik_smart" 
      },
      "created_at": { "type": "date" },
      "updated_at": { "type": "date" }
    }
  }
}`,
	}
}

// Init 初始化ElasticSearch索引
func (r *ElasticSearchRepository) Init(ctx context.Context) error {
	exists, err := r.client.IndexExists(r.indexName).Do(ctx)
	if err != nil {
		return err
	}
	
	if !exists {
		_, err := r.client.CreateIndex(r.indexName).
			Body(r.indexDefinition).
			Do(ctx)
		if err != nil {
			return err
		}
	}
	
	return nil
}

// convertProductToDoc 将商品实体转换为搜索文档
func (r *ElasticSearchRepository) convertProductToDoc(ctx context.Context, product *entity.Product) (*ElasticSearchProductDoc, error) {
	if product == nil {
		return nil, nil
	}
	
	doc := &ElasticSearchProductDoc{
		ID:              product.ID,
		CategoryID:      product.CategoryID,
		BrandID:         product.BrandsID,
		OnSale:          product.OnSale,
		ShipFree:        product.ShipFree,
		IsNew:           product.IsNew,
		IsHot:           product.IsHot,
		Name:            product.Name,
		GoodsSN:         product.GoodsSN,
		ClickNum:        product.ClickNum,
		SoldNum:         product.SoldNum,
		FavNum:          product.FavNum,
		MarketPrice:     product.MarketPrice,
		ShopPrice:       product.ShopPrice,
		GoodsBrief:      product.GoodsBrief,
		GoodsDesc:       product.GoodsDesc,
		GoodsFrontImage: product.GoodsFrontImage,
		CreatedAt:       product.CreatedAt,
		UpdatedAt:       product.UpdatedAt,
		Keywords:        generateKeywords(product),
	}
	
	// 获取分类名称
	if category, err := r.categoryRepo.GetCategoryByID(ctx, product.CategoryID); err == nil && category != nil {
		doc.CategoryName = category.Name
	}
	
	// 获取品牌名称
	if brand, err := r.brandRepo.GetBrandByID(ctx, product.BrandsID); err == nil && brand != nil {
		doc.BrandName = brand.Name
	}
	
	return doc, nil
}

// generateKeywords 生成商品关键词
func generateKeywords(product *entity.Product) []string {
	keywords := make([]string, 0)
	
	// 添加商品名称分词
	keywords = append(keywords, product.Name)
	
	// 商品简介作为关键词
	if product.GoodsBrief != "" {
		keywords = append(keywords, product.GoodsBrief)
	}
	
	// 分类和品牌名称也可以作为关键词，但需要从仓储获取
	// 这里已在convertProductToDoc中处理
	
	return keywords
}

// SearchProducts 搜索商品
func (r *ElasticSearchRepository) SearchProducts(ctx context.Context, params SearchParams) (*SearchResult, error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	
	if params.PageSize <= 0 {
		params.PageSize = 10
	}
	
	// 构建查询
	query := elastic.NewBoolQuery()
	
	// 关键词搜索
	if params.Keyword != "" {
		multiMatchQuery := elastic.NewMultiMatchQuery(params.Keyword,
			"name^3", // name字段权重高
			"goods_brief^2",
			"category_name",
			"brand_name",
			"goods_desc",
			"keywords",
		).Type("best_fields").TieBreaker(0.3)
		
		query = query.Must(multiMatchQuery)
	}
	
	// 分类过滤
	if params.CategoryID > 0 {
		query = query.Filter(elastic.NewTermQuery("category_id", params.CategoryID))
	}
	
	// 品牌过滤
	if params.BrandID > 0 {
		query = query.Filter(elastic.NewTermQuery("brand_id", params.BrandID))
	}
	
	// 价格区间过滤
	if params.PriceMin > 0 || params.PriceMax > 0 {
		rangeQuery := elastic.NewRangeQuery("shop_price")
		if params.PriceMin > 0 {
			rangeQuery = rangeQuery.Gte(params.PriceMin)
		}
		if params.PriceMax > 0 {
			rangeQuery = rangeQuery.Lte(params.PriceMax)
		}
		query = query.Filter(rangeQuery)
	}
	
	// 上架状态过滤
	if params.OnSale {
		query = query.Filter(elastic.NewTermQuery("on_sale", true))
	}
	
	// 其他条件过滤
	if params.IsNew {
		query = query.Filter(elastic.NewTermQuery("is_new", true))
	}
	
	if params.IsHot {
		query = query.Filter(elastic.NewTermQuery("is_hot", true))
	}
	
	if params.ShipFree {
		query = query.Filter(elastic.NewTermQuery("ship_free", true))
	}
	
	// 构建排序
	sorters := make([]elastic.Sorter, 0)
	if params.OrderBy != "" {
		parts := strings.Split(params.OrderBy, ":")
		field := parts[0]
		order := "desc" // 默认降序
		
		if len(parts) > 1 {
			order = strings.ToLower(parts[1])
		}
		
		// 允许排序的字段
		allowedSortFields := map[string]bool{
			"shop_price": true,
			"sold_num":   true,
			"click_num":  true,
			"fav_num":    true,
			"created_at": true,
		}
		
		if allowedSortFields[field] {
			sorter := elastic.NewFieldSort(field)
			if order == "asc" {
				sorter = sorter.Asc()
			} else {
				sorter = sorter.Desc()
			}
			sorters = append(sorters, sorter)
		}
	} else {
		// 默认排序
		sorters = append(sorters, elastic.NewFieldSort("updated_at").Desc())
	}
	
	// 构建搜索请求
	searchService := r.client.Search().
		Index(r.indexName).
		Query(query).
		SortBy(sorters...).
		From((params.Page - 1) * params.PageSize).
		Size(params.PageSize)
	
	// 执行搜索
	response, err := searchService.Do(ctx)
	if err != nil {
		return nil, err
	}
	
	// 解析搜索结果
	result := &SearchResult{
		Total: response.TotalHits(),
		Page:  params.Page,
		Size:  params.PageSize,
		Pages: int(response.TotalHits() / int64(params.PageSize)),
		Goods: make([]*entity.Product, 0),
	}
	
	if result.Pages*params.PageSize < int(response.TotalHits()) {
		result.Pages++
	}
	
	// 提取商品IDs
	productIDs := make([]int64, 0, len(response.Hits.Hits))
	for _, hit := range response.Hits.Hits {
		id, err := strconv.ParseInt(hit.Id, 10, 64)
		if err != nil {
			continue
		}
		productIDs = append(productIDs, id)
	}
	
	// 批量获取商品详情
	if len(productIDs) > 0 {
		products, err := r.productRepo.BatchGetProducts(ctx, productIDs)
		if err != nil {
			return nil, err
		}
		
		// 按照搜索结果的顺序组织商品
		productMap := make(map[int64]*entity.Product, len(products))
		for _, product := range products {
			productMap[product.ID] = product
		}
		
		for _, id := range productIDs {
			if product, ok := productMap[id]; ok {
				result.Goods = append(result.Goods, product)
			}
		}
	}
	
	return result, nil
}

// IndexProduct 索引单个商品
func (r *ElasticSearchRepository) IndexProduct(ctx context.Context, product *entity.Product) error {
	doc, err := r.convertProductToDoc(ctx, product)
	if err != nil {
		return err
	}
	
	// 商品被删除的情况下，需要从索引中删除
	if product.IsDeleted {
		_, err = r.client.Delete().
			Index(r.indexName).
			Id(strconv.FormatInt(product.ID, 10)).
			Refresh("true").
			Do(ctx)
		return err
	}
	
	// 更新或创建索引
	_, err = r.client.Index().
		Index(r.indexName).
		Id(strconv.FormatInt(product.ID, 10)).
		BodyJson(doc).
		Refresh("true").
		Do(ctx)
		
	return err
}

// BatchIndexProducts 批量索引商品
func (r *ElasticSearchRepository) BatchIndexProducts(ctx context.Context, products []*entity.Product) error {
	bulkRequest := r.client.Bulk()
	count := 0
	
	for _, product := range products {
		doc, err := r.convertProductToDoc(ctx, product)
		if err != nil {
			continue
		}
		
		if product.IsDeleted {
			// 删除操作
			req := elastic.NewBulkDeleteRequest().
				Index(r.indexName).
				Id(strconv.FormatInt(product.ID, 10))
			bulkRequest = bulkRequest.Add(req)
		} else {
			// 索引操作
			req := elastic.NewBulkIndexRequest().
				Index(r.indexName).
				Id(strconv.FormatInt(product.ID, 10)).
				Doc(doc)
			bulkRequest = bulkRequest.Add(req)
		}
		
		count++
		
		// 每批次处理一次
		if count >= ElasticSearchBatchSize {
			_, err := bulkRequest.Do(ctx)
			if err != nil {
				return err
			}
			bulkRequest = r.client.Bulk()
			count = 0
		}
	}
	
	// 处理剩余的请求
	if count > 0 {
		_, err := bulkRequest.Refresh("true").Do(ctx)
		if err != nil {
			return err
		}
	}
	
	return nil
}

// DeleteProductIndex 删除商品索引
func (r *ElasticSearchRepository) DeleteProductIndex(ctx context.Context, id int64) error {
	_, err := r.client.Delete().
		Index(r.indexName).
		Id(strconv.FormatInt(id, 10)).
		Refresh("true").
		Do(ctx)
		
	return err
}

// SyncProductIndex 同步所有商品到搜索引擎
func (r *ElasticSearchRepository) SyncProductIndex(ctx context.Context) error {
	// 分页获取所有商品
	page := 1
	pageSize := 100
	
	for {
		filter := ProductFilter{
			Page:     page,
			PageSize: pageSize,
		}
		
		products, total, err := r.productRepo.ListProducts(ctx, filter)
		if err != nil {
			return err
		}
		
		if len(products) == 0 {
			break
		}
		
		// 批量索引
		if err := r.BatchIndexProducts(ctx, products); err != nil {
			return err
		}
		
		// 判断是否已处理完所有商品
		if int64(page*pageSize) >= total {
			break
		}
		
		page++
	}
	
	return nil
}
