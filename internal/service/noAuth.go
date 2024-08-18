package service

import (
	"avito2024/internal/domain"
	"avito2024/pkg"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) DummyLoginService(status string) (string, error) {
	if status != "client" && status != "moderator" {
		return "", domain.NewCustomError(domain.InvalidInputError())
	}
	token, err := s.DummyLoginRepo(status)
	if err != nil {
		return "", err
	}
	return token, err
}

func (s Service) LoginService(user domain.User) (string, error) {
	if user.ID == "" || user.Password == "" {
		return "", domain.NewCustomError(domain.InvalidInputError())
	}
	token, err := s.LoginRepo(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s Service) RegisterService(user domain.User) (string, error) {
	if !(pkg.IsValidEmail(user.Email)) || user.Password == "" ||
		user.UserType != "client" && user.UserType != "moderator" {
		return "", domain.NewCustomError(domain.InvalidInputError())
	}

	bs, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return "", domain.NewCustomError(domain.InternalError(err))
	}

	user.Password = string(bs)

	token, err := s.RegisterRepo(user)
	if err != nil {
		return "", err
	}

	return token, nil
}
