package service

import (
	"context"
	"github.com/zhuravlev-pe/course-watch/internal/core/domain"
	"github.com/zhuravlev-pe/course-watch/internal/core/dto"
	"time"
	
	"github.com/zhuravlev-pe/course-watch/pkg/security"
	"golang.org/x/crypto/bcrypt"
)

type usersService struct {
	repo  UsersRepository
	idGen IDGenerator
}

func NewUsersService(repo UsersRepository, idGen IDGenerator) Users {
	return &usersService{
		repo:  repo,
		idGen: idGen,
	}
}

func (u *usersService) GetUserInfo(ctx context.Context, id string) (*dto.GetUserInfoOutput, error) {
	
	user, err := u.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	
	result := dto.GetUserInfoOutput{
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

func (u *usersService) UpdateUserInfo(ctx context.Context, id string, input *dto.UpdateUserInfoInput) error {
	
	if err := input.Validate(); err != nil {
		return err
	}
	
	_, err := u.repo.GetById(ctx, id)
	if err != nil {
		return err
	}
	
	upd := dto.UpdateUserInfoInput{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		DisplayName: input.DisplayName,
	}
	return u.repo.Update(ctx, id, &upd)
}

func (u *usersService) Signup(ctx context.Context, input *dto.SignupUserInput) error {
	
	user, err := u.repo.GetByEmail(ctx, input.Email)
	if err != nil && err != domain.ErrNotFound {
		//Any error other than ErrorNotFound should stop the Signup flow as ErrorNotFound is valid for the user Signup
		return err
	}
	
	if user != nil && user.Email == input.Email {
		return domain.ErrUserAlreadyExist
	}
	
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	
	user = &domain.User{
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

func (u *usersService) Login(ctx context.Context, input *dto.LoginInput) (*domain.User, error) {
	
	user, err := u.repo.GetByEmail(ctx, input.Email)
	if err != nil && err != domain.ErrNotFound {
		return nil, err
	}
	
	if err == domain.ErrNotFound {
		return nil, domain.ErrInvalidCredentials
	}
	
	if err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(input.Password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}
	return user, nil
}
