-- BENAR
CREATE UNIQUE INDEX product_categories_name_lower_idx 
ON product_categories (LOWER(name));