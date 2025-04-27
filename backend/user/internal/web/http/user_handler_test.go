package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"shop/backend/user/internal/domain/entity"
	"shop/backend/user/internal/domain/valueobject"
	userhttp "shop/backend/user/internal/web/http"
)

// MockUserService 模拟用户服务
type MockUserService struct {
	mock.Mock
}

// MockAuthService 模拟认证服务
type MockAuthService struct {
	mock.Mock
}

// 实现 UserService 接口
func (m *MockUserService) CreateUser(ctx context.Context, user *entity.User, password string) (*entity.User, error) {
	args := m.Called(ctx, user, password)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserService) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserService) GetUserByPhone(ctx context.Context, phone string) (*entity.User, error) {
	args := m.Called(ctx, phone)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserService) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	args := m.Called(ctx, userID, oldPassword, newPassword)
	return args.Error(0)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) LockUser(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) UnlockUser(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) ListUsers(ctx context.Context, page, pageSize int) ([]*entity.User, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*entity.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserService) UpdateUserStatus(ctx context.Context, id int64, status int) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockUserService) BindWechat(ctx context.Context, userID int64, openID, unionID string) error {
	args := m.Called(ctx, userID, openID, unionID)
	return args.Error(0)
}

func (m *MockUserService) UnbindWechat(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserService) GetUserPermissions(ctx context.Context, userID int64) ([]string, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]string), args.Error(1)
}

// 实现 AuthService 接口
func (m *MockAuthService) Login(ctx context.Context, username, password string) (*entity.User, *valueobject.Credential, error) {
	args := m.Called(ctx, username, password)
	return args.Get(0).(*entity.User), args.Get(1).(*valueobject.Credential), args.Error(2)
}

func (m *MockAuthService) Logout(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*valueobject.Credential, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(*valueobject.Credential), args.Error(1)
}

func (m *MockAuthService) GenerateToken(ctx context.Context, user *entity.User) (*valueobject.Credential, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(*valueobject.Credential), args.Error(1)
}

func (m *MockAuthService) ValidateToken(ctx context.Context, token string) (*valueobject.TokenClaims, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*valueobject.TokenClaims), args.Error(1)
}

func (m *MockAuthService) RegisterUser(ctx context.Context, user *entity.User, password string) (*entity.User, error) {
	args := m.Called(ctx, user, password)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockAuthService) VerifyPassword(ctx context.Context, userID int64, password string) (bool, error) {
	args := m.Called(ctx, userID, password)
	return args.Bool(0), args.Error(1)
}

func (m *MockAuthService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) GenerateVerificationCode(ctx context.Context, phone string, codeType string) (string, error) {
	args := m.Called(ctx, phone, codeType)
	return args.String(0), args.Error(1)
}

func (m *MockAuthService) VerifyVerificationCode(ctx context.Context, phone, code, codeType string) (bool, error) {
	args := m.Called(ctx, phone, code, codeType)
	return args.Bool(0), args.Error(1)
}

// TestUserHandler_Register 测试用户注册
func TestUserHandler_Register(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 创建mock服务
	mockUserService := new(MockUserService)
	mockAuthService := new(MockAuthService)

	// 创建测试用户
	testUser := &entity.User{
		ID:        1,
		Username:  "testuser",
		Nickname:  "testuser",
		Email:     "test@example.com",
		Phone:     "13800138000",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    1,
		Role:      1,
	}

	// 设置mock期望
	mockAuthService.On("RegisterUser", mock.Anything, mock.AnythingOfType("*entity.User"), "password123").Return(testUser, nil)

	// 创建处理器
	handler := userhttp.NewUserHandler(mockUserService, mockAuthService)

	// 创建路由
	router := gin.New()
	router.POST("/api/v1/users/register", handler.Register)

	// 创建请求
	reqBody, _ := json.Marshal(map[string]string{
		"username": "testuser",
		"password": "password123",
		"email":    "test@example.com",
		"phone":    "13800138000",
	})

	req, _ := http.NewRequest("POST", "/api/v1/users/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// 执行请求
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 检查结果
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "注册成功", resp["message"])

	data, _ := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["id"])
	assert.Equal(t, "testuser", data["username"])

	// 验证mock调用
	mockAuthService.AssertExpectations(t)
}

// TestUserHandler_Login 测试用户登录
func TestUserHandler_Login(t *testing.T) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 创建mock服务
	mockUserService := new(MockUserService)
	mockAuthService := new(MockAuthService)

	// 创建测试用户
	testUser := &entity.User{
		ID:        1,
		Username:  "testuser",
		Nickname:  "testuser",
		Email:     "test@example.com",
		Phone:     "13800138000",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    1,
		Role:      1,
	}

	// 创建凭证
	testCredential := &valueobject.Credential{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		ExpiresIn:    time.Hour,
	}

	// 设置mock期望
	mockAuthService.On("Login", mock.Anything, "testuser", "password123").Return(testUser, testCredential, nil)

	// 创建处理器
	handler := userhttp.NewUserHandler(mockUserService, mockAuthService)

	// 创建路由
	router := gin.New()
	router.POST("/api/v1/users/login", handler.Login)

	// 创建请求
	reqBody, _ := json.Marshal(map[string]string{
		"username": "testuser",
		"password": "password123",
	})

	req, _ := http.NewRequest("POST", "/api/v1/users/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// 执行请求
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 检查结果
	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "登录成功", resp["message"])

	data, _ := resp["data"].(map[string]interface{})
	assert.Equal(t, "test-access-token", data["access_token"])
	assert.Equal(t, "test-refresh-token", data["refresh_token"])

	// 验证mock调用
	mockAuthService.AssertExpectations(t)
}
