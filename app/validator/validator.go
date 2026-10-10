// Package validator membungkus go-playground/validator: aturan ditulis sebagai tag
// pada struct request, pesan error diterjemahkan ke bahasa Indonesia.
package validator

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	playground "github.com/go-playground/validator/v10"
)

var (
	nimPattern           = regexp.MustCompile(`^[0-9]{5}$`)
	tahunAkademikPattern = regexp.MustCompile(`^([0-9]{4})/([0-9]{4})-(Ganjil|Genap)$`)
)

// validate dibuat sekali karena validator.New() menyimpan cache refleksi.
var validate = newValidator()

func newValidator() *playground.Validate {
	v := playground.New()

	// Pesan error memakai nama field JSON ("nim"), bukan nama field Go ("NIM").
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
		if name == "" || name == "-" {
			return f.Name
		}
		return name
	})

	_ = v.RegisterValidation("nim", func(fl playground.FieldLevel) bool {
		return nimPattern.MatchString(fl.Field().String())
	})

	// angkatan tidak boleh melebihi tahun berjalan.
	_ = v.RegisterValidation("notfuture", func(fl playground.FieldLevel) bool {
		return fl.Field().Int() <= int64(time.Now().Year())
	})

	// Format 2026/2027-Ganjil: dua tahun berurutan, lalu Ganjil atau Genap.
	_ = v.RegisterValidation("tahunakademik", func(fl playground.FieldLevel) bool {
		return IsTahunAkademik(fl.Field().String())
	})
	return v
}

// IsTahunAkademik memeriksa format tahun akademik, contoh: 2026/2027-Ganjil.
func IsTahunAkademik(value string) bool {
	m := tahunAkademikPattern.FindStringSubmatch(value)
	if m == nil {
		return false
	}
	first, _ := strconv.Atoi(m[1])
	second, _ := strconv.Atoi(m[2])
	return second == first+1
}

// Validate menjalankan seluruh aturan pada tag struct.
// Mengembalikan nil bila lolos; bila gagal, peta field -> daftar pesan.
func Validate(s any) map[string][]string {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var invalid *playground.InvalidValidationError
	if errors.As(err, &invalid) {
		return map[string][]string{"body": {"Data tidak valid"}}
	}

	var fieldErrors playground.ValidationErrors
	if !errors.As(err, &fieldErrors) {
		return map[string][]string{"body": {"Validasi gagal"}}
	}

	result := make(map[string][]string, len(fieldErrors))
	for _, fe := range fieldErrors {
		result[fe.Field()] = append(result[fe.Field()], messageFor(fe))
	}
	return result
}

var labels = map[string]string{
	"nim":            "NIM",
	"nama":           "Nama",
	"email":          "Email",
	"password":       "Password",
	"prodi":          "Prodi",
	"angkatan":       "Angkatan",
	"ipk_terakhir":   "IPK",
	"course_id":      "course_id",
	"tahun_akademik": "Tahun akademik",
}

// Pesan khusus untuk kombinasi field + tag tertentu.
var specialMessages = map[string]string{
	"angkatan.gte":     "Angkatan harus 4 digit",
	"angkatan.lte":     "Angkatan harus 4 digit",
	"ipk_terakhir.gte": "IPK harus antara 0.00 sampai 4.00",
	"ipk_terakhir.lte": "IPK harus antara 0.00 sampai 4.00",
	"email.email":      "Format email tidak valid",
	"nim.nim":          "NIM harus terdiri dari 5 digit angka",
}

func messageFor(fe playground.FieldError) string {
	if msg, ok := specialMessages[fe.Field()+"."+fe.Tag()]; ok {
		return msg
	}

	label := labels[fe.Field()]
	if label == "" {
		label = fe.Field()
	}

	switch fe.Tag() {
	case "required":
		return label + " wajib diisi"
	case "email":
		return "Format email tidak valid"
	case "min":
		if fe.Kind() == reflect.String {
			return label + " minimal " + fe.Param() + " karakter"
		}
		return label + " minimal " + fe.Param()
	case "max":
		if fe.Kind() == reflect.String {
			return label + " maksimal " + fe.Param() + " karakter"
		}
		return label + " maksimal " + fe.Param()
	case "gt":
		return label + " harus lebih besar dari " + fe.Param()
	case "notfuture":
		return label + " tidak boleh melebihi tahun berjalan"
	case "tahunakademik":
		return "Format tahun akademik tidak valid, contoh: 2026/2027-Ganjil"
	default:
		return label + " tidak valid"
	}
}