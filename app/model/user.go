package model

import "time"

const (
	RoleAdmin     = "admin"
	RoleMahasiswa = "mahasiswa"
)

// User adalah baris pada tabel users.
type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // tidak pernah dikirim ke client
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuthUser adalah identitas yang disimpan di context request sesudah JWT diverifikasi.
type AuthUser struct {
	UserID int64
	Email  string
	Role   string
}

// MeResponse adalah payload GET /auth/me.
type MeResponse struct {
	ID      int64         `json:"id"`
	Email   string        `json:"email"`
	Role    string        `json:"role"`
	Student *StudentBrief `json:"student,omitempty"`
}

// StudentBrief ringkasan mahasiswa untuk GET /auth/me (role mahasiswa).
type StudentBrief struct {
	ID       int64  `json:"id"`
	NIM      string `json:"nim"`
	Nama     string `json:"nama"`
	Prodi    string `json:"prodi"`
	Angkatan int    `json:"angkatan"`
}
