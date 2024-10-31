package service

import (
	"GinChat/entity"
	"GinChat/repository"
	"GinChat/serializer"
	"errors"
	_ "math/rand/v2"
)

type AuthService interface {
	Register(serializer.RegisterRequest) (bool, error)
	Login(serializer.LoginRequest) (entity.User, error)
}

type authService struct {
	authRepository repository.AuthRepository
}

func NewAuthService(repo repository.AuthRepository) AuthService {
	return &authService{
		authRepository: repo,
	}
}

func (a authService) Register(registerRequest serializer.RegisterRequest) (bool, error) {
	user, err := a.authRepository.FindByPhone(registerRequest.PhoneNo)
	var isSignup = true
	if err == nil {
		if user.Name != nil {
			isSignup = false
		}
		if err = a.authRepository.CheckAndMakeOTP(user.PhoneNo); err != nil {
			return isSignup, err
		}
		return isSignup, nil
	}
	var newUser = entity.User{
		Name:     nil,
		Username: nil,
		PhoneNo:  registerRequest.PhoneNo,
	}
	if err = a.authRepository.NewUserAndMakeOTP(newUser); err != nil {
		return isSignup, err
	}
	return isSignup, nil
}

func (a authService) Login(loginRequest serializer.LoginRequest) (entity.User, error) {
	user, err := a.authRepository.FindByPhone(loginRequest.PhoneNo)
	var errorUser entity.User
	if err != nil {
		return errorUser, err
	}
	if err = a.authRepository.CheckOTP(loginRequest.PhoneNo, loginRequest.Token); err != nil {
		return errorUser, err
	}

	if user.Name == nil || *user.Name == "" {
		if loginRequest.Name == "" {
			return errorUser, errors.New("name_field_required")
		}
		user.Name = &loginRequest.Name
	}
	if err = a.authRepository.UserSave(user); err != nil {
		return errorUser, err
	}

	return user, nil
}
