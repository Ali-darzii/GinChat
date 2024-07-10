package api

import (
	"GinChat/service"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type ChatAPI interface {
	WsHandler(*gin.Context)
}

type chatAPI struct {
	service service.ChatService
}

func NewChatAPI(service service.ChatService) ChatAPI {
	return &chatAPI{
		service: service,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 1024 * 1024,
	WriteBufferSize: 1024 * 1024 * 1024,
	// remove CheckOrigin in production
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//var ctx = context.Background()
//var redisDb *redis.Client = db.ConnectRedis()

// rabbitmq localhost:15672 user & pass -> guest & guest
func (c chatAPI) WsHandler(request *gin.Context) {
	webSocket, err := upgrader.Upgrade(request.Writer, request.Request, nil)

	if err != nil {
		request.JSON(500, gin.H{"error": "can't reach ws connection !"})
		return
	}

	if err := webSocket.WriteMessage(websocket.TextMessage, []byte("connected from server")); err != nil {
		return
	}

}
