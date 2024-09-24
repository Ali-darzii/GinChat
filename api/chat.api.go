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
	SendPvMessage(request *gin.Context)
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
// @Summary connect to websocket and send message
// @Description Sends a chat message to a specific user or group
// @Description it's websocket connection not http post method (swagger doesn't support ws documentation)
// @Description type is either pv_message or group_message
// @Tags chat
// @Accept  json
// @Produce  json
// @Param   message  body  utils.DummyMessage  true  "Message body"
// @Success 200 {object}   serializer.SendPvMessage
// @Failure 500
// @Router /chat/ws [post]
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

// @Summary get all chat rooms
// @Description get all pv and gp chats that user have & need authentication
// @Description avatar --> if it's gp will be gp's avatar and if it's pv it will be user in chat avatar
// @Tags chat
// @Accept  json
// @Produce  json
// @Success 200 {object}   serializer.Room
// @Failure 401
// @Failure 400 {object} utils.ErrorResponse "Token_Expired_Or_Invalid(2) | Object_Not_Found(6)"
// @Router /chat/get-rooms [get]
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

// @Summary make pv chat
// @Description create private chat
// @Description you need to send 1 message too to create private chat
// @Description it has
// @Tags chat
// @Accept  json
// @Produce  json
// @Param   message  body  serializer.MakeNewChatRequest  true  "Message body"
// @Success 201 {object}   serializer.SendPvMessage "you're recipient going to receive the response from ws !"
// @Failure 400 {object}   utils.ErrorResponse "Token_Expired_Or_Invalid(2) | Object_Not_Found(6) | Bad_Format(5) | We_Don't_Know_What_Happened(8)"
// @Router /chat/make-private [post]
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
		request.JSON(http.StatusInternalServerError, utils.SomethingWentWrong)
		return
	}
	request.JSON(http.StatusCreated, message)
	return

}

// @Summary make gp chat
// @Description create group chat
// @Description send data in form-data
// @Description all users of group will receive data of created group by websocket (same as creator)
// @Description so on success creator wil receive nil
// @Tags chat
// @Accept  json
// @Produce  json
// @Param   message  body  utils.DummyMakeGroupChat  true  "Message body"
// @Success 201 {object}   nil
// @Failure 400 {object}   utils.ErrorResponse "Token_Expired_Or_Invalid(2) | Object_Not_Found(6) | Bad_Format(5)"
// @Failure 500 {object}   utils.ErrorResponse "We_Don't_Know_What_Happened(8)"
// @Router /chat/make-group [post]
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
	//NEED UPDATE BECAUSE OF SERIALIZER and METHOD CHANGING !!!!!!
	//websocketHandler.Manager.Broadcast <- message
	request.JSON(http.StatusCreated, nil)
	return
}

// @Summary send pv chat
// @Description send private chat
// @Description *send data in form-data !
// @Description all users will receive data by websocket (same as api creator)
// @Description so on success creator wil receive nil
// @Tags chat
// @Accept  json
// @Produce  json
// @Param   message  body  serializer.PvMessageRequest  true  "Message body"
// @Success 201 {object}   nil
// @Failure 400 {object}   utils.ErrorResponse "Token_Expired_Or_Invalid(2) | Object_Not_Found(6) | Bad_Format(5)"
// @Failure 500 {object}   utils.ErrorResponse "We_Don't_Know_What_Happened(8)"
// @Router /chat/send-pv-message [post]
func (c chatAPI) SendPvMessage(request *gin.Context) {
	var pvMessageRequest serializer.PvMessageRequest
	if err := request.ShouldBindWith(&pvMessageRequest, binding.FormMultipart); err != nil {
		request.JSON(http.StatusBadRequest, utils.BadFormat)
		return
	}
	if ok := pvMessageRequest.PvMessageValidate(); !ok {
		request.JSON(http.StatusBadRequest, utils.BadFormat)
		return
	}

	userPhoneNo, ok := request.Get("phoneNo")
	if !ok {
		request.JSON(http.StatusBadRequest, utils.TokenIsExpiredOrInvalid)
		return
	}
	message, err := c.service.SendPvMessage(pvMessageRequest, userPhoneNo.(string))
	if err != nil {
		if err.Error() == "bad_format" {
			request.JSON(http.StatusBadRequest, utils.BadFormat)
			return
		}
		if err.Error() == "room_id_issue" {
			request.JSON(http.StatusBadRequest, utils.ObjectNotFound)
			return
		}
		request.JSON(http.StatusInternalServerError, utils.SomethingWentWrong)
		return
	}
	if pvMessageRequest.Image != nil {
		pvMessageRequest.Image.Filename = message.PvMessage.Image[25:]
		if err = request.SaveUploadedFile(pvMessageRequest.Image, "assets/uploads/pvMessage/"+pvMessageRequest.Image.Filename); err != nil {
			request.JSON(http.StatusBadRequest, utils.SomethingWentWrong)
			return
		}
	}
	websocketHandler.Manager.Broadcast <- message

	request.JSON(http.StatusOK, nil)
	return
}
