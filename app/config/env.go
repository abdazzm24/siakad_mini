package config

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config menampung seluruh konfigurasi aplikasi yang dibaca dari environment.
type Config struct {
	AppPort string
	AppEnv  string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret        string
	JWTExpireMinutes int
}

// IsProduction dipakai untuk menyembunyikan detail error teknis dari client.
func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

// Load membaca file .env (bila ada) lalu menyusun Config.
// Aplikasi menolak menyala bila JWT_SECRET kosong atau terlalu pendek.
func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("peringatan: file .env tidak ditemukan, memakai environment sistem")
	}

	cfg := &Config{
		AppPort:          getEnv("APP_PORT", "3000"),
		AppEnv:           getEnv("APP_ENV", "development"),
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPassword:       getEnv("DB_PASSWORD", ""),
		DBName:           getEnv("DB_NAME", "siakad_mini"),
		DBSSLMode:        getEnv("DB_SSLMODE", "disable"),
		JWTSecret:        getEnv("JWT_SECRET", ""),
		JWTExpireMinutes: getEnvInt("JWT_EXPIRE_MINUTES", 15),
	}

	if len(cfg.JWTSecret) < 16 {
		return nil, errors.New("JWT_SECRET wajib diisi minimal 16 karakter")
	}
	if cfg.JWTExpireMinutes < 1 {
		return nil, errors.New("JWT_EXPIRE_MINUTES harus lebih dari 0")
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("peringatan: %s bukan angka (%q), memakai bawaan %d", key, value, fallback)
		return fallback
	}
	return parsed
}