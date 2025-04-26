package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"shop/backend/user/internal/domain/entity"
)

// RedisUserCache Redis用户缓存实现
type RedisUserCache struct {
	client     *redis.Client
	keyPrefix  string
	defaultTTL time.Duration
}

// NewRedisUserCache 创建Redis用户缓存
func NewRedisUserCache(client *redis.Client, keyPrefix string, defaultTTL time.Duration) *RedisUserCache {
	return &RedisUserCache{
		client:     client,
		keyPrefix:  keyPrefix,
		defaultTTL: defaultTTL,
	}
}

// getUserKey 生成用户信息缓存键
func (c *RedisUserCache) getUserKey(userID int64) string {
	return fmt.Sprintf("%s:user:%d", c.keyPrefix, userID)
}

// getTokenKey 生成令牌缓存键
func (c *RedisUserCache) getTokenKey(token string) string {
	return fmt.Sprintf("%s:token:%s", c.keyPrefix, token)
}

// getRefreshTokenKey 生成刷新令牌缓存键
func (c *RedisUserCache) getRefreshTokenKey(refreshToken string) string {
	return fmt.Sprintf("%s:refresh:%s", c.keyPrefix, refreshToken)
}

// GetUser 从缓存获取用户信息
func (c *RedisUserCache) GetUser(ctx context.Context, userID int64) (*entity.User, error) {
	key := c.getUserKey(userID)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var user entity.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// SetUser 缓存用户信息
func (c *RedisUserCache) SetUser(ctx context.Context, user *entity.User, ttl time.Duration) error {
	if ttl == 0 {
		ttl = c.defaultTTL
	}

	key := c.getUserKey(user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, ttl).Err()
}

// DeleteUser 删除缓存的用户信息
func (c *RedisUserCache) DeleteUser(ctx context.Context, userID int64) error {
	key := c.getUserKey(userID)
	return c.client.Del(ctx, key).Err()
}

// SaveToken 保存访问令牌与用户ID的映射
func (c *RedisUserCache) SaveToken(ctx context.Context, token string, userID int64, ttl time.Duration) error {
	key := c.getTokenKey(token)
	return c.client.Set(ctx, key, userID, ttl).Err()
}

// GetUserIDByToken 通过访问令牌获取用户ID
func (c *RedisUserCache) GetUserIDByToken(ctx context.Context, token string) (int64, error) {
	key := c.getTokenKey(token)
	userID, err := c.client.Get(ctx, key).Int64()
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// InvalidateToken 使访问令牌失效
func (c *RedisUserCache) InvalidateToken(ctx context.Context, token string) error {
	key := c.getTokenKey(token)
	return c.client.Del(ctx, key).Err()
}

// SaveRefreshToken 保存刷新令牌与用户ID的映射
func (c *RedisUserCache) SaveRefreshToken(ctx context.Context, refreshToken string, userID int64, ttl time.Duration) error {
	key := c.getRefreshTokenKey(refreshToken)
	return c.client.Set(ctx, key, userID, ttl).Err()
}

// GetUserIDByRefreshToken 通过刷新令牌获取用户ID
func (c *RedisUserCache) GetUserIDByRefreshToken(ctx context.Context, refreshToken string) (int64, error) {
	key := c.getRefreshTokenKey(refreshToken)
	userID, err := c.client.Get(ctx, key).Int64()
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// InvalidateRefreshToken 使刷新令牌失效
func (c *RedisUserCache) InvalidateRefreshToken(ctx context.Context, refreshToken string) error {
	key := c.getRefreshTokenKey(refreshToken)
	return c.client.Del(ctx, key).Err()
}

// InvalidateAllUserTokens 使用户的所有令牌失效
// 注意：这是一个昂贵的操作，需要通过模式匹配查找所有相关的键
func (c *RedisUserCache) InvalidateAllUserTokens(ctx context.Context, userID int64) error {
	// 查找与该用户关联的所有令牌（实际实现可能需要不同的策略）
	// 这里假设我们有一个单独的集合来跟踪用户的所有令牌
	userTokensKey := fmt.Sprintf("%s:user_tokens:%d", c.keyPrefix, userID)

	// 从集合中获取所有令牌
	tokens, err := c.client.SMembers(ctx, userTokensKey).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	// 删除每个令牌
	for _, token := range tokens {
		tokenKey := c.getTokenKey(token)
		if err := c.client.Del(ctx, tokenKey).Err(); err != nil {
			return err
		}
	}

	// 清空集合
	return c.client.Del(ctx, userTokensKey).Err()
}
