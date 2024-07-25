package repository

import (
	"GinChat/entity"
	"errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type AuthRepository interface {
	UserSave(entity.User) error
	PhoneSave(entity.Phone) error
	FindByPhone(string) (entity.User, error)
	NewUserSave(user entity.User) error
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
	res := a.postgresConn.Save(&phone)
	if res.Error != nil {
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
