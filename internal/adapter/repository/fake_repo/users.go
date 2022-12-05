package fake_repo

import (
	"context"
	"errors"
	"github.com/zhuravlev-pe/course-watch/internal/core"
	"github.com/zhuravlev-pe/course-watch/internal/core/service"
)

type users struct {
	byIds   map[string]*core.User
	byEmail map[string]*core.User
}

func newUsers() service.UsersRepository {
	return &users{
		byIds:   map[string]*core.User{},
		byEmail: map[string]*core.User{},
	}
}

func (u *users) GetById(ctx context.Context, id string) (*core.User, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	user, ok := u.byIds[id]
	if !ok {
		return nil, core.ErrNotFound
	}
	return user, nil
}

func (u *users) Insert(ctx context.Context, user *core.User) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	_, ok := u.byIds[user.Id]
	if ok {
		return errors.New("user with the specified id already exists")
	}
	u.byIds[user.Id] = user
	u.byEmail[user.Email] = user
	return nil
}

func (u *users) Update(ctx context.Context, id string, input *core.UpdateUserInfoInput) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	user, ok := u.byIds[id]
	if !ok {
		return core.ErrNotFound
	}
	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.DisplayName = input.DisplayName
	return nil
}

func (u *users) GetByEmail(ctx context.Context, email string) (*core.User, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	user, ok := u.byEmail[email]
	if !ok {
		return nil, core.ErrNotFound
	}
	return user, nil
}
