package request

import "strings"

// CreateEnrollmentRequest untuk POST /enrollments.
type CreateEnrollmentRequest struct {
	CourseID      int64  `json:"course_id" validate:"required,gt=0"`
	TahunAkademik string `json:"tahun_akademik" validate:"required,tahunakademik"`
}

func (r *CreateEnrollmentRequest) Normalize() {
	r.TahunAkademik = strings.TrimSpace(r.TahunAkademik)
}