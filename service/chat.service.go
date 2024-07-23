package service

import (
	"GinChat/repository"
	"GinChat/serializer"
	"fmt"
	"strconv"
)

type ChatService interface {
	GetAllUsers(serializer.PaginationRequest, string) (serializer.APIUserPagination, error)
	GetAllRooms(string) ([]serializer.UserInRoom, error)
	MakePvChat(string, uint) error
}

type chatService struct {
	chatRepository repository.ChatRepository
}

func NewChatService(repository repository.ChatRepository) ChatService {
	return &chatService{
		chatRepository: repository,
	}
}

func (c chatService) GetAllUsers(paginationRequest serializer.PaginationRequest, phoneNo string) (serializer.APIUserPagination, error) {
	var apiUserPagination serializer.APIUserPagination
	userId, err := c.chatRepository.FindByPhone(phoneNo)
	if err != nil {
		return apiUserPagination, err
	}
	users, userCount, err := c.chatRepository.GetAllUsers(paginationRequest, userId)
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
func (c chatService) GetAllRooms(phoneNo string) ([]serializer.UserInRoom, error) {
	userId, _ := c.chatRepository.FindByPhone(phoneNo)
	usersInSameRoom, err := c.chatRepository.GetAllRooms(userId)
	if err != nil {
		return []serializer.UserInRoom{}, err
	}
	return usersInSameRoom, nil
}
func (c chatService) MakePvChat(phoneNo string, recipientId uint) error {
	userId, err := c.chatRepository.FindByPhone(phoneNo)
	if err != nil {
		return err
	}
	roomId, err := c.chatRepository.MakePvChat(userId, recipientId)
	if err != nil {

	}
	fmt.Println(roomId)
	return nil
}
