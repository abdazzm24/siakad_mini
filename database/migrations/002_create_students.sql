-- 002_create_students.sql
-- Tabel students: relasi 1-1 ke users, mendukung soft delete lewat deleted_at.
CREATE TABLE IF NOT EXISTS students (
    id           BIGSERIAL     PRIMARY KEY,
    user_id      BIGINT        NOT NULL REFERENCES users(id),
    nim          VARCHAR(5)    NOT NULL,
    nama         VARCHAR(150)  NOT NULL,
    prodi        VARCHAR(100)  NOT NULL,
    angkatan     SMALLINT      NOT NULL,
    ipk_terakhir NUMERIC(3, 2) NOT NULL DEFAULT 0,
    deleted_at   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    CONSTRAINT students_user_id_key    UNIQUE (user_id), -- menjamin relasi 1-1
    CONSTRAINT students_nim_key        UNIQUE (nim),
    CONSTRAINT students_nim_check      CHECK (nim ~ '^[0-9]{5}$'),
    CONSTRAINT students_angkatan_check CHECK (angkatan BETWEEN 1000 AND 9999),
    CONSTRAINT students_ipk_check      CHECK (ipk_terakhir >= 0 AND ipk_terakhir <= 4)
);

-- students.nim dan students.user_id sudah ter-index lewat UNIQUE di atas.
CREATE INDEX IF NOT EXISTS students_prodi_idx      ON students (prodi);
CREATE INDEX IF NOT EXISTS students_angkatan_idx   ON students (angkatan);
CREATE INDEX IF NOT EXISTS students_deleted_at_idx ON students (deleted_at);