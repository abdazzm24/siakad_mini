package model

import "time"

// Enrollment adalah baris pada tabel enrollments (satu mata kuliah di KRS).
type Enrollment struct {
	ID            int64     `json:"id"`
	StudentID     int64     `json:"student_id"`
	CourseID      int64     `json:"course_id"`
	TahunAkademik string    `json:"tahun_akademik"`
	CreatedAt     time.Time `json:"created_at"`
}

// EnrollmentResult adalah hasil POST /enrollments.
type EnrollmentResult struct {
	Enrollment
	Course   EnrolledCourse `json:"course"`
	TotalSKS int            `json:"total_sks"`
	BatasSKS int            `json:"batas_sks"`
	SisaSKS  int            `json:"sisa_sks"`
}