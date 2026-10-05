package repository

import (
	"context"
	"fmt"

	"siakad-mini/app/model"
)

type EnrollmentRepository interface {
	Create(ctx context.Context, e model.Enrollment) (model.Enrollment, error)
	Exists(ctx context.Context, studentID, courseID int64, tahunAkademik string) (bool, error)
	FindByID(ctx context.Context, id int64) (model.Enrollment, error)
	Delete(ctx context.Context, id int64) error
	SumSKS(ctx context.Context, studentID int64, tahunAkademik string) (int, error)
	// ListByStudent mengembalikan KRS mahasiswa. tahunAkademik kosong = semua.
	ListByStudent(ctx context.Context, studentID int64, tahunAkademik string) ([]model.EnrolledCourse, error)
}

type enrollmentRepository struct{ db DBTX }

func NewEnrollmentRepository(db DBTX) EnrollmentRepository { return &enrollmentRepository{db: db} }

func (r *enrollmentRepository) Create(ctx context.Context, e model.Enrollment) (model.Enrollment, error) {
	err := r.db.QueryRow(ctx,
		`INSERT INTO enrollments (student_id, course_id, tahun_akademik)
		 VALUES ($1, $2, $3)
		 RETURNING id, created_at`,
		e.StudentID, e.CourseID, e.TahunAkademik).Scan(&e.ID, &e.CreatedAt)
	if err != nil {
		return model.Enrollment{}, translateErrorWrap("menyimpan enrollment", err)
	}
	return e, nil
}

func (r *enrollmentRepository) Exists(ctx context.Context, studentID, courseID int64, tahunAkademik string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS (
			SELECT 1 FROM enrollments
			WHERE student_id = $1 AND course_id = $2 AND tahun_akademik = $3
		)`, studentID, courseID, tahunAkademik).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("memeriksa duplikasi enrollment: %w", err)
	}
	return exists, nil
}

func (r *enrollmentRepository) FindByID(ctx context.Context, id int64) (model.Enrollment, error) {
	var e model.Enrollment
	err := r.db.QueryRow(ctx,
		`SELECT id, student_id, course_id, tahun_akademik, created_at
		 FROM enrollments WHERE id = $1`, id,
	).Scan(&e.ID, &e.StudentID, &e.CourseID, &e.TahunAkademik, &e.CreatedAt)
	if err != nil {
		return model.Enrollment{}, translateErrorWrap("mengambil enrollment", err)
	}
	return e, nil
}

func (r *enrollmentRepository) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM enrollments WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("menghapus enrollment: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *enrollmentRepository) SumSKS(ctx context.Context, studentID int64, tahunAkademik string) (int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(c.sks), 0)::int
		 FROM enrollments e
		 JOIN courses c ON c.id = e.course_id
		 WHERE e.student_id = $1 AND e.tahun_akademik = $2`,
		studentID, tahunAkademik).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("menghitung total sks: %w", err)
	}
	return total, nil
}

func (r *enrollmentRepository) ListByStudent(ctx context.Context, studentID int64, tahunAkademik string) ([]model.EnrolledCourse, error) {
	rows, err := r.db.Query(ctx,
		`SELECT e.id, e.tahun_akademik, c.id, c.kode_mk, c.nama_mk, c.sks, c.semester
		 FROM enrollments e
		 JOIN courses c ON c.id = e.course_id
		 WHERE e.student_id = $1 AND ($2::text = '' OR e.tahun_akademik = $2)
		 ORDER BY e.tahun_akademik, c.semester, c.kode_mk`,
		studentID, tahunAkademik)
	if err != nil {
		return nil, fmt.Errorf("mengambil KRS mahasiswa: %w", err)
	}
	defer rows.Close()

	courses := []model.EnrolledCourse{}
	for rows.Next() {
		var c model.EnrolledCourse
		if err := rows.Scan(&c.EnrollmentID, &c.TahunAkademik, &c.CourseID,
			&c.KodeMK, &c.NamaMK, &c.SKS, &c.Semester); err != nil {
			return nil, fmt.Errorf("membaca baris KRS: %w", err)
		}
		courses = append(courses, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("membaca hasil query KRS: %w", err)
	}
	return courses, nil
}