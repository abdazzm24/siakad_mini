package request

import "strings"

// CreateStudentRequest untuk POST /students.
type CreateStudentRequest struct {
	NIM         string   `json:"nim" validate:"required,nim"`
	Nama        string   `json:"nama" validate:"required,max=150"`
	Email       string   `json:"email" validate:"required,email,max=255"`
	Prodi       string   `json:"prodi" validate:"required,max=100"`
	Angkatan    int      `json:"angkatan" validate:"required,gte=1000,lte=9999,notfuture"`
	IPKTerakhir *float64 `json:"ipk_terakhir" validate:"omitnil,gte=0,lte=4"`
}

func (r *CreateStudentRequest) Normalize() {
	r.NIM = strings.TrimSpace(r.NIM)
	r.Nama = strings.TrimSpace(r.Nama)
	r.Email = strings.ToLower(strings.TrimSpace(r.Email))
	r.Prodi = strings.TrimSpace(r.Prodi)
}

// UpdateStudentRequest untuk PUT /students/{id}.
// Field nim sengaja tidak ada: NIM tidak boleh diubah, kalaupun dikirim diabaikan.
// Field yang tidak dikirim (nil) dibiarkan seperti data semula.
type UpdateStudentRequest struct {
	Nama        *string  `json:"nama" validate:"omitnil,min=1,max=150"`
	Prodi       *string  `json:"prodi" validate:"omitnil,min=1,max=100"`
	Angkatan    *int     `json:"angkatan" validate:"omitnil,gte=1000,lte=9999,notfuture"`
	IPKTerakhir *float64 `json:"ipk_terakhir" validate:"omitnil,gte=0,lte=4"`
}

func (r *UpdateStudentRequest) Normalize() {
	if r.Nama != nil {
		v := strings.TrimSpace(*r.Nama)
		r.Nama = &v
	}
	if r.Prodi != nil {
		v := strings.TrimSpace(*r.Prodi)
		r.Prodi = &v
	}
}

// IsEmpty true bila tidak ada satu pun field yang dikirim.
func (r UpdateStudentRequest) IsEmpty() bool {
	return r.Nama == nil && r.Prodi == nil && r.Angkatan == nil && r.IPKTerakhir == nil
}