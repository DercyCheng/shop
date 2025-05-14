package grpc

import (
	"context"
	"time"
	
	"shop/backend/profile/api/proto"
	"shop/backend/profile/internal/domain/entity"
	"shop/backend/profile/internal/service"
	
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ProfileServiceServer gRPC服务实现
type ProfileServiceServer struct {
	favService     service.FavoriteService
	addrService    service.AddressService
	feedbackService service.FeedbackService
	historyService service.BrowsingHistoryService
	proto.UnimplementedProfileServiceServer
}

// NewProfileServiceServer 创建个人信息服务gRPC处理器
func NewProfileServiceServer(
	favService service.FavoriteService,
	addrService service.AddressService,
	feedbackService service.FeedbackService,
	historyService service.BrowsingHistoryService,
) *ProfileServiceServer {
	return &ProfileServiceServer{
		favService:      favService,
		addrService:     addrService,
		feedbackService: feedbackService,
		historyService:  historyService,
	}
}

// RegisterWithServer 注册服务到gRPC服务器
func (s *ProfileServiceServer) RegisterWithServer(server *grpc.Server) {
	proto.RegisterProfileServiceServer(server, s)
}

// ======= 收藏相关接口 =======

// ListFavorites 获取收藏列表
func (s *ProfileServiceServer) ListFavorites(ctx context.Context, req *proto.ListFavoritesRequest) (*proto.ListFavoritesResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	
	favorites, total, err := s.favService.ListFavorites(ctx, req.UserId, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get favorites: %v", err)
	}
	
	// 转换为proto类型
	protoFavs := make([]*proto.Favorite, 0, len(favorites))
	for _, fav := range favorites {
		protoFavs = append(protoFavs, &proto.Favorite{
			Id:            fav.ID,
			UserId:        fav.UserID,
			GoodsId:       fav.GoodsID,
			CategoryId:    fav.CategoryID,
			Remark:        fav.Remark,
			PriceWhenFav:  fav.PriceWhenFav,
			Notification:  fav.Notification,
			CreatedAt:     fav.CreatedAt.Unix(),
		})
	}
	
	return &proto.ListFavoritesResponse{
		Favorites: protoFavs,
		Total:     total,
	}, nil
}

// AddFavorite 添加收藏
func (s *ProfileServiceServer) AddFavorite(ctx context.Context, req *proto.AddFavoriteRequest) (*proto.AddFavoriteResponse, error) {
	err := s.favService.AddFavorite(ctx, req.UserId, req.GoodsId, req.CategoryId, req.Remark)
	if err != nil {
		if err == service.ErrFavoriteExists {
			return &proto.AddFavoriteResponse{
				Success: false,
				Message: "Product already in favorites",
			}, nil
		}
		return nil, status.Errorf(codes.Internal, "Failed to add favorite: %v", err)
	}
	
	return &proto.AddFavoriteResponse{
		Success: true,
		Message: "Added to favorites successfully",
	}, nil
}

// RemoveFavorite 取消收藏
func (s *ProfileServiceServer) RemoveFavorite(ctx context.Context, req *proto.RemoveFavoriteRequest) (*proto.RemoveFavoriteResponse, error) {
	err := s.favService.RemoveFavorite(ctx, req.Id)
	if err != nil {
		if err == service.ErrFavoriteNotFound {
			return nil, status.Error(codes.NotFound, "Favorite not found")
		}
		return nil, status.Errorf(codes.Internal, "Failed to remove favorite: %v", err)
	}
	
	return &proto.RemoveFavoriteResponse{
		Success: true,
		Message: "Removed from favorites successfully",
	}, nil
}

// IsFavorite 判断是否已收藏
func (s *ProfileServiceServer) IsFavorite(ctx context.Context, req *proto.IsFavoriteRequest) (*proto.IsFavoriteResponse, error) {
	isFav, err := s.favService.IsFavorite(ctx, req.UserId, req.GoodsId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to check favorite status: %v", err)
	}
	
	return &proto.IsFavoriteResponse{
		IsFavorite: isFav,
	}, nil
}

// SetPriceNotification 设置价格变动通知
func (s *ProfileServiceServer) SetPriceNotification(ctx context.Context, req *proto.SetPriceNotificationRequest) (*proto.SetPriceNotificationResponse, error) {
	err := s.favService.SetPriceNotification(ctx, req.Id, req.Notify)
	if err != nil {
		if err == service.ErrFavoriteNotFound {
			return nil, status.Error(codes.NotFound, "Favorite not found")
		}
		return nil, status.Errorf(codes.Internal, "Failed to set price notification: %v", err)
	}
	
	return &proto.SetPriceNotificationResponse{
		Success: true,
		Message: "Price notification setting updated successfully",
	}, nil
}

// ======= 地址相关接口 =======

// ListAddresses 获取地址列表
func (s *ProfileServiceServer) ListAddresses(ctx context.Context, req *proto.ListAddressesRequest) (*proto.ListAddressesResponse, error) {
	addresses, err := s.addrService.ListAddresses(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get addresses: %v", err)
	}
	
	// 转换为proto类型
	protoAddrs := make([]*proto.Address, 0, len(addresses))
	for _, addr := range addresses {
		protoAddr := &proto.Address{
			Id:          addr.ID,
			UserId:      addr.UserID,
			Province:    addr.Province,
			City:        addr.City,
			District:    addr.District,
			Address:     addr.Address,
			SignerName:  addr.SignerName,
			SignerMobile: addr.SignerMobile,
			IsDefault:   addr.IsDefault,
			Label:       addr.Label,
			Postcode:    addr.Postcode,
			UsageCount:  int32(addr.UsageCount),
			CreatedAt:   addr.CreatedAt.Unix(),
		}
		
		if addr.LastUsedAt != nil {
			protoAddr.LastUsedAt = addr.LastUsedAt.Unix()
		}
		
		protoAddrs = append(protoAddrs, protoAddr)
	}
	
	return &proto.ListAddressesResponse{
		Addresses: protoAddrs,
	}, nil
}

// GetAddress 获取地址详情
func (s *ProfileServiceServer) GetAddress(ctx context.Context, req *proto.GetAddressRequest) (*proto.GetAddressResponse, error) {
	address, err := s.addrService.GetAddress(ctx, req.Id)
	if err != nil {
		if err == service.ErrAddressNotFound {
			return nil, status.Error(codes.NotFound, "Address not found")
		}
		return nil, status.Errorf(codes.Internal, "Failed to get address: %v", err)
	}
	
	// 检查是否是该用户的地址
	if address.UserID != req.UserId {
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	}
	
	// 转换为proto类型
	protoAddr := &proto.Address{
		Id:          address.ID,
		UserId:      address.UserID,
		Province:    address.Province,
		City:        address.City,
		District:    address.District,
		Address:     address.Address,
		SignerName:  address.SignerName,
		SignerMobile: address.SignerMobile,
		IsDefault:   address.IsDefault,
		Label:       address.Label,
		Postcode:    address.Postcode,
		UsageCount:  int32(address.UsageCount),
		CreatedAt:   address.CreatedAt.Unix(),
	}
	
	if address.LastUsedAt != nil {
		protoAddr.LastUsedAt = address.LastUsedAt.Unix()
	}
	
	return &proto.GetAddressResponse{
		Address: protoAddr,
	}, nil
}

// AddAddress 添加地址
func (s *ProfileServiceServer) AddAddress(ctx context.Context, req *proto.AddAddressRequest) (*proto.AddAddressResponse, error) {
	protoAddr := req.GetAddress()
	if protoAddr == nil {
		return nil, status.Error(codes.InvalidArgument, "Address is required")
	}
	
	// 转换为实体
	address := &entity.Address{
		UserID:      protoAddr.UserId,
		Province:    protoAddr.Province,
		City:        protoAddr.City,
		District:    protoAddr.District,
		Address:     protoAddr.Address,
		SignerName:  protoAddr.SignerName,
		SignerMobile: protoAddr.SignerMobile,
		IsDefault:   protoAddr.IsDefault,
		Label:       protoAddr.Label,
		Postcode:    protoAddr.Postcode,
	}
	
	err := s.addrService.AddAddress(ctx, address)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to add address: %v", err)
	}
	
	return &proto.AddAddressResponse{
		Id:      address.ID,
		Success: true,
		Message: "Address added successfully",
	}, nil
}

// UpdateAddress 更新地址
func (s *ProfileServiceServer) UpdateAddress(ctx context.Context, req *proto.UpdateAddressRequest) (*proto.UpdateAddressResponse, error) {
	protoAddr := req.GetAddress()
	if protoAddr == nil {
		return nil, status.Error(codes.InvalidArgument, "Address is required")
	}
	
	// 转换为实体
	address := &entity.Address{
		ID:          protoAddr.Id,
		UserID:      protoAddr.UserId,
		Province:    protoAddr.Province,
		City:        protoAddr.City,
		District:    protoAddr.District,
		Address:     protoAddr.Address,
		SignerName:  protoAddr.SignerName,
		SignerMobile: protoAddr.SignerMobile,
		IsDefault:   protoAddr.IsDefault,
		Label:       protoAddr.Label,
		Postcode:    protoAddr.Postcode,
	}
	
	err := s.addrService.UpdateAddress(ctx, address)
	if err != nil {
		if err == service.ErrAddressNotFound {
			return nil, status.Error(codes.NotFound, "Address not found")
		} else if err == service.ErrUserNotMatch {
			return nil, status.Error(codes.PermissionDenied, "Permission denied")
		}
		return nil, status.Errorf(codes.Internal, "Failed to update address: %v", err)
	}
	
	return &proto.UpdateAddressResponse{
		Success: true,
		Message: "Address updated successfully",
	}, nil
}

// DeleteAddress 删除地址
func (s *ProfileServiceServer) DeleteAddress(ctx context.Context, req *proto.DeleteAddressRequest) (*proto.DeleteAddressResponse, error) {
	// 先获取地址检查权限
	address, err := s.addrService.GetAddress(ctx, req.Id)
	if err != nil {
		if err == service.ErrAddressNotFound {
			return nil, status.Error(codes.NotFound, "Address not found")
		}
		return nil, status.Errorf(codes.Internal, "Failed to get address: %v", err)
	}
	
	// 检查是否是该用户的地址
	if address.UserID != req.UserId {
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	}
	
	err = s.addrService.DeleteAddress(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete address: %v", err)
	}
	
	return &proto.DeleteAddressResponse{
		Success: true,
		Message: "Address deleted successfully",
	}, nil
}

// SetDefaultAddress 设置默认地址
func (s *ProfileServiceServer) SetDefaultAddress(ctx context.Context, req *proto.SetDefaultAddressRequest) (*proto.SetDefaultAddressResponse, error) {
	err := s.addrService.SetDefault(ctx, req.UserId, req.AddressId)
	if err != nil {
		if err == service.ErrAddressNotFound {
			return nil, status.Error(codes.NotFound, "Address not found")
		} else if err == service.ErrUserNotMatch {
			return nil, status.Error(codes.PermissionDenied, "Permission denied")
		}
		return nil, status.Errorf(codes.Internal, "Failed to set default address: %v", err)
	}
	
	return &proto.SetDefaultAddressResponse{
		Success: true,
		Message: "Default address set successfully",
	}, nil
}

// GetDefaultAddress 获取默认地址
func (s *ProfileServiceServer) GetDefaultAddress(ctx context.Context, req *proto.GetDefaultAddressRequest) (*proto.GetDefaultAddressResponse, error) {
	address, err := s.addrService.GetDefaultAddress(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get default address: %v", err)
	}
	
	// 没有默认地址
	if address == nil {
		return &proto.GetDefaultAddressResponse{
			HasDefault: false,
		}, nil
	}
	
	// 转换为proto类型
	protoAddr := &proto.Address{
		Id:          address.ID,
		UserId:      address.UserID,
		Province:    address.Province,
		City:        address.City,
		District:    address.District,
		Address:     address.Address,
		SignerName:  address.SignerName,
		SignerMobile: address.SignerMobile,
		IsDefault:   address.IsDefault,
		Label:       address.Label,
		Postcode:    address.Postcode,
		UsageCount:  int32(address.UsageCount),
		CreatedAt:   address.CreatedAt.Unix(),
	}
	
	if address.LastUsedAt != nil {
		protoAddr.LastUsedAt = address.LastUsedAt.Unix()
	}
	
	return &proto.GetDefaultAddressResponse{
		Address:    protoAddr,
		HasDefault: true,
	}, nil
}

// 以下是简要实现的其他方法，实际项目中需要完整实现

// ListFeedbacks 获取反馈列表
func (s *ProfileServiceServer) ListFeedbacks(ctx context.Context, req *proto.ListFeedbacksRequest) (*proto.ListFeedbacksResponse, error) {
	// 实现省略，与前面类似
	return &proto.ListFeedbacksResponse{}, nil
}

// GetFeedback 获取反馈详情
func (s *ProfileServiceServer) GetFeedback(ctx context.Context, req *proto.GetFeedbackRequest) (*proto.GetFeedbackResponse, error) {
	// 实现省略
	return &proto.GetFeedbackResponse{}, nil
}

// SubmitFeedback 提交反馈
func (s *ProfileServiceServer) SubmitFeedback(ctx context.Context, req *proto.SubmitFeedbackRequest) (*proto.SubmitFeedbackResponse, error) {
	// 实现省略
	return &proto.SubmitFeedbackResponse{
		Success: true,
		Message: "Feedback submitted successfully",
	}, nil
}

// UpdateFeedback 更新反馈
func (s *ProfileServiceServer) UpdateFeedback(ctx context.Context, req *proto.UpdateFeedbackRequest) (*proto.UpdateFeedbackResponse, error) {
	// 实现省略
	return &proto.UpdateFeedbackResponse{
		Success: true,
		Message: "Feedback updated successfully",
	}, nil
}

// DeleteFeedback 删除反馈
func (s *ProfileServiceServer) DeleteFeedback(ctx context.Context, req *proto.DeleteFeedbackRequest) (*proto.DeleteFeedbackResponse, error) {
	// 实现省略
	return &proto.DeleteFeedbackResponse{
		Success: true,
		Message: "Feedback deleted successfully",
	}, nil
}

// GetBrowsingHistories 获取浏览历史
func (s *ProfileServiceServer) GetBrowsingHistories(ctx context.Context, req *proto.GetBrowsingHistoriesRequest) (*proto.GetBrowsingHistoriesResponse, error) {
	// 实现省略
	return &proto.GetBrowsingHistoriesResponse{}, nil
}

// AddBrowsingHistory 添加浏览历史
func (s *ProfileServiceServer) AddBrowsingHistory(ctx context.Context, req *proto.AddBrowsingHistoryRequest) (*proto.AddBrowsingHistoryResponse, error) {
	// 实现省略
	return &proto.AddBrowsingHistoryResponse{
		Success: true,
	}, nil
}

// RemoveBrowsingHistories 删除浏览历史
func (s *ProfileServiceServer) RemoveBrowsingHistories(ctx context.Context, req *proto.RemoveBrowsingHistoriesRequest) (*proto.RemoveBrowsingHistoriesResponse, error) {
	// 实现省略
	return &proto.RemoveBrowsingHistoriesResponse{
		Success: true,
		Message: "Browsing histories removed successfully",
	}, nil
}

// ClearBrowsingHistories 清空浏览历史
func (s *ProfileServiceServer) ClearBrowsingHistories(ctx context.Context, req *proto.ClearBrowsingHistoriesRequest) (*proto.ClearBrowsingHistoriesResponse, error) {
	// 实现省略
	return &proto.ClearBrowsingHistoriesResponse{
		Success: true,
		Message: "Browsing histories cleared successfully",
	}, nil
}
