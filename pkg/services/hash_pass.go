package services

import (
	"fmt"
	"golang_todo/pkg/types"

	"golang.org/x/crypto/bcrypt"
)

type UserServices struct{}

func NewUserServices() *UserServices {
	return &UserServices{}
}

func (s *UserServices) HashPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing the password: %v", err)
	}
	return string(hashedPass), nil
}

func (s *UserServices) CheckPassword(userPass types.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(userPass.password), []byte(password))
}
