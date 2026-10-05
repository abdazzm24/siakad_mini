package service

import (
	"context"

	"siakad-mini/app/apperror"
	"siakad-mini/app/model"
	"siakad-mini/app/repository"
)

type CourseService struct {
	repos *repository.Repositories
}

func NewCourseService(repos *repository.Repositories) *CourseService {
	return &CourseService{repos: repos}
}

func (s *CourseService) List(ctx context.Context, f model.CourseFilter) ([]model.Course, error) {
	courses, err := s.repos.Courses.List(ctx, f)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return courses, nil
}