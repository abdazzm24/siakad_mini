// Package router adalah peta seluruh endpoint beserta penjaganya. Baca file ini
// untuk tahu endpoint mana yang publik, mana yang butuh login, dan role apa.
package router

import (
	"github.com/gofiber/fiber/v2"

	"siakad-mini/app/handler"
	"siakad-mini/app/middleware"
	"siakad-mini/app/service"
)

type Dependencies struct {
	JWT         *service.JWTManager
	Auth        *handler.AuthHandler
	Students    *handler.StudentHandler
	Courses     *handler.CourseHandler
	Enrollments *handler.EnrollmentHandler
}

func Register(app *fiber.App, d Dependencies) {
	api := app.Group("/api/v1")

	// --- PUBLIK ---
	api.Post("/auth/login", middleware.LoginRateLimiter(), middleware.RenderErrors(), d.Auth.Login)

	// --- WAJIB LOGIN ---
	authed := api.Group("", middleware.AuthMiddleware(d.JWT))

	authed.Get("/auth/me", middleware.AnyRole(), d.Auth.Me)

	// Students
	authed.Get("/students", middleware.AdminOnly(), d.Students.List)
	authed.Post("/students", middleware.AdminOnly(), d.Students.Create)
	// admin ATAU pemilik data: ownership diperiksa di service
	authed.Get("/students/:id", middleware.AnyRole(), d.Students.Show)
	authed.Put("/students/:id", middleware.AdminOnly(), d.Students.Update)
	authed.Delete("/students/:id", middleware.AdminOnly(), d.Students.Delete)

	// Courses
	authed.Get("/courses", middleware.AnyRole(), d.Courses.List)

	// Enrollments (KRS): hanya mahasiswa, ownership diperiksa di service
	authed.Post("/enrollments", middleware.MahasiswaOnly(), d.Enrollments.Create)
	authed.Delete("/enrollments/:id", middleware.MahasiswaOnly(), d.Enrollments.Delete)
}