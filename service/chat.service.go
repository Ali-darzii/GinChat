package service

import (
	"GinChat/entity"
	"GinChat/repository"
	"GinChat/serializer"
	"GinChat/utils"
	"errors"
	"path/filepath"
)

type ChatService interface {
	GetAllRooms(string) ([]serializer.Room, error)
	MakePvChat(serializer.MakeNewChatRequest, string) (serializer.Message, error)
	MakeGroupChat(serializer.MakeGroupChatRequest, string) (serializer.Message, error)
	SendPvMessage(serializer.MessageRequest, string) (serializer.Message, error)
	SendGpMessage(serializer.MessageRequest, string) (serializer.Message, error)
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
	var message serializer.Message
	if err != nil {
		return message, err
	}
	var imagePath string
	if makeNewChatRequest.File != nil {
		ext := filepath.Ext(makeNewChatRequest.File.Filename)
		if ext != "mp3" {
			if ok := utils.ImageValidate(makeNewChatRequest.File); !ok {
				return message, errors.New("bad_format")
			}
		}
		imagePath = "assets/uploads/pvMessage/"
		imagePath = utils.FilePathController(imagePath, makeNewChatRequest.File.Filename)
	}

	privateRoom := entity.PrivateRoom{
		Users: []entity.User{
			{ID: userId},
			{ID: makeNewChatRequest.RecipientID},
		},
	}
	privateChat := entity.PrivateMessageRoom{
		Sender: userId,
		Body:   &makeNewChatRequest.Content,
		File:   &imagePath,
	}

	privateChat, err = c.chatRepository.MakePvChat(privateRoom, privateChat)
	if err != nil {
		return message, err
	}
	message.PvMessage.Type = "new_pv_message"
	message.PvMessage.File = imagePath
	message.PvMessage.RoomID = privateChat.PrivateID
	message.PvMessage.Sender = userId
	message.PvMessage.Content = *privateChat.Body
	message.Recipients = []uint{userId, makeNewChatRequest.RecipientID}

	return message, nil
}
func (c chatService) MakeGroupChat(makeGroupChatRequest serializer.MakeGroupChatRequest, phoneNo string) (serializer.Message, error) {
	userId, err := c.chatRepository.FindByPhone(phoneNo)
	var message serializer.Message
	if err != nil {
		return message, err
	}
	var imagePath string
	if makeGroupChatRequest.Avatar != nil {
		if ok := utils.ImageValidate(makeGroupChatRequest.Avatar); !ok {
			return message, errors.New("bad_format")
		}
		imagePath = "assets/uploads/groupAvatar/"
		imagePath = utils.FilePathController(imagePath, makeGroupChatRequest.Avatar.Filename)
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
		return message, err
	}

	message.Avatar = imagePath
	message.Recipients = append(makeGroupChatRequest.Recipients, userId)
	message.PvMessage.Type = "new_gp_message"
	message.PvMessage.RoomID = groupRoom.ID
	message.PvMessage.Sender = userId

	return message, nil

}
func (c chatService) SendPvMessage(pvMessage serializer.MessageRequest, phoneNo string) (serializer.Message, error) {
	userId, err := c.chatRepository.FindByPhone(phoneNo)
	var message serializer.Message
	if err != nil {
		return message, err
	}

	var imagePath string
	if pvMessage.File != nil {
		ext := filepath.Ext(pvMessage.File.Filename)
		if ext != "mp3" {
			if ok := utils.ImageValidate(pvMessage.File); !ok {
				return message, errors.New("bad_format")
			}
		}
		imagePath = "assets/uploads/pvMessage/"
		imagePath = utils.FilePathController(imagePath, pvMessage.File.Filename)

	}

	privateMessage := entity.PrivateMessageRoom{
		PrivateID: pvMessage.RoomID,
		Sender:    userId,
		Body:      &pvMessage.Content,
		File:      &imagePath,
	}
	recipientsId, err := c.chatRepository.SendPvMessage(privateMessage)
	if err != nil {
		return message, err
	}
	var sameRoom bool
	for _, item := range recipientsId {
		if userId == item {
			sameRoom = true
		}
	}
	if !sameRoom {
		return message, errors.New("room_id_issue")
	}

	message.PvMessage.Type = "pv_message"
	message.PvMessage.File = imagePath
	message.PvMessage.RoomID = pvMessage.RoomID
	message.PvMessage.Sender = userId
	message.PvMessage.Content = pvMessage.Content
	message.Recipients = recipientsId

	return message, nil

}
func (c chatService) SendGpMessage(gpMessage serializer.MessageRequest, phoneNo string) (serializer.Message, error) {
	userId, err := c.chatRepository.FindByPhone(phoneNo)
	var message serializer.Message
	if err != nil {
		return message, err
	}

	var imagePath string
	if gpMessage.File != nil {
		ext := filepath.Ext(gpMessage.File.Filename)
		if ext != "mp3" {
			if ok := utils.ImageValidate(gpMessage.File); !ok {
				return message, errors.New("bad_format")
			}
		}
		imagePath = "assets/uploads/pvMessage/"
		imagePath = utils.FilePathController(imagePath, gpMessage.File.Filename)

	}
	groupMessage := entity.GroupMessageRoom{
		GroupID: gpMessage.RoomID,
		Sender:  userId,
		Body:    &gpMessage.Content,
		File:    &imagePath,
	}
	recipientsId, err := c.chatRepository.SendGpMessage(groupMessage)
	if err != nil {
		return message, err
	}
	var sameRoom bool
	for _, item := range recipientsId {
		if userId == item {
			sameRoom = true
		}
	}
	if !sameRoom {
		return message, errors.New("room_id_issue")
	}
	message.PvMessage.Type = "gp_message"
	message.PvMessage.File = imagePath
	message.PvMessage.RoomID = gpMessage.RoomID
	message.PvMessage.Sender = userId
	message.PvMessage.Content = gpMessage.Content
	message.Recipients = recipientsId
	return message, nil
}
