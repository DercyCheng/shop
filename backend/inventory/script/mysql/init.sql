-- shop_inventory schema

-- Create the database if it doesn't exist
CREATE DATABASE IF NOT EXISTS shop_inventory;

-- Use the database
USE shop_inventory;

-- Create warehouses table
CREATE TABLE IF NOT EXISTS `warehouses` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(100) NOT NULL,
  `address` VARCHAR(255),
  `is_active` BOOLEAN NOT NULL DEFAULT TRUE,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_warehouse_name` (`name`)
);

-- Create stocks table
CREATE TABLE IF NOT EXISTS `stocks` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `product_id` BIGINT NOT NULL,
  `warehouse_id` BIGINT NOT NULL,
  `quantity` INT NOT NULL DEFAULT 0,
  `reserved` INT NOT NULL DEFAULT 0,
  `low_stock_threshold` INT NOT NULL DEFAULT 5,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_product_warehouse` (`product_id`, `warehouse_id`),
  INDEX `idx_warehouse_id` (`warehouse_id`),
  CONSTRAINT `fk_stock_warehouse` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouses` (`id`)
);

-- Create reservations table
CREATE TABLE IF NOT EXISTS `reservations` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `order_id` VARCHAR(50) NOT NULL,
  `status` ENUM('PENDING', 'COMMITTED', 'CANCELLED', 'EXPIRED') NOT NULL DEFAULT 'PENDING',
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `expires_at` TIMESTAMP NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_reservation_order_id` (`order_id`),
  INDEX `idx_reservation_status` (`status`),
  INDEX `idx_reservation_expires_at` (`expires_at`)
);

-- Create reservation items table
CREATE TABLE IF NOT EXISTS `reservation_items` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `reservation_id` BIGINT NOT NULL,
  `product_id` BIGINT NOT NULL,
  `warehouse_id` BIGINT NOT NULL,
  `quantity` INT NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_reservation_item_reservation_id` (`reservation_id`),
  INDEX `idx_reservation_item_product_id` (`product_id`),
  INDEX `idx_reservation_item_warehouse_id` (`warehouse_id`),
  CONSTRAINT `fk_reservation_item_reservation` FOREIGN KEY (`reservation_id`) REFERENCES `reservations` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_reservation_item_warehouse` FOREIGN KEY (`warehouse_id`) REFERENCES `warehouses` (`id`)
);

-- Insert some initial warehouse data if the table is empty
INSERT INTO `warehouses` (`name`, `address`, `is_active`)
SELECT 'Main Warehouse', '123 Main Street, City', TRUE
WHERE NOT EXISTS (SELECT 1 FROM `warehouses` LIMIT 1);

INSERT INTO `warehouses` (`name`, `address`, `is_active`)
SELECT 'Secondary Warehouse', '456 State Street, Other City', TRUE
WHERE EXISTS (SELECT 1 FROM `warehouses` LIMIT 1) AND COUNT(*) < 2 FROM `warehouses`;