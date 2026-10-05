package repository

import (
	"context"
	"fmt"
	"strings"

	"siakad-mini/app/model"
)

type CourseRepository interface {
	List(ctx context.Context, f model.CourseFilter) ([]model.Course, error)
	FindByID(ctx context.Context, id int64) (model.Course, error)
	// FindByIDForUpdate mengunci baris mata kuliah (SELECT ... FOR UPDATE)
	// agar pengecekan kuota aman dari race condition.
	FindByIDForUpdate(ctx context.Context, id int64) (model.Course, error)
	CountEnrolled(ctx context.Context, courseID int64) (int, error)
}

type courseRepository struct{ db DBTX }

func NewCourseRepository(db DBTX) CourseRepository { return &courseRepository{db: db} }

// terisi dihitung dari enrollments. Enrollment milik mahasiswa yang sudah
// di-soft delete tidak ikut dihitung sehingga kursinya kembali tersedia.
const enrolledCount = `(
	SELECT COUNT(*) FROM enrollments e
	JOIN students s ON s.id = e.student_id
	WHERE e.course_id = c.id AND s.deleted_at IS NULL
)::int`

func (r *courseRepository) List(ctx context.Context, f model.CourseFilter) ([]model.Course, error) {
	conds := []string{"1 = 1"}
	args := []any{}

	if f.Semester != nil {
		args = append(args, *f.Semester)
		conds = append(conds, fmt.Sprintf("t.semester = $%d", len(args)))
	}
	if f.Search != "" {
		args = append(args, "%"+escapeLike(f.Search)+"%")
		n := len(args)
		conds = append(conds, fmt.Sprintf("(t.kode_mk ILIKE $%d OR t.nama_mk ILIKE $%d)", n, n))
	}
	if f.AvailableOnly {
		conds = append(conds, "t.terisi < t.kuota")
	}

	sqlText := fmt.Sprintf(`
		SELECT t.id, t.kode_mk, t.nama_mk, t.sks, t.semester, t.kuota, t.terisi
		FROM (
			SELECT c.id, c.kode_mk, c.nama_mk, c.sks, c.semester, c.kuota, %s AS terisi
			FROM courses c
		) t
		WHERE %s
		ORDER BY t.semester ASC, t.kode_mk ASC`, enrolledCount, strings.Join(conds, " AND "))

	rows, err := r.db.Query(ctx, sqlText, args...)
	if err != nil {
		return nil, fmt.Errorf("mengambil daftar mata kuliah: %w", err)
	}
	defer rows.Close()

	courses := []model.Course{}
	for rows.Next() {
		var c model.Course
		if err := rows.Scan(&c.ID, &c.KodeMK, &c.NamaMK, &c.SKS, &c.Semester, &c.Kuota, &c.Terisi); err != nil {
			return nil, fmt.Errorf("membaca baris mata kuliah: %w", err)
		}
		c.SisaKuota = max(c.Kuota-c.Terisi, 0)
		courses = append(courses, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("membaca hasil query mata kuliah: %w", err)
	}
	return courses, nil
}

func (r *courseRepository) find(ctx context.Context, id int64, suffix string) (model.Course, error) {
	var c model.Course
	err := r.db.QueryRow(ctx,
		`SELECT c.id, c.kode_mk, c.nama_mk, c.sks, c.semester, c.kuota
		 FROM courses c WHERE c.id = $1`+suffix, id,
	).Scan(&c.ID, &c.KodeMK, &c.NamaMK, &c.SKS, &c.Semester, &c.Kuota)
	if err != nil {
		return model.Course{}, translateErrorWrap("mengambil mata kuliah", err)
	}
	return c, nil
}

func (r *courseRepository) FindByID(ctx context.Context, id int64) (model.Course, error) {
	return r.find(ctx, id, "")
}

func (r *courseRepository) FindByIDForUpdate(ctx context.Context, id int64) (model.Course, error) {
	return r.find(ctx, id, " FOR UPDATE")
}

func (r *courseRepository) CountEnrolled(ctx context.Context, courseID int64) (int, error) {
	var n int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM enrollments e
		 JOIN students s ON s.id = e.student_id
		 WHERE e.course_id = $1 AND s.deleted_at IS NULL`, courseID).Scan(&n)
	if err != nil {
		return 0, fmt.Errorf("menghitung peserta mata kuliah: %w", err)
	}
	return n, nil
}