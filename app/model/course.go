package model

// Course adalah mata kuliah beserta jumlah terisi yang dihitung dari enrollments.
type Course struct {
	ID        int64  `json:"id"`
	KodeMK    string `json:"kode_mk"`
	NamaMK    string `json:"nama_mk"`
	SKS       int    `json:"sks"`
	Semester  int    `json:"semester"`
	Kuota     int    `json:"kuota"`
	Terisi    int    `json:"terisi"`
	SisaKuota int    `json:"sisa_kuota"`
}

// CourseFilter menampung query GET /courses.
type CourseFilter struct {
	Semester      *int
	Search        string
	AvailableOnly bool
}

// EnrolledCourse adalah mata kuliah yang ada di KRS seorang mahasiswa.
type EnrolledCourse struct {
	EnrollmentID  int64  `json:"enrollment_id"`
	TahunAkademik string `json:"tahun_akademik"`
	CourseID      int64  `json:"course_id"`
	KodeMK        string `json:"kode_mk"`
	NamaMK        string `json:"nama_mk"`
	SKS           int    `json:"sks"`
	Semester      int    `json:"semester"`
}