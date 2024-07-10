package service

import "GinChat/repository"

type ChatService interface {
	Chat()
}

type chatService struct {
	chatRepository repository.ChatRepository
}

func NewChatService(repository repository.ChatRepository) ChatService {
	return &chatService{
		chatRepository: repository,
	}
}

func (c chatService) Chat() {

}
