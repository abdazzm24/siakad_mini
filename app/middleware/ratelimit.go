package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	"siakad-mini/app/apperror"
)

// LoginRateLimiter membatasi percobaan login GAGAL: maksimal 5 per menit per IP.
// Percobaan yang berhasil tidak dihitung (SkipSuccessfulRequests), sehingga
// percobaan ke-6 yang gagal dalam satu menit dijawab 429.
func LoginRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:                    5,
		Expiration:             1 * time.Minute,
		SkipSuccessfulRequests: true,
		KeyGenerator:           func(c *fiber.Ctx) string { return c.IP() },
		LimitReached: func(c *fiber.Ctx) error {
			c.Set(fiber.HeaderRetryAfter, "60")
			return apperror.TooManyRequests("Terlalu banyak percobaan login. Silakan coba lagi nanti.")
		},
	})
}

// RenderErrors mengubah error handler menjadi response SAAT ITU JUGA.
//
// Limiter Fiber memeriksa status response sesudah handler selesai. Tanpa ini,
// error yang dikembalikan sebagai `error` baru diubah menjadi JSON oleh
// ErrorHandler SETELAH limiter memeriksa, sehingga login gagal terbaca "200 OK"
// dan tidak pernah dihitung. Pasang middleware ini tepat di bawah limiter.
func RenderErrors() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := c.Next(); err != nil {
			return c.App().ErrorHandler(c, err)
		}
		return nil
	}
}