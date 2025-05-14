package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	
	"shop/backend/profile/internal/domain/entity"
	
	"github.com/go-redis/redis/v8"
)

// ProfileCache 个人信息缓存接口
type ProfileCache interface {
	// 收藏相关
	GetFavoriteCountByUser(ctx context.Context, userID int64) (int64, error)
	SetFavoriteCountByUser(ctx context.Context, userID int64, count int64) error
	DeleteFavoriteCache(ctx context.Context, userID int64) error
	
	// 地址相关
	GetAddressByID(ctx context.Context, id int64) (*entity.Address, error)
	SetAddress(ctx context.Context, address *entity.Address) error
	DeleteAddressCache(ctx context.Context, id int64) error
	GetDefaultAddressByUser(ctx context.Context, userID int64) (*entity.Address, error)
	SetDefaultAddressByUser(ctx context.Context, userID int64, address *entity.Address) error
	DeleteDefaultAddressCache(ctx context.Context, userID int64) error
}

// RedisProfileCache Redis实现的个人信息缓存
type RedisProfileCache struct {
	client           *redis.Client
	expiration       time.Duration
	favoritePrefix   string
	addressPrefix    string
	defaultAddrPrefix string
}

// NewRedisProfileCache 创建Redis个人信息缓存实例
func NewRedisProfileCache(client *redis.Client, expiration time.Duration) ProfileCache {
	return &RedisProfileCache{
		client:           client,
		expiration:       expiration,
		favoritePrefix:   "profile:fav:user:",
		addressPrefix:    "profile:address:id:",
		defaultAddrPrefix: "profile:address:default:user:",
	}
}

// GetFavoriteCountByUser 获取用户收藏数量
func (c *RedisProfileCache) GetFavoriteCountByUser(ctx context.Context, userID int64) (int64, error) {
	key := fmt.Sprintf("%s%d:count", c.favoritePrefix, userID)
	val, err := c.client.Get(ctx, key).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return val, nil
}

// SetFavoriteCountByUser 设置用户收藏数量
func (c *RedisProfileCache) SetFavoriteCountByUser(ctx context.Context, userID int64, count int64) error {
	key := fmt.Sprintf("%s%d:count", c.favoritePrefix, userID)
	return c.client.Set(ctx, key, count, c.expiration).Err()
}

// DeleteFavoriteCache 删除用户收藏缓存
func (c *RedisProfileCache) DeleteFavoriteCache(ctx context.Context, userID int64) error {
	key := fmt.Sprintf("%s%d:count", c.favoritePrefix, userID)
	return c.client.Del(ctx, key).Err()
}

// GetAddressByID 根据ID获取地址
func (c *RedisProfileCache) GetAddressByID(ctx context.Context, id int64) (*entity.Address, error) {
	key := fmt.Sprintf("%s%d", c.addressPrefix, id)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	
	var address entity.Address
	if err := json.Unmarshal(data, &address); err != nil {
		return nil, err
	}
	
	return &address, nil
}

// SetAddress 缓存地址
func (c *RedisProfileCache) SetAddress(ctx context.Context, address *entity.Address) error {
	key := fmt.Sprintf("%s%d", c.addressPrefix, address.ID)
	data, err := json.Marshal(address)
	if err != nil {
		return err
	}
	
	return c.client.Set(ctx, key, data, c.expiration).Err()
}

// DeleteAddressCache 删除地址缓存
func (c *RedisProfileCache) DeleteAddressCache(ctx context.Context, id int64) error {
	key := fmt.Sprintf("%s%d", c.addressPrefix, id)
	return c.client.Del(ctx, key).Err()
}

// GetDefaultAddressByUser 获取用户默认地址
func (c *RedisProfileCache) GetDefaultAddressByUser(ctx context.Context, userID int64) (*entity.Address, error) {
	key := fmt.Sprintf("%s%d", c.defaultAddrPrefix, userID)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	
	var address entity.Address
	if err := json.Unmarshal(data, &address); err != nil {
		return nil, err
	}
	
	return &address, nil
}

// SetDefaultAddressByUser 设置用户默认地址
func (c *RedisProfileCache) SetDefaultAddressByUser(ctx context.Context, userID int64, address *entity.Address) error {
	key := fmt.Sprintf("%s%d", c.defaultAddrPrefix, userID)
	data, err := json.Marshal(address)
	if err != nil {
		return err
	}
	
	return c.client.Set(ctx, key, data, c.expiration).Err()
}

// DeleteDefaultAddressCache 删除默认地址缓存
func (c *RedisProfileCache) DeleteDefaultAddressCache(ctx context.Context, userID int64) error {
	key := fmt.Sprintf("%s%d", c.defaultAddrPrefix, userID)
	return c.client.Del(ctx, key).Err()
}
