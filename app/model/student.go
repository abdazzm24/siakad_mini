package model

import "time"

// Student adalah baris pada tabel students (digabung email dari users).
type Student struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	NIM         string     `json:"nim"`
	Nama        string     `json:"nama"`
	Prodi       string     `json:"prodi"`
	Angkatan    int        `json:"angkatan"`
	IPKTerakhir float64    `json:"ipk_terakhir"`
	Email       string     `json:"email,omitempty"`
	DeletedAt   *time.Time `json:"-"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// StudentDetail dipakai oleh GET /students/{id}.
type StudentDetail struct {
	Student
	Courses  []EnrolledCourse `json:"courses"`
	TotalSKS int              `json:"total_sks"`
	BatasSKS int              `json:"batas_sks"`
}

// StudentListQuery menampung parameter GET /students yang sudah dinormalisasi.
type StudentListQuery struct {
	Page     int
	PerPage  int
	Prodi    string
	Angkatan *int
	Search   string
	Sort     string // nama, -nama, ipk_terakhir, -ipk_terakhir, nim, -nim, angkatan, -angkatan
}

func (q StudentListQuery) Offset() int {
	return (q.Page - 1) * q.PerPage
}

// StudentUpdate berisi field yang boleh diubah. Field nil dibiarkan seperti semula.
type StudentUpdate struct {
	Nama        *string
	Prodi       *string
	Angkatan    *int
	IPKTerakhir *float64
}