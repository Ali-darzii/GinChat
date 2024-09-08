package repository

import (
	"GinChat/entity"
	"GinChat/serializer"
	"errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type AuthRepository interface {
	UserSave(entity.User) error
	PhoneSave(entity.Phone) error
	FindByPhone(string) (entity.User, error)
	NewUserSave(user entity.User) error
	ProfileUpdate(serializer.ProfileUpdateRequest) (serializer.ProfileUpdateRequest, error)
}

type authRepository struct {
	postgresConn *gorm.DB
	redisConn    *redis.Client
}

func NewAuthRepository(postgresConnection *gorm.DB, redisConnection *redis.Client) AuthRepository {
	return &authRepository{
		postgresConn: postgresConnection,
		redisConn:    redisConnection,
	}
}
func (a authRepository) NewUserSave(user entity.User) error {
	user.UserLogins.UserID = user.ID
	if errs := a.postgresConn.Save(&user); errs.Error != nil {
		return errs.Error
	}
	a.redisConn.Del(ctx, "userCount")

	return nil
}
func (a authRepository) UserSave(user entity.User) error {
	if errs := a.postgresConn.Save(&user); errs.Error != nil {
		return errs.Error
	}
	a.redisConn.Del(ctx, "userCount")

	return nil
}
func (a authRepository) PhoneSave(phone entity.Phone) error {
	if res := a.postgresConn.Save(&phone); res.Error != nil {
		return res.Error
	}

	return nil
}
func (a authRepository) FindByPhone(phoneNo string) (entity.User, error) {
	var phone entity.Phone
	if res := a.postgresConn.Where("phone_no = ?", phoneNo).Take(&phone); res.Error != nil {
		return entity.User{}, errors.New("not_found")
	}
	var user entity.User
	if res := a.postgresConn.Where("id = ?", phone.UserID).Take(&user); res.Error != nil {
		return entity.User{}, res.Error
	}
	user.Phone = phone
	return user, nil
}
func (a authRepository) ProfileUpdate(profile serializer.ProfileUpdateRequest) (serializer.ProfileUpdateRequest, error) {
	var userCheck entity.User
	if res := a.postgresConn.Where("username = ?", profile.Username).Take(&userCheck); res.Error == nil {
		return serializer.ProfileUpdateRequest{}, errors.New("username_taken")
	}
	var updatedProfile serializer.UpdatedProfile
	res := a.postgresConn.Table("users").Take(&updatedProfile).
		Where("phone_no = ?", profile.PhoneNo).
		Updates(entity.User{
			Name:     &profile.Name,
			Username: &profile.Username,
			Avatar:   &profile.Avatar.Filename,
		})
	if res != nil {
		return serializer.ProfileUpdateRequest{}, res.Error
	}
	return profile, nil
}
