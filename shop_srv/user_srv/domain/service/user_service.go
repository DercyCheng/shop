package service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"nd/user_srv/proto"
)

// UserService 用户服务接口，定义与用户相关的业务逻辑
type UserService interface {
	// GetUserList 获取用户列表
	GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error)
	
	// GetUserByMobile 通过手机号查询用户
	GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error)
	
	// GetUserById 通过ID查询用户
	GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error)
	
	// CreateUser 创建新用户
	CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error)
	
	// UpdateUser 更新用户信息
	UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*empty.Empty, error)
	
	// CheckPassWord 校验密码
	CheckPassWord(ctx context.Context, req *proto.PasswordCheckInfo) (*proto.CheckResponse, error)
}