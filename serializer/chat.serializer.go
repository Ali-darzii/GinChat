package serializer

import (
	"mime/multipart"
	"time"
)

type ServerMessage struct {
	Content string `json:"content,omitempty"`
	RoomID  uint   `json:"room_id"`
	Status  bool   `json:"status"`
}

type Message struct {
	Type       string `json:"type" binding:"required"`
	Avatar     string `json:"avatar"`
	RoomID     uint   `json:"room_id" binding:"required"`
	Content    string `json:"content,omitempty" binding:"required"`
	Sender     uint   `json:"sender"`
	Recipients []uint `json:"recipients"`
}

type SendPvMessage struct {
	Type      string    `json:"type"`
	RoomID    uint      `json:"room_id" binding:"required"`
	Content   string    `json:"content,omitempty" binding:"required"`
	Sender    uint      `json:"sender"`
	TimeStamp time.Time `json:"timestamp"`
}

type NewGroupChat struct {
	Avatar  string `json:"avatar" form:"avatar"`
	Type    string `json:"type"`
	RoomID  uint   `json:"room_id"`
	Members []uint `json:"members"`
}

func (c Message) PrivateMessageValidate() bool {
	if c.RoomID == 0 {
		return false
	}
	return true
}
func (c Message) NewPrivateMessageValidate() bool {
	if c.Recipients[0] == 0 {
		return false
	}

	return true
}

type MakeGroupChatRequest struct {
	Avatar     *multipart.FileHeader `json:"avatar" form:"avatar"`
	Name       string                `binding:"required" json:"name" form:"name"`
	Recipients []uint                `binding:"required" json:"recipients_id" form:"recipients_id"`
}

type PaginationRequest struct {
	Limit  int `form:"limit" json:"limit" binding:"min=2"`
	Offset int `form:"offset" json:"offset"  binding:"min=0"`
}

type MakeNewChatRequest struct {
	RecipientID uint   `binding:"required,min=1" json:"recipient_id"`
	Content     string `binding:"required" json:"content"`
}

type UserInGpRoom struct {
	Avatar    string    `json:"avatar"`
	UserID    uint      `json:"user_id"`
	RoomID    uint      `json:"room_id"`
	GroupName string    `json:"group_name"`
	Name      *string   `json:"name"`
	Username  *string   `json:"username"`
	TimeStamp time.Time `json:"time_stamp"`
}

type UserAPI struct {
	ID       uint    `json:"id"`
	Name     *string `json:"name"`
	Username *string `json:"username"`
}

type Room struct {
	RoomType  string     `json:"room_type"`
	Avatar    string     `json:"avatar"`
	RoomID    uint       `json:"room_id"`
	Name      string     `json:"name"`
	Users     []UserAPI  `json:"users"`
	TimeStamp *time.Time `json:"time_stamp"`
}

type PvMessageRequest struct {
	RoomID  uint                  `json:"room_id" binding:"required" form:"room_id"`
	Content string                `json:"content,omitempty" form:"content"`
	Image   *multipart.FileHeader `json:"image" form:"image"`
	Sender  uint
}

// either we should have content or image or both
func (p *PvMessageRequest) PvMessageValidate() bool {

	if p.Content == "" && p.Image == nil {
		return false
	}
	return true
}
