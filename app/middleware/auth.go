package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"

	"siakad-mini/app/apperror"
	"siakad-mini/app/model"
	"siakad-mini/app/service"
)

// localsAuthUser adalah kunci penyimpanan identitas pada context request.
const localsAuthUser = "authUser"

// AuthMiddleware memeriksa header "Authorization: Bearer <token>".
// Bila sah, identitas disimpan di Locals; bila tidak, jawabannya 401.
func AuthMiddleware(jwt *service.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, err := bearerToken(c)
		if err != nil {
			c.Set(fiber.HeaderWWWAuthenticate, `Bearer realm="siakad-mini"`)
			return apperror.Unauthorized("Token tidak ditemukan atau format Authorization salah")
		}

		user, err := jwt.Parse(token)
		if err != nil {
			c.Set(fiber.HeaderWWWAuthenticate, `Bearer realm="siakad-mini"`)
			if errors.Is(err, service.ErrExpiredToken) {
				return apperror.Unauthorized("Token sudah kedaluwarsa")
			}
			return apperror.Unauthorized("Token tidak valid")
		}

		c.Locals(localsAuthUser, user)
		return c.Next()
	}
}

// CurrentUser membaca identitas yang disimpan AuthMiddleware.
func CurrentUser(c *fiber.Ctx) (model.AuthUser, bool) {
	user, ok := c.Locals(localsAuthUser).(model.AuthUser)
	return user, ok
}

func bearerToken(c *fiber.Ctx) (string, error) {
	header := c.Get(fiber.HeaderAuthorization)
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("format bukan Bearer")
	}
	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", errors.New("token kosong")
	}
	return token, nil
}