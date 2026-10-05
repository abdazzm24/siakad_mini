-- 003_create_courses.sql
-- Tabel courses: mata kuliah beserta kuota. Jumlah terisi dihitung dari enrollments.
CREATE TABLE IF NOT EXISTS courses (
    id         BIGSERIAL    PRIMARY KEY,
    kode_mk    VARCHAR(10)  NOT NULL,
    nama_mk    VARCHAR(150) NOT NULL,
    sks        SMALLINT     NOT NULL,
    semester   SMALLINT     NOT NULL,
    kuota      INTEGER      NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT courses_kode_mk_key    UNIQUE (kode_mk),
    CONSTRAINT courses_sks_check      CHECK (sks BETWEEN 1 AND 6),
    CONSTRAINT courses_semester_check CHECK (semester BETWEEN 1 AND 8),
    CONSTRAINT courses_kuota_check    CHECK (kuota >= 0)
);

-- courses.kode_mk sudah ter-index lewat UNIQUE di atas.
CREATE INDEX IF NOT EXISTS courses_semester_idx ON courses (semester);