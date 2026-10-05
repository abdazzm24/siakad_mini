package handler

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"siakad-mini/app/apperror"
	"siakad-mini/app/middleware"
	"siakad-mini/app/model"
	"siakad-mini/app/validator"
)

// requestContext memberi batas waktu pada setiap operasi database.
func requestContext(c *fiber.Ctx) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.UserContext(), 10*time.Second)
}

// bindAndValidate membaca body JSON ke dst lalu menjalankan validasi tag.
// Body yang bukan JSON sah dijawab 422, sama seperti kegagalan validasi lainnya.
func bindAndValidate(c *fiber.Ctx, dst any) error {
	if err := json.Unmarshal(c.Body(), dst); err != nil {
		return apperror.Unprocessable("Body request harus berupa JSON yang valid")
	}
	return nil
}

// validateStruct mengembalikan *AppError 422 bila ada aturan yang dilanggar.
func validateStruct(s any) error {
	if errs := validator.Validate(s); errs != nil {
		return apperror.Validation(errs)
	}
	return nil
}

// pathID membaca parameter :id. ID yang bukan angka positif dianggap tidak ditemukan.
func pathID(c *fiber.Ctx) (int64, error) {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, apperror.NotFound("Data tidak ditemukan")
	}
	return id, nil
}

// authUser mengambil identitas yang disimpan AuthMiddleware.
func authUser(c *fiber.Ctx) (model.AuthUser, error) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		return model.AuthUser{}, apperror.Unauthorized("Belum terautentikasi")
	}
	return user, nil
}

// emptyUpdateError dipakai PUT /students/{id} bila tidak ada field yang dikirim.
func emptyUpdateError() error {
	return apperror.Validation(map[string][]string{
		"body": {"Minimal satu field harus diisi: nama, prodi, angkatan, atau ipk_terakhir"},
	})
}