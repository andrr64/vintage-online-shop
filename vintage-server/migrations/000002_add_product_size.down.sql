-- Balikin tabel products ke kondisi semula
ALTER TABLE products DROP COLUMN size_id,
    ADD COLUMN size VARCHAR(8);
-- Hapus tabel product_size
DROP TABLE product_size; 