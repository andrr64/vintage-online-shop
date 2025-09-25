-- CORE TABLES
CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(64) UNIQUE NOT NULL,
    password VARCHAR(60) NOT NULL,
    email VARCHAR(64) UNIQUE,
    avatar_url VARCHAR(255),
    active BOOLEAN DEFAULT FALSE,
    role SMALLINT NOT NULL CHECK (role >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE product_conditions (
    id SMALLSERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE product_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE brands (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    logo_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- GEOGRAPHY TABLES
CREATE TABLE provinces (
    id VARCHAR(10) PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE regencies (
    id VARCHAR(10) PRIMARY KEY,
    province_id VARCHAR(10) NOT NULL,
    name VARCHAR(100) NOT NULL,
    FOREIGN KEY (province_id) REFERENCES provinces(id)
);

CREATE TABLE districts (
    id VARCHAR(10) PRIMARY KEY,
    regency_id VARCHAR(10) NOT NULL,
    name VARCHAR(100) NOT NULL,
    FOREIGN KEY (regency_id) REFERENCES regencies(id)
);

-- SHOP & PRODUCTS
CREATE TABLE shop (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT UNIQUE NOT NULL,
    name VARCHAR(64) UNIQUE NOT NULL,
    summary VARCHAR(255),
    description TEXT,
    active BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES accounts(id)
);

CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    shop_id BIGINT NOT NULL,
    condition_id SMALLINT NOT NULL,
    category_id INT NOT NULL,
    brand_id INT,
    name VARCHAR(128) NOT NULL,
    summary VARCHAR(255),
    description TEXT,
    price BIGINT NOT NULL CHECK (price >= 0),
    size VARCHAR(8),
    stock INT NOT NULL CHECK (stock >= 0),
    is_latest BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (shop_id) REFERENCES shop(id),
    FOREIGN KEY (condition_id) REFERENCES product_conditions(id),
    FOREIGN KEY (category_id) REFERENCES product_categories(id),
    FOREIGN KEY (brand_id) REFERENCES brands(id)
);

CREATE TABLE product_images (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    image_index SMALLINT NOT NULL CHECK (image_index >= 0),
    url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- ORDER & PAYMENT
CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL,
    total_price BIGINT NOT NULL CHECK (total_price >= 0),
    status SMALLINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES accounts(id)
);

CREATE TABLE order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    price_at_purchase BIGINT NOT NULL CHECK (price_at_purchase >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE TABLE order_status_logs (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    old_status SMALLINT,
    new_status SMALLINT NOT NULL,
    note TEXT,
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id)
);

CREATE TABLE payments (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT UNIQUE NOT NULL,
    payment_status VARCHAR(32) NOT NULL,
    midtrans_order_id VARCHAR(64) NOT NULL,
    midtrans_transaction_id VARCHAR(64),
    payment_method VARCHAR(32),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id)
);

-- USER FEATURES
CREATE TABLE addresses (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL,
    district_id VARCHAR(10) NOT NULL,
    regency_id VARCHAR(10) NOT NULL,
    province_id VARCHAR(10) NOT NULL,
    label VARCHAR(50) NOT NULL,
    recipient_name VARCHAR(100) NOT NULL,
    recipient_phone VARCHAR(20) NOT NULL,
    street TEXT NOT NULL,
    postal_code VARCHAR(10) NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES accounts(id),
    FOREIGN KEY (district_id) REFERENCES districts(id),
    FOREIGN KEY (regency_id) REFERENCES regencies(id),
    FOREIGN KEY (province_id) REFERENCES provinces(id)
);

CREATE TABLE cart (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES accounts(id)
);

CREATE TABLE cart_items (
    id BIGSERIAL PRIMARY KEY,
    cart_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (cart_id) REFERENCES cart(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id)
);

CREATE TABLE wishlist (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (account_id, product_id),
    FOREIGN KEY (account_id) REFERENCES accounts(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
);

-- REVIEWS & LOGS
CREATE TABLE reviews (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    account_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL,
    rating SMALLINT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (order_id, product_id),
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (account_id) REFERENCES accounts(id),
    FOREIGN KEY (order_id) REFERENCES orders(id)
);

CREATE TABLE shipments (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT UNIQUE NOT NULL,
    address_id BIGINT NOT NULL,
    courier VARCHAR(10) NOT NULL,
    service VARCHAR(50) NOT NULL,
    shipping_cost BIGINT NOT NULL CHECK (shipping_cost >= 0),
    tracking_number VARCHAR(100),
    status SMALLINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (address_id) REFERENCES addresses(id)
);

CREATE TABLE admin_logs (
    id BIGSERIAL PRIMARY KEY,
    admin_id BIGINT NOT NULL,
    action VARCHAR(100) NOT NULL,
    description TEXT,
    ip_address VARCHAR(45) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (admin_id) REFERENCES accounts(id)
);