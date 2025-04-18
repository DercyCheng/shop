package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"nd/userop_srv/global"
	"nd/userop_srv/model"
	"nd/userop_srv/proto"
)

func (*UserOpServer) GetAddressList(ctx context.Context, req *proto.AddressRequest) (*proto.AddressListResponse, error) {
	// 参数验证
	if req.UserId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "用户ID不能为空")
	}
	
	var addresses []model.Address
	var rsp proto.AddressListResponse
	
	// 查询地址列表
	result := global.DB.Where(&model.Address{User: req.UserId}).Find(&addresses)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "查询地址列表失败: %s", result.Error.Error())
	}
	
	rsp.Total = int32(result.RowsAffected)
	
	// 如果没有找到记录，返回空列表而不是错误
	if result.RowsAffected == 0 {
		rsp.Data = []*proto.AddressResponse{}
		return &rsp, nil
	}
	
	// 构建响应数据
	addressResponse := make([]*proto.AddressResponse, 0, len(addresses))
	for _, address := range addresses {
		addressResponse = append(addressResponse, &proto.AddressResponse{
			Id:           address.ID,
			UserId:       address.User,
			Province:     address.Province,
			City:         address.City,
			District:     address.District,
			Address:      address.Address,
			SignerName:   address.SignerName,
			SignerMobile: address.SignerMobile,
		})
	}
	rsp.Data = addressResponse
	
	return &rsp, nil
}

func (*UserOpServer) CreateAddress(ctx context.Context, req *proto.AddressRequest) (*proto.AddressResponse, error) {
	// 参数验证
	if req.UserId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "用户ID不能为空")
	}
	if req.Province == "" || req.City == "" || req.District == "" || req.Address == "" {
		return nil, status.Errorf(codes.InvalidArgument, "地址信息不完整")
	}
	if req.SignerName == "" || req.SignerMobile == "" {
		return nil, status.Errorf(codes.InvalidArgument, "收件人信息不完整")
	}

	// 创建地址
	address := model.Address{
		User:         req.UserId,
		Province:     req.Province,
		City:         req.City,
		District:     req.District,
		Address:      req.Address,
		SignerName:   req.SignerName,
		SignerMobile: req.SignerMobile,
	}

	// 保存并检查错误
	result := global.DB.Create(&address)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "创建地址失败: %s", result.Error.Error())
	}

	return &proto.AddressResponse{
		Id:           address.ID,
		UserId:       address.User,
		Province:     address.Province,
		City:         address.City,
		District:     address.District,
		Address:      address.Address,
		SignerName:   address.SignerName,
		SignerMobile: address.SignerMobile,
	}, nil
}

func (*UserOpServer) DeleteAddress(ctx context.Context, req *proto.AddressRequest) (*emptypb.Empty, error) {
	// 参数验证
	if req.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "地址ID不能为空")
	}
	if req.UserId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "用户ID不能为空")
	}
	
	// 执行删除操作
	result := global.DB.Where("id=? and user=?", req.Id, req.UserId).Delete(&model.Address{})
	
	// 检查是否有错误
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "删除地址失败: %s", result.Error.Error())
	}
	
	// 检查是否找到记录
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "收货地址不存在")
	}
	
	return &emptypb.Empty{}, nil
}

func (*UserOpServer) UpdateAddress(ctx context.Context, req *proto.AddressRequest) (*emptypb.Empty, error) {
	// 参数验证
	if req.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "地址ID不能为空")
	}
	if req.UserId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "用户ID不能为空")
	}
	
	// 查询要更新的地址
	var address model.Address
	result := global.DB.Where("id=? and user=?", req.Id, req.UserId).First(&address)
	
	// 检查是否有错误
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "查询地址失败: %s", result.Error.Error())
	}
	
	// 检查是否找到记录
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "地址不存在")
	}
	
	// 只更新提供了的字段
	if req.Province != "" {
		address.Province = req.Province
	}
	
	if req.City != "" {
		address.City = req.City
	}
	
	if req.District != "" {
		address.District = req.District
	}
	
	if req.Address != "" {
		address.Address = req.Address
	}
	
	if req.SignerName != "" {
		address.SignerName = req.SignerName
	}
	
	if req.SignerMobile != "" {
		address.SignerMobile = req.SignerMobile
	}
	
	// 保存更新并检查错误
	updateResult := global.DB.Save(&address)
	if updateResult.Error != nil {
		return nil, status.Errorf(codes.Internal, "更新地址失败: %s", updateResult.Error.Error())
	}
	
	return &emptypb.Empty{}, nil
}
