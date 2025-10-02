-- BENAR
CREATE UNIQUE INDEX product_condition_name_lower_idx 
ON product_conditions (LOWER(name));