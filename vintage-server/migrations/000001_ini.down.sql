-- =========================================
-- DOWN: Drop all tables in correct order (reverse of creation)
-- =========================================

-- Drop tables with foreign key dependencies first
DROP TABLE IF EXISTS admin_logs;
DROP TABLE IF EXISTS shipments;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS wishlist;
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS cart;
DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS order_status_logs;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS product_images;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS shop;

-- Drop reference tables
DROP TABLE IF EXISTS districts;
DROP TABLE IF EXISTS regencies;
DROP TABLE IF EXISTS provinces;
DROP TABLE IF EXISTS brands;
DROP TABLE IF EXISTS product_categories;
DROP TABLE IF EXISTS product_conditions;

-- Drop core tables last
DROP TABLE IF EXISTS accounts;

-- =========================================
-- ALTERNATIVE: Using SET FOREIGN_KEY_CHECKS (MySQL)
-- =========================================
/*
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS admin_logs;
DROP TABLE IF EXISTS shipments;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS wishlist;
DROP TABLE IF EXISTS cart_items;
DROP TABLE IF EXISTS cart;
DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS order_status_logs;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS product_images;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS shop;
DROP TABLE IF EXISTS districts;
DROP TABLE IF EXISTS regencies;
DROP TABLE IF EXISTS provinces;
DROP TABLE IF EXISTS brands;
DROP TABLE IF EXISTS product_categories;
DROP TABLE IF EXISTS product_conditions;
DROP TABLE IF EXISTS accounts;

SET FOREIGN_KEY_CHECKS = 1;
*/