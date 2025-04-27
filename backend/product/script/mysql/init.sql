-- Create database if not exists
CREATE DATABASE IF NOT EXISTS shop_product;
USE shop_product;

-- Create tables (if not already created by ORM)
-- For safety, we'll add IF NOT EXISTS clauses

-- Categories table
CREATE TABLE IF NOT EXISTS categories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    parent_category_id BIGINT,
    level INT DEFAULT 1,
    is_tab BOOLEAN DEFAULT FALSE,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME
);

-- Brands table
CREATE TABLE IF NOT EXISTS brands (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    logo VARCHAR(255),
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME
);

-- Category-Brand relation table
CREATE TABLE IF NOT EXISTS category_brands (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    category_id BIGINT NOT NULL,
    brands_id BIGINT NOT NULL,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    UNIQUE KEY unique_category_brand (category_id, brands_id)
);

-- Products table
CREATE TABLE IF NOT EXISTS products (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    goods_sn VARCHAR(50),
    category_id BIGINT NOT NULL,
    brands_id BIGINT NOT NULL,
    on_sale BOOLEAN DEFAULT TRUE,
    ship_free BOOLEAN DEFAULT FALSE,
    is_new BOOLEAN DEFAULT FALSE,
    is_hot BOOLEAN DEFAULT FALSE,
    click_num INT DEFAULT 0,
    sold_num INT DEFAULT 0,
    fav_num INT DEFAULT 0,
    market_price FLOAT DEFAULT 0,
    shop_price FLOAT DEFAULT 0,
    goods_brief TEXT,
    goods_desc TEXT,
    goods_front_image VARCHAR(255),
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME
);

-- Product images table
CREATE TABLE IF NOT EXISTS product_images (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    goods_id BIGINT NOT NULL,
    image VARCHAR(255) NOT NULL,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME
);

-- Banners table
CREATE TABLE IF NOT EXISTS banners (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    image VARCHAR(255) NOT NULL,
    url VARCHAR(255),
    `index` INT DEFAULT 0,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME
);

-- Insert sample data

-- Insert categories
INSERT INTO categories (name, parent_category_id, level, is_tab, created_at, updated_at)
VALUES
    ('Electronics', 0, 1, TRUE, NOW(), NOW()),
    ('Clothing', 0, 1, TRUE, NOW(), NOW()),
    ('Home & Garden', 0, 1, TRUE, NOW(), NOW()),
    ('Smartphones', 1, 2, FALSE, NOW(), NOW()),
    ('Laptops', 1, 2, FALSE, NOW(), NOW()),
    ('Accessories', 1, 2, FALSE, NOW(), NOW()),
    ('Men\'s Clothing', 2, 2, FALSE, NOW(), NOW()),
    ('Women\'s Clothing', 2, 2, FALSE, NOW(), NOW()),
    ('Furniture', 3, 2, FALSE, NOW(), NOW()),
    ('Decor', 3, 2, FALSE, NOW(), NOW());

-- Insert brands
INSERT INTO brands (name, logo, created_at, updated_at)
VALUES
    ('Apple', 'https://example.com/apple-logo.png', NOW(), NOW()),
    ('Samsung', 'https://example.com/samsung-logo.png', NOW(), NOW()),
    ('Dell', 'https://example.com/dell-logo.png', NOW(), NOW()),
    ('Nike', 'https://example.com/nike-logo.png', NOW(), NOW()),
    ('Adidas', 'https://example.com/adidas-logo.png', NOW(), NOW()),
    ('IKEA', 'https://example.com/ikea-logo.png', NOW(), NOW()),
    ('Xiaomi', 'https://example.com/xiaomi-logo.png', NOW(), NOW()),
    ('HP', 'https://example.com/hp-logo.png', NOW(), NOW());

-- Insert category-brand relations
INSERT INTO category_brands (category_id, brands_id, created_at, updated_at)
VALUES
    (1, 1, NOW(), NOW()),
    (1, 2, NOW(), NOW()),
    (1, 3, NOW(), NOW()),
    (1, 7, NOW(), NOW()),
    (1, 8, NOW(), NOW()),
    (2, 4, NOW(), NOW()),
    (2, 5, NOW(), NOW()),
    (3, 6, NOW(), NOW()),
    (4, 1, NOW(), NOW()),
    (4, 2, NOW(), NOW()),
    (4, 7, NOW(), NOW()),
    (5, 1, NOW(), NOW()),
    (5, 3, NOW(), NOW()),
    (5, 8, NOW(), NOW()),
    (6, 1, NOW(), NOW()),
    (6, 2, NOW(), NOW()),
    (6, 7, NOW(), NOW()),
    (7, 4, NOW(), NOW()),
    (7, 5, NOW(), NOW()),
    (8, 4, NOW(), NOW()),
    (8, 5, NOW(), NOW()),
    (9, 6, NOW(), NOW()),
    (10, 6, NOW(), NOW());

-- Insert products
INSERT INTO products (name, goods_sn, category_id, brands_id, on_sale, ship_free, is_new, is_hot, 
                    market_price, shop_price, goods_brief, goods_desc, goods_front_image, 
                    created_at, updated_at)
VALUES
    ('iPhone 14 Pro', 'IP14PRO001', 4, 1, TRUE, TRUE, TRUE, TRUE, 
     1199.99, 1099.99, 'Latest iPhone with advanced features', 
     'The iPhone 14 Pro features a 6.1-inch Super Retina XDR display, A16 Bionic chip, and professional camera system.', 
     'https://example.com/iphone14pro.jpg', NOW(), NOW()),

    ('Samsung Galaxy S23', 'SGS23001', 4, 2, TRUE, TRUE, TRUE, FALSE, 
     999.99, 949.99, 'Flagship Samsung smartphone with advanced camera', 
     'The Samsung Galaxy S23 comes with a powerful processor, vibrant display, and a versatile camera system.', 
     'https://example.com/galaxys23.jpg', NOW(), NOW()),

    ('Xiaomi Redmi Note 12', 'XM12001', 4, 7, TRUE, FALSE, TRUE, FALSE, 
     299.99, 249.99, 'Affordable smartphone with great features', 
     'The Xiaomi Redmi Note 12 offers excellent value with a large display, good camera, and long battery life.', 
     'https://example.com/redminote12.jpg', NOW(), NOW()),

    ('MacBook Pro 16"', 'MBP16001', 5, 1, TRUE, TRUE, TRUE, TRUE, 
     2499.99, 2399.99, 'Powerful laptop for professionals', 
     'The MacBook Pro 16" features Apple\'s M2 Pro chip, stunning Liquid Retina XDR display, and all-day battery life.', 
     'https://example.com/macbookpro16.jpg', NOW(), NOW()),

    ('Dell XPS 15', 'DX15001', 5, 3, TRUE, TRUE, FALSE, TRUE, 
     1899.99, 1799.99, 'Premium Windows laptop', 
     'The Dell XPS 15 features a beautiful InfinityEdge display, powerful Intel processors, and premium build quality.', 
     'https://example.com/dellxps15.jpg', NOW(), NOW()),

    ('HP Spectre x360', 'HPSX360001', 5, 8, TRUE, TRUE, FALSE, FALSE, 
     1499.99, 1399.99, 'Versatile convertible laptop', 
     'The HP Spectre x360 is a premium 2-in-1 laptop with a 360-degree hinge, allowing it to be used as a tablet.', 
     'https://example.com/spectrehp.jpg', NOW(), NOW()),

    ('Apple AirPods Pro', 'APP001', 6, 1, TRUE, FALSE, FALSE, TRUE, 
     249.99, 219.99, 'Premium wireless earbuds with noise cancellation', 
     'AirPods Pro feature Active Noise Cancellation, Transparency mode, and a customizable fit.', 
     'https://example.com/airpodspro.jpg', NOW(), NOW()),

    ('Samsung Galaxy Watch 5', 'SGW5001', 6, 2, TRUE, FALSE, TRUE, FALSE, 
     299.99, 279.99, 'Advanced smartwatch with health features', 
     'The Galaxy Watch 5 offers comprehensive health monitoring, fitness tracking, and smart notifications.', 
     'https://example.com/galaxywatch5.jpg', NOW(), NOW()),

    ('Nike Dri-FIT T-Shirt', 'NDF001', 7, 4, TRUE, TRUE, FALSE, FALSE, 
     39.99, 29.99, 'Men\'s performance t-shirt', 
     'Nike Dri-FIT technology helps you stay dry and comfortable during workouts.', 
     'https://example.com/nikedrifit.jpg', NOW(), NOW()),

    ('Adidas Ultraboost', 'AUB001', 7, 5, TRUE, FALSE, TRUE, TRUE, 
     189.99, 169.99, 'Premium running shoes for men', 
     'Ultraboost features responsive cushioning and a supportive fit for maximum comfort.', 
     'https://example.com/adidasultraboost.jpg', NOW(), NOW()),

    ('IKEA MALM Bed Frame', 'IKMALM001', 9, 6, TRUE, FALSE, FALSE, FALSE, 
     249.99, 229.99, 'Simple, stylish bed frame', 
     'The MALM bed frame has a timeless design that will look great in any bedroom.', 
     'https://example.com/ikeamalm.jpg', NOW(), NOW()),

    ('IKEA BILLY Bookcase', 'IKBILLY001', 9, 6, TRUE, FALSE, FALSE, TRUE, 
     99.99, 89.99, 'Versatile bookcase', 
     'The BILLY bookcase is a classic piece that provides ample storage for books and decorative items.', 
     'https://example.com/ikeabilly.jpg', NOW(), NOW());

-- Insert product images
INSERT INTO product_images (goods_id, image, created_at, updated_at)
VALUES
    (1, 'https://example.com/iphone14pro-1.jpg', NOW(), NOW()),
    (1, 'https://example.com/iphone14pro-2.jpg', NOW(), NOW()),
    (1, 'https://example.com/iphone14pro-3.jpg', NOW(), NOW()),
    (2, 'https://example.com/galaxys23-1.jpg', NOW(), NOW()),
    (2, 'https://example.com/galaxys23-2.jpg', NOW(), NOW()),
    (3, 'https://example.com/redminote12-1.jpg', NOW(), NOW()),
    (3, 'https://example.com/redminote12-2.jpg', NOW(), NOW()),
    (4, 'https://example.com/macbookpro16-1.jpg', NOW(), NOW()),
    (4, 'https://example.com/macbookpro16-2.jpg', NOW(), NOW()),
    (5, 'https://example.com/dellxps15-1.jpg', NOW(), NOW()),
    (5, 'https://example.com/dellxps15-2.jpg', NOW(), NOW()),
    (6, 'https://example.com/spectrehp-1.jpg', NOW(), NOW()),
    (7, 'https://example.com/airpodspro-1.jpg', NOW(), NOW()),
    (8, 'https://example.com/galaxywatch5-1.jpg', NOW(), NOW()),
    (9, 'https://example.com/nikedrifit-1.jpg', NOW(), NOW()),
    (10, 'https://example.com/adidasultraboost-1.jpg', NOW(), NOW()),
    (11, 'https://example.com/ikeamalm-1.jpg', NOW(), NOW()),
    (12, 'https://example.com/ikeabilly-1.jpg', NOW(), NOW());

-- Insert banners
INSERT INTO banners (image, url, `index`, created_at, updated_at)
VALUES
    ('https://example.com/banner-electronics.jpg', '/category/1', 1, NOW(), NOW()),
    ('https://example.com/banner-clothing.jpg', '/category/2', 2, NOW(), NOW()),
    ('https://example.com/banner-home.jpg', '/category/3', 3, NOW(), NOW()),
    ('https://example.com/banner-sale.jpg', '/sales', 4, NOW(), NOW()),
    ('https://example.com/banner-new-arrivals.jpg', '/new-arrivals', 5, NOW(), NOW());