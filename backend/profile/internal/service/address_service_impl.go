package service

import (
	"context"
	"errors"
	"time"
	
	"shop/backend/profile/internal/domain/entity"
)

// ErrAddressNotFound 地址不存在错误
var ErrAddressNotFound = errors.New("address not found")

// AddressServiceImpl 地址服务实现
type AddressServiceImpl struct {
	repo ProfileRepository
}

// NewAddressService 创建地址服务实例
func NewAddressService(repo ProfileRepository) AddressService {
	return &AddressServiceImpl{
		repo: repo,
	}
}

// ListAddresses 获取用户地址列表
func (s *AddressServiceImpl) ListAddresses(ctx context.Context, userID int64) ([]*entity.Address, error) {
	return s.repo.GetAddressesByUserID(ctx, userID)
}

// GetAddress 获取地址详情
func (s *AddressServiceImpl) GetAddress(ctx context.Context, id int64) (*entity.Address, error) {
	address, err := s.repo.GetAddressByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if address == nil {
		return nil, ErrAddressNotFound
	}
	
	return address, nil
}

// AddAddress 添加地址
func (s *AddressServiceImpl) AddAddress(ctx context.Context, address *entity.Address) error {
	now := time.Now()
	address.CreatedAt = now
	address.UpdatedAt = now
	
	// 如果是设置为默认地址
	if address.IsDefault {
		// 先重置该用户的其他默认地址
		oldDefault, err := s.repo.GetDefaultAddress(ctx, address.UserID)
		if err != nil {
			return err
		}
		
		if oldDefault != nil {
			oldDefault.IsDefault = false
			if err := s.repo.UpdateAddress(ctx, oldDefault); err != nil {
				return err
			}
		}
	}
	
	return s.repo.CreateAddress(ctx, address)
}

// UpdateAddress 更新地址
func (s *AddressServiceImpl) UpdateAddress(ctx context.Context, address *entity.Address) error {
	// 检查地址是否存在
	existingAddress, err := s.repo.GetAddressByID(ctx, address.ID)
	if err != nil {
		return err
	}
	
	if existingAddress == nil {
		return ErrAddressNotFound
	}
	
	// 如果不是同一个用户的地址，不允许修改
	if existingAddress.UserID != address.UserID {
		return ErrUserNotMatch
	}
	
	// 更新时间
	address.UpdatedAt = time.Now()
	address.CreatedAt = existingAddress.CreatedAt
	
	// 如果设置为默认地址，需要重置其他默认地址
	if address.IsDefault && !existingAddress.IsDefault {
		return s.SetDefault(ctx, address.UserID, address.ID)
	} else {
		return s.repo.UpdateAddress(ctx, address)
	}
}

// DeleteAddress 删除地址
func (s *AddressServiceImpl) DeleteAddress(ctx context.Context, id int64) error {
	// 检查地址是否存在
	address, err := s.repo.GetAddressByID(ctx, id)
	if err != nil {
		return err
	}
	
	if address == nil {
		return ErrAddressNotFound
	}
	
	return s.repo.DeleteAddress(ctx, id)
}

// SetDefault 设置默认地址
func (s *AddressServiceImpl) SetDefault(ctx context.Context, userID, addressID int64) error {
	// 检查地址是否存在
	address, err := s.repo.GetAddressByID(ctx, addressID)
	if err != nil {
		return err
	}
	
	if address == nil {
		return ErrAddressNotFound
	}
	
	// 检查地址是否属于该用户
	if address.UserID != userID {
		return ErrUserNotMatch
	}
	
	return s.repo.SetDefaultAddress(ctx, userID, addressID)
}

// GetDefaultAddress 获取默认地址
func (s *AddressServiceImpl) GetDefaultAddress(ctx context.Context, userID int64) (*entity.Address, error) {
	return s.repo.GetDefaultAddress(ctx, userID)
}

// UseAddress 使用地址（更新使用次数和最后使用时间）
func (s *AddressServiceImpl) UseAddress(ctx context.Context, id int64) error {
	// 检查地址是否存在
	address, err := s.repo.GetAddressByID(ctx, id)
	if err != nil {
		return err
	}
	
	if address == nil {
		return ErrAddressNotFound
	}
	
	// 更新使用次数和最后使用时间
	now := time.Now()
	address.UsageCount++
	address.LastUsedAt = &now
	address.UpdatedAt = now
	
	return s.repo.UpdateAddress(ctx, address)
}
