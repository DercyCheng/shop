-- Profile Service Database Initialization

-- User Favorites Table
CREATE TABLE IF NOT EXISTS `user_fav` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user` int(11) NOT NULL COMMENT '用户ID',
  `goods` int(11) NOT NULL COMMENT '商品ID',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_goods` (`user`, `goods`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- User Address Table
CREATE TABLE IF NOT EXISTS `address` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user` int(11) NOT NULL COMMENT '用户ID',
  `province` varchar(20) NOT NULL COMMENT '省',
  `city` varchar(20) NOT NULL COMMENT '市',
  `district` varchar(20) NOT NULL COMMENT '区/县',
  `address` varchar(100) NOT NULL COMMENT '详细地址',
  `signer_name` varchar(50) NOT NULL COMMENT '收货人姓名',
  `signer_mobile` varchar(11) NOT NULL COMMENT '收货人手机号',
  `is_default` tinyint(1) DEFAULT 0 COMMENT '是否默认地址',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_user` (`user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- User Messages Table
CREATE TABLE IF NOT EXISTS `leaving_messages` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user` int(11) NOT NULL COMMENT '用户ID',
  `message_type` int(11) DEFAULT 1 COMMENT '留言类型：1-留言，2-投诉，3-询问，4-售后，5-求购',
  `subject` varchar(100) NOT NULL COMMENT '主题',
  `message` text COMMENT '留言内容',
  `file` varchar(255) DEFAULT '' COMMENT '上传文件',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_user` (`user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;