package entity

import "time"

type Private struct {
	ID             uint           `gorm:"primary_key:auto_increment" json:"id"`
	Users          []User         `gorm:"many2many:pv_users;" json:"pv_users"`
	Timestamp      time.Time      `gorm:"default:current_timestamp" json:"timestamp"`
	PrivateMessage PrivateMessage `gorm:"foreignKey:PrivateID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"private_message"`
}

type PrivateMessage struct {
	ID        uint   `gorm:"primary_key:auto_increment" json:"id"`
	UserID    uint   `gorm:"uniqueIndex;NOT NULL" json:"user_id"`
	PrivateID uint   `gorm:"uniqueIndex;NOT NULL" json:"private_id"`
	Body      string `gorm:"type:text;NOT NULL" json:"phone_no"`
}
