package service

import (
	"GinChat/entity"
	"GinChat/repository"
	"GinChat/serializer"
	"GinChat/utils"
	"errors"
	"fmt"
	_ "math/rand/v2"
	"time"
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

// todo: we need celery for this
func (a authService) Register(registerRequest serializer.RegisterRequest) (bool, error) {
	user, err := a.authRepository.FindByPhone(registerRequest.PhoneNo)
	var isSignup = true
	if err == nil {
		if user.Name != nil {
			isSignup = false
		}

		if user.Phone.ExpTime != nil && user.Phone.ExpTime.After(time.Now()) {
			return isSignup, errors.New("too_many_request")
		}
		var expTime = utils.GetExpiryTime()
		var token = utils.SmsTokenGenerate()
		user.Phone.Token = &token
		user.Phone.ExpTime = &expTime
		if err = a.authRepository.PhoneSave(user.Phone); err != nil {
			return isSignup, err
		}
		fmt.Println(*user.Phone.Token)
		return isSignup, nil
	}

	var expTime = utils.GetExpiryTime()
	var token = utils.SmsTokenGenerate()
	var newUser = entity.User{
		Name:     nil,
		Username: nil,
		Phone: entity.Phone{
			PhoneNo: registerRequest.PhoneNo,
			Token:   &token,
			ExpTime: &expTime,
		},
	}
	fmt.Println(*newUser.Phone.Token)
	if err = a.authRepository.NewUserSave(newUser); err != nil {
		return isSignup, err
	}
	return isSignup, nil
}
func (a authService) Login(loginRequest serializer.LoginRequest) (entity.User, error) {
	user, err := a.authRepository.FindByPhone(loginRequest.PhoneNo)
	if err != nil {
		return entity.User{}, err
	}

	if user.Phone.ExpTime == nil || user.Phone.ExpTime.Before(time.Now()) {
		return entity.User{}, errors.New("expired_time")
	}
	if loginRequest.Token != *user.Phone.Token {
		return entity.User{}, errors.New("invalid_token")
	}
	if user.Name == nil || *user.Name == "" {
		if loginRequest.Name == "" {
			return entity.User{}, errors.New("name_field_required")
		}
		user.Name = &loginRequest.Name
	}
	user.Phone.ExpTime = nil
	user.Phone.Token = nil

	if err = a.authRepository.UserSave(user); err != nil {
		return entity.User{}, err
	}
	if err = a.authRepository.PhoneSave(user.Phone); err != nil {
		return entity.User{}, err
	}

	return user, nil
}
