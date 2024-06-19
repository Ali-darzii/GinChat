package repository

import (
	"GinChat/entity"
	"errors"
	"gorm.io/gorm"
)

type AuthRepository interface {
	UserSave(entity.User) error
	PhoneSave(entity.Phone) error
	FindByPhone(string) (entity.User, error)
}

type authRepository struct {
	conn *gorm.DB
}

func NewAuthRepository(connection *gorm.DB) AuthRepository {
	return &authRepository{
		conn: connection,
	}
}

func (a authRepository) UserSave(user entity.User) error {
	errs := a.conn.Save(&user)

	if errs.Error != nil {
		return errs.Error
	}

	return nil
}

func (a authRepository) PhoneSave(phone entity.Phone) error {
	res := a.conn.Save(&phone)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (a authRepository) FindByPhone(phoneNo string) (entity.User, error) {
	var phone entity.Phone
	if res := a.conn.Where("phone_no = ?", phoneNo).Take(&phone); res.Error != nil {
		return entity.User{}, errors.New("not_found")
	}
	var user entity.User
	if res := a.conn.Where("id = ?", phone.UserID).Take(&user); res.Error != nil {
		return entity.User{}, res.Error
	}
	user.Phone = phone
	return user, nil
}
