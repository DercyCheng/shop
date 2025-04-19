package persistence

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"gorm.io/gorm"
	"nd/user_srv/domain/repository"
	"nd/user_srv/model"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserRepositoryImpl 用户仓储接口的具体实现
type UserRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

// paginate 分页辅助函数
func paginate(page, pageSize int32) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(int(offset)).Limit(int(pageSize))
	}
}

// GetUserList 获取用户列表并支持分页
func (r *UserRepositoryImpl) GetUserList(ctx context.Context, page, pageSize int32) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	// 查询总数
	if err := r.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := r.db.Scopes(paginate(page, pageSize)).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUserByMobile 通过手机号查询用户
func (r *UserRepositoryImpl) GetUserByMobile(ctx context.Context, mobile string) (*model.User, error) {
	var user model.User
	result := r.db.Where(&model.User{Mobile: mobile}).First(&user)
	
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	return &user, nil
}

// GetUserById 通过ID查询用户
func (r *UserRepositoryImpl) GetUserById(ctx context.Context, id int32) (*model.User, error) {
	var user model.User
	result := r.db.First(&user, id)
	
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	return &user, nil
}

// CreateUser 创建新用户
func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user *model.User) error {
	// 检查手机号是否已存在
	var existUser model.User
	result := r.db.Where(&model.User{Mobile: user.Mobile}).First(&existUser)
	if result.RowsAffected == 1 {
		return status.Errorf(codes.AlreadyExists, "用户已存在")
	}

	// 创建用户
	if err := r.db.Create(user).Error; err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	return nil
}

// UpdateUser 更新用户信息
func (r *UserRepositoryImpl) UpdateUser(ctx context.Context, user *model.User) error {
	// 检查用户是否存在
	var existUser model.User
	result := r.db.First(&existUser, user.ID)
	if result.RowsAffected == 0 {
		return status.Errorf(codes.NotFound, "用户不存在")
	}

	// 更新用户
	if err := r.db.Save(user).Error; err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	return nil
}

// CheckPassword 验证密码
func (r *UserRepositoryImpl) CheckPassword(ctx context.Context, password, encryptedPassword string) (bool, error) {
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	passwordInfo := strings.Split(encryptedPassword, "$")
	
	if len(passwordInfo) != 4 {
		return false, status.Errorf(codes.Internal, "密码格式错误")
	}
	
	check := password.Verify(password, passwordInfo[2], passwordInfo[3], options)
	return check, nil
}

// EncryptPassword 加密密码
func EncryptPassword(plainPassword string) string {
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode(plainPassword, options)
	return fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)
}