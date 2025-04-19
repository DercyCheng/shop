package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"nd/user_srv/application/service"
	"nd/user_srv/model"
	"nd/user_srv/proto"
	"nd/user_srv/tests/mocks"
	"testing"
	"time"
)

func TestUserService_GetUserById(t *testing.T) {
	// 准备测试数据
	mockUser := &model.User{
		ID:       1,
		Mobile:   "13800138000",
		Password: "encrypted_password",
		NickName: "测试用户",
		Gender:   "male",
		Role:     1,
	}

	// 创建仓储层的Mock
	mockRepo := new(mocks.MockUserRepository)
	mockRepo.On("GetUserById", mock.Anything, int32(1)).Return(mockUser, nil)

	// 创建服务实例，注入Mock的仓储
	userService := service.NewUserService(mockRepo)

	// 执行被测试的方法
	resp, err := userService.GetUserById(context.Background(), &proto.IdRequest{Id: 1})

	// 断言结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(1), resp.Id)
	assert.Equal(t, "测试用户", resp.NickName)
	assert.Equal(t, "13800138000", resp.Mobile)

	// 验证Mock的方法被正确调用
	mockRepo.AssertExpectations(t)
}

func TestUserService_CreateUser(t *testing.T) {
	// 创建仓储层的Mock
	mockRepo := new(mocks.MockUserRepository)
	
	// 设置Mock行为：保存用户时不返回错误
	mockRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(user *model.User) bool {
		return user.Mobile == "13900139000" && user.NickName == "新用户"
	})).Return(nil)

	// 创建服务实例，注入Mock的仓储
	userService := service.NewUserService(mockRepo)

	// 执行被测试的方法
	resp, err := userService.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: "新用户",
		PassWord: "password123",
		Mobile:   "13900139000",
	})

	// 断言结果
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "新用户", resp.NickName)
	assert.Equal(t, "13900139000", resp.Mobile)

	// 验证Mock的方法被正确调用
	mockRepo.AssertExpectations(t)
}

func TestUserService_UpdateUser(t *testing.T) {
	// 准备测试数据
	birthday := time.Now()
	mockUser := &model.User{
		ID:       2,
		Mobile:   "13800138001",
		NickName: "旧用户名",
		Gender:   "female",
	}

	// 创建仓储层的Mock
	mockRepo := new(mocks.MockUserRepository)
	mockRepo.On("GetUserById", mock.Anything, int32(2)).Return(mockUser, nil)
	mockRepo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(user *model.User) bool {
		return user.ID == 2 && user.NickName == "新用户名" && user.Gender == "male"
	})).Return(nil)

	// 创建服务实例，注入Mock的仓储
	userService := service.NewUserService(mockRepo)

	// 执行被测试的方法
	_, err := userService.UpdateUser(context.Background(), &proto.UpdateUserInfo{
		Id:       2,
		NickName: "新用户名",
		Gender:   "male",
		BirthDay: uint64(birthday.Unix()),
	})

	// 断言结果
	assert.NoError(t, err)

	// 验证Mock的方法被正确调用
	mockRepo.AssertExpectations(t)
}