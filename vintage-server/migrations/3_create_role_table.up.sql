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

-- Fungsi trigger
CREATE OR REPLACE FUNCTION check_unique_email_username_per_role()
RETURNS TRIGGER AS $$
DECLARE
    new_email TEXT;
    new_username TEXT;
BEGIN
    -- Ambil email & username dari akun yang mau dikasih role
    SELECT email, username INTO new_email, new_username
    FROM accounts
    WHERE id = NEW.account_id;

    -- Cek duplikat email untuk role yang sama
    IF EXISTS (
        SELECT 1
        FROM account_roles ar
        JOIN accounts a ON a.id = ar.account_id
        WHERE a.email = new_email
          AND ar.role_id = NEW.role_id
          AND a.email IS NOT NULL
    ) THEN
        RAISE EXCEPTION 'Email sudah digunakan untuk role ini';
    END IF;

    -- Cek duplikat username untuk role yang sama
    IF EXISTS (
        SELECT 1
        FROM account_roles ar
        JOIN accounts a ON a.id = ar.account_id
        WHERE a.username = new_username
          AND ar.role_id = NEW.role_id
    ) THEN
        RAISE EXCEPTION 'Username sudah digunakan untuk role ini';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger
CREATE TRIGGER trg_unique_email_username_per_role
BEFORE INSERT ON account_roles
FOR EACH ROW
EXECUTE FUNCTION check_unique_email_username_per_role();