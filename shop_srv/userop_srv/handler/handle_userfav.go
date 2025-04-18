package handler

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"nd/userop_srv/global"
	"nd/userop_srv/model"
	"nd/userop_srv/proto"
)

func (*UserOpServer) GetFavList(ctx context.Context, req *proto.UserFavRequest) (*proto.UserFavListResponse, error) {
	var rsp proto.UserFavListResponse
	var userFavs []model.UserFav
	var userFavList []*proto.UserFavResponse

	// 构建查询条件
	query := global.DB.Model(&model.UserFav{})
	
	// 根据请求参数构建不同的查询
	if req.UserId != 0 && req.GoodsId != 0 {
		// 查询特定用户是否收藏了特定商品
		query = query.Where("user = ? AND goods = ?", req.UserId, req.GoodsId)
	} else if req.UserId != 0 {
		// 查询用户的所有收藏
		query = query.Where("user = ?", req.UserId)
	} else if req.GoodsId != 0 {
		// 查询商品被哪些用户收藏
		query = query.Where("goods = ?", req.GoodsId)
	}

	// 执行查询
	result := query.Find(&userFavs)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "查询收藏列表失败: %s", result.Error.Error())
	}
	
	rsp.Total = int32(result.RowsAffected)

	// 构建响应数据
	userFavList = make([]*proto.UserFavResponse, 0, len(userFavs))
	for _, userFav := range userFavs {
		userFavList = append(userFavList, &proto.UserFavResponse{
			UserId:  userFav.User,
			GoodsId: userFav.Goods,
		})
	}

	rsp.Data = userFavList
	return &rsp, nil
}

func (*UserOpServer) AddUserFav(ctx context.Context, req *proto.UserFavRequest) (*emptypb.Empty, error) {
	// 先检查是否已经收藏过
	var count int64
	global.DB.Model(&model.UserFav{}).Where("goods=? and user=?", req.GoodsId, req.UserId).Count(&count)
	if count > 0 {
		return nil, status.Errorf(codes.AlreadyExists, "该商品已收藏")
	}

	// 创建新收藏
	userFav := model.UserFav{
		User:  req.UserId,
		Goods: req.GoodsId,
	}

	// 保存并检查错误
	if result := global.DB.Create(&userFav); result.Error != nil {
		return nil, status.Errorf(codes.Internal, "添加收藏失败: %s", result.Error.Error())
	}

	return &emptypb.Empty{}, nil
}

func (*UserOpServer) DeleteUserFav(ctx context.Context, req *proto.UserFavRequest) (*emptypb.Empty, error) {
	// 首先检查参数是否有效
	if req.UserId == 0 || req.GoodsId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "用户ID和商品ID不能为空")
	}
	
	// 执行删除操作
	result := global.DB.Unscoped().Where("goods=? and user=?", req.GoodsId, req.UserId).Delete(&model.UserFav{})
	
	// 检查是否发生错误
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "删除收藏记录失败: %s", result.Error.Error())
	}
	
	// 检查是否删除了记录
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "收藏记录不存在")
	}
	
	return &emptypb.Empty{}, nil
}

func (*UserOpServer) GetUserFavDetail(ctx context.Context, req *proto.UserFavRequest) (*proto.UserFavResponse, error) {
	var userfav model.UserFav
	if result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).Find(&userfav); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "收藏记录不存在")
	}
	return &proto.UserFavResponse{
		UserId:  userfav.User,
		GoodsId: userfav.Goods,
	}, nil
}
