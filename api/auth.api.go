package api

import (
	"GinChat/serializer"
	"GinChat/service"
	"GinChat/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthAPI interface {
	Register(*gin.Context)
	CreateToken() error
	Login(*gin.Context)
}

type authAPI struct {
	service service.AuthService
}

func NewAuthAPI(service service.AuthService) AuthAPI {
	return &authAPI{
		service: service,
	}
}

func (a authAPI) Register(request *gin.Context) {
	/*  send phone msg with db data creation  */
	var registerRequest serializer.RegisterRequest
	if err := request.ShouldBind(&registerRequest); err != nil {
		request.JSON(http.StatusBadRequest, utils.BadFormat)
		return
	}
	if err := registerRequest.PhoneNoValidate(registerRequest.PhoneNo); err != nil {
		request.JSON(http.StatusBadRequest, utils.BadFormat)
		return
	}
	if registerRequest.UserName != "" {
		if err := registerRequest.UsernameValidate(registerRequest.UserName); err != nil {
			request.JSON(http.StatusBadRequest, utils.BadFormat)
			return
		}
	}
	if err := a.service.Register(registerRequest); err != nil {
		if err.Error() == "unique_field" {
			request.JSON(http.StatusBadRequest, utils.UniqueField)
			return
		}
		request.JSON(http.StatusBadRequest, err.Error())
		return
	}
	request.JSON(http.StatusCreated, gin.H{"data": "Token generated", "status": true})
	return
}

func (a authAPI) Login(request *gin.Context) {
	var loginRequest serializer.LoginRequest
	if err := request.ShouldBind(&loginRequest); err != nil {
		request.JSON(http.StatusBadRequest, utils.BadFormat)
		return
	}
	if err := loginRequest.PhoneNoValidate(loginRequest.PhoneNo); err != nil {
		request.JSON(http.StatusBadRequest, utils.BadFormat)
		return
	}
	user, err := a.service.Login(loginRequest)
	if err != nil {
		if err.Error() == "not_found" {
			request.JSON(http.StatusBadRequest, utils.ObjectNotFound)
			return
		}
		request.JSON(http.StatusBadRequest, err.Error())
		return
	}
	request.JSON(http.StatusOK, user)
	return
}

func (a authAPI) CreateToken() error {
	return nil
}
