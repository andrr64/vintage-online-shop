-- Add column 'village_id' on 'addresses' with refferences to villages(id)
ALTER TABLE addresses
ADD COLUMN village_id VARCHAR(16),  -- sesuaikan tipe dengan id di villages
ADD CONSTRAINT fk_village
    FOREIGN KEY (village_id)
    REFERENCES villages(id)
    ON UPDATE CASCADE
    ON DELETE SET NULL;
