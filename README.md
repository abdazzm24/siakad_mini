# SIAKAD Mini RESTful API

RESTful API untuk Sistem Informasi Akademik sederhana (UTS Pemrograman Backend):
mengelola data mahasiswa, mata kuliah, dan Kartu Rencana Studi (KRS) dengan dua role: `admin` dan `mahasiswa`.

## Features

- JWT authentication (`Authorization: Bearer <token>`)
- Role authorization (`AdminOnly`, `MahasiswaOnly`) + ownership check (`student.user_id == user login`)
- Student management (CRUD, soft delete)
- Course management (sisa kuota dihitung dari tabel enrollments)
- KRS management (ambil / batalkan mata kuliah)
- Pagination, search, filter, sort
- Transaction + row locking (`SELECT ... FOR UPDATE`) pada POST students dan POST enrollments
- Business rule SKS (24 / 21 / 18 berdasarkan IPK), anti-duplikasi, kuota
- Rate limiting login (5 percobaan gagal / menit / IP)
- Password di-hash bcrypt, query terparameterisasi, error handling terpusat

## Requirements

- Go 1.22+
- PostgreSQL 14+
- Git
- Postman atau curl

## Installation

```bash
git clone <repository-url>
cd siakad-mini
go mod tidy
```

1. Buat database:

```bash
psql -U postgres -c "CREATE DATABASE siakad_mini;"
```

2. Salin environment lalu sesuaikan `DB_PASSWORD` dan `JWT_SECRET`:

```bash
cp .env.example .env
```

3. Migration dan seeder (aman diulang):

```bash
go run . migrate
go run . seed
```

4. Jalankan server:

```bash
go run .
```

Server berjalan di `http://localhost:3000`, prefix API `/api/v1`.

## Akun Seeder

| Role | Email | Password |
|---|---|---|
| admin | admin@siakad.test | password123 |
| mahasiswa | mahasiswa01@siakad.test … mahasiswa20@siakad.test | NIM masing-masing (`1872210000NN`, contoh `187221000001`) |

Seeder membuat 1 admin, 20 mahasiswa, 10 mata kuliah, dan 35 enrollment awal
(tahun akademik `2026/2027-Ganjil`). Mata kuliah `IF108` sengaja penuh (kuota 3, terisi 3).

## Struktur Project

```text
app/
  apperror/    error aplikasi (status + pesan aman)
  config/      env + koneksi database
  handler/     HTTP: baca request, panggil service, tulis response
  middleware/  auth JWT, role, rate limit, error handler
  model/       struct data
  repository/  seluruh query SQL
  request/     DTO + tag validasi
  response/    format response JSON seragam
  router/      peta endpoint
  service/     business logic (SKS, kuota, duplikasi, ownership, JWT)
  validator/   go-playground/validator + pesan Indonesia
database/
  migrations/  001-004 (.sql, di-embed ke binary)
  seeders/     seed.sql
tests/api_test.sh   pengujian end-to-end otomatis
```

## Endpoint

| # | Method | Path | Akses |
|---|---|---|---|
| 1 | POST | /api/v1/auth/login | Publik |
| 2 | GET | /api/v1/auth/me | Admin, Mahasiswa |
| 3 | GET | /api/v1/students | Admin |
| 4 | POST | /api/v1/students | Admin |
| 5 | GET | /api/v1/students/{id} | Admin, Mahasiswa pemilik |
| 6 | PUT | /api/v1/students/{id} | Admin |
| 7 | DELETE | /api/v1/students/{id} | Admin |
| 8 | GET | /api/v1/courses | Admin, Mahasiswa |
| 9 | POST | /api/v1/enrollments | Mahasiswa |
| 10 | DELETE | /api/v1/enrollments/{id} | Mahasiswa pemilik |

Contoh curl dan skenario pengujian lengkap: [docs/API_TESTING.md](docs/API_TESTING.md).

## Pengujian

```bash
go vet ./...
go test ./...            # unit test business rule + validator
bash tests/api_test.sh   # end-to-end (server hidup, DB baru di-migrate + seed)
```

`tests/api_test.sh` memuat jeda 61 detik di akhir untuk menguji rate limit 429.
Bila dijalankan ulang, reset dulu database (`DROP DATABASE`, `CREATE DATABASE`, migrate, seed).

## Catatan Desain

- Password awal mahasiswa baru = NIM (di-hash bcrypt). `password` tidak pernah keluar di JSON (`json:"-"`).
- Soft delete: semua query mahasiswa memakai `deleted_at IS NULL`; mahasiswa terhapus tidak bisa login,
  tidak muncul di daftar, dan detailnya 404. NIM-nya tetap tidak bisa dipakai ulang.
- `POST /enrollments` mengunci baris student lalu course (selalu urutan itu, agar tidak deadlock)
  sehingga dua request bersamaan tidak bisa sama-sama lolos cek kuota maupun batas SKS.
- Batas SKS dihitung per tahun akademik.
- `APP_ENV=production` menyembunyikan detail error; detail teknis hanya masuk log server.