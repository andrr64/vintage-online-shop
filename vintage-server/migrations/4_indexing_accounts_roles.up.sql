-- Supaya pencarian cepat kalau cek role by user
CREATE INDEX idx_account_roles_account_id ON account_roles (account_id);

-- Supaya pencarian cepat kalau cari semua user dengan role tertentu
CREATE INDEX idx_account_roles_role_id ON account_roles (role_id);


CREATE UNIQUE INDEX idx_accounts_email ON accounts (email);
CREATE UNIQUE INDEX idx_accounts_username ON accounts (username);
