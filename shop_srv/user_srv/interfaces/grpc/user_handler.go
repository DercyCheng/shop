package grpc

import (
	"context"
	"nd/user_srv/domain/service"
	"nd/user_srv/proto"

	"github.com/golang/protobuf/ptypes/empty"
)

// UserHandler 实现proto.UserServer接口，作为gRPC服务的入口点
type UserHandler struct {
	proto.UnimplementedUserServer
	userService service.UserService
}

// NewUserHandler 创建gRPC用户处理器
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUserList 获取用户列表
func (h *UserHandler) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	return h.userService.GetUserList(ctx, req)
}

// GetUserByMobile 通过手机号获取用户
func (h *UserHandler) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	return h.userService.GetUserByMobile(ctx, req)
}

// GetUserById 通过ID获取用户
func (h *UserHandler) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	return h.userService.GetUserById(ctx, req)
}

// CreateUser 创建用户
func (h *UserHandler) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	return h.userService.CreateUser(ctx, req)
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*empty.Empty, error) {
	return h.userService.UpdateUser(ctx, req)
}

// CheckPassWord 校验密码
func (h *UserHandler) CheckPassWord(ctx context.Context, req *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	return h.userService.CheckPassWord(ctx, req)
}