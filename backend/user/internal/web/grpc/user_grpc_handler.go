package grpc

import (
	"context"
	
	"shop/backend/user/api/proto"
	"shop/backend/user/internal/service"
	"shop/backend/user/internal/domain/valueobject"
	
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserServiceServer gRPC服务实现
type UserServiceServer struct {
	userService service.UserService
	authService service.AuthService
	proto.UnimplementedUserServiceServer
}

// NewUserServiceServer 创建用户服务gRPC处理器
func NewUserServiceServer(userService service.UserService, authService service.AuthService) *UserServiceServer {
	return &UserServiceServer{
		userService: userService,
		authService: authService,
	}
}

// RegisterWithServer 注册服务到gRPC服务器
func (s *UserServiceServer) RegisterWithServer(server *grpc.Server) {
	proto.RegisterUserServiceServer(server, s)
}

// ValidateToken 验证令牌
func (s *UserServiceServer) ValidateToken(ctx context.Context, req *proto.ValidateTokenRequest) (*proto.ValidateTokenResponse, error) {
	// 调用认证服务验证令牌
	userID, role, err := s.authService.ValidateToken(ctx, req.Token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token: %v", err)
	}
	
	return &proto.ValidateTokenResponse{
		UserId: userID,
		Role:   role,
		Valid:  true,
	}, nil
}

// GetUserInfo 获取用户信息
func (s *UserServiceServer) GetUserInfo(ctx context.Context, req *proto.GetUserInfoRequest) (*proto.GetUserInfoResponse, error) {
	// 根据用户ID获取用户信息
	user, err := s.userService.GetUserByID(ctx, req.UserId)
	if err != nil {
		if err == service.ErrUserNotFound {
			return nil, status.Error(codes.NotFound, "User not found")
		}
		return nil, status.Errorf(codes.Internal, "Failed to get user: %v", err)
	}
	
	// 构造响应
	return &proto.GetUserInfoResponse{
		UserId:    user.ID,
		Mobile:    user.Mobile,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Gender:    user.Gender,
		Role:      int32(user.Role),
		Status:    int32(user.Status),
		CreatedAt: user.CreatedAt.Unix(),
	}, nil
}

// RefreshToken 刷新令牌
func (s *UserServiceServer) RefreshToken(ctx context.Context, req *proto.RefreshTokenRequest) (*proto.RefreshTokenResponse, error) {
	// 调用认证服务刷新令牌
	response, err := s.authService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid refresh token: %v", err)
	}
	
	// 构造响应
	return &proto.RefreshTokenResponse{
		Token:        response.Token,
		RefreshToken: response.RefreshToken,
		ExpiresIn:    response.ExpireTime,
	}, nil
}

// Login 用户登录
func (s *UserServiceServer) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	// 创建凭证对象
	credential := &valueobject.Credential{
		Username: req.Mobile,
		Password: req.Password,
		Type:     valueobject.CredentialTypePassword,
	}
	
	// 调用认证服务进行登录
	response, err := s.authService.Login(ctx, credential)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			return nil, status.Error(codes.NotFound, "User not found")
		case service.ErrInvalidCredentials:
			return nil, status.Error(codes.Unauthenticated, "Invalid credentials")
		case service.ErrUserDisabled:
			return nil, status.Error(codes.PermissionDenied, "User is disabled")
		case service.ErrUserLocked:
			return nil, status.Error(codes.PermissionDenied, "User is locked")
		default:
			return nil, status.Errorf(codes.Internal, "Login failed: %v", err)
		}
	}
	
	// 获取用户信息
	userInfo := response.UserInfo.(map[string]interface{})
	
	// 构造响应
	return &proto.LoginResponse{
		Token:        response.Token,
		RefreshToken: response.RefreshToken,
		ExpiresIn:    response.ExpireTime,
		UserInfo: &proto.UserInfo{
			UserId:   int64(userInfo["id"].(int64)),
			Mobile:   userInfo["mobile"].(string),
			Nickname: userInfo["nickname"].(string),
			Avatar:   userInfo["avatar"].(string),
			Role:     int32(userInfo["role"].(int)),
		},
	}, nil
}

// Register 用户注册
func (s *UserServiceServer) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	// 调用用户服务进行注册
	user, err := s.userService.RegisterUser(ctx, req.Mobile, req.Password, req.Nickname)
	if err != nil {
		if err == service.ErrMobileExists {
			return nil, status.Error(codes.AlreadyExists, "Mobile already exists")
		}
		return nil, status.Errorf(codes.Internal, "Registration failed: %v", err)
	}
	
	// 构造响应
	return &proto.RegisterResponse{
		UserId:   user.ID,
		Mobile:   user.Mobile,
		Nickname: user.Nickname,
	}, nil
}
