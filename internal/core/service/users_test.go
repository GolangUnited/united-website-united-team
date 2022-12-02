package service

import (
	"context"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zhuravlev-pe/course-watch/internal/core/domain"
	"github.com/zhuravlev-pe/course-watch/internal/core/dto"
	repoMocks "github.com/zhuravlev-pe/course-watch/internal/core/service/mocks"
	"github.com/zhuravlev-pe/course-watch/pkg/idgen"
	"github.com/zhuravlev-pe/course-watch/pkg/security"
	"testing"
	"time"
)

var someDatabaseError = errors.New("some database error")
var someDate = time.Now()

func noError(t *testing.T, err error) {
	assert.NoError(t, err)
}

func TestUsersService_GetUserInfo(t *testing.T) {
	cases := map[string]struct {
		id         string
		setupMocks func(context.Context, *repoMocks.MockUsersRepository)
		output     *dto.GetUserInfoOutput
		checkError func(*testing.T, error)
	}{
		"success": {
			id: "1111111",
			setupMocks: func(ctx context.Context, mockUsersRepo *repoMocks.MockUsersRepository) {
				u := &domain.User{
					Id:               "1111111",
					Email:            "doe.h@example.com",
					FirstName:        "John",
					LastName:         "Doe",
					DisplayName:      "JohnnyD",
					RegistrationDate: someDate,
					Roles:            []security.Role{security.Student},
				}
				mockUsersRepo.EXPECT().GetById(ctx, "1111111").Return(u, nil).Times(1)
			},
			output: &dto.GetUserInfoOutput{
				Id:               "1111111",
				Email:            "doe.h@example.com",
				FirstName:        "John",
				LastName:         "Doe",
				DisplayName:      "JohnnyD",
				RegistrationDate: someDate,
				Roles:            []security.Role{security.Student},
			},
			checkError: noError,
		},
		"not_found": {
			id: "1111111",
			setupMocks: func(ctx context.Context, mockUsersRepo *repoMocks.MockUsersRepository) {
				mockUsersRepo.EXPECT().GetById(ctx, "1111111").Return(nil, domain.ErrNotFound).Times(1)
			},
			output: nil,
			checkError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, domain.ErrNotFound)
			},
		},
		"db_error": {
			id: "1111111",
			setupMocks: func(ctx context.Context, mockUsersRepo *repoMocks.MockUsersRepository) {
				mockUsersRepo.EXPECT().GetById(ctx, "1111111").Return(nil, someDatabaseError).Times(1)
			},
			output: nil,
			checkError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, someDatabaseError)
			},
		},
	}
	
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockUsersRepo := repoMocks.NewMockUsersRepository(mockCtrl)
			gen, err := idgen.New(1)
			assert.NoError(t, err)
			s := NewUsersService(mockUsersRepo, gen)
			ctx := context.Background()
			tc.setupMocks(ctx, mockUsersRepo)
			
			out, err := s.GetUserInfo(ctx, tc.id)
			
			assert.Equal(t, tc.output, out)
			tc.checkError(t, err)
		})
	}
}

func TestUsersService_UpdateUserInfo(t *testing.T) {
	cases := map[string]struct {
		id         string
		input      *dto.UpdateUserInfoInput
		setupMocks func(context.Context, *repoMocks.MockUsersRepository)
		checkError func(*testing.T, error)
	}{
		"success": {
			id: "1111111",
			input: &dto.UpdateUserInfoInput{
				FirstName:   "John",
				LastName:    "Doe",
				DisplayName: "JohnnyD",
			},
			setupMocks: func(ctx context.Context, mockUsersRepo *repoMocks.MockUsersRepository) {
				mockUsersRepo.EXPECT().GetById(ctx, "1111111").Return(nil, nil).Times(1)
				var upd dto.UpdateUserInfoInput
				upd.FirstName = "John"
				upd.LastName = "Doe"
				upd.DisplayName = "JohnnyD"
				mockUsersRepo.EXPECT().Update(ctx, "1111111", &upd).Return(nil).Times(1)
			},
			checkError: noError,
		},
		"not_found": {
			id: "1111111",
			input: &dto.UpdateUserInfoInput{
				FirstName:   "John",
				LastName:    "Doe",
				DisplayName: "JohnnyD",
			},
			setupMocks: func(ctx context.Context, mockUsersRepo *repoMocks.MockUsersRepository) {
				mockUsersRepo.EXPECT().GetById(ctx, "1111111").Return(nil, domain.ErrNotFound).Times(1)
			},
			checkError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, domain.ErrNotFound)
			},
		},
		"db_error_on_get": {
			id: "1111111",
			input: &dto.UpdateUserInfoInput{
				FirstName:   "John",
				LastName:    "Doe",
				DisplayName: "JohnnyD",
			},
			setupMocks: func(ctx context.Context, mockUsersRepo *repoMocks.MockUsersRepository) {
				mockUsersRepo.EXPECT().GetById(ctx, "1111111").Return(nil, someDatabaseError).Times(1)
			},
			checkError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, someDatabaseError)
			},
		},
		"db_error_on_update": {
			id: "1111111",
			input: &dto.UpdateUserInfoInput{
				FirstName:   "John",
				LastName:    "Doe",
				DisplayName: "JohnnyD",
			},
			setupMocks: func(ctx context.Context, mockUsersRepo *repoMocks.MockUsersRepository) {
				mockUsersRepo.EXPECT().GetById(ctx, "1111111").Return(nil, nil).Times(1)
				var upd dto.UpdateUserInfoInput
				upd.FirstName = "John"
				upd.LastName = "Doe"
				upd.DisplayName = "JohnnyD"
				mockUsersRepo.EXPECT().Update(ctx, "1111111", &upd).Return(someDatabaseError).Times(1)
			},
			checkError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, someDatabaseError)
			},
		},
		"validation_firstName_required": {
			id: "1111111",
			input: &dto.UpdateUserInfoInput{
				LastName: "Doe",
			},
			setupMocks: func(ctx context.Context, mockUsersRepo *repoMocks.MockUsersRepository) {},
			checkError: func(t *testing.T, err error) {
				var errs validation.Errors
				ok := errors.As(err, &errs)
				require.True(t, ok)
				assert.Equal(t, 1, len(errs))
				_, ok = errs["first_name"]
				assert.True(t, ok)
			},
		},
		"validation_lastName_required": {
			id: "1111111",
			input: &dto.UpdateUserInfoInput{
				FirstName: "John",
			},
			setupMocks: func(ctx context.Context, mockUsersRepo *repoMocks.MockUsersRepository) {},
			checkError: func(t *testing.T, err error) {
				var errs validation.Errors
				ok := errors.As(err, &errs)
				require.True(t, ok)
				assert.Equal(t, 1, len(errs))
				_, ok = errs["last_name"]
				assert.True(t, ok)
			},
		},
	}
	
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			mockUsersRepo := repoMocks.NewMockUsersRepository(mockCtrl)
			gen, err := idgen.New(1)
			assert.NoError(t, err)
			s := NewUsersService(mockUsersRepo, gen)
			ctx := context.Background()
			tc.setupMocks(ctx, mockUsersRepo)
			
			err = s.UpdateUserInfo(ctx, tc.id, tc.input)
			
			tc.checkError(t, err)
		})
	}
}
