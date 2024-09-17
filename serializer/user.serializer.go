package serializer

import (
	"mime/multipart"
	"time"
)

type ProfileUpdateRequest struct {
	ID       uint                  `json:"id" form:"id"`
	Avatar   *multipart.FileHeader `json:"avatar" form:"avatar"`
	Name     string                `binding:"required,name_validator" json:"name" form:"name"`
	Username string                `json:"username" form:"username"`
}
type UserInRoom struct {
	Avatar    string    `json:"avatar" form:"avatar"`
	UserID    uint      `json:"user_id"`
	Name      *string   `json:"name"`
	Username  *string   `json:"username"`
	RoomID    uint      `json:"room_id"`
	TimeStamp time.Time `json:"time_stamp"`
}

type APIUserPagination struct {
	Count    int64        `json:"count"`
	Previous string       `json:"previous"`
	Next     string       `json:"next"`
	Results  []UserInRoom `json:"results"`
}

type UpdatedProfile struct {
	ID       uint   `json:"id" form:"id"`
	Avatar   string `json:"avatar" form:"avatar"`
	Name     string `json:"name" form:"name"`
	Username string `json:"username" form:"username"`
}
type ProfileAPI struct {
	Avatar   *string `json:"avatar"`
	Name     *string `json:"name"`
	Username *string `json:"username"`
}
