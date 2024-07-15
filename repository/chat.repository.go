package repository

import (
	"GinChat/entity"
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var ctx = context.Background()

type ChatRepository interface {
	WsHandler(any) (uint, error)
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

func (c chatRepository) WsHandler(phoneNo any) (uint, error) {
	var phone entity.Phone
	if res := c.postgresConn.Where("phone_no = ?", phoneNo).Take(&phone); res.Error != nil {
		return 0, errors.New("not_found")
	}

	return phone.UserID, nil
}
