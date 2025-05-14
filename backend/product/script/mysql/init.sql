-- 商品表
CREATE TABLE `goods` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `category_id` int(11) NOT NULL COMMENT '分类ID',
  `brands_id` int(11) NOT NULL COMMENT '品牌ID',
  `on_sale` tinyint(1) DEFAULT 1 COMMENT '是否上架',
  `ship_free` tinyint(1) DEFAULT 1 COMMENT '是否免运费',
  `is_new` tinyint(1) DEFAULT 0 COMMENT '是否新品',
  `is_hot` tinyint(1) DEFAULT 0 COMMENT '是否热销',
  `name` varchar(100) NOT NULL COMMENT '商品名称',
  `goods_sn` varchar(50) DEFAULT '' COMMENT '商品编号',
  `click_num` int(11) DEFAULT 0 COMMENT '点击数',
  `sold_num` int(11) DEFAULT 0 COMMENT '销量',
  `fav_num` int(11) DEFAULT 0 COMMENT '收藏数',
  `market_price` float DEFAULT 0 COMMENT '市场价',
  `shop_price` float DEFAULT 0 COMMENT '本店价格',
  `goods_brief` varchar(255) DEFAULT '' COMMENT '商品简短描述',
  `goods_desc` text COMMENT '商品详情',
  `goods_front_image` varchar(255) DEFAULT '' COMMENT '商品封面图',
  `is_deleted` tinyint(1) DEFAULT 0 COMMENT '是否已删除',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_category_id` (`category_id`),
  INDEX `idx_brands_id` (`brands_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 商品SKU表
CREATE TABLE `goods_sku` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `goods` int(11) NOT NULL COMMENT '商品ID',
  `sku_name` varchar(100) NOT NULL COMMENT 'SKU名称',
  `sku_code` varchar(50) NOT NULL COMMENT 'SKU编码',
  `bar_code` varchar(50) DEFAULT '' COMMENT '条形码',
  `price` decimal(10,2) NOT NULL COMMENT '价格',
  `promotion_price` decimal(10,2) DEFAULT 0 COMMENT '促销价',
  `points` int(11) DEFAULT 0 COMMENT '积分',
  `stocks` int(11) NOT NULL COMMENT '库存',
  `image` varchar(255) DEFAULT '' COMMENT '图片',
  `original_stock` int(11) DEFAULT 0 COMMENT '原始库存',
  `spec_values` json DEFAULT NULL COMMENT '规格值JSON',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_goods` (`goods`),
  UNIQUE KEY `idx_sku_code` (`sku_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 分类表
CREATE TABLE `category` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL COMMENT '分类名称',
  `parent_category_id` int(11) DEFAULT 0 COMMENT '父分类ID',
  `level` int(11) DEFAULT 1 COMMENT '分类级别',
  `is_tab` tinyint(1) DEFAULT 0 COMMENT '是否显示在首页tab',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_parent_id` (`parent_category_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 品牌表
CREATE TABLE `brands` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL COMMENT '品牌名称',
  `logo` varchar(255) DEFAULT '' COMMENT '品牌logo',
  `desc` varchar(255) DEFAULT '' COMMENT '品牌描述',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 品牌分类关系表
CREATE TABLE `goods_category_brand` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `category_id` int(11) NOT NULL COMMENT '分类ID',
  `brands_id` int(11) NOT NULL COMMENT '品牌ID',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_category_brand` (`category_id`, `brands_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 轮播图表
CREATE TABLE `banner` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `image` varchar(255) NOT NULL COMMENT '轮播图片地址',
  `url` varchar(255) DEFAULT '' COMMENT '跳转链接',
  `index` int(11) DEFAULT 0 COMMENT '排序索引',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 商品属性表
CREATE TABLE `goods_attribute` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `goods` int(11) NOT NULL COMMENT '商品ID',
  `attr_name` varchar(50) NOT NULL COMMENT '属性名',
  `attr_value` varchar(255) NOT NULL COMMENT '属性值',
  `attr_sort` int(11) DEFAULT 0 COMMENT '排序',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_goods` (`goods`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 商品图片表
CREATE TABLE `goods_image` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `goods` int(11) NOT NULL COMMENT '商品ID',
  `image_url` varchar(255) NOT NULL COMMENT '图片地址',
  `is_main` tinyint(1) DEFAULT 0 COMMENT '是否主图',
  `sort` int(11) DEFAULT 0 COMMENT '排序',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_goods` (`goods`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 商品规格表
CREATE TABLE `goods_spec` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `goods` int(11) NOT NULL COMMENT '商品ID',
  `spec_name` varchar(50) NOT NULL COMMENT '规格名',
  `spec_values` json NOT NULL COMMENT '规格值列表JSON',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_goods` (`goods`),
  UNIQUE KEY `idx_goods_spec` (`goods`, `spec_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 热搜词表
CREATE TABLE `hot_search` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `keyword` varchar(50) NOT NULL COMMENT '关键词',
  `count` int(11) DEFAULT 0 COMMENT '搜索次数',
  `is_active` tinyint(1) DEFAULT 1 COMMENT '是否启用',
  `index` int(11) DEFAULT 0 COMMENT '排序',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_keyword` (`keyword`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 搜索历史表
CREATE TABLE `search_history` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user` int(11) NOT NULL COMMENT '用户ID',
  `keyword` varchar(50) NOT NULL COMMENT '搜索关键词',
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_user_created` (`user`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
