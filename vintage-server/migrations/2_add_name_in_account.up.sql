-- 0000002 add 'firstname' and 'lastname' to 'accounts table'
ALTER TABLE accounts
    ADD COLUMN firstname VARCHAR(64),
    ADD COLUMN lastname VARCHAR(64);