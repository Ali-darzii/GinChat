package repository

import (
	"errors"
	"gorm.io/gorm"
)

type AuthRepository interface {
	Register() (string, error)
}

type authRepository struct {
	conn *gorm.DB
}

func NewAuthRepository(connection *gorm.DB) AuthRepository {
	return &authRepository{
		conn: connection,
	}
}

func (a authRepository) Register() (string, error) {
	return "", errors.New("")
}
