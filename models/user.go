package models

import (
	"fmt"

	"github.com/PAY-HERO-CONSULTING/gh-tools/custom_types"
	"github.com/PAY-HERO-CONSULTING/gh-tools/null"
)

type User struct {
	custom_types.UUIDIdentifier
	FirstName       *string                 `json:"first_name"`
	LastName        *string                 `json:"last_name"`
	DialCode        string                  `json:"dial_code"`
	Number          string                  `json:"number"`
	UserName        string                  `json:"username"`
	Email           string                  `json:"email"`
	Password        string                  `json:"_"`
	IsAdmin         bool                    `json:"is_admin"`
	ProfileFileName *string                 `json:"profile_filename"`
	IsPinSet        bool                    `json:"is_pin_set"`
	PhoneVerified   bool                    `json:"phone_verified"`
	EmailVefifed    bool                    `json:"email_verified"`
	Status          custom_types.UserStatus `json:"status"`
	custom_types.Timestamps
}

func NewUser(username string) *User {
	return &User{
		FirstName: null.NullValue(custom_types.GenerateRandomString(70)),
		LastName:  null.NullValue(custom_types.GenerateRandomString(70)),
		UserName:  custom_types.GenerateRandomString(100),
		Email:     fmt.Sprintf("%s@mail.com", username),
		Password:  username,
		IsAdmin:   false,
		Status:    custom_types.UserStatusVerified,
	}
}

func NewAdmin(username string) *User {
	return &User{
		FirstName: null.NullValue(custom_types.GenerateRandomString(70)),
		LastName:  null.NullValue(custom_types.GenerateRandomString(70)),
		UserName:  custom_types.GenerateRandomString(100),
		Email:     fmt.Sprintf("%s@mail.com", username),
		Password:  username,
		IsAdmin:   true,
		Status:    custom_types.UserStatusVerified,
	}
}
