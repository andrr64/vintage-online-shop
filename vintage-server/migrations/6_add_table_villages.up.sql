CREATE TABLE villages (
    id VARCHAR(24) PRIMARY KEY,
    district_id VARCHAR(16) REFERENCES districts(id),
    name VARCHAR(100) NOT NULL
);
