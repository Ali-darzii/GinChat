package repository

import (
	"GinChat/entity"
	"GinChat/serializer"
	"errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"
)

type UserRepository interface {
	FindByPhone(string) (uint, error)
	GetAllUsers(serializer.PaginationRequest, uint) ([]serializer.UserInRoom, int64, error)
	ProfileUpdate(entity.User) (serializer.UpdatedProfile, error)
	GetUserProfile(entity.User) (serializer.ProfileAPI, error)
}

type userRepository struct {
	postgresConn *gorm.DB
	redisConn    *redis.Client
}

func NewUserRepository(postgres *gorm.DB, redis *redis.Client) UserRepository {
	return &userRepository{
		postgresConn: postgres,
		redisConn:    redis,
	}
}
func (u userRepository) GetAllUsers(paginationRequest serializer.PaginationRequest, userId uint) ([]serializer.UserInRoom, int64, error) {
	var allUsers []serializer.UserInRoom
	u.postgresConn.
		Limit(paginationRequest.Limit).
		Offset(paginationRequest.Offset).
		Table("users").
		Select("users.avatar, users.id as user_id, users.name, users.username, COALESCE(private_rooms.id, 0) as room_id, CASE WHEN private_rooms.id IS NULL THEN ? ELSE 0 END as my_user_id", userId).
		Joins("LEFT JOIN pv_users ON users.id = pv_users.user_id").
		Joins("LEFT JOIN private_rooms ON private_rooms.id = pv_users.private_room_id").
		Where("private_rooms.id IN (?) OR private_rooms.id IS NULL", u.postgresConn.Table("pv_users").
			Select("private_room_id").
			Where("user_id = ?", userId)).
		Where("users.id != ?", userId).
		Group("users.id, private_rooms.id").
		Scan(&allUsers)

	//(User Count), using redis for less query
	redisUserCount, _ := u.redisConn.Get(ctx, "userCount").Result()
	var userCount int64
	if redisUserCount == "" {
		u.postgresConn.Model(&entity.User{}).Count(&userCount)
		u.redisConn.Set(ctx, "userCount", userCount, time.Hour)
	} else {
		userCount, _ = strconv.ParseInt(redisUserCount, 10, 64)
	}
	// exclude self user
	return allUsers, userCount - 1, nil
}

func (u userRepository) FindByPhone(phoneNo string) (uint, error) {
	var phone entity.Phone
	if res := u.postgresConn.Where("phone_no = ?", phoneNo).Take(&phone); res.Error != nil {
		return 0, errors.New("not_found")
	}
	return phone.UserID, nil
}

func (u userRepository) ProfileUpdate(user entity.User) (serializer.UpdatedProfile, error) {
	//unique username check completely
	var userUsernameCheck entity.User
	if res := u.postgresConn.Where("username = ? AND id != ?", user.Username, user.ID).Take(&userUsernameCheck); res.Error == nil {
		return serializer.UpdatedProfile{}, errors.New("username_taken")
	}
	// remove old avatar if exist
	var userImageRemove entity.User
	if res := u.postgresConn.Where("id = ?", user.ID).Take(&userImageRemove); res.Error != nil {
		return serializer.UpdatedProfile{}, res.Error
	}
	if *userImageRemove.Avatar != "" {
		os.Remove(*userImageRemove.Avatar)
	}

	var updatedProfile serializer.UpdatedProfile
	res := u.postgresConn.Save(&user).Find(&updatedProfile)
	if res.Error != nil {
		return serializer.UpdatedProfile{}, res.Error
	}
	return updatedProfile, nil
}

func (u userRepository) GetUserProfile(user entity.User) (serializer.ProfileAPI, error) {
	var userProfile serializer.ProfileAPI
	if res := u.postgresConn.First(&user).Find(userProfile); res.Error != nil {
		return userProfile, res.Error
	}
	return userProfile, nil

}
