-- seed.sql
-- Data awal: 1 admin, 20 mahasiswa, 10 mata kuliah, dan beberapa enrollment.
-- Seluruh password di-hash dengan bcrypt. Aman dijalankan berulang kali (ON CONFLICT DO NOTHING).

-- Admin: admin@siakad.test / password123
INSERT INTO users (email, password, role) VALUES
    ('admin@siakad.test', '$2a$10$XBdvj1JubKWedC6c14RU7uK8kH6n6kNqS2/rcZmY56ZfY/P7Inm8e', 'admin')
ON CONFLICT (email) DO NOTHING;

-- Mahasiswa: mahasiswaNN@siakad.test, password awal = NIM masing-masing
INSERT INTO users (email, password, role) VALUES
    ('mahasiswa01@siakad.test', '$2a$10$L60LYz5asvUPKKgxAanPZ.cIdIv4GxS1lOWwXKu7YZcj3DQMZMtNC', 'mahasiswa'),
    ('mahasiswa02@siakad.test', '$2a$10$ATjXEtIt7ugBs2dz5Ev.HuyItycFP3Wl8GU/GnuhOrJlfBJnU9IZO', 'mahasiswa'),
    ('mahasiswa03@siakad.test', '$2a$10$zWWWf5o2ZQXxA9j7S465HOTIHtzlp.jARPwQWaDM1By2iWda0voIe', 'mahasiswa'),
    ('mahasiswa04@siakad.test', '$2a$10$6wpMGPR89jh9GuCvCojMSu3MF2OchBRxvslxE49SMAg.zexJ0o5ne', 'mahasiswa'),
    ('mahasiswa05@siakad.test', '$2a$10$aScCxNBYvEYjxUOjCKCFNuuPSRnu1is56fnVia9yRMtUO/522yJK6', 'mahasiswa'),
    ('mahasiswa06@siakad.test', '$2a$10$nkizAo5tDCOaxh5tb9O5peaoAVcMgKY8y5OYBVbAcnOf7nqqgT/QK', 'mahasiswa'),
    ('mahasiswa07@siakad.test', '$2a$10$5pCHu/KbGvU5go27UZ6MxOjtKqCSFQB6caAtE9trg/4r1eyy1ZKa6', 'mahasiswa'),
    ('mahasiswa08@siakad.test', '$2a$10$gRxuYfDtqnXDhw0vO7aHG.hvIX3ttmDw8h0ljrQE4/8AQO45g/Muy', 'mahasiswa'),
    ('mahasiswa09@siakad.test', '$2a$10$jtE9aap5y6TfkijHg0jjzuCToqEO0H04GDNyKPp7N06jlYyCGcG5q', 'mahasiswa'),
    ('mahasiswa10@siakad.test', '$2a$10$r6BPp4VA2MMOeXp4rn74f.NcN2tSwxi7MnwVHYTTzXhyOqp.otQ5m', 'mahasiswa'),
    ('mahasiswa11@siakad.test', '$2a$10$c3yOeo.HKVi8aMDTTcaQJeMTfjvmYQpWgIVeSDA.e4AeKujw8ucVW', 'mahasiswa'),
    ('mahasiswa12@siakad.test', '$2a$10$vwB6a6TTK/lauCOCCvPYVOpzAxqFJ6RosJIJvNQyBmeM25U4.csXa', 'mahasiswa'),
    ('mahasiswa13@siakad.test', '$2a$10$KgAA3YLr8qZvKl2oDqkX3O4ZLY2iANe4r766ZaCwSujiwfKMjcW12', 'mahasiswa'),
    ('mahasiswa14@siakad.test', '$2a$10$5P5cqCqD2XoPEnspxaobNu/aW2NgSo92eLEgKEYHxoegCQPnkjTV6', 'mahasiswa'),
    ('mahasiswa15@siakad.test', '$2a$10$ggdaU1MhB0xpYBvKijcEAehCMHZPuS4SE/RYj97ea8xwD9lVeSQfe', 'mahasiswa'),
    ('mahasiswa16@siakad.test', '$2a$10$37QIVkn9vVGV3IPegBijyeqYQaXj.DywPBictBKRJugLgS6X0LfnK', 'mahasiswa'),
    ('mahasiswa17@siakad.test', '$2a$10$eElEACdkPoDiUkz2Rn17IeSAPgQyNf.TdkjLhy17na.b8zjqMBMqO', 'mahasiswa'),
    ('mahasiswa18@siakad.test', '$2a$10$6rr9FcifTBSM.D.eyU6Bve1mDM7nVsg6Q3mJBBx90rHwMQKKHAUvy', 'mahasiswa'),
    ('mahasiswa19@siakad.test', '$2a$10$zBj0r9u4bXXdDXNyBg8IqehwKkmwwV3rieAxsQWC0OsXGxu68hbn6', 'mahasiswa'),
    ('mahasiswa20@siakad.test', '$2a$10$eDhmG8fB1vOUNEQQfdytDuhexICp9qyx0KAQzgEdf3.vRW7u35jAa', 'mahasiswa')
ON CONFLICT (email) DO NOTHING;

INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
SELECT u.id, v.nim, v.nama, v.prodi, v.angkatan, v.ipk
FROM (VALUES
    ('mahasiswa01@siakad.test', '187221000001', 'Rina Putri', 'Sistem Informasi', 2022, 3.45),
    ('mahasiswa02@siakad.test', '187221000002', 'Budi Santoso', 'Teknik Informatika', 2022, 3.50),
    ('mahasiswa03@siakad.test', '187221000003', 'Siti Aisyah', 'Teknik Informatika', 2022, 2.75),
    ('mahasiswa04@siakad.test', '187221000004', 'Ahmad Fauzi', 'Teknik Informatika', 2023, 2.00),
    ('mahasiswa05@siakad.test', '187221000005', 'Dewi Lestari', 'Sistem Informasi', 2023, 3.20),
    ('mahasiswa06@siakad.test', '187221000006', 'Eko Prasetyo', 'Teknik Informatika', 2023, 2.90),
    ('mahasiswa07@siakad.test', '187221000007', 'Fitri Handayani', 'Sistem Informasi', 2023, 3.80),
    ('mahasiswa08@siakad.test', '187221000008', 'Gilang Ramadhan', 'Teknik Informatika', 2023, 2.40),
    ('mahasiswa09@siakad.test', '187221000009', 'Hana Pertiwi', 'Sistem Informasi', 2023, 3.10),
    ('mahasiswa10@siakad.test', '187221000010', 'Irfan Hakim', 'Teknik Informatika', 2023, 2.60),
    ('mahasiswa11@siakad.test', '187221000011', 'Juwita Sari', 'Sistem Informasi', 2024, 3.55),
    ('mahasiswa12@siakad.test', '187221000012', 'Kevin Wijaya', 'Teknik Informatika', 2024, 2.30),
    ('mahasiswa13@siakad.test', '187221000013', 'Laila Nur Azizah', 'Sistem Informasi', 2024, 3.00),
    ('mahasiswa14@siakad.test', '187221000014', 'Muhammad Rizky', 'Teknik Informatika', 2024, 2.85),
    ('mahasiswa15@siakad.test', '187221000015', 'Nadia Putri', 'Teknik Informatika', 2024, 3.70),
    ('mahasiswa16@siakad.test', '187221000016', 'Oka Pratama', 'Sistem Informasi', 2024, 2.10),
    ('mahasiswa17@siakad.test', '187221000017', 'Putri Maharani', 'Teknik Informatika', 2024, 3.35),
    ('mahasiswa18@siakad.test', '187221000018', 'Rizal Firmansyah', 'Sistem Informasi', 2024, 2.55),
    ('mahasiswa19@siakad.test', '187221000019', 'Salsabila Zahra', 'Teknik Informatika', 2024, 3.90),
    ('mahasiswa20@siakad.test', '187221000020', 'Tegar Wicaksono', 'Teknik Informatika', 2024, 2.95)
) AS v(email, nim, nama, prodi, angkatan, ipk)
JOIN users u ON u.email = v.email
ORDER BY v.nim
ON CONFLICT (nim) DO NOTHING;

-- IF108 sengaja berkuota 3 dan diisi 3 mahasiswa (penuh) untuk menguji aturan kuota.
INSERT INTO courses (kode_mk, nama_mk, sks, semester, kuota) VALUES
    ('IF101', 'Pemrograman Web', 3, 3, 40),
    ('IF102', 'Basis Data', 3, 3, 40),
    ('IF103', 'Algoritma dan Struktur Data', 3, 2, 35),
    ('IF104', 'Jaringan Komputer', 3, 4, 30),
    ('IF105', 'Sistem Operasi', 3, 3, 30),
    ('IF106', 'Pemrograman Backend Lanjut', 3, 5, 25),
    ('IF107', 'Rekayasa Perangkat Lunak', 3, 4, 30),
    ('IF108', 'Kecerdasan Buatan', 3, 5, 3),
    ('IF109', 'Etika Profesi', 2, 6, 50),
    ('IF110', 'Bahasa Inggris Teknik', 2, 1, 50)
ON CONFLICT (kode_mk) DO NOTHING;

-- Enrollment awal (tahun akademik 2026/2027-Ganjil). Mahasiswa 02, 03, 04 sengaja kosong
-- agar dapat dipakai menguji batas SKS (IPK 3.50 / 2.75 / 2.00).
INSERT INTO enrollments (student_id, course_id, tahun_akademik)
SELECT s.id, c.id, '2026/2027-Ganjil'
FROM (VALUES
    (1, '187221000001', 'IF101'),
    (2, '187221000001', 'IF102'),
    (3, '187221000001', 'IF103'),
    (4, '187221000001', 'IF105'),
    (5, '187221000005', 'IF101'),
    (6, '187221000005', 'IF102'),
    (7, '187221000005', 'IF108'),
    (8, '187221000006', 'IF108'),
    (9, '187221000007', 'IF108'),
    (10, '187221000006', 'IF101'),
    (11, '187221000007', 'IF101'),
    (12, '187221000008', 'IF101'),
    (13, '187221000009', 'IF101'),
    (14, '187221000010', 'IF101'),
    (15, '187221000011', 'IF101'),
    (16, '187221000012', 'IF101'),
    (17, '187221000008', 'IF102'),
    (18, '187221000009', 'IF102'),
    (19, '187221000010', 'IF102'),
    (20, '187221000011', 'IF102'),
    (21, '187221000012', 'IF102'),
    (22, '187221000013', 'IF102'),
    (23, '187221000014', 'IF102'),
    (24, '187221000013', 'IF104'),
    (25, '187221000014', 'IF104'),
    (26, '187221000015', 'IF104'),
    (27, '187221000016', 'IF104'),
    (28, '187221000017', 'IF106'),
    (29, '187221000018', 'IF106'),
    (30, '187221000019', 'IF106'),
    (31, '187221000020', 'IF106'),
    (32, '187221000009', 'IF109'),
    (33, '187221000010', 'IF109'),
    (34, '187221000011', 'IF109'),
    (35, '187221000012', 'IF109')
) AS v(ord, nim, kode_mk)
JOIN students s ON s.nim = v.nim
JOIN courses c  ON c.kode_mk = v.kode_mk
ORDER BY v.ord
ON CONFLICT (student_id, course_id, tahun_akademik) DO NOTHING;