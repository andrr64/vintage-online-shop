-- =========================================
-- UP: Create all tables
-- =========================================
-- =====================
-- LOCATION / ADDRESS
-- =====================
CREATE TABLE provinsis (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);
CREATE TABLE kabupaten_kotas (
    id BIGSERIAL PRIMARY KEY,
    provinsi_id BIGINT NOT NULL REFERENCES provinsis(id),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);
CREATE TABLE kecamatans (
    id BIGSERIAL PRIMARY KEY,
    kabupaten_kota_id BIGINT NOT NULL REFERENCES kabupaten_kotas(id),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);
-- =====================
-- USER DOMAIN
-- =====================
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    profile_picture VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);
CREATE TABLE addresses (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    kecamatan_id BIGINT NOT NULL REFERENCES kecamatans(id),
    label VARCHAR(50),
    recipient_name VARCHAR(100),
    phone_number VARCHAR(20),
    full_address TEXT,
    postal_code VARCHAR(10),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);
CREATE TABLE carts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE NOT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);
CREATE TABLE cart_details (
    cart_id BIGINT NOT NULL REFERENCES carts(id),
    product_variant_id BIGINT NOT NULL,
    quantity INT NOT NULL,
    added_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (cart_id, product_variant_id)
);
CREATE TABLE wishlists (
    user_id BIGINT NOT NULL REFERENCES users(id),
    product_variant_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, product_variant_id)
);
-- =====================
-- PRODUCT DOMAIN
-- =====================
CREATE TABLE brands (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);
CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    brand_id BIGINT NOT NULL REFERENCES brands(id),
    base_name VARCHAR(255) NOT NULL,
    base_description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);
CREATE TABLE product_categories (
    product_id BIGINT NOT NULL REFERENCES products(id),
    category_id BIGINT NOT NULL REFERENCES categories(id),
    PRIMARY KEY (product_id, category_id)
);
CREATE TABLE product_variants (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id),
    variant_name_extension VARCHAR(255),
    price NUMERIC(12, 2),
    stock_quantity INT
);
CREATE TABLE product_images (
    id BIGSERIAL PRIMARY KEY,
    product_variant_id BIGINT NOT NULL REFERENCES product_variants(id),
    image_url VARCHAR(255) NOT NULL,
    alt_text VARCHAR(255),
    display_order INT
);
CREATE TABLE options (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);
CREATE TABLE option_values (
    id BIGSERIAL PRIMARY KEY,
    option_id BIGINT NOT NULL REFERENCES options(id),
    name VARCHAR(255) NOT NULL
);
CREATE TABLE variant_options (
    product_variant_id BIGINT NOT NULL REFERENCES product_variants(id),
    option_value_id BIGINT NOT NULL REFERENCES option_values(id),
    PRIMARY KEY (product_variant_id, option_value_id)
);
CREATE TABLE reviews (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    product_id BIGINT NOT NULL REFERENCES products(id),
    rating INT NOT NULL,
    comment TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- =====================
-- ORDER DOMAIN
-- =====================
CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    shipping_address_id BIGINT NOT NULL REFERENCES addresses(id),
    order_date TIMESTAMP NOT NULL DEFAULT NOW(),
    current_status VARCHAR(50)
);
CREATE TABLE order_line_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id),
    product_variant_id BIGINT NOT NULL REFERENCES product_variants(id),
    quantity INT NOT NULL,
    price_at_transaction NUMERIC(12, 2)
);
CREATE TABLE order_status_histories (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id),
    status VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- =====================
-- ADMIN DOMAIN
-- =====================
CREATE TABLE admins (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);
CREATE TABLE admin_logs (
    id BIGSERIAL PRIMARY KEY,
    admin_id BIGINT NOT NULL REFERENCES admins(id),
    action TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT NOW()
);