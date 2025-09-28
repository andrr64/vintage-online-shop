ALTER TABLE districts
ALTER COLUMN id TYPE character varying(16);
-- Kalau mau, bisa sekalian pastikan foreign key ke regencies masih valid
-- (tidak perlu diubah kalau id regency tetap sama)