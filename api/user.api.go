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
}
type userAPI struct {
	service service.UserService
}

func NewUserAPI(service service.UserService) UserAPI {
	return &userAPI{
		service: service,
	}
}

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
