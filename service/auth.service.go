package service

import (
	"GinChat/repository"
	"errors"
)

type AuthService interface {
	Register() (string, error)
}

type authService struct {
	authRepository repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{
		authRepository: repo,
	}
}

func (a authService) Register() (string, error) {
	return "", errors.New("")
}
