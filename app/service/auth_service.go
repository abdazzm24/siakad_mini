package service

import (
	"context"
	"errors"

	"siakad-mini/app/apperror"
	"siakad-mini/app/model"
	"siakad-mini/app/repository"
	"siakad-mini/app/request"
)

const msgLoginFailed = "Email atau password salah"

type AuthService struct {
	repos *repository.Repositories
	jwt   *JWTManager
}

func NewAuthService(repos *repository.Repositories, jwt *JWTManager) *AuthService {
	return &AuthService{repos: repos, jwt: jwt}
}

type LoginResult struct {
	AccessToken string
	TokenType   string
	ExpiresIn   int
	User        model.User
}

func (s *AuthService) Login(ctx context.Context, req request.LoginRequest) (*LoginResult, error) {
	user, err := s.repos.Users.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			verifyDummyPassword(req.Password)
			return nil, apperror.Unauthorized(msgLoginFailed)
		}
		return nil, apperror.Internal(err)
	}
	if !VerifyPassword(user.Password, req.Password) {
		return nil, apperror.Unauthorized(msgLoginFailed)
	}

	// Mahasiswa yang sudah di-soft delete tidak boleh login.
	if user.Role == model.RoleMahasiswa {
		if _, err := s.repos.Students.FindByUserID(ctx, user.ID); err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, apperror.Unauthorized(msgLoginFailed)
			}
			return nil, apperror.Internal(err)
		}
	}

	token, err := s.jwt.Generate(user)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return &LoginResult{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   s.jwt.ExpiresInSeconds(),
		User:        user,
	}, nil
}

func (s *AuthService) Me(ctx context.Context, au model.AuthUser) (*model.MeResponse, error) {
	user, err := s.repos.Users.FindByID(ctx, au.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperror.Unauthorized("Akun tidak ditemukan")
		}
		return nil, apperror.Internal(err)
	}

	me := &model.MeResponse{ID: user.ID, Email: user.Email, Role: user.Role}
	if user.Role == model.RoleMahasiswa {
		st, err := s.repos.Students.FindByUserID(ctx, user.ID)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				return nil, apperror.Unauthorized("Akun mahasiswa tidak aktif")
			}
			return nil, apperror.Internal(err)
		}
		me.Student = &model.StudentBrief{
			ID: st.ID, NIM: st.NIM, Nama: st.Nama, Prodi: st.Prodi, Angkatan: st.Angkatan,
		}
	}
	return me, nil
}