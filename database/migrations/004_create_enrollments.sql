-- 004_create_enrollments.sql
-- Tabel enrollments: KRS. Satu baris = satu mata kuliah yang diambil seorang mahasiswa.
CREATE TABLE IF NOT EXISTS enrollments (
    id             BIGSERIAL   PRIMARY KEY,
    student_id     BIGINT      NOT NULL REFERENCES students(id),
    course_id      BIGINT      NOT NULL REFERENCES courses(id),
    tahun_akademik VARCHAR(20) NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    -- Mencegah mata kuliah yang sama diambil dua kali pada tahun akademik yang sama.
    CONSTRAINT enrollments_student_course_term_key
        UNIQUE (student_id, course_id, tahun_akademik)
);

CREATE INDEX IF NOT EXISTS enrollments_student_id_idx ON enrollments (student_id);
CREATE INDEX IF NOT EXISTS enrollments_course_id_idx  ON enrollments (course_id);