package service

import (
	"context"
	"github.com/zhuravlev-pe/course-watch/internal/core/domain"
	"github.com/zhuravlev-pe/course-watch/internal/core/dto"
)

type CoursesService struct {
	repo  CoursesRepository
	idGen IDGenerator
}

func NewCoursesService(repo CoursesRepository, idGen IDGenerator) *CoursesService {
	return &CoursesService{
		repo:  repo,
		idGen: idGen,
	}
}

func (s *CoursesService) GetById(ctx context.Context, id string) (*domain.Course, error) {
	return s.repo.GetById(ctx, id)
}

func (s *CoursesService) Create(ctx context.Context, input dto.CreateCourseInput) (*domain.Course, error) {
	course := &domain.Course{
		Id:          s.idGen.Generate(),
		Title:       input.Title,
		Description: input.Description,
	}
	if err := s.repo.Insert(ctx, course); err != nil {
		return nil, err
	}
	return course, nil
}
