-- 001_create_users.sql
-- Tabel users: akun login untuk admin dan mahasiswa.
CREATE TABLE IF NOT EXISTS users (
    id         BIGSERIAL    PRIMARY KEY,
    email      VARCHAR(255) NOT NULL,
    password   VARCHAR(255) NOT NULL, -- selalu hash bcrypt, tidak pernah plaintext
    role       VARCHAR(20)  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    -- UNIQUE otomatis membuat index pada kolom email
    CONSTRAINT users_email_key  UNIQUE (email),
    CONSTRAINT users_role_check CHECK (role IN ('admin', 'mahasiswa'))
);