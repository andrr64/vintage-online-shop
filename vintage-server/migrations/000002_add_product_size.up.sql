-- Buat tabel product_size
CREATE TABLE product_size (
    id SERIAL PRIMARY KEY,
    size_name VARCHAR(32) NOT NULL UNIQUE
);

-- Ubah tabel products
ALTER TABLE products
    DROP COLUMN size,
    ADD COLUMN size_id INT REFERENCES product_size(id);
