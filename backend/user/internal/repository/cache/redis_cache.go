package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	
	"shop/backend/user/internal/domain/entity"
	
	"github.com/go-redis/redis/v8"
)

// UserCache 用户缓存接口
type UserCache interface {
	// 根据ID获取用户
	GetUser(ctx context.Context, id int64) (*entity.User, error)
	
	// 根据手机号获取用户
	GetUserByMobile(ctx context.Context, mobile string) (*entity.User, error)
	
	// 设置用户缓存
	SetUser(ctx context.Context, user *entity.User) error
	
	// 删除用户缓存
	DeleteUser(ctx context.Context, id int64) error
}

// RedisUserCache 基于Redis的用户缓存实现
type RedisUserCache struct {
	client       *redis.Client
	expiration   time.Duration
	keyPrefix    string
	mobilePrefix string
}

// NewRedisUserCache 创建Redis用户缓存实例
func NewRedisUserCache(client *redis.Client, expiration time.Duration) UserCache {
	return &RedisUserCache{
		client:       client,
		expiration:   expiration,
		keyPrefix:    "user:id:",
		mobilePrefix: "user:mobile:",
	}
}

// GetUser 根据ID从缓存获取用户
func (c *RedisUserCache) GetUser(ctx context.Context, id int64) (*entity.User, error) {
	key := fmt.Sprintf("%s%d", c.keyPrefix, id)
	
	// 从Redis获取用户JSON数据
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存未命中
		}
		return nil, err
	}
	
	// 将JSON数据反序列化为用户对象
	var user entity.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}
	
	return &user, nil
}

// GetUserByMobile 根据手机号从缓存获取用户
func (c *RedisUserCache) GetUserByMobile(ctx context.Context, mobile string) (*entity.User, error) {
	key := fmt.Sprintf("%s%s", c.mobilePrefix, mobile)
	
	// 从索引中获取用户ID
	idStr, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	
	// 从用户ID获取用户对象
	var id int64
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		return nil, err
	}
	
	return c.GetUser(ctx, id)
}

// SetUser 设置用户缓存
func (c *RedisUserCache) SetUser(ctx context.Context, user *entity.User) error {
	// 序列化用户对象为JSON
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	
	// 设置用户ID到对象的映射
	key := fmt.Sprintf("%s%d", c.keyPrefix, user.ID)
	if err := c.client.Set(ctx, key, data, c.expiration).Err(); err != nil {
		return err
	}
	
	// 设置手机号到用户ID的映射（二级索引）
	if user.Mobile != "" {
		mobileKey := fmt.Sprintf("%s%s", c.mobilePrefix, user.Mobile)
		if err := c.client.Set(ctx, mobileKey, user.ID, c.expiration).Err(); err != nil {
			return err
		}
	}
	
	return nil
}

// DeleteUser 删除用户缓存
func (c *RedisUserCache) DeleteUser(ctx context.Context, id int64) error {
	// 先获取用户信息，以便删除手机号索引
	user, err := c.GetUser(ctx, id)
	if err == nil && user != nil && user.Mobile != "" {
		mobileKey := fmt.Sprintf("%s%s", c.mobilePrefix, user.Mobile)
		if err := c.client.Del(ctx, mobileKey).Err(); err != nil {
			return err
		}
	}
	
	// 删除ID索引
	key := fmt.Sprintf("%s%d", c.keyPrefix, id)
	return c.client.Del(ctx, key).Err()
}
