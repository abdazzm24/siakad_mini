package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgxpool"

	"siakad-mini/app/config"
	"siakad-mini/app/handler"
	"siakad-mini/app/middleware"
	"siakad-mini/app/repository"
	"siakad-mini/app/router"
	"siakad-mini/app/service"
	"siakad-mini/database"
)

// Cara pakai:
//
//	go run . migrate   -> membuat tabel
//	go run . seed      -> mengisi data awal
//	go run .           -> menjalankan server (sama dengan: go run . serve)
func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		log.Error("konfigurasi tidak valid", slog.String("error", err.Error()))
		os.Exit(1)
	}

	command := "serve"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	if command != "migrate" && command != "seed" && command != "serve" {
		fmt.Println("perintah tidak dikenal. Gunakan: migrate | seed | serve")
		os.Exit(1)
	}

	ctx := context.Background()
	pool, err := config.NewPool(ctx, cfg)
	if err != nil {
		log.Error("database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	switch command {
	case "migrate":
		if err := database.Migrate(ctx, pool, log); err != nil {
			log.Error("migrate gagal", slog.String("error", err.Error()))
			os.Exit(1)
		}
	case "seed":
		if err := database.Seed(ctx, pool, log); err != nil {
			log.Error("seed gagal", slog.String("error", err.Error()))
			os.Exit(1)
		}
	default:
		serve(cfg, pool, log)
	}
}

// serve merakit dependency dari dalam ke luar, lalu menjalankan server
// dengan graceful shutdown.
func serve(cfg *config.Config, pool *pgxpool.Pool, log *slog.Logger) {
	// repository -> service -> handler
	repos := repository.NewRepositories(pool)
	txManager := repository.NewTxManager(pool)

	jwtManager := service.NewJWTManager(cfg.JWTSecret, time.Duration(cfg.JWTExpireMinutes)*time.Minute)

	deps := router.Dependencies{
		JWT:         jwtManager,
		Auth:        handler.NewAuthHandler(service.NewAuthService(repos, jwtManager)),
		Students:    handler.NewStudentHandler(service.NewStudentService(repos, txManager)),
		Courses:     handler.NewCourseHandler(service.NewCourseService(repos)),
		Enrollments: handler.NewEnrollmentHandler(service.NewEnrollmentService(repos, txManager)),
	}

	app := fiber.New(fiber.Config{
		AppName:      "SIAKAD Mini",
		ErrorHandler: middleware.ErrorHandler(log, cfg.IsProduction()),
		BodyLimit:    1 * 1024 * 1024, // 1 MB
	})
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip} ${method} ${path} ${status} ${latency}\n",
	}))
	router.Register(app, deps)

	// Jalur yang tidak dikenal -> 404 lewat ErrorHandler (format JSON seragam).
	app.Use(func(c *fiber.Ctx) error { return fiber.ErrNotFound })

	go func() {
		log.Info("server berjalan", slog.String("url", "http://localhost:"+cfg.AppPort))
		if err := app.Listen(":" + cfg.AppPort); err != nil {
			log.Error("server berhenti", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("sinyal berhenti diterima, menutup server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Error("gagal menutup server dengan rapi", slog.String("error", err.Error()))
	}
}