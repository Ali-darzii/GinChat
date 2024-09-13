package service

import (
	"GinChat/entity"
	"GinChat/repository"
	"GinChat/serializer"
	"GinChat/utils"
	"errors"
)

type ChatService interface {
	GetAllRooms(string) ([]serializer.Room, error)
	MakePvChat(serializer.MakeNewChatRequest, string) (serializer.Message, error)
	MakeGroupChat(serializer.MakeGroupChatRequest, string) (serializer.Message, error)
}

type chatService struct {
	chatRepository repository.ChatRepository
}

func NewChatService(repository repository.ChatRepository) ChatService {
	return &chatService{
		chatRepository: repository,
	}
}

func (c chatService) GetAllRooms(phoneNo string) ([]serializer.Room, error) {
	userId, _ := c.chatRepository.FindByPhone(phoneNo)
	usersInSameRoom, err := c.chatRepository.GetAllRooms(userId)
	if err != nil {
		return []serializer.Room{}, err
	}
	return usersInSameRoom, nil
}
func (c chatService) MakePvChat(makeNewChatRequest serializer.MakeNewChatRequest, phoneNo string) (serializer.Message, error) {
	userId, err := c.chatRepository.FindByPhone(phoneNo)
	if err != nil {
		return serializer.Message{}, err
	}
	message, err := c.chatRepository.MakePvChat(makeNewChatRequest, userId)
	if err != nil {
		return serializer.Message{}, err
	}
	return message, nil
}
func (c chatService) MakeGroupChat(makeGroupChatRequest serializer.MakeGroupChatRequest, phoneNo string) (serializer.Message, error) {
	userId, err := c.chatRepository.FindByPhone(phoneNo)

	if err != nil {
		return serializer.Message{}, err
	}
	var imagePath string
	if makeGroupChatRequest.Avatar != nil {
		if ok := utils.ImageValidate(makeGroupChatRequest.Avatar); !ok {
			return serializer.Message{}, errors.New("bad_format")
		}
		imagePath = "assets/uploads/groupAvatar/"
		imagePath = utils.ImageController(imagePath, makeGroupChatRequest.Avatar.Filename)
	}
	groupRoom := entity.GroupRoom{
		Avatar: &imagePath,
		Name:   makeGroupChatRequest.Name,
		Users:  []entity.User{{ID: userId}},
		Admins: []entity.User{{ID: userId}},
	}

	for _, id := range makeGroupChatRequest.Recipients {
		groupRoom.Users = append(groupRoom.Users, entity.User{ID: id})
	}

	groupRoom, err = c.chatRepository.MakeGroupChat(groupRoom)
	if err != nil {
		return serializer.Message{}, err
	}

	message := serializer.Message{
		Type:       "new_group_message",
		Avatar:     imagePath,
		RoomID:     groupRoom.ID,
		Sender:     userId,
		Recipients: append(makeGroupChatRequest.Recipients, userId),
	}

	//websocketHandler.Manager.Broadcast <- message, moved to api

	return message, nil

}
