package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	
	"shop/backend/inventory/internal/domain/entity"
)

// InventoryCache 库存缓存接口
type InventoryCache interface {
	// 获取商品库存
	GetInventory(ctx context.Context, productID int64, warehouseID int) (*entity.Inventory, error)
	// 批量获取库存，返回找到的库存和缺失的ID
	BatchGetInventory(ctx context.Context, productIDs []int64, warehouseID int) ([]*entity.Inventory, []int64)
	// 设置库存缓存
	SetInventory(ctx context.Context, inventory *entity.Inventory) error
	// 删除库存缓存
	DeleteInventory(ctx context.Context, productID int64, warehouseID int) error
}

const (
	// 库存缓存键前缀
	inventoryKeyPrefix = "inventory:"
	// 默认缓存过期时间
	defaultCacheTTL = 5 * time.Minute
)

// RedisInventoryCache Redis库存缓存实现
type RedisInventoryCache struct {
	client  *redis.Client
	logger  *zap.Logger
	cacheTTL time.Duration
}

// NewRedisInventoryCache 创建Redis库存缓存
func NewRedisInventoryCache(client *redis.Client, logger *zap.Logger, ttlSeconds int) InventoryCache {
	ttl := defaultCacheTTL
	if ttlSeconds > 0 {
		ttl = time.Duration(ttlSeconds) * time.Second
	}
	
	return &RedisInventoryCache{
		client:  client,
		logger:  logger,
		cacheTTL: ttl,
	}
}

// 构建库存缓存键
func buildInventoryKey(productID int64, warehouseID int) string {
	return fmt.Sprintf("%s%d:%d", inventoryKeyPrefix, productID, warehouseID)
}

// GetInventory 从Redis获取库存
func (c *RedisInventoryCache) GetInventory(ctx context.Context, productID int64, warehouseID int) (*entity.Inventory, error) {
	key := buildInventoryKey(productID, warehouseID)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			// 缓存未命中
			return nil, err
		}
		c.logger.Warn("Failed to get inventory from Redis",
			zap.Int64("product_id", productID),
			zap.Int("warehouse_id", warehouseID),
			zap.Error(err))
		return nil, err
	}
	
	var inventory entity.Inventory
	if err := json.Unmarshal(data, &inventory); err != nil {
		c.logger.Error("Failed to unmarshal inventory data",
			zap.Int64("product_id", productID),
			zap.Int("warehouse_id", warehouseID),
			zap.Error(err))
		return nil, err
	}
	
	return &inventory, nil
}

// BatchGetInventory 批量获取库存
func (c *RedisInventoryCache) BatchGetInventory(ctx context.Context, productIDs []int64, warehouseID int) ([]*entity.Inventory, []int64) {
	if len(productIDs) == 0 {
		return []*entity.Inventory{}, []int64{}
	}
	
	// 批量构建缓存键
	keys := make([]string, 0, len(productIDs))
	for _, productID := range productIDs {
		keys = append(keys, buildInventoryKey(productID, warehouseID))
	}
	
	// 批量获取缓存数据
	pipe := c.client.Pipeline()
	cmds := make([]*redis.StringCmd, len(keys))
	for i, key := range keys {
		cmds[i] = pipe.Get(ctx, key)
	}
	
	// 执行批量命令
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		c.logger.Error("Failed to exec pipeline", zap.Error(err))
	}
	
	// 处理结果
	inventories := make([]*entity.Inventory, 0, len(productIDs))
	missingIDs := make([]int64, 0)
	
	for i, cmd := range cmds {
		data, err := cmd.Bytes()
		if err != nil {
			missingIDs = append(missingIDs, productIDs[i])
			continue
		}
		
		var inventory entity.Inventory
		if err := json.Unmarshal(data, &inventory); err != nil {
			c.logger.Error("Failed to unmarshal inventory data",
				zap.Int64("product_id", productIDs[i]),
				zap.Error(err))
			missingIDs = append(missingIDs, productIDs[i])
			continue
		}
		
		inventories = append(inventories, &inventory)
	}
	
	return inventories, missingIDs
}

// SetInventory 设置库存缓存
func (c *RedisInventoryCache) SetInventory(ctx context.Context, inventory *entity.Inventory) error {
	if inventory == nil {
		return fmt.Errorf("inventory cannot be nil")
	}
	
	data, err := json.Marshal(inventory)
	if err != nil {
		c.logger.Error("Failed to marshal inventory data",
			zap.Int64("product_id", inventory.ProductID),
			zap.Int("warehouse_id", inventory.WarehouseID),
			zap.Error(err))
		return err
	}
	
	key := buildInventoryKey(inventory.ProductID, inventory.WarehouseID)
	err = c.client.Set(ctx, key, data, c.cacheTTL).Err()
	if err != nil {
		c.logger.Error("Failed to set inventory cache",
			zap.Int64("product_id", inventory.ProductID),
			zap.Int("warehouse_id", inventory.WarehouseID),
			zap.Error(err))
		return err
	}
	
	return nil
}

// DeleteInventory 删除库存缓存
func (c *RedisInventoryCache) DeleteInventory(ctx context.Context, productID int64, warehouseID int) error {
	key := buildInventoryKey(productID, warehouseID)
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		c.logger.Warn("Failed to delete inventory cache",
			zap.Int64("product_id", productID),
			zap.Int("warehouse_id", warehouseID),
			zap.Error(err))
		return err
	}
	
	return nil
}
