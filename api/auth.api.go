package api

import (
	"GinChat/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthAPI interface {
	Register(*gin.Context)
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
	request.JSON(http.StatusOK, gin.H{})
	return
}

func (a authAPI) Login(request *gin.Context) {
	request.JSON(200, gin.H{})
	return
}
