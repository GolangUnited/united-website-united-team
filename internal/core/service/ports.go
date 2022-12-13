package service

//go:generate mockgen -source=$GOFILE -package=mock_adapter -destination=mocks/adapters.go

import (
	"context"
	"github.com/zhuravlev-pe/course-watch/internal/core"
)

type UsersRepository interface {
	GetById(ctx context.Context, id string) (*core.User, error)
	Insert(ctx context.Context, user *core.User) error
	Update(ctx context.Context, id string, input *core.UpdateUserInfoInput) error
	GetByEmail(ctx context.Context, email string) (*core.User, error)
}

type CoursesRepository interface {
	GetById(ctx context.Context, id string) (*core.Course, error)
	Insert(ctx context.Context, course *core.Course) error
}

type IdGenerator interface {
	Generate() string
}
