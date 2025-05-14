-- 用户收藏表
CREATE TABLE `user_fav` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user` int(11) NOT NULL COMMENT '用户ID',
  `goods` int(11) NOT NULL COMMENT '商品ID',
  `category_id` int(11) DEFAULT NULL COMMENT '商品分类ID',
  `remark` varchar(255) DEFAULT NULL COMMENT '收藏备注',
  `price_when_fav` decimal(10,2) DEFAULT NULL COMMENT '收藏时价格',
  `notification` tinyint(1) DEFAULT 0 COMMENT '价格变动通知',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_goods` (`user`, `goods`),
  INDEX `idx_user_category` (`user`, `category_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 用户地址表
CREATE TABLE `address` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user` int(11) NOT NULL COMMENT '用户ID',
  `province` varchar(20) NOT NULL COMMENT '省',
  `city` varchar(20) NOT NULL COMMENT '市',
  `district` varchar(20) NOT NULL COMMENT '区/县',
  `address` varchar(100) NOT NULL COMMENT '详细地址',
  `signer_name` varchar(50) NOT NULL COMMENT '收货人姓名',
  `signer_mobile` varchar(11) NOT NULL COMMENT '收货人手机号',
  `is_default` tinyint(1) DEFAULT 0 COMMENT '是否默认地址',
  `label` varchar(20) DEFAULT NULL COMMENT '地址标签：家、公司等',
  `postcode` varchar(10) DEFAULT NULL COMMENT '邮政编码',
  `usage_count` int(11) DEFAULT 0 COMMENT '使用次数',
  `last_used_at` datetime DEFAULT NULL COMMENT '最后使用时间',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_user` (`user`),
  INDEX `idx_user_default` (`user`, `is_default`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 用户反馈表
CREATE TABLE `user_feedback` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user` int(11) NOT NULL COMMENT '用户ID',
  `feedback_type` int(11) DEFAULT 1 COMMENT '反馈类型：1-留言，2-投诉，3-询问，4-售后，5-求购',
  `subject` varchar(100) NOT NULL COMMENT '主题',
  `content` text COMMENT '反馈内容',
  `file_urls` json DEFAULT NULL COMMENT '附件URLs',
  `status` tinyint(1) DEFAULT 0 COMMENT '状态：0-待处理，1-处理中，2-已解决，3-已关闭',
  `order_sn` varchar(50) DEFAULT NULL COMMENT '相关订单号',
  `admin_reply` text DEFAULT NULL COMMENT '管理员回复',
  `reply_at` datetime DEFAULT NULL COMMENT '回复时间',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_user` (`user`),
  INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 浏览历史表
CREATE TABLE `browsing_history` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user` int(11) NOT NULL COMMENT '用户ID',
  `goods` int(11) NOT NULL COMMENT '商品ID',
  `source` varchar(50) DEFAULT NULL COMMENT '来源',
  `stay_time` int(11) DEFAULT 0 COMMENT '停留时间(秒)',
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_user_created` (`user`, `created_at`),
  INDEX `idx_goods` (`goods`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 用户偏好设置表
CREATE TABLE `user_setting` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user` int(11) NOT NULL COMMENT '用户ID',
  `notify_new_order` tinyint(1) DEFAULT 1 COMMENT '新订单通知',
  `notify_promotion` tinyint(1) DEFAULT 1 COMMENT '促销活动通知',
  `notify_system` tinyint(1) DEFAULT 1 COMMENT '系统通知',
  `privacy_show_fav` tinyint(1) DEFAULT 1 COMMENT '隐私-展示收藏',
  `theme_color` varchar(20) DEFAULT 'default' COMMENT '主题颜色',
  `language` varchar(10) DEFAULT 'zh_CN' COMMENT '语言偏好',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user` (`user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
