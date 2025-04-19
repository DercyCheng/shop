package repository

import (
	"context"
	"nd/user_srv/model"
	"nd/user_srv/proto"
)

// UserRepository 用户仓储接口，定义与用户相关的数据访问方法
type UserRepository interface {
	// GetUserList 获取用户列表并支持分页
	GetUserList(ctx context.Context, page, pageSize int32) ([]*model.User, int64, error)
	
	// GetUserByMobile 通过手机号查询用户
	GetUserByMobile(ctx context.Context, mobile string) (*model.User, error)
	
	// GetUserById 通过ID查询用户
	GetUserById(ctx context.Context, id int32) (*model.User, error)
	
	// CreateUser 创建新用户
	CreateUser(ctx context.Context, user *model.User) error
	
	// UpdateUser 更新用户信息
	UpdateUser(ctx context.Context, user *model.User) error
	
	// CheckPassword 验证密码
	CheckPassword(ctx context.Context, password, encryptedPassword string) (bool, error)
}