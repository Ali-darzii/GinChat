package entity

import "time"

type PrivateRoom struct {
	ID                 uint               `gorm:"primary_key:auto_increment" json:"id"`
	Users              []User             `gorm:"many2many:pv_users;" json:"pv_users"`
	Timestamp          time.Time          `gorm:"default:current_timestamp" json:"timestamp"`
	PrivateMessageRoom PrivateMessageRoom `gorm:"foreignKey:PrivateID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"private_message"`
}

type PrivateMessageRoom struct {
	ID        uint   `gorm:"primary_key:auto_increment" json:"id"`
	PrivateID uint   `gorm:"NOT NULL" json:"private_id"`
	Sender    uint   `gorm:"NOT NULL" json:"sender"`
	Body      string `gorm:"type:text;NOT NULL" json:"body"`
}
