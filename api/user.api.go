package api

import (
	"GinChat/serializer"
	"GinChat/service"
	"GinChat/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
)

type UserAPI interface {
	GetAllUsers(request *gin.Context)
	ProfileUpdate(*gin.Context)
	GetUserProfile(*gin.Context)
}
type userAPI struct {
	service service.UserService
}

func NewUserAPI(service service.UserService) UserAPI {
	return &userAPI{
		service: service,
	}
}

// @Summary get all users
// @Description get all users
// @Description if their have a room in pv chat it will come with it
// @Description this url need get-users?offset=0&limit=0
// @Tags user
// @Accept  json
// @Produce  json
// @Success 200 {object}   serializer.APIUserPagination
// @Failure 400 {object}   utils.ErrorResponse "Token_Expired_Or_Invalid(2) | Bad_Format(5)"
// @Failure 500 {object}   utils.ErrorResponse "We_Don't_Know_What_Happened(8)"
// @Router /chat/get-users/ [get]
func (u userAPI) GetAllUsers(request *gin.Context) {
	var paginationRequest serializer.PaginationRequest
	if err := request.ShouldBindQuery(&paginationRequest); err != nil {
		request.JSON(http.StatusBadRequest, utils.BadFormat)
		return
	}
	userPhoneNo, ok := request.Get("phoneNo")
	if !ok {
		request.JSON(http.StatusBadRequest, utils.TokenIsExpiredOrInvalid)
		return
	}
	apiUserPagination, err := u.service.GetAllUsers(paginationRequest, userPhoneNo.(string))
	if err != nil {
		request.JSON(http.StatusInternalServerError, utils.SomethingWentWrong)
		return
	}

	request.JSON(http.StatusOK, apiUserPagination)
	return
}

// @Summary update user profile
// @Description send this in form-data
// @Description it has
// @Tags user
// @Accept  json
// @Produce  json
// @Param   message  body  utils.DummyProfileUpdate  true  "Message body"
// @Success 201 {object}   serializer.UpdatedProfile
// @Failure 400 {object}   utils.ErrorResponse "Token_Expired_Or_Invalid(2) | Object_Not_Found(6) | Bad_Format(5) | User_Name_Is_Taken(11)"
// @Failure 500 {object}   utils.ErrorResponse "We_Don't_Know_What_Happened(8)"
// @Router /user/profile-update/:id/ [get]
func (u userAPI) ProfileUpdate(request *gin.Context) {
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

	updatedProfile, err := u.service.ProfileUpdate(profileUpdateRequest)
	if err != nil {
		if err.Error() == "username_taken" {
			request.JSON(http.StatusBadRequest, utils.UserNameIsTaken)
			return
		}
		if err.Error() == "bad_format" {
			request.JSON(http.StatusBadRequest, utils.BadFormat)
			return
		}
		request.JSON(http.StatusBadRequest, utils.SomethingWentWrong)
		return
	}
	if profileUpdateRequest.Avatar != nil {
		profileUpdateRequest.Avatar.Filename = updatedProfile.Avatar[26:]
		if err = request.SaveUploadedFile(profileUpdateRequest.Avatar, "assets/uploads/userAvatar/"+profileUpdateRequest.Avatar.Filename); err != nil {
			request.JSON(http.StatusBadRequest, utils.SomethingWentWrong)
			return
		}
	}

	request.JSON(http.StatusOK, updatedProfile)
	return
}

// todo:need docs!
func (u userAPI) GetUserProfile(request *gin.Context) {
	id, err := strconv.ParseInt(request.Param("id"), 10, 32)
	if err != nil {
		request.JSON(http.StatusBadRequest, utils.BadFormat)
		return
	}
	userProfile, err := u.service.GetUserProfile(uint(id))
	if err != nil {
		request.JSON(http.StatusNotFound, utils.ObjectNotFound)
		return
	}
	request.JSON(http.StatusOK, userProfile)
}
