package entity

import "time"

type PrivateRoom struct {
	ID                 uint               `gorm:"primary_key:auto_increment" json:"id"`
	Users              []User             `gorm:"many2many:pv_users;" json:"pv_users"`
	Timestamp          time.Time          `gorm:"default:current_timestamp" json:"timestamp"`
	PrivateMessageRoom PrivateMessageRoom `gorm:"foreignKey:PrivateID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"private_message"`
}

type PrivateMessageRoom struct {
	ID        uint      `gorm:"primary_key:auto_increment" json:"id"`
	PrivateID uint      `gorm:"NOT NULL" json:"private_id"`
	Sender    uint      `gorm:"NOT NULL" json:"sender"`
	Body      string    `gorm:"type:text;NOT NULL" json:"body"`
	Timestamp time.Time `gorm:"default:current_timestamp" json:"timestamp"`
}

type GroupRoom struct {
	ID               uint             `gorm:"primary_key:auto_increment" json:"id"`
	Avatar           *string          `gorm:"type:varchar(100);NULL" json:"avatar"`
	Name             string           `gorm:"type:varchar(50);NOT NULL" json:"name"`
	Users            []User           `gorm:"many2many:group_users;" json:"group_users"`
	Admins           []User           `gorm:"many2many:group_admins;" json:"admins"`
	Timestamp        time.Time        `gorm:"default:current_timestamp" json:"timestamp"`
	GroupMessageRoom GroupMessageRoom `gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"group_message"`
}

type GroupMessageRoom struct {
	ID        uint      `gorm:"primary_key:auto_increment" json:"id"`
	GroupID   uint      `gorm:"NOT NULL" json:"group_id"`
	Sender    uint      `gorm:"NOT NULL" json:"sender"`
	Body      string    `gorm:"type:text;NOT NULL" json:"body"`
	TimeStamp time.Time `json:"timestamp"`
}
