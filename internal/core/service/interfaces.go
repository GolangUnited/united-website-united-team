package service

//go:generate mockgen -source=$GOFILE -destination=mocks/services.go

import (
	"context"
	"github.com/zhuravlev-pe/course-watch/internal/core/domain"
	"github.com/zhuravlev-pe/course-watch/internal/core/dto"
)

type UsersRepository interface {
	GetById(ctx context.Context, id string) (*domain.User, error)
	Insert(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, id string, input *dto.UpdateUserInfoInput) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type CoursesRepository interface {
	GetById(ctx context.Context, id string) (*domain.Course, error)
	Insert(ctx context.Context, course *domain.Course) error
}

type Users interface {
	GetUserInfo(ctx context.Context, id string) (*dto.GetUserInfoOutput, error)
	UpdateUserInfo(ctx context.Context, id string, input *dto.UpdateUserInfoInput) error
	Login(ctx context.Context, input *dto.LoginInput) (*domain.User, error)
	Signup(ctx context.Context, input *dto.SignupUserInput) error
}

type Courses interface {
	GetById(ctx context.Context, id string) (*domain.Course, error)
	Create(ctx context.Context, input dto.CreateCourseInput) (*domain.Course, error)
}

type IDGenerator interface {
	Generate() string
}

type Services struct {
	Courses Courses
	Users   Users
}

func NewServices(deps Deps) *Services {
	return &Services{
		Courses: NewCoursesService(deps.Repos.Courses, deps.IdGen),
		Users:   NewUsersService(deps.Repos.Users, deps.IdGen),
	}
}
