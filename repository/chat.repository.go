package repository

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type ChatRepository interface {
	Chat()
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

func (c chatRepository) Chat() {

}
