package api

import (
	"GinChat/db"
	"GinChat/entity"
	"GinChat/serializer"
	"GinChat/service"
	"GinChat/utils"
	"GinChat/websocketHandler"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"net/http"
)

type ChatAPI interface {
	ChatWs(*gin.Context)
	GetAllRooms(ctx *gin.Context)
	MakePvChat(*gin.Context)
	MakeGroupChat(request *gin.Context)
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

// *only this ChatWs don't use api service repo structure
// @Summary Send Chat Message
// @Description Sends a chat message to a specific user
// @Tags chat
// @Accept  json
// @Produce  json
// sss@Param   message  body  MessageRequest  true  "Message body"
// sss@Success 200 {object} MessageResponse
// sss@Failure 400 {object} ErrorResponse
// @Router /chat/send [post]
func (c chatAPI) ChatWs(request *gin.Context) {
	webSocket, err := upgrader.Upgrade(request.Writer, request.Request, nil)
	phoneNo, exist := request.Get("phoneNo")
	if !exist {
		request.JSON(http.StatusInternalServerError, utils.TokenIsExpiredOrInvalid)
	}
	if err != nil {
		request.JSON(http.StatusInternalServerError, utils.CanNotReachWsConnection)
		return
	}

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
func (c chatAPI) GetAllRooms(request *gin.Context) {
	userPhoneNo, ok := request.Get("phoneNo")
	if !ok {
		request.JSON(http.StatusBadRequest, utils.TokenIsExpiredOrInvalid)
		return
	}
	usersInSameRoom, err := c.service.GetAllRooms(userPhoneNo.(string))
	if err != nil {
		request.JSON(http.StatusNoContent, utils.ObjectNotFound)
		return
	}
	request.JSON(http.StatusOK, usersInSameRoom)
}
func (c chatAPI) MakePvChat(request *gin.Context) {
	var makeNewChatRequest serializer.MakeNewChatRequest
	if err := request.ShouldBind(&makeNewChatRequest); err != nil {
		request.JSON(http.StatusBadRequest, utils.BadFormat)
	}

	userPhoneNo, ok := request.Get("phoneNo")
	if !ok {
		request.JSON(http.StatusBadRequest, utils.TokenIsExpiredOrInvalid)
		return
	}

	message, err := c.service.MakePvChat(makeNewChatRequest, userPhoneNo.(string))
	if err != nil {
		request.JSON(http.StatusInternalServerError, err)
		return
	}
	request.JSON(http.StatusCreated, message)
	return

}
func (c chatAPI) MakeGroupChat(request *gin.Context) {
	phoneNo, ok := request.Get("phoneNo")
	if !ok {
		request.JSON(http.StatusBadRequest, utils.BadFormat)
		return
	}
	var makeChatRequest serializer.MakeGroupChatRequest
	if err := request.ShouldBindWith(&makeChatRequest, binding.FormMultipart); err != nil {
		request.JSON(http.StatusBadRequest, utils.BadFormat)
		return
	}
	message, err := c.service.MakeGroupChat(makeChatRequest, phoneNo.(string))
	if err != nil {
		if err.Error() == "bad_format" {
			request.JSON(http.StatusBadRequest, utils.BadFormat)
			return
		}
	}
	if makeChatRequest.Avatar != nil {
		makeChatRequest.Avatar.Filename = message.Avatar[27:]
		if err = request.SaveUploadedFile(makeChatRequest.Avatar, "assets/uploads/groupAvatar/"+makeChatRequest.Avatar.Filename); err != nil {
			request.JSON(http.StatusBadRequest, utils.SomethingWentWrong)
			return
		}
	}
	if err != nil {
		request.JSON(http.StatusInternalServerError, utils.SomethingWentWrong)
		return
	}
	websocketHandler.Manager.Broadcast <- message
	request.JSON(http.StatusCreated, nil)
	return
}
