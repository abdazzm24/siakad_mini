# Pengujian API SIAKAD Mini

Base URL: `http://localhost:3000/api/v1`. Isi `TOKEN` dengan `access_token` hasil login.
Data di bawah memakai hasil seeder pada database baru.

| Akun | Email | Password | IPK | Batas SKS |
|---|---|---|---|---|
| Admin | admin@siakad.test | password123 | - | - |
| Mahasiswa A (Rina, id 1) | mahasiswa01@siakad.test | 187221000001 | 3.45 | 24 |
| Mahasiswa B (Dewi, id 5) | mahasiswa05@siakad.test | 187221000005 | 3.20 | 24 |
| IPK 3.50 (id 2) | mahasiswa02@siakad.test | 187221000002 | 3.50 | 24 |
| IPK 2.75 (id 3) | mahasiswa03@siakad.test | 187221000003 | 2.75 | 21 |
| IPK 2.00 (id 4) | mahasiswa04@siakad.test | 187221000004 | 2.00 | 18 |

Mahasiswa 02, 03, 04 sengaja belum punya KRS agar mudah dipakai menguji batas SKS.

## curl: 10 endpoint

### 1. Login (publik) -> 200
```bash
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@siakad.test", "password": "password123"}'
```
Response:
```json
{"success":true,"message":"Login berhasil","access_token":"<jwt>","token_type":"Bearer","expires_in":900,"data":{"id":1,"email":"admin@siakad.test","role":"admin"}}
```

### 2. Profil login -> 200
```bash
curl http://localhost:3000/api/v1/auth/me -H "Authorization: Bearer TOKEN"
```
Mahasiswa mendapat objek `student` berisi `nim`, `nama`, `prodi`, `angkatan`.

### 3. Daftar mahasiswa (admin) -> 200
```bash
curl "http://localhost:3000/api/v1/students?page=1&per_page=10&prodi=Teknik%20Informatika&angkatan=2024&search=a&sort=-ipk_terakhir" \
  -H "Authorization: Bearer ADMIN_TOKEN"
```
`meta` = `{"current_page":1,"per_page":10,"total":20,"last_page":2}` (tanpa filter).

### 4. Tambah mahasiswa (admin) -> 201
```bash
curl -X POST http://localhost:3000/api/v1/students \
  -H "Authorization: Bearer ADMIN_TOKEN" -H "Content-Type: application/json" \
  -d '{"nim":"187221000021","nama":"Mahasiswa Baru","email":"mahasiswa21@siakad.test","prodi":"Teknik Informatika","angkatan":2024,"ipk_terakhir":3.50}'
```
Password awal akun baru = `187221000021`.

### 5. Detail mahasiswa -> 200
```bash
curl http://localhost:3000/api/v1/students/1 -H "Authorization: Bearer TOKEN"
```
Berisi `courses`, `total_sks`, `batas_sks`. Opsional `?tahun_akademik=2026/2027-Ganjil`.

### 6. Update mahasiswa (admin) -> 200
```bash
curl -X PUT http://localhost:3000/api/v1/students/1 \
  -H "Authorization: Bearer ADMIN_TOKEN" -H "Content-Type: application/json" \
  -d '{"nama":"Rina Putri","prodi":"Sistem Informasi","angkatan":2022,"ipk_terakhir":3.60}'
```
Field `nim` yang ikut dikirim diabaikan (NIM tidak bisa diubah).

### 7. Soft delete (admin) -> 204
```bash
curl -X DELETE http://localhost:3000/api/v1/students/21 -H "Authorization: Bearer ADMIN_TOKEN"
```

### 8. Daftar mata kuliah -> 200
```bash
curl "http://localhost:3000/api/v1/courses?semester=3&search=Pemrograman&available=true" \
  -H "Authorization: Bearer TOKEN"
```
Tiap item memuat `terisi` dan `sisa_kuota`.

### 9. Ambil mata kuliah (mahasiswa) -> 201
```bash
curl -X POST http://localhost:3000/api/v1/enrollments \
  -H "Authorization: Bearer MAHASISWA_TOKEN" -H "Content-Type: application/json" \
  -d '{"course_id": 1, "tahun_akademik": "2026/2027-Ganjil"}'
```

### 10. Batalkan mata kuliah (mahasiswa pemilik) -> 204
```bash
curl -X DELETE http://localhost:3000/api/v1/enrollments/ENROLLMENT_ID -H "Authorization: Bearer MAHASISWA_TOKEN"
```
`ENROLLMENT_ID` = `courses[].enrollment_id` dari `GET /students/{id}`.

## Daftar pengujian (Test 1-23)

| # | Skenario | Hasil |
|---|---|---|
| 1 | Login admin | 200 |
| 2 | Login password salah | 401 `Email atau password salah` |
| 3 | Login body tidak valid (`{"email":"x","password":"1"}`) | 422 + `errors` |
| 4 | 6x login gagal dalam 1 menit | 429 + header `Retry-After: 60` |
| 5 | `GET /auth/me` tanpa token | 401 |
| 6 | Admin `GET /students` | 200 |
| 7 | Mahasiswa `GET /students` | 403 |
| 8 | Admin `POST /students` | 201 |
| 9 | `POST /students` NIM duplikat | 422 `NIM sudah terdaftar` |
| 10 | Admin `GET /students/1` | 200 |
| 11 | Mahasiswa A `GET /students/1` | 200 |
| 12 | Mahasiswa A `GET /students/5` | 403 |
| 13 | Admin `PUT /students/1` | 200 |
| 14 | Mahasiswa `PUT /students/1` | 403 |
| 15 | Admin `DELETE /students/{id}` | 204 |
| 16 | Login mahasiswa yang sudah di-soft delete | 401 |
| 17 | `GET /courses` | 200 |
| 18 | Mahasiswa `POST /enrollments` | 201 |
| 19 | `POST /enrollments` course yang sama | 409 `Mata kuliah sudah diambil pada tahun akademik tersebut` |
| 20 | Course penuh (`course_id` 8 / IF108) | 422 `Kuota mata kuliah sudah penuh` |
| 21 | Total SKS melebihi batas | 422 `Batas SKS terlampaui. Sisa SKS yang dapat diambil: N` |
| 22 | Mahasiswa A hapus enrollment miliknya | 204 |
| 23 | Mahasiswa A hapus enrollment mahasiswa B | 403 |

## Business rule SKS

Semua mata kuliah seed bernilai 3 SKS kecuali IF109 dan IF110 (2 SKS). Ambil `course_id` 1-7 (7 x 3 = 21 SKS), lalu `9` (2 SKS).

- **IPK 3.50 (mahasiswa02, batas 24):** ambil course 1-7 -> total 21; ambil course 9 -> 201 (total 23); ambil course 10 (2 SKS) -> **422 `Sisa SKS yang dapat diambil: 1`**.
- **IPK 2.75 (mahasiswa03, batas 21):** ambil course 1-7 -> total 21; ambil course 9 -> **422 `Sisa SKS yang dapat diambil: 0`**.
- **IPK 2.00 (mahasiswa04, batas 18):** ambil course 1-6 -> total 18; ambil course 7 -> **422 `Sisa SKS yang dapat diambil: 0`**.

## Ownership

Mahasiswa A = mahasiswa01 (id 1), Mahasiswa B = mahasiswa05 (id 5).

- A `GET /students/1` -> 200; A `GET /students/5` -> **403**
- Ambil `enrollment_id` A dan B dari `GET /students/1` dan `GET /students/5`.
- A `DELETE /enrollments/{id-A}` -> 204; A `DELETE /enrollments/{id-B}` -> **403**

## Pagination, search, filter, course filter

```text
GET /students?page=1&per_page=10          -> meta total 20, last_page 2
GET /students?page=2&per_page=10          -> current_page 2
GET /students?per_page=999                -> per_page dibatasi 50
GET /students?search=Rina                 -> cocok nama
GET /students?search=187221000001         -> cocok NIM
GET /students?prodi=Teknik%20Informatika  -> filter prodi
GET /students?angkatan=2024               -> filter angkatan
GET /students?sort=nama                   -> urut nama A-Z
GET /students?sort=-ipk_terakhir          -> IPK tertinggi dulu
GET /courses?semester=3
GET /courses?search=Pemrograman
GET /courses?available=true               -> IF108 (penuh) tidak muncul
```

## Pengujian otomatis

Seluruh skenario di atas (63 pemeriksaan) dijalankan oleh `bash tests/api_test.sh`.