-- Admin & Logs
CREATE TABLE admins (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    username VARCHAR(100) NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE admin_logs (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    admin_id BIGINT NOT NULL REFERENCES admins(id) ON DELETE CASCADE,
    action VARCHAR(255)
);

-- Master Data
CREATE TABLE brands (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    name VARCHAR(255) NOT NULL
);

-- Product & Versioning
CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    name VARCHAR(255) NOT NULL,
    brand_id BIGINT NOT NULL REFERENCES brands(id) ON DELETE CASCADE,
    category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE CASCADE
);

CREATE TABLE product_versions (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    previous_version_id BIGINT REFERENCES product_versions(id),
    size VARCHAR(50),
    color VARCHAR(50),
    description TEXT,
    price DOUBLE PRECISION NOT NULL,
    stock INT NOT NULL,
    active_from TIMESTAMP,
    active_to TIMESTAMP
);

-- User Related
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    username VARCHAR(100) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    fullname TEXT NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    profile_picture VARCHAR(255)
);

CREATE TABLE addresses (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    full_address TEXT,
    kel INT,
    kec INT,
    kab INT,
    prov INT,
    kode_pos VARCHAR(10)
);

CREATE TABLE wishlists (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    version_id BIGINT NOT NULL REFERENCES product_versions(id) ON DELETE CASCADE
);

CREATE TABLE carts (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE cart_details (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    cart_id BIGINT NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
    version_id BIGINT NOT NULL REFERENCES product_versions(id) ON DELETE CASCADE,
    quantity INT
);

-- Transactions
CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    total_amount DOUBLE PRECISION,
    total_quantity INT,
    status VARCHAR(50)
);

CREATE TABLE transaction_details (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    transaction_id BIGINT NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    version_id BIGINT NOT NULL REFERENCES product_versions(id) ON DELETE CASCADE,
    quantity INT,
    total DOUBLE PRECISION
);
