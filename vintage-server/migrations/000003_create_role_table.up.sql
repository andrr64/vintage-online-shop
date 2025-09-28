CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(32) UNIQUE NOT NULL
);


INSERT INTO roles (name)
VALUES 
  ('customer'),
  ('seller'),
  ('admin');


CREATE TABLE account_roles (
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    role_id INT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (account_id, role_id)
);


INSERT INTO account_roles (account_id, role_id)
SELECT id, role
FROM accounts;

ALTER TABLE accounts DROP COLUMN role;