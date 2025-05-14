-- 库存服务数据库初始化脚本

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 创建库存表
DROP TABLE IF EXISTS `inventory`;
CREATE TABLE `inventory` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `goods` bigint(20) NOT NULL COMMENT '商品ID',
  `stocks` int(11) NOT NULL DEFAULT 0 COMMENT '库存数量',
  `version` int(11) NOT NULL DEFAULT 0 COMMENT '乐观锁版本号',
  `warehouse_id` int(11) NOT NULL DEFAULT 1 COMMENT '仓库ID',
  `lock_stocks` int(11) NOT NULL DEFAULT 0 COMMENT '锁定库存数量',
  `alert_threshold` int(11) DEFAULT 10 COMMENT '预警阈值',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_goods_warehouse` (`goods`, `warehouse_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='库存表';

-- 创建库存记录表
DROP TABLE IF EXISTS `stock_sell_detail`;
CREATE TABLE `stock_sell_detail` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `order_sn` varchar(50) NOT NULL COMMENT '订单号',
  `status` int(11) NOT NULL DEFAULT 1 COMMENT '状态：1:锁定，2:已扣减，3:已归还',
  `detail` json DEFAULT NULL COMMENT '库存扣减明细，结构为[{goods_id:1, num:2, warehouse_id:1}]',
  `lock_time` datetime(3) DEFAULT NULL COMMENT '锁定时间',
  `confirm_time` datetime(3) DEFAULT NULL COMMENT '确认时间',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_order_sn` (`order_sn`),
  KEY `idx_status` (`status`),
  KEY `idx_lock_time` (`lock_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='库存操作明细表';

-- 创建仓库表
DROP TABLE IF EXISTS `warehouse`;
CREATE TABLE `warehouse` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '仓库名称',
  `address` varchar(255) NOT NULL COMMENT '仓库地址',
  `contact` varchar(50) DEFAULT NULL COMMENT '联系人',
  `phone` varchar(20) DEFAULT NULL COMMENT '联系电话',
  `status` tinyint(1) DEFAULT 1 COMMENT '状态：1-正常，0-禁用',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='仓库表';

-- 创建库存变更历史表
DROP TABLE IF EXISTS `inventory_history`;
CREATE TABLE `inventory_history` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `goods` bigint(20) NOT NULL COMMENT '商品ID',
  `warehouse_id` int(11) NOT NULL COMMENT '仓库ID',
  `quantity` int(11) NOT NULL COMMENT '变更数量（正数增加，负数减少）',
  `operation_type` varchar(20) NOT NULL COMMENT '操作类型：lock, unlock, decrease, increase, adjust',
  `operator` varchar(50) DEFAULT NULL COMMENT '操作人',
  `order_sn` varchar(50) DEFAULT NULL COMMENT '相关订单号',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_goods` (`goods`),
  KEY `idx_order_sn` (`order_sn`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='库存变更历史表';

-- 初始化默认仓库
INSERT INTO `warehouse` (`name`, `address`, `contact`, `phone`, `status`, `created_at`, `updated_at`)
VALUES ('默认仓库', '默认地址', '系统管理员', '10000000000', 1, NOW(), NOW());

SET FOREIGN_KEY_CHECKS = 1;
