package v1

//go:generate mockgen -source=$GOFILE -package=mock_service -destination=mocks/services.go

import (
	"context"
	"github.com/zhuravlev-pe/course-watch/internal/core"
)

type UserService interface {
	GetUserInfo(ctx context.Context, id string) (*core.GetUserInfoOutput, error)
	UpdateUserInfo(ctx context.Context, id string, input *core.UpdateUserInfoInput) error
	Login(ctx context.Context, input *core.LoginInput) (*core.User, error)
	Signup(ctx context.Context, input *core.SignupUserInput) error
}

type CourseService interface {
	GetById(ctx context.Context, id string) (*core.Course, error)
	Create(ctx context.Context, input core.CreateCourseInput) (*core.Course, error)
}
