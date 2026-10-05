package middleware

import (
	"github.com/gofiber/fiber/v2"

	"siakad-mini/app/apperror"
	"siakad-mini/app/model"
)

// RequireRole hanya meloloskan user dengan salah satu role yang diizinkan.
// Harus dipasang SETELAH AuthMiddleware. Prinsipnya fail closed: bila identitas
// tidak ada atau role tidak dikenal, akses ditolak.
func RequireRole(roles ...string) fiber.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *fiber.Ctx) error {
		user, ok := CurrentUser(c)
		if !ok {
			return apperror.Unauthorized("Belum terautentikasi")
		}
		if _, granted := allowed[user.Role]; !granted {
			return apperror.Forbidden("Anda tidak memiliki akses ke endpoint ini")
		}
		return c.Next()
	}
}

// AdminOnly: hanya role admin.
func AdminOnly() fiber.Handler { return RequireRole(model.RoleAdmin) }

// MahasiswaOnly: hanya role mahasiswa.
func MahasiswaOnly() fiber.Handler { return RequireRole(model.RoleMahasiswa) }

// AnyRole: admin maupun mahasiswa (tetap wajib login).
func AnyRole() fiber.Handler { return RequireRole(model.RoleAdmin, model.RoleMahasiswa) }