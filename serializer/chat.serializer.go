package serializer

import (
	"mime/multipart"
	"path/filepath"
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
	Image      string `json:"image"`
	Sender     uint   `json:"sender"`
	Recipients []uint `json:"recipients"`
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

type MessageV2 struct {
	Avatar     string `json:"avatar"`
	Recipients []uint `json:"recipients"`
	PvMessage  SendPvMessage
}

type SendPvMessage struct {
	Type      string    `json:"type"`
	RoomID    uint      `json:"room_id" binding:"required"`
	Content   string    `json:"content,omitempty" binding:"required"`
	Sender    uint      `json:"sender"`
	File      string    `json:"file"`
	TimeStamp time.Time `json:"timestamp"`
}

type NewGroupChat struct {
	Avatar  string `json:"avatar" form:"avatar"`
	Type    string `json:"type"`
	RoomID  uint   `json:"room_id"`
	Members []uint `json:"members"`
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
	RecipientID uint                  `binding:"required,min=1" json:"recipient_id" form:"recipient_id"`
	Content     string                `json:"content" form:"content"`
	File        *multipart.FileHeader `json:"file" form:"file"`
}

func (m *MakeNewChatRequest) PvMessageValidate() bool {
	ext := filepath.Ext(m.File.Filename)
	if ext == "mp3" && m.Content != "" {
		return false
	}
	if m.Content == "" && m.File == nil {
		return false
	}
	return true
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

type MessageRequest struct {
	RoomID  uint                  `json:"room_id" binding:"required" form:"room_id"`
	Content string                `json:"content,omitempty" form:"content"`
	File    *multipart.FileHeader `json:"file" form:"image"`
}

// either we should have content or File(image) if there is voice u can't send voice
func (m *MessageRequest) PvMessageValidate() bool {
	ext := filepath.Ext(m.File.Filename)
	if ext == "mp3" && m.Content != "" {
		return false
	}
	if m.Content == "" && m.File == nil {
		return false
	}
	return true
}
