package service

import (
	"context"
	"github.com/zhuravlev-pe/course-watch/internal/core"
	"time"
	
	"github.com/zhuravlev-pe/course-watch/pkg/security"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo  UsersRepository
	idGen IdGenerator
}

func NewUserService(repo UsersRepository, idGen IdGenerator) *UserService {
	return &UserService{
		repo:  repo,
		idGen: idGen,
	}
}

func (u *UserService) GetUserInfo(ctx context.Context, id string) (*core.GetUserInfoOutput, error) {
	
	user, err := u.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	
	result := core.GetUserInfoOutput{
		Id:               user.Id,
		Email:            user.Email,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		DisplayName:      user.DisplayName,
		RegistrationDate: user.RegistrationDate,
		Roles:            user.Roles,
	}
	return &result, nil
}

func (u *UserService) UpdateUserInfo(ctx context.Context, id string, input *core.UpdateUserInfoInput) error {
	
	if err := input.Validate(); err != nil {
		return err
	}
	
	_, err := u.repo.GetById(ctx, id)
	if err != nil {
		return err
	}
	
	upd := core.UpdateUserInfoInput{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		DisplayName: input.DisplayName,
	}
	return u.repo.Update(ctx, id, &upd)
}

func (u *UserService) Signup(ctx context.Context, input *core.SignupUserInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	user, err := u.repo.GetByEmail(ctx, input.Email)
	if err != nil && err != core.ErrNotFound {
		//Any error other than ErrorNotFound should stop the Signup flow as ErrorNotFound is valid for the user Signup
		return err
	}
	
	if user != nil && user.Email == input.Email {
		return core.ErrUserAlreadyExist
	}
	
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user = &core.User{
		Id:               u.idGen.Generate(),
		Email:            input.Email,
		FirstName:        input.FirstName,
		LastName:         input.LastName,
		DisplayName:      input.DisplayName,
		RegistrationDate: time.Now(),
		HashedPassword:   hashPassword,
		Roles:            []security.Role{security.Student},
	}
	
	if err = u.repo.Insert(ctx, user); err != nil {
		return err
	}
	return nil
}

func (u *UserService) Login(ctx context.Context, input *core.LoginInput) (*core.User, error) {
	
	user, err := u.repo.GetByEmail(ctx, input.Email)
	if err != nil && err != core.ErrNotFound {
		return nil, err
	}
	
	if err == core.ErrNotFound {
		return nil, core.ErrInvalidCredentials
	}
	
	if err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(input.Password)); err != nil {
		return nil, core.ErrInvalidCredentials
	}
	return user, nil
}
