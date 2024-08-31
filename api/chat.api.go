package api

import (
	"GinChat/db"
	"GinChat/entity"
	"GinChat/serializer"
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
	GetAllUsers(*gin.Context)
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

// *only this ChatWs don't use api repo service structure
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
func (c chatAPI) GetAllUsers(request *gin.Context) {
	var paginationRequest serializer.PaginationRequest
	if err := request.ShouldBindQuery(&paginationRequest); err != nil {
		request.JSON(http.StatusBadRequest, utils.BadFormat)
		return
	}
	userPhoneNo, ok := request.Get("phoneNo")
	if !ok {
		request.JSON(http.StatusBadRequest, utils.TokenIsExpiredOrInvalid)
		return
	}
	apiUserPagination, err := c.service.GetAllUsers(paginationRequest, userPhoneNo.(string))
	if err != nil {
		request.JSON(http.StatusInternalServerError, utils.SomethingWentWrong)
		return
	}

	request.JSON(http.StatusOK, apiUserPagination)
	return
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
	if err := request.ShouldBind(&makeChatRequest); err != nil {
		request.JSON(http.StatusBadRequest, utils.BadFormat)
		return
	}
	err := c.service.MakeGroupChat(makeChatRequest, phoneNo.(string))
	if err != nil {
		request.JSON(http.StatusInternalServerError, utils.SomethingWentWrong)
		return
	}
	request.JSON(http.StatusCreated, nil)
	return
}
