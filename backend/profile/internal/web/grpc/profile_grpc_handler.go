package grpc

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"shop/backend/profile/api/proto"
	"shop/backend/profile/internal/domain/entity"
	"shop/backend/profile/internal/service"
)

// ProfileGRPCServer implements the ProfileService gRPC interface
type ProfileGRPCServer struct {
	proto.UnimplementedProfileServiceServer
	userFavService service.UserFavService
	addressService service.AddressService
	messageService service.MessageService
	logger         *zap.Logger
}

// NewProfileGRPCServer creates a new instance of the ProfileService gRPC server
func NewProfileGRPCServer(
	userFavService service.UserFavService,
	addressService service.AddressService,
	messageService service.MessageService,
	logger *zap.Logger,
) *ProfileGRPCServer {
	return &ProfileGRPCServer{
		userFavService: userFavService,
		addressService: addressService,
		messageService: messageService,
		logger:         logger,
	}
}

// GetFavList returns the list of user favorites
func (s *ProfileGRPCServer) GetFavList(ctx context.Context, req *proto.UserFavRequest) (*proto.UserFavListResponse, error) {
	s.logger.Info("GetFavList called", zap.Int64("user_id", req.UserId))

	// Set default page size if not provided
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	// Set default page if not provided
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}

	// Get user favorites from service
	favs, total, err := s.userFavService.GetUserFavList(ctx, req.UserId, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to get user favorites", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get user favorites: %v", err)
	}

	// Convert domain entities to proto messages
	favResponses := make([]*proto.UserFavResponse, 0, len(favs))
	for _, fav := range favs {
		favResp := &proto.UserFavResponse{
			Id:        fav.ID,
			UserId:    fav.UserID,
			GoodsId:   fav.GoodsID,
			CreatedAt: fav.CreatedAt.Format(time.RFC3339),
		}

		// Add goods info if available
		if fav.GoodsInfo != nil {
			favResp.GoodsInfo = &proto.GoodsInfoResponse{
				Id:           fav.GoodsInfo.ID,
				Name:         fav.GoodsInfo.Name,
				ShopPrice:    fav.GoodsInfo.ShopPrice,
				Image:        fav.GoodsInfo.Image,
				CategoryName: fav.GoodsInfo.CategoryName,
				BrandName:    fav.GoodsInfo.BrandName,
			}
		}

		favResponses = append(favResponses, favResp)
	}

	return &proto.UserFavListResponse{
		Total: int32(total),
		Data:  favResponses,
	}, nil
}

// AddUserFav adds a product to user favorites
func (s *ProfileGRPCServer) AddUserFav(ctx context.Context, req *proto.UserFavRequest) (*empty.Empty, error) {
	s.logger.Info("AddUserFav called",
		zap.Int64("user_id", req.UserId),
		zap.Int64("goods_id", req.GoodsId))

	if err := s.userFavService.AddUserFav(ctx, req.UserId, req.GoodsId); err != nil {
		s.logger.Error("Failed to add user favorite", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to add user favorite: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// DeleteUserFav removes a product from user favorites
func (s *ProfileGRPCServer) DeleteUserFav(ctx context.Context, req *proto.UserFavRequest) (*empty.Empty, error) {
	s.logger.Info("DeleteUserFav called",
		zap.Int64("user_id", req.UserId),
		zap.Int64("goods_id", req.GoodsId))

	if err := s.userFavService.DeleteUserFav(ctx, req.UserId, req.GoodsId); err != nil {
		s.logger.Error("Failed to delete user favorite", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to delete user favorite: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// GetUserFavDetail checks if a product is in user favorites
func (s *ProfileGRPCServer) GetUserFavDetail(ctx context.Context, req *proto.UserFavRequest) (*proto.UserFavDetailResponse, error) {
	s.logger.Info("GetUserFavDetail called",
		zap.Int64("user_id", req.UserId),
		zap.Int64("goods_id", req.GoodsId))

	isFav, err := s.userFavService.CheckUserFav(ctx, req.UserId, req.GoodsId)
	if err != nil {
		s.logger.Error("Failed to check user favorite", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to check user favorite: %v", err)
	}

	return &proto.UserFavDetailResponse{
		IsFav: isFav,
	}, nil
}

// GetAddressList returns the list of user addresses
func (s *ProfileGRPCServer) GetAddressList(ctx context.Context, req *proto.AddressRequest) (*proto.AddressListResponse, error) {
	s.logger.Info("GetAddressList called", zap.Int64("user_id", req.UserId))

	// Set default page size if not provided
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	// Set default page if not provided
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}

	// Get addresses from service
	addresses, total, err := s.addressService.GetAddressList(ctx, req.UserId, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to get addresses", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get addresses: %v", err)
	}

	// Convert domain entities to proto messages
	addressResponses := make([]*proto.AddressResponse, 0, len(addresses))
	for _, addr := range addresses {
		addressResponses = append(addressResponses, &proto.AddressResponse{
			Id:           addr.ID,
			UserId:       addr.UserID,
			Province:     addr.Province,
			City:         addr.City,
			District:     addr.District,
			Address:      addr.Address,
			SignerName:   addr.SignerName,
			SignerMobile: addr.SignerMobile,
			IsDefault:    addr.IsDefault,
		})
	}

	return &proto.AddressListResponse{
		Total: int32(total),
		Data:  addressResponses,
	}, nil
}

// CreateAddress creates a new address
func (s *ProfileGRPCServer) CreateAddress(ctx context.Context, req *proto.AddressRequest) (*proto.AddressResponse, error) {
	s.logger.Info("CreateAddress called", zap.Int64("user_id", req.UserId))

	// Convert proto message to domain entity
	address := &entity.Address{
		UserID:       req.UserId,
		Province:     req.Province,
		City:         req.City,
		District:     req.District,
		Address:      req.Address,
		SignerName:   req.SignerName,
		SignerMobile: req.SignerMobile,
		IsDefault:    req.IsDefault,
	}

	// Create address
	createdAddress, err := s.addressService.CreateAddress(ctx, address)
	if err != nil {
		s.logger.Error("Failed to create address", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create address: %v", err)
	}

	// Convert domain entity to proto message
	return &proto.AddressResponse{
		Id:           createdAddress.ID,
		UserId:       createdAddress.UserID,
		Province:     createdAddress.Province,
		City:         createdAddress.City,
		District:     createdAddress.District,
		Address:      createdAddress.Address,
		SignerName:   createdAddress.SignerName,
		SignerMobile: createdAddress.SignerMobile,
		IsDefault:    createdAddress.IsDefault,
	}, nil
}

// DeleteAddress deletes an address
func (s *ProfileGRPCServer) DeleteAddress(ctx context.Context, req *proto.AddressRequest) (*empty.Empty, error) {
	s.logger.Info("DeleteAddress called",
		zap.Int64("address_id", req.Id),
		zap.Int64("user_id", req.UserId))

	if err := s.addressService.DeleteAddress(ctx, req.Id, req.UserId); err != nil {
		s.logger.Error("Failed to delete address", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to delete address: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// UpdateAddress updates an address
func (s *ProfileGRPCServer) UpdateAddress(ctx context.Context, req *proto.AddressRequest) (*empty.Empty, error) {
	s.logger.Info("UpdateAddress called",
		zap.Int64("address_id", req.Id),
		zap.Int64("user_id", req.UserId))

	// Convert proto message to domain entity
	address := &entity.Address{
		ID:           req.Id,
		UserID:       req.UserId,
		Province:     req.Province,
		City:         req.City,
		District:     req.District,
		Address:      req.Address,
		SignerName:   req.SignerName,
		SignerMobile: req.SignerMobile,
		IsDefault:    req.IsDefault,
	}

	// Update address
	if err := s.addressService.UpdateAddress(ctx, address); err != nil {
		s.logger.Error("Failed to update address", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to update address: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// GetMessageList returns the list of user messages
func (s *ProfileGRPCServer) GetMessageList(ctx context.Context, req *proto.MessageRequest) (*proto.MessageListResponse, error) {
	s.logger.Info("GetMessageList called", zap.Int64("user_id", req.UserId))

	// Set default page size if not provided
	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	// Set default page if not provided
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}

	// Get messages from service
	messages, total, err := s.messageService.GetMessagesByUserID(ctx, req.UserId, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to get messages", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to get messages: %v", err)
	}

	// Convert domain entities to proto messages
	messageResponses := make([]*proto.MessageResponse, 0, len(messages))
	for _, msg := range messages {
		messageResponses = append(messageResponses, &proto.MessageResponse{
			Id:          msg.ID,
			UserId:      msg.UserID,
			MessageType: int32(msg.MessageType),
			Subject:     msg.Subject,
			Message:     msg.Content,
			File:        msg.File,
			Images:      msg.Images,
			CreatedAt:   msg.CreatedAt.Format(time.RFC3339),
		})
	}

	return &proto.MessageListResponse{
		Total: int32(total),
		Data:  messageResponses,
	}, nil
}

// CreateMessage creates a new message
func (s *ProfileGRPCServer) CreateMessage(ctx context.Context, req *proto.MessageRequest) (*proto.MessageResponse, error) {
	s.logger.Info("CreateMessage called",
		zap.Int64("user_id", req.UserId),
		zap.Int32("message_type", req.MessageType))

	// Convert proto message to domain entity
	message := &entity.Message{
		UserID:      req.UserId,
		MessageType: int(req.MessageType),
		Subject:     req.Subject,
		Content:     req.Message,
		File:        req.File,
		Images:      req.Images,
	}

	// Create message
	if err := s.messageService.CreateMessage(ctx, message); err != nil {
		s.logger.Error("Failed to create message", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create message: %v", err)
	}

	// Convert domain entity to proto message
	return &proto.MessageResponse{
		Id:          message.ID,
		UserId:      message.UserID,
		MessageType: int32(message.MessageType),
		Subject:     message.Subject,
		Message:     message.Content,
		File:        message.File,
		Images:      message.Images,
		CreatedAt:   message.CreatedAt.Format(time.RFC3339),
	}, nil
}
