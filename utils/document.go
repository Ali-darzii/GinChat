package utils

type RegisterResponse struct {
	Detail   string
	IsSignup bool
}
type DummyMessage struct {
	Type    string `json:"type" binding:"required"`
	RoomID  uint   `json:"room_id" binding:"required"`
	Content string `json:"content,omitempty" binding:"required"`
}
type DummyMakeGroupChat struct {
	Avatar     string `json:"avatar"`
	Name       string `binding:"required" json:"name"`
	Recipients []uint `binding:"required" json:"recipients_id"`
}
type DummyProfileUpdate struct {
	ID       uint   `json:"id"`
	Avatar   string `json:"avatar" form:"avatar"`
	Name     string `binding:"required" json:"name"`
	Username string `json:"username"`
}
type DummyMakeNewChatRequest struct {
	RecipientID uint   `binding:"required,min=1" json:"recipient_id"`
	Content     string `json:"content"`
	File        string `json:"file"`
}
type DummyMessageRequest struct {
	RoomID  uint   `json:"room_id" binding:"required"`
	Content string `json:"content"`
	File    string `json:"file"`
}
