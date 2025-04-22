package handler

import (
	"context"

	"goods_srv/global"
	"goods_srv/model"
	proto "goods_srv/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// //品牌和轮播图
func (s *GoodsServer) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	brandListResponse := proto.BrandListResponse{}

	var brands []model.Brands
	//result := global.DB.Find(&brands)

	result := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brands)
	if result.Error != nil {
		return nil, result.Error
	}
	var total int64
	global.DB.Model(&model.Brands{}).Count(&total)
	brandListResponse.Total = int32(total)

	var brandResponses []*proto.BrandInfoResponse
	for _, brand := range brands {
		brandResponses = append(brandResponses, &proto.BrandInfoResponse{
			Id:   brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}
	brandListResponse.Data = brandResponses
	return &brandListResponse, nil
}

func (s *GoodsServer) CreateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	brand := &model.Brands{Name: req.Name}
	if result := global.DB.Where(&brand).First(brand); result.RowsAffected == 1 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌已存在")
	}
	brand.Logo = req.Logo
	if result := global.DB.Save(&brand); result.Error != nil {
		return nil, status.Errorf(codes.Internal, "创建品牌失败")
	}
	return &proto.BrandInfoResponse{Id: brand.ID, Logo: brand.Logo, Name: brand.Name}, nil

}

func (s *GoodsServer) DeleteBrand(ctx context.Context, req *proto.BrandRequest) (*proto.Empty, error) {
	if result := global.DB.Delete(&model.Brands{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "品牌不存在")
	}
	return &proto.Empty{}, nil
}

func (s *GoodsServer) UpdateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.Empty, error) {
	brands := model.Brands{}
	if result := global.DB.Where(&model.Brands{BaseModel: model.BaseModel{ID: req.Id}}).First(&brands); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}
	if req.Name != "" {
		brands.Name = req.Name
	}
	if req.Logo != "" {
		brands.Logo = req.Logo
	}

	if result := global.DB.Save(&brands); result.Error != nil {
		return nil, status.Errorf(codes.Internal, "创建品牌失败")
	}

	return &proto.Empty{}, nil
}
