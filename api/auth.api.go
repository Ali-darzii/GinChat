package api

import (
	"GinChat/pkg/JWT"
	"GinChat/serializer"
	"GinChat/service"
	"GinChat/utils"
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

// @Summary send token
// @Description 1 min for every request, not authenticated, and returns a JWT token
// @Tags Authenticate
// @Accept  json
// @Produce  json
// @Param   Register  body  serializer.RegisterRequest  true  "Register details"
// @Success 201 {object} utils.RegisterResponse
// @Success 200 {object} utils.RegisterResponse
// @Failure 400 {object} utils.ErrorResponse "Too_Many_Token_Request(7) | Token_Expired_Or_Invalid(2) | We_Don't_Know_What_Happened(8) | MUST_NOT_AUTHENTICATED(1)"
// @Router /auth/ [post]
func (a authAPI) Register(request *gin.Context) {
	/*  send phone msg with db data creation  */
	var registerRequest serializer.RegisterRequest
	if err := request.ShouldBind(&registerRequest); err != nil {
		request.JSON(http.StatusBadRequest, utils.BadFormat)
		return
	}
	isSignup, err := a.service.Register(registerRequest)
	if err != nil {
		if err.Error() == "too_many_request" {
			request.JSON(http.StatusTooManyRequests, utils.ExpiredTimeBlocked)
			return
		}

		request.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if isSignup {
		request.JSON(http.StatusCreated, utils.RegisterResponse{Detail: "sent", IsSignup: isSignup})
		return
	}
	request.JSON(http.StatusOK, utils.RegisterResponse{Detail: "sent", IsSignup: isSignup})
	return
}

// @Summary check token
// @Description Authenticates a user and returns a JWT token
// @Tags Authenticate
// @Accept  json
// @Produce  json
// @Param   login  body  serializer.LoginRequest  true  "Login details"
// @Success 200 {object} serializer.Token
// @Failure 400 {object} utils.ErrorResponse "Object_Not_Found(6) | Token_Expired_Or_Invalid(2) | Name_Field_Required_For_Register(12) | We_Don't_Know_What_Happened(8) | MUST_NOT_AUTHENTICATED(1)"
// @Router /auth/ [put]
func (a authAPI) Login(request *gin.Context) {
	var loginRequest serializer.LoginRequest

	if err := request.ShouldBind(&loginRequest); err != nil {
		request.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err := a.service.Login(loginRequest)
	if err != nil {
		switch err.Error() {
		case "not_found":
			request.JSON(http.StatusBadRequest, utils.ObjectNotFound)
			return
		case "expired_time", "invalid_token":
			request.JSON(http.StatusBadRequest, utils.TokenIsExpiredOrInvalid)
			return
		case "name_field_required":
			request.JSON(http.StatusBadRequest, utils.NameFieldRequired)
			return
		default:
			request.JSON(http.StatusBadRequest, err.Error())
			return
		}
	}

	// todo: test issue
	request.SetCookie(
		"loginAttempt",
		"", -1,
		"/",
		"localhost",
		false,
		true,
	)
	if err := utils.UserLoggedIn(request, user); err != nil {
		request.JSON(http.StatusBadRequest, err.Error())
		return
	}
	jwt := JWT.Jwt{}
	token, err := jwt.CreateToken(user)
	if err != nil {
		request.JSON(http.StatusBadRequest, utils.SomethingWentWrong)
		return
	}
	request.JSON(http.StatusOK, token)
	return
}
