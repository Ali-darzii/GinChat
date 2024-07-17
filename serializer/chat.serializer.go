package serializer

type ServerMessage struct {
	Content string `json:"content,omitempty"`
	Status  bool   `json:"status"`
}

type Message struct {
	Type       string `json:"type" binding:"required"`
	RoomID     uint   `json:"room_id" binding:"required"`
	Sender     uint   `json:"sender"`
	Recipients []uint `json:"recipients"`
	Content    string `json:"content,omitempty" binding:"required"`
}

func (c Message) Validate() bool {
	if c.RoomID == 0 {
		return false
	}
	return true
}
