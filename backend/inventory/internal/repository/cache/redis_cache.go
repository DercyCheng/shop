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

const (
	// 库存缓存前缀
	inventoryCacheKeyPrefix = "inventory:item:"
	// 过期时间
	inventoryCacheTTL = 24 * time.Hour
)

// InventoryCache 库存缓存接口
type InventoryCache interface {
	GetInventory(ctx context.Context, productID int64, warehouseID int) (*entity.Inventory, error)
	SetInventory(ctx context.Context, inventory *entity.Inventory) error
	BatchGetInventory(ctx context.Context, productIDs []int64, warehouseID int) ([]*entity.Inventory, []int64) // 返回结果和缓存未命中的ID
	DeleteInventory(ctx context.Context, productID int64, warehouseID int) error
}

// RedisInventoryCache Redis实现的库存缓存
type RedisInventoryCache struct {
	client *redis.Client
	logger *zap.Logger
}

// NewRedisInventoryCache 创建Redis库存缓存
func NewRedisInventoryCache(client *redis.Client, logger *zap.Logger) InventoryCache {
	return &RedisInventoryCache{
		client: client,
		logger: logger,
	}
}

// getCacheKey 生成缓存键
func (c *RedisInventoryCache) getCacheKey(productID int64, warehouseID int) string {
	return fmt.Sprintf("%s%d:%d", inventoryCacheKeyPrefix, productID, warehouseID)
}

// GetInventory 从缓存获取库存
func (c *RedisInventoryCache) GetInventory(ctx context.Context, productID int64, warehouseID int) (*entity.Inventory, error) {
	key := c.getCacheKey(productID, warehouseID)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("inventory cache miss")
		}
		c.logger.Error("Failed to get inventory from redis", 
			zap.Error(err), 
			zap.Int64("product_id", productID),
			zap.Int("warehouse_id", warehouseID),
			zap.String("key", key))
		return nil, err
	}
	
	var inventory entity.Inventory
	if err := json.Unmarshal(data, &inventory); err != nil {
		c.logger.Error("Failed to unmarshal inventory from redis", 
			zap.Error(err), 
			zap.Int64("product_id", productID),
			zap.Int("warehouse_id", warehouseID))
		return nil, err
	}
	
	return &inventory, nil
}

// SetInventory 缓存库存信息
func (c *RedisInventoryCache) SetInventory(ctx context.Context, inventory *entity.Inventory) error {
	key := c.getCacheKey(inventory.ProductID, inventory.WarehouseID)
	data, err := json.Marshal(inventory)
	if err != nil {
		c.logger.Error("Failed to marshal inventory for redis", 
			zap.Error(err), 
			zap.Int64("product_id", inventory.ProductID),
			zap.Int("warehouse_id", inventory.WarehouseID))
		return err
	}
	
	if err := c.client.Set(ctx, key, data, inventoryCacheTTL).Err(); err != nil {
		c.logger.Error("Failed to set inventory to redis", 
			zap.Error(err), 
			zap.Int64("product_id", inventory.ProductID),
			zap.Int("warehouse_id", inventory.WarehouseID),
			zap.String("key", key))
		return err
	}
	
	return nil
}

// BatchGetInventory 批量获取库存信息
func (c *RedisInventoryCache) BatchGetInventory(ctx context.Context, productIDs []int64, warehouseID int) ([]*entity.Inventory, []int64) {
	var inventories []*entity.Inventory
	var missingIDs []int64
	
	// 创建管道批量操作
	pipe := c.client.Pipeline()
	commands := make(map[int64]*redis.StringCmd)
	
	// 构建批量查询
	for _, productID := range productIDs {
		key := c.getCacheKey(productID, warehouseID)
		commands[productID] = pipe.Get(ctx, key)
	}
	
	// 执行管道
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		c.logger.Error("Failed to execute redis pipeline for batch get inventory", 
			zap.Error(err))
		// 如果Redis失败，所有ID都需要从数据库查询
		return []*entity.Inventory{}, productIDs
	}
	
	// 处理结果
	for productID, cmd := range commands {
		data, err := cmd.Bytes()
		if err != nil {
			missingIDs = append(missingIDs, productID)
			continue
		}
		
		var inventory entity.Inventory
		if err := json.Unmarshal(data, &inventory); err != nil {
			c.logger.Error("Failed to unmarshal inventory from redis", 
				zap.Error(err), 
				zap.Int64("product_id", productID),
				zap.Int("warehouse_id", warehouseID))
			missingIDs = append(missingIDs, productID)
			continue
		}
		
		inventories = append(inventories, &inventory)
	}
	
	return inventories, missingIDs
}

// DeleteInventory 删除库存缓存
func (c *RedisInventoryCache) DeleteInventory(ctx context.Context, productID int64, warehouseID int) error {
	key := c.getCacheKey(productID, warehouseID)
	if err := c.client.Del(ctx, key).Err(); err != nil {
		c.logger.Error("Failed to delete inventory from redis", 
			zap.Error(err), 
			zap.Int64("product_id", productID),
			zap.Int("warehouse_id", warehouseID),
			zap.String("key", key))
		return err
	}
	
	return nil
}
