
-- Create databases
CREATE DATABASE IF NOT EXISTS mxshop_user_srv DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
CREATE DATABASE IF NOT EXISTS mxshop_goods_srv DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
CREATE DATABASE IF NOT EXISTS mxshop_inventory_srv DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
CREATE DATABASE IF NOT EXISTS mxshop_order_srv DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
CREATE DATABASE IF NOT EXISTS mxshop_userop_srv DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

-- User service tables
USE mxshop_user_srv;

CREATE TABLE IF NOT EXISTS `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `mobile` varchar(11) NOT NULL,
  `password` varchar(100) NOT NULL,
  `nickname` varchar(20) DEFAULT '',
  `birthday` datetime DEFAULT NULL,
  `gender` varchar(6) DEFAULT