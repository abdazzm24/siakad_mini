// Package database berisi migration runner dan seeder sederhana.
// File SQL di-embed ke dalam binary sehingga `go run .` tidak bergantung pada
// lokasi folder saat dijalankan.
package database

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

//go:embed seeders/*.sql
var seederFiles embed.FS

// Migrate menjalankan seluruh file migrations/*.sql secara berurutan.
// Migration yang sudah pernah dijalankan dicatat di tabel schema_migrations
// dan dilewati, sehingga aman dijalankan berulang kali.
func Migrate(ctx context.Context, pool *pgxpool.Pool, logger *slog.Logger) error {
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version    VARCHAR(100) PRIMARY KEY,
			applied_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
		)`)
	if err != nil {
		return fmt.Errorf("membuat schema_migrations: %w", err)
	}

	names, err := sortedNames(migrationFiles, "migrations")
	if err != nil {
		return err
	}

	for _, name := range names {
		var applied bool
		if err := pool.QueryRow(ctx,
			`SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version = $1)`, name,
		).Scan(&applied); err != nil {
			return fmt.Errorf("memeriksa migration %s: %w", name, err)
		}
		if applied {
			logger.Info("migration dilewati (sudah dijalankan)", slog.String("file", name))
			continue
		}

		content, err := migrationFiles.ReadFile("migrations/" + name)
		if err != nil {
			return fmt.Errorf("membaca %s: %w", name, err)
		}
		if err := runInTx(ctx, pool, string(content), name); err != nil {
			return err
		}
		logger.Info("migration berhasil", slog.String("file", name))
	}
	return nil
}

// runInTx menjalankan satu file SQL beserta pencatatannya dalam satu transaction.
func runInTx(ctx context.Context, pool *pgxpool.Pool, sqlText, name string) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("memulai transaction %s: %w", name, err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, sqlText); err != nil {
		return fmt.Errorf("menjalankan %s: %w", name, err)
	}
	if _, err := tx.Exec(ctx,
		`INSERT INTO schema_migrations (version) VALUES ($1)`, name); err != nil {
		return fmt.Errorf("mencatat %s: %w", name, err)
	}
	return tx.Commit(ctx)
}

// Seed menjalankan seeders/seed.sql. Seeder bersifat idempotent (ON CONFLICT DO
// NOTHING) sehingga aman dijalankan berulang kali.
func Seed(ctx context.Context, pool *pgxpool.Pool, logger *slog.Logger) error {
	content, err := seederFiles.ReadFile("seeders/seed.sql")
	if err != nil {
		return fmt.Errorf("membaca seed.sql: %w", err)
	}
	if err := runSeed(ctx, pool, string(content)); err != nil {
		return err
	}
	logger.Info("seeder berhasil dijalankan")
	return nil
}

func runSeed(ctx context.Context, pool *pgxpool.Pool, sqlText string) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("memulai transaction seed: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if _, err := tx.Exec(ctx, sqlText); err != nil {
		return fmt.Errorf("menjalankan seed.sql: %w", err)
	}
	return tx.Commit(ctx)
}

func sortedNames(fsys embed.FS, dir string) ([]string, error) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, fmt.Errorf("membaca folder %s: %w", dir, err)
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}