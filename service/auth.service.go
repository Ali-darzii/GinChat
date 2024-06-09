package service

import (
	"GinChat/entity"
	"GinChat/repository"
	"GinChat/serializer"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand/v2"
	_ "math/rand/v2"
	"time"
)

type AuthService interface {
	Register(serializer.RegisterRequest) (gin.H, error)
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

func (a authService) Register(registerRequest serializer.RegisterRequest) (gin.H, error) {
	user, err := a.authRepository.FindByPhone(registerRequest.PhoneNo)

	if err == nil {
		timeNow := time.Now()
		var expTime time.Time = timeNow.Add(time.Second * 60)
		var token int = rand.IntN(8999) + 1000
		user.Phone.Token = &token
		user.Phone.ExpTime = &expTime
		fmt.Println(*user.Phone.Token)
		return gin.H{"detail": "send", "is_signup": false}, nil
	}

	timeNow := time.Now()
	var expTime time.Time = timeNow.Add(time.Second * 60)
	var token int = rand.IntN(8999) + 1000
	var newUser = entity.User{
		Name:     registerRequest.Name,
		Username: nil,
		Phone: entity.Phone{
			PhoneNo: registerRequest.PhoneNo,
			Token:   &token,
			ExpTime: &expTime,
		},
	}
	fmt.Println(*newUser.Phone.Token)
	if err := a.authRepository.Register(newUser); err != nil {
		return gin.H{}, err
	}
	return gin.H{"detail": "send", "is_signup": true}, nil
}
func (a authService) Login(loginRequest serializer.LoginRequest) (entity.User, error) {
	user, err := a.authRepository.FindByPhone(loginRequest.PhoneNo)
	if err != nil {
		return entity.User{}, err
	}
	if user.Phone.ExpTime.Before(time.Now()) {
		return entity.User{}, errors.New("expired_time")
	}
	if loginRequest.Token != user.Phone.Token {
		return entity.User{}, errors.New("invalid_token")
	}
	user.Phone.Token = nil

	return user, nil
}
