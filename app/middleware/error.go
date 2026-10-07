package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"siakad-mini/app/apperror"
	"siakad-mini/app/response"
)

// ErrorHandler adalah SATU-SATUNYA tempat response error dibentuk.
// Handler dan service cukup mengembalikan error; bentuk JSON-nya seragam di sini.
// Detail teknis (pesan database, stack trace) hanya masuk log, tidak ke client.
func ErrorHandler(logger *slog.Logger, production bool) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			if appErr.Status >= 500 {
				logger.Error("internal_error",
					slog.String("method", c.Method()),
					slog.String("path", c.Path()),
					slog.String("error", appErr.Error()))
			}
			return response.Error(c, appErr.Status, appErr.Message, appErr.Errors)
		}

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			return response.Error(c, fiberErr.Code, safeFiberMessage(fiberErr), nil)
		}

		// Error tak terduga (termasuk panic yang ditangkap recover middleware).
		logger.Error("unhandled_error",
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.String("error", err.Error()))
		message := "Terjadi kesalahan pada server"
		if !production {
			message = "Terjadi kesalahan pada server (lihat log server untuk detail)"
		}
		return response.Error(c, http.StatusInternalServerError, message, nil)
	}
}

func safeFiberMessage(e *fiber.Error) string {
	switch e.Code {
	case fiber.StatusNotFound:
		return "Endpoint tidak ditemukan"
	case fiber.StatusMethodNotAllowed:
		return "Method tidak diizinkan"
	case fiber.StatusRequestEntityTooLarge:
		return "Ukuran request terlalu besar"
	default:
		return e.Message
	}
}