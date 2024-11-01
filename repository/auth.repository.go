package repository

import (
	"GinChat/entity"
	"GinChat/utils"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

var ctx = context.Background()

type AuthRepository interface {
	CheckAndMakeOTP(string) error
	NewUserAndMakeOTP(user entity.User) error
	CheckOTP(string, string) error
	UserSave(entity.User) error
	FindByPhone(string) (entity.User, error)
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

func (a authRepository) CheckAndMakeOTP(phoneNo string) error {
	lastCode, err := a.redisConn.Get(ctx, fmt.Sprintf("otp_"+phoneNo)).Result()
	if err != nil {
		return errors.New("something_went_wrong")
	}
	if lastCode != "" {
		return errors.New("too_many_request")
	}
	expTime := utils.GetExpiryTime()
	//handle user test
	if strings.HasPrefix(phoneNo, "0950") {
		token := "5555"
		a.redisConn.Set(ctx, "otp_"+phoneNo, token, expTime)
		return nil
	}
	token := strconv.Itoa(utils.SmsTokenGenerate())
	a.redisConn.Set(ctx, "otp_"+phoneNo, token, expTime)
	go utils.SendSMS(token, phoneNo)
	fmt.Println(token)
	return nil
}
func (a authRepository) NewUserAndMakeOTP(user entity.User) error {
	if res := a.postgresConn.Save(&user); res.Error != nil {
		return res.Error
	}
	go func() {
		user.UserLogins.UserID = user.ID
		a.postgresConn.Save(&user)
		a.redisConn.Del(ctx, "userCount")
	}()
	expTime := utils.GetExpiryTime()
	if strings.HasPrefix(user.PhoneNo, "0950") {
		token := "5555"
		a.redisConn.Set(ctx, "otp_"+user.PhoneNo, token, expTime)
		return nil
	}
	token := strconv.Itoa(utils.SmsTokenGenerate())
	a.redisConn.Set(ctx, "otp_"+user.PhoneNo, token, expTime)
	go utils.SendSMS(token, user.PhoneNo)
	fmt.Println(token)
	return nil
}
func (a authRepository) CheckOTP(PhoneNo string, token string) error {
	lastCode, err := a.redisConn.Get(ctx, "otp_"+PhoneNo).Result()
	if err != nil {
		return err
	}
	if lastCode == "" {
		return errors.New("expired_time")
	}
	if lastCode != token {
		return errors.New("invalid_token")
	}
	return nil
}
func (a authRepository) UserSave(user entity.User) error {
	if errs := a.postgresConn.Save(&user); errs.Error != nil {
		return errs.Error
	}
	a.redisConn.Del(ctx, "userCount")

	return nil
}
func (a authRepository) FindByPhone(phoneNo string) (entity.User, error) {
	var user entity.User
	if res := a.postgresConn.Where("phone_no = ?", phoneNo).Take(&user); res.Error != nil {
		return user, errors.New("not_found")
	}
	return user, nil
}
