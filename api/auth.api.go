package api

import (
	"GinChat/pkg/JWT"
	"GinChat/serializer"
	"GinChat/service"
	"GinChat/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
)

type AuthAPI interface {
	Register(*gin.Context)
	Login(*gin.Context)
	ProfileUpdate(request *gin.Context)
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
		request.JSON(http.StatusCreated, gin.H{"data": "sent", "is_signup": isSignup})
		return
	}
	request.JSON(http.StatusOK, gin.H{"data": "sent", "is_signup": isSignup})
	return
}
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
		case "expired_time":
			request.JSON(http.StatusBadRequest, utils.TokenIsExpiredOrInvalid)
			return
		case "invalid_token":
			request.JSON(http.StatusBadRequest, utils.TokenIsExpiredOrInvalid)
			return
		case "name_field_required":
			request.JSON(http.StatusBadRequest, gin.H{"error": "name field for the first login is required", "status": false})
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
func (a authAPI) ProfileUpdate(request *gin.Context) {
	var profileUpdateRequest serializer.ProfileUpdateRequest
	if err := request.ShouldBindWith(&profileUpdateRequest, binding.FormMultipart); err != nil {
		request.JSON(http.StatusBadRequest, utils.BadFormat)
		return
	}
	id, err := strconv.ParseInt(request.Param("id"), 10, 32)
	if err != nil {
		request.JSON(http.StatusBadRequest, utils.BadFormat)
		return
	}
	profileUpdateRequest.ID = uint(id)

	updatedProfile, err := a.service.ProfileUpdate(profileUpdateRequest)
	if err != nil {
		if err.Error() == "username_taken" {
			request.JSON(http.StatusBadRequest, utils.UserNameIsTaken)
			return
		}
		request.JSON(http.StatusBadRequest, utils.SomethingWentWrong)
	}
	profileUpdateRequest.Avatar.Filename = updatedProfile.Avatar[26:]

	if err = request.SaveUploadedFile(profileUpdateRequest.Avatar, "assets/uploads/"+profileUpdateRequest.Avatar.Filename); err != nil {
		request.JSON(http.StatusBadRequest, utils.SomethingWentWrong)
		return
	}
	request.JSON(http.StatusOK, updatedProfile)
	return
}
