package models

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost     = 12
	minUserNameLen = 2
	minPasswordLen = 7
)

type CreateUserParams struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (params CreateUserParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(params.UserName) < minUserNameLen {
		errors["username"] = fmt.Sprintf("username length should be at least %d characters", minUserNameLen)
	}
	if len(params.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("password length should be at least %d characters", minPasswordLen)
	}
	if !IsEmailValid(params.Email) {
		errors["email"] = fmt.Sprintf("email %s is invalid", params.Email)
	}
	return errors
}

func IsEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}

	userID := NewUUID()

	return &User{
		ID:                userID,
		UserName:          params.UserName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}

func IsValidPassword(encpw, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(pw)) == nil
}

type UpdateUserParams struct {
	UserName string `json:"username"`
}

func (p UpdateUserParams) ToFieldsMap() map[string]interface{} {
	fields := map[string]interface{}{}
	if len(p.UserName) > 0 {
		fields["username"] = p.UserName
	}
	return fields
}
