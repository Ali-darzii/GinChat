package service

import (
	"GinChat/repository"
)

type ChatService interface {
	WsHandler(any) (uint, error)
}

type chatService struct {
	chatRepository repository.ChatRepository
}

func NewChatService(repository repository.ChatRepository) ChatService {
	return &chatService{
		chatRepository: repository,
	}
}

func (c chatService) WsHandler(phoneNo any) (uint, error) {
	userId, err := c.chatRepository.WsHandler(phoneNo)
	if err != nil {
		return 0, err
	}
	return userId, nil

}
