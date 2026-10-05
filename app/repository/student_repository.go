package repository

import (
	"context"
	"fmt"
	"strings"

	"siakad-mini/app/model"
)

type StudentRepository interface {
	Create(ctx context.Context, s model.Student) (model.Student, error)
	FindByID(ctx context.Context, id int64) (model.Student, error)
	FindByUserID(ctx context.Context, userID int64) (model.Student, error)
	// FindByUserIDForUpdate mengunci baris mahasiswa (SELECT ... FOR UPDATE).
	FindByUserIDForUpdate(ctx context.Context, userID int64) (model.Student, error)
	NIMExists(ctx context.Context, nim string) (bool, error)
	List(ctx context.Context, q model.StudentListQuery) ([]model.Student, int, error)
	Update(ctx context.Context, id int64, u model.StudentUpdate) (model.Student, error)
	SoftDelete(ctx context.Context, id int64) error
}

type studentRepository struct{ db DBTX }

func NewStudentRepository(db DBTX) StudentRepository { return &studentRepository{db: db} }

// Semua query mahasiswa memakai s.deleted_at IS NULL (soft delete).
const studentSelect = `
	SELECT s.id, s.user_id, s.nim, s.nama, s.prodi, s.angkatan,
	       s.ipk_terakhir::float8, u.email, s.deleted_at, s.created_at, s.updated_at
	FROM students s
	JOIN users u ON u.id = s.user_id`

func scanStudent(row interface{ Scan(dest ...any) error }) (model.Student, error) {
	var s model.Student
	err := row.Scan(&s.ID, &s.UserID, &s.NIM, &s.Nama, &s.Prodi, &s.Angkatan,
		&s.IPKTerakhir, &s.Email, &s.DeletedAt, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func (r *studentRepository) findOne(ctx context.Context, where, suffix string, arg any) (model.Student, error) {
	s, err := scanStudent(r.db.QueryRow(ctx,
		studentSelect+" WHERE "+where+" AND s.deleted_at IS NULL"+suffix, arg))
	if err != nil {
		return model.Student{}, translateErrorWrap("mengambil mahasiswa", err)
	}
	return s, nil
}

func (r *studentRepository) FindByID(ctx context.Context, id int64) (model.Student, error) {
	return r.findOne(ctx, "s.id = $1", "", id)
}

func (r *studentRepository) FindByUserID(ctx context.Context, userID int64) (model.Student, error) {
	return r.findOne(ctx, "s.user_id = $1", "", userID)
}

func (r *studentRepository) FindByUserIDForUpdate(ctx context.Context, userID int64) (model.Student, error) {
	return r.findOne(ctx, "s.user_id = $1", " FOR UPDATE OF s", userID)
}

func (r *studentRepository) NIMExists(ctx context.Context, nim string) (bool, error) {
	// Baris yang sudah di-soft delete tetap menghitung: NIM tidak boleh dipakai ulang.
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM students WHERE nim = $1)`, nim).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("memeriksa nim: %w", err)
	}
	return exists, nil
}

func (r *studentRepository) Create(ctx context.Context, s model.Student) (model.Student, error) {
	var id int64
	err := r.db.QueryRow(ctx,
		`INSERT INTO students (user_id, nim, nama, prodi, angkatan, ipk_terakhir)
		 VALUES ($1, $2, $3, $4, $5, ROUND($6::float8::numeric, 2))
		 RETURNING id`,
		s.UserID, s.NIM, s.Nama, s.Prodi, s.Angkatan, s.IPKTerakhir).Scan(&id)
	if err != nil {
		return model.Student{}, translateErrorWrap("menyimpan mahasiswa", err)
	}
	return r.FindByID(ctx, id)
}

// sortClauses adalah daftar putih ORDER BY. Nilai dari client hanya dipakai
// sebagai KUNCI pencarian, tidak pernah disambung ke teks SQL.
var sortClauses = map[string]string{
	"nama":          "s.nama ASC",
	"-nama":         "s.nama DESC",
	"nim":           "s.nim ASC",
	"-nim":          "s.nim DESC",
	"ipk_terakhir":  "s.ipk_terakhir ASC",
	"-ipk_terakhir": "s.ipk_terakhir DESC",
	"angkatan":      "s.angkatan ASC",
	"-angkatan":     "s.angkatan DESC",
}

func buildStudentFilter(q model.StudentListQuery) (string, []any) {
	conds := []string{"s.deleted_at IS NULL"}
	args := []any{}

	if q.Search != "" {
		args = append(args, "%"+escapeLike(q.Search)+"%")
		n := len(args)
		conds = append(conds, fmt.Sprintf("(s.nim ILIKE $%d OR s.nama ILIKE $%d)", n, n))
	}
	if q.Prodi != "" {
		args = append(args, q.Prodi)
		conds = append(conds, fmt.Sprintf("LOWER(s.prodi) = LOWER($%d)", len(args)))
	}
	if q.Angkatan != nil {
		args = append(args, *q.Angkatan)
		conds = append(conds, fmt.Sprintf("s.angkatan = $%d", len(args)))
	}
	return " WHERE " + strings.Join(conds, " AND "), args
}

func (r *studentRepository) List(ctx context.Context, q model.StudentListQuery) ([]model.Student, int, error) {
	where, args := buildStudentFilter(q)

	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM students s`+where, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("menghitung mahasiswa: %w", err)
	}

	order := "s.id ASC"
	if clause, ok := sortClauses[q.Sort]; ok {
		order = clause + ", s.id ASC"
	}

	args = append(args, q.PerPage, q.Offset())
	sqlText := fmt.Sprintf(`%s%s ORDER BY %s LIMIT $%d OFFSET $%d`,
		studentSelect, where, order, len(args)-1, len(args))

	rows, err := r.db.Query(ctx, sqlText, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("mengambil daftar mahasiswa: %w", err)
	}
	defer rows.Close()

	students := []model.Student{}
	for rows.Next() {
		s, err := scanStudent(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("membaca baris mahasiswa: %w", err)
		}
		students = append(students, s)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("membaca hasil query mahasiswa: %w", err)
	}
	return students, total, nil
}

func (r *studentRepository) Update(ctx context.Context, id int64, u model.StudentUpdate) (model.Student, error) {
	var updatedID int64
	err := r.db.QueryRow(ctx,
		`UPDATE students SET
			nama         = COALESCE($2, nama),
			prodi        = COALESCE($3, prodi),
			angkatan     = COALESCE($4, angkatan),
			ipk_terakhir = COALESCE(ROUND($5::float8::numeric, 2), ipk_terakhir),
			updated_at   = NOW()
		 WHERE id = $1 AND deleted_at IS NULL
		 RETURNING id`,
		id, u.Nama, u.Prodi, u.Angkatan, u.IPKTerakhir).Scan(&updatedID)
	if err != nil {
		return model.Student{}, translateErrorWrap("memperbarui mahasiswa", err)
	}
	return r.FindByID(ctx, updatedID)
}

func (r *studentRepository) SoftDelete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE students SET deleted_at = NOW(), updated_at = NOW()
		 WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return fmt.Errorf("soft delete mahasiswa: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}