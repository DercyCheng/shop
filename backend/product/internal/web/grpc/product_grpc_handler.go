package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "shop/product/api/proto"
	"shop/product/internal/domain/entity"
	"shop/product/internal/service"
)

// ProductGRPCServer implements the gRPC ProductService
type ProductGRPCServer struct {
	productService       service.ProductService
	categoryService      service.CategoryService
	brandService         service.BrandService
	bannerService        service.BannerService
	categoryBrandService service.CategoryBrandService
	pb.UnimplementedProductServiceServer
}

// NewProductGRPCServer creates a new ProductGRPCServer
func NewProductGRPCServer(
	productService service.ProductService,
	categoryService service.CategoryService,
	brandService service.BrandService,
	bannerService service.BannerService,
	categoryBrandService service.CategoryBrandService,
) *ProductGRPCServer {
	return &ProductGRPCServer{
		productService:       productService,
		categoryService:      categoryService,
		brandService:         brandService,
		bannerService:        bannerService,
		categoryBrandService: categoryBrandService,
	}
}

// GetProduct retrieves a product by ID
func (s *ProductGRPCServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductInfo, error) {
	product, err := s.productService.GetProductWithRelations(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "product not found: %v", err)
	}

	return convertProductToProto(product), nil
}

// ListProducts retrieves products based on filters
func (s *ProductGRPCServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ProductListResponse, error) {
	filter := &service.ProductFilter{
		Page:       int(req.Page),
		PageSize:   int(req.PageSize),
		CategoryID: req.CategoryId,
		BrandID:    req.BrandId,
		Keyword:    req.Keyword,
		MinPrice:   float64(req.MinPrice),
		MaxPrice:   float64(req.MaxPrice),
	}

	// Set boolean filters
	onSale := req.OnSale
	if onSale {
		filter.OnSale = &onSale
	}

	shipFree := req.ShipFree
	if shipFree {
		filter.ShipFree = &shipFree
	}

	isNew := req.IsNew
	if isNew {
		filter.IsNew = &isNew
	}

	isHot := req.IsHot
	if isHot {
		filter.IsHot = &isHot
	}

	// Set ordering based on OrderBy enum
	switch req.OrderBy {
	case pb.ListProductsRequest_PRICE_ASC:
		filter.OrderBy = "shop_price"
		filter.Order = "asc"
	case pb.ListProductsRequest_PRICE_DESC:
		filter.OrderBy = "shop_price"
		filter.Order = "desc"
	case pb.ListProductsRequest_SOLD_DESC:
		filter.OrderBy = "sold_num"
		filter.Order = "desc"
	case pb.ListProductsRequest_NEW_DESC:
		filter.OrderBy = "created_at"
		filter.Order = "desc"
	}

	products, total, err := s.productService.ListProducts(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list products: %v", err)
	}

	productInfos := make([]*pb.ProductInfo, 0, len(products))
	for _, product := range products {
		productInfos = append(productInfos, convertProductToProto(product))
	}

	return &pb.ProductListResponse{
		Total:    int32(total),
		Products: productInfos,
	}, nil
}

// BatchGetProducts retrieves multiple products by IDs
func (s *ProductGRPCServer) BatchGetProducts(ctx context.Context, req *pb.BatchGetProductsRequest) (*pb.BatchGetProductsResponse, error) {
	products, err := s.productService.GetProductsByIDs(ctx, req.ProductIds)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get products: %v", err)
	}

	productInfos := make([]*pb.ProductInfo, 0, len(products))
	for _, product := range products {
		productInfos = append(productInfos, convertProductToProto(product))
	}

	return &pb.BatchGetProductsResponse{
		Products: productInfos,
	}, nil
}

// CreateProduct adds a new product
func (s *ProductGRPCServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductInfo, error) {
	product := &entity.Product{
		Name:            req.Name,
		GoodsSN:         req.GoodsSn,
		CategoryID:      req.CategoryId,
		BrandID:         req.BrandId,
		OnSale:          req.OnSale,
		ShipFree:        req.ShipFree,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		MarketPrice:     float64(req.MarketPrice),
		ShopPrice:       float64(req.ShopPrice),
		GoodsBrief:      req.GoodsBrief,
		GoodsDesc:       req.GoodsDesc,
		GoodsFrontImage: req.GoodsFrontImage,
	}

	createdProduct, err := s.productService.CreateProduct(ctx, product, req.Images)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}

	return convertProductToProto(createdProduct), nil
}

// UpdateProduct modifies an existing product
func (s *ProductGRPCServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductInfo, error) {
	product := &entity.Product{
		ID:              req.Id,
		Name:            req.Name,
		GoodsSN:         req.GoodsSn,
		CategoryID:      req.CategoryId,
		BrandID:         req.BrandId,
		OnSale:          req.OnSale,
		ShipFree:        req.ShipFree,
		IsNew:           req.IsNew,
		IsHot:           req.IsHot,
		MarketPrice:     float64(req.MarketPrice),
		ShopPrice:       float64(req.ShopPrice),
		GoodsBrief:      req.GoodsBrief,
		GoodsDesc:       req.GoodsDesc,
		GoodsFrontImage: req.GoodsFrontImage,
	}

	if err := s.productService.UpdateProduct(ctx, product, req.Images); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update product: %v", err)
	}

	// Fetch the updated product
	updatedProduct, err := s.productService.GetProductWithRelations(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "product updated but failed to retrieve: %v", err)
	}

	return convertProductToProto(updatedProduct), nil
}

// DeleteProduct removes a product by ID
func (s *ProductGRPCServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*emptypb.Empty, error) {
	if err := s.productService.DeleteProduct(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete product: %v", err)
	}
	return &emptypb.Empty{}, nil
}

// GetCategory retrieves a category by ID
func (s *ProductGRPCServer) GetCategory(ctx context.Context, req *pb.GetCategoryRequest) (*pb.CategoryInfo, error) {
	category, err := s.categoryService.GetCategoryByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "category not found: %v", err)
	}

	return convertCategoryToProto(category), nil
}

// ListCategories retrieves categories based on filters
func (s *ProductGRPCServer) ListCategories(ctx context.Context, req *pb.ListCategoriesRequest) (*pb.CategoryListResponse, error) {
	var categories []*entity.Category
	var err error

	if req.ParentId > 0 {
		categories, err = s.categoryService.GetCategoriesByParentID(ctx, req.ParentId)
	} else if req.Level > 0 {
		categories, err = s.categoryService.GetCategoriesByLevel(ctx, int(req.Level))
	} else {
		categories, err = s.categoryService.GetAllCategories(ctx)
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list categories: %v", err)
	}

	categoryInfos := make([]*pb.CategoryInfo, 0, len(categories))
	for _, category := range categories {
		categoryInfos = append(categoryInfos, convertCategoryToProto(category))
	}

	return &pb.CategoryListResponse{
		Categories: categoryInfos,
	}, nil
}

// CreateCategory adds a new category
func (s *ProductGRPCServer) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CategoryInfo, error) {
	category := &entity.Category{
		Name:             req.Name,
		ParentCategoryID: req.ParentCategoryId,
		Level:            int(req.Level),
		IsTab:            req.IsTab,
	}

	createdCategory, err := s.categoryService.CreateCategory(ctx, category)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create category: %v", err)
	}

	return convertCategoryToProto(createdCategory), nil
}

// UpdateCategory modifies an existing category
func (s *ProductGRPCServer) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.CategoryInfo, error) {
	category := &entity.Category{
		ID:               req.Id,
		Name:             req.Name,
		ParentCategoryID: req.ParentCategoryId,
		Level:            int(req.Level),
		IsTab:            req.IsTab,
	}

	if err := s.categoryService.UpdateCategory(ctx, category); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update category: %v", err)
	}

	// Fetch the updated category
	updatedCategory, err := s.categoryService.GetCategoryByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "category updated but failed to retrieve: %v", err)
	}

	return convertCategoryToProto(updatedCategory), nil
}

// DeleteCategory removes a category by ID
func (s *ProductGRPCServer) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*emptypb.Empty, error) {
	if err := s.categoryService.DeleteCategory(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete category: %v", err)
	}
	return &emptypb.Empty{}, nil
}

// GetBrand retrieves a brand by ID
func (s *ProductGRPCServer) GetBrand(ctx context.Context, req *pb.GetBrandRequest) (*pb.BrandInfo, error) {
	brand, err := s.brandService.GetBrandByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "brand not found: %v", err)
	}

	return convertBrandToProto(brand), nil
}

// ListBrands retrieves brands with pagination
func (s *ProductGRPCServer) ListBrands(ctx context.Context, req *pb.ListBrandsRequest) (*pb.BrandListResponse, error) {
	brands, total, err := s.brandService.ListBrands(ctx, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list brands: %v", err)
	}

	brandInfos := make([]*pb.BrandInfo, 0, len(brands))
	for _, brand := range brands {
		brandInfos = append(brandInfos, convertBrandToProto(brand))
	}

	return &pb.BrandListResponse{
		Total:  int32(total),
		Brands: brandInfos,
	}, nil
}

// CreateBrand adds a new brand
func (s *ProductGRPCServer) CreateBrand(ctx context.Context, req *pb.CreateBrandRequest) (*pb.BrandInfo, error) {
	brand := &entity.Brand{
		Name: req.Name,
		Logo: req.Logo,
	}

	createdBrand, err := s.brandService.CreateBrand(ctx, brand)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create brand: %v", err)
	}

	return convertBrandToProto(createdBrand), nil
}

// UpdateBrand modifies an existing brand
func (s *ProductGRPCServer) UpdateBrand(ctx context.Context, req *pb.UpdateBrandRequest) (*pb.BrandInfo, error) {
	brand := &entity.Brand{
		ID:   req.Id,
		Name: req.Name,
		Logo: req.Logo,
	}

	if err := s.brandService.UpdateBrand(ctx, brand); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update brand: %v", err)
	}

	// Fetch the updated brand
	updatedBrand, err := s.brandService.GetBrandByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "brand updated but failed to retrieve: %v", err)
	}

	return convertBrandToProto(updatedBrand), nil
}

// DeleteBrand removes a brand by ID
func (s *ProductGRPCServer) DeleteBrand(ctx context.Context, req *pb.DeleteBrandRequest) (*emptypb.Empty, error) {
	if err := s.brandService.DeleteBrand(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete brand: %v", err)
	}
	return &emptypb.Empty{}, nil
}

// ListBanners retrieves all banners ordered by index
func (s *ProductGRPCServer) ListBanners(ctx context.Context, req *pb.ListBannersRequest) (*pb.BannerListResponse, error) {
	banners, err := s.bannerService.ListBanners(ctx, int(req.Limit))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list banners: %v", err)
	}

	bannerInfos := make([]*pb.BannerInfo, 0, len(banners))
	for _, banner := range banners {
		bannerInfos = append(bannerInfos, convertBannerToProto(banner))
	}

	return &pb.BannerListResponse{
		Banners: bannerInfos,
	}, nil
}

// CreateBanner adds a new banner
func (s *ProductGRPCServer) CreateBanner(ctx context.Context, req *pb.CreateBannerRequest) (*pb.BannerInfo, error) {
	banner := &entity.Banner{
		Image: req.Image,
		URL:   req.Url,
		Index: int(req.Index),
	}

	createdBanner, err := s.bannerService.CreateBanner(ctx, banner)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create banner: %v", err)
	}

	return convertBannerToProto(createdBanner), nil
}

// UpdateBanner modifies an existing banner
func (s *ProductGRPCServer) UpdateBanner(ctx context.Context, req *pb.UpdateBannerRequest) (*pb.BannerInfo, error) {
	banner := &entity.Banner{
		ID:    req.Id,
		Image: req.Image,
		URL:   req.Url,
		Index: int(req.Index),
	}

	if err := s.bannerService.UpdateBanner(ctx, banner); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update banner: %v", err)
	}

	// Fetch the updated banner
	updatedBanner, err := s.bannerService.GetBannerByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "banner updated but failed to retrieve: %v", err)
	}

	return convertBannerToProto(updatedBanner), nil
}

// DeleteBanner removes a banner by ID
func (s *ProductGRPCServer) DeleteBanner(ctx context.Context, req *pb.DeleteBannerRequest) (*emptypb.Empty, error) {
	if err := s.bannerService.DeleteBanner(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete banner: %v", err)
	}
	return &emptypb.Empty{}, nil
}

// ListCategoryBrands retrieves category-brand relations with pagination
func (s *ProductGRPCServer) ListCategoryBrands(ctx context.Context, req *pb.ListCategoryBrandsRequest) (*pb.CategoryBrandListResponse, error) {
	categoryBrands, total, err := s.categoryBrandService.ListCategoryBrands(ctx, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list category-brand relations: %v", err)
	}

	categoryBrandInfos := make([]*pb.CategoryBrandInfo, 0, len(categoryBrands))
	for _, cb := range categoryBrands {
		categoryBrandInfos = append(categoryBrandInfos, convertCategoryBrandToProto(cb))
	}

	return &pb.CategoryBrandListResponse{
		Total:          int32(total),
		CategoryBrands: categoryBrandInfos,
	}, nil
}

// GetCategoryBrandList retrieves all brands for a category
func (s *ProductGRPCServer) GetCategoryBrandList(ctx context.Context, req *pb.GetCategoryBrandListRequest) (*pb.BrandListResponse, error) {
	brands, err := s.categoryBrandService.GetBrandsByCategoryID(ctx, req.CategoryId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get brands for category: %v", err)
	}

	brandInfos := make([]*pb.BrandInfo, 0, len(brands))
	for _, brand := range brands {
		brandInfos = append(brandInfos, convertBrandToProto(brand))
	}

	return &pb.BrandListResponse{
		Total:  int32(len(brandInfos)),
		Brands: brandInfos,
	}, nil
}

// CreateCategoryBrand adds a new category-brand relation
func (s *ProductGRPCServer) CreateCategoryBrand(ctx context.Context, req *pb.CreateCategoryBrandRequest) (*pb.CategoryBrandInfo, error) {
	categoryBrand := &entity.CategoryBrand{
		CategoryID: req.CategoryId,
		BrandID:    req.BrandId,
	}

	createdCategoryBrand, err := s.categoryBrandService.CreateCategoryBrand(ctx, categoryBrand)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create category-brand relation: %v", err)
	}

	return convertCategoryBrandToProto(createdCategoryBrand), nil
}

// UpdateCategoryBrand modifies an existing category-brand relation
func (s *ProductGRPCServer) UpdateCategoryBrand(ctx context.Context, req *pb.UpdateCategoryBrandRequest) (*pb.CategoryBrandInfo, error) {
	categoryBrand := &entity.CategoryBrand{
		ID:         req.Id,
		CategoryID: req.CategoryId,
		BrandID:    req.BrandId,
	}

	if err := s.categoryBrandService.UpdateCategoryBrand(ctx, categoryBrand); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update category-brand relation: %v", err)
	}

	// Fetch the updated category-brand relation
	updatedCategoryBrand, err := s.categoryBrandService.GetCategoryBrandByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "category-brand relation updated but failed to retrieve: %v", err)
	}

	return convertCategoryBrandToProto(updatedCategoryBrand), nil
}

// DeleteCategoryBrand removes a category-brand relation by ID
func (s *ProductGRPCServer) DeleteCategoryBrand(ctx context.Context, req *pb.DeleteCategoryBrandRequest) (*emptypb.Empty, error) {
	if err := s.categoryBrandService.DeleteCategoryBrand(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete category-brand relation: %v", err)
	}
	return &emptypb.Empty{}, nil
}

// Helper functions to convert entity objects to protobuf messages
func convertProductToProto(product *entity.Product) *pb.ProductInfo {
	if product == nil {
		return nil
	}

	images := make([]string, 0)
	if product.Images != nil {
		images = product.Images
	}

	return &pb.ProductInfo{
		Id:              product.ID,
		Name:            product.Name,
		GoodsSn:         product.GoodsSN,
		CategoryId:      product.CategoryID,
		BrandId:         product.BrandID,
		OnSale:          product.OnSale,
		ShipFree:        product.ShipFree,
		IsNew:           product.IsNew,
		IsHot:           product.IsHot,
		ClickNum:        int32(product.ClickNum),
		SoldNum:         int32(product.SoldNum),
		FavNum:          int32(product.FavNum),
		MarketPrice:     float32(product.MarketPrice),
		ShopPrice:       float32(product.ShopPrice),
		GoodsBrief:      product.GoodsBrief,
		GoodsDesc:       product.GoodsDesc,
		GoodsFrontImage: product.GoodsFrontImage,
		Images:          images,
		CreatedAt:       timestamppb.New(product.CreatedAt),
		UpdatedAt:       timestamppb.New(product.UpdatedAt),
	}
}

func convertCategoryToProto(category *entity.Category) *pb.CategoryInfo {
	if category == nil {
		return nil
	}

	return &pb.CategoryInfo{
		Id:               category.ID,
		Name:             category.Name,
		ParentCategoryId: category.ParentCategoryID,
		Level:            int32(category.Level),
		IsTab:            category.IsTab,
		CreatedAt:        timestamppb.New(category.CreatedAt),
		UpdatedAt:        timestamppb.New(category.UpdatedAt),
	}
}

func convertBrandToProto(brand *entity.Brand) *pb.BrandInfo {
	if brand == nil {
		return nil
	}

	return &pb.BrandInfo{
		Id:        brand.ID,
		Name:      brand.Name,
		Logo:      brand.Logo,
		CreatedAt: timestamppb.New(brand.CreatedAt),
		UpdatedAt: timestamppb.New(brand.UpdatedAt),
	}
}

func convertBannerToProto(banner *entity.Banner) *pb.BannerInfo {
	if banner == nil {
		return nil
	}

	return &pb.BannerInfo{
		Id:        banner.ID,
		Image:     banner.Image,
		Url:       banner.URL,
		Index:     int32(banner.Index),
		CreatedAt: timestamppb.New(banner.CreatedAt),
		UpdatedAt: timestamppb.New(banner.UpdatedAt),
	}
}

func convertCategoryBrandToProto(categoryBrand *entity.CategoryBrand) *pb.CategoryBrandInfo {
	if categoryBrand == nil {
		return nil
	}

	result := &pb.CategoryBrandInfo{
		Id:         categoryBrand.ID,
		CategoryId: categoryBrand.CategoryID,
		BrandId:    categoryBrand.BrandID,
		CreatedAt:  timestamppb.New(categoryBrand.CreatedAt),
		UpdatedAt:  timestamppb.New(categoryBrand.UpdatedAt),
	}

	if categoryBrand.Category != nil {
		result.Category = convertCategoryToProto(categoryBrand.Category)
	}

	if categoryBrand.Brand != nil {
		result.Brand = convertBrandToProto(categoryBrand.Brand)
	}

	return result
}
