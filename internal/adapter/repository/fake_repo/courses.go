package fake_repo

import (
	"context"
	"github.com/zhuravlev-pe/course-watch/internal/core/domain"
	"github.com/zhuravlev-pe/course-watch/internal/core/service"
)

type courses struct {
	data map[string]*domain.Course
}

func NewCourses() service.CoursesRepository {
	return &courses{
		data: map[string]*domain.Course{},
	}
}

func (c *courses) GetById(_ context.Context, id string) (*domain.Course, error) {
	course, ok := c.data[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return course, nil
}

func (c *courses) Insert(_ context.Context, course *domain.Course) error {
	c.data[course.Id] = course
	return nil
}
