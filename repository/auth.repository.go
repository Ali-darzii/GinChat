package repository

import (
	"GinChat/entity"
	"errors"
	"gorm.io/gorm"
)

type AuthRepository interface {
	Register(entity.User) error
	Login(entity.User) (entity.User, error)
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

func (a authRepository) Register(user entity.User) error {
	//todo: check that phone_no doesn't exist | generate SmS token
	if _, err := a.FindByPhone(user.Phone.PhoneNo); err != nil {

		if err.Error() == "not_found" {
			if errs := a.conn.Save(&user); errs != nil {
				return errs.Error
			}
			return nil
		}
		return err
	}
	return errors.New("unique_field")
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

func (a authRepository) Login(user entity.User) (entity.User, error) {
	return user, nil
}
