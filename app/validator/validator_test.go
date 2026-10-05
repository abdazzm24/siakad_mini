package validator

import (
	"testing"
	"time"

	"siakad-mini/app/request"
)

func TestIsTahunAkademik(t *testing.T) {
	valid := []string{"2026/2027-Ganjil", "2025/2026-Genap"}
	invalid := []string{"", "2026-2027-Ganjil", "2026/2028-Ganjil", "2026/2027-ganjil", "2026/2027"}
	for _, v := range valid {
		if !IsTahunAkademik(v) {
			t.Errorf("%q seharusnya valid", v)
		}
	}
	for _, v := range invalid {
		if IsTahunAkademik(v) {
			t.Errorf("%q seharusnya tidak valid", v)
		}
	}
}

func TestCreateStudentValidation(t *testing.T) {
	ok := request.CreateStudentRequest{
		NIM: "187221000021", Nama: "Mahasiswa Baru", Email: "baru@siakad.test",
		Prodi: "Teknik Informatika", Angkatan: 2024,
	}
	if errs := Validate(&ok); errs != nil {
		t.Fatalf("seharusnya lolos, dapat %v", errs)
	}

	bad := request.CreateStudentRequest{
		NIM: "123", Email: "bukan-email", Angkatan: time.Now().Year() + 1,
	}
	errs := Validate(&bad)
	for _, field := range []string{"nim", "nama", "email", "prodi", "angkatan"} {
		if len(errs[field]) == 0 {
			t.Errorf("field %s seharusnya gagal validasi, errs=%v", field, errs)
		}
	}

	ipk := 4.5
	ok.IPKTerakhir = &ipk
	if errs := Validate(&ok); len(errs["ipk_terakhir"]) == 0 {
		t.Error("IPK 4.5 seharusnya ditolak")
	}
}