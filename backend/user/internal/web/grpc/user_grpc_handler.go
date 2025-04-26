package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "shop/backend/user/api/proto/user"
	"shop/backend/user/internal/domain/entity"
	"shop/backend/user/internal/service"
)

// UserGRPCServer 用户服务gRPC实现
type UserGRPCServer struct {
	pb.UnimplementedUserServiceServer
	userService service.UserService
	authService service.AuthService
}

// NewUserGRPCServer 创建用户服务gRPC服务器
func NewUserGRPCServer(userService service.UserService, authService service.AuthService) *UserGRPCServer {
	return &UserGRPCServer{
		userService: userService,
		authService: authService,
	}
}

// Register 用户注册
func (s *UserGRPCServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.UserInfo, error) {
	// 构建用户实体
	user := &entity.User{
		Username: req.Username,
		Nickname: req.Username, // 默认昵称与用户名相同
		Email:    req.Email,
		Phone:    req.Phone,
	}

	// 调用注册服务
	createdUser, err := s.authService.RegisterUser(ctx, user, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "注册失败: %v", err)
	}

	// 转换为响应对象
	return entityToUserInfo(createdUser), nil
}

// Login 用户登录
func (s *UserGRPCServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// 调用登录服务
	user, credential, err := s.authService.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "登录失败: %v", err)
	}

	// 构建响应
	response := &pb.LoginResponse{
		UserInfo: entityToUserInfo(user),
		Token: &pb.TokenResponse{
			AccessToken:  credential.AccessToken,
			RefreshToken: credential.RefreshToken,
			ExpiresIn:    int64(credential.ExpiresIn.Seconds()),
		},
	}

	return response, nil
}

// Logout 用户登出
func (s *UserGRPCServer) Logout(ctx context.Context, req *pb.LogoutRequest) (*emptypb.Empty, error) {
	// 调用登出服务
	err := s.authService.Logout(ctx, req.AccessToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "登出失败: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// RefreshToken 刷新令牌
func (s *UserGRPCServer) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.TokenResponse, error) {
	// 调用刷新令牌服务
	credential, err := s.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "刷新令牌失败: %v", err)
	}

	// 构建响应
	response := &pb.TokenResponse{
		AccessToken:  credential.AccessToken,
		RefreshToken: credential.RefreshToken,
		ExpiresIn:    int64(credential.ExpiresIn.Seconds()),
	}

	return response, nil
}

// GetUserInfo 获取用户信息
func (s *UserGRPCServer) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.UserInfo, error) {
	// 获取用户信息
	user, err := s.userService.GetUserByID(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "用户不存在: %v", err)
	}

	// 转换为响应对象
	return entityToUserInfo(user), nil
}

// UpdateUserInfo 更新用户信息
func (s *UserGRPCServer) UpdateUserInfo(ctx context.Context, req *pb.UpdateUserInfoRequest) (*pb.UserInfo, error) {
	// 首先获取当前用户信息
	user, err := s.userService.GetUserByID(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "用户不存在: %v", err)
	}

	// 更新用户信息
	if req.Nickname != nil {
		user.Nickname = req.Nickname.Value
	}
	if req.Avatar != nil {
		user.Avatar = req.Avatar.Value
	}
	if req.Email != nil {
		user.Email = req.Email.Value
	}
	if req.Phone != nil {
		user.Phone = req.Phone.Value
	}
	if req.Gender != nil {
		switch req.Gender.Value {
		case pb.Gender_GENDER_MALE:
			user.Gender = "male"
		case pb.Gender_GENDER_FEMALE:
			user.Gender = "female"
		default:
			user.Gender = "unknown"
		}
	}
	if req.Birthday != nil {
		user.Birthday = req.Birthday.AsTime()
	}

	// 调用更新服务
	updatedUser, err := s.userService.UpdateUser(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "更新用户信息失败: %v", err)
	}

	// 转换为响应对象
	return entityToUserInfo(updatedUser), nil
}

// ChangePassword 修改密码
func (s *UserGRPCServer) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*emptypb.Empty, error) {
	// 调用修改密码服务
	err := s.userService.ChangePassword(ctx, req.UserId, req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "修改密码失败: %v", err)
	}

	return &emptypb.Empty{}, nil
}

// ValidateToken 验证令牌
func (s *UserGRPCServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	// 调用验证令牌服务
	claims, err := s.authService.ValidateToken(ctx, req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	// 获取用户权限
	permissions, err := s.userService.GetUserPermissions(ctx, claims.UserID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取用户权限失败: %v", err)
	}

	// 构建响应
	return &pb.ValidateTokenResponse{
		Valid:       true,
		UserId:      claims.UserID,
		Permissions: permissions,
	}, nil
}

// GetUserByID 通过ID获取用户
func (s *UserGRPCServer) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.UserInfo, error) {
	// 获取用户信息
	user, err := s.userService.GetUserByID(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "用户不存在: %v", err)
	}

	// 转换为响应对象
	return entityToUserInfo(user), nil
}

// 将用户实体转换为gRPC响应对象
func entityToUserInfo(user *entity.User) *pb.UserInfo {
	var gender pb.Gender
	switch user.Gender {
	case "male":
		gender = pb.Gender_GENDER_MALE
	case "female":
		gender = pb.Gender_GENDER_FEMALE
	default:
		gender = pb.Gender_GENDER_UNKNOWN
	}

	return &pb.UserInfo{
		Id:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Email:     user.Email,
		Phone:     user.Phone,
		Gender:    gender,
		Birthday:  timestamppb.New(user.Birthday),
		Status:    int32(user.Status),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
