package serializer

type ServerMessage struct {
	Content string `json:"content,omitempty"`
	Status  bool   `json:"status"`
}

type Message struct {
	Type       string `json:"type"`
	RoomID     uint   `json:"room_id" binding:"required"`
	Content    string `json:"content,omitempty" binding:"required"`
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

type PaginationRequest struct {
	Limit  int `form:"limit" json:"limit" binding:"min=2"`
	Offset int `form:"offset" json:"offset"  binding:"min=0"`
}

type UserInRoom struct {
	UserID   uint    `json:"user_id"`
	Name     *string `json:"name"`
	Username *string `json:"username"`
	RoomID   uint    `json:"room_id"`
}

type APIUserPagination struct {
	Count    int64        `json:"count"`
	Previous string       `json:"previous"`
	Next     string       `json:"next"`
	Results  []UserInRoom `json:"results"`
}
type PvUsers struct {
	RoomID uint `json:"private_room_id"`
	UserID uint `json:"user_id"`
}
