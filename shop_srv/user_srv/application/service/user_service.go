package service

import (
	"context"
	"nd/user_srv/domain/repository"
	"nd/user_srv/domain/service"
	"nd/user_srv/infrastructure/persistence"
	"nd/user_srv/model"
	"nd/user_srv/proto"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
)

// UserServiceImpl 实现UserService接口
type UserServiceImpl struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repository.UserRepository) service.UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

// ModelToResponse 将User模型转换为响应
func ModelToResponse(user *model.User) *proto.UserInfoResponse {
	userInfoRsp := &proto.UserInfoResponse{
		Id:       user.ID,
		PassWord: user.Password,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
		Mobile:   user.Mobile,
	}
	if user.Birthday != nil {
		userInfoRsp.BirthDay = uint64(user.Birthday.Unix())
	}
	return userInfoRsp
}

// GetUserList 获取用户列表
func (s *UserServiceImpl) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	// 使用仓储接口获取用户列表
	users, total, err := s.userRepo.GetUserList(ctx, req.Pn, req.PSize)
	if err != nil {
		return nil, err
	}

	// 构造响应
	rsp := &proto.UserListResponse{
		Total: int32(total),
	}

	// 转换用户数据
	for _, user := range users {
		userInfoRsp := ModelToResponse(user)
		rsp.Data = append(rsp.Data, userInfoRsp)
	}

	return rsp, nil
}

// GetUserByMobile 通过手机号查询用户
func (s *UserServiceImpl) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	// 使用仓储接口查询用户
	user, err := s.userRepo.GetUserByMobile(ctx, req.Mobile)
	if err != nil {
		return nil, err
	}

	// 转换并返回用户信息
	userInfoRsp := ModelToResponse(user)
	return userInfoRsp, nil
}

// GetUserById 通过ID查询用户
func (s *UserServiceImpl) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	// 使用仓储接口查询用户
	user, err := s.userRepo.GetUserById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// 转换并返回用户信息
	userInfoRsp := ModelToResponse(user)
	return userInfoRsp, nil
}

// CreateUser 创建新用户
func (s *UserServiceImpl) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	// 创建用户对象
	user := &model.User{
		Mobile:   req.Mobile,
		NickName: req.NickName,
		Password: persistence.EncryptPassword(req.PassWord),
	}

	// 调用仓储接口创建用户
	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	// 转换并返回用户信息
	return ModelToResponse(user), nil
}

// UpdateUser 更新用户信息
func (s *UserServiceImpl) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*empty.Empty, error) {
	// 先查询用户是否存在
	user, err := s.userRepo.GetUserById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// 更新用户信息
	birthDay := time.Unix(int64(req.BirthDay), 0)
	user.NickName = req.NickName
	user.Birthday = &birthDay
	user.Gender = req.Gender

	// 调用仓储接口更新用户
	if err := s.userRepo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

// CheckPassWord 校验密码
func (s *UserServiceImpl) CheckPassWord(ctx context.Context, req *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	// 调用仓储接口校验密码
	check, err := s.userRepo.CheckPassword(ctx, req.Password, req.EncryptedPassword)
	if err != nil {
		return nil, err
	}

	return &proto.CheckResponse{Success: check}, nil
}