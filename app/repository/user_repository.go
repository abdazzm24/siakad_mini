package repository

import (
	"context"
	"fmt"

	"siakad-mini/app/model"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (model.User, error)
	FindByID(ctx context.Context, id int64) (model.User, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, u model.User) (model.User, error)
}

type userRepository struct{ db DBTX }

func NewUserRepository(db DBTX) UserRepository { return &userRepository{db: db} }

const userColumns = `id, email, password, role, created_at, updated_at`

func scanUser(row interface{ Scan(dest ...any) error }) (model.User, error) {
	var u model.User
	err := row.Scan(&u.ID, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	return u, err
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (model.User, error) {
	u, err := scanUser(r.db.QueryRow(ctx,
		`SELECT `+userColumns+` FROM users WHERE email = $1`, email))
	if err != nil {
		return model.User{}, translateErrorWrap("mengambil user by email", err)
	}
	return u, nil
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (model.User, error) {
	u, err := scanUser(r.db.QueryRow(ctx,
		`SELECT `+userColumns+` FROM users WHERE id = $1`, id))
	if err != nil {
		return model.User{}, translateErrorWrap("mengambil user", err)
	}
	return u, nil
}

func (r *userRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("memeriksa email: %w", err)
	}
	return exists, nil
}

func (r *userRepository) Create(ctx context.Context, u model.User) (model.User, error) {
	created, err := scanUser(r.db.QueryRow(ctx,
		`INSERT INTO users (email, password, role) VALUES ($1, $2, $3)
		 RETURNING `+userColumns, u.Email, u.Password, u.Role))
	if err != nil {
		return model.User{}, translateErrorWrap("menyimpan user", err)
	}
	return created, nil
}

// translateErrorWrap menerjemahkan error pgx; error lain dibungkus dengan konteks.
func translateErrorWrap(context string, err error) error {
	translated := translateError(err)
	if translated == ErrNotFound {
		return ErrNotFound
	}
	if _, ok := translated.(*DuplicateError); ok {
		return translated
	}
	return fmt.Errorf("%s: %w", context, err)
}