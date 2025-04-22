package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"userop_srv/global"
	"userop_srv/model"
	proto "userop_srv/proto"
)

func (s *UserOpServer) GetAddressList(ctx context.Context, req *proto.AddressRequest) (*proto.AddressListResponse, error) {
	var addressInfo []*model.Address
	var addressRsp proto.AddressListResponse
	var addressResponse []*proto.AddressResponse
	if result := global.DB.Where(&model.Address{User: req.UserId}).Find(&addressInfo); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "不存在收货地址")
	} else {
		addressRsp.Total = int32(result.RowsAffected)
	}
	for _, address := range addressInfo {
		addressResponse = append(addressResponse, &proto.AddressResponse{
			UserId:       address.User,
			Province:     address.Province,
			City:         address.City,
			District:     address.District,
			Address:      address.Address,
			SignerName:   address.SignerName,
			SignerMobile: address.SignerMobile,
		})
	}
	addressRsp.Data = addressResponse
	return &addressRsp, nil
}
func (s *UserOpServer) CreateAddress(ctx context.Context, req *proto.AddressRequest) (*proto.AddressResponse, error) {
	addressInfo := model.Address{
		User:         req.UserId,
		Province:     req.Province,
		City:         req.City,
		District:     req.District,
		Address:      req.Address,
		SignerName:   req.SignerName,
		SignerMobile: req.SignerMobile,
	}
	if result := global.DB.Create(&addressInfo); result.Error != nil {
		return nil, result.Error
	}
	return &proto.AddressResponse{Id: addressInfo.ID}, nil
}
func (s *UserOpServer) DeleteAddress(ctx context.Context, req *proto.AddressRequest) (*proto.Empty, error) {
	if result := global.DB.Where(&model.Address{BaseModel: model.BaseModel{ID: req.Id}}).First(&model.Address{}); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "收货地址不存在")
	}
	if result := global.DB.Where(&model.Address{BaseModel: model.BaseModel{ID: req.Id}}).Delete(&model.Address{}); result.Error != nil {
		return nil, result.Error
	}
	return &proto.Empty{}, nil
}
func (s *UserOpServer) UpdateAddress(ctx context.Context, req *proto.AddressRequest) (*proto.Empty, error) {
	var address model.Address
	if result := global.DB.Where(&model.Address{BaseModel: model.BaseModel{ID: req.Id}, User: req.UserId}).First(&address); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "收货地址不存在")
	}
	address.Province = req.Province
	address.City = req.City
	address.District = req.District
	address.Address = req.Address
	address.SignerName = req.SignerName
	address.SignerMobile = req.SignerMobile
	if result := global.DB.Save(&address); result.Error != nil {
		return nil, result.Error
	}
	return &proto.Empty{}, nil
}
