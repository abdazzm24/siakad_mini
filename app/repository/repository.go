// Package repository berisi seluruh akses database (query SQL).
// Setiap repository menerima DBTX sehingga bisa berjalan dalam transaction maupun tidak.
package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound dikembalikan saat baris tidak ditemukan (pgx.ErrNoRows).
var ErrNotFound = errors.New("data tidak ditemukan")

// ErrDuplicate dikembalikan saat terjadi pelanggaran UNIQUE constraint.
var ErrDuplicate = errors.New("data sudah ada")

// DuplicateError membawa nama constraint agar service bisa memberi pesan spesifik.
type DuplicateError struct {
	Constraint string
}

func (e *DuplicateError) Error() string {
	return fmt.Sprintf("duplicate: %s", e.Constraint)
}

func (e *DuplicateError) Is(target error) bool { return target == ErrDuplicate }

// DBTX adalah interface yang dipenuhi oleh *pgxpool.Pool maupun pgx.Tx.
// Semua repository menerima DBTX bukan Pool langsung agar bisa dipakai dalam transaction.
type DBTX interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

// Repositories menyatukan semua repository. Dibuat dari pool untuk operasi biasa,
// atau dari tx untuk operasi dalam transaction.
type Repositories struct {
	Users       UserRepository
	Students    StudentRepository
	Courses     CourseRepository
	Enrollments EnrollmentRepository
}

// NewRepositories membuat Repositories dari pool (tanpa transaction).
func NewRepositories(pool *pgxpool.Pool) *Repositories {
	return &Repositories{
		Users:       NewUserRepository(pool),
		Students:    NewStudentRepository(pool),
		Courses:     NewCourseRepository(pool),
		Enrollments: NewEnrollmentRepository(pool),
	}
}

// newRepositoriesFromTx membuat Repositories dari tx (dalam transaction).
func newRepositoriesFromTx(tx pgx.Tx) *Repositories {
	return &Repositories{
		Users:       NewUserRepository(tx),
		Students:    NewStudentRepository(tx),
		Courses:     NewCourseRepository(tx),
		Enrollments: NewEnrollmentRepository(tx),
	}
}

// TxManager menjalankan fungsi dalam satu transaction; rollback otomatis bila error.
type TxManager struct{ pool *pgxpool.Pool }

func NewTxManager(pool *pgxpool.Pool) *TxManager { return &TxManager{pool: pool} }

func (m *TxManager) Run(ctx context.Context, fn func(*Repositories) error) error {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("memulai transaction: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if err := fn(newRepositoriesFromTx(tx)); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// translateError mengubah error pgx menjadi domain error.
func translateError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// 23505 = unique_violation
		if pgErr.Code == "23505" {
			return &DuplicateError{Constraint: pgErr.ConstraintName}
		}
	}
	return err
}

// escapeLike escapes wildcard characters in a LIKE pattern.
func escapeLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}