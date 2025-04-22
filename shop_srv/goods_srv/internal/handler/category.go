package handler

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"goods_srv/global"
	"goods_srv/model"
	proto "goods_srv/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

// 商品分类
func (s *GoodsServer) GetAllCategorysList(context.Context, *proto.Empty) (*proto.CategoryListResponse, error) {
	/*
		[
			{
				"id":xxx,
				"name":"",
				"level":1,
				"is_tab":false,
				"parent":"",
				"sub_category":[{
					"id":xxx,
					"name":"",
					"level":2,
					"parent":13xxx,
					"is_tab":false,
					"sub_category":[{
						"id":xxx,
						"name":"",
						"level":3,
						"is_tab":false,
						"parent":13xxx,
						"sub_category":[]
					}]
				}]
			}
		]
	*/
	categoryRsp := proto.CategoryListResponse{}

	var categorys []model.Category

	if result := global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys); result.Error != nil {
		return nil, result.Error
	}
	b, _ := json.Marshal(&categorys)
	categoryRsp.JsonData = string(b)

	var categorys_proto []model.Category
	result := global.DB.Find(&categorys_proto)
	if result.Error != nil {
		return nil, result.Error
	}
	categoryRsp.Total = int32(result.RowsAffected)
	for _, category := range categorys_proto {
		categoryInfo := proto.CategoryInfoResponse{}
		categoryInfo.Id = category.ID
		categoryInfo.Name = category.Name
		categoryInfo.ParentCategory = category.ParentCategoryID
		categoryInfo.Level = category.Level
		categoryInfo.IsTab = category.IsTab
		categoryRsp.Data = append(categoryRsp.Data, &categoryInfo)
	}
	return &categoryRsp, nil
}

// 获取子分类
func (s *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	categoryListResponse := proto.SubCategoryListResponse{}
	var category model.Category
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	categoryListResponse.Info = &proto.CategoryInfoResponse{
		Id:             category.ID,
		Name:           category.Name,
		Level:          category.Level,
		IsTab:          category.IsTab,
		ParentCategory: category.ParentCategoryID,
	}
	var subCategoorys []model.Category
	var subCategoryResponse []*proto.CategoryInfoResponse
	//preloads := "SubCategory"
	//if category.Level == 1 {
	//	preloads = "SubCategory.SubCategory"
	//}

	//if result := global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Preload(preloads).Find(&subCategoorys); result.Error != nil {
	//	return nil, result.Error
	//}
	if result := global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Find(&subCategoorys); result.Error != nil {
		return nil, result.Error
	}
	for _, subCategoory := range subCategoorys {
		subCategoryResponse = append(subCategoryResponse, &proto.CategoryInfoResponse{
			Id:             subCategoory.ID,
			Name:           subCategoory.Name,
			Level:          subCategoory.Level,
			IsTab:          subCategoory.IsTab,
			ParentCategory: subCategoory.ParentCategoryID,
		})
	}
	categoryListResponse.SubCategorys = subCategoryResponse
	return &categoryListResponse, nil
}

// 新建商品分类
func (s *GoodsServer) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	category := model.Category{}
	cMap := map[string]interface{}{}
	cMap["add_time"] = time.Now()
	cMap["name"] = req.Name
	cMap["level"] = req.Level
	cMap["is_tab"] = req.IsTab
	if req.Level != 1 {
		//去查询父类目是否存在
		cMap["parent_category_id"] = req.ParentCategory
	}
	if result := global.DB.Model(&category).Create(cMap); result.Error != nil {
		zap.S().Error("新建商品分类失败！")
	}
	//fmt.Println(tx)
	rsp := proto.CategoryInfoResponse{
		Id:             category.ID,
		Name:           category.Name,
		Level:          category.Level,
		IsTab:          category.IsTab,
		ParentCategory: category.ParentCategoryID,
	}
	return &rsp, nil
}

// 删除商品分类
func (s *GoodsServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*proto.Empty, error) {
	if result := global.DB.Delete(&model.Category{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	return &proto.Empty{}, nil
}

// 更新商品分类
func (s *GoodsServer) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.Empty, error) {

	var category model.Category
	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	//sqlMap:=make(map[string]int)
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.ParentCategory != 0 {
		category.ParentCategoryID = req.ParentCategory
	}
	if req.Level != 0 {
		category.Level = req.Level
	}
	if req.IsTab {
		category.IsTab = req.IsTab
	}
	if req.ParentCategory != 0 {
		if result := global.DB.Omit("ParentCategoryID").Save(&category); result.Error != nil {
			zap.S().Error("更新商品分类失败", result.Error)
		}
	} else {
		if result := global.DB.Save(&category); result.Error != nil {
			zap.S().Error("更新商品分类失败", result.Error)
		}
	}

	return &proto.Empty{}, nil
}
