package repository

import (
	"GinChat/entity"
	"GinChat/serializer"
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"sort"
)

var ctx = context.Background()

type ChatRepository interface {
	FindByPhone(string) (uint, error)
	GetAllRooms(uint) ([]serializer.Room, error)
	MakePvChat(entity.PrivateRoom, entity.PrivateMessageRoom) (entity.PrivateMessageRoom, error)
	MakeGroupChat(entity.GroupRoom) (entity.GroupRoom, error)
	SendPvMessage(entity.PrivateMessageRoom) ([]uint, error)
	SendGpMessage(entity.GroupMessageRoom) ([]uint, error)
}
type chatRepository struct {
	postgresConn *gorm.DB
	redisConn    *redis.Client
}

func NewChatRepository(postgres *gorm.DB, redisConnection *redis.Client) ChatRepository {
	return &chatRepository{
		postgresConn: postgres,
		redisConn:    redisConnection,
	}
}

func (c chatRepository) FindByPhone(phoneNo string) (uint, error) {
	var phone entity.Phone
	if res := c.postgresConn.Where("phone_no = ?", phoneNo).Take(&phone); res.Error != nil {
		return 0, errors.New("not_found")
	}
	return phone.UserID, nil
}
func (c chatRepository) GetAllRooms(userId uint) ([]serializer.Room, error) {
	var allRooms []serializer.Room
	var pvRooms []serializer.UserInRoom
	c.postgresConn.Table("users").
		Select("users.avatar, users.id as user_id, users.name, users.username, COALESCE(private_rooms.id, 0) as room_id, MAX(private_message_rooms.timestamp) as time_stamp").
		Joins("LEFT JOIN pv_users ON users.id = pv_users.user_id").
		Joins("LEFT JOIN private_rooms ON private_rooms.id = pv_users.private_room_id").
		Joins("LEFT JOIN private_message_rooms ON private_message_rooms.private_id = private_rooms.id").
		Where("private_rooms.id IN (?) OR private_rooms.id IS NULL", c.postgresConn.Table("pv_users").
			Select("private_room_id").
			Where("user_id = ?", userId)).
		Where("users.id != ?", userId).
		Group("users.id, private_rooms.id").
		Scan(&pvRooms)
	for _, room := range pvRooms {
		allRooms = append(allRooms, serializer.Room{
			RoomType: "pv_room",
			Users: []serializer.UserAPI{{
				ID:       room.UserID,
				Username: room.Username,
				Name:     room.Name,
			}},
			Avatar:    room.Avatar,
			RoomID:    room.RoomID,
			TimeStamp: &room.TimeStamp,
		})
	}
	var gpRooms []serializer.UserInGpRoom
	c.postgresConn.Table("users").
		Select("users.avatar, users.id as user_id, users.name, users.username, COALESCE(group_rooms.id, 0) as room_id, group_rooms.name as group_name, MAX(group_message_rooms.timestamp) as time_stamp").
		Joins("LEFT JOIN group_users ON users.id = group_users.user_id").
		Joins("LEFT JOIN group_rooms ON group_rooms.id = group_users.group_room_id").
		Joins("LEFT JOIN group_message_rooms ON group_rooms.id = group_message_rooms.group_id").
		Where("group_rooms.id IN (?) OR group_rooms.id IS NULL", c.postgresConn.Table("group_users").
			Select("group_room_id").
			Where("user_id = ?", userId)).
		Where("users.id != ?", userId).
		Group("users.id, group_rooms.id, group_rooms.name").
		Scan(&gpRooms)

	roomMap := make(map[uint]*serializer.Room)
	for _, userInRoom := range gpRooms {
		if room, exists := roomMap[userInRoom.RoomID]; exists {
			room.Users = append(room.Users, serializer.UserAPI{
				ID:       userInRoom.UserID,
				Name:     userInRoom.Name,
				Username: userInRoom.Username,
			})
		} else {
			roomMap[userInRoom.RoomID] = &serializer.Room{
				RoomType: "gp_room",
				Avatar:   userInRoom.Avatar,
				RoomID:   userInRoom.RoomID,
				Name:     userInRoom.GroupName,
				Users: []serializer.UserAPI{
					{
						ID:       userInRoom.UserID,
						Name:     userInRoom.Name,
						Username: userInRoom.Username,
					},
				},
				TimeStamp: &userInRoom.TimeStamp,
			}
		}
	}
	for _, room := range roomMap {
		allRooms = append(allRooms, *room)
	}
	sort.Slice(allRooms, func(i, j int) bool {
		return allRooms[i].TimeStamp.After(*allRooms[j].TimeStamp)
	})

	return allRooms, nil
}
func (c chatRepository) MakePvChat(privateRoom entity.PrivateRoom, privateMessage entity.PrivateMessageRoom) (entity.PrivateMessageRoom, error) {
	// check users have no pv room together
	var count int64
	res := c.postgresConn.Table("pv_users as u1").
		Select("COUNT(DISTINCT u1.private_room_id)").
		Joins("JOIN pv_users as u2 ON u1.private_room_id = u2.private_room_id").
		Where("u1.user_id = ? AND u2.user_id = ?", privateRoom.Users[0].ID, privateRoom.Users[1].ID).
		Count(&count)
	if res.Error != nil {
		return privateMessage, res.Error
	}
	if count > 0 {
		return privateMessage, errors.New("room_exist")
	}
	if res := c.postgresConn.Save(&privateRoom); res.Error != nil {
		return privateMessage, res.Error
	}
	privateMessage.PrivateID = privateRoom.ID
	if res := c.postgresConn.Save(&privateMessage); res.Error != nil {
		return privateMessage, res.Error
	}

	return privateMessage, nil
}
func (c chatRepository) MakeGroupChat(groupRoom entity.GroupRoom) (entity.GroupRoom, error) {
	if res := c.postgresConn.Save(&groupRoom); res.Error != nil {
		return entity.GroupRoom{}, res.Error
	}
	
	return groupRoom, nil

}
func (c chatRepository) SendPvMessage(pvMessage entity.PrivateMessageRoom) ([]uint, error) {
	var recipientsId []uint
	if res := c.postgresConn.Save(&pvMessage); res.Error != nil {
		return recipientsId, res.Error
	}

	if res := c.postgresConn.Table("pv_users").Select("user_id").Where("private_room_id = ?", pvMessage.PrivateID).Pluck("user_id", &recipientsId); res.Error != nil {
		return recipientsId, res.Error
	}

	return recipientsId, nil
}
func (c chatRepository) SendGpMessage(groupMessage entity.GroupMessageRoom) ([]uint, error) {
	var recipientsId []uint
	if res := c.postgresConn.Save(&groupMessage); res.Error != nil {
		return recipientsId, res.Error
	}

	if res := c.postgresConn.Table("group_users").Select("user_id").Where("group_room_id = ?", groupMessage.GroupID).Pluck("user_id", &recipientsId); res.Error != nil {
		return recipientsId, res.Error
	}
	return recipientsId, nil
}
