package repository

import (
	"GinChat/entity"
	"GinChat/serializer"
	"GinChat/websocketHandler"
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"sort"
	"strconv"
	"time"
)

var ctx = context.Background()

type ChatRepository interface {
	FindByPhone(string) (uint, error)
	GetAllUsers(serializer.PaginationRequest, uint) ([]serializer.UserInRoom, int64, error)
	GetAllRooms(uint) ([]serializer.Room, error)
	MakePvChat(serializer.MakeNewChatRequest, uint) (serializer.Message, error)
	MakeGroupChat(entity.GroupRoom) (entity.GroupRoom, error)
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
func (c chatRepository) GetAllUsers(paginationRequest serializer.PaginationRequest, userId uint) ([]serializer.UserInRoom, int64, error) {
	var allUsers []serializer.UserInRoom
	c.postgresConn.
		Limit(paginationRequest.Limit).
		Offset(paginationRequest.Offset).
		Table("users").
		Select("users.avatar, users.id as user_id, users.name, users.username, COALESCE(private_rooms.id, 0) as room_id, CASE WHEN private_rooms.id IS NULL THEN ? ELSE 0 END as my_user_id", userId).
		Joins("LEFT JOIN pv_users ON users.id = pv_users.user_id").
		Joins("LEFT JOIN private_rooms ON private_rooms.id = pv_users.private_room_id").
		Where("private_rooms.id IN (?) OR private_rooms.id IS NULL", c.postgresConn.Table("pv_users").
			Select("private_room_id").
			Where("user_id = ?", userId)).
		Where("users.id != ?", userId).
		Group("users.id, private_rooms.id").
		Scan(&allUsers)

	//(User Count), using redis for less query
	redisUserCount, _ := c.redisConn.Get(ctx, "userCount").Result()
	var userCount int64
	if redisUserCount == "" {
		c.postgresConn.Model(&entity.User{}).Count(&userCount)
		c.redisConn.Set(ctx, "userCount", userCount, time.Hour)
	} else {
		userCount, _ = strconv.ParseInt(redisUserCount, 10, 64)
	}
	// exclude self user
	return allUsers, userCount - 1, nil
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
func (c chatRepository) MakePvChat(makeNewChatRequest serializer.MakeNewChatRequest, userId uint) (serializer.Message, error) {
	var message serializer.Message
	privateRoom := entity.PrivateRoom{
		Users: []entity.User{
			{ID: userId},
			{ID: makeNewChatRequest.RecipientID},
		},
	}
	if res := c.postgresConn.Create(&privateRoom); res.Error != nil {
		return message, res.Error
	}
	message = serializer.Message{
		Type:       "new_pv_message",
		RoomID:     privateRoom.ID,
		Sender:     userId,
		Content:    makeNewChatRequest.Content,
		Recipients: []uint{makeNewChatRequest.RecipientID},
	}
	websocketHandler.Manager.Broadcast <- message

	return message, nil
}
func (c chatRepository) MakeGroupChat(groupRoom entity.GroupRoom) (entity.GroupRoom, error) {
	if res := c.postgresConn.Save(&groupRoom); res.Error != nil {
		return entity.GroupRoom{}, res.Error
	}

	return groupRoom, nil

}
