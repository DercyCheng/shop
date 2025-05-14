-- 用户表
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `mobile` varchar(11) NOT NULL COMMENT '手机号码',
  `password` varchar(100) NOT NULL COMMENT '密码',
  `nickname` varchar(20) DEFAULT '' COMMENT '昵称',
  `avatar` varchar(200) DEFAULT '' COMMENT '用户头像URL',
  `birthday` datetime DEFAULT NULL COMMENT '生日',
  `gender` varchar(6) DEFAULT 'male' COMMENT '性别',
  `role` int(11) DEFAULT 1 COMMENT '角色，1表示普通用户，2表示管理员',
  `status` int(11) DEFAULT 1 COMMENT '用户状态：1正常、2禁用、3锁定',
  `login_fail_count` int(11) DEFAULT 0 COMMENT '连续登录失败次数',
  `last_login_at` datetime DEFAULT NULL COMMENT '最后登录时间',
  `wechat_open_id` varchar(100) DEFAULT NULL COMMENT '微信OpenID',
  `wechat_union_id` varchar(100) DEFAULT NULL COMMENT '微信UnionID',
  `session_id` varchar(128) DEFAULT NULL COMMENT '会话ID',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_mobile` (`mobile`),
  KEY `idx_wechat_open_id` (`wechat_open_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 短信验证码表
CREATE TABLE `verification_code` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `mobile` varchar(11) NOT NULL COMMENT '手机号码',
  `code` varchar(6) NOT NULL COMMENT '验证码',
  `type` int(11) NOT NULL COMMENT '类型：1注册，2登录，3重置密码',
  `expire_at` datetime NOT NULL COMMENT '过期时间',
  `used` tinyint(1) DEFAULT 0 COMMENT '是否已使用',
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_mobile_type` (`mobile`, `type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 初始化管理员用户 (密码: admin123)
INSERT INTO `user` (`mobile`, `password`, `nickname`, `role`, `status`, `created_at`, `updated_at`)
VALUES ('13800138000', '$2a$10$1qAz2wSx3eDc4rFv5tGb5edGWh99ybLcDkL8E6LZ.VU8t3Udywbyi', '系统管理员', 2, 1, NOW(), NOW());
