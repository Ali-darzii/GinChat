package repository

import (
	"GinChat/entity"
	"GinChat/serializer"
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"strconv"
	"time"
)

var ctx = context.Background()

type ChatRepository interface {
	FindByPhone(string) (uint, error)
	GetAllUsers(serializer.PaginationRequest, uint) ([]serializer.UserInRoom, int64, error)
	GetAllRooms(uint) ([]serializer.UserInRoom, error)
	MakePvChat(uint, uint) (entity.PrivateRoom, error)
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
		Select("users.id as user_id, users.name, users.username, COALESCE(private_rooms.id, 0) as room_id, CASE WHEN private_rooms.id IS NULL THEN ? ELSE 0 END as my_user_id", userId).
		Joins("LEFT JOIN pv_users ON users.id = pv_users.user_id").
		Joins("LEFT JOIN private_rooms ON private_rooms.id = pv_users.private_room_id").
		Where("private_rooms.id IN (?) OR private_rooms.id IS NULL", c.postgresConn.Table("pv_users").
			Select("private_room_id").
			Where("user_id = ?", userId)).
		Where("users.id != ?", userId).
		Group("users.id, private_rooms.id").
		Scan(&allUsers)
	//for more security
	for index, item := range allUsers {
		if item.RoomID != 0 {
			allUsers[index].UserID = 0
		}
	}

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
func (c chatRepository) GetAllRooms(userId uint) ([]serializer.UserInRoom, error) {
	var usersInRoom []serializer.UserInRoom
	//if res := c.postgresConn.Table("users").
	//	Select("users.name, users.username, private_rooms.id as room_id").
	//	Joins("JOIN pv_users ON users.id = pv_users.user_id").
	//	Joins("JOIN private_rooms ON private_rooms.id = pv_users.private_room_id").
	//	Where("private_rooms.id IN (?)", c.postgresConn.Table("pv_users").
	//		Select("private_room_id").
	//		Where("user_id = ?", userId)).
	//	Where("users.id != ?", userId).
	//	Scan(&usersInRoom); res.Error != nil {
	//	return []serializer.UserInRoom{}, res.Error
	//}
	//userId = 11
	c.postgresConn.Table("users").
		Select("users.name, users.username, COALESCE(private_rooms.id, 0) as room_id").
		Joins("LEFT JOIN pv_users ON users.id = pv_users.user_id").
		Joins("LEFT JOIN private_rooms ON private_rooms.id = pv_users.private_room_id").
		Where("private_rooms.id IN (?) OR private_rooms.id IS NULL", c.postgresConn.Table("pv_users").
			Select("private_room_id").
			Where("user_id = ?", userId)).
		Where("users.id != ?", userId).
		Group("users.id, private_rooms.id").
		Scan(&usersInRoom)

	return usersInRoom, nil
}
func (c chatRepository) MakePvChat(userId uint, recipientId uint) (entity.PrivateRoom, error) {
	privateRoom := entity.PrivateRoom{
		Users: []entity.User{
			{ID: userId},
			{ID: recipientId},
		},
	}
	if res := c.postgresConn.Create(&privateRoom); res.Error != nil {
		return privateRoom, res.Error
	}
	return privateRoom, nil
}
