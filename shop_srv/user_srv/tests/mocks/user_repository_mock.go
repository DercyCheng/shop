package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"nd/user_srv/domain/repository"
	"nd/user_srv/model"
)

// MockUserRepository 实现 UserRepository 接口的模拟类
type MockUserRepository struct {
	mock.Mock
}

// Ensure MockUserRepository implements UserRepository interface
var _ repository.UserRepository = (*MockUserRepository)(nil)

// GetUserList 模拟获取用户列表
func (m *MockUserRepository) GetUserList(ctx context.Context, page, pageSize int32) ([]*model.User, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*model.User), args.Get(1).(int64), args.Error(2)
}

// GetUserByMobile 模拟通过手机号查询用户
func (m *MockUserRepository) GetUserByMobile(ctx context.Context, mobile string) (*model.User, error) {
	args := m.Called(ctx, mobile)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// GetUserById 模拟通过ID查询用户
func (m *MockUserRepository) GetUserById(ctx context.Context, id int32) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// CreateUser 模拟创建新用户
func (m *MockUserRepository) CreateUser(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// UpdateUser 模拟更新用户信息
func (m *MockUserRepository) UpdateUser(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// CheckPassword 模拟验证密码
func (m *MockUserRepository) CheckPassword(ctx context.Context, password, encryptedPassword string) (bool, error) {
	args := m.Called(ctx, password, encryptedPassword)
	return args.Bool(0), args.Error(1)
}