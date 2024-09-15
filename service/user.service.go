package service

import (
	"GinChat/entity"
	"GinChat/repository"
	"GinChat/serializer"
	"GinChat/utils"
	"errors"
	"fmt"
	"strconv"
)

type UserService interface {
	GetAllUsers(serializer.PaginationRequest, string) (serializer.APIUserPagination, error)
	ProfileUpdate(serializer.ProfileUpdateRequest) (serializer.UpdatedProfile, error)
	GetUserProfile(uint) (serializer.ProfileAPI, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		userRepository: repo,
	}
}
func (u userService) GetAllUsers(paginationRequest serializer.PaginationRequest, phoneNo string) (serializer.APIUserPagination, error) {
	var apiUserPagination serializer.APIUserPagination
	userId, err := u.userRepository.FindByPhone(phoneNo)
	if err != nil {
		return apiUserPagination, err
	}
	users, userCount, err := u.userRepository.GetAllUsers(paginationRequest, userId)
	if err != nil {
		return apiUserPagination, err
	}

	var next string
	if int64(paginationRequest.Limit+paginationRequest.Offset) > userCount {
		next = "no next"
	} else {
		next = fmt.Sprintf(
			"http://localhost:8080/api/v1/chat/get-users?limit=%s&offset=%s",
			strconv.Itoa(paginationRequest.Limit),
			strconv.Itoa(paginationRequest.Limit+paginationRequest.Offset),
		)
	}
	var previous string
	if paginationRequest.Offset >= paginationRequest.Limit {
		previous = fmt.Sprintf(
			"http://localhost:8080/api/v1/chat/get-users?limit=%s&offset=%s",
			strconv.Itoa(paginationRequest.Limit),
			strconv.Itoa(paginationRequest.Offset-paginationRequest.Limit),
		)
	} else {
		previous = "no previous"
	}
	apiUserPagination.Count = userCount
	apiUserPagination.Results = users
	apiUserPagination.Next = next
	apiUserPagination.Previous = previous

	return apiUserPagination, nil
}
func (u userService) ProfileUpdate(profile serializer.ProfileUpdateRequest) (serializer.UpdatedProfile, error) {
	// unique image name check
	var imagePath string
	if profile.Avatar != nil {
		if ok := utils.ImageValidate(profile.Avatar); !ok {
			return serializer.UpdatedProfile{}, errors.New("bad_format")
		}
		imagePath = "assets/uploads/userAvatar/"
		imagePath = utils.ImageController(imagePath, profile.Avatar.Filename)
	}

	user := entity.User{
		ID:       profile.ID,
		Name:     &profile.Name,
		Username: &profile.Username,
		Avatar:   &imagePath,
	}

	updatedProfile, err := u.userRepository.ProfileUpdate(user)
	if err != nil {
		return serializer.UpdatedProfile{}, err
	}
	return updatedProfile, nil

}
func (u userService) GetUserProfile(id uint) (serializer.ProfileAPI, error) {
	var user entity.User
	user.ID = id
	userProfile, err := u.userRepository.GetUserProfile(user)
	if err != nil {
		return serializer.ProfileAPI{}, err
	}
	return userProfile, nil
}
