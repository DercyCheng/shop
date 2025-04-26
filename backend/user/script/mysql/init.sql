-- 创建用户表
CREATE TABLE IF NOT EXISTS `users` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(50) NOT NULL,
    `password` VARCHAR(100) NOT NULL,
    `nickname` VARCHAR(50) DEFAULT NULL,
    `email` VARCHAR(100) DEFAULT NULL,
    `phone` VARCHAR(20) DEFAULT NULL,
    `avatar` VARCHAR(255) DEFAULT NULL,
    `gender` VARCHAR(10) DEFAULT 'unknown',
    `birthday` DATETIME DEFAULT NULL,
    `status` TINYINT DEFAULT 1 COMMENT '1: 正常, 2: 禁用, 3: 锁定',
    `role` TINYINT DEFAULT 1 COMMENT '1: 普通用户, 2: 管理员',
    `login_fail_count` INT DEFAULT 0,
    `last_login_at` DATETIME DEFAULT NULL,
    `wechat_open_id` VARCHAR(50) DEFAULT NULL,
    `wechat_union_id` VARCHAR(50) DEFAULT NULL,
    `session_id` VARCHAR(100) DEFAULT NULL,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    `deleted_at` DATETIME DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_username` (`username`),
    UNIQUE KEY `idx_phone` (`phone`),
    KEY `idx_email` (`email`),
    KEY `idx_wechat_open_id` (`wechat_open_id`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 创建初始管理员用户（密码为 admin123：使用MD5哈希 + 随机盐值）
-- 注意：实际生产中应该使用更强的哈希算法，这里仅作为示例
INSERT INTO `users` (
    `username`, 
    `password`, 
    `nickname`, 
    `email`, 
    `status`, 
    `role`, 
    `created_at`, 
    `updated_at`
) VALUES (
    'admin', 
    '5f4dcc3b5aa765d61d8327deb882cf99:randomsalt', -- 实际密码会是真正哈希后的值
    '系统管理员', 
    'admin@example.com', 
    1, 
    2, 
    NOW(), 
    NOW()
) ON DUPLICATE KEY UPDATE `updated_at` = NOW();