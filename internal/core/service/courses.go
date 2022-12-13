package service

import (
	"context"
	"github.com/zhuravlev-pe/course-watch/internal/core"
)

type CourseService struct {
	repo  CoursesRepository
	idGen IdGenerator
}

func NewCoursesService(repo CoursesRepository, idGen IdGenerator) *CourseService {
	return &CourseService{
		repo:  repo,
		idGen: idGen,
	}
}

func (s *CourseService) GetById(ctx context.Context, id string) (*core.Course, error) {
	return s.repo.GetById(ctx, id)
}

func (s *CourseService) Create(ctx context.Context, input core.CreateCourseInput) (*core.Course, error) {
	course := &core.Course{
		Id:          s.idGen.Generate(),
		Title:       input.Title,
		Description: input.Description,
	}
	if err := s.repo.Insert(ctx, course); err != nil {
		return nil, err
	}
	return course, nil
}
