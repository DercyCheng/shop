package grpc

import (
	"context"
	
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	
	"shop/backend/product/api/proto"
	"shop/backend/product/internal/domain/entity"
	"shop/backend/product/internal/service"
)

// ProductHandler 商品服务gRPC处理器
type ProductHandler struct {
	proto.UnimplementedProductServiceServer
	productService  service.ProductService
	categoryService service.CategoryService
	brandService    service.BrandService
	bannerService   service.BannerService
	searchService   service.SearchService
}

// NewProductHandler 创建商品服务gRPC处理器
func NewProductHandler(
	productService service.ProductService,
	categoryService service.CategoryService,
	brandService service.BrandService,
	bannerService service.BannerService,
	searchService service.SearchService,
) *ProductHandler {
	return &ProductHandler{
		productService:  productService,
		categoryService: categoryService,
		brandService:    brandService,
		bannerService:   bannerService,
		searchService:   searchService,
	}
}

// GoodsList 获取商品列表
func (h *ProductHandler) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	// 构建过滤条件
	filter := service.ProductFilter{
		Page:       int(req.Page),
		PageSize:   int(req.PageSize),
		PriceMin:   floatPtr(float64(req.PriceMin)),
		PriceMax:   floatPtr(float64(req.PriceMax)),
		CategoryID: req.CategoryId,
		BrandID:    req.BrandId,
		IsHot:      boolPtr(req.IsHot),
		IsNew:      boolPtr(req.IsNew),
		OrderBy:    req.OrderBy,
	}
	
	// 关键词搜索
	if req.Keywords != "" {
		filter.Name = req.Keywords
	}
	
	// 获取商品列表
	products, total, err := h.productService.ListProducts(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取商品列表失败: %v", err)
	}
	
	// 转换为响应格式
	goodsList := make([]*proto.GoodsInfoResponse, 0, len(products))
	for _, product := range products {
		goodsList = append(goodsList, convertProductToProto(product))
	}
	
	return &proto.GoodsListResponse{
		Total: total,
		Goods: goodsList,
	}, nil
}

// BatchGetGoods 批量获取商品信息
func (h *ProductHandler) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	// 批量获取商品
	products, err := h.productService.BatchGetProducts(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "批量获取商品信息失败: %v", err)
	}
	
	// 转换为响应格式
	goodsList := make([]*proto.GoodsInfoResponse, 0, len(products))
	for _, product := range products {
		goodsList = append(goodsList, convertProductToProto(product))
	}
	
	return &proto.GoodsListResponse{
		Total: int64(len(goodsList)),
		Goods: goodsList,
	}, nil
}

// GetGoodsDetail 获取商品详情
func (h *ProductHandler) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	// 获取商品详情
	product, err := h.productService.GetProductByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取商品详情失败: %v", err)
	}
	
	if product == nil {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	
	// 转换为响应格式
	return convertProductToProto(product), nil
}

// CreateGoods 创建商品
func (h *ProductHandler) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	// 构建商品实体
	product := &entity.Product{
		Name:            req.Name,
		GoodsSN:         req.GoodsSn,
		CategoryID:      req.CategoryId,
		BrandsID:        req.BrandId,
		MarketPrice:     req.MarketPrice,
		ShopPrice:       req.ShopPrice,
		GoodsBrief:      req.GoodsBrief,
		GoodsDesc:       req.GoodsDesc,
		ShipFree:        req.ShipFree,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		OnSale:          req.OnSale,
		GoodsFrontImage: req.GoodsFrontImage,
	}
	
	// 创建商品
	createdProduct, err := h.productService.CreateProduct(ctx, product, nil, nil, nil, req.Images)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "创建商品失败: %v", err)
	}
	
	// 转换为响应格式
	return convertProductToProto(createdProduct), nil
}

// UpdateGoods 更新商品
func (h *ProductHandler) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	// 检查商品是否存在
	existingProduct, err := h.productService.GetProductByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询商品失败: %v", err)
	}
	
	if existingProduct == nil {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	
	// 更新商品信息
	existingProduct.Name = req.Name
	existingProduct.GoodsSN = req.GoodsSn
	existingProduct.CategoryID = req.CategoryId
	existingProduct.BrandsID = req.BrandId
	existingProduct.MarketPrice = req.MarketPrice
	existingProduct.ShopPrice = req.ShopPrice
	existingProduct.GoodsBrief = req.GoodsBrief
	existingProduct.GoodsDesc = req.GoodsDesc
	existingProduct.ShipFree = req.ShipFree
	existingProduct.IsNew = req.IsNew
	existingProduct.IsHot = req.IsHot
	existingProduct.OnSale = req.OnSale
	existingProduct.GoodsFrontImage = req.GoodsFrontImage
	
	// 更新商品
	if err := h.productService.UpdateProduct(ctx, existingProduct, nil, nil, nil, req.Images); err != nil {
		return nil, status.Errorf(codes.Internal, "更新商品失败: %v", err)
	}
	
	return &emptypb.Empty{}, nil
}

// DeleteGoods 删除商品
func (h *ProductHandler) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	// 删除商品
	if err := h.productService.DeleteProduct(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "删除商品失败: %v", err)
	}
	
	return &emptypb.Empty{}, nil
}

// GetAllCategorysList 获取所有分类
func (h *ProductHandler) GetAllCategorysList(ctx context.Context, _ *emptypb.Empty) (*proto.CategoryListResponse, error) {
	// 获取所有分类
	categories, err := h.categoryService.GetAllCategories(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取分类列表失败: %v", err)
	}
	
	// 转换为响应格式
	categoryList := make([]*proto.CategoryInfoResponse, 0, len(categories))
	for _, category := range categories {
		categoryList = append(categoryList, convertCategoryToProto(category))
	}
	
	return &proto.CategoryListResponse{
		Total: int64(len(categoryList)),
		Data:  categoryList,
	}, nil
}

// GetSubCategory 获取子分类
func (h *ProductHandler) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	// 获取分类详情
	category, err := h.categoryService.GetCategoryByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取分类详情失败: %v", err)
	}
	
	if category == nil {
		return nil, status.Errorf(codes.NotFound, "分类不存在")
	}
	
	// 获取子分类
	subCategories, err := h.categoryService.GetSubCategories(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取子分类失败: %v", err)
	}
	
	// 转换为响应格式
	subCategoryList := make([]*proto.CategoryInfoResponse, 0, len(subCategories))
	for _, subCategory := range subCategories {
		subCategoryList = append(subCategoryList, convertCategoryToProto(subCategory))
	}
	
	return &proto.SubCategoryListResponse{
		Total:        int64(len(subCategoryList)),
		Info:         convertCategoryToProto(category),
		SubCategories: subCategoryList,
	}, nil
}

// CreateCategory 创建分类
func (h *ProductHandler) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	// 构建分类实体
	category := &entity.Category{
		Name:             req.Name,
		ParentCategoryID: req.ParentCategoryId,
		Level:            int(req.Level),
		IsTab:            req.IsTab,
	}
	
	// 创建分类
	createdCategory, err := h.categoryService.CreateCategory(ctx, category)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "创建分类失败: %v", err)
	}
	
	// 转换为响应格式
	return convertCategoryToProto(createdCategory), nil
}

// UpdateCategory 更新分类
func (h *ProductHandler) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	// 检查分类是否存在
	existingCategory, err := h.categoryService.GetCategoryByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询分类失败: %v", err)
	}
	
	if existingCategory == nil {
		return nil, status.Errorf(codes.NotFound, "分类不存在")
	}
	
	// 更新分类信息
	existingCategory.Name = req.Name
	existingCategory.ParentCategoryID = req.ParentCategoryId
	existingCategory.Level = int(req.Level)
	existingCategory.IsTab = req.IsTab
	
	// 更新分类
	if err := h.categoryService.UpdateCategory(ctx, existingCategory); err != nil {
		return nil, status.Errorf(codes.Internal, "更新分类失败: %v", err)
	}
	
	return &emptypb.Empty{}, nil
}

// DeleteCategory 删除分类
func (h *ProductHandler) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryInfo) (*emptypb.Empty, error) {
	// 删除分类
	if err := h.categoryService.DeleteCategory(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "删除分类失败: %v", err)
	}
	
	return &emptypb.Empty{}, nil
}

// BrandList 获取品牌列表
func (h *ProductHandler) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	// 构建过滤条件
	filter := service.BrandFilter{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
	}
	
	// 获取品牌列表
	brands, total, err := h.brandService.ListBrands(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取品牌列表失败: %v", err)
	}
	
	// 转换为响应格式
	brandList := make([]*proto.BrandInfoResponse, 0, len(brands))
	for _, brand := range brands {
		brandList = append(brandList, convertBrandToProto(brand))
	}
	
	return &proto.BrandListResponse{
		Total: total,
		Data:  brandList,
	}, nil
}

// CreateBrand 创建品牌
func (h *ProductHandler) CreateBrand(ctx context.Context, req *proto.BrandInfoRequest) (*proto.BrandInfoResponse, error) {
	// 构建品牌实体
	brand := &entity.Brand{
		Name: req.Name,
		Logo: req.Logo,
	}
	
	// 创建品牌
	createdBrand, err := h.brandService.CreateBrand(ctx, brand)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "创建品牌失败: %v", err)
	}
	
	// 转换为响应格式
	return convertBrandToProto(createdBrand), nil
}

// UpdateBrand 更新品牌
func (h *ProductHandler) UpdateBrand(ctx context.Context, req *proto.BrandInfoRequest) (*emptypb.Empty, error) {
	// 检查品牌是否存在
	existingBrand, err := h.brandService.GetBrandByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询品牌失败: %v", err)
	}
	
	if existingBrand == nil {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	
	// 更新品牌信息
	existingBrand.Name = req.Name
	existingBrand.Logo = req.Logo
	
	// 更新品牌
	if err := h.brandService.UpdateBrand(ctx, existingBrand); err != nil {
		return nil, status.Errorf(codes.Internal, "更新品牌失败: %v", err)
	}
	
	return &emptypb.Empty{}, nil
}

// DeleteBrand 删除品牌
func (h *ProductHandler) DeleteBrand(ctx context.Context, req *proto.BrandInfoRequest) (*emptypb.Empty, error) {
	// 删除品牌
	if err := h.brandService.DeleteBrand(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "删除品牌失败: %v", err)
	}
	
	return &emptypb.Empty{}, nil
}

// BannerList 获取轮播图列表
func (h *ProductHandler) BannerList(ctx context.Context, _ *emptypb.Empty) (*proto.BannerListResponse, error) {
	// 获取轮播图列表
	banners, err := h.bannerService.ListBanners(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取轮播图列表失败: %v", err)
	}
	
	// 转换为响应格式
	bannerList := make([]*proto.BannerResponse, 0, len(banners))
	for _, banner := range banners {
		bannerList = append(bannerList, &proto.BannerResponse{
			Id:    banner.ID,
			Index: int32(banner.Index),
			Image: banner.Image,
			Url:   banner.Url,
		})
	}
	
	return &proto.BannerListResponse{
		Total: int64(len(bannerList)),
		Data:  bannerList,
	}, nil
}

// CreateBanner 创建轮播图
func (h *ProductHandler) CreateBanner(ctx context.Context, req *proto.BannerRequest) (*proto.BannerResponse, error) {
	// 构建轮播图实体
	banner := &entity.Banner{
		Index: int(req.Index),
		Image: req.Image,
		Url:   req.Url,
	}
	
	// 创建轮播图
	createdBanner, err := h.bannerService.CreateBanner(ctx, banner)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "创建轮播图失败: %v", err)
	}
	
	// 转换为响应格式
	return &proto.BannerResponse{
		Id:    createdBanner.ID,
		Index: int32(createdBanner.Index),
		Image: createdBanner.Image,
		Url:   createdBanner.Url,
	}, nil
}

// UpdateBanner 更新轮播图
func (h *ProductHandler) UpdateBanner(ctx context.Context, req *proto.BannerRequest) (*emptypb.Empty, error) {
	// 检查轮播图是否存在
	existingBanner, err := h.bannerService.GetBannerByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "查询轮播图失败: %v", err)
	}
	
	if existingBanner == nil {
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}
	
	// 更新轮播图信息
	existingBanner.Index = int(req.Index)
	existingBanner.Image = req.Image
	existingBanner.Url = req.Url
	
	// 更新轮播图
	if err := h.bannerService.UpdateBanner(ctx, existingBanner); err != nil {
		return nil, status.Errorf(codes.Internal, "更新轮播图失败: %v", err)
	}
	
	return &emptypb.Empty{}, nil
}

// DeleteBanner 删除轮播图
func (h *ProductHandler) DeleteBanner(ctx context.Context, req *proto.BannerRequest) (*emptypb.Empty, error) {
	// 删除轮播图
	if err := h.bannerService.DeleteBanner(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "删除轮播图失败: %v", err)
	}
	
	return &emptypb.Empty{}, nil
}

// SearchGoods 商品搜索
func (h *ProductHandler) SearchGoods(ctx context.Context, req *proto.GoodsSearchRequest) (*proto.GoodsSearchResponse, error) {
	// 构建搜索参数
	searchParams := &service.SearchParams{
		Keyword:   req.Keywords,
		CategoryID: req.CategoryId,
		BrandID:   req.BrandId,
		PriceMin:  req.PriceMin,
		PriceMax:  req.PriceMax,
		IsNew:     req.IsNew,
		IsHot:     req.IsHot,
		Page:      int(req.Page),
		PageSize:  int(req.PageSize),
		OrderBy:   req.OrderBy,
	}
	
	// 执行搜索
	result, err := h.searchService.SearchProducts(ctx, searchParams)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "商品搜索失败: %v", err)
	}
	
	// 转换为响应格式
	goodsList := make([]*proto.GoodsInfoResponse, 0, len(result.Goods))
	for _, product := range result.Goods {
		goodsList = append(goodsList, convertProductToProto(product))
	}
	
	return &proto.GoodsSearchResponse{
		Total: result.Total,
		Goods: goodsList,
	}, nil
}

// 工具函数：转换商品实体为proto响应
func convertProductToProto(product *entity.Product) *proto.GoodsInfoResponse {
	if product == nil {
		return nil
	}
	
	goodsInfo := &proto.GoodsInfoResponse{
		Id:              product.ID,
		Name:            product.Name,
		GoodsSn:         product.GoodsSN,
		MarketPrice:     product.MarketPrice,
		ShopPrice:       product.ShopPrice,
		GoodsBrief:      product.GoodsBrief,
		GoodsDesc:       product.GoodsDesc,
		ShipFree:        product.ShipFree,
		GoodsFrontImage: product.GoodsFrontImage,
		IsNew:           product.IsNew,
		IsHot:           product.IsHot,
		OnSale:          product.OnSale,
		CategoryId:      product.CategoryID,
		BrandId:         product.BrandsID,
		CreatedAt:       timestamppb.New(product.CreatedAt),
	}
	
	// 添加分类信息
	if product.Category != nil {
		goodsInfo.Category = convertCategoryToProto(product.Category)
	}
	
	// 添加品牌信息
	if product.Brand != nil {
		goodsInfo.Brand = convertBrandToProto(product.Brand)
	}
	
	// 商品库存
	if len(product.SKUs) > 0 {
		var totalStock int32
		for _, sku := range product.SKUs {
			totalStock += int32(sku.Stock)
		}
		goodsInfo.Stocks = totalStock
	}
	
	// 商品图片
	if len(product.Images) > 0 {
		images := make([]string, 0, len(product.Images))
		for _, img := range product.Images {
			images = append(images, img.Url)
		}
		goodsInfo.Images = images
	}
	
	return goodsInfo
}

// 工具函数：转换分类实体为proto响应
func convertCategoryToProto(category *entity.Category) *proto.CategoryInfoResponse {
	if category == nil {
		return nil
	}
	
	return &proto.CategoryInfoResponse{
		Id:               category.ID,
		Name:             category.Name,
		ParentCategoryId: category.ParentCategoryID,
		Level:            int32(category.Level),
		IsTab:            category.IsTab,
	}
}

// 工具函数：转换品牌实体为proto响应
func convertBrandToProto(brand *entity.Brand) *proto.BrandInfoResponse {
	if brand == nil {
		return nil
	}
	
	return &proto.BrandInfoResponse{
		Id:   brand.ID,
		Name: brand.Name,
		Logo: brand.Logo,
	}
}

// 工具函数：创建float64指针
func floatPtr(v float64) *float64 {
	return &v
}

// 工具函数：创建bool指针
func boolPtr(v bool) *bool {
	return &v
}
