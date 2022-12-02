package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/zhuravlev-pe/course-watch/pkg/security"
	"time"
)

type SignupUserInput struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DisplayName string `json:"display_name"`
}

type LoginInput struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	Persistent bool   `json:"persistent"`
}

type UpdateUserInfoInput struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DisplayName string `json:"display_name"`
}

func (i *UpdateUserInfoInput) Validate() error {
	return validation.ValidateStruct(i,
		validation.Field(&i.FirstName, validation.Required),
		validation.Field(&i.LastName, validation.Required),
	)
}

type GetUserInfoOutput struct {
	Id               string          `json:"id"`
	Email            string          `json:"email"`
	FirstName        string          `json:"first_name"`
	LastName         string          `json:"last_name"`
	DisplayName      string          `json:"display_name"`
	RegistrationDate time.Time       `json:"registration_date"`
	Roles            []security.Role `json:"roles"`
}
