package service

import (
	"context"
	"errors"
	"fmt"

	"siakad-mini/app/apperror"
	"siakad-mini/app/model"
	"siakad-mini/app/repository"
	"siakad-mini/app/request"
)

type EnrollmentService struct {
	repos *repository.Repositories
	tx    *repository.TxManager
}

func NewEnrollmentService(repos *repository.Repositories, tx *repository.TxManager) *EnrollmentService {
	return &EnrollmentService{repos: repos, tx: tx}
}

// Create mengambil mata kuliah (menambah KRS) dalam satu transaction:
//  1. kunci & ambil student      4. cek duplikasi
//  2. kunci & ambil course       5. cek kuota
//  3. (row locking FOR UPDATE)   6. cek batas SKS lalu INSERT
//
// Baris student dan course dikunci dengan SELECT ... FOR UPDATE sehingga dua
// request bersamaan tidak bisa sama-sama lolos cek kuota maupun cek SKS.
// Urutan kunci selalu student -> course agar tidak terjadi deadlock.
func (s *EnrollmentService) Create(ctx context.Context, au model.AuthUser, req request.CreateEnrollmentRequest) (*model.EnrollmentResult, error) {
	var result *model.EnrollmentResult

	err := s.tx.Run(ctx, func(r *repository.Repositories) error {
		student, err := r.Students.FindByUserIDForUpdate(ctx, au.UserID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return apperror.Unauthorized("Akun mahasiswa tidak aktif")
			}
			return apperror.Internal(err)
		}

		course, err := r.Courses.FindByIDForUpdate(ctx, req.CourseID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return apperror.FieldError("course_id", "Mata kuliah tidak ditemukan")
			}
			return apperror.Internal(err)
		}

		exists, err := r.Enrollments.Exists(ctx, student.ID, course.ID, req.TahunAkademik)
		if err != nil {
			return apperror.Internal(err)
		}
		if exists {
			return apperror.Conflict("Mata kuliah sudah diambil pada tahun akademik tersebut")
		}

		terisi, err := r.Courses.CountEnrolled(ctx, course.ID)
		if err != nil {
			return apperror.Internal(err)
		}
		if IsCourseFull(terisi, course.Kuota) {
			return apperror.Unprocessable("Kuota mata kuliah sudah penuh")
		}

		current, err := r.Enrollments.SumSKS(ctx, student.ID, req.TahunAkademik)
		if err != nil {
			return apperror.Internal(err)
		}
		limit := MaxSKS(student.IPKTerakhir)
		if ExceedsSKSLimit(limit, current, course.SKS) {
			return apperror.Unprocessable(fmt.Sprintf(
				"Batas SKS terlampaui. Sisa SKS yang dapat diambil: %d",
				RemainingSKS(limit, current)))
		}

		enrollment, err := r.Enrollments.Create(ctx, model.Enrollment{
			StudentID: student.ID, CourseID: course.ID, TahunAkademik: req.TahunAkademik,
		})
		if err != nil {
			if errors.Is(err, repository.ErrDuplicate) {
				return apperror.Conflict("Mata kuliah sudah diambil pada tahun akademik tersebut")
			}
			return apperror.Internal(err)
		}

		total := current + course.SKS
		result = &model.EnrollmentResult{
			Enrollment: enrollment,
			Course: model.EnrolledCourse{
				EnrollmentID: enrollment.ID, TahunAkademik: enrollment.TahunAkademik,
				CourseID: course.ID, KodeMK: course.KodeMK, NamaMK: course.NamaMK,
				SKS: course.SKS, Semester: course.Semester,
			},
			TotalSKS: total,
			BatasSKS: limit,
			SisaSKS:  RemainingSKS(limit, total),
		}
		return nil
	})
	if err != nil {
		return nil, asAppError(err)
	}
	return result, nil
}

// Delete membatalkan mata kuliah dari KRS. Hanya pemilik enrollment yang boleh.
// Kuota otomatis kembali karena "terisi" dihitung dari jumlah baris enrollments.
func (s *EnrollmentService) Delete(ctx context.Context, au model.AuthUser, id int64) error {
	student, err := s.repos.Students.FindByUserID(ctx, au.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apperror.Unauthorized("Akun mahasiswa tidak aktif")
		}
		return apperror.Internal(err)
	}

	enrollment, err := s.repos.Enrollments.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apperror.NotFound("Enrollment tidak ditemukan")
		}
		return apperror.Internal(err)
	}

	// Ownership: enrollment harus milik mahasiswa yang sedang login.
	if enrollment.StudentID != student.ID {
		return apperror.Forbidden("Anda tidak berhak mengubah KRS mahasiswa lain")
	}

	if err := s.repos.Enrollments.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apperror.NotFound("Enrollment tidak ditemukan")
		}
		return apperror.Internal(err)
	}
	return nil
}