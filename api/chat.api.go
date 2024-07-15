package api

import (
	"GinChat/db"
	"GinChat/entity"
	"GinChat/service"
	"GinChat/utils"
	"GinChat/websocketHandler"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"net/http"
)

type ChatAPI interface {
	ChatWs(*gin.Context)
}

type chatAPI struct {
	service service.ChatService
}

func NewChatAPI(service service.ChatService) ChatAPI {
	return &chatAPI{
		service: service,
	}
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024 * 1024 * 1024,
		WriteBufferSize: 1024 * 1024 * 1024,
		// remove CheckOrigin in production
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	postDb *gorm.DB = db.ConnectPostgres()
)

func (c chatAPI) ChatWs(request *gin.Context) {
	go websocketHandler.Manager.Start()
	webSocket, err := upgrader.Upgrade(request.Writer, request.Request, nil)
	phoneNo, exist := request.Get("phoneNo")
	if !exist {
		request.JSON(http.StatusInternalServerError, utils.TokenIsExpiredOrInvalid)
	}
	if err != nil {
		request.JSON(http.StatusInternalServerError, utils.CanNotReachWsConnection)
		return
	}

	if err = webSocket.WriteMessage(websocket.TextMessage, []byte("connected from server")); err != nil {
		request.JSON(http.StatusInternalServerError, utils.CanNotReachWsConnection)
		return
	}
	//userId, err := c.service.WsHandler(phoneNo)
	var phone entity.Phone
	if res := postDb.Where("phone_no = ?", phoneNo).Take(&phone); res.Error != nil {
		request.JSON(http.StatusInternalServerError, utils.SomethingWentWrong)
		return
	}

	client := &websocketHandler.Client{
		Id:     phone.UserID,
		Socket: webSocket,
		Send:   make(chan []byte),
	}
	websocketHandler.Manager.Register <- client
	go client.Read()
	go client.Write()
}
