package service

import (
	"context"
	"errors"

	"siakad-mini/app/apperror"
	"siakad-mini/app/model"
	"siakad-mini/app/repository"
	"siakad-mini/app/request"
)

type StudentService struct {
	repos *repository.Repositories
	tx    *repository.TxManager
}

func NewStudentService(repos *repository.Repositories, tx *repository.TxManager) *StudentService {
	return &StudentService{repos: repos, tx: tx}
}

func (s *StudentService) List(ctx context.Context, q model.StudentListQuery) ([]model.Student, int, error) {
	students, total, err := s.repos.Students.List(ctx, q)
	if err != nil {
		return nil, 0, apperror.Internal(err)
	}
	return students, total, nil
}

// Create membuat user (role mahasiswa) dan student dalam SATU transaction.
// Password awal = NIM yang di-hash. Bila salah satu langkah gagal: ROLLBACK.
func (s *StudentService) Create(ctx context.Context, req request.CreateStudentRequest) (model.Student, error) {
	// bcrypt sengaja dihitung di luar transaction agar koneksi tidak tertahan lama.
	hash, err := HashPassword(req.NIM)
	if err != nil {
		return model.Student{}, apperror.Internal(err)
	}

	ipk := 0.0
	if req.IPKTerakhir != nil {
		ipk = *req.IPKTerakhir
	}

	var created model.Student
	err = s.tx.Run(ctx, func(r *repository.Repositories) error {
		// Pemeriksaan awal agar pesan error lengkap (NIM dan email sekaligus).
		errs := map[string][]string{}
		if exists, err := r.Students.NIMExists(ctx, req.NIM); err != nil {
			return apperror.Internal(err)
		} else if exists {
			errs["nim"] = []string{"NIM sudah terdaftar"}
		}
		if exists, err := r.Users.EmailExists(ctx, req.Email); err != nil {
			return apperror.Internal(err)
		} else if exists {
			errs["email"] = []string{"Email sudah terdaftar"}
		}
		if len(errs) > 0 {
			return apperror.Validation(errs)
		}

		user, err := r.Users.Create(ctx, model.User{
			Email: req.Email, Password: hash, Role: model.RoleMahasiswa,
		})
		if err != nil {
			return mapStudentWriteError(err)
		}

		created, err = r.Students.Create(ctx, model.Student{
			UserID: user.ID, NIM: req.NIM, Nama: req.Nama, Prodi: req.Prodi,
			Angkatan: req.Angkatan, IPKTerakhir: ipk,
		})
		if err != nil {
			return mapStudentWriteError(err)
		}
		return nil
	})
	if err != nil {
		return model.Student{}, asAppError(err)
	}
	return created, nil
}

// mapStudentWriteError menangani race condition: dua request bersamaan yang lolos
// pemeriksaan awal tetap ditolak database lewat UNIQUE constraint.
func mapStudentWriteError(err error) error {
	var dup *repository.DuplicateError
	if errors.As(err, &dup) {
		switch dup.Constraint {
		case "students_nim_key":
			return apperror.FieldError("nim", "NIM sudah terdaftar")
		case "users_email_key":
			return apperror.FieldError("email", "Email sudah terdaftar")
		}
	}
	return apperror.Internal(err)
}

// GetDetail mengembalikan mahasiswa beserta KRS, total SKS, dan batas SKS.
// Admin boleh melihat siapa pun; mahasiswa hanya data miliknya sendiri.
func (s *StudentService) GetDetail(ctx context.Context, au model.AuthUser, id int64, tahunAkademik string) (*model.StudentDetail, error) {
	student, err := s.repos.Students.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperror.NotFound("Mahasiswa tidak ditemukan")
		}
		return nil, apperror.Internal(err)
	}

	// Ownership: student.user_id harus sama dengan user yang sedang login.
	if au.Role != model.RoleAdmin && student.UserID != au.UserID {
		return nil, apperror.Forbidden("Anda tidak berhak mengakses data mahasiswa lain")
	}

	courses, err := s.repos.Enrollments.ListByStudent(ctx, student.ID, tahunAkademik)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	total := 0
	for _, c := range courses {
		total += c.SKS
	}

	return &model.StudentDetail{
		Student:  student,
		Courses:  courses,
		TotalSKS: total,
		BatasSKS: MaxSKS(student.IPKTerakhir),
	}, nil
}

func (s *StudentService) Update(ctx context.Context, id int64, req request.UpdateStudentRequest) (model.Student, error) {
	student, err := s.repos.Students.Update(ctx, id, model.StudentUpdate{
		Nama: req.Nama, Prodi: req.Prodi, Angkatan: req.Angkatan, IPKTerakhir: req.IPKTerakhir,
	})
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return model.Student{}, apperror.NotFound("Mahasiswa tidak ditemukan")
		}
		return model.Student{}, apperror.Internal(err)
	}
	return student, nil
}

func (s *StudentService) Delete(ctx context.Context, id int64) error {
	if err := s.repos.Students.SoftDelete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return apperror.NotFound("Mahasiswa tidak ditemukan")
		}
		return apperror.Internal(err)
	}
	return nil
}

// asAppError memastikan error dari transaction selalu berupa *AppError.
func asAppError(err error) error {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return apperror.Internal(err)
}