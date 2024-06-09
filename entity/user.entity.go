package entity

import (
	"time"
)

type User struct {
	ID       uint    `gorm:"primary_key:auto_increment" json:"id"`
	Name     string  `gorm:"type:varchar(50);NOT NULL" json:"name"`
	Username *string `gorm:"type:varchar(50);min=5;unique;NULL" json:"username"`
	Phone    Phone   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Phone struct {
	ID      uint   `gorm:"primary_key:auto_increment" json:"id"`
	UserID  uint   `gorm:"uniqueIndex;NOT NULL" json:"user_id"`
	PhoneNo string `gorm:"type:varchar(11);min=11;unique;NOT NULL" json:"phone_no"`
	Token   *int   `gorm:"type:int;min=4,max=4" json:"token"`
	ExpTime *time.Time
}
